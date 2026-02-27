package internal

import (
	"context"
	"fmt"
	"net/url"
	"path"
	"strings"
	"time"

	"wb-techschool/L2/2.16/internal/fetcher"
	"wb-techschool/L2/2.16/internal/processor"
	"wb-techschool/L2/2.16/internal/storage"
)

// MirrorConfig описывает параметры мирroring-а.
type MirrorConfig struct {
	RootDir      string        // куда сохраняем локальную копию
	MaxDepth     int           // глубина рекурсии по ссылкам (0 — только начальная страница)
	RequestTTL   time.Duration // общий таймаут на один HTTP-запрос
	MaxPages     int           // защитный лимит числа страниц (0 — без лимита)
	SameHostOnly bool          // ограничиваться ли одним хостом
}

// RunMirror запускает процесс мирroring-а сайта.
func RunMirror(ctx context.Context, startURL string, cfg MirrorConfig) error {
	if cfg.MaxDepth < 0 {
		cfg.MaxDepth = 0
	}
	if cfg.RequestTTL <= 0 {
		cfg.RequestTTL = 20 * time.Second
	}

	start, err := url.Parse(startURL)
	if err != nil {
		return fmt.Errorf("invalid start url: %w", err)
	}

	store := &storage.LocalStorage{Path: cfg.RootDir}
	proc := processor.New(store)
	f := fetcher.New()

	type queueItem struct {
		url   string
		depth int
	}

	visited := make(map[string]bool)
	queue := []queueItem{{url: startURL, depth: 0}}
	pagesProcessed := 0

	normalize := func(u *url.URL) string {
		u2 := *u
		u2.Fragment = ""
		return u2.String()
	}

	isAsset := func(u *url.URL) bool {
		ext := strings.ToLower(path.Ext(u.Path))
		switch ext {
		case ".css", ".js", ".png", ".jpg", ".jpeg", ".gif", ".webp", ".svg", ".ico", ".woff", ".woff2", ".ttf", ".otf":
			return true
		default:
			return false
		}
	}

	for len(queue) > 0 {
		item := queue[0]
		queue = queue[1:]

		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		u, err := url.Parse(item.url)
		if err != nil {
			continue
		}

		if cfg.SameHostOnly && u.Host != start.Host {
			continue
		}

		norm := normalize(u)
		if visited[norm] {
			continue
		}
		visited[norm] = true

		reqCtx, cancel := context.WithTimeout(ctx, cfg.RequestTTL)
		res, err := f.Fetch(reqCtx, item.url)
		cancel()
		if err != nil {
			fmt.Println("fetch error:", err)
			continue
		}

		links, err := proc.Process(res)
		if err != nil {
			fmt.Println("process error:", err)
			continue
		}

		pagesProcessed++

		for _, link := range links {
			linkURL, err := url.Parse(link)
			if err != nil {
				continue
			}
			n := normalize(linkURL)
			if visited[n] {
				continue
			}

			asset := isAsset(linkURL)

			// Ограничение по домену оставляем только для HTML-страниц.
			if cfg.SameHostOnly && !asset && linkURL.Host != start.Host {
				continue
			}

			// Ограничение по глубине применяем только к страницам.
			if !asset && item.depth >= cfg.MaxDepth {
				continue
			}

			if cfg.MaxPages > 0 && pagesProcessed >= cfg.MaxPages && !asset {
				break
			}

			queue = append(queue, queueItem{
				url:   n,
				depth: item.depth + 1,
			})
		}
	}

	return nil
}

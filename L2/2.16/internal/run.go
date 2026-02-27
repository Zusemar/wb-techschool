package processor

import (
	"bytes"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/net/html"

	"wb-techschool/L2/2.16/internal/fetcher"
)

// Processor обрабатывает содержимое страниц.
type Processor struct {
	storage Storage
}

// New создает новый процессор.
func New(storage Storage) *Processor {
	return &Processor{storage: storage}
}

func (p *Processor) Process(res *fetcher.Result) ([]string, error) {
	// Сначала сохраняем файл и получаем его локальный путь
	localPath, err := p.storage.Save(res.FinalURL, res.Body)
	if err != nil {
		return nil, err
	}

	// Если это не HTML — просто выходим
	if res.StatusCode >= 400 || !strings.HasPrefix(res.ContentType, "text/html") {
		return nil, nil
	}

	// Парсим HTML
	doc, err := html.Parse(bytes.NewReader(res.Body))
	if err != nil {
		return nil, err
	}

	baseURL, err := url.Parse(res.FinalURL)
	if err != nil {
		return nil, err
	}

	var links []string

	var crawl func(*html.Node)
	crawl = func(n *html.Node) {
		if n.Type == html.ElementNode {
			var attrKey string

			switch n.Data {
			case "a", "link":
				attrKey = "href"
			case "img", "script":
				attrKey = "src"
			}

			if attrKey != "" {
				for i, attr := range n.Attr {
					if attr.Key == attrKey {
						ref, err := url.Parse(attr.Val)
						if err != nil {
							continue
						}

						abs := baseURL.ResolveReference(ref)
						if abs.Scheme != "http" && abs.Scheme != "https" {
							continue
						}

						links = append(links, abs.String())

						// вычисляем локальный путь для target
						targetLocal, err := p.storage.BuildLocalPath(abs.String())
						if err != nil {
							continue
						}

						rel, err := filepath.Rel(filepath.Dir(localPath), targetLocal)
						if err != nil {
							continue
						}

						n.Attr[i].Val = filepath.ToSlash(rel)
					}
				}
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			crawl(c)
		}
	}

	crawl(doc)

	// Сериализуем HTML обратно
	var buf bytes.Buffer
	if err := html.Render(&buf, doc); err != nil {
		return nil, err
	}

	// Перезаписываем HTML уже с rewrite ссылок
	if err := os.WriteFile(localPath, buf.Bytes(), 0644); err != nil {
		return nil, err
	}

	return links, nil
}

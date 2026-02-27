package processor

import (
	"bytes"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/net/html"

	"wb-techschool/L2/2.16/internal/fetcher"
	"wb-techschool/L2/2.16/internal/storage"
)

type Processor struct {
	storage *storage.LocalStorage
}

func New(s *storage.LocalStorage) *Processor {
	return &Processor{storage: s}
}

// Process сохраняет полученный ресурс и, если это HTML со статусом < 400,
// извлекает новые ссылки для обхода.
func (p *Processor) Process(res *fetcher.Result) ([]string, error) {
	if res.StatusCode >= 400 {
		// даже при ошибочном статусе сохраняем тело как есть
		if _, err := p.storage.Save(res.FinalURL, res.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}

	if !strings.HasPrefix(res.ContentType, "text/html") {
		// не HTML — просто сохраняем как есть, без извлечения ссылок
		if _, err := p.storage.Save(res.FinalURL, res.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}

	localPath, err := p.storage.LocalPath(res.FinalURL)
	if err != nil {
		return nil, err
	}

	links, rewritten, err := p.rewriteHTML(res.FinalURL, localPath, res.Body)
	if err != nil {
		return nil, err
	}

	dir := filepath.Dir(localPath)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, err
	}

	if err := os.WriteFile(localPath, rewritten, 0o644); err != nil {
		return nil, err
	}

	return links, nil
}

// rewriteHTML:
//   - собирает ссылки для дальнейшей загрузки
//   - переписывает href/src на относительные пути до локальных файлов.
func (p *Processor) rewriteHTML(baseURL, localPath string, body []byte) ([]string, []byte, error) {
	doc, err := html.Parse(bytes.NewReader(body))
	if err != nil {
		return nil, nil, err
	}

	base, err := url.Parse(baseURL)
	if err != nil {
		return nil, nil, err
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
					if attr.Key != attrKey {
						continue
					}

					ref, err := url.Parse(attr.Val)
					if err != nil {
						continue
					}

					abs := base.ResolveReference(ref)
					if abs.Scheme != "http" && abs.Scheme != "https" {
						continue
					}

					absStr := abs.String()
					links = append(links, absStr)

					targetLocal, err := p.storage.LocalPath(absStr)
					if err != nil {
						continue
					}

					rel, err := filepath.Rel(filepath.Dir(localPath), targetLocal)
					if err != nil {
						continue
					}

					// приводим к URL-формату (слеши вперёд)
					n.Attr[i].Val = filepath.ToSlash(rel)
				}
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			crawl(c)
		}
	}

	crawl(doc)

	var buf bytes.Buffer
	if err := html.Render(&buf, doc); err != nil {
		return nil, nil, err
	}

	return links, buf.Bytes(), nil
}


package processor

import (
	"bytes"
	"net/url"
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

func (p *Processor) Process(res *fetcher.Result) ([]string, error) {
	// сохраняем файл
	_, err := p.storage.Save(res.FinalURL, res.Body)
	if err != nil {
		return nil, err
	}

	// если не html — ничего не извлекаем
	if res.StatusCode >= 400 {
		return nil, nil
	}

	if !strings.HasPrefix(res.ContentType, "text/html") {
		return nil, nil
	}

	// извлекаем ссылки
	links, err := extractLinks(res.FinalURL, res.Body)
	if err != nil {
		return nil, err
	}

	return links, nil
}

func extractLinks(baseURL string, body []byte) ([]string, error) {
	doc, err := html.Parse(bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	base, err := url.Parse(baseURL)
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
				for _, attr := range n.Attr {
					if attr.Key == attrKey {
						ref, err := url.Parse(attr.Val)
						if err == nil {
							abs := base.ResolveReference(ref)

							// только http/https
							if abs.Scheme == "http" || abs.Scheme == "https" {
								links = append(links, abs.String())
							}
						}
					}
				}
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			crawl(c)
		}
	}

	crawl(doc)

	return links, nil
}

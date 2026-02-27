package storage

import (
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type LocalStorage struct {
	Path string
}

// сохраняет body полученный по url, самостоятельно собирая путь
// возвращает локальный путь сохраненного файла
func (s *LocalStorage) Save(url string, body []byte) (string, error) {
	fullPath, err := buildLocalPath(s.Path, url)
	if err != nil {
		return "", err
	}

	dir := filepath.Dir(fullPath)

	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}

	if err := os.WriteFile(fullPath, body, 0644); err != nil {
		return "", err
	}

	return fullPath, nil
}

// возвращает итоговый путь по которому нужно сохранить файл
func buildLocalPath(rootDir string, rawURL string) (string, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}

	host := u.Host
	p := u.Path

	// Нормализация path
	if p == "" || p == "/" {
		p = "index.html"
	} else {
		// убираем ведущий слеш, чтобы не ломать Join
		if strings.HasPrefix(p, "/") {
			p = p[1:]
		}

		if strings.HasSuffix(p, "/") {
			p = p + "index.html"
		} else if path.Ext(p) == "" {
			p = p + ".html"
		}
	}

	// Собираем итоговый путь (OS-безопасно)
	fullPath := filepath.Join(rootDir, host, filepath.FromSlash(p))

	return fullPath, nil
}

// LocalPath возвращает путь на диске, по которому будет сохранён ресурс.
func (s *LocalStorage) LocalPath(rawURL string) (string, error) {
	return buildLocalPath(s.Path, rawURL)
}

package fetcher

import (
	"context"
	"errors"
	"io"
	"net"
	"net/http"
	"strings"
	"time"
)

// MaxBodyBytes ограничивает размер тела, которое мы читаем в память.
// Увеличен, чтобы можно было скачивать реальные страницы и ресурсы.
const MaxBodyBytes int64 = 20 * 1024 * 1024 // 20 MiB

// Result содержит минимально необходимую информацию о HTTP-ответе.
type Result struct {
	FinalURL    string
	StatusCode  int
	ContentType string
	Body        []byte
}

// Fetcher инкапсулирует общий http.Client.
type Fetcher struct {
	client *http.Client
}

// New возвращает Fetcher с одним настроенным http.Client на все приложение.
func New() *Fetcher {
	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   10 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		MaxIdleConnsPerHost:   10,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		ResponseHeaderTimeout: 10 * time.Second,
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   0, // контролируем временем через context
	}

	return &Fetcher{client: client}
}

// Fetch выполняет HTTP GET с учетом таймаутов, лимита размера и базовых заголовков.
func (f *Fetcher) Fetch(ctx context.Context, rawURL string) (*Result, error) {
	if f == nil || f.client == nil {
		return nil, errors.New("fetcher is not initialized")
	}

	// Разрешаем только http/https.
	if !strings.HasPrefix(rawURL, "http://") && !strings.HasPrefix(rawURL, "https://") {
		return nil, errors.New("only http and https schemes are supported")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "wgetlite/0.1")

	resp, err := f.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	limited := io.LimitReader(resp.Body, MaxBodyBytes)
	body, err := io.ReadAll(limited)
	if err != nil {
		return nil, err
	}

	finalURL := ""
	if resp.Request != nil && resp.Request.URL != nil {
		finalURL = resp.Request.URL.String()
	}

	contentType := resp.Header.Get("Content-Type")
	if idx := strings.Index(contentType, ";"); idx != -1 {
		contentType = strings.TrimSpace(contentType[:idx])
	}

	return &Result{
		FinalURL:    finalURL,
		StatusCode:  resp.StatusCode,
		ContentType: contentType,
		Body:        body,
	}, nil
}


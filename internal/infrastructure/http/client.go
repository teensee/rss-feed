package http

import (
	"context"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"rss-feed/internal/domain/logging"
	"rss-feed/internal/interfaces/rest/middleware"
	"time"
)

var userAgents = []string{
	"Mozilla/5.0 (Windows; U; Windows NT 10.1; Win64; x64; en-US) Gecko/20130401 Firefox/62.6",
	"Mozilla/5.0 (Linux; Android 5.1.1; XT1021 Build/LPC23) AppleWebKit/533.21 (KHTML, like Gecko)  Chrome/48.0.2178.136 Mobile Safari/603.7",
	"Mozilla/5.0 (U; Linux i582 ; en-US) AppleWebKit/534.45 (KHTML, like Gecko) Chrome/47.0.1108.138 Safari/600",
	"Mozilla/5.0 (iPhone; CPU iPhone OS 9_9_4; like Mac OS X) AppleWebKit/601.25 (KHTML, like Gecko)  Chrome/48.0.1616.219 Mobile Safari/600.9",
	"Mozilla/5.0 (Windows NT 10.1; WOW64; en-US) AppleWebKit/603.27 (KHTML, like Gecko) Chrome/53.0.3269.114 Safari/536.6 Edge/16.16842",
	"Mozilla/5.0 (compatible; MSIE 8.0; Windows NT 6.2; Trident/4.0)",
	"Mozilla/5.0 (iPhone; CPU iPhone OS 7_8_2; like Mac OS X) AppleWebKit/534.9 (KHTML, like Gecko)  Chrome/50.0.3983.120 Mobile Safari/600.4",
	"Mozilla/5.0 (Linux; U; Linux x86_64; en-US) Gecko/20100101 Firefox/60.7",
	"Mozilla/5.0 (U; Linux i670 ; en-US) Gecko/20100101 Firefox/46.6",
	"Mozilla/5.0 (Linux; Linux x86_64) Gecko/20130401 Firefox/50.0",
}

type UnexpectedStatusCode struct {
	statusCode int
	body       []byte
}

func (u *UnexpectedStatusCode) Error() string {
	return fmt.Sprintf("unexpected status code: %d (%q)", u.statusCode, u.body)
}

func (u *UnexpectedStatusCode) StatusCode() int {
	return u.statusCode
}

func (u *UnexpectedStatusCode) GetBody() []byte {
	return u.body
}

var _ HttpClient = &Client{}

type HttpClient interface {
	Do(ctx context.Context, req *http.Request) ([]byte, error)
	GET(ctx context.Context, relativePath string, headers map[string]string) ([]byte, error)
	POST(ctx context.Context, relativePath string, body io.Reader, headers map[string]string) ([]byte, error)
}

type Client struct {
	baseUrl *url.URL
	c       *http.Client
	l       logging.Logger
}

func NewClient(baseUrl *url.URL, timeout time.Duration, l logging.Logger) HttpClient {
	return &Client{
		l:       l,
		baseUrl: baseUrl,
		c: &http.Client{
			Timeout: timeout,
		},
	}
}

func (c *Client) do(ctx context.Context, req *http.Request) ([]byte, error) {
	c.l.Debug(ctx, fmt.Sprintf("Do %s", req.URL.String()), "headers", req.Header)

	resp, err := c.c.Do(req)

	if err != nil {
		return nil, err
	}

	defer func() {
		if err = resp.Body.Close(); err != nil {
			log.Printf("failed to close response body: %v", err)
		}
	}()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return []byte{}, err
	}

	switch resp.StatusCode {
	case http.StatusOK:
		return body, nil
	default:
		err = &UnexpectedStatusCode{resp.StatusCode, body}
		c.l.Error(ctx, err.Error())

		return nil, err
	}
}

func (c *Client) Do(ctx context.Context, req *http.Request) ([]byte, error) {
	req.Header.Set("User-Agent", userAgents[rand.Intn(len(userAgents))]) // nolint:gosec // изменить потом на crypto/rand

	if req.Header.Get(middleware.TraceIdHeader.ToString()) == "" {
		if v, ok := ctx.Value(middleware.TraceIdHeader.ToString()).(string); ok {
			req.Header.Set(
				middleware.TraceIdHeader.ToString(),
				v,
			)
		}
	}

	return c.do(ctx, req)
}

func (c *Client) GET(ctx context.Context, relativePath string, headers map[string]string) ([]byte, error) {
	fullURL, err := c.buildURL(ctx, relativePath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodGet, fullURL, http.NoBody)
	if err != nil {
		c.l.Error(ctx, fmt.Sprintf("Failed to create GET request for %s, error: %s", fullURL, err.Error()))
		return nil, fmt.Errorf("failed to create GET request: %w", err)
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	return c.Do(ctx, req)
}

func (c *Client) POST(ctx context.Context, relativePath string, body io.Reader, headers map[string]string) ([]byte, error) {
	fullURL, err := c.buildURL(ctx, relativePath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, fullURL, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create POST request: %w", err)
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	return c.Do(ctx, req)
}

func (c *Client) buildURL(ctx context.Context, relativePath string) (string, error) {
	rel, err := url.Parse(relativePath)
	if err != nil {
		c.l.Error(ctx, fmt.Sprintf("invalid relative path: %s", err))
		return "", fmt.Errorf("invalid relative path: %w", err)
	}

	fullURL := c.baseUrl.ResolveReference(rel)
	resultUrl := fullURL.String()

	c.l.Debug(ctx, fmt.Sprintf("result request url: %s", resultUrl))

	return resultUrl, nil
}

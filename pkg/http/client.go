package http

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"time"
)

type UnexpectedStatusCode struct {
	statusCode int
	body       []byte
}

func (u *UnexpectedStatusCode) Error() string {
	return "unexpected status code"
}

func (u *UnexpectedStatusCode) StatusCode() int {
	return u.statusCode
}

func (u *UnexpectedStatusCode) GetBody() []byte {
	return u.body
}

type Client struct {
	baseUrl url.URL
	c       *http.Client
	l       *slog.Logger
}

func NewClient(baseUrl url.URL, l *slog.Logger) *Client {
	return &Client{
		l:       l,
		baseUrl: baseUrl,
		c: &http.Client{
			Timeout: time.Minute,
		},
	}
}

func (c *Client) Do(ctx context.Context, req *http.Request) ([]byte, error) {
	resp, err := c.c.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return []byte{}, err
	}

	switch resp.StatusCode {
	case http.StatusOK:
		return body, nil
	default:
		return nil, &UnexpectedStatusCode{
			resp.StatusCode,
			body,
		}
	}
}

func (c *Client) GET(ctx context.Context, relativePath string, headers map[string]string) ([]byte, error) {
	fullURL, err := c.buildURL(ctx, relativePath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodGet, fullURL, nil)
	if err != nil {
		c.l.ErrorContext(ctx, fmt.Sprintf("Failed to create GET request for %s, error: %s", fullURL, err.Error()))
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
		c.l.ErrorContext(ctx, fmt.Sprintf("invalid relative path: %s", err))
		return "", fmt.Errorf("invalid relative path: %w", err)
	}

	fullURL := c.baseUrl.ResolveReference(rel)
	resultUrl := fullURL.String()

	c.l.DebugContext(ctx, fmt.Sprintf("result request url: %s", resultUrl))

	return resultUrl, nil
}

package flyone

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"sync"
	"time"

	"github.com/resoul/travel/internal/domain"
)

const (
	baseURL   = "https://api4.flyone.eu"
	siteURL   = "https://flyone.eu/en/"
	originHdr = "https://flyone.eu"
)

var (
	tokenRegex = regexp.MustCompile(`loadCookieToken\('([^']+)'\)`)
)

var _ domain.FlyOneProvider = (*Client)(nil)

// Client handles communication with FlyOne API endpoints.
type Client struct {
	http       *http.Client
	mu         sync.Mutex
	token      string
	tokenUntil time.Time
}

// NewClient creates a new FlyOne API client with transport support.
func NewClient(transport ...http.RoundTripper) *Client {
	var tr http.RoundTripper
	if len(transport) > 0 && transport[0] != nil {
		tr = transport[0]
	}

	return &Client{
		http: &http.Client{
			Transport: tr,
			Timeout:   30 * time.Second,
		},
	}
}

// getToken returns a cached token or extracts a fresh one from flyone.eu HTML.
func (c *Client) getToken(ctx context.Context) (string, error) {
	c.mu.Lock()
	if c.token != "" && time.Now().Before(c.tokenUntil) {
		token := c.token
		c.mu.Unlock()
		return token, nil
	}
	c.mu.Unlock()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, siteURL, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create site request: %w", err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")

	resp, err := c.http.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to fetch FlyOne homepage: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read FlyOne homepage: %w", err)
	}

	matches := tokenRegex.FindSubmatch(body)
	if len(matches) < 2 {
		return "", fmt.Errorf("failed to extract token from FlyOne homepage")
	}

	token := string(matches[1])

	c.mu.Lock()
	c.token = token
	c.tokenUntil = time.Now().Add(25 * time.Minute)
	c.mu.Unlock()

	return token, nil
}

// postJSON performs an authenticated JSON POST request with proper Origin and Referer headers.
func (c *Client) postJSON(ctx context.Context, endpoint string, payload []byte) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Origin", originHdr)
	req.Header.Set("Referer", siteURL)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request to %s failed: %w", endpoint, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("FlyOne API error at %s: status %d", endpoint, resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}

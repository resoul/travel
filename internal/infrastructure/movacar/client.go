package movacar

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/resoul/travel/internal/domain"
)

const (
	baseURL = "https://crowd-api-production-615013621295.europe-west1.run.app"
)

var _ domain.MovacarProvider = (*Client)(nil)

// Client handles communication with Movacar Cloud Run API endpoints.
type Client struct {
	http *http.Client
}

// NewClient creates a new Movacar API client with optional transport caching.
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

// get performs an HTTP GET request with standard headers.
func (c *Client) get(ctx context.Context, endpoint string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request to %s failed: %w", endpoint, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Movacar API error at %s: status %d", endpoint, resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}

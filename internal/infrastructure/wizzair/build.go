package wizzair

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// FetchBuildURL retrieves the current dynamic API base URL for Wizzair.
func (c *Client) FetchBuildURL(ctx context.Context) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://www.wizzair.com/buildnumber", nil)
	if err != nil {
		return "", fmt.Errorf("failed to create buildnumber request: %w", err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "*/*")

	resp, err := c.http.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to fetch buildnumber: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("wizzair buildnumber returned HTTP %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read buildnumber body: %w", err)
	}

	text := strings.TrimSpace(string(body))
	parts := strings.Split(text, " ")
	if len(parts) < 2 {
		return "", fmt.Errorf("invalid buildnumber response format: %q", text)
	}

	return parts[1], nil
}

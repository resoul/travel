package wizzair

import (
	"context"
	"io"
	"net/http"
	"regexp"
	"strings"
)

const defaultBuildURL = "https://be.wizzair.com/29.13.0"

var apiUrlRegex = regexp.MustCompile(`https://be\.wizzair\.com/([0-9\.]+)`)

// FetchBuildURL retrieves the current dynamic API base URL for Wizzair.
func (c *Client) FetchBuildURL(ctx context.Context) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://www.wizzair.com/en-gb", nil)
	if err != nil {
		return defaultBuildURL, nil
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")

	resp, err := c.http.Do(req)
	if err != nil {
		return defaultBuildURL, nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return defaultBuildURL, nil
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return defaultBuildURL, nil
	}

	matches := apiUrlRegex.FindStringSubmatch(string(body))
	if len(matches) > 1 && strings.TrimSpace(matches[1]) != "" {
		return "https://be.wizzair.com/" + strings.TrimSpace(matches[1]), nil
	}

	return defaultBuildURL, nil
}

package cruise

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/resoul/travel/internal/domain"
)

const (
	baseURL   = "https://cruise-web-api-wa-prod.azurewebsites.net"
	tokenURL  = "https://cruise-web-api-wa-prod.azurewebsites.net/api/auth/get-access-token"
	matrixURL = "https://cruise-web-api-wa-prod.azurewebsites.net/api/search/get-search-matrix"
	searchURL = "https://cruise-web-api-wa-prod.azurewebsites.net/api/search/get-search-results"

	apiKey        = "44efd947-fb1a-47f6-bf0f-a9fd7fde6db4"
	cobrandID     = "2059479"
	pin           = "2627"
	partnerID     = "309"
	appDomain     = "cruises.airasiabig.com"
	originHeader  = "https://cruises.airasiabig.com"
	refererHeader = "https://cruises.airasiabig.com/"
	userAgent     = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/128.0.0.0 Safari/537.36"
)

var _ domain.CruiseProvider = (*Client)(nil)

// Client handles communication with AirAsia / Arrivia cruise platform API.
type Client struct {
	http      *http.Client
	mu        sync.Mutex
	token     string
	expiresAt time.Time
}

// NewClient creates a new Cruise API client with transport support.
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

// getToken retrieves a cached JWT token or requests a new one from the authentication endpoint.
func (c *Client) getToken(ctx context.Context) (string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.token != "" && time.Now().Before(c.expiresAt) {
		return c.token, nil
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, tokenURL, bytes.NewBufferString("{}"))
	if err != nil {
		return "", fmt.Errorf("failed to create token request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Origin", originHeader)
	req.Header.Set("Referer", refererHeader)
	req.Header.Set("CruiseWebApiKey", apiKey)
	req.Header.Set("AppDomain", appDomain)

	resp, err := c.http.Do(req)
	if err != nil {
		return "", fmt.Errorf("token request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("token API returned status %d", resp.StatusCode)
	}

	var tokenResp TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return "", fmt.Errorf("failed to decode token response: %w", err)
	}

	if tokenResp.Token == "" {
		return "", fmt.Errorf("token response was empty")
	}

	c.token = tokenResp.Token
	// Set cache duration to 25 minutes
	c.expiresAt = time.Now().Add(25 * time.Minute)

	return c.token, nil
}

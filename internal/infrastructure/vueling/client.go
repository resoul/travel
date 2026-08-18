package vueling

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
	amsBaseURL     = "https://ams.vueling.com"
	airtrfxBaseURL = "https://openair-california.airtrfx.com"
	defaultProfile = "e8ffa738-cb67-4a02-b501-9bfd975a4b65"
	airtrfxAPIKey  = "HeQpRjsFI5xlAaSx2onkjc1HTK0ukqA1IrVvd5fvaMhNtzLTxInTpeYB1MK93pah"
)

var _ domain.VuelingProvider = (*Client)(nil)

// Client handles communication with Vueling AMS and AirTRFX APIs.
type Client struct {
	http           *http.Client
	profileID      string
	apiKey         string
	token          string
	tokenExpiresAt time.Time
	tokenMu        sync.Mutex
}

// NewClient creates a new Vueling API client with token caching and transport support.
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
		profileID: defaultProfile,
		apiKey:    airtrfxAPIKey,
	}
}

// GetToken retrieves an existing cached Bearer token or requests a new one from /asm/v1/Auth.
func (c *Client) GetToken(ctx context.Context) (string, error) {
	c.tokenMu.Lock()
	defer c.tokenMu.Unlock()

	// If token is still valid for at least 30 seconds, return it
	if c.token != "" && time.Now().Add(30*time.Second).Before(c.tokenExpiresAt) {
		return c.token, nil
	}

	url := fmt.Sprintf("%s/asm/v1/Auth", amsBaseURL)
	payload := map[string]string{
		"profileId": c.profileID,
	}

	bodyBytes, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal auth payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return "", fmt.Errorf("failed to create auth request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")

	resp, err := c.http.Do(req)
	if err != nil {
		return "", fmt.Errorf("auth request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("vueling auth API error: status %d", resp.StatusCode)
	}

	var authResp authResponseDTO
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		return "", fmt.Errorf("failed to decode auth response: %w", err)
	}

	c.token = authResp.AccessToken
	expirationSec := authResp.Expiration
	if expirationSec <= 0 {
		expirationSec = 1199
	}
	c.tokenExpiresAt = time.Now().Add(time.Duration(expirationSec) * time.Second)

	return c.token, nil
}

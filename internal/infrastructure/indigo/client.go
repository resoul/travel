package indigo

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
	tokenURL         = "https://api-prod-session-skyplus6e.goindigo.in/v1/token/create"
	fareCalendarURL  = "https://api-prod-booking-skyplus6e.goindigo.in/v1/getfarecalendar"
	fareRadarBaseURL = "https://6ewai.goindigo.in/r10next/web/fare-radar"

	tokenUserKey    = "654a6a3cc4998e498e5c0c8ead072915"
	subscriptionKey = "S9pIpbp4QxCTs98Nzrmy0A=="
	bookingUserKey  = "15faf8ddf1e8354e90e54fa098e8b1a8"

	userAgent     = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/128.0.0.0 Safari/537.36"
	originHeader  = "https://www.goindigo.in"
	refererHeader = "https://www.goindigo.in/"
)

var _ domain.IndiGoProvider = (*Client)(nil)

// Client handles communication with IndiGo API endpoints.
type Client struct {
	http *http.Client

	mu        sync.RWMutex
	token     string
	expiresAt time.Time
}

// NewClient creates a new IndiGo API client.
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

// getToken retrieves a valid JWT session token, requesting a new one if expired.
func (c *Client) getToken(ctx context.Context) (string, error) {
	c.mu.RLock()
	if c.token != "" && time.Now().Add(2*time.Minute).Before(c.expiresAt) {
		tok := c.token
		c.mu.RUnlock()
		return tok, nil
	}
	c.mu.RUnlock()

	c.mu.Lock()
	defer c.mu.Unlock()

	// Double-check after acquiring write lock
	if c.token != "" && time.Now().Add(2*time.Minute).Before(c.expiresAt) {
		return c.token, nil
	}

	reqBody, err := json.Marshal(TokenCreateRequest{
		StrToken:        "",
		SubscriptionKey: subscriptionKey,
	})
	if err != nil {
		return "", fmt.Errorf("failed to marshal token request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, tokenURL, bytes.NewBuffer(reqBody))
	if err != nil {
		return "", fmt.Errorf("failed to create token request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Origin", originHeader)
	req.Header.Set("Referer", refererHeader)
	req.Header.Set("user_key", tokenUserKey)
	req.Header.Set("apikey", tokenUserKey)

	resp, err := c.http.Do(req)
	if err != nil {
		return "", fmt.Errorf("token request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("token request failed with status: %d", resp.StatusCode)
	}

	var tokenResp TokenCreateResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return "", fmt.Errorf("failed to decode token response: %w", err)
	}

	if tokenResp.Data.Token.Token == "" {
		return "", fmt.Errorf("empty token received from session endpoint")
	}

	timeoutMinutes := tokenResp.Data.Token.IdleTimeoutInMinutes
	if timeoutMinutes <= 0 {
		timeoutMinutes = 15
	}

	c.token = tokenResp.Data.Token.Token
	c.expiresAt = time.Now().Add(time.Duration(timeoutMinutes) * time.Minute)

	return c.token, nil
}

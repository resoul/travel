package driiveme

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/resoul/travel/internal/domain"
	"golang.org/x/net/publicsuffix"
)

const (
	defaultBaseURL = "https://www.driiveme.com"
	defaultLocale  = "en-GB"
	userAgent      = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"
)

var _ domain.DriiveMeProvider = (*Client)(nil)

// Client handles communication with DriiveMe website, API endpoints and session authentication.
type Client struct {
	http       *http.Client
	jar        http.CookieJar
	baseURL    string
	locale     string
	mu         sync.RWMutex
	userEmail  string
	isLoggedIn bool
}

// NewClient creates a new DriiveMe API client with cookie jar and optional transport.
func NewClient(transport ...http.RoundTripper) *Client {
	jar, _ := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})

	var tr http.RoundTripper
	if len(transport) > 0 && transport[0] != nil {
		tr = transport[0]
	}

	return &Client{
		http: &http.Client{
			Transport: tr,
			Jar:       jar,
			Timeout:   30 * time.Second,
		},
		jar:     jar,
		baseURL: defaultBaseURL,
		locale:  defaultLocale,
	}
}

// IsAuthenticated returns true if client has successfully logged in.
func (c *Client) IsAuthenticated() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.isLoggedIn
}

// UserEmail returns logged-in user email or empty.
func (c *Client) UserEmail() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.userEmail
}

// Login authenticates with DriiveMe using email and password, setting session cookies.
func (c *Client) Login(ctx context.Context, email, password string) error {
	if email == "" || password == "" {
		return fmt.Errorf("email and password must not be empty")
	}

	// Prepare security cookies required by Symfony CSRF protection and platform check
	baseParsed, err := url.Parse(c.baseURL)
	if err != nil {
		return fmt.Errorf("invalid base url: %w", err)
	}

	c.jar.SetCookies(baseParsed, []*http.Cookie{
		{
			Name:   "csrf-token",
			Value:  "csrf-token",
			Path:   "/",
			Domain: baseParsed.Host,
		},
		{
			Name:   "DEVICE_TOKEN",
			Value:  "MTkyMHgxMDgw", // Base64 1920x1080 screen size
			Path:   "/",
			Domain: baseParsed.Host,
		},
	})

	formData := url.Values{}
	formData.Set("login_form[email]", email)
	formData.Set("login_form[password]", password)
	formData.Set("login_form[_token]", "csrf-token")
	formData.Set("login_form[_remember_me]", "1")

	endpoint := fmt.Sprintf("/%s/login.html", c.locale)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+endpoint, strings.NewReader(formData.Encode()))
	if err != nil {
		return fmt.Errorf("failed to create login request: %w", err)
	}

	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("Origin", c.baseURL)
	req.Header.Set("Referer", fmt.Sprintf("%s/%s/authentication.html", c.baseURL, c.locale))
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Dest", "empty")

	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("login request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read login response: %w", err)
	}

	var loginResp loginResponseDTO
	if err := json.Unmarshal(body, &loginResp); err != nil {
		return fmt.Errorf("unexpected login response (HTTP %d): %s", resp.StatusCode, string(body))
	}

	if resp.StatusCode != http.StatusOK || loginResp.User == "" {
		if loginResp.Message != "" {
			return fmt.Errorf("login failed: %s", loginResp.Message)
		}
		if loginResp.Error != "" {
			return fmt.Errorf("login failed: %s", loginResp.Error)
		}
		return fmt.Errorf("login failed with status %d: %s", resp.StatusCode, string(body))
	}

	c.mu.Lock()
	c.isLoggedIn = true
	c.userEmail = loginResp.User
	c.mu.Unlock()

	return nil
}

// get performs an HTTP GET request to a relative or absolute path.
func (c *Client) get(ctx context.Context, path string) ([]byte, error) {
	fullURL := path
	if strings.HasPrefix(path, "/") {
		fullURL = c.baseURL + path
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,application/json,*/*;q=0.8")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request to %s failed: %w", fullURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("DriiveMe API error at %s: status %d", fullURL, resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}

// post performs an HTTP POST form request.
func (c *Client) post(ctx context.Context, path string, data url.Values) ([]byte, error) {
	fullURL := path
	if strings.HasPrefix(path, "/") {
		fullURL = c.baseURL + path
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fullURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request to %s failed: %w", fullURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("DriiveMe API error at %s: status %d", fullURL, resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}

package ryanair

import (
	"net/http"
	"net/http/cookiejar"
	"time"

	"github.com/resoul/travel/internal/domain"
	"golang.org/x/net/publicsuffix"
)

const baseURL = "https://www.ryanair.com"

var _ domain.RyanairProvider = (*Client)(nil)

// Client handles communication with the Ryanair API.
type Client struct {
	http *http.Client
}

// NewClient creates a new Ryanair API client with cookie support.
func NewClient(transport ...http.RoundTripper) *Client {
	jar, _ := cookiejar.New(&cookiejar.Options{
		PublicSuffixList: publicsuffix.List,
	})

	var tr http.RoundTripper
	if len(transport) > 0 && transport[0] != nil {
		tr = transport[0]
	}

	return &Client{
		http: &http.Client{
			Transport: tr,
			Timeout:   30 * time.Second,
			Jar:       jar,
		},
	}
}

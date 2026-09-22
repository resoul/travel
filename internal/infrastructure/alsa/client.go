package alsa

import (
	"net/http"
	"net/http/cookiejar"
	"time"

	"golang.org/x/net/publicsuffix"
)

const baseURL = "https://www.alsa.com/en"

// Client handles communication with ALSA API (HTTP) and web search (Chromedp).
type Client struct {
	http *http.Client
}

// NewClient creates a new ALSA client with both HTTP and browser capabilities.
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

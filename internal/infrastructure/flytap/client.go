package flytap

import (
	"net/http"
	"time"

	"github.com/resoul/travel/internal/domain"
)

const (
	baseURL         = "https://www.flytap.com"
	originSearchURL = "https://www.flytap.com/api/flight?functionName=originSearch"
	destSearchURL   = "https://www.flytap.com/api/flight?functionName=destinationSearch"
	calendarURL     = "https://www.flytap.com/api/calendar?functionName=calendar"

	userAgent     = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/128.0.0.0 Safari/537.36"
	originHeader  = "https://www.flytap.com"
	refererHeader = "https://www.flytap.com/en-us"
)

var _ domain.FlyTapProvider = (*Client)(nil)

// Client handles communication with TAP Air Portugal API endpoints.
type Client struct {
	http *http.Client
}

// NewClient creates a new TAP Air Portugal API client with transport support.
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

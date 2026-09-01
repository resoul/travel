package vueling

import (
	"context"
	"net/http"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/resoul/travel/internal/domain"
)

const (
	defaultUserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/128.0.0.0 Safari/537.36"
)

var _ domain.VuelingProvider = (*Client)(nil)

// Client handles communication with Vueling via CDN assets and Chromedp headless browser automation.
type Client struct {
	http *http.Client
}

// NewClient creates a new Vueling API client.
func NewClient(transport ...http.RoundTripper) *Client {
	var tr http.RoundTripper
	if len(transport) > 0 && transport[0] != nil {
		tr = transport[0]
	}

	return &Client{
		http: &http.Client{
			Transport: tr,
			Timeout:   15 * time.Second,
		},
	}
}

// executeInBrowser runs Chromedp actions with standard stealth options.
func (c *Client) executeInBrowser(ctx context.Context, actions ...chromedp.Action) error {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.UserAgent(defaultUserAgent),
		chromedp.Flag("headless", true),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-blink-features", "AutomationControlled"),
	)

	allocCtx, cancelAlloc := chromedp.NewExecAllocator(ctx, opts...)
	defer cancelAlloc()

	chromCtx, cancelChrom := chromedp.NewContext(allocCtx)
	defer cancelChrom()

	timeoutCtx, cancelTimeout := context.WithTimeout(chromCtx, 35*time.Second)
	defer cancelTimeout()

	return chromedp.Run(timeoutCtx, actions...)
}

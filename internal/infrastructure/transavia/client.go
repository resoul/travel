package transavia

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/resoul/travel/internal/domain"
)

var _ domain.TransaviaProvider = (*Client)(nil)

// Client handles Transavia low-fare calendar extraction using Chromedp headless automation.
type Client struct {
	http *http.Client
}

// NewClient creates a new Transavia client.
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

// GetFareCalendar extracts daily flight fares for a specific month on Transavia.
func (c *Client) GetFareCalendar(ctx context.Context, origin, destination string, year, month int, adults int) ([]domain.FlightOffer, error) {
	origin = strings.ToUpper(strings.TrimSpace(origin))
	destination = strings.ToUpper(strings.TrimSpace(destination))
	if origin == "" || destination == "" {
		return nil, fmt.Errorf("origin and destination are required")
	}

	if year <= 0 {
		year = time.Now().Year()
	}
	if month <= 0 || month > 12 {
		month = int(time.Now().Month())
	}
	if adults <= 0 {
		adults = 1
	}

	monthStr := fmt.Sprintf("%04d-%02d", year, month)
	targetURL := fmt.Sprintf("https://www.transavia.com/start/api/calendar-fares?dr=%s&ac=%d&cc=0&ic=0&ds=%s&as=%s&lf=Monetary",
		monthStr, adults, origin, destination)

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.UserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"),
		chromedp.Flag("headless", true),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-blink-features", "AutomationControlled"),
	)

	allocCtx, cancelAlloc := chromedp.NewExecAllocator(ctx, opts...)
	defer cancelAlloc()

	chromCtx, cancelChrom := chromedp.NewContext(allocCtx)
	defer cancelChrom()

	timeoutCtx, cancelTimeout := context.WithTimeout(chromCtx, 30*time.Second)
	defer cancelTimeout()

	var bodyText string

	err := chromedp.Run(timeoutCtx,
		chromedp.Navigate(targetURL),
		chromedp.Sleep(4*time.Second),
		chromedp.Evaluate("document.body.innerText", &bodyText),
	)

	if err != nil {
		return nil, fmt.Errorf("failed to fetch Transavia calendar in browser: %w", err)
	}

	var resp CalendarFaresResponseDTO
	if err := json.Unmarshal([]byte(bodyText), &resp); err != nil {
		return nil, fmt.Errorf("failed to parse Transavia calendar JSON: %w", err)
	}

	offers := make([]domain.FlightOffer, 0, len(resp.Data))
	for _, item := range resp.Data {
		var depT *time.Time
		if t, err := time.Parse("2006-01-02", item.Date); err == nil {
			depT = &t
		}

		status := "REGULAR_FARE"
		if item.Type == "lowFare" {
			status = "LOW_FARE 🔥"
		}

		offers = append(offers, domain.FlightOffer{
			TransportType:    domain.TransportTypeFlight,
			Airline:          "Transavia",
			FlightNumber:     fmt.Sprintf("HV %s", status),
			DepartureStation: origin,
			ArrivalStation:   destination,
			DepartureTime:    depT,
			DepartureRaw:     item.Date,
			Price: domain.Price{
				Amount:   item.Price,
				Currency: "EUR",
			},
			IsAvailable: item.Price > 0,
			Status:      status,
		})
	}

	return offers, nil
}

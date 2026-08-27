package eurowings

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

var _ domain.EurowingsProvider = (*Client)(nil)

// Client handles communication with Eurowings via Chromedp headless automation.
type Client struct {
	http *http.Client
}

// NewClient creates a new Eurowings client.
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

// executeInBrowser navigates to a URL with Chromedp and returns the document body text.
func (c *Client) executeInBrowser(ctx context.Context, targetURL string) (string, error) {
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
		return "", fmt.Errorf("failed to fetch %s in browser: %w", targetURL, err)
	}

	return bodyText, nil
}

// GetAirports retrieves all airports in the Eurowings network.
func (c *Client) GetAirports(ctx context.Context) ([]domain.Airport, error) {
	url := "https://apps.eurowings.com/flightsearch/v1/search/airports/list?locale=en-GB&airlineCodes=EW,BCS&showCityClusters=true"

	body, err := c.executeInBrowser(ctx, url)
	if err != nil {
		return nil, err
	}

	var resp AirportsListResponseDTO
	if err := json.Unmarshal([]byte(body), &resp); err != nil {
		return nil, fmt.Errorf("failed to parse Eurowings airports JSON: %w", err)
	}

	countryMap := make(map[string]CountryDTO)
	for _, cnt := range resp.Countries {
		countryMap[cnt.CountryCode] = cnt
	}

	airports := make([]domain.Airport, 0, len(resp.Stations))
	for _, a := range resp.Stations {
		if a.TLC == "" {
			continue
		}

		cnt := countryMap[a.CountryCode]
		countryName := cnt.Name
		if countryName == "" {
			countryName = a.CountryCode
		}

		airports = append(airports, domain.Airport{
			Code: a.TLC,
			Name: a.Name,
			Country: domain.Country{
				Code:     a.CountryCode,
				Name:     countryName,
				Currency: cnt.CurrencyCode,
			},
			Coordinates: domain.Coordinates{
				Latitude:  a.Latitude,
				Longitude: a.Longitude,
			},
		})
	}

	return airports, nil
}

// GetRoutesFromOrigin retrieves all direct destination airport codes available from an origin airport.
func (c *Client) GetRoutesFromOrigin(ctx context.Context, origin string) ([]string, error) {
	origin = strings.ToUpper(strings.TrimSpace(origin))
	if origin == "" {
		return nil, fmt.Errorf("origin airport code is required")
	}

	url := fmt.Sprintf("https://apps.eurowings.com/flightsearch/v1/search/airports?origin=%s&airlineCodes=EW,BCS&locale=en-GB&showCityClusters=true", origin)

	body, err := c.executeInBrowser(ctx, url)
	if err != nil {
		return nil, err
	}

	var resp RoutesResponseDTO
	if err := json.Unmarshal([]byte(body), &resp); err != nil {
		return nil, fmt.Errorf("failed to parse Eurowings routes JSON: %w", err)
	}

	return resp.MatchedTLCs, nil
}

// GetFlightDates retrieves all bookable flight dates for a given route on Eurowings.
func (c *Client) GetFlightDates(ctx context.Context, origin, destination string) ([]string, error) {
	origin = strings.ToUpper(strings.TrimSpace(origin))
	destination = strings.ToUpper(strings.TrimSpace(destination))
	if origin == "" || destination == "" {
		return nil, fmt.Errorf("origin and destination are required")
	}

	url := fmt.Sprintf("https://apps.eurowings.com/flightsearch/v1/search/flight-schedule/dates?origin=%s&destination=%s&locale=en-GB&airlineCodes=EW,BCS", origin, destination)

	body, err := c.executeInBrowser(ctx, url)
	if err != nil {
		return nil, err
	}

	var resp ScheduleDatesResponseDTO
	if err := json.Unmarshal([]byte(body), &resp); err != nil {
		return nil, fmt.Errorf("failed to parse Eurowings flight schedule JSON: %w", err)
	}

	dates := make([]string, 0)
	for _, sec := range resp.Sections {
		for _, m := range sec.BookableMonths {
			if !m.Bookable {
				continue
			}
			for _, d := range m.Dates {
				dates = append(dates, fmt.Sprintf("%04d-%02d-%02d", m.Year, m.Month, d))
			}
		}
	}

	return dates, nil
}

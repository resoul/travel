package tictactrip

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/resoul/travel/internal/domain"
)

const (
	defaultBaseURL   = "https://api.tictactrip.eu"
	defaultUserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"
)

var _ domain.TictactripProvider = (*Client)(nil)

// Client handles communication with the Tictactrip API.
type Client struct {
	http    *http.Client
	baseURL string
}

// NewClient creates a new Tictactrip client.
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
		baseURL: defaultBaseURL,
	}
}

// AutocompleteCities searches for European cities and stations matching a query string.
func (c *Client) AutocompleteCities(ctx context.Context, query string) ([]domain.TictactripCity, error) {
	if query == "" {
		return nil, fmt.Errorf("query cannot be empty")
	}

	reqURL := fmt.Sprintf("%s/cities/autocomplete/?q=%s", c.baseURL, url.QueryEscape(query))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", defaultUserAgent)
	req.Header.Set("Accept", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request to %s failed: %w", reqURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("tictactrip API error: status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var dtos []CityDTO
	if err := json.Unmarshal(body, &dtos); err != nil {
		return nil, fmt.Errorf("failed to parse autocomplete response: %w", err)
	}

	cities := make([]domain.TictactripCity, len(dtos))
	for i, d := range dtos {
		cities[i] = domain.TictactripCity{
			ID:         d.CityID,
			LocalName:  d.LocalName,
			UniqueName: d.UniqueName,
			Latitude:   d.Latitude,
			Longitude:  d.Longitude,
			Score:      d.Score,
			Serviced:   d.Serviced,
		}
	}

	return cities, nil
}

// GetPopularDestinations retrieves popular destinations from a given origin or across Europe.
func (c *Client) GetPopularDestinations(ctx context.Context, fromCity string, limit int) ([]domain.TictactripCity, error) {
	if limit <= 0 {
		limit = 7
	}

	var reqURL string
	if fromCity != "" {
		reqURL = fmt.Sprintf("%s/cities/popular/from/%s/%d", c.baseURL, url.PathEscape(strings.ToLower(fromCity)), limit)
	} else {
		reqURL = fmt.Sprintf("%s/cities/popular/%d", c.baseURL, limit)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", defaultUserAgent)
	req.Header.Set("Accept", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request to %s failed: %w", reqURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("tictactrip API error: status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var dtos []PopularCityDTO
	if err := json.Unmarshal(body, &dtos); err != nil {
		return nil, fmt.Errorf("failed to parse popular cities response: %w", err)
	}

	cities := make([]domain.TictactripCity, len(dtos))
	for i, d := range dtos {
		cities[i] = domain.TictactripCity{
			ID:         d.CityID,
			LocalName:  d.LocalName,
			UniqueName: d.UniqueName,
			Latitude:   d.Latitude,
			Longitude:  d.Longitude,
			NbSearch:   d.NbSearch,
			Serviced:   true,
		}
	}

	return cities, nil
}

// GetMonthlyPriceCalendar retrieves the lowest train and bus prices for every day of a requested month.
func (c *Client) GetMonthlyPriceCalendar(ctx context.Context, originID, destinationID int, month string) ([]domain.TictactripCalendarDay, error) {
	if originID <= 0 || destinationID <= 0 {
		return nil, fmt.Errorf("originID and destinationID must be positive integers")
	}
	if month == "" {
		month = time.Now().Format("2006-01")
	}

	reqURL := fmt.Sprintf("%s/priceCalendar/month?originId=%d&destinationId=%d&month=%s",
		c.baseURL, originID, destinationID, url.QueryEscape(month))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", defaultUserAgent)
	req.Header.Set("Accept", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request to %s failed: %w", reqURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("tictactrip price calendar API error: status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var dtos []PriceCalendarResponseDTO
	if err := json.Unmarshal(body, &dtos); err != nil {
		return nil, fmt.Errorf("failed to parse price calendar response: %w", err)
	}

	days := make([]domain.TictactripCalendarDay, len(dtos))
	for i, d := range dtos {
		day := domain.TictactripCalendarDay{
			Date:    d.Date,
			HasTrip: d.Trip != nil,
		}

		if d.Trip != nil {
			day.TransportType = d.Trip.TransportType
			day.Companies = d.Trip.Companies
			day.DurationMinutes = d.Trip.DurationMinutes
			day.NumberOfStops = d.Trip.NumberOfStops
			day.Price = domain.Price{
				Amount:   float64(d.Trip.PriceCents) / 100.0,
				Currency: "EUR",
			}
			if d.Trip.DepartureUnixUTC > 0 {
				day.DepartureTime = time.Unix(d.Trip.DepartureUnixUTC, 0)
			}
			if d.Trip.ArrivalUnixUTC > 0 {
				day.ArrivalTime = time.Unix(d.Trip.ArrivalUnixUTC, 0)
			}
		}

		days[i] = day
	}

	return days, nil
}

package vueling

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/resoul/travel/internal/domain"
)

// GetSchedule fetches all scheduled flights between origin and destination using Vueling AvailabilityServices.
func (c *Client) GetSchedule(ctx context.Context, origin, destination string, year, month, monthsRange int) ([]domain.FlightOffer, error) {
	origin = strings.ToUpper(origin)
	destination = strings.ToUpper(destination)

	now := time.Now()
	if year <= 0 {
		year = now.Year()
	}
	if month <= 0 {
		month = int(now.Month())
	}
	if monthsRange <= 0 {
		monthsRange = 12
	}

	token, err := c.GetToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get vueling auth token for schedule: %w", err)
	}

	url := fmt.Sprintf("%s/avy/v3/AvailabilityServices/allFlights", amsBaseURL)

	payload := map[string]any{
		"originCode":      origin,
		"destinationCode": destination,
		"year":            year,
		"month":           month,
		"currencyCode":    "EUR",
		"monthsRange":     monthsRange,
		"flightType":      "OW",
	}

	bodyBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal availability payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create availability request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch vueling availability: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("vueling availability API error: status %d", resp.StatusCode)
	}

	var rawFlights []availabilityFlightDTO
	if err := json.NewDecoder(resp.Body).Decode(&rawFlights); err != nil {
		return nil, fmt.Errorf("failed to decode vueling availability: %w", err)
	}

	var offers []domain.FlightOffer
	for _, f := range rawFlights {
		if !f.IsAvailableDay && f.Price <= 0 {
			continue
		}
		flightCopy := f
		offers = append(offers, flightCopy.toDomain())
	}

	return offers, nil
}

// GetDates returns all unique scheduled departure dates between origin and destination.
func (c *Client) GetDates(ctx context.Context, origin, destination string, year, month, monthsRange int) ([]string, error) {
	offers, err := c.GetSchedule(ctx, origin, destination, year, month, monthsRange)
	if err != nil {
		return nil, err
	}

	seen := make(map[string]bool)
	var dates []string

	for _, o := range offers {
		if o.DepartureTime != nil {
			d := o.DepartureTime.Format("2006-01-02")
			if !seen[d] {
				seen[d] = true
				dates = append(dates, d)
			}
		} else if len(o.DepartureRaw) >= 10 {
			d := o.DepartureRaw[:10]
			if !seen[d] {
				seen[d] = true
				dates = append(dates, d)
			}
		}
	}

	sort.Strings(dates)
	return dates, nil
}

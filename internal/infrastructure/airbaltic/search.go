package airbaltic

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/resoul/travel/internal/domain"
)

// fetchFares queries /api/fsf/outbound for fares between origin and destination.
func (c *Client) fetchFares(ctx context.Context, origin, destination, startDate, endDate string) ([]fareOfferDTO, error) {
	origin = strings.ToUpper(origin)
	destination = strings.ToUpper(destination)

	params := url.Values{}
	params.Set("flightMode", "oneway")
	params.Set("origin", origin)
	params.Set("destin", destination)

	if startDate != "" {
		params.Set("startDate", startDate)
	}
	if endDate != "" {
		params.Set("endDate", endDate)
	}

	endpoint := fmt.Sprintf("%s/api/fsf/outbound?%s", baseURL, params.Encode())

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create fares request: %w", err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch airBaltic fares: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("airBaltic fares API error: status %d", resp.StatusCode)
	}

	var response fsfResponseDTO
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode airBaltic fares: %w", err)
	}

	return response.Data, nil
}

// SearchFlights searches for airBaltic flight offers matching criteria.
func (c *Client) SearchFlights(ctx context.Context, criteria domain.FlightSearchCriteria) ([]domain.FlightOffer, error) {
	if criteria.Origin == "" || criteria.Destination == "" {
		return nil, fmt.Errorf("origin and destination are required")
	}

	startDate := criteria.DepartureDate
	endDate := ""

	var targetTime time.Time
	hasTargetTime := false

	if criteria.DepartureDate != "" {
		if t, err := time.Parse("2006-01-02", criteria.DepartureDate); err == nil {
			targetTime = t
			hasTargetTime = true
			if criteria.FlexDaysBefore > 0 {
				startDate = t.AddDate(0, 0, -criteria.FlexDaysBefore).Format("2006-01-02")
			}
			if criteria.FlexDaysAfter > 0 {
				endDate = t.AddDate(0, 0, criteria.FlexDaysAfter).Format("2006-01-02")
			} else {
				// Search within the target month / 30 days window if flex is 0
				endDate = t.AddDate(0, 1, 0).Format("2006-01-02")
			}
		}
	}

	rawFares, err := c.fetchFares(ctx, criteria.Origin, criteria.Destination, startDate, endDate)
	if err != nil {
		return nil, err
	}

	flexBefore := criteria.FlexDaysBefore
	flexAfter := criteria.FlexDaysAfter

	var offers []domain.FlightOffer
	for _, f := range rawFares {
		offer := f.toDomain(criteria.Origin, criteria.Destination)

		if hasTargetTime && offer.DepartureTime != nil {
			offerDate := time.Date(
				offer.DepartureTime.Year(),
				offer.DepartureTime.Month(),
				offer.DepartureTime.Day(),
				0, 0, 0, 0, time.UTC,
			)

			minDate := targetTime.AddDate(0, 0, -flexBefore)
			maxDate := targetTime.AddDate(0, 0, flexAfter)

			if offerDate.Before(minDate) || offerDate.After(maxDate) {
				continue
			}
		} else if criteria.DepartureDate != "" {
			if offer.DepartureRaw != criteria.DepartureDate {
				continue
			}
		}

		offers = append(offers, offer)
	}

	sort.Slice(offers, func(i, j int) bool {
		return offers[i].DepartureRaw < offers[j].DepartureRaw
	})

	return offers, nil
}

// GetDates returns all available scheduled departure dates between origin and destination.
func (c *Client) GetDates(ctx context.Context, origin, destination string) ([]string, error) {
	rawFares, err := c.fetchFares(ctx, origin, destination, "", "")
	if err != nil {
		return nil, err
	}

	seen := make(map[string]bool)
	var dates []string

	for _, f := range rawFares {
		if f.Date != "" && !seen[f.Date] {
			seen[f.Date] = true
			dates = append(dates, f.Date)
		}
	}

	sort.Strings(dates)
	return dates, nil
}

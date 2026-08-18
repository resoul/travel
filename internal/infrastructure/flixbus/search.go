package flixbus

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"time"

	"github.com/resoul/travel/internal/domain"
)

// SearchTrips searches for FlixBus and FlixTrain trips between two cities on a given date.
func (c *Client) SearchTrips(ctx context.Context, criteria domain.FlightSearchCriteria) ([]domain.FlightOffer, error) {
	if criteria.Origin == "" || criteria.Destination == "" {
		return nil, fmt.Errorf("origin and destination (city name or UUID) are required")
	}

	fromUUID, _, err := c.ResolveCityID(ctx, criteria.Origin)
	if err != nil {
		return nil, fmt.Errorf("origin resolution error: %w", err)
	}

	toUUID, _, err := c.ResolveCityID(ctx, criteria.Destination)
	if err != nil {
		return nil, fmt.Errorf("destination resolution error: %w", err)
	}

	depDateStr := criteria.DepartureDate
	if depDateStr == "" {
		depDateStr = time.Now().Format("2006-01-02")
	}

	// Format to DD.MM.YYYY required by FlixBus v4 API
	var flixDate string
	if t, err := time.Parse("2006-01-02", depDateStr); err == nil {
		flixDate = t.Format("02.01.2006")
	} else if t, err := time.Parse("02.01.2006", depDateStr); err == nil {
		flixDate = t.Format("02.01.2006")
	} else {
		flixDate = depDateStr
	}

	adults := criteria.Adults
	if adults <= 0 {
		adults = 1
	}

	currency := "EUR"

	productsJSON := fmt.Sprintf(`{"adult":%d}`, adults)

	endpoint := fmt.Sprintf(
		"%s/search/service/v4/search?from_city_id=%s&to_city_id=%s&departure_date=%s&products=%s&currency=%s&locale=en_US&search_by=cities&include_after_midnight_rides=1&disable_distribusion_trips=0&disable_global_trips=0&disable_trips=%%5B%%5D",
		baseURL,
		fromUUID,
		toUUID,
		url.QueryEscape(flixDate),
		url.QueryEscape(productsJSON),
		currency,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create search request: %w", err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to perform flixbus search: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("flixbus search error: status %d", resp.StatusCode)
	}

	var searchResp searchResponseDTO
	if err := json.NewDecoder(resp.Body).Decode(&searchResp); err != nil {
		return nil, fmt.Errorf("failed to decode flixbus search response: %w", err)
	}

	var offers []domain.FlightOffer
	for _, trip := range searchResp.Trips {
		for _, ride := range trip.Results {
			if ride.Status != "available" && ride.Price.Total <= 0 && ride.Price.TotalWithPlatformFee <= 0 {
				continue
			}
			offers = append(offers, ride.toDomain(currency, searchResp.Cities, searchResp.Stations))
		}
	}

	sort.Slice(offers, func(i, j int) bool {
		if offers[i].DepartureTime != nil && offers[j].DepartureTime != nil {
			return offers[i].DepartureTime.Before(*offers[j].DepartureTime)
		}
		return offers[i].DepartureRaw < offers[j].DepartureRaw
	})

	return offers, nil
}

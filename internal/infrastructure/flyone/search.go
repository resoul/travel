package flyone

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/resoul/travel/internal/domain"
)

// fetchFareCalendar calls /api/routes/cms-route-fare to obtain all available flight offers.
func (c *Client) fetchFareCalendar(ctx context.Context, origin, destination string) ([]domain.FlightOffer, error) {
	origin = strings.ToUpper(origin)
	destination = strings.ToUpper(destination)

	token, err := c.getToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get FlyOne auth token: %w", err)
	}

	reqPayload := cmsRouteFareRequestDTO{
		DepCity: origin,
		ArrCity: destination,
		Token:   token,
	}

	payloadBytes, err := json.Marshal(reqPayload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	url := fmt.Sprintf("%s/api/routes/cms-route-fare", baseURL)
	body, err := c.postJSON(ctx, url, payloadBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch FlyOne fares: %w", err)
	}

	var response cmsRouteFareResponseDTO
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to decode FlyOne fares: %w", err)
	}

	var offers []domain.FlightOffer
	for _, calYear := range response.Calendar {
		currency := calYear.Currency
		for _, calMonth := range calYear.Month {
			for _, dayItem := range calMonth.Days {
				if dayItem.Price > 0 {
					offer := formatFlightOffer(origin, destination, calYear.Year, calMonth.Month, dayItem.Day, dayItem.Price, currency)
					offers = append(offers, offer)
				}
			}
		}
	}

	return offers, nil
}

// SearchFlights searches for FlyOne flight offers matching criteria.
func (c *Client) SearchFlights(ctx context.Context, criteria domain.FlightSearchCriteria) ([]domain.FlightOffer, error) {
	if criteria.Origin == "" || criteria.Destination == "" {
		return nil, fmt.Errorf("origin and destination are required")
	}

	allOffers, err := c.fetchFareCalendar(ctx, criteria.Origin, criteria.Destination)
	if err != nil {
		return nil, err
	}

	var targetTime time.Time
	hasTargetTime := false

	if criteria.DepartureDate != "" {
		if t, err := time.Parse("2006-01-02", criteria.DepartureDate); err == nil {
			targetTime = t
			hasTargetTime = true
		}
	}

	flexBefore := criteria.FlexDaysBefore
	flexAfter := criteria.FlexDaysAfter

	var filtered []domain.FlightOffer
	for _, offer := range allOffers {
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

		filtered = append(filtered, offer)
	}

	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].DepartureRaw < filtered[j].DepartureRaw
	})

	return filtered, nil
}

// GetDates returns all available departure dates between origin and destination.
func (c *Client) GetDates(ctx context.Context, origin, destination string) ([]string, error) {
	offers, err := c.fetchFareCalendar(ctx, origin, destination)
	if err != nil {
		return nil, err
	}

	seen := make(map[string]bool)
	var dates []string

	for _, offer := range offers {
		if offer.DepartureRaw != "" && !seen[offer.DepartureRaw] {
			seen[offer.DepartureRaw] = true
			dates = append(dates, offer.DepartureRaw)
		}
	}

	sort.Strings(dates)
	return dates, nil
}

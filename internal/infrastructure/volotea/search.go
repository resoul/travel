package volotea

import (
	"context"
	"strings"
	"time"

	"github.com/resoul/travel/internal/domain"
)

// SearchFlights searches for Volotea flights matching the requested criteria.
func (c *Client) SearchFlights(ctx context.Context, criteria domain.FlightSearchCriteria) ([]domain.FlightOffer, error) {
	offers, err := c.GetSchedule(ctx, criteria.Origin, criteria.Destination)
	if err != nil {
		return nil, err
	}

	targetDate := criteria.DepartureDate // e.g. 2026-08-25
	flexBefore := criteria.FlexDaysBefore
	flexAfter := criteria.FlexDaysAfter

	var targetTime time.Time
	hasTargetTime := false
	if targetDate != "" {
		if t, err := time.Parse("2006-01-02", targetDate); err == nil {
			targetTime = t
			hasTargetTime = true
		}
	}

	requiredSeats := criteria.Adults + criteria.Teens + criteria.Children
	if requiredSeats <= 0 {
		requiredSeats = 1
	}

	var results []domain.FlightOffer
	for _, offer := range offers {
		if offer.SeatsLeft > 0 && offer.SeatsLeft < requiredSeats {
			continue
		}

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
		} else if targetDate != "" {
			if !strings.HasPrefix(offer.DepartureRaw, targetDate) {
				continue
			}
		}

		results = append(results, offer)
	}

	return results, nil
}

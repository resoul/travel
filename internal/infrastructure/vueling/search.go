package vueling

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/resoul/travel/internal/domain"
)

// SearchFlights searches for Vueling flights matching search criteria.
func (c *Client) SearchFlights(ctx context.Context, criteria domain.FlightSearchCriteria) ([]domain.FlightOffer, error) {
	if criteria.Origin == "" || criteria.Destination == "" {
		return nil, fmt.Errorf("origin and destination are required")
	}

	targetDate := criteria.DepartureDate
	var targetTime time.Time
	hasTargetTime := false
	year := time.Now().Year()
	month := int(time.Now().Month())

	if targetDate != "" {
		if t, err := time.Parse("2006-01-02", targetDate); err == nil {
			targetTime = t
			hasTargetTime = true
			year = t.Year()
			month = int(t.Month())
		}
	}

	monthsRange := 3
	offers, err := c.GetSchedule(ctx, criteria.Origin, criteria.Destination, year, month, monthsRange)
	if err != nil {
		return nil, err
	}

	flexBefore := criteria.FlexDaysBefore
	flexAfter := criteria.FlexDaysAfter

	var results []domain.FlightOffer
	for _, offer := range offers {
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

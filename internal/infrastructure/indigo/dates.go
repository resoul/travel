package indigo

import (
	"context"
	"time"
)

// GetDates returns available scheduled dates for flights between origin and destination.
func (c *Client) GetDates(ctx context.Context, origin, destination string) ([]string, error) {
	startDate := time.Now().Format("2006-01-02")
	endDate := time.Now().AddDate(0, 2, 0).Format("2006-01-02")

	offers, err := c.GetFareCalendar(ctx, origin, destination, startDate, endDate, "INR")
	if err != nil {
		return nil, err
	}

	dates := make([]string, 0, len(offers))
	for _, offer := range offers {
		if offer.IsAvailable && offer.DepartureRaw != "" {
			dates = append(dates, offer.DepartureRaw)
		}
	}

	return dates, nil
}

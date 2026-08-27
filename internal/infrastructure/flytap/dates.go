package flytap

import (
	"context"
	"time"
)

// GetDates returns available scheduled dates for flights between origin and destination.
func (c *Client) GetDates(ctx context.Context, origin, destination string) ([]string, error) {
	now := time.Now()
	dates := make([]string, 0)
	seen := make(map[string]bool)

	// Check current month and next 2 months
	for i := 0; i < 3; i++ {
		t := now.AddDate(0, i, 0)
		offers, err := c.GetCalendar(ctx, origin, destination, t.Year(), int(t.Month()), "PT")
		if err != nil {
			continue
		}

		for _, offer := range offers {
			if offer.IsAvailable && offer.DepartureRaw != "" && !seen[offer.DepartureRaw] {
				seen[offer.DepartureRaw] = true
				dates = append(dates, offer.DepartureRaw)
			}
		}
	}

	return dates, nil
}

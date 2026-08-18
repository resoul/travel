package volotea

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strings"

	"github.com/resoul/travel/internal/domain"
)

// GetSchedule fetches the flight schedule for a route between origin and destination.
func (c *Client) GetSchedule(ctx context.Context, origin, destination string) ([]domain.FlightOffer, error) {
	origin = strings.ToUpper(origin)
	destination = strings.ToUpper(destination)

	pair := fmt.Sprintf("%s-%s", origin, destination)
	url := fmt.Sprintf("%s/dist/schedule/%s_schedule.json?v=1", baseURL, pair)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create schedule request: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch schedule for %s: %w", pair, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("no scheduled flights found for route %s", pair)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("volotea schedule API error for %s: status %d", pair, resp.StatusCode)
	}

	var scheduleData map[string][]flightScheduleDTO
	if err := json.NewDecoder(resp.Body).Decode(&scheduleData); err != nil {
		return nil, fmt.Errorf("failed to decode schedule data: %w", err)
	}

	flightsDTO, exists := scheduleData[pair]
	if !exists {
		// Try looking at any key in the map
		for _, v := range scheduleData {
			flightsDTO = v
			break
		}
	}

	offers := make([]domain.FlightOffer, 0, len(flightsDTO))
	for _, f := range flightsDTO {
		flightCopy := f
		offers = append(offers, flightCopy.toDomain(origin, destination))
	}

	return offers, nil
}

// GetDates returns unique scheduled flight dates (YYYY-MM-DD) between origin and destination.
func (c *Client) GetDates(ctx context.Context, origin, destination string) ([]string, error) {
	offers, err := c.GetSchedule(ctx, origin, destination)
	if err != nil {
		return nil, err
	}

	seen := make(map[string]bool)
	var dates []string

	for _, offer := range offers {
		if offer.DepartureTime != nil {
			d := offer.DepartureTime.Format("2006-01-02")
			if !seen[d] {
				seen[d] = true
				dates = append(dates, d)
			}
		} else if len(offer.DepartureRaw) >= 10 {
			d := offer.DepartureRaw[:10]
			if !seen[d] {
				seen[d] = true
				dates = append(dates, d)
			}
		}
	}

	sort.Strings(dates)
	return dates, nil
}

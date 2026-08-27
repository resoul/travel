package indigo

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/resoul/travel/internal/domain"
)

// GetFareCalendar queries daily lowest flight fares between two airports over a specified date range.
func (c *Client) GetFareCalendar(ctx context.Context, origin, destination, startDate, endDate, currency string) ([]domain.FlightOffer, error) {
	origin = strings.ToUpper(strings.TrimSpace(origin))
	destination = strings.ToUpper(strings.TrimSpace(destination))
	currency = strings.ToUpper(strings.TrimSpace(currency))

	if origin == "" || destination == "" {
		return nil, fmt.Errorf("origin and destination are required")
	}

	if currency == "" {
		currency = "INR"
	}

	if startDate == "" {
		startDate = time.Now().Format("2006-01-02")
	}

	if endDate == "" {
		t, err := time.Parse("2006-01-02", startDate)
		if err != nil {
			t = time.Now()
		}
		endDate = t.AddDate(0, 1, 0).Format("2006-01-02")
	}

	token, err := c.getToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to obtain IndiGo session token: %w", err)
	}

	reqPayload := FareCalendarRequest{
		StartDate:    startDate,
		EndDate:      endDate,
		Origin:       origin,
		Destination:  destination,
		CurrencyCode: currency,
		PromoCode:    "",
		LowestIn:     "M",
	}

	bodyBytes, err := json.Marshal(reqPayload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal fare calendar request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fareCalendarURL, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create fare calendar request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Origin", originHeader)
	req.Header.Set("Referer", refererHeader)
	req.Header.Set("Authorization", token)
	req.Header.Set("user_key", bookingUserKey)
	req.Header.Set("apikey", bookingUserKey)

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fare calendar request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("fare calendar request returned status %d", resp.StatusCode)
	}

	var calResp FareCalendarResponse
	if err := json.NewDecoder(resp.Body).Decode(&calResp); err != nil {
		return nil, fmt.Errorf("failed to decode fare calendar response: %w", err)
	}

	offers := make([]domain.FlightOffer, 0, len(calResp.Data.LowFares))
	for _, item := range calResp.Data.LowFares {
		var depTime *time.Time
		if parsed, err := time.Parse("2006-01-02T15:04:05", item.Date); err == nil {
			depTime = &parsed
		} else if parsed, err := time.Parse("2006-01-02", item.Date); err == nil {
			depTime = &parsed
		}

		dateStr := item.Date
		if depTime != nil {
			dateStr = depTime.Format("2006-01-02")
		}

		if item.Price <= 0 {
			continue
		}

		offers = append(offers, domain.FlightOffer{
			TransportType:    domain.TransportTypeFlight,
			Airline:          "IndiGo",
			DepartureStation: origin,
			ArrivalStation:   destination,
			DepartureTime:    depTime,
			DepartureRaw:     dateStr,
			Price: domain.Price{
				Amount:   item.Price,
				Currency: currency,
			},
			IsAvailable: true,
			Status:      item.Category,
		})
	}

	return offers, nil
}

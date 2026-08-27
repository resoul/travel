package flytap

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/resoul/travel/internal/domain"
)

// GetCalendar queries daily lowest flight fares between two airports for a specific month and year.
func (c *Client) GetCalendar(ctx context.Context, origin, destination string, year, month int, market string) ([]domain.FlightOffer, error) {
	origin = strings.ToUpper(strings.TrimSpace(origin))
	destination = strings.ToUpper(strings.TrimSpace(destination))
	market = strings.ToUpper(strings.TrimSpace(market))

	if origin == "" || destination == "" {
		return nil, fmt.Errorf("origin and destination are required")
	}

	if market == "" {
		market = "PT"
	}

	now := time.Now()
	if year <= 0 {
		year = now.Year()
	}
	if month <= 0 {
		month = int(now.Month())
	}

	reqPayload := CalendarRequest{
		CabinClass:   "E",
		Destination:  destination,
		Market:       market,
		Month:        strconv.Itoa(month),
		Origin:       origin,
		PaxType:      "ADT",
		PayWithMiles: false,
		StarAlliance: false,
		TripType:     "O",
		Year:         strconv.Itoa(year),
	}

	bodyBytes, err := json.Marshal(reqPayload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal calendar request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, calendarURL, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create calendar request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Origin", originHeader)
	req.Header.Set("Referer", refererHeader)

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("calendar request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("calendar request returned status %d", resp.StatusCode)
	}

	var calResp CalendarResponse
	if err := json.NewDecoder(resp.Body).Decode(&calResp); err != nil {
		return nil, fmt.Errorf("failed to decode calendar response: %w", err)
	}

	offers := make([]domain.FlightOffer, 0, len(calResp.Data.BestPriceForDates))
	for _, item := range calResp.Data.BestPriceForDates {
		var depTime *time.Time
		dateStr := item.DepartureDate
		if parsed, err := time.Parse("2006-01-02T15:04:05", item.DepartureDate); err == nil {
			depTime = &parsed
			dateStr = parsed.Format("2006-01-02")
		} else if parsed, err := time.Parse("2006-01-02", item.DepartureDate); err == nil {
			depTime = &parsed
			dateStr = parsed.Format("2006-01-02")
		}

		var priceAmount float64
		if item.BestTotalPrice != nil {
			priceAmount = *item.BestTotalPrice
		}

		currency := item.Currency
		if currency == "" {
			if market == "US" {
				currency = "USD"
			} else {
				currency = "EUR"
			}
		}

		isNoFlight := item.NoFlight != nil && *item.NoFlight
		isAvail := !item.SoldOut && !isNoFlight && priceAmount > 0

		if !isAvail {
			continue
		}

		status := "available"
		if item.SoldOut {
			status = "sold out"
		}

		offers = append(offers, domain.FlightOffer{
			TransportType:    domain.TransportTypeFlight,
			Airline:          "TAP Air Portugal",
			DepartureStation: origin,
			ArrivalStation:   destination,
			DepartureTime:    depTime,
			DepartureRaw:     dateStr,
			Price: domain.Price{
				Amount:   priceAmount,
				Currency: currency,
			},
			IsAvailable: isAvail,
			Status:      status,
		})
	}

	return offers, nil
}

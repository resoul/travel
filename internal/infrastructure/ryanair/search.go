package ryanair

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/resoul/travel/internal/domain"
)

const farfndBaseURL = "https://services-api.ryanair.com/farfnd/3"

// SearchFlights fetches available fares and returns domain FlightOffers.
func (c *Client) SearchFlights(ctx context.Context, criteria domain.FlightSearchCriteria) ([]domain.FlightOffer, error) {
	dateFrom := criteria.DepartureDate
	dateTo := criteria.DepartureDate

	if parsed, err := time.Parse("2006-01-02", criteria.DepartureDate); err == nil {
		if criteria.FlexDaysBefore > 0 {
			dateFrom = parsed.AddDate(0, 0, -criteria.FlexDaysBefore).Format("2006-01-02")
		}
		if criteria.FlexDaysAfter > 0 {
			dateTo = parsed.AddDate(0, 0, criteria.FlexDaysAfter).Format("2006-01-02")
		}
	}

	adults := criteria.Adults
	if adults <= 0 {
		adults = 1
	}

	params := url.Values{}
	params.Set("departureAirportIataCode", criteria.Origin)
	params.Set("arrivalAirportIataCode", criteria.Destination)
	params.Set("outboundDepartureDateFrom", dateFrom)
	params.Set("outboundDepartureDateTo", dateTo)
	params.Set("currency", "EUR")

	apiURL := fmt.Sprintf("%s/oneWayFares?%s", farfndBaseURL, params.Encode())

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create search request: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute search request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ryanair farfnd API error: status %d", resp.StatusCode)
	}

	var raw farfndResponseDTO
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, fmt.Errorf("failed to decode search response: %w", err)
	}

	var offers []domain.FlightOffer
	for _, item := range raw.Fares {
		out := item.Outbound
		var depTime, arrTime *time.Time
		var durationStr string

		if dt, err := time.Parse("2006-01-02T15:04:05", out.DepartureDate); err == nil {
			depTime = &dt
			if at, err := time.Parse("2006-01-02T15:04:05", out.ArrivalDate); err == nil {
				arrTime = &at
				dur := at.Sub(dt)
				hours := int(dur.Hours())
				mins := int(dur.Minutes()) % 60
				durationStr = fmt.Sprintf("%02d:%02d", hours, mins)
			}
		}

		priceVal := out.Price.Value
		if adults > 1 {
			priceVal *= float64(adults)
		}

		offers = append(offers, domain.FlightOffer{
			TransportType:    domain.TransportTypeFlight,
			Airline:          "Ryanair",
			DepartureStation: out.DepartureAirport.IataCode,
			ArrivalStation:   out.ArrivalAirport.IataCode,
			FlightNumber:     out.FlightNumber,
			DepartureTime:    depTime,
			ArrivalTime:      arrTime,
			DepartureRaw:     out.DepartureDate,
			ArrivalRaw:       out.ArrivalDate,
			Duration:         durationStr,
			Price: domain.Price{
				Amount:   priceVal,
				Currency: out.Price.CurrencyCode,
			},
			IsAvailable: true,
			Status:      "available",
		})
	}

	return offers, nil
}

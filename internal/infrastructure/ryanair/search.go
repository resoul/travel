package ryanair

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/resoul/travel/internal/domain"
)

// SearchFlights fetches available fares and returns domain FlightOffers.
func (c *Client) SearchFlights(ctx context.Context, criteria domain.FlightSearchCriteria) ([]domain.FlightOffer, error) {
	flexBefore := criteria.FlexDaysBefore
	flexAfter := criteria.FlexDaysAfter
	if flexBefore == 0 && flexAfter == 0 {
		flexBefore, flexAfter = 2, 2
	}

	adults := criteria.Adults
	if adults <= 0 {
		adults = 1
	}

	params := url.Values{}
	params.Set("ADT", strconv.Itoa(adults))
	params.Set("TEEN", strconv.Itoa(criteria.Teens))
	params.Set("CHD", strconv.Itoa(criteria.Children))
	params.Set("INF", strconv.Itoa(criteria.Infants))
	params.Set("Origin", criteria.Origin)
	params.Set("Destination", criteria.Destination)
	params.Set("DateOut", criteria.DepartureDate)
	params.Set("FlexDaysBeforeOut", strconv.Itoa(flexBefore))
	params.Set("FlexDaysOut", strconv.Itoa(flexAfter))
	params.Set("RoundTrip", strconv.FormatBool(criteria.RoundTrip))
	params.Set("IncludeConnectingFlights", "false")
	params.Set("IncludePrimeFares", "false")
	params.Set("promoCode", "")
	params.Set("ToUs", "AGREED")

	if criteria.RoundTrip && criteria.ReturnDate != "" {
		params.Set("DateIn", criteria.ReturnDate)
		params.Set("FlexDaysBeforeIn", strconv.Itoa(flexBefore))
		params.Set("FlexDaysIn", strconv.Itoa(flexAfter))
	}

	apiURL := fmt.Sprintf(
		"%s/api/booking/v4/en-us/availability?%s",
		baseURL,
		params.Encode(),
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create search request: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute search request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ryanair availability API error: status %d", resp.StatusCode)
	}

	var raw availabilityResponseDTO
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, fmt.Errorf("failed to decode search response: %w", err)
	}

	var offers []domain.FlightOffer
	for _, trip := range raw.Trips {
		for _, date := range trip.Dates {
			for _, flight := range date.Flights {
				price := extractADTPrice(flight.RegularFare.Fares)
				if price <= 0 {
					continue
				}

				dep, arr := "", ""
				if len(flight.Time) >= 2 {
					dep, arr = flight.Time[0], flight.Time[1]
				}

				offers = append(offers, domain.FlightOffer{
					TransportType:    domain.TransportTypeFlight,
					Airline:          "Ryanair",
					DepartureStation: trip.Origin,
					ArrivalStation:   trip.Destination,
					FlightNumber:     flight.FlightNumber,
					DepartureRaw:     dep,
					ArrivalRaw:       arr,
					Duration:         flight.Duration,
					Price: domain.Price{
						Amount:   price,
						Currency: raw.Currency,
					},
					SeatsLeft:   flight.FaresLeft,
					IsAvailable: true,
					Status:      "available",
				})
			}
		}
	}

	return offers, nil
}

func extractADTPrice(fares []fareDTO) float64 {
	for _, f := range fares {
		if f.Type == "ADT" {
			return f.Amount
		}
	}
	if len(fares) > 0 {
		return fares[0].Amount
	}
	return 0
}

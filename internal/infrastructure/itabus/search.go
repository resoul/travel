package itabus

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/resoul/travel/internal/domain"
)

// SearchTrips searches for bus journeys on a specific date.
func (c *Client) SearchTrips(ctx context.Context, criteria domain.FlightSearchCriteria) ([]domain.FlightOffer, error) {
	origin := criteria.Origin
	destination := criteria.Destination
	departDate := criteria.DepartureDate

	if origin == "" || destination == "" || departDate == "" {
		return nil, fmt.Errorf("origin, destination, and departure date are required")
	}

	// Parse and validate date format
	if _, err := time.Parse("2006-01-02", departDate); err != nil {
		return nil, fmt.Errorf("invalid departure date format (expected YYYY-MM-DD): %w", err)
	}

	params := url.Values{}
	params.Set("origin", origin)
	params.Set("destination", destination)
	params.Set("datestart", departDate)
	params.Set("adults", "1")
	params.Set("children", "0")

	apiURL := fmt.Sprintf("%s/Api-Travels?%s", baseURL, params.Encode())

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create travels request: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch travels: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ItaBus travels API error: status %d", resp.StatusCode)
	}

	var raw TravelsResponseDTO
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, fmt.Errorf("failed to decode travels response: %w", err)
	}

	var offers []domain.FlightOffer

	// Extract outbound routes
	if raw.Data.Outbound != nil {
		for _, route := range raw.Data.Outbound.Routes {
			if !route.Available || !route.Rates.Available {
				continue
			}

			depTime, _ := time.Parse("2006-01-02T15:04:05-0700", route.DepartureTimestamp)
			arrTime, _ := time.Parse("2006-01-02T15:04:05-0700", route.ArrivalTimestamp)

			// Get minimum price across all available classes
			minPrice := getMinPrice(route.Rates.ITABUS)
			if minPrice <= 0 {
				continue
			}

			offers = append(offers, domain.FlightOffer{
				TransportType:    domain.TransportTypeBus,
				Airline:          "ItaBus",
				FlightNumber:     route.ServiceName,
				DepartureStation: route.Origin.Code,
				ArrivalStation:   route.Destination.Code,
				DepartureRaw:     route.DepartureTimestamp,
				ArrivalRaw:       route.ArrivalTimestamp,
				DepartureTime:    &depTime,
				ArrivalTime:      &arrTime,
				Duration:         route.TravelDuration,
				Price: domain.Price{
					Amount:   minPrice,
					Currency: "EUR",
				},
				IsAvailable: true,
				Status:      fmt.Sprintf("Transfers: %d, Services: %d", route.TransferNumber, len(route.Services)),
			})
		}
	}

	// Extract return routes if available
	if raw.Data.Return != nil {
		for _, route := range raw.Data.Return.Routes {
			if !route.Available || !route.Rates.Available {
				continue
			}

			depTime, _ := time.Parse("2006-01-02T15:04:05-0700", route.DepartureTimestamp)
			arrTime, _ := time.Parse("2006-01-02T15:04:05-0700", route.ArrivalTimestamp)

			minPrice := getMinPrice(route.Rates.ITABUS)
			if minPrice <= 0 {
				continue
			}

			offers = append(offers, domain.FlightOffer{
				TransportType:    domain.TransportTypeBus,
				Airline:          "ItaBus",
				FlightNumber:     route.ServiceName,
				DepartureStation: route.Origin.Code,
				ArrivalStation:   route.Destination.Code,
				DepartureRaw:     route.DepartureTimestamp,
				ArrivalRaw:       route.ArrivalTimestamp,
				DepartureTime:    &depTime,
				ArrivalTime:      &arrTime,
				Duration:         route.TravelDuration,
				Price: domain.Price{
					Amount:   minPrice,
					Currency: "EUR",
				},
				IsAvailable: true,
				Status:      fmt.Sprintf("Return - Transfers: %d", route.TransferNumber),
			})
		}
	}

	return offers, nil
}

// getMinPrice finds the minimum price across all available rate categories and classes
func getMinPrice(rates map[string]RateCategoryDTO) float64 {
	minPrice := float64(9999)

	for _, category := range rates {
		if category.COMFORT.Price > 0 && category.COMFORT.Price < minPrice {
			minPrice = category.COMFORT.Price
		}
	}

	if minPrice == 9999 {
		return 0
	}
	return minPrice
}

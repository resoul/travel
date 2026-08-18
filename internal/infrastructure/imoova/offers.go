package imoova

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/resoul/travel/internal/domain"
)

const (
	relocationsTableQuery = `
query RelocationsTable {
  relocationsTable {
    id
    departure_city_name
    delivery_city_name
    available_from_date
    available_to_date
    latest_departure_date
    vehicle_name
    vehicle_type
    vehicle_sleeps
    vehicle_seatbelts
    currency
    fuel_amount
    extra_units
  }
}`
)

// GetOffers retrieves and filters campervan relocation offers from imoova GraphQL API.
func (c *Client) GetOffers(ctx context.Context, criteria domain.FlightSearchCriteria) ([]domain.FlightOffer, error) {
	req := graphQLRequest{
		Query:         relocationsTableQuery,
		OperationName: "RelocationsTable",
	}

	body, err := c.query(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch imoova offers: %w", err)
	}

	var resp relocationsTableResponseDTO
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to decode imoova offers: %w", err)
	}

	originFilter := strings.ToLower(strings.TrimSpace(criteria.Origin))
	destFilter := strings.ToLower(strings.TrimSpace(criteria.Destination))
	dateFilter := strings.TrimSpace(criteria.DepartureDate)

	var offers []domain.FlightOffer
	for _, item := range resp.Data.RelocationsTable {
		depCity := item.DepartureCityName
		arrCity := item.DeliveryCityName

		if originFilter != "" && !strings.Contains(strings.ToLower(depCity), originFilter) {
			continue
		}
		if destFilter != "" && !strings.Contains(strings.ToLower(arrCity), destFilter) {
			continue
		}

		if dateFilter != "" {
			fromD := item.AvailableFromDate
			toD := item.AvailableToDate
			if fromD != "" && dateFilter < fromD {
				continue
			}
			if toD != "" && dateFilter > toD {
				continue
			}
		}

		offers = append(offers, item.toDomain())
	}

	sort.Slice(offers, func(i, j int) bool {
		if offers[i].DepartureRaw == offers[j].DepartureRaw {
			return offers[i].DepartureStation < offers[j].DepartureStation
		}
		return offers[i].DepartureRaw < offers[j].DepartureRaw
	})

	return offers, nil
}

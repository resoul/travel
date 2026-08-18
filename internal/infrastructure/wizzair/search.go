package wizzair

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/resoul/travel/internal/domain"
)

// SearchFlights performs a timetable search against Wizzair API and returns domain FlightOffers.
func (c *Client) SearchFlights(ctx context.Context, criteria domain.WizzairSearchCriteria) ([]domain.FlightOffer, error) {
	buildURL, err := c.FetchBuildURL(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get build URL for search: %w", err)
	}

	priceType := criteria.PriceType
	if priceType == "" {
		priceType = "regular"
	}
	adults := criteria.Adults
	if adults <= 0 {
		adults = 1
	}

	payload := flightRequestDTO{
		FlightList: []flightItemDTO{
			{
				DepartureStation: criteria.DepartureStation,
				ArrivalStation:   criteria.ArrivalStation,
				From:             criteria.FromDate,
				To:               criteria.ToDate,
			},
		},
		PriceType:          priceType,
		AdultCount:         adults,
		ChildCount:         criteria.Children,
		InfantCount:        criteria.Infants,
		MacStationsAllowed: false,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal search payload: %w", err)
	}

	apiURL := buildURL + "/Api/search/timetableV2"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create search request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute search: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("wizzair search API error: status %d", resp.StatusCode)
	}

	var result timetableResponseDTO
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode search response: %w", err)
	}

	var offers []domain.FlightOffer
	for _, f := range result.OutboundFlights {
		offers = append(offers, domain.FlightOffer{
			TransportType:    domain.TransportTypeFlight,
			Airline:          "Wizzair",
			DepartureStation: f.DepartureStation,
			ArrivalStation:   f.ArrivalStation,
			DepartureRaw:     f.DepartureDate,
			Price: domain.Price{
				Amount:   f.Price.Amount,
				Currency: f.Price.CurrencyCode,
			},
			IsAvailable: true,
			Status:      "available",
		})
	}

	return offers, nil
}

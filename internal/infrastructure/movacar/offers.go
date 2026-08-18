package movacar

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/resoul/travel/internal/domain"
)

// GetOffers retrieves and filters car/campervan relocation offers from Movacar API.
func (c *Client) GetOffers(ctx context.Context, criteria domain.FlightSearchCriteria) ([]domain.FlightOffer, error) {
	now := time.Now()
	dateFrom := now.Format("2006-01-02")
	dateTo := now.AddDate(0, 2, 0).Format("2006-01-02") // 60-day window

	if criteria.DepartureDate != "" {
		if t, err := time.Parse("2006-01-02", criteria.DepartureDate); err == nil {
			dateFrom = t.Format("2006-01-02")
			dateTo = t.AddDate(0, 1, 0).Format("2006-01-02")
		}
	}

	params := url.Values{}
	params.Set("pickupDateFrom", dateFrom)
	params.Set("pickupDateTo", dateTo)

	endpoint := fmt.Sprintf("%s/v1/offers?%s", baseURL, params.Encode())
	body, err := c.get(ctx, endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch Movacar offers: %w", err)
	}

	var resp offersResponseDTO
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to decode Movacar offers: %w", err)
	}

	// 1. Build lookup tables from included
	stations := make(map[string]string)
	prices := make(map[string]domain.Price)

	for _, item := range resp.Included {
		switch item.Type {
		case "station":
			name := ""
			if city, ok := item.Attributes["city"].(string); ok && city != "" {
				name = city
			} else if alt, ok := item.Attributes["alternative_city"].(string); ok && alt != "" {
				name = alt
			} else if street, ok := item.Attributes["street"].(string); ok && street != "" {
				name = street
			}
			if name != "" {
				stations[item.ID] = name
			}
		case "monetary_amount":
			amount := 1.0
			currency := "EUR"
			if a, ok := item.Attributes["amount"].(float64); ok {
				amount = a / 100.0
			}
			if curr, ok := item.Attributes["currency"].(string); ok && curr != "" {
				currency = curr
			}
			prices[item.ID] = domain.Price{
				Amount:   amount,
				Currency: currency,
			}
		}
	}

	originFilter := strings.ToLower(strings.TrimSpace(criteria.Origin))
	destFilter := strings.ToLower(strings.TrimSpace(criteria.Destination))

	var offers []domain.FlightOffer
	for _, o := range resp.Data {
		origID := o.Relationships.Origin.Data.ID
		destID := o.Relationships.Destination.Data.ID
		priceID := o.Relationships.BasePrice.Data.ID

		origName := stations[origID]
		if origName == "" {
			origName = origID
		}

		destName := stations[destID]
		if destName == "" {
			destName = destID
		}

		// Apply Origin & Destination text filtering if supplied
		if originFilter != "" && !strings.Contains(strings.ToLower(origName), originFilter) {
			continue
		}
		if destFilter != "" && !strings.Contains(strings.ToLower(destName), destFilter) {
			continue
		}

		price, ok := prices[priceID]
		if !ok {
			price = domain.Price{Amount: 1.0, Currency: "EUR"}
		}

		attrs := o.Attributes
		vehicleName := strings.TrimSpace(fmt.Sprintf("%s %s", attrs.Make, attrs.Model))
		if attrs.VehicleCategoryName != "" {
			vehicleName = fmt.Sprintf("%s (%s)", vehicleName, attrs.VehicleCategoryName)
		}
		if vehicleName == "" {
			vehicleName = "Car / Campervan"
		}

		days := attrs.Period / 24
		if days <= 0 {
			days = 1
		}
		durationStr := fmt.Sprintf("%d days (free %d km)", days, attrs.FreeKM)

		depParsed, depFormatted := parseTime(attrs.StartDate)
		arrParsed, arrFormatted := parseTime(attrs.EndDate)

		flightOffer := domain.FlightOffer{
			TransportType:    domain.TransportTypeCar,
			Airline:          "Movacar",
			FlightNumber:     vehicleName,
			DepartureStation: origName,
			ArrivalStation:   destName,
			DepartureTime:    depParsed,
			ArrivalTime:      arrParsed,
			DepartureRaw:     depFormatted,
			ArrivalRaw:       arrFormatted,
			Duration:         durationStr,
			Price:            price,
			IsAvailable:      true,
			Status:           "available",
		}

		offers = append(offers, flightOffer)
	}

	sort.Slice(offers, func(i, j int) bool {
		return offers[i].DepartureRaw < offers[j].DepartureRaw
	})

	return offers, nil
}

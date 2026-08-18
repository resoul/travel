package movacar

import (
	"fmt"
	"time"

	"github.com/resoul/travel/internal/domain"
)

// --- Locations Offers DTO ---

type locationSummaryAttrs struct {
	Name         string  `json:"name"`
	LocationType string  `json:"location_type"`
	OfferCount   int     `json:"offer_count"`
	Latitude     float64 `json:"latitude"`
	Longitude    float64 `json:"longitude"`
	Reference    string  `json:"reference"`
}

type locationSummaryDTO struct {
	Type       string               `json:"type"`
	ID         string               `json:"id"`
	Attributes locationSummaryAttrs `json:"attributes"`
}

type locationsOffersResponseDTO struct {
	Included []locationSummaryDTO `json:"included"`
}

func (l *locationSummaryDTO) toDomain() domain.Airport {
	return domain.Airport{
		Code: l.ID,
		Name: fmt.Sprintf("%s (%d offers)", l.Attributes.Name, l.Attributes.OfferCount),
		City: domain.City{
			Code: l.ID,
			Name: l.Attributes.Name,
		},
		Coordinates: domain.Coordinates{
			Latitude:  l.Attributes.Latitude,
			Longitude: l.Attributes.Longitude,
		},
	}
}

// --- Offers DTO ---

type offerAttrs struct {
	Make                string `json:"make"`
	Model               string `json:"model"`
	VehicleCategoryName string `json:"vehicle_category_name"`
	StartDate           string `json:"start_date"`
	EndDate             string `json:"end_date"`
	Period              int    `json:"period"` // hours
	FreeKM              int    `json:"free_km"`
	CentsPerExtraKM     int    `json:"cents_per_extra_km"`
	ValidUntil          string `json:"valid_until"`
}

type relationshipData struct {
	Type string `json:"type"`
	ID   string `json:"id"`
}

type relationshipItem struct {
	Data relationshipData `json:"data"`
}

type offerRelationships struct {
	Origin      relationshipItem `json:"origin"`
	Destination relationshipItem `json:"destination"`
	BasePrice   relationshipItem `json:"base_price"`
}

type offerDTO struct {
	Type          string             `json:"type"`
	ID            string             `json:"id"`
	Attributes    offerAttrs         `json:"attributes"`
	Relationships offerRelationships `json:"relationships"`
}

type includedItem struct {
	Type       string                 `json:"type"`
	ID         string                 `json:"id"`
	Attributes map[string]interface{} `json:"attributes"`
}

type offersResponseDTO struct {
	Data     []offerDTO     `json:"data"`
	Included []includedItem `json:"included"`
}

func parseTime(raw string) (*time.Time, string) {
	if len(raw) >= 10 {
		t, err := time.Parse("2006-01-02", raw[:10])
		if err == nil {
			return &t, raw[:10]
		}
	}
	return nil, raw
}

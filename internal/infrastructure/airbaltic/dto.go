package airbaltic

import (
	"strconv"
	"time"

	"github.com/resoul/travel/internal/domain"
)

// --- Orig-Dest DTO ---

type airportDTO struct {
	Code      string `json:"code"`
	Type      string `json:"type"`
	Country   string `json:"country"`
	City      string `json:"city"`
	Apt       string `json:"apt"`
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
}

func (a *airportDTO) toDomain() domain.Airport {
	lat, _ := strconv.ParseFloat(a.Latitude, 64)
	lon, _ := strconv.ParseFloat(a.Longitude, 64)

	name := a.Apt
	if name == "" {
		name = a.City
	}

	return domain.Airport{
		Code: a.Code,
		Name: name,
		City: domain.City{
			Code: a.Code,
			Name: a.City,
		},
		Country: domain.Country{
			Name: a.Country,
		},
		Coordinates: domain.Coordinates{
			Latitude:  lat,
			Longitude: lon,
		},
	}
}

type origDataDTO struct {
	BtOrigins map[string]airportDTO `json:"btOrigins"`
}

type destinItemDTO struct {
	BtDest    map[string]airportDTO `json:"btDest"`
	NonBtDest map[string]airportDTO `json:"nonBtDest"`
}

type origDestResponseDTO struct {
	OrigData   origDataDTO              `json:"origData"`
	DestinData map[string]destinItemDTO `json:"destinData"`
}

// --- FSF Outbound / Search DTO ---

type fareOfferDTO struct {
	Price        float64  `json:"price"`
	UpdatedPrice *float64 `json:"updatedPrice"`
	Date         string   `json:"date"` // e.g. "2026-08-21"
	IsDirect     bool     `json:"isDirect"`
}

type fsfResponseDTO struct {
	Success bool           `json:"success"`
	Data    []fareOfferDTO `json:"data"`
}

func (f *fareOfferDTO) toDomain(origin, destination string) domain.FlightOffer {
	price := f.Price
	if f.UpdatedPrice != nil && *f.UpdatedPrice > 0 {
		price = *f.UpdatedPrice
	}

	depParsed, depFormatted := parseAirBalticDate(f.Date)

	transferType := "Direct"
	if !f.IsDirect {
		transferType = "Connecting"
	}

	return domain.FlightOffer{
		TransportType:    domain.TransportTypeFlight,
		Airline:          "airBaltic",
		FlightNumber:     transferType,
		DepartureStation: origin,
		ArrivalStation:   destination,
		DepartureTime:    depParsed,
		DepartureRaw:     depFormatted,
		Price: domain.Price{
			Amount:   price,
			Currency: "EUR",
		},
		IsAvailable: true,
		Status:      "available",
	}
}

func parseAirBalticDate(raw string) (*time.Time, string) {
	// e.g. "2026-08-21"
	if len(raw) >= 10 {
		t, err := time.Parse("2006-01-02", raw[:10])
		if err == nil {
			return &t, t.Format("2006-01-02")
		}
	}
	return nil, raw
}

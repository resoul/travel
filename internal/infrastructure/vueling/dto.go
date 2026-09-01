package vueling

import (
	"fmt"
	"time"

	"github.com/resoul/travel/internal/domain"
)

// --- Auth DTO ---

type authResponseDTO struct {
	TokenType   string `json:"tokenType"`
	AccessToken string `json:"accessToken"`
	Expiration  int    `json:"expiration"`
}

// --- CDN Station & Country DTOs ---

type stationLocationDetailsDTO struct {
	CityCode    string `json:"cityCode"`
	CountryCode string `json:"countryCode"`
}

type stationItemDTO struct {
	StationCode     string                    `json:"stationCode"`
	FullName        string                    `json:"fullName"`
	ShortName       string                    `json:"shortName"`
	IcaoCode        string                    `json:"icaoCode"`
	Allowed         bool                      `json:"allowed"`
	InActive        bool                      `json:"inActive"`
	LocationDetails stationLocationDetailsDTO `json:"locationDetails"`
}

type countryItemDTO struct {
	CountryCode  string `json:"countryCode"`
	CountryCode3 string `json:"countryCode3C"`
	Name         string `json:"name"`
	InActive     bool   `json:"inActive"`
}

// --- AirTRFX Airports DTO ---

type cityNameDTO struct {
	Name string `json:"name"`
}

type countryInfoDTO struct {
	ISOCode string `json:"isoCode"`
	Name    string `json:"name"`
}

type airtrfxAirportDTO struct {
	Name          string         `json:"name"`
	IATACode      string         `json:"iataCode"`
	City          cityNameDTO    `json:"city"`
	Country       countryInfoDTO `json:"country"`
	LocationLabel string         `json:"locationLabel"`
}

func (a *airtrfxAirportDTO) toDomain() domain.Airport {
	return domain.Airport{
		Code: a.IATACode,
		Name: a.Name,
		City: domain.City{
			Name: a.City.Name,
		},
		Country: domain.Country{
			Code: a.Country.ISOCode,
			Name: a.Country.Name,
		},
	}
}

// --- AMS Markets DTO ---

type marketDTO struct {
	FromCode   string `json:"fromCode"`
	ToCode     string `json:"toCode"`
	Connection string `json:"connection"`
}

// --- AMS Availability Flights DTO ---

type availabilityFlightDTO struct {
	DepartureStation   string  `json:"departureStation"`
	ArrivalStation     string  `json:"arrivalStation"`
	DepartureDate      string  `json:"departureDate"` // e.g. "2026-08-18T16:55:00"
	ArrivalDate        string  `json:"arrivalDate"`   // e.g. "2026-08-18T00:00:00"
	FlightID           string  `json:"flightID"`
	Price              float64 `json:"price"`
	Currency           string  `json:"currency"`
	CarrierCode        string  `json:"carrierCode"`
	IsConnectionFlight bool    `json:"isConnectionFlight"`
	IsAvailableDay     bool    `json:"isAvailableDay"`
}

func (f *availabilityFlightDTO) toDomain() domain.FlightOffer {
	airline := "Vueling"
	if f.CarrierCode != "" && f.CarrierCode != "VY" {
		airline = f.CarrierCode
	}

	flightNumber := f.FlightID
	if f.CarrierCode != "" && flightNumber != "" {
		flightNumber = fmt.Sprintf("%s %s", f.CarrierCode, f.FlightID)
	}

	depParsed, depFormatted := parseVuelingTime(f.DepartureDate)
	arrParsed, arrFormatted := parseVuelingTime(f.ArrivalDate)

	var duration string
	if depParsed != nil && arrParsed != nil && arrParsed.After(*depParsed) {
		diff := arrParsed.Sub(*depParsed)
		hours := int(diff.Hours())
		mins := int(diff.Minutes()) % 60
		duration = fmt.Sprintf("%02d:%02d", hours, mins)
	}

	currency := f.Currency
	if currency == "" {
		currency = "EUR"
	}

	return domain.FlightOffer{
		TransportType:    domain.TransportTypeFlight,
		Airline:          airline,
		FlightNumber:     flightNumber,
		DepartureStation: f.DepartureStation,
		ArrivalStation:   f.ArrivalStation,
		DepartureTime:    depParsed,
		ArrivalTime:      arrParsed,
		DepartureRaw:     depFormatted,
		ArrivalRaw:       arrFormatted,
		Duration:         duration,
		Price: domain.Price{
			Amount:   f.Price,
			Currency: currency,
		},
		IsAvailable: f.IsAvailableDay,
		Status:      "available",
	}
}

func parseVuelingTime(raw string) (*time.Time, string) {
	// e.g. "2026-08-18T16:55:00"
	if len(raw) >= 19 {
		t, err := time.Parse("2006-01-02T15:04:05", raw[:19])
		if err == nil {
			return &t, t.Format("2006-01-02 15:04")
		}
	}
	if len(raw) >= 10 {
		return nil, raw[:10]
	}
	return nil, raw
}

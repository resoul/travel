package volotea

import (
	"fmt"
	"time"

	"github.com/resoul/travel/internal/domain"
)

// --- Countries DTO ---

type cultureNameDTO struct {
	Name string `json:"Name"`
}

type currencyDTO struct {
	Code    string `json:"Code"`
	Symbol  string `json:"Symbol"`
	Visible bool   `json:"Visible"`
}

type countryDTO struct {
	Code     string                    `json:"Code"`
	Prefix   string                    `json:"Prefix"`
	Culture  map[string]cultureNameDTO `json:"Culture"`
	Currency currencyDTO               `json:"Currency"`
}

func (c *countryDTO) toDomain() domain.Country {
	name := c.Code
	if en, ok := c.Culture["en-GB"]; ok && en.Name != "" {
		name = en.Name
	} else {
		for _, cult := range c.Culture {
			if cult.Name != "" {
				name = cult.Name
				break
			}
		}
	}

	return domain.Country{
		Code:     c.Code,
		Name:     name,
		Currency: c.Currency.Code,
	}
}

// --- Stations / Airports DTO ---

type stationCultureDTO struct {
	Name     string `json:"Name"`
	FullName string `json:"FullName"`
}

type marketDTO struct {
	Price         float64  `json:"Price"`
	Enabled       bool     `json:"Enabled"`
	FlightType    string   `json:"FlightType"`
	MinFlightDate string   `json:"MinFlightDate"`
	MaxFlightDate string   `json:"MaxFlightDate"`
	FareTypes     []string `json:"AvailableFareTypes"`
}

type stationDTO struct {
	StationCode string                       `json:"StationCode"`
	Culture     map[string]stationCultureDTO `json:"Culture"`
	Markets     map[string]marketDTO         `json:"Markets"`
	Lat         float64                      `json:"Lat"`
	Long        float64                      `json:"Long"`
	Enabled     bool                         `json:"Enabled"`
	IsBase      bool                         `json:"IsBase"`
	Country     string                       `json:"Country"`
	City        string                       `json:"City"`
}

func (s *stationDTO) toDomain(code string) domain.Airport {
	name := code
	if en, ok := s.Culture["en-GB"]; ok && en.Name != "" {
		name = en.Name
	} else {
		for _, cult := range s.Culture {
			if cult.Name != "" {
				name = cult.Name
				break
			}
		}
	}

	cityName := s.City
	if cityName == "" {
		cityName = name
	}

	return domain.Airport{
		Code: code,
		Name: name,
		Base: s.IsBase,
		City: domain.City{
			Name: cityName,
		},
		Country: domain.Country{
			Code: s.Country,
		},
		Coordinates: domain.Coordinates{
			Latitude:  s.Lat,
			Longitude: s.Long,
		},
	}
}

// --- Schedule DTO ---

type schedulePriceDTO struct {
	Price        float64 `json:"Price"`
	PriceWithFee float64 `json:"PriceWithFee"`
	FareType     string  `json:"FareType"`
	Currency     string  `json:"Currency"`
}

type flightScheduleDTO struct {
	Departure      string             `json:"Departure"` // YYYYMMDDHHMM
	Arrival        string             `json:"Arrival"`   // YYYYMMDDHHMM
	FlightNumber   string             `json:"FlightNumber"`
	Prices         []schedulePriceDTO `json:"Prices"`
	AvailableSeats int                `json:"AvailableSeats"`
	CarrierCode    string             `json:"CarrierCode"`
}

func (f *flightScheduleDTO) toDomain(origin, destination string) domain.FlightOffer {
	carrier := f.CarrierCode
	if carrier == "" {
		carrier = "Volotea"
	} else if carrier == "V7" {
		carrier = "Volotea"
	}

	depParsed, depFormatted := parseVoloteaTime(f.Departure)
	arrParsed, arrFormatted := parseVoloteaTime(f.Arrival)

	var duration string
	if depParsed != nil && arrParsed != nil {
		diff := arrParsed.Sub(*depParsed)
		hours := int(diff.Hours())
		mins := int(diff.Minutes()) % 60
		duration = fmt.Sprintf("%02d:%02d", hours, mins)
	}

	var priceAmount float64
	var currency string
	for _, p := range f.Prices {
		if p.Price > 0 {
			priceAmount = p.Price
			currency = p.Currency
			break
		}
	}
	if currency == "" {
		currency = "EUR"
	}

	flightNum := f.FlightNumber
	if f.CarrierCode != "" {
		flightNum = fmt.Sprintf("%s %s", f.CarrierCode, f.FlightNumber)
	}

	return domain.FlightOffer{
		TransportType:    domain.TransportTypeFlight,
		Airline:          carrier,
		FlightNumber:     flightNum,
		DepartureStation: origin,
		ArrivalStation:   destination,
		DepartureTime:    depParsed,
		ArrivalTime:      arrParsed,
		DepartureRaw:     depFormatted,
		ArrivalRaw:       arrFormatted,
		Duration:         duration,
		Price: domain.Price{
			Amount:   priceAmount,
			Currency: currency,
		},
		SeatsLeft:   f.AvailableSeats,
		IsAvailable: f.AvailableSeats > 0 || priceAmount > 0,
		Status:      "available",
	}
}

func parseVoloteaTime(raw string) (*time.Time, string) {
	// raw format: "202608250605" (12 digits)
	if len(raw) == 12 {
		t, err := time.Parse("200601021504", raw)
		if err == nil {
			return &t, t.Format("2006-01-02 15:04")
		}
	}
	return nil, raw
}

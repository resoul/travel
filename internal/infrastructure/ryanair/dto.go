package ryanair

import "github.com/resoul/travel/internal/domain"

// --- Airports DTOs ---

type coordinatesDTO struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type cityDTO struct {
	Name string `json:"name"`
	Code string `json:"code"`
}

type regionDTO struct {
	Name string `json:"name"`
	Code string `json:"code"`
}

type countryDTO struct {
	Code               string `json:"code"`
	ISO3Code           string `json:"iso3code"`
	Name               string `json:"name"`
	Currency           string `json:"currency"`
	DefaultAirportCode string `json:"defaultAirportCode"`
	Schengen           bool   `json:"schengen"`
}

type airportDTO struct {
	Code        string         `json:"code"`
	Name        string         `json:"name"`
	SeoName     string         `json:"seoName"`
	Aliases     []string       `json:"aliases"`
	Base        bool           `json:"base"`
	City        cityDTO        `json:"city"`
	Region      regionDTO      `json:"region"`
	Country     countryDTO     `json:"country"`
	Coordinates coordinatesDTO `json:"coordinates"`
	TimeZone    string         `json:"timeZone"`
}

type routeItemDTO struct {
	ArrivalAirport airportDTO `json:"arrivalAirport"`
}

func (a *airportDTO) toDomain() domain.Airport {
	return domain.Airport{
		Code:    a.Code,
		Name:    a.Name,
		SeoName: a.SeoName,
		Aliases: a.Aliases,
		Base:    a.Base,
		City: domain.City{
			Code: a.City.Code,
			Name: a.City.Name,
		},
		Country: domain.Country{
			Code:               a.Country.Code,
			ISO3Code:           a.Country.ISO3Code,
			Name:               a.Country.Name,
			Currency:           a.Country.Currency,
			DefaultAirportCode: a.Country.DefaultAirportCode,
			Schengen:           a.Country.Schengen,
		},
		Coordinates: domain.Coordinates{
			Latitude:  a.Coordinates.Latitude,
			Longitude: a.Coordinates.Longitude,
		},
		TimeZone: a.TimeZone,
	}
}

// --- Availability API DTOs ---

type fareDTO struct {
	Type             string  `json:"type"`
	Amount           float64 `json:"amount"`
	Count            int     `json:"count"`
	HasDiscount      bool    `json:"hasDiscount"`
	PublishedFare    float64 `json:"publishedFare"`
	DiscountAmount   float64 `json:"discountAmount"`
	DiscountInPct    int     `json:"discountInPercent"`
	HasPromoDiscount bool    `json:"hasPromoDiscount"`
	IsPrime          bool    `json:"isPrime"`
	HasBogof         bool    `json:"hasBogof"`
}

type regularFareDTO struct {
	FareKey string    `json:"fareKey"`
	Fares   []fareDTO `json:"fares"`
}

type segmentDTO struct {
	SegmentNr    int    `json:"segmentNr"`
	Origin       string `json:"origin"`
	Destination  string `json:"destination"`
	FlightNumber string `json:"flightNumber"`
}

type rawFlightDTO struct {
	FaresLeft    int            `json:"faresLeft"`
	FlightKey    string         `json:"flightKey"`
	FlightNumber string         `json:"flightNumber"`
	InfantsLeft  int            `json:"infantsLeft"`
	Duration     string         `json:"duration"`
	IsSSIMLoad   bool           `json:"isSSIMLoad"`
	OperatedBy   string         `json:"operatedBy"`
	RegularFare  regularFareDTO `json:"regularFare"`
	Segments     []segmentDTO   `json:"segments"`
	Time         []string       `json:"time"`    // local: ["2026-06-22T10:55:00.000", ...]
	TimeUTC      []string       `json:"timeUTC"` // utc:   ["2026-06-22T07:55:00.000Z", ...]
}

type rawDateDTO struct {
	DateOut string         `json:"dateOut"`
	Flights []rawFlightDTO `json:"flights"`
}

type rawTripDTO struct {
	Origin          string       `json:"origin"`
	OriginName      string       `json:"originName"`
	Destination     string       `json:"destination"`
	DestinationName string       `json:"destinationName"`
	Dates           []rawDateDTO `json:"dates"`
}

type availabilityResponseDTO struct {
	Currency string       `json:"currency"`
	Trips    []rawTripDTO `json:"trips"`
}

// --- Farfnd API DTOs ---

type farfndAirportDTO struct {
	CountryName string `json:"countryName"`
	IataCode    string `json:"iataCode"`
	Name        string `json:"name"`
	SeoName     string `json:"seoName"`
}

type farfndPriceDTO struct {
	Value        float64 `json:"value"`
	CurrencyCode string  `json:"currencyCode"`
}

type farfndOutboundDTO struct {
	DepartureAirport farfndAirportDTO `json:"departureAirport"`
	ArrivalAirport   farfndAirportDTO `json:"arrivalAirport"`
	DepartureDate    string           `json:"departureDate"`
	ArrivalDate      string           `json:"arrivalDate"`
	Price            farfndPriceDTO   `json:"price"`
	FlightKey        string           `json:"flightKey"`
	FlightNumber     string           `json:"flightNumber"`
}

type farfndFareItemDTO struct {
	Outbound farfndOutboundDTO `json:"outbound"`
}

type farfndResponseDTO struct {
	Fares []farfndFareItemDTO `json:"fares"`
	Total int                 `json:"total"`
}

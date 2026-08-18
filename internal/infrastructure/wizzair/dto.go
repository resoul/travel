package wizzair

import "github.com/resoul/travel/internal/domain"

type connectionDTO struct {
	IATA string `json:"iata"`
}

type cityDTO struct {
	IATA        string          `json:"iata"`
	ShortName   string          `json:"shortName"`
	CountryName string          `json:"countryName"`
	Connections []connectionDTO `json:"connections"`
}

func (c *cityDTO) toDomain() domain.City {
	return domain.City{
		Code: c.IATA,
		Name: c.ShortName,
	}
}

type mapResponseDTO struct {
	Cities []cityDTO `json:"cities"`
}

type flightItemDTO struct {
	DepartureStation string `json:"departureStation"`
	ArrivalStation   string `json:"arrivalStation"`
	From             string `json:"from"`
	To               string `json:"to"`
}

type flightRequestDTO struct {
	FlightList         []flightItemDTO `json:"flightList"`
	PriceType          string          `json:"priceType"`
	AdultCount         int             `json:"adultCount"`
	ChildCount         int             `json:"childCount"`
	InfantCount        int             `json:"infantCount"`
	MacStationsAllowed bool            `json:"macStationsAllowed"`
}

type priceDTO struct {
	Amount       float64 `json:"amount"`
	CurrencyCode string  `json:"currencyCode"`
}

type flightResultDTO struct {
	ArrivalStation   string   `json:"arrivalStation"`
	DepartureStation string   `json:"departureStation"`
	DepartureDate    string   `json:"departureDate"`
	Price            priceDTO `json:"price"`
}

type timetableResponseDTO struct {
	OutboundFlights []flightResultDTO `json:"outboundFlights"`
	ReturnFlights   []flightResultDTO `json:"returnFlights"`
}

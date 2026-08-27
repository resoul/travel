package flytap

import "github.com/resoul/travel/internal/domain"

// AirportSearchRequest represents the request payload for originSearch and destinationSearch.
type AirportSearchRequest struct {
	AirlineIds   []string `json:"airlineIds"`
	Language     string   `json:"language"`
	Market       string   `json:"market"`
	Origin       string   `json:"origin,omitempty"`
	PayWithMiles bool     `json:"payWithMiles"`
	TripType     string   `json:"tripType"`
}

// AirportItemDTO represents an airport item in originSearch or destinationSearch.
type AirportItemDTO struct {
	Airport     string `json:"airport"`
	AirportName string `json:"airportName"`
	City        string `json:"city"`
	CityName    string `json:"cityName"`
	Country     string `json:"country"`
	TapRoute    bool   `json:"tapRoute"`
	Direct      *bool  `json:"direct,omitempty"`
}

func (dto AirportItemDTO) toDomain() domain.Airport {
	name := dto.AirportName
	if name == "" {
		name = dto.CityName
	}
	if name == "" {
		name = dto.Airport
	}

	cityName := dto.CityName
	if cityName == "" {
		cityName = dto.City
	}

	return domain.Airport{
		Code: dto.Airport,
		Name: name,
		City: domain.City{
			Code: dto.City,
			Name: cityName,
		},
		Country: domain.Country{
			Code: dto.Country,
			Name: dto.Country,
		},
	}
}

// CalendarRequest represents the request payload for calendar.
type CalendarRequest struct {
	CabinClass   string `json:"cabinClass"`
	Destination  string `json:"destination"`
	Market       string `json:"market"`
	Month        string `json:"month"`
	Origin       string `json:"origin"`
	PaxType      string `json:"paxType"`
	PayWithMiles bool   `json:"payWithMiles"`
	StarAlliance bool   `json:"starAlliance"`
	TripType     string `json:"tripType"`
	Year         string `json:"year"`
}

// CalendarResponse represents the response from calendar endpoint.
type CalendarResponse struct {
	Data struct {
		BestPriceForDates []CalendarDateItem `json:"bestPriceForDates"`
	} `json:"data"`
}

// CalendarDateItem represents a single day fare in the calendar.
type CalendarDateItem struct {
	DepartureDate    string   `json:"departureDate"`
	DepartureAirport string   `json:"departureAirport"`
	ArrivalAirport   string   `json:"arrivalAirport"`
	BestTotalPrice   *float64 `json:"bestTotalPrice"`
	Currency         string   `json:"currency"`
	CabinClass       string   `json:"cabinClass"`
	SoldOut          bool     `json:"soldOut"`
	NoFlight         *bool    `json:"noFlight"`
}

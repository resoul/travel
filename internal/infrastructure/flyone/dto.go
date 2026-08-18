package flyone

import (
	"fmt"
	"time"

	"github.com/resoul/travel/internal/domain"
)

// --- Routes DTO ---

type routeItemDTO struct {
	DepCode        string   `json:"depCode"`
	DepAirportName string   `json:"depAirportName"`
	CountryCode    string   `json:"countryCode"`
	CountryName    string   `json:"countryName"`
	ArrCodes       []string `json:"arrCodes"`
}

type getRoutesResponseDTO struct {
	Routes []routeItemDTO `json:"routes"`
}

func (r *routeItemDTO) toDomain() domain.Airport {
	return domain.Airport{
		Code: r.DepCode,
		Name: r.DepAirportName,
		City: domain.City{
			Code: r.DepCode,
			Name: r.DepAirportName,
		},
		Country: domain.Country{
			Code: r.CountryCode,
			Name: r.CountryName,
		},
	}
}

// --- CMS Route Fare DTO ---

type cmsRouteFareRequestDTO struct {
	DepCity string `json:"depCity"`
	ArrCity string `json:"arrCity"`
	Token   string `json:"token"`
}

type calendarDayDTO struct {
	Day   int     `json:"day"`
	Price float64 `json:"price"`
}

type calendarMonthDTO struct {
	Month int              `json:"month"`
	Days  []calendarDayDTO `json:"days"`
}

type calendarYearDTO struct {
	Year     int                `json:"year"`
	Currency string             `json:"currency"`
	Month    []calendarMonthDTO `json:"month"`
}

type cmsRouteFareResponseDTO struct {
	Calendar []calendarYearDTO `json:"calendar"`
	Result   struct {
		IsSuccess bool `json:"isSuccess"`
	} `json:"result"`
}

func formatFlightOffer(origin, destination string, year, month, day int, price float64, currency string) domain.FlightOffer {
	if currency == "" {
		currency = "EUR"
	}

	dateStr := fmt.Sprintf("%04d-%02d-%02d", year, month, day)
	t, err := time.Parse("2006-01-02", dateStr)

	var tPtr *time.Time
	if err == nil {
		tPtr = &t
	}

	return domain.FlightOffer{
		TransportType:    domain.TransportTypeFlight,
		Airline:          "FlyOne",
		FlightNumber:     "Direct",
		DepartureStation: origin,
		ArrivalStation:   destination,
		DepartureTime:    tPtr,
		DepartureRaw:     dateStr,
		Price: domain.Price{
			Amount:   price,
			Currency: currency,
		},
		IsAvailable: true,
		Status:      "available",
	}
}

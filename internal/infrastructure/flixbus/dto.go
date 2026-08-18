package flixbus

import (
	"fmt"
	"strings"
	"time"

	"github.com/resoul/travel/internal/domain"
)

// --- Autocomplete DTO ---

type locationDTO struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

type autocompleteCityDTO struct {
	ID              string      `json:"id"`
	Name            string      `json:"name"`
	Country         string      `json:"country"`
	District        string      `json:"district"`
	Location        locationDTO `json:"location"`
	HasTrainStation bool        `json:"has_train_station"`
	IsFlixbusCity   bool        `json:"is_flixbus_city"`
}

func (a *autocompleteCityDTO) toDomain() domain.Airport {
	return domain.Airport{
		Code: a.ID,
		Name: a.Name,
		City: domain.City{
			Code: a.ID,
			Name: a.Name,
		},
		Country: domain.Country{
			Code: strings.ToUpper(a.Country),
		},
		Coordinates: domain.Coordinates{
			Latitude:  a.Location.Lat,
			Longitude: a.Location.Lon,
		},
	}
}

// --- Reachable Cities DTO ---

type reachablePriceItemDTO struct {
	Min          float64 `json:"min"`
	Max          float64 `json:"max"`
	Avg          float64 `json:"avg"`
	FormattedMin string  `json:"formatted_min"`
	FormattedAvg string  `json:"formatted_avg"`
	FormattedMax string  `json:"formatted_max"`
}

type reachableCityDTO struct {
	UUID                   string                           `json:"uuid"`
	Name                   string                           `json:"name"`
	Country                string                           `json:"country"`
	Price                  map[string]reachablePriceItemDTO `json:"price"`
	TransportationCategory []string                         `json:"transportation_category"`
	Location               locationDTO                      `json:"location"`
}

type reachableResponseDTO struct {
	Result []reachableCityDTO `json:"result"`
}

func (r *reachableCityDTO) toDomain() domain.Airport {
	name := r.Name
	if len(r.Price) > 0 {
		for curr, p := range r.Price {
			if p.Min > 0 {
				name = fmt.Sprintf("%s (from %.2f %s)", r.Name, p.Min, curr)
				break
			}
		}
	}

	return domain.Airport{
		Code: r.UUID,
		Name: name,
		City: domain.City{
			Code: r.UUID,
			Name: r.Name,
		},
		Country: domain.Country{
			Code: strings.ToUpper(r.Country),
		},
		Coordinates: domain.Coordinates{
			Latitude:  r.Location.Lat,
			Longitude: r.Location.Lon,
		},
	}
}

// --- Search / Trips DTO ---

type rideLocationDTO struct {
	Date      string `json:"date"` // e.g. "2026-08-27T05:55:00+03:00"
	CityID    string `json:"city_id"`
	StationID string `json:"station_id"`
}

type durationDTO struct {
	Hours   int `json:"hours"`
	Minutes int `json:"minutes"`
}

type priceDTO struct {
	Total                  float64 `json:"total"`
	TotalWithPlatformFee   float64 `json:"total_with_platform_fee"`
	Average                float64 `json:"average"`
	AverageWithPlatformFee float64 `json:"average_with_platform_fee"`
}

type rideResultDTO struct {
	UID          string          `json:"uid"`
	Status       string          `json:"status"`
	TransferType string          `json:"transfer_type"` // e.g. "Direct"
	Provider     string          `json:"provider"`      // e.g. "flixbus"
	Departure    rideLocationDTO `json:"departure"`
	Arrival      rideLocationDTO `json:"arrival"`
	Duration     durationDTO     `json:"duration"`
	Price        priceDTO        `json:"price"`
}

type tripItemDTO struct {
	DepartureCityID string                   `json:"departure_city_id"`
	ArrivalCityID   string                   `json:"arrival_city_id"`
	Date            string                   `json:"date"`
	Results         map[string]rideResultDTO `json:"results"`
}

type namedEntityDTO struct {
	Name        string `json:"name"`
	CountryCode string `json:"country_code"`
	CityID      string `json:"city_id"`
}

type searchResponseDTO struct {
	Trips    []tripItemDTO             `json:"trips"`
	Cities   map[string]namedEntityDTO `json:"cities"`
	Stations map[string]namedEntityDTO `json:"stations"`
}

func (r *rideResultDTO) toDomain(currency string, cities, stations map[string]namedEntityDTO) domain.FlightOffer {
	provider := strings.Title(r.Provider)
	if provider == "" {
		provider = "FlixBus"
	}

	depStationName := ""
	if st, ok := stations[r.Departure.StationID]; ok && st.Name != "" {
		depStationName = st.Name
	} else if ct, ok := cities[r.Departure.CityID]; ok && ct.Name != "" {
		depStationName = ct.Name
	}

	arrStationName := ""
	if st, ok := stations[r.Arrival.StationID]; ok && st.Name != "" {
		arrStationName = st.Name
	} else if ct, ok := cities[r.Arrival.CityID]; ok && ct.Name != "" {
		arrStationName = ct.Name
	}

	depParsed, depFormatted := parseFlixBusTime(r.Departure.Date)
	arrParsed, arrFormatted := parseFlixBusTime(r.Arrival.Date)

	durationStr := fmt.Sprintf("%02d:%02d", r.Duration.Hours, r.Duration.Minutes)

	// User requested to display price with platform_fee as shown on the website
	priceAmount := r.Price.TotalWithPlatformFee
	if priceAmount <= 0 {
		priceAmount = r.Price.Total
	}

	if currency == "" {
		currency = "EUR"
	}

	flightType := r.TransferType
	if flightType == "" {
		flightType = "Direct"
	}

	transportType := domain.TransportTypeBus
	if strings.Contains(strings.ToLower(r.Provider), "train") {
		transportType = domain.TransportTypeTrain
	}

	isAvailable := r.Status == "available" || (r.Status == "" && priceAmount > 0)

	return domain.FlightOffer{
		TransportType:    transportType,
		Airline:          provider,
		FlightNumber:     flightType,
		DepartureStation: depStationName,
		ArrivalStation:   arrStationName,
		DepartureTime:    depParsed,
		ArrivalTime:      arrParsed,
		DepartureRaw:     depFormatted,
		ArrivalRaw:       arrFormatted,
		Duration:         durationStr,
		Price: domain.Price{
			Amount:   priceAmount,
			Currency: currency,
		},
		IsAvailable: isAvailable,
		Status:      r.Status,
	}
}

func parseFlixBusTime(raw string) (*time.Time, string) {
	// e.g. "2026-08-27T05:55:00+03:00"
	t, err := time.Parse(time.RFC3339, raw)
	if err == nil {
		return &t, t.Format("2006-01-02 15:04")
	}
	if len(raw) >= 16 {
		return nil, strings.Replace(raw[:16], "T", " ", 1)
	}
	return nil, raw
}

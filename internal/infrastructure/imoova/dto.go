package imoova

import (
	"fmt"
	"time"

	"github.com/resoul/travel/internal/domain"
)

// --- Relocations Table DTO ---

type relocationTableItemDTO struct {
	ID                  string   `json:"id"`
	DepartureCityName   string   `json:"departure_city_name"`
	DeliveryCityName    string   `json:"delivery_city_name"`
	AvailableFromDate   string   `json:"available_from_date"`
	AvailableToDate     string   `json:"available_to_date"`
	LatestDepartureDate string   `json:"latest_departure_date"`
	VehicleName         string   `json:"vehicle_name"`
	VehicleType         string   `json:"vehicle_type"`
	VehicleSleeps       *int     `json:"vehicle_sleeps"`
	VehicleSeatbelts    *int     `json:"vehicle_seatbelts"`
	Currency            string   `json:"currency"`
	FuelAmount          *float64 `json:"fuel_amount"`
	ExtraUnits          *int     `json:"extra_units"`
}

type relocationsTableResponseDTO struct {
	Data struct {
		RelocationsTable []relocationTableItemDTO `json:"relocationsTable"`
	} `json:"data"`
}

func (item *relocationTableItemDTO) toDomain() domain.FlightOffer {
	curr := item.Currency
	if curr == "" {
		curr = "USD"
	}

	depParsed, depFormatted := parseTime(item.AvailableFromDate)
	arrParsed, arrFormatted := parseTime(item.AvailableToDate)

	vehicleDescription := item.VehicleName
	if vehicleDescription == "" {
		vehicleDescription = "Campervan"
	}

	details := ""
	if item.VehicleSleeps != nil && *item.VehicleSleeps > 0 {
		details += fmt.Sprintf("%d berths", *item.VehicleSleeps)
	}
	if item.FuelAmount != nil && *item.FuelAmount > 0 {
		if details != "" {
			details += ", "
		}
		details += fmt.Sprintf("Fuel bonus: %.0f %s", *item.FuelAmount, curr)
	}

	if details != "" {
		vehicleDescription = fmt.Sprintf("%s (%s)", vehicleDescription, details)
	}

	durationStr := "Relocation Deal"
	if depFormatted != "" && arrFormatted != "" {
		durationStr = fmt.Sprintf("%s — %s", depFormatted, arrFormatted)
	}

	return domain.FlightOffer{
		TransportType:    domain.TransportTypeCar,
		Airline:          "imoova",
		FlightNumber:     vehicleDescription,
		DepartureStation: item.DepartureCityName,
		ArrivalStation:   item.DeliveryCityName,
		DepartureTime:    depParsed,
		ArrivalTime:      arrParsed,
		DepartureRaw:     depFormatted,
		ArrivalRaw:       arrFormatted,
		Duration:         durationStr,
		Price: domain.Price{
			Amount:   1.0,
			Currency: curr,
		},
		IsAvailable: true,
		Status:      "available",
	}
}

// --- Relocation Points DTO ---

type secondaryPointDTO struct {
	CityName      string   `json:"city_name"`
	Count         int      `json:"count"`
	RelocationIDs []string `json:"relocation_ids"`
}

type relocationPointDTO struct {
	CityName        string              `json:"city_name"`
	CitySlug        string              `json:"city_slug"`
	Region          string              `json:"region"`
	Lat             float64             `json:"lat"`
	Lng             float64             `json:"lng"`
	Count           int                 `json:"count"`
	RelocationIDs   []string            `json:"relocation_ids"`
	SecondaryPoints []secondaryPointDTO `json:"secondary_points"`
}

type relocationPointsResponseDTO struct {
	Data struct {
		RelocationPoints []relocationPointDTO `json:"relocationPoints"`
	} `json:"data"`
}

func (p *relocationPointDTO) toDomain() domain.Airport {
	return domain.Airport{
		Code: p.CitySlug,
		Name: fmt.Sprintf("%s [%s] (%d deals)", p.CityName, p.Region, p.Count),
		City: domain.City{
			Code: p.CitySlug,
			Name: p.CityName,
		},
		Country: domain.Country{
			Code: p.Region,
			Name: p.Region,
		},
		Coordinates: domain.Coordinates{
			Latitude:  p.Lat,
			Longitude: p.Lng,
		},
	}
}

// --- Regions DTO ---

type regionItemDTO struct {
	Count  int    `json:"count"`
	Region string `json:"region"`
}

type homepageQueryResponseDTO struct {
	Data struct {
		RelocationsByRegion []regionItemDTO `json:"relocationsByRegion"`
	} `json:"data"`
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

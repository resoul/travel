package tictactrip

// CityDTO represents an autocomplete city/station item from Tictactrip API.
type CityDTO struct {
	CityID     int     `json:"city_id"`
	StationID  int     `json:"station_id"`
	LocalName  string  `json:"local_name"`
	Latitude   float64 `json:"latitude"`
	Longitude  float64 `json:"longitude"`
	UniqueName string  `json:"unique_name"`
	IsCity     bool    `json:"iscity"`
	Score      int     `json:"score"`
	Serviced   bool    `json:"serviced"`
}

// PopularCityDTO represents popular destination city from Tictactrip API.
type PopularCityDTO struct {
	ID         int     `json:"id"`
	UniqueName string  `json:"unique_name"`
	LocalName  string  `json:"local_name"`
	Latitude   float64 `json:"latitude"`
	Longitude  float64 `json:"longitude"`
	CityID     int     `json:"city_id"`
	NbSearch   string  `json:"nb_search"`
	Popular    bool    `json:"popular"`
	IsCity     bool    `json:"iscity"`
}

// PriceCalendarResponseDTO represents daily lowest price from priceCalendar/month.
type PriceCalendarResponseDTO struct {
	Date string           `json:"date"`
	Trip *CalendarTripDTO `json:"trip"`
}

// CalendarTripDTO represents trip details in price calendar.
type CalendarTripDTO struct {
	DepartureUnixUTC int64    `json:"departureUnixUtc"`
	DepartureOffset  string   `json:"departureOffset"`
	ArrivalUnixUTC   int64    `json:"arrivalUnixUtc"`
	ArrivalOffset    string   `json:"arrivalOffset"`
	DurationMinutes  int      `json:"durationMinutes"`
	PriceCents       int      `json:"priceCents"`
	TransportType    string   `json:"transportType"`
	Companies        []string `json:"companies"`
	NumberOfStops    int      `json:"numberOfStops"`
}

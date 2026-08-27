package sata

// AirportDTO represents an airport returned by the Azores Airlines API.
type AirportDTO struct {
	Code         string   `json:"code"`
	Name         string   `json:"name"`
	CityCode     string   `json:"cityCode"`
	CityName     string   `json:"cityName"`
	CountryCode  string   `json:"countryCode"`
	CountryName  string   `json:"countryName"`
	RegionName   string   `json:"regionName"`
	Destinations []string `json:"destinations"`
}

// CalendarItemDTO represents a daily minimum price in Azores Airlines calendar.
type CalendarItemDTO struct {
	From     string  `json:"from"`
	To       string  `json:"to"`
	FromDate string  `json:"fromDate"`
	ToDate   string  `json:"toDate"`
	Min      float64 `json:"min"`
	Cur      string  `json:"cur"`
}

// CalendarWrapperDTO represents an optional wrapper object for some international routes.
type CalendarWrapperDTO struct {
	Dates map[string]CalendarItemDTO `json:"Dates"`
}

// RealtimeStatusDTO represents Azores Airlines realtime flight status.
type RealtimeStatusDTO struct {
	Airport    string            `json:"airport"`
	Arrivals   []FlightStatusDTO `json:"arrivals"`
	Departures []FlightStatusDTO `json:"departures"`
}

// FlightStatusDTO represents a single flight in the realtime status board.
type FlightStatusDTO struct {
	DateTime     string `json:"datetime"`
	FlightNumber string `json:"flightNumber"`
	Airport      string `json:"airport"`
	Status       string `json:"status"`
}

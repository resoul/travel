package eurowings

// AirportsListResponseDTO represents response from /search/airports/list.
type AirportsListResponseDTO struct {
	Countries []CountryDTO `json:"countries"`
	Stations  []StationDTO `json:"stations"`
}

// CountryDTO represents a country in the Eurowings network.
type CountryDTO struct {
	CountryCode  string `json:"countryCode"`
	CurrencyCode string `json:"currencyCode"`
	Name         string `json:"name"`
}

// StationDTO represents an airport or station in the Eurowings network.
type StationDTO struct {
	CountryCode string  `json:"countryCode"`
	TLC         string  `json:"tlc"`
	Name        string  `json:"name"`
	Type        string  `json:"type"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
}

// RoutesResponseDTO represents response from /search/airports?origin={code}.
type RoutesResponseDTO struct {
	MatchedTLCs []string `json:"matched-tlcs"`
}

// ScheduleDatesResponseDTO represents response from /search/flight-schedule/dates.
type ScheduleDatesResponseDTO struct {
	Meta     ScheduleMetaDTO   `json:"meta"`
	Sections []ScheduleSection `json:"sections"`
}

// ScheduleMetaDTO contains route metadata.
type ScheduleMetaDTO struct {
	Stations StationsDTO `json:"stations"`
}

// StationsDTO contains airport stations info.
type StationsDTO struct {
	AirlineCode     string `json:"airlineCode"`
	Origin          string `json:"origin"`
	OriginName      string `json:"originName"`
	Destination     string `json:"destination"`
	DestinationName string `json:"destinationName"`
}

// ScheduleSection contains month batches.
type ScheduleSection struct {
	BookableMonths []BookableMonthDTO `json:"bookableMonths"`
}

// BookableMonthDTO contains active flight dates for a month.
type BookableMonthDTO struct {
	Year     int   `json:"year"`
	Month    int   `json:"month"`
	Dates    []int `json:"dates"`
	Bookable bool  `json:"bookable"`
}

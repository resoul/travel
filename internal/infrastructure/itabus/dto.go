package itabus

// StationDTO represents a city/station from ItaBus API
type StationDTO struct {
	Code           string   `json:"code"`
	Name           string   `json:"name"`
	CountryCode    string   `json:"country_code"`
	Synonyms       []string `json:"synonyms"`
	SequenceNumber int      `json:"sequence_number"`
}

// StationsResponseDTO is the wrapper for stations endpoint
type StationsResponseDTO struct {
	Action string       `json:"action"`
	Locale string       `json:"locale"`
	Data   []StationDTO `json:"data"`
}

// DateDTO represents an available date with fare
type DateDTO struct {
	Date   string  `json:"date"`
	Amount float64 `json:"amount"`
}

// DatesResponseDTO is the wrapper for dates endpoint
type DatesResponseDTO struct {
	Action string    `json:"action"`
	Locale string    `json:"locale"`
	Data   []DateDTO `json:"data"`
}

// LocationDTO represents origin/destination in routes
type LocationDTO struct {
	Name string `json:"name"`
	Code string `json:"code"`
}

// ServiceDTO represents onboard services (WiFi, WC, etc)
type ServiceDTO struct {
	Code        string `json:"code"`
	Description string `json:"description"`
}

// PassengerFareDTO represents fare for a passenger
type PassengerFareDTO struct {
	BucketCode     string  `json:"bucket_code"`
	OriginalPrice  float64 `json:"original_price"`
	PassengerID    string  `json:"passenger_id"`
	InventoryClass string  `json:"inventory_class"`
	Price          float64 `json:"price"`
	TariffCode     string  `json:"tariff_code"`
}

// ProductItemDTO represents a ticket product
type ProductItemDTO struct {
	ProductCode     string             `json:"product_code"`
	ProductFamilyID string             `json:"product_family_id"`
	ProductType     string             `json:"product_type"`
	LegIds          []string           `json:"leg_ids"`
	Name            string             `json:"name"`
	PassengerFares  []PassengerFareDTO `json:"passenger_fares"`
}

// RateClassDTO represents a rate class (WOW, ECONOMY, etc)
type RateClassDTO struct {
	ProductFamilyID      string           `json:"product_family_id"`
	Price                float64          `json:"price"`
	OriginalPrice        float64          `json:"original_price"`
	Items                []ProductItemDTO `json:"items"`
	PassengersCount      int              `json:"passengers_count"`
	AveragePrice         float64          `json:"average_price"`
	AverageOriginalPrice float64          `json:"average_original_price"`
}

// RateCategoryDTO represents rate category (COMFORT, etc)
type RateCategoryDTO struct {
	COMFORT RateClassDTO `json:"COMFORT"`
}

// RatesDTO represents all available rates for a route
type RatesDTO struct {
	Available bool                       `json:"available"`
	ITABUS    map[string]RateCategoryDTO `json:"ITABUS"`
}

// RouteDTO represents a single travel route
type RouteDTO struct {
	Direction          string       `json:"direction"`
	TravelDuration     string       `json:"travel_duration"`
	ID                 string       `json:"id"`
	Origin             LocationDTO  `json:"origin"`
	OriginCode         string       `json:"origin_code"`
	Destination        LocationDTO  `json:"destination"`
	DepartureTimestamp string       `json:"departure_timestamp"`
	ArrivalTimestamp   string       `json:"arrival_timestamp"`
	ServiceIdentifier  string       `json:"service_identifier"`
	ServiceName        string       `json:"service_name"`
	Services           []ServiceDTO `json:"services"`
	TransferNumber     int          `json:"transferNumber"`
	Available          bool         `json:"available"`
	Rates              RatesDTO     `json:"rates"`
}

// DirectionDTO represents routes for a single direction
type DirectionDTO struct {
	Direction   string     `json:"direction"`
	Origin      string     `json:"origin"`
	Destination string     `json:"destination"`
	Routes      []RouteDTO `json:"routes"`
}

// TravelsDataDTO contains outbound/return routes
type TravelsDataDTO struct {
	Outbound *DirectionDTO `json:"outbound"`
	Return   *DirectionDTO `json:"return"`
}

// TravelsResponseDTO is the wrapper for travels endpoint
type TravelsResponseDTO struct {
	Action string         `json:"action"`
	Locale string         `json:"locale"`
	Data   TravelsDataDTO `json:"data"`
}

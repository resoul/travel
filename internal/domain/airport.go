package domain

// Coordinates represents geographic latitude and longitude coordinates.
type Coordinates struct {
	Latitude  float64
	Longitude float64
}

// City represents a city entity.
type City struct {
	Code string
	Name string
}

// Country represents a country entity.
type Country struct {
	Code               string
	ISO3Code           string
	Name               string
	Currency           string
	DefaultAirportCode string
	Schengen           bool
}

// Airport represents an airport entity.
type Airport struct {
	Code        string
	Name        string
	SeoName     string
	Aliases     []string
	Base        bool
	City        City
	Country     Country
	Coordinates Coordinates
	TimeZone    string
}

// Route represents a flight connection between two airports.
type Route struct {
	OriginCode      string
	DestinationCode string
	Destination     Airport
}

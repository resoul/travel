package domain

import "time"

// Price represents monetary amount and currency.
type Price struct {
	Amount   float64
	Currency string
}

// TransportType represents the type of transportation mode.
type TransportType string

const (
	TransportTypeFlight TransportType = "Flight"
	TransportTypeBus    TransportType = "Bus"
	TransportTypeTrain  TransportType = "Train"
	TransportTypeCar    TransportType = "Car"
	TransportTypeCruise TransportType = "Cruise"
	TransportTypeHotel  TransportType = "Hotel"
)

// HotelOffer represents a hotel/accommodation deal.
type HotelOffer struct {
	ID          string
	Name        string
	City        string
	Address     string
	Stars       float64
	Rating      float64
	ReviewCount int
	Price       Price
	Nights      int
	RoomType    string
	URL         string
	ImageURL    string
}

// TripRoomOffer represents detailed room information from a hotel page.
type TripRoomOffer struct {
	ID        string
	Name      string
	Area      string
	Beds      string
	Guests    int
	HasWindow string
	Smoking   string
	Amenities []string
	ImageURL  string
	Price     Price
}

// CarHireOffer represents a rental car option found across rental agencies.
type CarHireOffer struct {
	ID           string
	Model        string
	Category     string
	Transmission string
	Seats        int
	Doors        int
	Bags         int
	Supplier     string
	SupplierLogo string
	PricePerDay  Price
	TotalPrice   Price
	PickupDate   string
	ReturnDate   string
	PickupPlace  string
	ReturnPlace  string
	Features     []string
	ImageURL     string
	BookingURL   string
}

// CarHireCriteria defines parameters for rental car searches.
type CarHireCriteria struct {
	PickupCityID   string // e.g. "39050"
	PickupCityName string // e.g. "Otopeni"
	PickupCode     string // e.g. "OTP"
	PickupAddress  string // e.g. "Bucharest Henri Coandă International Airport (OTP)"
	PickupDate     string // "YYYY-MM-DD HH:mm"
	ReturnCityID   string
	ReturnCityName string
	ReturnCode     string
	ReturnAddress  string
	ReturnDate     string // "YYYY-MM-DD HH:mm"
	CountryID      string // e.g. "63" for Romania, "95" for Spain
	CountryName    string // e.g. "Romania"
	DriverAge      string // "30-60"
	Currency       string // "USD"
	Limit          int
}

// HotelSearchCriteria defines search parameters for accommodations.
type HotelSearchCriteria struct {
	CityID   string // e.g. "19216"
	CityName string // e.g. "Mamaia"
	CheckIn  string // "YYYY-MM-DD"
	CheckOut string // "YYYY-MM-DD"
	Rooms    int
	Adults   int
	Children int
	Currency string
	Sort     string // "priceLowToHigh", "rank", etc.
	Limit    int
}

// FlightOffer represents a travel offer option found across airlines and ground operators.
type FlightOffer struct {
	TransportType    TransportType
	Airline          string
	FlightNumber     string
	DepartureStation string
	ArrivalStation   string
	DepartureTime    *time.Time
	ArrivalTime      *time.Time
	DepartureRaw     string
	ArrivalRaw       string
	Duration         string
	Price            Price
	SeatsLeft        int
	IsAvailable      bool
	Status           string
}

// FlightSearchCriteria defines search parameters for flights.
type FlightSearchCriteria struct {
	Origin         string
	Destination    string
	DepartureDate  string // Format: YYYY-MM-DD
	ReturnDate     string // Format: YYYY-MM-DD, optional
	Adults         int
	Teens          int
	Children       int
	Infants        int
	RoundTrip      bool
	FlexDaysBefore int
	FlexDaysAfter  int
}

// WizzairSearchCriteria defines search parameters specific to Wizzair searches.
type WizzairSearchCriteria struct {
	DepartureStation string
	ArrivalStation   string
	FromDate         string // Format: YYYY-MM-DD
	ToDate           string // Format: YYYY-MM-DD
	Adults           int
	Children         int
	Infants          int
	PriceType        string
}

// CruiseSearchCriteria defines search parameters for cruise searches.
type CruiseSearchCriteria struct {
	DestinationID string
	CruiseLineID  string
	Month         int
	Year          int
	DurationMin   int
	DurationMax   int
	Limit         int
}

// CruiseLine represents a cruise operator in the matrix.
type CruiseLine struct {
	ID   int
	Name string
}

// CruiseDestination represents a destination region and its ID.
type CruiseDestination struct {
	ID   int
	Name string
}

// TictactripCity represents a city or railway station in Tictactrip's European network.
type TictactripCity struct {
	ID         int
	LocalName  string
	UniqueName string
	Latitude   float64
	Longitude  float64
	Score      int
	Serviced   bool
	NbSearch   string
}

// TictactripCalendarDay represents the lowest train/bus fare for a specific day.
type TictactripCalendarDay struct {
	Date            string
	HasTrip         bool
	TransportType   string
	Companies       []string
	Price           Price
	DurationMinutes int
	DepartureTime   time.Time
	ArrivalTime     time.Time
	NumberOfStops   int
}

// TrenitaliaStation represents an Italian railway/bus station in the Trenitalia network.
type TrenitaliaStation struct {
	ID         int
	Name       string
	Value      string
	IsFrecce   bool
	IsEurocity bool
}

// TrenitaliaSearchCriteria defines search parameters for Trenitalia train journeys.
type TrenitaliaSearchCriteria struct {
	OriginID        int
	DestinationID   int
	OriginName      string
	DestinationName string
	DepartureDate   string // "YYYY-MM-DD"
	DepartureTime   string // "HH:mm"
	Adults          int
	Children        int
	FrecceOnly      bool
	RegionalOnly    bool
	NoChanges       bool
	Limit           int
}

// OBiletSearchCriteria defines search parameters for oBilet bus journeys.
type OBiletSearchCriteria struct {
	OriginID        int
	DestinationID   int
	OriginName      string
	DestinationName string
	DepartureDate   string // "YYYY-MM-DD"
	Limit           int
}

// PitchupSearchCriteria defines search parameters for Pitchup campsites and outdoor holiday stays.
type PitchupSearchCriteria struct {
	Country    string // e.g. "France", "England", "Spain", "Italy", "Germany", "USA"
	Region     string // optional region/location
	ArriveDate string // "YYYY-MM-DD"
	DepartDate string // "YYYY-MM-DD"
	Adults     int
	Limit      int
}

// HipcampSearchCriteria defines search parameters for Hipcamp glamping and outdoor spots.
type HipcampSearchCriteria struct {
	Country string // e.g. "united-states", "united-kingdom", "canada", "australia", "france"
	Region  string // e.g. "california", "england", "ontario", "new-south-wales"
	Limit   int
}

// CampspaceSearchCriteria defines search parameters for Campspace micro-camping spots.
type CampspaceSearchCriteria struct {
	Category string // e.g. "tent-pitches", "camper-sites", "glamping", "treehouses", "yurts"
	Limit    int
}

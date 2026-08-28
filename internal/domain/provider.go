package domain

import "context"

// RyanairProvider defines operations supported by Ryanair integration.
type RyanairProvider interface {
	SearchFlights(ctx context.Context, criteria FlightSearchCriteria) ([]FlightOffer, error)
	GetAirports(ctx context.Context) ([]Airport, error)
	GetRoutes(ctx context.Context, originIATA string) ([]Airport, error)
	GetAvailabilities(ctx context.Context, origin, destination string) ([]string, error)
}

// WizzairProvider defines operations supported by Wizzair integration.
type WizzairProvider interface {
	SearchFlights(ctx context.Context, criteria WizzairSearchCriteria) ([]FlightOffer, error)
	GetMap(ctx context.Context) ([]City, error)
}

// VoloteaProvider defines operations supported by Volotea integration.
type VoloteaProvider interface {
	SearchFlights(ctx context.Context, criteria FlightSearchCriteria) ([]FlightOffer, error)
	GetAirports(ctx context.Context) ([]Airport, error)
	GetRoutes(ctx context.Context, originIATA string) ([]Airport, error)
	GetSchedule(ctx context.Context, origin, destination string) ([]FlightOffer, error)
	GetDates(ctx context.Context, origin, destination string) ([]string, error)
	GetCountries(ctx context.Context) ([]Country, error)
}

// VuelingProvider defines operations supported by Vueling integration.
type VuelingProvider interface {
	SearchFlights(ctx context.Context, criteria FlightSearchCriteria) ([]FlightOffer, error)
	GetAirports(ctx context.Context) ([]Airport, error)
	GetRoutes(ctx context.Context, originIATA string) ([]Airport, error)
	GetSchedule(ctx context.Context, origin, destination string, year, month, monthsRange int) ([]FlightOffer, error)
	GetDates(ctx context.Context, origin, destination string, year, month, monthsRange int) ([]string, error)
}

// FlixBusProvider defines operations supported by FlixBus integration.
type FlixBusProvider interface {
	SearchTrips(ctx context.Context, criteria FlightSearchCriteria) ([]FlightOffer, error)
	GetCities(ctx context.Context, query string) ([]Airport, error)
	GetReachable(ctx context.Context, cityQueryOrID string, limit int) ([]Airport, error)
}

// AirBalticProvider defines operations supported by airBaltic integration.
type AirBalticProvider interface {
	SearchFlights(ctx context.Context, criteria FlightSearchCriteria) ([]FlightOffer, error)
	GetAirports(ctx context.Context) ([]Airport, error)
	GetRoutes(ctx context.Context, originIATA string) ([]Airport, error)
	GetDates(ctx context.Context, origin, destination string) ([]string, error)
}

// FlyOneProvider defines operations supported by FlyOne integration.
type FlyOneProvider interface {
	SearchFlights(ctx context.Context, criteria FlightSearchCriteria) ([]FlightOffer, error)
	GetAirports(ctx context.Context) ([]Airport, error)
	GetRoutes(ctx context.Context, originIATA string) ([]Airport, error)
	GetDates(ctx context.Context, origin, destination string) ([]string, error)
}

// MovacarProvider defines operations supported by Movacar integration.
type MovacarProvider interface {
	GetOffers(ctx context.Context, criteria FlightSearchCriteria) ([]FlightOffer, error)
	GetLocations(ctx context.Context) ([]Airport, error)
}

// ImoovaProvider defines operations supported by imoova campervan relocation integration.
type ImoovaProvider interface {
	GetOffers(ctx context.Context, criteria FlightSearchCriteria) ([]FlightOffer, error)
	GetLocations(ctx context.Context) ([]Airport, error)
	GetRegions(ctx context.Context) ([]Country, error)
}

// DriiveMeProvider defines operations supported by DriiveMe 1-euro car relocation integration.
type DriiveMeProvider interface {
	Login(ctx context.Context, email, password string) error
	IsAuthenticated() bool
	GetOffers(ctx context.Context, criteria FlightSearchCriteria) ([]FlightOffer, error)
	GetCities(ctx context.Context, query string) ([]Airport, error)
	GetAvailabilities(ctx context.Context, transportID string) ([]string, error)
}

// IndiGoRadarFare represents a destination fare returned by IndiGo fare radar.
type IndiGoRadarFare struct {
	Origin      string
	OriginCity  string
	Destination string
	DestCity    string
	TravelDate  string
	FlightTime  string
	Price       Price
}

// IndiGoProvider defines operations supported by IndiGo integration.
type IndiGoProvider interface {
	GetFareRadar(ctx context.Context, originIATA string) ([]IndiGoRadarFare, error)
	GetFareCalendar(ctx context.Context, origin, destination, startDate, endDate, currency string) ([]FlightOffer, error)
	GetDates(ctx context.Context, origin, destination string) ([]string, error)
}

// FlyTapProvider defines operations supported by TAP Air Portugal integration.
type FlyTapProvider interface {
	GetAirports(ctx context.Context) ([]Airport, error)
	GetRoutes(ctx context.Context, originIATA string) ([]Airport, error)
	GetCalendar(ctx context.Context, origin, destination string, year, month int, market string) ([]FlightOffer, error)
	GetDates(ctx context.Context, origin, destination string) ([]string, error)
}

// CruiseProvider defines operations supported by Cruise integration.
type CruiseProvider interface {
	SearchCruises(ctx context.Context, criteria CruiseSearchCriteria) ([]FlightOffer, error)
	GetCruiseLines(ctx context.Context) ([]CruiseLine, error)
	GetCruiseDestinations(ctx context.Context) ([]CruiseDestination, error)
}

// AgodaProvider defines operations supported by Agoda accommodation integration.
type AgodaProvider interface {
	SearchHotels(ctx context.Context, criteria HotelSearchCriteria) ([]HotelOffer, error)
	GetCountries(ctx context.Context, languageID int) ([]Country, error)
}

// TripProvider defines operations supported by Trip.com accommodation and car hire integration.
type TripProvider interface {
	SearchHotels(ctx context.Context, criteria HotelSearchCriteria) ([]HotelOffer, error)
	GetHotelDetails(ctx context.Context, hotelID, checkIn, checkOut string, adults, rooms int, currency string) ([]TripRoomOffer, error)
	SearchCars(ctx context.Context, criteria CarHireCriteria) ([]CarHireOffer, error)
}

// TictactripProvider defines operations for European train, bus, and multimodal routes.
type TictactripProvider interface {
	AutocompleteCities(ctx context.Context, query string) ([]TictactripCity, error)
	GetPopularDestinations(ctx context.Context, fromCity string, limit int) ([]TictactripCity, error)
	GetMonthlyPriceCalendar(ctx context.Context, originID, destinationID int, month string) ([]TictactripCalendarDay, error)
}

// TrenitaliaProvider defines operations for Italian rail and Le Frecce high-speed train searches.
type TrenitaliaProvider interface {
	GetStations(ctx context.Context) ([]TrenitaliaStation, error)
	SearchTrains(ctx context.Context, criteria TrenitaliaSearchCriteria) ([]FlightOffer, error)
}

// NorwegianProvider defines operations for Norwegian Air Shuttle low-fare calendars.
type NorwegianProvider interface {
	GetFareCalendar(ctx context.Context, origin, destination string, year, month int, currency string) ([]FlightOffer, error)
}

// OBiletProvider defines operations for Turkish and regional bus searches on oBilet.
type OBiletProvider interface {
	SearchBuses(ctx context.Context, criteria OBiletSearchCriteria) ([]FlightOffer, error)
}

// EurowingsProvider defines operations for Eurowings airports, route networks, and flight schedules.
type EurowingsProvider interface {
	GetAirports(ctx context.Context) ([]Airport, error)
	GetRoutesFromOrigin(ctx context.Context, origin string) ([]string, error)
	GetFlightDates(ctx context.Context, origin, destination string) ([]string, error)
}

// TransaviaProvider defines operations for Transavia (Air France-KLM Group) low-fare calendars.
type TransaviaProvider interface {
	GetFareCalendar(ctx context.Context, origin, destination string, year, month int, adults int) ([]FlightOffer, error)
}

// PitchupProvider defines operations for campsites, glamping, and outdoor stays on Pitchup.
type PitchupProvider interface {
	SearchCampsites(ctx context.Context, criteria PitchupSearchCriteria) ([]FlightOffer, error)
}

// HipcampProvider defines operations for glamping, nature stays, and outdoor camping on Hipcamp.
type HipcampProvider interface {
	SearchSpots(ctx context.Context, criteria HipcampSearchCriteria) ([]FlightOffer, error)
}

// CampspaceProvider defines operations for sustainable micro-camping and nature spots on Campspace.
type CampspaceProvider interface {
	SearchSpots(ctx context.Context, criteria CampspaceSearchCriteria) ([]FlightOffer, error)
}

// SATAProvider defines operations for Azores Airlines / SATA Air Açores route network, fare calendar, and live flight status.
type SATAProvider interface {
	GetAirports(ctx context.Context) ([]Airport, error)
	GetRoutesFromOrigin(ctx context.Context, origin string) ([]string, error)
	GetFareCalendar(ctx context.Context, origin, destination string) ([]FlightOffer, error)
}

// LevelProvider defines operations for LEVEL (IAG Group long-haul low-cost airline).
type LevelProvider interface {
	GetAirports(ctx context.Context) ([]Airport, error)
	GetRoutesFromOrigin(ctx context.Context, origin string) ([]string, error)
	GetFlightDates(ctx context.Context, origin, destination string) ([]string, error)
	GetFareCalendar(ctx context.Context, origin, destination string, year, month int) ([]FlightOffer, error)
}

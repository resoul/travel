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

package usecase

import (
	"context"
	"fmt"

	"github.com/resoul/travel/internal/domain"
)

// ListAirportsUseCase handles fetching airports and route maps.
type ListAirportsUseCase struct {
	ryanair   domain.RyanairProvider
	wizzair   domain.WizzairProvider
	volotea   domain.VoloteaProvider
	vueling   domain.VuelingProvider
	flixbus   domain.FlixBusProvider
	airbaltic domain.AirBalticProvider
	flyone    domain.FlyOneProvider
	movacar   domain.MovacarProvider
	imoova    domain.ImoovaProvider
}

// NewListAirportsUseCase creates a new ListAirportsUseCase.
func NewListAirportsUseCase(
	ryanair domain.RyanairProvider,
	wizzair domain.WizzairProvider,
	volotea domain.VoloteaProvider,
	vueling domain.VuelingProvider,
	flixbus domain.FlixBusProvider,
	airbaltic domain.AirBalticProvider,
	flyone domain.FlyOneProvider,
	movacar domain.MovacarProvider,
	imoova domain.ImoovaProvider,
) *ListAirportsUseCase {
	return &ListAirportsUseCase{
		ryanair:   ryanair,
		wizzair:   wizzair,
		volotea:   volotea,
		vueling:   vueling,
		flixbus:   flixbus,
		airbaltic: airbaltic,
		flyone:    flyone,
		movacar:   movacar,
		imoova:    imoova,
	}
}

// GetRyanairAirports returns all active Ryanair airports.
func (uc *ListAirportsUseCase) GetRyanairAirports(ctx context.Context) ([]domain.Airport, error) {
	return uc.ryanair.GetAirports(ctx)
}

// GetRyanairRoutes returns airports reachable from a given origin IATA code.
func (uc *ListAirportsUseCase) GetRyanairRoutes(ctx context.Context, originIATA string) ([]domain.Airport, error) {
	if originIATA == "" {
		return nil, fmt.Errorf("origin IATA code is required")
	}
	return uc.ryanair.GetRoutes(ctx, originIATA)
}

// GetWizzairMap returns connected cities from Wizzair map.
func (uc *ListAirportsUseCase) GetWizzairMap(ctx context.Context) ([]domain.City, error) {
	return uc.wizzair.GetMap(ctx)
}

// GetVoloteaAirports returns all airports in Volotea network.
func (uc *ListAirportsUseCase) GetVoloteaAirports(ctx context.Context) ([]domain.Airport, error) {
	return uc.volotea.GetAirports(ctx)
}

// GetVoloteaRoutes returns airports reachable from a given origin IATA code in Volotea network.
func (uc *ListAirportsUseCase) GetVoloteaRoutes(ctx context.Context, originIATA string) ([]domain.Airport, error) {
	if originIATA == "" {
		return nil, fmt.Errorf("origin IATA code is required")
	}
	return uc.volotea.GetRoutes(ctx, originIATA)
}

// GetVoloteaCountries returns all countries from Volotea dist endpoint.
func (uc *ListAirportsUseCase) GetVoloteaCountries(ctx context.Context) ([]domain.Country, error) {
	return uc.volotea.GetCountries(ctx)
}

// GetVuelingAirports returns all active airports in Vueling network.
func (uc *ListAirportsUseCase) GetVuelingAirports(ctx context.Context) ([]domain.Airport, error) {
	return uc.vueling.GetAirports(ctx)
}

// GetVuelingRoutes returns airports reachable from a given origin IATA code in Vueling network.
func (uc *ListAirportsUseCase) GetVuelingRoutes(ctx context.Context, originIATA string) ([]domain.Airport, error) {
	if originIATA == "" {
		return nil, fmt.Errorf("origin IATA code is required")
	}
	return uc.vueling.GetRoutes(ctx, originIATA)
}

// GetFlixBusCities returns FlixBus cities matching a query.
func (uc *ListAirportsUseCase) GetFlixBusCities(ctx context.Context, query string) ([]domain.Airport, error) {
	return uc.flixbus.GetCities(ctx, query)
}

// GetFlixBusRoutes returns reachable destinations from a city in FlixBus network.
func (uc *ListAirportsUseCase) GetFlixBusRoutes(ctx context.Context, cityQueryOrID string, limit int) ([]domain.Airport, error) {
	if cityQueryOrID == "" {
		return nil, fmt.Errorf("city name or UUID is required")
	}
	return uc.flixbus.GetReachable(ctx, cityQueryOrID, limit)
}

// GetAirBalticAirports returns all primary airports in airBaltic network.
func (uc *ListAirportsUseCase) GetAirBalticAirports(ctx context.Context) ([]domain.Airport, error) {
	return uc.airbaltic.GetAirports(ctx)
}

// GetAirBalticRoutes returns reachable destinations from an origin airport in airBaltic network.
func (uc *ListAirportsUseCase) GetAirBalticRoutes(ctx context.Context, originIATA string) ([]domain.Airport, error) {
	if originIATA == "" {
		return nil, fmt.Errorf("origin IATA code is required")
	}
	return uc.airbaltic.GetRoutes(ctx, originIATA)
}

// GetFlyOneAirports returns all departure airports in FlyOne network.
func (uc *ListAirportsUseCase) GetFlyOneAirports(ctx context.Context) ([]domain.Airport, error) {
	return uc.flyone.GetAirports(ctx)
}

// GetFlyOneRoutes returns reachable destinations from an origin airport in FlyOne network.
func (uc *ListAirportsUseCase) GetFlyOneRoutes(ctx context.Context, originIATA string) ([]domain.Airport, error) {
	if originIATA == "" {
		return nil, fmt.Errorf("origin IATA code is required")
	}
	return uc.flyone.GetRoutes(ctx, originIATA)
}

// GetMovacarLocations returns all active cities/stations in Movacar network with their offer counts.
func (uc *ListAirportsUseCase) GetMovacarLocations(ctx context.Context) ([]domain.Airport, error) {
	return uc.movacar.GetLocations(ctx)
}

// GetImoovaLocations returns all active departure cities and route counts in imoova network.
func (uc *ListAirportsUseCase) GetImoovaLocations(ctx context.Context) ([]domain.Airport, error) {
	return uc.imoova.GetLocations(ctx)
}

// GetImoovaRegions returns deal counts per region (US, CA, EU, AU, NZ, SA) from imoova.
func (uc *ListAirportsUseCase) GetImoovaRegions(ctx context.Context) ([]domain.Country, error) {
	return uc.imoova.GetRegions(ctx)
}

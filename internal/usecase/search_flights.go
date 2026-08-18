package usecase

import (
	"context"
	"fmt"

	"github.com/resoul/travel/internal/domain"
)

// SearchFlightsUseCase handles flight search workflows.
type SearchFlightsUseCase struct {
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

// NewSearchFlightsUseCase creates a new SearchFlightsUseCase.
func NewSearchFlightsUseCase(
	ryanair domain.RyanairProvider,
	wizzair domain.WizzairProvider,
	volotea domain.VoloteaProvider,
	vueling domain.VuelingProvider,
	flixbus domain.FlixBusProvider,
	airbaltic domain.AirBalticProvider,
	flyone domain.FlyOneProvider,
	movacar domain.MovacarProvider,
	imoova domain.ImoovaProvider,
) *SearchFlightsUseCase {
	return &SearchFlightsUseCase{
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

// SearchRyanair searches Ryanair flights matching given criteria.
func (uc *SearchFlightsUseCase) SearchRyanair(ctx context.Context, criteria domain.FlightSearchCriteria) ([]domain.FlightOffer, error) {
	if criteria.Origin == "" || criteria.Destination == "" {
		return nil, fmt.Errorf("origin and destination are required")
	}
	if criteria.DepartureDate == "" {
		return nil, fmt.Errorf("departure date is required")
	}

	return uc.ryanair.SearchFlights(ctx, criteria)
}

// SearchWizzair searches Wizzair flights matching given criteria.
func (uc *SearchFlightsUseCase) SearchWizzair(ctx context.Context, criteria domain.WizzairSearchCriteria) ([]domain.FlightOffer, error) {
	if criteria.DepartureStation == "" || criteria.ArrivalStation == "" {
		return nil, fmt.Errorf("origin and destination are required")
	}
	if criteria.FromDate == "" {
		return nil, fmt.Errorf("start date is required")
	}

	return uc.wizzair.SearchFlights(ctx, criteria)
}

// SearchVolotea searches Volotea flights matching given criteria.
func (uc *SearchFlightsUseCase) SearchVolotea(ctx context.Context, criteria domain.FlightSearchCriteria) ([]domain.FlightOffer, error) {
	if criteria.Origin == "" || criteria.Destination == "" {
		return nil, fmt.Errorf("origin and destination are required")
	}

	return uc.volotea.SearchFlights(ctx, criteria)
}

// SearchVueling searches Vueling flights matching given criteria.
func (uc *SearchFlightsUseCase) SearchVueling(ctx context.Context, criteria domain.FlightSearchCriteria) ([]domain.FlightOffer, error) {
	if criteria.Origin == "" || criteria.Destination == "" {
		return nil, fmt.Errorf("origin and destination are required")
	}

	return uc.vueling.SearchFlights(ctx, criteria)
}

// SearchFlixBus searches FlixBus trips matching given criteria.
func (uc *SearchFlightsUseCase) SearchFlixBus(ctx context.Context, criteria domain.FlightSearchCriteria) ([]domain.FlightOffer, error) {
	if criteria.Origin == "" || criteria.Destination == "" {
		return nil, fmt.Errorf("origin and destination are required")
	}

	return uc.flixbus.SearchTrips(ctx, criteria)
}

// SearchAirBaltic searches airBaltic flights matching given criteria.
func (uc *SearchFlightsUseCase) SearchAirBaltic(ctx context.Context, criteria domain.FlightSearchCriteria) ([]domain.FlightOffer, error) {
	if criteria.Origin == "" || criteria.Destination == "" {
		return nil, fmt.Errorf("origin and destination are required")
	}

	return uc.airbaltic.SearchFlights(ctx, criteria)
}

// SearchFlyOne searches FlyOne flights matching given criteria.
func (uc *SearchFlightsUseCase) SearchFlyOne(ctx context.Context, criteria domain.FlightSearchCriteria) ([]domain.FlightOffer, error) {
	if criteria.Origin == "" || criteria.Destination == "" {
		return nil, fmt.Errorf("origin and destination are required")
	}

	return uc.flyone.SearchFlights(ctx, criteria)
}

// SearchMovacar searches Movacar relocation offers matching given criteria.
func (uc *SearchFlightsUseCase) SearchMovacar(ctx context.Context, criteria domain.FlightSearchCriteria) ([]domain.FlightOffer, error) {
	return uc.movacar.GetOffers(ctx, criteria)
}

// SearchImoova searches imoova campervan relocation offers matching given criteria.
func (uc *SearchFlightsUseCase) SearchImoova(ctx context.Context, criteria domain.FlightSearchCriteria) ([]domain.FlightOffer, error) {
	return uc.imoova.GetOffers(ctx, criteria)
}

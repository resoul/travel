package usecase

import (
	"context"
	"fmt"

	"github.com/resoul/travel/internal/domain"
)

// SearchFlightsUseCase handles flight search workflows.
type SearchFlightsUseCase struct {
	ryanair    domain.RyanairProvider
	wizzair    domain.WizzairProvider
	volotea    domain.VoloteaProvider
	vueling    domain.VuelingProvider
	flixbus    domain.FlixBusProvider
	airbaltic  domain.AirBalticProvider
	flyone     domain.FlyOneProvider
	movacar    domain.MovacarProvider
	imoova     domain.ImoovaProvider
	driiveme   domain.DriiveMeProvider
	cruise     domain.CruiseProvider
	agoda      domain.AgodaProvider
	trip       domain.TripProvider
	trenitalia domain.TrenitaliaProvider
	obilet     domain.OBiletProvider
	pitchup    domain.PitchupProvider
	hipcamp    domain.HipcampProvider
	campspace  domain.CampspaceProvider
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
	driiveme domain.DriiveMeProvider,
	cruise domain.CruiseProvider,
	agoda domain.AgodaProvider,
	trip domain.TripProvider,
	trenitalia domain.TrenitaliaProvider,
	obilet domain.OBiletProvider,
	pitchup domain.PitchupProvider,
	hipcamp domain.HipcampProvider,
	campspace domain.CampspaceProvider,
) *SearchFlightsUseCase {
	return &SearchFlightsUseCase{
		ryanair:    ryanair,
		wizzair:    wizzair,
		volotea:    volotea,
		vueling:    vueling,
		flixbus:    flixbus,
		airbaltic:  airbaltic,
		flyone:     flyone,
		movacar:    movacar,
		imoova:     imoova,
		driiveme:   driiveme,
		cruise:     cruise,
		agoda:      agoda,
		trip:       trip,
		trenitalia: trenitalia,
		obilet:     obilet,
		pitchup:    pitchup,
		hipcamp:    hipcamp,
		campspace:  campspace,
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

// SearchDriiveMe searches DriiveMe 1-euro car relocation offers matching given criteria.
func (uc *SearchFlightsUseCase) SearchDriiveMe(ctx context.Context, criteria domain.FlightSearchCriteria) ([]domain.FlightOffer, error) {
	return uc.driiveme.GetOffers(ctx, criteria)
}

// LoginDriiveMe authenticates with DriiveMe using email and password.
func (uc *SearchFlightsUseCase) LoginDriiveMe(ctx context.Context, email, password string) error {
	return uc.driiveme.Login(ctx, email, password)
}

// GetDriiveMeAvailabilities retrieves booking slot availabilities for a DriiveMe transport ID.
func (uc *SearchFlightsUseCase) GetDriiveMeAvailabilities(ctx context.Context, transportID string) ([]string, error) {
	return uc.driiveme.GetAvailabilities(ctx, transportID)
}

// SearchCruises searches available cruises matching given criteria.
func (uc *SearchFlightsUseCase) SearchCruises(ctx context.Context, criteria domain.CruiseSearchCriteria) ([]domain.FlightOffer, error) {
	return uc.cruise.SearchCruises(ctx, criteria)
}

// GetCruiseLines retrieves all available cruise lines from the matrix.
func (uc *SearchFlightsUseCase) GetCruiseLines(ctx context.Context) ([]domain.CruiseLine, error) {
	return uc.cruise.GetCruiseLines(ctx)
}

// GetCruiseDestinations retrieves all available cruising destination regions from the matrix.
func (uc *SearchFlightsUseCase) GetCruiseDestinations(ctx context.Context) ([]domain.CruiseDestination, error) {
	return uc.cruise.GetCruiseDestinations(ctx)
}

// SearchAgodaHotels searches accommodations on Agoda matching given criteria.
func (uc *SearchFlightsUseCase) SearchAgodaHotels(ctx context.Context, criteria domain.HotelSearchCriteria) ([]domain.HotelOffer, error) {
	return uc.agoda.SearchHotels(ctx, criteria)
}

// SearchTripHotels searches accommodations on Trip.com matching given criteria.
func (uc *SearchFlightsUseCase) SearchTripHotels(ctx context.Context, criteria domain.HotelSearchCriteria) ([]domain.HotelOffer, error) {
	return uc.trip.SearchHotels(ctx, criteria)
}

// GetTripHotelDetails retrieves all room options and details for a Trip.com hotel ID.
func (uc *SearchFlightsUseCase) GetTripHotelDetails(ctx context.Context, hotelID, checkIn, checkOut string, adults, rooms int, currency string) ([]domain.TripRoomOffer, error) {
	return uc.trip.GetHotelDetails(ctx, hotelID, checkIn, checkOut, adults, rooms, currency)
}

// SearchTripCars searches rental cars on Trip.com matching given criteria.
func (uc *SearchFlightsUseCase) SearchTripCars(ctx context.Context, criteria domain.CarHireCriteria) ([]domain.CarHireOffer, error) {
	return uc.trip.SearchCars(ctx, criteria)
}

// SearchTrenitaliaTrains searches Italian trains and Le Frecce high-speed rail on Trenitalia.
func (uc *SearchFlightsUseCase) SearchTrenitaliaTrains(ctx context.Context, criteria domain.TrenitaliaSearchCriteria) ([]domain.FlightOffer, error) {
	return uc.trenitalia.SearchTrains(ctx, criteria)
}

// SearchOBiletBuses searches bus journeys in Turkey and regional routes on oBilet.
func (uc *SearchFlightsUseCase) SearchOBiletBuses(ctx context.Context, criteria domain.OBiletSearchCriteria) ([]domain.FlightOffer, error) {
	return uc.obilet.SearchBuses(ctx, criteria)
}

// SearchPitchupCampsites searches campsites, glamping, and holiday parks on Pitchup.
func (uc *SearchFlightsUseCase) SearchPitchupCampsites(ctx context.Context, criteria domain.PitchupSearchCriteria) ([]domain.FlightOffer, error) {
	return uc.pitchup.SearchCampsites(ctx, criteria)
}

// SearchHipcampSpots searches outdoor camping and glamping spots on Hipcamp.
func (uc *SearchFlightsUseCase) SearchHipcampSpots(ctx context.Context, criteria domain.HipcampSearchCriteria) ([]domain.FlightOffer, error) {
	return uc.hipcamp.SearchSpots(ctx, criteria)
}

// SearchCampspaceSpots searches sustainable micro-camping and nature spots on Campspace.
func (uc *SearchFlightsUseCase) SearchCampspaceSpots(ctx context.Context, criteria domain.CampspaceSearchCriteria) ([]domain.FlightOffer, error) {
	return uc.campspace.SearchSpots(ctx, criteria)
}

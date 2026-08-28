package usecase

import (
	"context"
	"fmt"

	"github.com/resoul/travel/internal/domain"
)

// FlightDatesUseCase handles queries for available flight schedules.
type FlightDatesUseCase struct {
	ryanair    domain.RyanairProvider
	volotea    domain.VoloteaProvider
	vueling    domain.VuelingProvider
	airbaltic  domain.AirBalticProvider
	flyone     domain.FlyOneProvider
	indigo     domain.IndiGoProvider
	flytap     domain.FlyTapProvider
	tictactrip domain.TictactripProvider
	norwegian  domain.NorwegianProvider
	eurowings  domain.EurowingsProvider
	transavia  domain.TransaviaProvider
	sata       domain.SATAProvider
	level      domain.LevelProvider
}

// NewFlightDatesUseCase creates a new FlightDatesUseCase.
func NewFlightDatesUseCase(
	ryanair domain.RyanairProvider,
	volotea domain.VoloteaProvider,
	vueling domain.VuelingProvider,
	airbaltic domain.AirBalticProvider,
	flyone domain.FlyOneProvider,
	indigo domain.IndiGoProvider,
	flytap domain.FlyTapProvider,
	tictactrip domain.TictactripProvider,
	norwegian domain.NorwegianProvider,
	eurowings domain.EurowingsProvider,
	transavia domain.TransaviaProvider,
	sata domain.SATAProvider,
	level domain.LevelProvider,
) *FlightDatesUseCase {
	return &FlightDatesUseCase{
		ryanair:    ryanair,
		volotea:    volotea,
		vueling:    vueling,
		airbaltic:  airbaltic,
		flyone:     flyone,
		indigo:     indigo,
		flytap:     flytap,
		tictactrip: tictactrip,
		norwegian:  norwegian,
		eurowings:  eurowings,
		transavia:  transavia,
		sata:       sata,
		level:      level,
	}
}

// GetRyanairDates returns scheduled flight dates between origin and destination.
func (uc *FlightDatesUseCase) GetRyanairDates(ctx context.Context, origin, destination string) ([]string, error) {
	if origin == "" || destination == "" {
		return nil, fmt.Errorf("origin and destination are required")
	}
	return uc.ryanair.GetAvailabilities(ctx, origin, destination)
}

// GetVoloteaDates returns scheduled flight dates between origin and destination.
func (uc *FlightDatesUseCase) GetVoloteaDates(ctx context.Context, origin, destination string) ([]string, error) {
	if origin == "" || destination == "" {
		return nil, fmt.Errorf("origin and destination are required")
	}
	return uc.volotea.GetDates(ctx, origin, destination)
}

// GetVoloteaSchedule returns full scheduled flight offers between origin and destination.
func (uc *FlightDatesUseCase) GetVoloteaSchedule(ctx context.Context, origin, destination string) ([]domain.FlightOffer, error) {
	if origin == "" || destination == "" {
		return nil, fmt.Errorf("origin and destination are required")
	}
	return uc.volotea.GetSchedule(ctx, origin, destination)
}

// GetVuelingDates returns scheduled flight dates between origin and destination.
func (uc *FlightDatesUseCase) GetVuelingDates(ctx context.Context, origin, destination string, year, month, monthsRange int) ([]string, error) {
	if origin == "" || destination == "" {
		return nil, fmt.Errorf("origin and destination are required")
	}
	return uc.vueling.GetDates(ctx, origin, destination, year, month, monthsRange)
}

// GetVuelingSchedule returns full scheduled flight offers between origin and destination.
func (uc *FlightDatesUseCase) GetVuelingSchedule(ctx context.Context, origin, destination string, year, month, monthsRange int) ([]domain.FlightOffer, error) {
	if origin == "" || destination == "" {
		return nil, fmt.Errorf("origin and destination are required")
	}
	return uc.vueling.GetSchedule(ctx, origin, destination, year, month, monthsRange)
}

// GetAirBalticDates returns scheduled flight dates between origin and destination.
func (uc *FlightDatesUseCase) GetAirBalticDates(ctx context.Context, origin, destination string) ([]string, error) {
	if origin == "" || destination == "" {
		return nil, fmt.Errorf("origin and destination are required")
	}
	return uc.airbaltic.GetDates(ctx, origin, destination)
}

// GetFlyOneDates returns scheduled flight dates between origin and destination.
func (uc *FlightDatesUseCase) GetFlyOneDates(ctx context.Context, origin, destination string) ([]string, error) {
	if origin == "" || destination == "" {
		return nil, fmt.Errorf("origin and destination are required")
	}
	return uc.flyone.GetDates(ctx, origin, destination)
}

// GetIndiGoDates returns scheduled flight dates between origin and destination.
func (uc *FlightDatesUseCase) GetIndiGoDates(ctx context.Context, origin, destination string) ([]string, error) {
	if origin == "" || destination == "" {
		return nil, fmt.Errorf("origin and destination are required")
	}
	return uc.indigo.GetDates(ctx, origin, destination)
}

// GetIndiGoFareCalendar returns date-by-date low fares between origin and destination.
func (uc *FlightDatesUseCase) GetIndiGoFareCalendar(ctx context.Context, origin, destination, startDate, endDate, currency string) ([]domain.FlightOffer, error) {
	if origin == "" || destination == "" {
		return nil, fmt.Errorf("origin and destination are required")
	}
	return uc.indigo.GetFareCalendar(ctx, origin, destination, startDate, endDate, currency)
}

// GetIndiGoFareRadar returns lowest fares and destination suggestions from an origin airport.
func (uc *FlightDatesUseCase) GetIndiGoFareRadar(ctx context.Context, origin string) ([]domain.IndiGoRadarFare, error) {
	if origin == "" {
		return nil, fmt.Errorf("origin is required")
	}
	return uc.indigo.GetFareRadar(ctx, origin)
}

// GetTictactripMonthlyPriceCalendar retrieves the lowest train/bus prices for every day of a month on Tictactrip.
func (uc *FlightDatesUseCase) GetTictactripMonthlyPriceCalendar(ctx context.Context, originID, destinationID int, month string) ([]domain.TictactripCalendarDay, error) {
	return uc.tictactrip.GetMonthlyPriceCalendar(ctx, originID, destinationID, month)
}

// GetFlyTapCalendar returns date-by-date low fares for a specific month and year from TAP Air Portugal.
func (uc *FlightDatesUseCase) GetFlyTapCalendar(ctx context.Context, origin, destination string, year, month int, market string) ([]domain.FlightOffer, error) {
	if origin == "" || destination == "" {
		return nil, fmt.Errorf("origin and destination are required")
	}
	return uc.flytap.GetCalendar(ctx, origin, destination, year, month, market)
}

// GetNorwegianFareCalendar retrieves the lowest flight fares across a month from Norwegian Air Shuttle.
func (uc *FlightDatesUseCase) GetNorwegianFareCalendar(ctx context.Context, origin, destination string, year, month int, currency string) ([]domain.FlightOffer, error) {
	if origin == "" || destination == "" {
		return nil, fmt.Errorf("origin and destination are required")
	}
	return uc.norwegian.GetFareCalendar(ctx, origin, destination, year, month, currency)
}

// GetFlyTapDates returns scheduled flight dates between origin and destination from TAP Air Portugal.
func (uc *FlightDatesUseCase) GetFlyTapDates(ctx context.Context, origin, destination string) ([]string, error) {
	if origin == "" || destination == "" {
		return nil, fmt.Errorf("origin and destination are required")
	}
	return uc.flytap.GetDates(ctx, origin, destination)
}

// GetEurowingsDates retrieves scheduled flight dates between origin and destination from Eurowings.
func (uc *FlightDatesUseCase) GetEurowingsDates(ctx context.Context, origin, destination string) ([]string, error) {
	if origin == "" || destination == "" {
		return nil, fmt.Errorf("origin and destination are required")
	}
	return uc.eurowings.GetFlightDates(ctx, origin, destination)
}

// GetTransaviaFareCalendar returns daily lowest fares for a given route and month from Transavia.
func (uc *FlightDatesUseCase) GetTransaviaFareCalendar(ctx context.Context, origin, destination string, year, month, adults int) ([]domain.FlightOffer, error) {
	if origin == "" || destination == "" {
		return nil, fmt.Errorf("origin and destination are required")
	}
	return uc.transavia.GetFareCalendar(ctx, origin, destination, year, month, adults)
}

// GetSATAFareCalendar returns daily lowest fares across available future dates from Azores Airlines / SATA.
func (uc *FlightDatesUseCase) GetSATAFareCalendar(ctx context.Context, origin, destination string) ([]domain.FlightOffer, error) {
	if origin == "" || destination == "" {
		return nil, fmt.Errorf("origin and destination are required")
	}
	return uc.sata.GetFareCalendar(ctx, origin, destination)
}

// GetLevelFlightDates returns scheduled flight dates across 365 days on LEVEL.
func (uc *FlightDatesUseCase) GetLevelFlightDates(ctx context.Context, origin, destination string) ([]string, error) {
	if origin == "" || destination == "" {
		return nil, fmt.Errorf("origin and destination are required")
	}
	return uc.level.GetFlightDates(ctx, origin, destination)
}

// GetLevelFareCalendar returns daily lowest fares for a given route, year, and month from LEVEL.
func (uc *FlightDatesUseCase) GetLevelFareCalendar(ctx context.Context, origin, destination string, year, month int) ([]domain.FlightOffer, error) {
	if origin == "" || destination == "" {
		return nil, fmt.Errorf("origin and destination are required")
	}
	return uc.level.GetFareCalendar(ctx, origin, destination, year, month)
}

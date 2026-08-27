package usecase

import (
	"context"
	"fmt"

	"github.com/resoul/travel/internal/domain"
)

// FlightDatesUseCase handles queries for available flight schedules.
type FlightDatesUseCase struct {
	ryanair   domain.RyanairProvider
	volotea   domain.VoloteaProvider
	vueling   domain.VuelingProvider
	airbaltic domain.AirBalticProvider
	flyone    domain.FlyOneProvider
	indigo    domain.IndiGoProvider
	flytap    domain.FlyTapProvider
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
) *FlightDatesUseCase {
	return &FlightDatesUseCase{
		ryanair:   ryanair,
		volotea:   volotea,
		vueling:   vueling,
		airbaltic: airbaltic,
		flyone:    flyone,
		indigo:    indigo,
		flytap:    flytap,
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

// GetFlyTapCalendar returns date-by-date low fares for a specific month and year from TAP Air Portugal.
func (uc *FlightDatesUseCase) GetFlyTapCalendar(ctx context.Context, origin, destination string, year, month int, market string) ([]domain.FlightOffer, error) {
	if origin == "" || destination == "" {
		return nil, fmt.Errorf("origin and destination are required")
	}
	return uc.flytap.GetCalendar(ctx, origin, destination, year, month, market)
}

// GetFlyTapDates returns scheduled flight dates between origin and destination from TAP Air Portugal.
func (uc *FlightDatesUseCase) GetFlyTapDates(ctx context.Context, origin, destination string) ([]string, error) {
	if origin == "" || destination == "" {
		return nil, fmt.Errorf("origin and destination are required")
	}
	return uc.flytap.GetDates(ctx, origin, destination)
}

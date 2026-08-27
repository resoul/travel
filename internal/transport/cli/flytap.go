package cli

import (
	"fmt"
	"time"

	"github.com/resoul/travel/internal/usecase"
	"github.com/spf13/cobra"
)

func newFlyTapCmd(
	airportsUC *usecase.ListAirportsUseCase,
	datesUC *usecase.FlightDatesUseCase,
	presenter *Presenter,
) *cobra.Command {
	flytapCmd := &cobra.Command{
		Use:   "tap",
		Short: "TAP Air Portugal (TP) flight search, route map, and low fare calendar commands",
	}

	// tap airports
	airportsCmd := &cobra.Command{
		Use:   "airports",
		Short: "List all operating origin airports in the TAP Air Portugal network",
		RunE: func(cmd *cobra.Command, args []string) error {
			airports, err := airportsUC.GetFlyTapAirports(cmd.Context())
			if err != nil {
				return fmt.Errorf("failed to get TAP airports: %w", err)
			}

			presenter.PrintAirports(airports)
			return nil
		},
	}

	// tap routes [ORIGIN]
	routesCmd := &cobra.Command{
		Use:   "routes [ORIGIN]",
		Short: "List destination airports reachable from an origin airport in TAP network (e.g. LIS, OPO)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			destinations, err := airportsUC.GetFlyTapRoutes(cmd.Context(), args[0])
			if err != nil {
				return fmt.Errorf("failed to get TAP routes: %w", err)
			}

			presenter.PrintAirports(destinations)
			return nil
		},
	}

	// tap calendar
	var (
		calOrigin      string
		calDestination string
		calYear        int
		calMonth       int
		calMarket      string
	)

	now := time.Now()
	calendarCmd := &cobra.Command{
		Use:   "calendar",
		Short: "Search date-by-date low fares between two airports for a specific month",
		RunE: func(cmd *cobra.Command, args []string) error {
			offers, err := datesUC.GetFlyTapCalendar(cmd.Context(), calOrigin, calDestination, calYear, calMonth, calMarket)
			if err != nil {
				return fmt.Errorf("failed to get TAP fare calendar: %w", err)
			}

			presenter.PrintFlightOffers(offers)
			return nil
		},
	}

	calendarCmd.Flags().StringVar(&calOrigin, "origin", "LIS", "Origin airport IATA code")
	calendarCmd.Flags().StringVar(&calDestination, "destination", "BCN", "Destination airport IATA code")
	calendarCmd.Flags().IntVar(&calYear, "year", now.Year(), "Flight year (YYYY)")
	calendarCmd.Flags().IntVar(&calMonth, "month", int(now.Month()), "Flight month (1-12)")
	calendarCmd.Flags().StringVar(&calMarket, "market", "PT", "Market / currency code (PT = EUR, US = USD, GB = GBP)")

	// tap dates [ORIGIN] [DESTINATION]
	datesCmd := &cobra.Command{
		Use:   "dates [ORIGIN] [DESTINATION]",
		Short: "List scheduled flight dates between two airports in upcoming months",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			dates, err := datesUC.GetFlyTapDates(cmd.Context(), args[0], args[1])
			if err != nil {
				return fmt.Errorf("failed to get TAP flight dates: %w", err)
			}

			presenter.PrintDates(args[0], args[1], dates)
			return nil
		},
	}

	flytapCmd.AddCommand(airportsCmd)
	flytapCmd.AddCommand(routesCmd)
	flytapCmd.AddCommand(calendarCmd)
	flytapCmd.AddCommand(datesCmd)

	return flytapCmd
}

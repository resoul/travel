package cli

import (
	"fmt"

	"github.com/resoul/travel/internal/usecase"
	"github.com/spf13/cobra"
)

func newIndiGoCmd(
	datesUC *usecase.FlightDatesUseCase,
	presenter *Presenter,
) *cobra.Command {
	indigoCmd := &cobra.Command{
		Use:   "indigo",
		Short: "IndiGo (6E) flight search, fare calendar, and fare radar commands",
	}

	// indigo radar [ORIGIN]
	radarCmd := &cobra.Command{
		Use:   "radar [ORIGIN]",
		Short: "Get lowest fare recommendations and destinations from an airport (e.g. DEL, BOM, BLR, HYD)",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			origin := "DEL"
			if len(args) > 0 && args[0] != "" {
				origin = args[0]
			}

			fares, err := datesUC.GetIndiGoFareRadar(cmd.Context(), origin)
			if err != nil {
				return fmt.Errorf("failed to get fare radar: %w", err)
			}

			presenter.PrintIndiGoRadar(fares)
			return nil
		},
	}

	// indigo calendar
	var (
		calOrigin      string
		calDestination string
		calStartDate   string
		calEndDate     string
		calCurrency    string
	)

	calendarCmd := &cobra.Command{
		Use:   "calendar",
		Short: "Search date-by-date low fares between two airports over a date range",
		RunE: func(cmd *cobra.Command, args []string) error {
			offers, err := datesUC.GetIndiGoFareCalendar(cmd.Context(), calOrigin, calDestination, calStartDate, calEndDate, calCurrency)
			if err != nil {
				return fmt.Errorf("failed to get fare calendar: %w", err)
			}

			presenter.PrintFlightOffers(offers)
			return nil
		},
	}

	calendarCmd.Flags().StringVar(&calOrigin, "origin", "DXB", "Origin airport IATA code")
	calendarCmd.Flags().StringVar(&calDestination, "destination", "DEL", "Destination airport IATA code")
	calendarCmd.Flags().StringVar(&calStartDate, "start-date", "", "Start date (YYYY-MM-DD, defaults to today)")
	calendarCmd.Flags().StringVar(&calEndDate, "end-date", "", "End date (YYYY-MM-DD, defaults to start-date + 30 days)")
	calendarCmd.Flags().StringVar(&calCurrency, "currency", "INR", "Currency code (e.g. INR, AED, EUR, USD)")

	// indigo dates [ORIGIN] [DESTINATION]
	datesCmd := &cobra.Command{
		Use:   "dates [ORIGIN] [DESTINATION]",
		Short: "List scheduled flight dates between two airports",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			dates, err := datesUC.GetIndiGoDates(cmd.Context(), args[0], args[1])
			if err != nil {
				return fmt.Errorf("failed to get flight dates: %w", err)
			}

			presenter.PrintDates(args[0], args[1], dates)
			return nil
		},
	}

	indigoCmd.AddCommand(radarCmd)
	indigoCmd.AddCommand(calendarCmd)
	indigoCmd.AddCommand(datesCmd)

	return indigoCmd
}

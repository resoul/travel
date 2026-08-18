package cli

import (
	"time"

	"github.com/resoul/travel/internal/domain"
	"github.com/resoul/travel/internal/usecase"
	"github.com/spf13/cobra"
)

func newVuelingCmd(
	searchUC *usecase.SearchFlightsUseCase,
	airportsUC *usecase.ListAirportsUseCase,
	datesUC *usecase.FlightDatesUseCase,
	presenter *Presenter,
) *cobra.Command {
	vuelingCmd := &cobra.Command{
		Use:   "vueling",
		Short: "Vueling flight search and lookup commands",
	}

	// vueling search
	var (
		origin      string
		destination string
		dateOut     string
		adults      int
		flexBefore  int
		flexAfter   int
	)

	searchCmd := &cobra.Command{
		Use:   "search",
		Short: "Search available Vueling fares",
		RunE: func(cmd *cobra.Command, args []string) error {
			criteria := domain.FlightSearchCriteria{
				Origin:         origin,
				Destination:    destination,
				DepartureDate:  dateOut,
				Adults:         adults,
				FlexDaysBefore: flexBefore,
				FlexDaysAfter:  flexAfter,
			}

			results, err := searchUC.SearchVueling(cmd.Context(), criteria)
			if err != nil {
				return err
			}

			presenter.PrintFlightOffers(results)
			return nil
		},
	}

	searchCmd.Flags().StringVar(&origin, "origin", "BCN", "Origin airport IATA code")
	searchCmd.Flags().StringVar(&destination, "destination", "FCO", "Destination airport IATA code")
	searchCmd.Flags().StringVar(&dateOut, "date", "", "Departure date (YYYY-MM-DD, optional)")
	searchCmd.Flags().IntVar(&adults, "adults", 1, "Number of adults")
	searchCmd.Flags().IntVar(&flexBefore, "flex-before", 0, "Flex days before")
	searchCmd.Flags().IntVar(&flexAfter, "flex-after", 0, "Flex days after")

	// vueling airports
	airportsCmd := &cobra.Command{
		Use:   "airports",
		Short: "List all active Vueling airports",
		RunE: func(cmd *cobra.Command, args []string) error {
			airports, err := airportsUC.GetVuelingAirports(cmd.Context())
			if err != nil {
				return err
			}

			presenter.PrintAirports(airports)
			return nil
		},
	}

	// vueling routes [IATA]
	routesCmd := &cobra.Command{
		Use:   "routes [IATA]",
		Short: "List destinations reachable from a given airport in Vueling network",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			destinations, err := airportsUC.GetVuelingRoutes(cmd.Context(), args[0])
			if err != nil {
				return err
			}

			if len(destinations) == 0 {
				return nil
			}

			presenter.PrintAirports(destinations)
			return nil
		},
	}

	// vueling schedule [ORIGIN] [DESTINATION]
	var (
		scheduleYear   int
		scheduleMonth  int
		scheduleMonths int
	)

	scheduleCmd := &cobra.Command{
		Use:   "schedule [ORIGIN] [DESTINATION]",
		Short: "List full flight schedule and fares between two airports",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			offers, err := datesUC.GetVuelingSchedule(cmd.Context(), args[0], args[1], scheduleYear, scheduleMonth, scheduleMonths)
			if err != nil {
				return err
			}

			presenter.PrintFlightOffers(offers)
			return nil
		},
	}

	now := time.Now()
	scheduleCmd.Flags().IntVar(&scheduleYear, "year", now.Year(), "Start year")
	scheduleCmd.Flags().IntVar(&scheduleMonth, "month", int(now.Month()), "Start month (1-12)")
	scheduleCmd.Flags().IntVar(&scheduleMonths, "range", 12, "Number of months range")

	// vueling dates [ORIGIN] [DESTINATION]
	var (
		datesYear   int
		datesMonth  int
		datesMonths int
	)

	datesCmd := &cobra.Command{
		Use:   "dates [ORIGIN] [DESTINATION]",
		Short: "List scheduled flight dates between two airports",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			dates, err := datesUC.GetVuelingDates(cmd.Context(), args[0], args[1], datesYear, datesMonth, datesMonths)
			if err != nil {
				return err
			}

			presenter.PrintDates(args[0], args[1], dates)
			return nil
		},
	}

	datesCmd.Flags().IntVar(&datesYear, "year", now.Year(), "Start year")
	datesCmd.Flags().IntVar(&datesMonth, "month", int(now.Month()), "Start month (1-12)")
	datesCmd.Flags().IntVar(&datesMonths, "range", 12, "Number of months range")

	vuelingCmd.AddCommand(searchCmd)
	vuelingCmd.AddCommand(airportsCmd)
	vuelingCmd.AddCommand(routesCmd)
	vuelingCmd.AddCommand(scheduleCmd)
	vuelingCmd.AddCommand(datesCmd)

	return vuelingCmd
}

package cli

import (
	"github.com/resoul/travel/internal/domain"
	"github.com/resoul/travel/internal/usecase"
	"github.com/spf13/cobra"
)

func newFlyOneCmd(
	searchUC *usecase.SearchFlightsUseCase,
	airportsUC *usecase.ListAirportsUseCase,
	datesUC *usecase.FlightDatesUseCase,
	presenter *Presenter,
) *cobra.Command {
	flyoneCmd := &cobra.Command{
		Use:   "flyone",
		Short: "FlyOne flight search and lookup commands",
	}

	// flyone search
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
		Short: "Search available FlyOne flight offers and fares",
		RunE: func(cmd *cobra.Command, args []string) error {
			criteria := domain.FlightSearchCriteria{
				Origin:         origin,
				Destination:    destination,
				DepartureDate:  dateOut,
				Adults:         adults,
				FlexDaysBefore: flexBefore,
				FlexDaysAfter:  flexAfter,
			}

			results, err := searchUC.SearchFlyOne(cmd.Context(), criteria)
			if err != nil {
				return err
			}

			presenter.PrintFlightOffers(results)
			return nil
		},
	}

	searchCmd.Flags().StringVar(&origin, "origin", "OTP", "Origin airport IATA code")
	searchCmd.Flags().StringVar(&destination, "destination", "BRU", "Destination airport IATA code")
	searchCmd.Flags().StringVar(&dateOut, "date", "", "Departure date (YYYY-MM-DD, optional)")
	searchCmd.Flags().IntVar(&adults, "adults", 1, "Number of adults")
	searchCmd.Flags().IntVar(&flexBefore, "flex-before", 0, "Flex days before")
	searchCmd.Flags().IntVar(&flexAfter, "flex-after", 0, "Flex days after")

	// flyone airports
	airportsCmd := &cobra.Command{
		Use:   "airports",
		Short: "List all departure airports in the FlyOne network",
		RunE: func(cmd *cobra.Command, args []string) error {
			airports, err := airportsUC.GetFlyOneAirports(cmd.Context())
			if err != nil {
				return err
			}

			presenter.PrintAirports(airports)
			return nil
		},
	}

	// flyone routes [IATA]
	routesCmd := &cobra.Command{
		Use:   "routes [IATA]",
		Short: "List destinations reachable from a given airport in the FlyOne network",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			destinations, err := airportsUC.GetFlyOneRoutes(cmd.Context(), args[0])
			if err != nil {
				return err
			}

			presenter.PrintAirports(destinations)
			return nil
		},
	}

	// flyone dates [ORIGIN] [DESTINATION]
	datesCmd := &cobra.Command{
		Use:   "dates [ORIGIN] [DESTINATION]",
		Short: "List scheduled flight dates between two airports in the FlyOne network",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			dates, err := datesUC.GetFlyOneDates(cmd.Context(), args[0], args[1])
			if err != nil {
				return err
			}

			presenter.PrintDates(args[0], args[1], dates)
			return nil
		},
	}

	flyoneCmd.AddCommand(searchCmd)
	flyoneCmd.AddCommand(airportsCmd)
	flyoneCmd.AddCommand(routesCmd)
	flyoneCmd.AddCommand(datesCmd)

	return flyoneCmd
}

package cli

import (
	"github.com/resoul/travel/internal/domain"
	"github.com/resoul/travel/internal/usecase"
	"github.com/spf13/cobra"
)

func newRyanairCmd(
	searchUC *usecase.SearchFlightsUseCase,
	airportsUC *usecase.ListAirportsUseCase,
	datesUC *usecase.FlightDatesUseCase,
	presenter *Presenter,
) *cobra.Command {
	ryanairCmd := &cobra.Command{
		Use:   "ryanair",
		Short: "Ryanair flight search and lookup commands",
	}

	// ryanair search
	var (
		origin      string
		destination string
		dateOut     string
		dateIn      string
		adults      int
		teens       int
		children    int
		infants     int
		roundTrip   bool
		flexBefore  int
		flexAfter   int
	)

	searchCmd := &cobra.Command{
		Use:   "search",
		Short: "Search available Ryanair fares",
		RunE: func(cmd *cobra.Command, args []string) error {
			criteria := domain.FlightSearchCriteria{
				Origin:         origin,
				Destination:    destination,
				DepartureDate:  dateOut,
				ReturnDate:     dateIn,
				Adults:         adults,
				Teens:          teens,
				Children:       children,
				Infants:        infants,
				RoundTrip:      roundTrip,
				FlexDaysBefore: flexBefore,
				FlexDaysAfter:  flexAfter,
			}

			results, err := searchUC.SearchRyanair(cmd.Context(), criteria)
			if err != nil {
				return err
			}

			presenter.PrintFlightOffers(results)
			return nil
		},
	}

	searchCmd.Flags().StringVar(&origin, "origin", "BBU", "Origin airport IATA code")
	searchCmd.Flags().StringVar(&destination, "destination", "GRO", "Destination airport IATA code")
	searchCmd.Flags().StringVar(&dateOut, "date", "2026-06-24", "Departure date (YYYY-MM-DD)")
	searchCmd.Flags().StringVar(&dateIn, "return-date", "", "Return date (YYYY-MM-DD, optional)")
	searchCmd.Flags().IntVar(&adults, "adults", 1, "Number of adults")
	searchCmd.Flags().IntVar(&teens, "teens", 0, "Number of teenagers")
	searchCmd.Flags().IntVar(&children, "children", 0, "Number of children")
	searchCmd.Flags().IntVar(&infants, "infants", 0, "Number of infants")
	searchCmd.Flags().BoolVar(&roundTrip, "roundtrip", false, "Is round trip")
	searchCmd.Flags().IntVar(&flexBefore, "flex-before", 2, "Flex days before")
	searchCmd.Flags().IntVar(&flexAfter, "flex-after", 2, "Flex days after")

	// ryanair airports
	airportsCmd := &cobra.Command{
		Use:   "airports",
		Short: "List all active Ryanair airports",
		RunE: func(cmd *cobra.Command, args []string) error {
			airports, err := airportsUC.GetRyanairAirports(cmd.Context())
			if err != nil {
				return err
			}

			presenter.PrintAirports(airports)
			return nil
		},
	}

	// ryanair routes [IATA]
	routesCmd := &cobra.Command{
		Use:   "routes [IATA]",
		Short: "List airports reachable from a given airport",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			airports, err := airportsUC.GetRyanairRoutes(cmd.Context(), args[0])
			if err != nil {
				return err
			}

			presenter.PrintAirports(airports)
			return nil
		},
	}

	// ryanair dates [ORIGIN] [DESTINATION]
	datesCmd := &cobra.Command{
		Use:   "dates [ORIGIN] [DESTINATION]",
		Short: "List scheduled flight dates between two airports",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			dates, err := datesUC.GetRyanairDates(cmd.Context(), args[0], args[1])
			if err != nil {
				return err
			}

			presenter.PrintDates(args[0], args[1], dates)
			return nil
		},
	}

	ryanairCmd.AddCommand(searchCmd)
	ryanairCmd.AddCommand(airportsCmd)
	ryanairCmd.AddCommand(routesCmd)
	ryanairCmd.AddCommand(datesCmd)

	return ryanairCmd
}

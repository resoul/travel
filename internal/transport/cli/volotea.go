package cli

import (
	"github.com/resoul/travel/internal/domain"
	"github.com/resoul/travel/internal/usecase"
	"github.com/spf13/cobra"
)

func newVoloteaCmd(
	searchUC *usecase.SearchFlightsUseCase,
	airportsUC *usecase.ListAirportsUseCase,
	datesUC *usecase.FlightDatesUseCase,
	presenter *Presenter,
) *cobra.Command {
	voloteaCmd := &cobra.Command{
		Use:   "volotea",
		Short: "Volotea flight search and lookup commands",
	}

	// volotea search
	var (
		origin      string
		destination string
		dateOut     string
		adults      int
		teens       int
		children    int
		infants     int
		flexBefore  int
		flexAfter   int
	)

	searchCmd := &cobra.Command{
		Use:   "search",
		Short: "Search available Volotea fares",
		RunE: func(cmd *cobra.Command, args []string) error {
			criteria := domain.FlightSearchCriteria{
				Origin:         origin,
				Destination:    destination,
				DepartureDate:  dateOut,
				Adults:         adults,
				Teens:          teens,
				Children:       children,
				Infants:        infants,
				FlexDaysBefore: flexBefore,
				FlexDaysAfter:  flexAfter,
			}

			results, err := searchUC.SearchVolotea(cmd.Context(), criteria)
			if err != nil {
				return err
			}

			presenter.PrintFlightOffers(results)
			return nil
		},
	}

	searchCmd.Flags().StringVar(&origin, "origin", "NTE", "Origin airport IATA code")
	searchCmd.Flags().StringVar(&destination, "destination", "VAR", "Destination airport IATA code")
	searchCmd.Flags().StringVar(&dateOut, "date", "", "Departure date (YYYY-MM-DD, optional)")
	searchCmd.Flags().IntVar(&adults, "adults", 1, "Number of adults")
	searchCmd.Flags().IntVar(&teens, "teens", 0, "Number of teenagers")
	searchCmd.Flags().IntVar(&children, "children", 0, "Number of children")
	searchCmd.Flags().IntVar(&infants, "infants", 0, "Number of infants")
	searchCmd.Flags().IntVar(&flexBefore, "flex-before", 0, "Flex days before")
	searchCmd.Flags().IntVar(&flexAfter, "flex-after", 0, "Flex days after")

	// volotea airports
	airportsCmd := &cobra.Command{
		Use:   "airports",
		Short: "List all active Volotea airports",
		RunE: func(cmd *cobra.Command, args []string) error {
			airports, err := airportsUC.GetVoloteaAirports(cmd.Context())
			if err != nil {
				return err
			}

			presenter.PrintAirports(airports)
			return nil
		},
	}

	// volotea routes [IATA]
	routesCmd := &cobra.Command{
		Use:   "routes [IATA]",
		Short: "List airports reachable from a given Volotea airport",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			airports, err := airportsUC.GetVoloteaRoutes(cmd.Context(), args[0])
			if err != nil {
				return err
			}

			presenter.PrintAirports(airports)
			return nil
		},
	}

	// volotea schedule [ORIGIN] [DESTINATION]
	scheduleCmd := &cobra.Command{
		Use:   "schedule [ORIGIN] [DESTINATION]",
		Short: "List full flight schedule and fares between two airports",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			offers, err := datesUC.GetVoloteaSchedule(cmd.Context(), args[0], args[1])
			if err != nil {
				return err
			}

			presenter.PrintFlightOffers(offers)
			return nil
		},
	}

	// volotea dates [ORIGIN] [DESTINATION]
	datesCmd := &cobra.Command{
		Use:   "dates [ORIGIN] [DESTINATION]",
		Short: "List scheduled flight dates between two airports",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			dates, err := datesUC.GetVoloteaDates(cmd.Context(), args[0], args[1])
			if err != nil {
				return err
			}

			presenter.PrintDates(args[0], args[1], dates)
			return nil
		},
	}

	// volotea countries
	countriesCmd := &cobra.Command{
		Use:   "countries",
		Short: "List all countries from Volotea",
		RunE: func(cmd *cobra.Command, args []string) error {
			countries, err := airportsUC.GetVoloteaCountries(cmd.Context())
			if err != nil {
				return err
			}

			presenter.PrintCountries(countries)
			return nil
		},
	}

	voloteaCmd.AddCommand(searchCmd)
	voloteaCmd.AddCommand(airportsCmd)
	voloteaCmd.AddCommand(routesCmd)
	voloteaCmd.AddCommand(scheduleCmd)
	voloteaCmd.AddCommand(datesCmd)
	voloteaCmd.AddCommand(countriesCmd)

	return voloteaCmd
}

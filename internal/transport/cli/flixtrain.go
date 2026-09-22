package cli

import (
	"fmt"
	"time"

	"github.com/resoul/travel/internal/domain"
	"github.com/resoul/travel/internal/usecase"
	"github.com/spf13/cobra"
)

func newFlixTrainCmd(
	searchUC *usecase.SearchFlightsUseCase,
	airportsUC *usecase.ListAirportsUseCase,
	presenter *Presenter,
) *cobra.Command {
	flixtrainCmd := &cobra.Command{
		Use:   "flixtrain",
		Short: "FlixTrain European railway trip search and routes lookup",
	}

	// flixtrain search
	var (
		fromCity string
		toCity   string
		date     string
		adults   int
	)

	searchCmd := &cobra.Command{
		Use:   "search",
		Short: "Search available FlixTrain trips with final prices (including platform fee)",
		RunE: func(cmd *cobra.Command, args []string) error {
			if fromCity == "" || toCity == "" {
				return fmt.Errorf("both --from and --to city parameters are required")
			}

			criteria := domain.FlightSearchCriteria{
				Origin:        fromCity,
				Destination:   toCity,
				DepartureDate: date,
				Adults:        adults,
			}

			results, err := searchUC.SearchFlixTrain(cmd.Context(), criteria)
			if err != nil {
				return err
			}

			presenter.PrintFlightOffers(results)
			return nil
		},
	}

	searchCmd.Flags().StringVar(&fromCity, "from", "", "Departure city name or UUID (e.g. Berlin)")
	searchCmd.Flags().StringVar(&toCity, "to", "", "Arrival city name or UUID (e.g. Paris)")
	searchCmd.Flags().StringVar(&date, "date", time.Now().Format("2006-01-02"), "Departure date (YYYY-MM-DD or DD.MM.YYYY)")
	searchCmd.Flags().IntVar(&adults, "adults", 1, "Number of adult passengers")
	_ = searchCmd.MarkFlagRequired("from")
	_ = searchCmd.MarkFlagRequired("to")

	// flixtrain cities [QUERY]
	citiesCmd := &cobra.Command{
		Use:   "cities [QUERY]",
		Short: "Search FlixTrain cities and stations by name query",
		RunE: func(cmd *cobra.Command, args []string) error {
			query := "Berlin"
			if len(args) > 0 {
				query = args[0]
			}

			cities, err := airportsUC.GetFlixTrainCities(cmd.Context(), query)
			if err != nil {
				return err
			}

			presenter.PrintAirports(cities)
			return nil
		},
	}

	// flixtrain routes [CITY]
	var limit int
	routesCmd := &cobra.Command{
		Use:   "routes [CITY]",
		Short: "List reachable destinations and starting prices from a city",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			destinations, err := airportsUC.GetFlixTrainRoutes(cmd.Context(), args[0], limit)
			if err != nil {
				return err
			}

			presenter.PrintAirports(destinations)
			return nil
		},
	}

	routesCmd.Flags().IntVar(&limit, "limit", 30, "Maximum number of reachable destinations")

	flixtrainCmd.AddCommand(searchCmd)
	flixtrainCmd.AddCommand(citiesCmd)
	flixtrainCmd.AddCommand(routesCmd)

	return flixtrainCmd
}

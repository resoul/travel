package cli

import (
	"github.com/resoul/travel/internal/domain"
	"github.com/resoul/travel/internal/usecase"
	"github.com/spf13/cobra"
)

func newWizzairCmd(
	searchUC *usecase.SearchFlightsUseCase,
	airportsUC *usecase.ListAirportsUseCase,
	presenter *Presenter,
) *cobra.Command {
	wizzairCmd := &cobra.Command{
		Use:   "wizzair",
		Short: "Wizzair flight search and lookup commands",
	}

	var (
		origin      string
		destination string
		fromDate    string
		toDate      string
		adults      int
		children    int
		infants     int
		priceType   string
	)

	searchCmd := &cobra.Command{
		Use:   "search",
		Short: "Search available Wizzair fares",
		RunE: func(cmd *cobra.Command, args []string) error {
			criteria := domain.WizzairSearchCriteria{
				DepartureStation: origin,
				ArrivalStation:   destination,
				FromDate:         fromDate,
				ToDate:           toDate,
				Adults:           adults,
				Children:         children,
				Infants:          infants,
				PriceType:        priceType,
			}

			results, err := searchUC.SearchWizzair(cmd.Context(), criteria)
			if err != nil {
				return err
			}

			presenter.PrintFlightOffers(results)
			return nil
		},
	}

	searchCmd.Flags().StringVar(&origin, "origin", "OTP", "Departure airport IATA code")
	searchCmd.Flags().StringVar(&destination, "destination", "CRL", "Arrival airport IATA code")
	searchCmd.Flags().StringVar(&fromDate, "from", "2026-06-01", "Start date (YYYY-MM-DD)")
	searchCmd.Flags().StringVar(&toDate, "to", "2026-07-05", "End date (YYYY-MM-DD)")
	searchCmd.Flags().IntVar(&adults, "adults", 1, "Number of adults")
	searchCmd.Flags().IntVar(&children, "children", 0, "Number of children")
	searchCmd.Flags().IntVar(&infants, "infants", 0, "Number of infants")
	searchCmd.Flags().StringVar(&priceType, "price-type", "regular", "Price type (e.g. regular, wdc)")

	mapCmd := &cobra.Command{
		Use:   "map",
		Short: "List all connected cities from Wizzair map",
		RunE: func(cmd *cobra.Command, args []string) error {
			cities, err := airportsUC.GetWizzairMap(cmd.Context())
			if err != nil {
				return err
			}

			presenter.PrintCities(cities)
			return nil
		},
	}

	wizzairCmd.AddCommand(searchCmd)
	wizzairCmd.AddCommand(mapCmd)

	return wizzairCmd
}

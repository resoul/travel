package cli

import (
	"fmt"

	"github.com/resoul/travel/internal/domain"
	"github.com/resoul/travel/internal/usecase"
	"github.com/spf13/cobra"
)

func newImoovaCmd(
	searchUC *usecase.SearchFlightsUseCase,
	airportsUC *usecase.ListAirportsUseCase,
	presenter *Presenter,
) *cobra.Command {
	imoovaCmd := &cobra.Command{
		Use:   "imoova",
		Short: "imoova 1-dollar campervan and motorhome relocation search commands",
	}

	// imoova search / offers
	var (
		from string
		to   string
		date string
	)

	searchCmd := &cobra.Command{
		Use:     "search",
		Aliases: []string{"offers", "deals"},
		Short:   "Search available 1-dollar/day campervan relocation deals",
		RunE: func(cmd *cobra.Command, args []string) error {
			criteria := domain.FlightSearchCriteria{
				Origin:        from,
				Destination:   to,
				DepartureDate: date,
			}

			results, err := searchUC.SearchImoova(cmd.Context(), criteria)
			if err != nil {
				return err
			}

			presenter.PrintFlightOffers(results)
			return nil
		},
	}

	searchCmd.Flags().StringVar(&from, "from", "", "Departure city name filter (optional)")
	searchCmd.Flags().StringVar(&to, "to", "", "Delivery city name filter (optional)")
	searchCmd.Flags().StringVar(&date, "date", "", "Departure date (YYYY-MM-DD, optional)")

	// imoova locations
	locationsCmd := &cobra.Command{
		Use:     "locations",
		Aliases: []string{"cities", "routes"},
		Short:   "List all 68 active departure cities with route counts",
		RunE: func(cmd *cobra.Command, args []string) error {
			locations, err := airportsUC.GetImoovaLocations(cmd.Context())
			if err != nil {
				return err
			}

			presenter.PrintAirports(locations)
			return nil
		},
	}

	// imoova regions
	regionsCmd := &cobra.Command{
		Use:     "regions",
		Aliases: []string{"countries"},
		Short:   "List active campervan relocation deals per region (US, CA, EU, AU, NZ, SA)",
		RunE: func(cmd *cobra.Command, args []string) error {
			regions, err := airportsUC.GetImoovaRegions(cmd.Context())
			if err != nil {
				return err
			}

			fmt.Println("imoova Relocations by Region:")
			for _, r := range regions {
				fmt.Printf("  [%s] %s\n", r.Code, r.Name)
			}
			return nil
		},
	}

	imoovaCmd.AddCommand(searchCmd)
	imoovaCmd.AddCommand(locationsCmd)
	imoovaCmd.AddCommand(regionsCmd)

	return imoovaCmd
}

package cli

import (
	"github.com/resoul/travel/internal/domain"
	"github.com/resoul/travel/internal/usecase"
	"github.com/spf13/cobra"
)

func newMovacarCmd(
	searchUC *usecase.SearchFlightsUseCase,
	airportsUC *usecase.ListAirportsUseCase,
	presenter *Presenter,
) *cobra.Command {
	movacarCmd := &cobra.Command{
		Use:   "movacar",
		Short: "Movacar 1-euro car and campervan relocation search commands",
	}

	// movacar search / offers
	var (
		from string
		to   string
		date string
	)

	searchCmd := &cobra.Command{
		Use:     "search",
		Aliases: []string{"offers"},
		Short:   "Search available 1-euro car and campervan relocation offers",
		RunE: func(cmd *cobra.Command, args []string) error {
			criteria := domain.FlightSearchCriteria{
				Origin:        from,
				Destination:   to,
				DepartureDate: date,
			}

			results, err := searchUC.SearchMovacar(cmd.Context(), criteria)
			if err != nil {
				return err
			}

			presenter.PrintFlightOffers(results)
			return nil
		},
	}

	searchCmd.Flags().StringVar(&from, "from", "", "Pickup city or station name filter (optional)")
	searchCmd.Flags().StringVar(&to, "to", "", "Dropoff city or station name filter (optional)")
	searchCmd.Flags().StringVar(&date, "date", "", "Pickup date (YYYY-MM-DD, optional)")

	// movacar locations
	locationsCmd := &cobra.Command{
		Use:     "locations",
		Aliases: []string{"cities", "stations"},
		Short:   "List all active cities and stations with available vehicles",
		RunE: func(cmd *cobra.Command, args []string) error {
			locations, err := airportsUC.GetMovacarLocations(cmd.Context())
			if err != nil {
				return err
			}

			presenter.PrintAirports(locations)
			return nil
		},
	}

	movacarCmd.AddCommand(searchCmd)
	movacarCmd.AddCommand(locationsCmd)

	return movacarCmd
}

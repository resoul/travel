package cli

import (
	"fmt"
	"time"

	"github.com/resoul/travel/internal/domain"
	"github.com/resoul/travel/internal/usecase"
	"github.com/spf13/cobra"
)

func newAlsaCmd(
	searchUC *usecase.SearchFlightsUseCase,
	airportsUC *usecase.ListAirportsUseCase,
	presenter *Presenter,
) *cobra.Command {
	alsaCmd := &cobra.Command{
		Use:   "alsa",
		Short: "ALSA Spanish bus network trip search and stations lookup",
	}

	// alsa search
	var (
		fromStationID string
		toStationID   string
		date          string
		seats         int
	)

	searchCmd := &cobra.Command{
		Use:   "search",
		Short: "Search available ALSA bus trips",
		RunE: func(cmd *cobra.Command, args []string) error {
			if fromStationID == "" || toStationID == "" {
				return fmt.Errorf("both --from and --to station IDs are required")
			}

			criteria := domain.FlightSearchCriteria{
				Origin:        fromStationID,
				Destination:   toStationID,
				DepartureDate: date,
				Adults:        seats,
			}

			results, err := searchUC.SearchAlsa(cmd.Context(), criteria)
			if err != nil {
				return err
			}

			presenter.PrintFlightOffers(results)
			return nil
		},
	}

	searchCmd.Flags().StringVar(&fromStationID, "from", "", "Departure station ID (e.g. 90595 for Barcelona)")
	searchCmd.Flags().StringVar(&toStationID, "to", "", "Arrival station ID (e.g. 90155 for Madrid)")
	searchCmd.Flags().StringVar(&date, "date", time.Now().Format("2006-01-02"), "Departure date (YYYY-MM-DD)")
	searchCmd.Flags().IntVar(&seats, "seats", 1, "Number of passengers")
	_ = searchCmd.MarkFlagRequired("from")
	_ = searchCmd.MarkFlagRequired("to")

	// alsa origins
	originsCmd := &cobra.Command{
		Use:   "origins",
		Short: "List all ALSA departure stations",
		RunE: func(cmd *cobra.Command, args []string) error {
			origins, err := airportsUC.GetAlsaOrigins(cmd.Context())
			if err != nil {
				return err
			}

			presenter.PrintAirports(origins)
			return nil
		},
	}

	// alsa destinations [FROM_STATION_ID]
	destinationsCmd := &cobra.Command{
		Use:   "destinations [FROM_STATION_ID]",
		Short: "List available destinations from a station",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			destinations, err := airportsUC.GetAlsaDestinations(cmd.Context(), args[0])
			if err != nil {
				return err
			}

			presenter.PrintAirports(destinations)
			return nil
		},
	}

	alsaCmd.AddCommand(searchCmd)
	alsaCmd.AddCommand(originsCmd)
	alsaCmd.AddCommand(destinationsCmd)

	return alsaCmd
}

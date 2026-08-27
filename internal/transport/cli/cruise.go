package cli

import (
	"fmt"
	"time"

	"github.com/resoul/travel/internal/domain"
	"github.com/resoul/travel/internal/usecase"
	"github.com/spf13/cobra"
)

func newCruiseCmd(
	searchUC *usecase.SearchFlightsUseCase,
	presenter *Presenter,
) *cobra.Command {
	cruiseCmd := &cobra.Command{
		Use:   "cruise",
		Short: "Cruise search and line/destination lookup commands (Arrivia / AirAsia Cruises)",
	}

	var (
		destinationID string
		cruiseLineID  string
		month         int
		year          int
		durationMin   int
		durationMax   int
		limit         int
	)

	now := time.Now()

	// cruise search
	searchCmd := &cobra.Command{
		Use:   "search",
		Short: "Search available cruises and starting cabin prices across cruise lines",
		RunE: func(cmd *cobra.Command, args []string) error {
			criteria := domain.CruiseSearchCriteria{
				DestinationID: destinationID,
				CruiseLineID:  cruiseLineID,
				Month:         month,
				Year:          year,
				DurationMin:   durationMin,
				DurationMax:   durationMax,
				Limit:         limit,
			}

			offers, err := searchUC.SearchCruises(cmd.Context(), criteria)
			if err != nil {
				return fmt.Errorf("failed to search cruises: %w", err)
			}

			presenter.PrintFlightOffers(offers)
			return nil
		},
	}

	searchCmd.Flags().StringVar(&destinationID, "destination", "", "Destination ID (see 'travel cruise destinations')")
	searchCmd.Flags().StringVar(&cruiseLineID, "cruise-line", "", "Cruise line ID (see 'travel cruise lines')")
	searchCmd.Flags().IntVar(&month, "month", int(now.Month()), "Sailing month (1-12)")
	searchCmd.Flags().IntVar(&year, "year", now.Year(), "Sailing year (YYYY)")
	searchCmd.Flags().IntVar(&durationMin, "duration-min", 0, "Minimum duration in nights")
	searchCmd.Flags().IntVar(&durationMax, "duration-max", 0, "Maximum duration in nights")
	searchCmd.Flags().IntVar(&limit, "limit", 25, "Maximum number of results to fetch")

	// cruise lines
	linesCmd := &cobra.Command{
		Use:   "lines",
		Short: "List all available cruise operators and their matrix IDs",
		RunE: func(cmd *cobra.Command, args []string) error {
			lines, err := searchUC.GetCruiseLines(cmd.Context())
			if err != nil {
				return fmt.Errorf("failed to get cruise lines: %w", err)
			}

			presenter.PrintCruiseLines(lines)
			return nil
		},
	}

	// cruise destinations
	destinationsCmd := &cobra.Command{
		Use:   "destinations",
		Short: "List all available cruise destination regions and their matrix IDs",
		RunE: func(cmd *cobra.Command, args []string) error {
			destinations, err := searchUC.GetCruiseDestinations(cmd.Context())
			if err != nil {
				return fmt.Errorf("failed to get cruise destinations: %w", err)
			}

			presenter.PrintCruiseDestinations(destinations)
			return nil
		},
	}

	cruiseCmd.AddCommand(searchCmd)
	cruiseCmd.AddCommand(linesCmd)
	cruiseCmd.AddCommand(destinationsCmd)

	return cruiseCmd
}

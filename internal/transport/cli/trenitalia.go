package cli

import (
	"fmt"
	"strings"
	"time"

	"github.com/resoul/travel/internal/domain"
	"github.com/resoul/travel/internal/usecase"
	"github.com/spf13/cobra"
)

func newTrenitaliaCmd(searchUC *usecase.SearchFlightsUseCase, airportsUC *usecase.ListAirportsUseCase, presenter *Presenter) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "trenitalia",
		Short: "Trenitalia (Le Frecce) Italian high-speed rail, Frecciarossa, and Intercity search",
	}

	cmd.AddCommand(newTrenitaliaSearchCmd(searchUC, presenter))
	cmd.AddCommand(newTrenitaliaStationsCmd(airportsUC))

	return cmd
}

func newTrenitaliaSearchCmd(searchUC *usecase.SearchFlightsUseCase, presenter *Presenter) *cobra.Command {
	var (
		origin        string
		destination   string
		originID      int
		destinationID int
		departureDate string
		departureTime string
		adults        int
		children      int
		frecceOnly    bool
		regionalOnly  bool
		noChanges     bool
		limit         int
	)

	cmd := &cobra.Command{
		Use:   "search",
		Short: "Search live Trenitalia & Le Frecce train tickets and schedules",
		Long: `Search train journeys across Italy (Frecciarossa, Frecciargento, Frecciabianca, Intercity, Regionale).

Examples:
  travel trenitalia search --origin "Roma Termini" --destination "Milano Centrale" --date 2026-09-10 --time 08:00 --frecce-only
  travel trenitalia search --origin "Firenze SMN" --destination "Venezia Santa Lucia" --date 2026-09-15 --no-changes
  travel trenitalia search --origin "Napoli Centrale" --destination "Torino Porta Nuova" --date 2026-09-20
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			if departureDate == "" {
				departureDate = time.Now().AddDate(0, 0, 14).Format("2006-01-02")
			}

			criteria := domain.TrenitaliaSearchCriteria{
				OriginID:        originID,
				DestinationID:   destinationID,
				OriginName:      origin,
				DestinationName: destination,
				DepartureDate:   departureDate,
				DepartureTime:   departureTime,
				Adults:          adults,
				Children:        children,
				FrecceOnly:      frecceOnly,
				RegionalOnly:    regionalOnly,
				NoChanges:       noChanges,
				Limit:           limit,
			}

			fmt.Printf("🚆 Searching Trenitalia trains from %s to %s on %s at %s (%d adults)...\n\n",
				origin, destination, departureDate, departureTime, adults)

			offers, err := searchUC.SearchTrenitaliaTrains(ctx, criteria)
			if err != nil {
				return fmt.Errorf("failed to search Trenitalia trains: %w", err)
			}

			if len(offers) == 0 {
				fmt.Println("No train journeys found matching criteria.")
				return nil
			}

			fmt.Printf("Found %d train options on Trenitalia:\n", len(offers))
			fmt.Println(strings.Repeat("-", 100))

			for i, t := range offers {
				depStr := "N/A"
				if t.DepartureTime != nil {
					depStr = t.DepartureTime.Format("15:04")
				}
				arrStr := "N/A"
				if t.ArrivalTime != nil {
					arrStr = t.ArrivalTime.Format("15:04")
				}

				priceStr := "Price unavailable"
				if t.Price.Amount > 0 {
					priceStr = fmt.Sprintf("%.2f %s", t.Price.Amount, t.Price.Currency)
				}

				fmt.Printf("[%2d] %s\n", i+1, t.FlightNumber)
				fmt.Printf("     📍 Route: %s -> %s\n", t.DepartureStation, t.ArrivalStation)
				fmt.Printf("     🕒 Departure: %s | Arrival: %s | ⏱️ Duration: %s\n", depStr, arrStr, t.Duration)
				fmt.Printf("     💰 Price: %s | Status: %s\n", priceStr, t.Status)
				fmt.Println()
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&origin, "origin", "Roma Termini", "Departure station name")
	cmd.Flags().StringVar(&destination, "destination", "Milano Centrale", "Arrival station name")
	cmd.Flags().IntVar(&originID, "origin-id", 0, "Optional numeric Origin Station ID")
	cmd.Flags().IntVar(&destinationID, "dest-id", 0, "Optional numeric Destination Station ID")
	cmd.Flags().StringVar(&departureDate, "date", "", "Departure date (YYYY-MM-DD)")
	cmd.Flags().StringVar(&departureTime, "time", "08:00", "Departure time (HH:mm)")
	cmd.Flags().IntVar(&adults, "adults", 1, "Number of adult passengers")
	cmd.Flags().IntVar(&children, "children", 0, "Number of children")
	cmd.Flags().BoolVar(&frecceOnly, "frecce-only", false, "Filter only high-speed Freccia trains")
	cmd.Flags().BoolVar(&regionalOnly, "regional-only", false, "Filter only Regional trains")
	cmd.Flags().BoolVar(&noChanges, "no-changes", false, "Filter only direct journeys without transfers")
	cmd.Flags().IntVar(&limit, "limit", 10, "Maximum number of results to display")

	return cmd
}

func newTrenitaliaStationsCmd(airportsUC *usecase.ListAirportsUseCase) *cobra.Command {
	var (
		query      string
		frecceOnly bool
	)

	cmd := &cobra.Command{
		Use:   "stations",
		Short: "List and search Italian railway and bus stations in Trenitalia catalog",
		Example: `  travel trenitalia stations --query roma
  travel trenitalia stations --query milano
  travel trenitalia stations --frecce`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			fmt.Println("🔍 Fetching Trenitalia stations catalog...")

			stations, err := airportsUC.GetTrenitaliaStations(ctx)
			if err != nil {
				return fmt.Errorf("failed to get stations: %w", err)
			}

			qLower := strings.ToLower(strings.TrimSpace(query))
			filtered := make([]domain.TrenitaliaStation, 0)

			for _, s := range stations {
				if frecceOnly && !s.IsFrecce {
					continue
				}
				if qLower != "" {
					nameLower := strings.ToLower(s.Name)
					valLower := strings.ToLower(s.Value)
					if !strings.Contains(nameLower, qLower) && !strings.Contains(valLower, qLower) {
						continue
					}
				}
				filtered = append(filtered, s)
			}

			if len(filtered) == 0 {
				fmt.Println("No stations found matching filters.")
				return nil
			}

			fmt.Printf("Found %d stations:\n", len(filtered))
			fmt.Println(strings.Repeat("-", 90))

			for i, s := range filtered {
				if i >= 40 {
					fmt.Printf("... and %d more stations (use --query to filter).\n", len(filtered)-40)
					break
				}

				frecceTag := ""
				if s.IsFrecce {
					frecceTag = " [🚄 Freccia]"
				}
				if s.IsEurocity {
					frecceTag += " [🇪🇺 Eurocity]"
				}

				idStr := ""
				if s.ID > 0 {
					idStr = fmt.Sprintf(" | ID: %d", s.ID)
				}

				fmt.Printf("[%2d] %s (%s%s)%s\n", i+1, s.Name, s.Value, idStr, frecceTag)
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&query, "query", "q", "", "Search query for station name")
	cmd.Flags().BoolVar(&frecceOnly, "frecce", false, "Show only stations served by Freccia trains")

	return cmd
}

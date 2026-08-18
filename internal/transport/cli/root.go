package cli

import (
	"context"

	"github.com/resoul/travel/internal/infrastructure/cache"
	"github.com/resoul/travel/internal/usecase"
	"github.com/spf13/cobra"
)

// CLI represents the command line application.
type CLI struct {
	rootCmd *cobra.Command
}

// NewCLI builds and wires all CLI commands with their use case dependencies.
func NewCLI(
	searchUC *usecase.SearchFlightsUseCase,
	airportsUC *usecase.ListAirportsUseCase,
	datesUC *usecase.FlightDatesUseCase,
	fileCache *cache.FileCache,
	cacheDir string,
	presenter *Presenter,
) *CLI {
	if presenter == nil {
		presenter = NewPresenter(nil)
	}

	rootCmd := &cobra.Command{
		Use:   "travel",
		Short: "Travel CLI — flight search and airline information tool",
	}

	rootCmd.AddCommand(newRyanairCmd(searchUC, airportsUC, datesUC, presenter))
	rootCmd.AddCommand(newWizzairCmd(searchUC, airportsUC, presenter))
	rootCmd.AddCommand(newVoloteaCmd(searchUC, airportsUC, datesUC, presenter))
	rootCmd.AddCommand(newVuelingCmd(searchUC, airportsUC, datesUC, presenter))
	rootCmd.AddCommand(newFlixBusCmd(searchUC, airportsUC, presenter))
	rootCmd.AddCommand(newAirBalticCmd(searchUC, airportsUC, datesUC, presenter))
	rootCmd.AddCommand(newFlyOneCmd(searchUC, airportsUC, datesUC, presenter))
	rootCmd.AddCommand(newMovacarCmd(searchUC, airportsUC, presenter))
	rootCmd.AddCommand(newImoovaCmd(searchUC, airportsUC, presenter))

	if fileCache != nil {
		rootCmd.AddCommand(newCacheCmd(fileCache, cacheDir))
	}

	return &CLI{
		rootCmd: rootCmd,
	}
}

// Execute runs the CLI application with context.
func (c *CLI) Execute(ctx context.Context) error {
	return c.rootCmd.ExecuteContext(ctx)
}

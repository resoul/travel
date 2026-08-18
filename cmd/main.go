package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/resoul/travel/internal/infrastructure/airbaltic"
	"github.com/resoul/travel/internal/infrastructure/cache"
	"github.com/resoul/travel/internal/infrastructure/flixbus"
	"github.com/resoul/travel/internal/infrastructure/flyone"
	"github.com/resoul/travel/internal/infrastructure/imoova"
	"github.com/resoul/travel/internal/infrastructure/movacar"
	"github.com/resoul/travel/internal/infrastructure/ryanair"
	"github.com/resoul/travel/internal/infrastructure/volotea"
	"github.com/resoul/travel/internal/infrastructure/vueling"
	"github.com/resoul/travel/internal/infrastructure/wizzair"
	"github.com/resoul/travel/internal/transport/cli"
	"github.com/resoul/travel/internal/usecase"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	// 1. Cache Layer (File Cache with 1-hour default TTL)
	cacheDir := ".cache"
	fileCache, err := cache.NewFileCache(cacheDir, 1*time.Hour)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: cache initialization failed: %v\n", err)
	}

	var cachedTransport http.RoundTripper
	if fileCache != nil {
		cachedTransport = cache.NewCachedTransport(http.DefaultTransport, fileCache, 1*time.Hour)
	}

	// 2. Infrastructure Layer (Providers / API clients with HTTP caching)
	ryanairClient := ryanair.NewClient(cachedTransport)
	wizzairClient := wizzair.NewClient(cachedTransport)
	voloteaClient := volotea.NewClient(cachedTransport)
	vuelingClient := vueling.NewClient(cachedTransport)
	flixbusClient := flixbus.NewClient(cachedTransport)
	airbalticClient := airbaltic.NewClient(cachedTransport)
	flyoneClient := flyone.NewClient(cachedTransport)
	movacarClient := movacar.NewClient(cachedTransport)
	imoovaClient := imoova.NewClient(cachedTransport)

	// 3. Use Case Layer (Application Business Logic)
	searchUC := usecase.NewSearchFlightsUseCase(ryanairClient, wizzairClient, voloteaClient, vuelingClient, flixbusClient, airbalticClient, flyoneClient, movacarClient, imoovaClient)
	airportsUC := usecase.NewListAirportsUseCase(ryanairClient, wizzairClient, voloteaClient, vuelingClient, flixbusClient, airbalticClient, flyoneClient, movacarClient, imoovaClient)
	datesUC := usecase.NewFlightDatesUseCase(ryanairClient, voloteaClient, vuelingClient, airbalticClient, flyoneClient)

	// 4. Transport Layer (CLI)
	presenter := cli.NewPresenter(os.Stdout)
	appCLI := cli.NewCLI(searchUC, airportsUC, datesUC, fileCache, cacheDir, presenter)

	// 4. Execution
	if err := appCLI.Execute(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

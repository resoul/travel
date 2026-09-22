package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/resoul/travel/internal/infrastructure/agoda"
	"github.com/resoul/travel/internal/infrastructure/airbaltic"
	"github.com/resoul/travel/internal/infrastructure/alsa"
	"github.com/resoul/travel/internal/infrastructure/cache"
	"github.com/resoul/travel/internal/infrastructure/campspace"
	"github.com/resoul/travel/internal/infrastructure/cruise"
	"github.com/resoul/travel/internal/infrastructure/driiveme"
	"github.com/resoul/travel/internal/infrastructure/eurowings"
	"github.com/resoul/travel/internal/infrastructure/flixbus"
	"github.com/resoul/travel/internal/infrastructure/flixtrain"
	"github.com/resoul/travel/internal/infrastructure/flyone"
	"github.com/resoul/travel/internal/infrastructure/flytap"
	"github.com/resoul/travel/internal/infrastructure/hipcamp"
	"github.com/resoul/travel/internal/infrastructure/imoova"
	"github.com/resoul/travel/internal/infrastructure/indigo"
	"github.com/resoul/travel/internal/infrastructure/itabus"
	"github.com/resoul/travel/internal/infrastructure/movacar"
	"github.com/resoul/travel/internal/infrastructure/norwegian"
	"github.com/resoul/travel/internal/infrastructure/obilet"
	"github.com/resoul/travel/internal/infrastructure/pitchup"
	"github.com/resoul/travel/internal/infrastructure/ryanair"
	"github.com/resoul/travel/internal/infrastructure/sata"
	"github.com/resoul/travel/internal/infrastructure/tictactrip"
	"github.com/resoul/travel/internal/infrastructure/transavia"
	"github.com/resoul/travel/internal/infrastructure/trenitalia"
	"github.com/resoul/travel/internal/infrastructure/trip"
	"github.com/resoul/travel/internal/infrastructure/volotea"
	"github.com/resoul/travel/internal/infrastructure/vueling"
	"github.com/resoul/travel/internal/infrastructure/wizzair"
	"github.com/resoul/travel/internal/transport/cli"
	"github.com/resoul/travel/internal/usecase"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	cacheDir := ".cache"
	fileCache, err := cache.NewFileCache(cacheDir, 1*time.Hour)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: cache initialization failed: %v\n", err)
	}

	var cachedTransport http.RoundTripper
	if fileCache != nil {
		cachedTransport = cache.NewCachedTransport(http.DefaultTransport, fileCache, 1*time.Hour)
	}

	ryanairClient := ryanair.NewClient(cachedTransport)
	wizzairClient := wizzair.NewClient(cachedTransport)
	voloteaClient := volotea.NewClient(cachedTransport)
	vuelingClient := vueling.NewClient(cachedTransport)
	flixbusClient := flixbus.NewClient(cachedTransport)
	flixtrainClient := flixtrain.NewClient(cachedTransport)
	airbalticClient := airbaltic.NewClient(cachedTransport)
	flyoneClient := flyone.NewClient(cachedTransport)
	movacarClient := movacar.NewClient(cachedTransport)
	imoovaClient := imoova.NewClient(cachedTransport)
	driivemeClient := driiveme.NewClient(cachedTransport)
	indigoClient := indigo.NewClient(cachedTransport)
	flytapClient := flytap.NewClient(cachedTransport)
	cruiseClient := cruise.NewClient(cachedTransport)
	agodaClient := agoda.NewClient(cachedTransport)
	tripClient := trip.NewClient(cachedTransport)
	tictactripClient := tictactrip.NewClient(cachedTransport)
	trenitaliaClient := trenitalia.NewClient(cachedTransport)
	norwegianClient := norwegian.NewClient(cachedTransport)
	obiletClient := obilet.NewClient(cachedTransport)
	eurowingsClient := eurowings.NewClient(cachedTransport)
	transaviaClient := transavia.NewClient(cachedTransport)
	pitchupClient := pitchup.NewClient(cachedTransport)
	hipcampClient := hipcamp.NewClient(cachedTransport)
	campspaceClient := campspace.NewClient(cachedTransport)
	sataClient := sata.NewClient(cachedTransport)
	alsaClient := alsa.NewClient(cachedTransport)
	itabusClient := itabus.NewClient(cachedTransport)

	searchUC := usecase.NewSearchFlightsUseCase(ryanairClient, wizzairClient, voloteaClient, vuelingClient, flixbusClient, flixtrainClient, airbalticClient, flyoneClient, movacarClient, imoovaClient, driivemeClient, cruiseClient, agodaClient, tripClient, trenitaliaClient, obiletClient, pitchupClient, hipcampClient, campspaceClient, alsaClient, itabusClient)
	airportsUC := usecase.NewListAirportsUseCase(ryanairClient, wizzairClient, voloteaClient, vuelingClient, flixbusClient, flixtrainClient, airbalticClient, flyoneClient, movacarClient, imoovaClient, driivemeClient, flytapClient, agodaClient, tictactripClient, trenitaliaClient, eurowingsClient, sataClient, alsaClient, itabusClient)
	datesUC := usecase.NewFlightDatesUseCase(ryanairClient, voloteaClient, vuelingClient, airbalticClient, flyoneClient, indigoClient, flytapClient, tictactripClient, norwegianClient, eurowingsClient, transaviaClient, sataClient)

	presenter := cli.NewPresenter(os.Stdout)
	appCLI := cli.NewCLI(searchUC, airportsUC, datesUC, fileCache, cacheDir, presenter)

	if err := appCLI.Execute(ctx); err != nil {
		_, err := fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		if err != nil {
			return
		}
		os.Exit(1)
	}
}

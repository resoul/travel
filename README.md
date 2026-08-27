# Travel CLI ✈️ 🚌 🚗 🚐 🚢

A command-line flight, bus, vehicle relocation, and cruise search tool for low-cost carriers, rental platforms, and cruise operators (**Ryanair**, **Wizzair**, **Volotea**, **Vueling**, **airBaltic**, **FlyOne**, **IndiGo**, **TAP Air Portugal**, **Cruises**, **FlixBus**, **Movacar**, **imoova**, and **DriiveMe**).

The project is built following **Clean Architecture** principles in Go.

---

## 🏛 Architecture

The codebase is organized into independent layers:

- **`internal/domain`** — Pure business entities (`FlightOffer`, `Airport`, `Country`, `FlightSearchCriteria`, `CruiseSearchCriteria`, `CruiseLine`, `CruiseDestination`, `IndiGoRadarFare`) and ports (interfaces `RyanairProvider`, `WizzairProvider`, `VoloteaProvider`, `VuelingProvider`, `FlixBusProvider`, `AirBalticProvider`, `FlyOneProvider`, `IndiGoProvider`, `FlyTapProvider`, `CruiseProvider`, `MovacarProvider`, `ImoovaProvider`, `DriiveMeProvider`). Has zero external dependencies.
- **`internal/usecase`** — Application business logic (`SearchFlightsUseCase`, `ListAirportsUseCase`, `FlightDatesUseCase`).
- **`internal/infrastructure`** — External API adapters (`infrastructure/ryanair`, `infrastructure/wizzair`, `infrastructure/volotea`, `infrastructure/vueling`, `infrastructure/airbaltic`, `infrastructure/flyone`, `infrastructure/indigo`, `infrastructure/flytap`, `infrastructure/cruise`, `infrastructure/movacar`, `infrastructure/imoova`, `infrastructure/driiveme`, `infrastructure/flixbus`), disk-based caching (`infrastructure/cache`), DTOs, and domain mappers.
- **`internal/transport/cli`** — CLI user interface built with Cobra (`root`, `ryanair`, `wizzair`, `volotea`, `vueling`, `airbaltic`, `flyone`, `indigo`, `tap`, `cruise`, `movacar`, `imoova`, `driiveme`, `flixbus`, `cache`, `presenter`).
- **`cmd/main.go`** — Single application entry point: handles Dependency Injection (DI), wires layers, and starts the CLI with graceful shutdown support (`signal.NotifyContext`).



```
travel/
├── cmd/
│   └── main.go                               # Single application entrypoint
├── internal/
│   ├── domain/                               # Domain models & interfaces
│   │   ├── airport.go
│   │   ├── flight.go
│   │   └── provider.go
│   ├── usecase/                              # Application use cases
│   │   ├── search_flights.go
│   │   ├── list_airports.go
│   │   └── flight_dates.go
│   ├── infrastructure/                       # External adapters & cache
│   │   ├── cache/
│   │   ├── ryanair/
│   │   ├── wizzair/
│   │   ├── volotea/
│   │   ├── vueling/
│   │   ├── airbaltic/
│   │   ├── flyone/
│   │   ├── movacar/
│   │   ├── imoova/
│   │   ├── driiveme/
│   │   ├── agoda/
│   │   ├── trip/
│   │   ├── tictactrip/
│   │   ├── trenitalia/
│   │   ├── norwegian/
│   │   ├── obilet/
│   │   ├── eurowings/
│   │   ├── transavia/
│   │   ├── pitchup/
│   │   ├── sata/
│   │   ├── indigo/
│   │   ├── flytap/
│   │   ├── cruise/
│   │   └── flixbus/
│   └── transport/
│       └── cli/                              # Cobra CLI transport
│           ├── root.go
│           ├── ryanair.go
│           ├── wizzair.go
│           ├── volotea.go
│           ├── vueling.go
│           ├── airbaltic.go
│           ├── flyone.go
│           ├── movacar.go
│           ├── imoova.go
│           ├── driiveme.go
│           ├── agoda.go
│           ├── trip.go
│           ├── tictactrip.go
│           ├── trenitalia.go
│           ├── norwegian.go
│           ├── obilet.go
│           ├── eurowings.go
│           ├── transavia.go
│           ├── pitchup.go
│           ├── sata.go
│           ├── indigo.go
│           ├── flytap.go
│           ├── cruise.go
│           ├── flixbus.go
│           ├── cache.go
│           └── presenter.go
├── README.md
├── Makefile
├── go.mod
└── go.sum
```

---

## 🚀 Installation & Getting Started

### Prerequisites
- Go 1.22+ installed

### Run directly with `go run`:
```bash
go run ./cmd --help
```

### Build binary:
```bash
make build
./bin/travel --help
```

### Quick Commands via Makefile:
```bash
# View all available make targets
make help

# Run any travel command
make run ARGS="flyone airports"

# Quick operator searches
make ryanair-search ORIGIN=BBU DEST=GRO DATE=2026-09-02
make volotea-search ORIGIN=NTE DEST=VAR DATE=2026-09-01
make vueling-search ORIGIN=BCN DEST=FCO DATE=2026-08-30
make airbaltic-search ORIGIN=ALC DEST=RIX DATE=2026-08-28
make flyone-search ORIGIN=OTP DEST=BRU DATE=2026-08-28
make flixbus-search FROM="Bucharest" TO="Brasov" DATE=2026-08-28
make movacar-offers FROM=Berlin TO=Reutlingen
make imoova-offers FROM=Vancouver TO="San Francisco"
make indigo-radar ORIGIN=DEL
make indigo-calendar ORIGIN=DXB DEST=DEL CURRENCY=AED
make tap-airports
make tap-calendar ORIGIN=LIS DEST=BCN MONTH=9 YEAR=2026
make cruise-search MONTH=11 YEAR=2026 LIMIT=10

# Cache management
make cache-info
make cache-clear
```


---

## 📖 Command Reference

### 1. Ryanair

```bash
# Display Ryanair help
travel ryanair --help

# 1. Search available fares
travel ryanair search --origin BBU --destination GRO --date 2026-08-22

# Search with custom passenger and date flexibility options:
travel ryanair search \
  --origin BBU \
  --destination GRO \
  --date 2026-08-22 \
  --adults 2 \
  --teens 1 \
  --children 0 \
  --flex-before 3 \
  --flex-after 3

# 2. List all active Ryanair airports
travel ryanair airports

# 3. List reachable destinations from a specific airport
travel ryanair routes BBU

# 4. List scheduled flight dates between two airports
travel ryanair dates BBU GRO
```

#### Flags for `travel ryanair search`:
| Flag | Description | Default |
|------|-------------|---------|
| `--origin` | Origin airport IATA code | `BBU` |
| `--destination` | Destination airport IATA code | `GRO` |
| `--date` | Departure date (`YYYY-MM-DD`) | `2026-06-24` |
| `--return-date` | Return date (for round trips) | `""` |
| `--adults` | Number of adult passengers | `1` |
| `--teens` | Number of teenagers (12-15 yrs) | `0` |
| `--children` | Number of children (2-11 yrs) | `0` |
| `--infants` | Number of infants (< 2 yrs) | `0` |
| `--roundtrip` | Is round trip | `false` |
| `--flex-before` | Flex search window before departure date (days) | `2` |
| `--flex-after` | Flex search window after departure date (days) | `2` |

---

### 2. Volotea

```bash
# Display Volotea help
travel volotea --help

# 1. Search flights for a route and date
travel volotea search --origin NTE --destination VAR --date 2026-08-25

# 2. View full flight schedule and fares between two airports
travel volotea schedule NTE VAR

# 3. List scheduled flight dates
travel volotea dates NTE VAR

# 4. List reachable destinations from an airport
travel volotea routes NTE

# 5. List all active Volotea airports
travel volotea airports

# 6. List supported countries and currencies
travel volotea countries
```

#### Flags for `travel volotea search`:
| Flag | Description | Default |
|------|-------------|---------|
| `--origin` | Origin airport IATA code | `NTE` |
| `--destination` | Destination airport IATA code | `VAR` |
| `--date` | Departure date (`YYYY-MM-DD`) | `""` (all available) |
| `--adults` | Number of adult passengers | `1` |
| `--flex-before` | Flex days before target date | `0` |
| `--flex-after` | Flex days after target date | `0` |

---

### 3. Vueling

```bash
# Display Vueling help
travel vueling --help

# 1. Search flights for a route and date
travel vueling search --origin BCN --destination FCO --date 2026-08-19

# 2. View flight schedule and fares across months
travel vueling schedule BCN FCO --range 6

# 3. List scheduled flight dates between two airports
travel vueling dates BCN FCO

# 4. List reachable destinations from an airport in Vueling network
travel vueling routes BCN

# 5. List all 129 active Vueling airports
travel vueling airports
```

#### Flags for `travel vueling search`:
| Flag | Description | Default |
|------|-------------|---------|
| `--origin` | Origin airport IATA code | `BCN` |
| `--destination` | Destination airport IATA code | `FCO` |
| `--date` | Departure date (`YYYY-MM-DD`) | `""` (all available) |
| `--adults` | Number of adult passengers | `1` |
| `--flex-before` | Flex days before target date | `0` |
| `--flex-after` | Flex days after target date | `0` |

---

### 4. airBaltic

```bash
# Display airBaltic help
travel airbaltic --help

# 1. Search flights for a route and date
travel airbaltic search --origin ALC --destination RIX --date 2026-08-28

# 2. List scheduled flight dates and fares between two airports
travel airbaltic dates ALC RIX

# 3. List reachable destinations from an airport in airBaltic network
travel airbaltic routes RIX

# 4. List all 85 primary origin airports
travel airbaltic airports
```

#### Flags for `travel airbaltic search`:
| Flag | Description | Default |
|------|-------------|---------|
| `--origin` | Origin airport IATA code | `ALC` |
| `--destination` | Destination airport IATA code | `RIX` |
| `--date` | Departure date (`YYYY-MM-DD`) | `""` (all available) |
| `--adults` | Number of adult passengers | `1` |
| `--flex-before` | Flex days before target date | `0` |
| `--flex-after` | Flex days after target date | `0` |

---

### 5. FlyOne

```bash
# Display FlyOne help
travel flyone --help

# 1. Search flights and fares for a route and date
travel flyone search --origin OTP --destination BRU --date 2026-08-28

# 2. List scheduled flight dates between two airports
travel flyone dates OTP BRU

# 3. List reachable destinations from an airport in FlyOne network
travel flyone routes OTP

# 4. List all 77 departure airports in FlyOne network
travel flyone airports
```

#### Flags for `travel flyone search`:
| Flag | Description | Default |
|------|-------------|---------|
| `--origin` | Origin airport IATA code | `OTP` |
| `--destination` | Destination airport IATA code | `BRU` |
| `--date` | Departure date (`YYYY-MM-DD`) | `""` (all available) |
| `--adults` | Number of adult passengers | `1` |
| `--flex-before` | Flex days before target date | `0` |
| `--flex-after` | Flex days after target date | `0` |

---

### 6. FlixBus & FlixTrain 🚌

```bash
# Display FlixBus help
travel flixbus --help

# 1. Search trips between cities (prices include platform fee)
travel flixbus search --from Bucharest --to Brasov --date 2026-08-27

# Search using city names or UUIDs:
travel flixbus search --from Berlin --to Paris --date 2026-08-20 --adults 2

# 2. Search FlixBus cities and stations by name
travel flixbus cities Berlin

# 3. List reachable destinations and starting prices from a city
travel flixbus routes Berlin --limit 20
```

#### Flags for `travel flixbus search`:
| Flag | Description | Default |
|------|-------------|---------|
| `--from` | Departure city name or UUID *(required)* | `""` |
| `--to` | Arrival city name or UUID *(required)* | `""` |
| `--date` | Departure date (`YYYY-MM-DD` or `DD.MM.YYYY`) | Today |
| `--adults` | Number of adult passengers | `1` |

---

### 7. Wizzair

```bash
# Display Wizzair help
travel wizzair --help

# 1. Search Wizzair timetable fares
travel wizzair search --origin OTP --destination CRL --from 2026-06-01 --to 2026-07-05

# 2. Fetch connected cities map
travel wizzair map
```

#### Flags for `travel wizzair search`:
| Flag | Description | Default |
|------|-------------|---------|
| `--origin` | Origin airport IATA code | `OTP` |
| `--destination` | Destination airport IATA code | `CRL` |
| `--from` | Start date (`YYYY-MM-DD`) | `2026-06-01` |
| `--to` | End date (`YYYY-MM-DD`) | `2026-07-05` |
| `--adults` | Number of adult passengers | `1` |
| `--price-type` | Fare type (`regular`, `wdc`) | `regular` |

### 8. Movacar 🚗🚐

```bash
# Display Movacar help
travel movacar --help

# 1. Search all available 1-euro car and campervan relocation offers
travel movacar search

# Filter by departure city or destination:
travel movacar search --from Berlin --to Reutlingen

# 2. List all active cities and stations with available vehicles
travel movacar locations
```

#### Flags for `travel movacar search`:
| Flag | Description | Default |
|------|-------------|---------|
| `--from` | Pickup city or station name filter | `""` |
| `--to` | Dropoff city or station name filter | `""` |
| `--date` | Pickup date (`YYYY-MM-DD`) | `""` |

---

### 9. imoova 🚐🌍

```bash
# Display imoova help
travel imoova --help

# 1. Search all available 1-dollar/day campervan relocation deals
travel imoova search

# Filter by departure city, destination, or date:
travel imoova search --from Vancouver --to "San Francisco"
travel imoova search --from Dallas --to Nashville

# 2. List all 68 active departure cities with route counts
travel imoova locations

# 3. List active campervan relocation deals per region (US, CA, EU, AU, NZ, SA)
travel imoova regions
```

#### Flags for `travel imoova search`:
| Flag | Description | Default |
|------|-------------|---------|
| `--from` | Departure city name filter | `""` |
| `--to` | Delivery city name filter | `""` |
| `--date` | Departure date (`YYYY-MM-DD`) | `""` |

---

### 10. DriiveMe 🚗🇪🇺

```bash
# Display DriiveMe help
travel driiveme --help

# 1. Login to DriiveMe
travel driiveme login --email resoul.developer@gmail.com --password "YOUR_PASSWORD"

# 2. Search 1-euro car & van relocation offers (with optional login for enriched vehicle model details)
travel driiveme search
travel driiveme search --from Barcelona --to Torremolinos
travel driiveme search --from 38262 --email resoul.developer@gmail.com --password "YOUR_PASSWORD"

# 3. Autocomplete cities
travel driiveme cities --query Barcelona
travel driiveme cities --query London

# 4. View available booking slots for a transport
travel driiveme availabilities 5252575
```

#### Flags for `travel driiveme search`:
| Flag | Description | Default |
|------|-------------|---------|
| `--from` | Pickup city name or ID | `""` |
| `--to` | Dropoff city name or ID | `""` |
| `--date` | Pickup min date (`YYYY-MM-DD`) | `""` |
| `-e, --email` | DriiveMe account email | `""` (or `DRIIVEME_EMAIL`) |
| `-p, --password` | DriiveMe account password | `""` (or `DRIIVEME_PASSWORD`) |

---

### 11. IndiGo (6E) ✈️🇮🇳

```bash
# Display IndiGo help
travel indigo --help

# 1. Get lowest fare recommendations and destinations from an airport
travel indigo radar DEL
travel indigo radar BOM
travel indigo radar BLR

# 2. Search date-by-date fare calendar between two airports
travel indigo calendar --origin DXB --destination DEL --currency AED
travel indigo calendar --origin DEL --destination BOM --currency INR --start-date 2026-09-01 --end-date 2026-10-01

# 3. List scheduled flight dates between two airports
travel indigo dates DXB DEL
```

#### Flags for `travel indigo calendar`:
| Flag | Description | Default |
|------|-------------|---------|
| `--origin` | Origin airport IATA code | `DXB` |
| `--destination` | Destination airport IATA code | `DEL` |
| `--start-date` | Start date (`YYYY-MM-DD`) | Today |
| `--end-date` | End date (`YYYY-MM-DD`) | Start date + 30 days |
| `--currency` | Currency code (`INR`, `AED`, `EUR`, `USD`) | `INR` |

---

### 11. TAP Air Portugal (TP) ✈️🇵🇹

```bash
# Display TAP help
travel tap --help

# 1. List all operating departure airports in TAP network
travel tap airports

# 2. List reachable destination airports from an origin (e.g. LIS, OPO)
travel tap routes LIS
travel tap routes OPO

# 3. Search date-by-date low fare calendar for a specific month
travel tap calendar --origin LIS --destination BCN --year 2026 --month 9 --market PT
travel tap calendar --origin OPO --destination MAD --year 2026 --month 10 --market US

# 4. List scheduled flight dates between two airports
travel tap dates LIS BCN
```

#### Flags for `travel tap calendar`:
| Flag | Description | Default |
|------|-------------|---------|
| `--origin` | Origin airport IATA code | `LIS` |
| `--destination` | Destination airport IATA code | `BCN` |
| `--year` | Flight year (`YYYY`) | Current year |
| `--month` | Flight month (`1-12`) | Current month |
| `--market` | Market / currency code (`PT` = EUR, `US` = USD, `GB` = GBP) | `PT` |

---

### 12. Cruises (AirAsia / Arrivia) 🚢🌊

Search ocean and river cruises across major global operators (*Royal Caribbean, Carnival, MSC, Norwegian Cruise Line, Princess Cruises, Holland America, Celebrity, Disney, Costa, AmaWaterways, etc.*).

```bash
# Display cruise help
travel cruise --help

# 1. List all cruise lines and their matrix IDs
travel cruise lines

# 2. List all cruise destination regions and their matrix IDs
travel cruise destinations

# 3. Search cruises across all lines for a specific month
travel cruise search --month 11 --year 2026 --limit 10

# 4. Search cruises filtered by destination and cruise line
# Example: Asia cruises (Destination ID: 54)
travel cruise search --destination 54 --month 11 --year 2026 --limit 10

# Example: Carnival Cruises (Cruise Line ID: 1)
travel cruise search --cruise-line 1 --month 11 --year 2026 --duration-min 3 --duration-max 7
```

#### Flags for `travel cruise search`:
| Flag | Description | Default |
|------|-------------|---------|
| `--destination` | Destination region ID (see `travel cruise destinations`) | None |
| `--cruise-line` | Cruise line ID (see `travel cruise lines`) | None |
| `--month` | Sailing month (`1-12`) | Current month |
| `--year` | Sailing year (`YYYY`) | Current year |
| `--duration-min` | Minimum duration in nights | None |
| `--duration-max` | Maximum duration in nights | None |
| `--limit` | Maximum results to fetch | `25` |

---

### 13. Agoda (Accommodations & Hotels)

Search live hotel and accommodation deals on Agoda with full prices and ratings using Chromedp headless browser automation:

```bash
# 1. Search hotel deals in Mamaia (City ID: 19216)
travel agoda search --city-id 19216 --city Mamaia --check-in 2026-09-06 --check-out 2026-09-08 --adults 2 --currency EUR --limit 10

# 2. Search hotel deals in Thessaloniki (City ID: 17336)
travel agoda search --city-id 17336 --city Thessaloniki --check-in 2026-09-10 --check-out 2026-09-12 --adults 2 --currency EUR

# 3. List full world countries directory from Agoda CDN
travel agoda countries --lang 11  # Russian
travel agoda countries --lang 1   # English
```

#### Flags for `travel agoda search`:
| Flag | Description | Default |
|------|-------------|---------|
| `--city-id` | Agoda City ID (e.g. `19216` for Mamaia, `17336` for Thessaloniki) | `19216` |
| `--city` | City name text query | `Mamaia` |
| `--check-in` | Check-in date (`YYYY-MM-DD`) | `2026-09-06` |
| `--check-out` | Check-out date (`YYYY-MM-DD`) | `2026-09-08` |
| `--rooms` | Number of rooms | `1` |
| `--adults` | Number of adult guests | `2` |
| `--children` | Number of children | `0` |
| `--currency` | Display currency code | `EUR` |
| `--sort` | Sort order (`priceLowToHigh`, `rank`, etc.) | `priceLowToHigh` |
| `--limit` | Maximum results to display | `20` |

---

### 14. Trip.com (Hotels & Rental Cars)

Search live hotel deals via fast SSR HTTP requests, or live rental cars across global car hire agencies (Avis, Hertz, Europcar, Sixt, Klass Wagen) via headless automation:

```bash
# 1. Search hotel deals in Barcelona (City ID: 40795)
travel trip search --city-id 40795 --city Barcelona --check-in 2026-11-19 --check-out 2026-11-23 --adults 2 --currency USD --limit 10

# 2. View full room breakdown, bed configurations, area (m²), and amenities for a hotel
travel trip details 134130947 --check-in 2026-11-19 --check-out 2026-11-23

# 3. Search live rental cars at Bucharest Otopeni Airport (OTP)
travel trip cars --country-id 63 --city-id 39050 --city Otopeni --code OTP --pickup 2026-08-29 --return 2026-09-01 --currency USD --limit 10
```

#### Flags for `travel trip cars`:
| Flag | Description | Default |
|------|-------------|---------|
| `--country-id` | Country ID (e.g. `63` for Romania, `95` for Spain, `75` for France) | `63` |
| `--city-id` | City/Airport ID (e.g. `39050` for Otopeni, `40795` for Barcelona) | `39050` |
| `--city` | City / Airport text query | `Otopeni` |
| `--code` | IATA airport code (e.g. `OTP`, `BCN`, `ATH`) | `OTP` |
| `--pickup` | Pickup date & time (`YYYY-MM-DD HH:mm`) | `2026-08-29 10:00` |
| `--return` | Return date & time (`YYYY-MM-DD HH:mm`) | `2026-09-01 10:00` |
| `--age` | Driver age category | `30-60` |
| `--currency` | Display currency code | `USD` |
| `--limit` | Maximum results to display | `15` |

---

### 15. Tictactrip (European Trains, Buses & Price Calendar)

Explore European train (TGV Inoui, TER, Ouigo, Renfe, Trenitalia) and bus routes (FlixBus, BlaBlaCar Bus) with date-by-date lowest fares:

```bash
# 1. Search and autocomplete European cities and railway/bus stations
travel tictactrip cities --query bar
travel tictactrip cities --query paris

# 2. List top popular travel destinations from a city
travel tictactrip popular --from barcelone --limit 7
travel tictactrip popular --limit 10

# 3. View date-by-date lowest train & bus fares for a whole month (e.g. Barcelona [76] to Montpellier [542] in Nov 2026)
travel tictactrip calendar --origin-id 76 --dest-id 542 --month 2026-11
```

---

### 16. Trenitalia (Italian High-Speed Rail & Le Frecce)

Search live high-speed trains (*Frecciarossa 1000, Frecciargento, Frecciabianca, Intercity, Regionale*) across Italy via Trenitalia's official REST BFF API:

```bash
# 1. Search high-speed trains from Rome to Milan
travel trenitalia search --origin "Roma Termini" --destination "Milano Centrale" --date 2026-09-10 --time 08:00 --frecce-only

# 2. Search direct trains from Florence to Venice
travel trenitalia search --origin "Firenze SMN" --destination "Venezia Santa Lucia" --date 2026-09-15 --no-changes

# 3. Explore and filter the Italian railway stations catalog
travel trenitalia stations --query roma
travel trenitalia stations --query milano
travel trenitalia stations --frecce
```

---

### 17. Norwegian Air Shuttle (Low-Fare Calendar)

Explore Scandinavian and European low-cost flight fares date-by-date across full months with Norwegian Air Shuttle via Chromedp headless extraction:

```bash
# 1. Search low fare calendar from Oslo (OSL) to Barcelona (BCN) in September 2026
travel norwegian calendar --origin OSL --destination BCN --year 2026 --month 9

# 2. Search Stockholm (ARN) to London Gatwick (LGW) in October 2026
travel norwegian calendar --origin ARN --destination LGW --year 2026 --month 10

# 3. Search Copenhagen (CPH) to Rome (FCO) in November 2026 with currency EUR
travel norwegian calendar --origin CPH --destination FCO --year 2026 --month 11 --currency EUR
```

---

### 18. oBilet (Turkish & Regional Intercity Buses)

Search intercity bus tickets across Turkey and neighboring regions (*Metro Turizm, Kamil Koç, Pamukkale, Ali Osman Ulusoy, etc.*) via Chromedp headless extraction:

```bash
# 1. Search bus tickets from Istanbul to Ankara
travel obilet search --origin istanbul --destination ankara --date 2026-09-10

# 2. Search bus tickets from Izmir to Antalya
travel obilet search --origin izmir --destination antalya --date 2026-09-15

# 3. Search bus tickets from Antalya to Cappadocia (Kayseri)
travel obilet search --origin antalya --destination cappadocia --date 2026-09-20
```

---

### 19. Eurowings (Lufthansa Group Route Network & Flight Schedules)

Explore European route maps, direct destinations, and full flight schedules across Eurowings and Eurowings Discover via Chromedp headless extraction:

```bash
# 1. List active airports in Germany on Eurowings
travel eurowings airports --country DE

# 2. List all direct routes available from Bucharest (OTP)
travel eurowings routes --origin OTP

# 3. List all scheduled flight dates for Bucharest (OTP) to Dusseldorf (DUS)
travel eurowings dates --origin OTP --destination DUS
```

---

### 20. Transavia (Air France-KLM Low-Fare Calendar)

Explore Dutch, French, and Mediterranean low-cost flight fares date-by-date across full months with Transavia via Chromedp headless extraction:

```bash
# 1. Search low fare calendar from Amsterdam (AMS) to Barcelona (BCN) in September 2026
travel transavia calendar --origin AMS --destination BCN --year 2026 --month 9

# 2. Search Paris Orly (ORY) to Madrid (MAD) in October 2026
travel transavia calendar --origin ORY --destination MAD --year 2026 --month 10

# 3. Search Rotterdam (RTM) to Alicante (ALC) in November 2026
travel transavia calendar --origin RTM --destination ALC --year 2026 --month 11
```

---

### 21. Pitchup (Campsites, Glamping & Outdoor Stays)

Search outdoor holiday accommodations, tent pitches, campervan pitches, and luxury glamping across 67+ countries (*France, England, Spain, Italy, Germany, etc.*) via Chromedp headless extraction:

```bash
# 1. Search campsites in France for 2 adults
travel pitchup search --country france --arrive 2026-09-10 --depart 2026-09-12

# 2. Search campsites in England
travel pitchup search --country england --arrive 2026-09-15 --depart 2026-09-18

# 3. Search campsites & glamping in Italy for 2 adults
travel pitchup search --country italy --arrive 2026-10-01 --depart 2026-10-05 --adults 2
```

---

### 22. Hipcamp (Worldwide Glamping & Outdoor Stays)

Search outdoor accommodations, private land camping, yurts, treehouses, and glamping across the US, UK, Canada, Australia, and Europe via Chromedp headless extraction:

```bash
# 1. Search outdoor spots in California, US
travel hipcamp search --country united-states --region california

# 2. Search outdoor spots in England, UK
travel hipcamp search --country united-kingdom --region england

# 3. Search outdoor spots in Ontario, Canada
travel hipcamp search --country canada --region ontario
```

---

### 23. Campspace (European Sustainable Micro-Camping)

Search unique nature spots, private farm pitches, treehouses, and campervan sites across the Netherlands, Germany, Belgium, France, Denmark, and Southern Europe via Chromedp headless extraction:

```bash
# 1. Search tent pitch micro-campsites
travel campspace search --category tent-pitches

# 2. Search campervan / RV micro-campsites
travel campspace search --category camper-sites

# 3. Search luxury glamping spots
travel campspace search --category glamping

# 4. Search unique treehouses in nature
travel campspace search --category treehouses
```

---

### 24. Azores Airlines / SATA Air Açores

Search the Portuguese island route network and 365-day low-fare calendars across the Azores archipelago, mainland Portugal (*Lisbon, Porto, Faro*), Europe (*Frankfurt, Paris*), and North America (*Boston, New York JFK, Toronto, Montreal*) via open public JSON endpoints:

```bash
# 1. List all airports in the Azores Airlines / SATA network
travel sata airports

# 2. List direct destination codes from Ponta Delgada (Azores hub)
travel sata routes --origin PDL

# 3. View 365-day low-fare calendar from Lisbon (LIS) to Ponta Delgada (PDL)
travel sata calendar --origin LIS --destination PDL --limit 15

# 4. View 365-day low-fare calendar from Boston (BOS) to Ponta Delgada (PDL)
travel sata calendar --origin BOS --destination PDL --limit 15
```

---

### 25. Cache Management

By default, all outgoing GET/GraphQL requests across providers are cached locally on disk in the `.cache` directory for **1 hour**. Repeated requests are returned instantly without making redundant network calls.

```bash
# Show cache directory, file count, and total size
travel cache info

# Clear all cached responses
travel cache clear
```

---

## 🌐 Provider Support & WAF Status

| Provider | Mode | Status | Real-Time Prices | Notes |
| :--- | :---: | :---: | :---: | :--- |
| **Ryanair** | ✈️ Flight | ✅ Supported | ✅ Yes | Direct JSON API |
| **Volotea** | ✈️ Flight | ✅ Supported | ✅ Yes | Public CDN JSON endpoints |
| **Vueling** | ✈️ Flight | ✅ Supported | ✅ Yes | EveryMundo AirTRFX + AMS JWT Auth |
| **airBaltic** | ✈️ Flight | ✅ Supported | ✅ Yes | Open FSF Availability & Network API |
| **FlyOne** | ✈️ Flight | ✅ Supported | ✅ Yes | Direct Routes & CMS Route Fare API with auto-token |
| **IndiGo (6E)** | ✈️ Flight | ✅ Supported | ✅ Yes | Automated JWT session token + `getfarecalendar` & `fare-radar` API |
| **TAP Air Portugal** | ✈️ Flight | ✅ Supported | ✅ Yes | Direct Public JSON API for airports, routes, and monthly low fare calendar |
| **Azores Airlines / SATA** | ✈️ Flight | ✅ Supported | ✅ Yes | Direct Public JSON API for airport directories, route maps, and 365-day low fare calendars |
| **Norwegian Air Shuttle** | ✈️ Flight | ✅ Supported | ✅ Yes | Chromedp Headless Browser Live Low-Fare Calendar Extraction |
| **Eurowings (Lufthansa Group)** | ✈️ Flight | ✅ Supported | ✅ Yes | Chromedp Headless Browser Live Airport Directory, Route Network, and Flight Schedule Dates |
| **Transavia (Air France-KLM)** | ✈️ Flight | ✅ Supported | ✅ Yes | Chromedp Headless Browser Live Monthly Low-Fare Calendar Extraction |
| **AirAsia / Arrivia Cruises** | 🚢 Cruise | ✅ Supported | ✅ Yes | Automated JWT auth + matrix & global cruise fare search |
| **Movacar** | 🚗 Car / 🚐 Campervan | ✅ Supported | ✅ Yes | Direct Cloud Run API for 1€ vehicle relocations |
| **imoova** | 🚐 Campervan / 🚗 Car | ✅ Supported | ✅ Yes | Direct GraphQL API for $1/day worldwide RV relocations |
| **DriiveMe** | 🚗 Car / 🚐 Van | ✅ Supported | ✅ Yes | Symfony Session + detail enrichment & availabilities API |
| **Trip.com** | 🏨 Hotel / 🚗 Car | ✅ Supported | ✅ Yes | Direct Next.js SSR HTTP GET for live hotel lists, final prices, room details, and rental cars |
| **Agoda** | 🏨 Hotel / Stay | ✅ Supported | ✅ Yes | Chromedp Headless Browser Live Extraction + Static CDN Directories |
| **Pitchup (`pitchup.com`)** | ⛺ Campsite / Glamping | ✅ Supported | ✅ Yes | Chromedp Headless Browser Live Search for outdoor stays across 67+ countries |
| **Hipcamp (`hipcamp.com`)** | 🌲 Glamping / Outdoor Stay | ✅ Supported | ✅ Yes | Chromedp Headless Browser Live Search for worldwide outdoor & glamping spots |
| **Campspace (`campspace.com`)** | 🌿 Micro-Camping / Nature | ✅ Supported | ✅ Yes | Chromedp Headless Browser Live Search for European micro-camping & eco stays |
| **Tictactrip (ComparaBUS)** | 🚆 Train / 🚌 Bus | ✅ Supported | ✅ Yes | Direct Public JSON API for European trains (SNCF/TGV/TER), buses (FlixBus/BlaBlaCar), and monthly price calendars |
| **Trenitalia & Le Frecce** | 🚆 Train | ✅ Supported | ✅ Yes | Direct BFF JSON API for Italian high-speed rail (Frecciarossa, Frecciargento, Intercity) and station directories |
| **oBilet (`obilet.com`)** | 🚌 Bus | ✅ Supported | ✅ Yes | Chromedp Headless Browser Live Search for Turkish intercity buses (Metro, Kamil Koç, Pamukkale) |
| **FlixBus & FlixTrain** | 🚌 Bus / Train | ✅ Supported | ✅ Yes | Global API with final platform fee |
| **Wizzair** | ✈️ Flight | ✅ Supported | ✅ Yes | Dynamic Build Timetable API |
| **easyJet** | ✈️ Flight | ⚠️ Protected | ❌ WAF Block | Protected by **Akamai Bot Manager** (`403 Access Denied`). Requires `_abck` sensor cookies or residential IP. |
| **Pegasus Airlines** | ✈️ Flight | ⚠️ Protected | ❌ WAF Block | Pricing & search endpoints (`web.flypgs.com`) are protected by **Akamai WAF**. Public endpoints only provide airports and route maps. |
| **Omio (`omio.com`)** | 🚌 Bus / 🚆 Train / ✈️ Flight | ⚠️ Protected | ❌ WAF Block | Protected by **Cloudflare Turnstile / Managed Challenge**. Direct HTTP calls return `Just a moment...` JS challenge. |
| **Scoot (`flyscoot.com`)** | ✈️ Flight | ⚠️ Protected | ❌ WAF Block | Protected by **Akamai Bot Manager WAF** (`403 Access Denied`). Requires browser environment / sensor cookies. |
| **Norse Atlantic (`flynorse.com`)** | ✈️ Flight | ⚠️ Protected | ❌ WAF Block | Route map and flight schedule dates (`services.flynorse.com`) are public, but fare availability (`/availability/lowfare`) is protected by **Cloudflare Bot Management** (`403 Access Restricted`). |

---

## 🛠 Development & Testing

Run static analysis and build verification:
```bash
go vet ./...
go test ./...
go build ./...
```




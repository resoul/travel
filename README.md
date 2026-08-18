# Travel CLI ✈️ 🚌 🚗 🚐

A command-line flight, bus, and vehicle relocation search and schedule lookup tool for low-cost carriers and rental platforms (**Ryanair**, **Wizzair**, **Volotea**, **Vueling**, **airBaltic**, **FlyOne**, **FlixBus**, **Movacar**, and **imoova**).

The project is built following **Clean Architecture** principles in Go.

---

## 🏛 Architecture

The codebase is organized into independent layers:

- **`internal/domain`** — Pure business entities (`FlightOffer`, `Airport`, `Country`, `FlightSearchCriteria`) and ports (interfaces `RyanairProvider`, `WizzairProvider`, `VoloteaProvider`, `VuelingProvider`, `FlixBusProvider`, `AirBalticProvider`, `FlyOneProvider`, `MovacarProvider`, `ImoovaProvider`). Has zero external dependencies.
- **`internal/usecase`** — Application business logic (`SearchFlightsUseCase`, `ListAirportsUseCase`, `FlightDatesUseCase`).
- **`internal/infrastructure`** — External API adapters (`infrastructure/ryanair`, `infrastructure/wizzair`, `infrastructure/volotea`, `infrastructure/vueling`, `infrastructure/airbaltic`, `infrastructure/flyone`, `infrastructure/movacar`, `infrastructure/imoova`, `infrastructure/flixbus`), disk-based caching (`infrastructure/cache`), DTOs, and domain mappers.
- **`internal/transport/cli`** — CLI user interface built with Cobra (`root`, `ryanair`, `wizzair`, `volotea`, `vueling`, `airbaltic`, `flyone`, `movacar`, `imoova`, `flixbus`, `cache`, `presenter`).
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
│           ├── flixbus.go
│           ├── cache.go
│           └── presenter.go
├── README.md
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
make ryanair-search ORIGIN=BBU DEST=GRO DATE=2026-08-22
make volotea-search ORIGIN=NTE DEST=VAR DATE=2026-08-25
make vueling-search ORIGIN=BCN DEST=FCO DATE=2026-08-19
make airbaltic-search ORIGIN=ALC DEST=RIX DATE=2026-08-28
make flyone-search ORIGIN=OTP DEST=BRU DATE=2026-08-28
make flixbus-search FROM="Bucharest" TO="Brasov" DATE=2026-08-27

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

### 10. Cache Management

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
| **Movacar** | 🚗 Car / 🚐 Campervan | ✅ Supported | ✅ Yes | Direct Cloud Run API for 1€ vehicle relocations |
| **imoova** | 🚐 Campervan / 🚗 Car | ✅ Supported | ✅ Yes | Direct GraphQL API for $1/day worldwide RV relocations |
| **FlixBus & FlixTrain** | 🚌 Bus / Train | ✅ Supported | ✅ Yes | Global API with final platform fee |
| **Wizzair** | ✈️ Flight | ✅ Supported | ✅ Yes | Dynamic Build Timetable API |
| **easyJet** | ✈️ Flight | ⚠️ Protected | ❌ WAF Block | Protected by **Akamai Bot Manager** (`403 Access Denied`). Requires `_abck` sensor cookies or headless browser. |
| **Pegasus Airlines** | ✈️ Flight | ⚠️ Protected | ❌ WAF Block | Pricing & search endpoints (`web.flypgs.com`) are protected by **Akamai WAF**. Public endpoints only provide airports and route maps. |
| **oBilet (`obilet.com`)** | 🚌 Bus / ✈️ Flight | ⚠️ Protected | ❌ WAF Block | Protected by **Cloudflare Bot Management / Managed Challenge**. Requires `cf_clearance` or headless browser. |
| **Omio (`omio.com`)** | 🚌 Bus / 🚆 Train / ✈️ Flight | ⚠️ Protected | ❌ WAF Block | Protected by **Cloudflare Turnstile / Managed Challenge**. Direct HTTP calls return `Just a moment...` JS challenge. |
| **Scoot (`flyscoot.com`)** | ✈️ Flight | ⚠️ Protected | ❌ WAF Block | Protected by **Akamai Bot Manager WAF** (`403 Access Denied`). Requires browser environment / sensor cookies. |

---

## 🛠 Development & Testing

Run static analysis and build verification:
```bash
go vet ./...
go test ./...
go build ./...
```




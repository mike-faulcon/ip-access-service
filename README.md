# ip-access-service

A lightweight Go service that determines whether an IP address originates from an allowed country using MaxMind GeoLite2 data.

## Configuration

Configuration is loaded from environment variables at startup. On launch, the service logs the resolved values:

    INFO config loaded http_port=8080 grpc_port=9090 geoip_path=data/GeoLite2-Country.mmdb

| Variable | Default | Description |
|----------|---------|-------------|
| `HTTP_PORT` | `8080` | HTTP server listen port |
| `GRPC_PORT` | `9090` | gRPC server listen port |
| `GEOIP_DB_PATH` | `data/GeoLite2-Country.mmdb` | Path to the MaxMind GeoLite2 Country database |

Invalid non-empty values (e.g. HTTP_PORT=abc) prevent startup.

### Examples

Run locally with custom ports:

```bash
HTTP_PORT=3001 GRPC_PORT=4001 go run ./cmd/server
```

Override the GeoIP database path:

```bash
GEOIP_DB_PATH=/path/to/GeoLite2-Country.mmdb go run ./cmd/server
```

## Running

### Prerequisites

This repo does not include the MaxMind database. Before running locally or with Docker, download `GeoLite2-Country.mmdb` into `data/` (see [GeoIP database](#geoip-database) below). Without it, the service fails at startup with an error opening the database file.

### Local

```bash
make run
```

### Docker

```bash
make docker
```

`docker-compose.yml` sets all three environment variables and mounts `./data` to `/data` for the GeoIP database.

## API

### Health check

```bash
curl http://localhost:8080/health
```

### Check IP access (HTTP)

```bash
curl -X POST http://localhost:8080/v1/check \
  -H "Content-Type: application/json" \
  -d '{"ip":"68.184.123.121","allowedCountries":["US","CA"]}'
```

### Check IP access (gRPC)

Requires [grpcurl](https://github.com/fullstorydev/grpcurl).

```bash
grpcurl -plaintext \
  -d '{"ip":"8.8.8.8","allowedCountries":["US","CA"]}' \
  localhost:9090 \
  access.v1.AccessService/CheckAccess
```

## Development

```bash
make test          # run all tests
make test-verbose  # run tests with verbose output
make test-http     # smoke test HTTP endpoint
make test-grpc     # smoke test gRPC endpoint
make vet           # run go vet
make coverage      # run tests with coverage
```

## GeoIP database

This service uses the [MaxMind GeoLite2 Country](https://dev.maxmind.com/geoip/geolite2-free-geolocation-data) database (`GeoLite2-Country.mmdb`). The file is not committed to the repo (see `.gitignore`).

### Obtaining the database

1. Create a free MaxMind account and generate a license key.
2. Download the GeoLite2 Country edition and extract `GeoLite2-Country.mmdb` into `data/` (or another path and set `GEOIP_DB_PATH`).
3. Follow MaxMind's license terms, including required attribution.

GeoLite2 Country is updated periodically (often weekly). A stale database can misclassify recently allocated IP ranges.

### Maintenance plan

We do not yet automate database updates. The phases below describe how we intend to improve maintenance over time.

**Phase 1 — Document and bootstrap (near term)**  
Document the manual download steps (above) so developers and operators can obtain the file without guesswork. Optionally add a download script or Makefile target that reads `MAXMIND_LICENSE_KEY` from the environment, downloads the tarball, and installs the `.mmdb` via atomic rename so partial downloads never corrupt the active file.

**Phase 2 — Scheduled refresh outside the app (medium term)**  
Refresh the file on a schedule without changing application code. Examples: host cron, CI scheduled job, or Kubernetes CronJob writing to a shared volume. After a successful download, restart the service (or redeploy an image built with a fresh DB). Store the license key in a secret manager; alert if the database is older than a chosen threshold (e.g. 14 days).

**Phase 3 — In-process hot reload (later, optional)**  
Periodically download a new database, validate it, atomically swap files, and reload the GeoIP reader in-process so HTTP/gRPC keep running without a restart. Would likely add configuration such as `GEOIP_UPDATE_INTERVAL` and optional health metadata (e.g. database last-modified time).

**Phase 4 — Production hardening (optional)**  
Startup checks for missing or overly stale databases, clearer operational visibility, and revisiting paid GeoIP2 Country if accuracy requirements increase.

| Phase | Effort | Requires code changes |
|-------|--------|------------------------|
| 1 — Document / download helper | Low | No (docs only) or minimal (script) |
| 2 — External scheduled refresh | Low–medium | No |
| 3 — Hot reload | Medium | Yes |
| 4 — Hardening | Varies | Yes |
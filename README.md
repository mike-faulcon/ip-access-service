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

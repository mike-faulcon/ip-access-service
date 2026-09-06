build:
	go build ./cmd/server

run:
	go run ./cmd/server

test:
	go test ./...

coverage:
	go test ./... -cover

vet:
	go vet ./...

tidy:
	go mod tidy

docker-build:
	docker build -t ip-access-service .

docker-run:
	docker run --rm \
		-p 8080:8080 \
		-e GEOIP_DB_PATH=/data/GeoLite2-Country.mmdb \
		-v "$$(pwd)/data:/data:ro" \
		ip-access-service
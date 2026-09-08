build:
	go build ./cmd/server

run:
	go run ./cmd/server

test:
	go test ./...

test-verbose:
	go test ./... -v

test-http:
	scripts/test-http.sh

test-grpc:
	scripts/test-grpc.sh

coverage:
	go test ./... -cover

vet:
	go vet ./...

tidy:
	go mod tidy

docker:
	docker compose up --build
build:
	go build ./cmd/server

run:
	go run ./cmd/server

test:
	go test ./...

vet:
	go vet ./...

tidy:
	go mod tidy
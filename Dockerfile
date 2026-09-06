# Build stage
FROM golang:1.27 AS builder

WORKDIR /app

# Download dependencies first so this layer can be cached
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build a statically linked Linux binary
RUN CGO_ENABLED=0 GOOS=linux go build \
    -o /app/ip-access-service \
    ./cmd/server


# Runtime stage
FROM gcr.io/distroless/static-debian13:nonroot

WORKDIR /app

COPY --from=builder /app/ip-access-service .

USER nonroot:nonroot

EXPOSE 8080

ENTRYPOINT ["/app/ip-access-service"]
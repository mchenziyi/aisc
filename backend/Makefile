.PHONY: build run test vet clean migrate

# Build the server binary
build:
	go build -o bin/server ./cmd/server

# Build the migrate binary
build-migrate:
	go build -o bin/migrate ./cmd/migrate

# Run the server
run:
	go run ./cmd/server

# Run all tests
test:
	go test ./... -v

# Run go vet
vet:
	go vet ./...

# Run migrations
migrate:
	go run ./cmd/migrate

# Clean build artifacts
clean:
	rm -rf bin/

# Build for production
production: build
	@echo "Build complete: bin/server"

# Run with air for hot-reload (requires air)
dev:
	air

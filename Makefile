.PHONY: ensure-mod

run-server:
	go run cmd/server/server.go

run-web:
	go run web/web.go

lint:
	golangci-lint run

test:
	echo "Running unit tests..."; \
	go test --race -v ./... | { command -v gocol >/dev/null && gocol || cat; };

test-coverage:
	go test -tags=integration -v -coverprofile=coverage.out ./... | { command -v gocol >/dev/null && gocol || cat; }
	go tool cover -html=coverage.out

ensure-mod:
	go mod download

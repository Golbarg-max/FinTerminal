.PHONY: build run test clean install help

BINARY_NAME=fin
BINARY_PATH=./$(BINARY_NAME)

help:
	@echo "FinTerminal - Makefile targets:"
	@echo "  make build    - Build the binary"
	@echo "  make run      - Run the application"
	@echo "  make test     - Run tests"
	@echo "  make clean    - Clean build artifacts"
	@echo "  make install  - Install binary to /usr/local/bin"
	@echo "  make deps     - Download dependencies"

build:
	@echo "Building FinTerminal..."
	@go build -o $(BINARY_NAME) -v

run: build
	@echo "Running FinTerminal..."
	@$(BINARY_PATH) $(ARGS)

test:
	@echo "Running tests..."
	@go test -v ./...

clean:
	@echo "Cleaning build artifacts..."
	@rm -f $(BINARY_NAME)
	@go clean

install: build
	@echo "Installing to /usr/local/bin/$(BINARY_NAME)..."
	@sudo mv $(BINARY_PATH) /usr/local/bin/
	@echo "✓ Installed. Run 'fin' from anywhere"

deps:
	@echo "Downloading dependencies..."
	@go mod download
	@go mod tidy

fmt:
	@echo "Formatting code..."
	@go fmt ./...

lint:
	@echo "Running linter..."
	@golangci-lint run ./...

quote:
	@echo "Getting quote for AAPL..."
	@$(BINARY_PATH) quote AAPL

portfolio:
	@echo "Showing portfolio..."
	@$(BINARY_PATH) portfolio

example:
	@echo "Adding example holdings to portfolio..."
	@$(BINARY_PATH) portfolio add AAPL 10 150.50
	@$(BINARY_PATH) portfolio add MSFT 25 250.00
	@$(BINARY_PATH) portfolio add GOOGL 5 2800.00

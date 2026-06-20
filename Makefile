BINARY := fret
BUILD_DIR := bin
MAIN := ./cmd/fret

.PHONY: help build run test fmt clean

help:
	@printf "Targets:\n"
	@printf "  make build   Build $(BUILD_DIR)/$(BINARY)\n"
	@printf "  make run     Run the TUI\n"
	@printf "  make test    Run tests\n"
	@printf "  make fmt     Format Go files\n"
	@printf "  make clean   Remove build output\n"

build:
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY) $(MAIN)

run:
	go run $(MAIN)

test:
	go test ./...

fmt:
	gofmt -w cmd internal

clean:
	rm -rf $(BUILD_DIR)

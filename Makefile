GO := go
SRC := .
BIN := _netmon
BIN_DIR := ./bin
BIN_LOC := $(BIN_DIR)/$(BIN)

.PHONY: build clean start

build:
	@mkdir -p ./bin
	$(GO) build -o $(BIN_LOC) $(SRC)

start:build
	$(BIN_LOC)

clean:
	@rm -r $(BIN_DIR)
	@echo "Removed `./bin` directory"

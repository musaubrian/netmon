GO := go
SRC := .
BIN := netmon
BIN_DIR := ./bin
LOG_DIR := ./logs
BIN_LOC := $(BIN_DIR)/$(BIN)

.PHONY: build clean start test

build:
	@mkdir -p $(BIN_DIR)
	@mkdir -p $(LOG_DIR)
	@$(GO) build -o $(BIN_LOC) $(SRC)

start:build
	@$(BIN_LOC)
	
test:
	@go test -v ./... 

clean:
	@rm -r $(BIN_DIR)
	@rm -r $(LOG_DIR)
	@echo "Removed './bin && ./logs' directory"

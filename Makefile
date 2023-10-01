PROJ_DIR=$(PWD)
BUILD_DIR=$(PROJ_DIR)/_build
MAIN_FILE=cmd/main.go

.PHONY: build
build:
	go build -o $(BUILD_DIR)/pm $(MAIN_FILE)

.PHONY: test
test:
	go test -v ./internal/packet/validator
	go test -v ./internal/packet/archiver
	go test -v ./internal/packet/files
	go test -v ./internal/packet/parser
	go test -v ./internal/adapter/pmssh



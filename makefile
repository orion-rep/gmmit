BINARY_NAME=build/gmmit
MAIN_GO_FILE=./cmd/gmmit/
INSTALL_DIR:=/usr/local/bin
APP_VERSION := $(shell git describe --tags --always --dirty)

.PHONY: clean build test cover

build:
	@echo "Building version $(APP_VERSION)"
	go mod tidy
	go fmt ./...
	go build -ldflags "-X main.Version=$(APP_VERSION)" -o $(BINARY_NAME) $(MAIN_GO_FILE)

run: build
	./$(BINARY_NAME)

install: build
	@echo "Installing $(BINARY_NAME) to $(INSTALL_DIR)..."
	sudo install -m 755 $(BINARY_NAME) $(INSTALL_DIR)

test:
	go test -race -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out

cover: test
	go tool cover -html=coverage.out

clean:
	go clean
	rm -f ${BINARY_NAME}

tapes:
	vhs gif-commit.tape
	vhs gif-pull-request.tape
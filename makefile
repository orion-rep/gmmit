BINARY_NAME=build/gmmit
MAIN_GO_FILE=./cmd/gmmit/
APP_VERSION!=git describe --tags --always --dirty

.PHONY : clean build

build:
	@echo "Building version $(APP_VERSION)"
	go mod tidy
	go build -ldflags "-X main.Version=$(APP_VERSION)" -o $(BINARY_NAME) $(MAIN_GO_FILE)

run: build
	./$(BINARY_NAME)

clean:
	go clean
	rm -f ${BINARY_NAME}
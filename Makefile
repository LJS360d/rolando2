.PHONY: all build test lint clean dev build-dev run

ifeq ($(OS),Windows_NT)
  EXE 	:= .exe
  RM 		= del
else
  EXE 	:=
  RM 		= rm -f
endif

VERSION         := 3.0.0
BUILD_DIR       := bin
MAIN_PACKAGE    := ./app
BUILD						:= $(shell git rev-parse --short HEAD)
ENV 						?= production

ifeq ($(ENV), production)
	BINARY_NAME   = rolando@$(VERSION)
else
	BINARY_NAME   = dev
endif

BUILDPATH       = $(BUILD_DIR)/$(BINARY_NAME)$(EXE)
LDFLAGS 				= -ldflags "-w -s -X main.Version=$(VERSION) -X main.Build=$(BUILD) -X main.Env=$(ENV)"

all: dev

build:
	go build $(LDFLAGS) -o $(BUILDPATH) $(MAIN_PACKAGE)

test:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

lint:
	go fmt ./...
	staticcheck ./...

clean:
	go clean
	$(RM) $(BUILD_DIR)

run: build
	./$(BUILDPATH)

dev:
	air

build-dev: ENV=development
build-dev:
	go build $(LDFLAGS) -o $(BUILDPATH) $(MAIN_PACKAGE)

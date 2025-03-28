ifeq ($(OS),Windows_NT)
	EXE := .exe
	SET_ENV := set
	SEP := &
else
	EXE :=
	SET_ENV :=
	SEP := ;
endif

BUILD_DIR 		= bin
BINARY_NAME 	= save-parser
BUILDPATH 		= $(BUILD_DIR)/$(BINARY_NAME)$(EXE)
MAIN_PACKAGE 	= ./cmd

WASM_BINARY_NAME 	= save-parser.wasm
WASM_BUILDPATH 		= ../docs/public/$(WASM_BINARY_NAME)
WASM_MAIN_PACKAGE 	= ./export

WASM_BUILDFLAGS = $(SET_ENV) GOOS=js$(SEP) $(SET_ENV) GOARCH=wasm$(SEP)

all: wasm

build:
	go build -o $(BUILDPATH) $(MAIN_PACKAGE)

wasm:
	$(WASM_BUILDFLAGS) go build -o $(WASM_BUILDPATH) $(WASM_MAIN_PACKAGE)

test:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

lint:
	go fmt ./...

clean:
	go clean
	rm -rf $(BUILD_DIR)
	rm -rf build

run: build
	./$(BUILDPATH)

.PHONY: all build wasm test lint clean run

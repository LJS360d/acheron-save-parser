ifeq ($(OS),Windows_NT)
	EXE := .exe
	DEL := del /f
	SET_ENV := set
	SEP := &
else
	EXE :=
	DEL := rm -f
	SET_ENV :=
	SEP := ;
endif

BUILD_DIR 		= bin
BINARY_NAME 	= acheron-save-parser
BUILDPATH 		= $(BUILD_DIR)/$(BINARY_NAME)$(EXE)
MAIN_PACKAGE 	= ./cmd

WASM_BINARY_NAME 	= main.wasm
WASM_BUILDPATH 		= ../docs/public/$(WASM_BINARY_NAME)
WASM_MAIN_PACKAGE 	= ./export

all: build-wasm

build:
	go build -o $(BUILDPATH) $(MAIN_PACKAGE)

build-wasm:
	$(SET_ENV) GOOS=js$(SEP) $(SET_ENV) GOARCH=wasm$(SEP) go build -o $(WASM_BUILDPATH) $(WASM_MAIN_PACKAGE)

test:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

lint:
	go fmt ./...

clean:
	go clean
	$(DEL) $(BUILDPATH)
	$(DEL) $(WASM_BUILDPATH)

run: build
	./$(BUILDPATH)

.PHONY: all build build-wasm test lint clean run

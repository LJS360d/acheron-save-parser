# Makefile
ifeq ($(OS),Windows_NT)
	EXE := .exe
  DEL = del /f
else
	EXE :=
  DEL = rm -f
endif
BUILD_DIR 		= bin
BINARY_NAME 	= acheron-save-parser
BUILDPATH 		= $(BUILD_DIR)/$(BINARY_NAME)$(EXE)
MAIN_PACKAGE 	= ./cmd


all: build

build:
	go build -o $(BUILDPATH) $(MAIN_PACKAGE)

test:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

lint:
	go fmt ./...

staticcheck:
	staticcheck ./...

clean:
	go clean
	$(DEL) $(BINARY_NAME)

run: build
	./$(BUILDPATH)

.PHONY: all build test lint staticcheck clean run

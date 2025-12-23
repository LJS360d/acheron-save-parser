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

all: cmp

USE_GBA_FILE_OUR := test/gba/acheron-emerald.gba
USE_MAP_FILE_OUR := test/memmap/acheron-emerald.map
USE_OUTPUTS := species,evolutions,moves,learnsets,items,encounters,trainers
USE_JBP_OUR := our_
VEXTRA := acheron-emerald

USE_GBA_FILE_RHH := test/gba/rhh1.11.8.gba
USE_MAP_FILE_RHH := test/memmap/rhh1.11.8.map
USE_JBP_OUR := rhh_

use: build
	$(BUILDPATH) -gba='$(USE_GBA_FILE_OUR)' -map='$(USE_MAP_FILE_OUR)' -o='$(USE_OUTPUTS)' -jbp=$(USE_JBP_OUR) -vextra=$(VEXTRA)
	$(BUILDPATH) -gba='$(USE_GBA_FILE_RHH)' -map='$(USE_MAP_FILE_RHH)' -o='$(USE_OUTPUTS)' -jbp=$(USE_JBP_RHH)

BUILD_COMPARE := scripts/build_compare.go
RES_DIR := build
cmp: use
	go run $(BUILD_COMPARE) -new='$(RES_DIR)/$(USE_JBP_OUR)species.json' -old='$(RES_DIR)/$(USE_JBP_RHH)species.json' -id='id' -o='$(RES_DIR)/species.json' -d
	go run $(BUILD_COMPARE) -new='$(RES_DIR)/$(USE_JBP_OUR)learnsets.json' -old='$(RES_DIR)/$(USE_JBP_RHH)learnsets.json' -id='species' -o='$(RES_DIR)/learnsets.json' -d
	go run $(BUILD_COMPARE) -new='$(RES_DIR)/$(USE_JBP_OUR)evolutions.json' -old='$(RES_DIR)/$(USE_JBP_RHH)evolutions.json' -id='family' -o='$(RES_DIR)/evolutions.json' -d
	go run $(BUILD_COMPARE) -new='$(RES_DIR)/$(USE_JBP_OUR)items.json' -old='$(RES_DIR)/$(USE_JBP_RHH)items.json' -id='id' -o='$(RES_DIR)/items.json' -d
	go run $(BUILD_COMPARE) -new='$(RES_DIR)/$(USE_JBP_OUR)moves.json' -old='$(RES_DIR)/$(USE_JBP_RHH)moves.json' -id='id' -o='$(RES_DIR)/moves.json' -d
	go run $(BUILD_COMPARE) -new='$(RES_DIR)/$(USE_JBP_OUR)wild_encounters.json' -old='$(RES_DIR)/$(USE_JBP_RHH)wild_encounters.json' -id='mapNum' -o='$(RES_DIR)/wild_encounters.json' -d


$(BUILDPATH):
	go build -o $(BUILDPATH) $(MAIN_PACKAGE)

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

package gba

import "runtime"

const (
	wasm_build     = runtime.GOARCH == "wasm" && runtime.GOOS == "js"
	POINTER_OFFSET = 0x08000000
	NULL_POINTER   = 0x0f8000000
)

var (
	Config GbaConfig
	Data   []byte

	// parsed data
	Header         *GbaHeader
	Abilities      []*AbilityData
	Species        []*SpeciesData
	Items          []*ItemData
	Natures        []*NatureData
	Moves          []*MoveData
	WildEncounters []*WildPokemonHeader
)

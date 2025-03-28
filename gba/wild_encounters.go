package gba

import (
	"encoding/binary"
	"log"
)

type WildPokemon struct {
	MinLevel uint8  `json:"minLevel"`
	MaxLevel uint8  `json:"maxLevel"`
	Species  uint16 `json:"species"`
}

type WildPokemonInfo struct {
	EncounterRate uint8         `json:"encounterRate"`
	WildPokemon   []WildPokemon `json:"wildPokemon"`
}

type WildPokemonHeader struct {
	MapGroup          uint8            `json:"mapGroup"`
	MapNum            uint8            `json:"mapNum"`
	LandMonsInfo      *WildPokemonInfo `json:"landMonsInfo"`
	WaterMonsInfo     *WildPokemonInfo `json:"waterMonsInfo"`
	RockSmashMonsInfo *WildPokemonInfo `json:"rockSmashMonsInfo"`
	FishingMonsInfo   *WildPokemonInfo `json:"fishingMonsInfo"`
}

func ParseWildEncounters(startOffset int) []*WildPokemonHeader {
	encounters := make([]*WildPokemonHeader, 0)
	if startOffset == NULL_POINTER || startOffset >= len(Data) {
		log.Printf("[WARN] Wild encounters offset [%x] is null or out of bounds", startOffset)
		return encounters
	}
	for i := 0; i < Config.WildEncounterCount; i++ {
		offset := startOffset + i*Config.WildEncounterSize
		encounter := &WildPokemonHeader{}
		encounter.MapGroup = uint8(Data[offset])
		encounter.MapNum = uint8(Data[offset+1])
		// align pointer, 2 bytes of padding
		landMonsInfoPtr := binary.LittleEndian.Uint32(Data[offset+4:]) - POINTER_OFFSET
		encounter.LandMonsInfo = parseWildPokemonInfo(landMonsInfoPtr)

		waterMonsInfoPtr := binary.LittleEndian.Uint32(Data[offset+8:]) - POINTER_OFFSET
		encounter.WaterMonsInfo = parseWildPokemonInfo(waterMonsInfoPtr)

		rockSmashMonsInfoPtr := binary.LittleEndian.Uint32(Data[offset+12:]) - POINTER_OFFSET
		encounter.RockSmashMonsInfo = parseWildPokemonInfo(rockSmashMonsInfoPtr)

		fishingMonsInfoPtr := binary.LittleEndian.Uint32(Data[offset+16:]) - POINTER_OFFSET
		encounter.FishingMonsInfo = parseWildPokemonInfo(fishingMonsInfoPtr)

		encounters = append(encounters, encounter)
	}
	return encounters
}

func parseWildPokemonInfo(offset uint32) *WildPokemonInfo {
	if offset == NULL_POINTER || offset >= uint32(len(Data)) {
		return nil
	}
	info := &WildPokemonInfo{
		EncounterRate: Data[offset],
		WildPokemon:   make([]WildPokemon, 0),
	}
	// align pointer, 3 bytes of padding
	wildPokemonPtr := binary.LittleEndian.Uint32(Data[offset+4:]) - POINTER_OFFSET
	const WILD_POKEMON_SIZE = 4
	for i := wildPokemonPtr; i+WILD_POKEMON_SIZE < uint32(len(Data)); i += WILD_POKEMON_SIZE {
		wildPokemon := &WildPokemon{
			MinLevel: Data[i],
			MaxLevel: Data[i+1],
			Species:  binary.LittleEndian.Uint16(Data[i+2:]),
		}
		if wildPokemon.MinLevel == 0 || wildPokemon.MaxLevel == 0 || wildPokemon.Species == 0 {
			// in the original C code the array is fixed size, so theres no standard way to know when to stop reading
			// so we just check that an entry that was parsed makes no sense
			break
		}
		info.WildPokemon = append(info.WildPokemon, *wildPokemon)
	}
	return info
}

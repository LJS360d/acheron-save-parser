package main

import (
	jsonconvert "acheron-save-parser/export/json"
	"acheron-save-parser/gba"
	"acheron-save-parser/utils"
	"bufio"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"image/color"
	"os"
	"strings"
)

type JSON = map[string]interface{}

var (
	TYPES = []string{
		"NONE",
		"NORMAL",
		"FIGHTING",
		"FLYING",
		"POISON",
		"GROUND",
		"ROCK",
		"BUG",
		"GHOST",
		"STEEL",
		"MYSTERY",
		"FIRE",
		"WATER",
		"GRASS",
		"ELECTRIC",
		"PSYCHIC",
		"ICE",
		"DRAGON",
		"DARK",
		"FAIRY",
		"STELLAR",
	}
	EGG_GROUPS = []string{
		"NONE",
		"MONSTER",
		"WATER_1",
		"BUG",
		"FLYING",
		"FIELD",
		"FAIRY",
		"GRASS",
		"HUMAN_LIKE",
		"WATER_3",
		"MINERAL",
		"AMORPHOUS",
		"WATER_2",
		"DITTO",
		"DRAGON",
		"NO_EGGS_DISCOVERED",
	}
)

const (
	ICON_PALETTES_COUNT     = 6
	PALETTE_SIZE            = 32
	COMPRESSED_PALETTE_SIZE = 40 // yes compressed is bigger than uncompressed, gamefreak probably had a reason to compress otherwise they just had the big stupid
)

func SaveItemsIcons(data []byte, items []*gba.ItemData) error {
	for i := 0; i < len(items); i++ {
		if items[i].IconPalettePtr == gba.BAD_POINTER {
			fmt.Println("MISSING ITEM PALETTE FOR", i, items[i].Name)
			continue
		}
		if items[i].IconPicPtr == gba.BAD_POINTER {
			fmt.Println("MISSING ITEM PIC FOR", i, items[i].Name)
			continue
		}
		compressedPalBytes := data[items[i].IconPalettePtr : items[i].IconPalettePtr+COMPRESSED_PALETTE_SIZE]
		palBytes, err := utils.DecompressLZ77(compressedPalBytes)
		if err != nil {
			return fmt.Errorf("ERROR DECOMPRESSING PALETTE FOR %d: %w", i, err)
		}
		pal := utils.ParsePaletteBytes(palBytes)
		compressedIconBytes := data[items[i].IconPicPtr : items[i].IconPicPtr+292]
		iconBytes, err := utils.DecompressLZ77(compressedIconBytes)
		if err != nil {
			return fmt.Errorf("ERROR DECOMPRESSING ICON FOR %d: %w", i, err)
		}
		err = utils.Save4bppImageBytes(iconBytes, "build/images/items/icons/"+fmt.Sprint(i), pal, 24, 24, true)
		if err != nil {
			return fmt.Errorf("ERROR SAVING ITEM ICON FOR %d: %w", i, err)
		}
	}
	return nil
}

func SaveSpeciesIcons(data []byte, s []*gba.SpeciesData, iconPalettesPtr uint32) error {
	iconPalettes := make([][]color.Color, 0)
	palPtr := binary.LittleEndian.Uint32(data[iconPalettesPtr:iconPalettesPtr+4]) - gba.POINTER_OFFSET
	for i := 0; i < ICON_PALETTES_COUNT; i++ {
		palBytes := data[palPtr+uint32(i*PALETTE_SIZE) : palPtr+uint32(i*PALETTE_SIZE+PALETTE_SIZE)]
		iconPalettes = append(iconPalettes, utils.ParsePaletteBytes(palBytes))
	}
	for i := 0; i < len(s); i++ {
		if s[i].IconSpritePtr == gba.BAD_POINTER {
			continue
		}
		pal := iconPalettes[s[i].IconPalIndex]
		iconBytes := data[s[i].IconSpritePtr : s[i].IconSpritePtr+1024]
		err := utils.Save4bppImageBytes(iconBytes, "build/images/pokemon/icons/"+fmt.Sprint(i), pal, 32, 32, true)
		if err != nil {
			return fmt.Errorf("ERROR SAVING POKEMON ICON FOR %d: %w", i, err)
		}
	}
	return nil
}

func SaveSpeciesSprites(data []byte, s []*gba.SpeciesData) error {
	for i := 0; i < len(s); i++ {
		if s[i].FrontPicPtr == gba.BAD_POINTER {
			continue
		}
		decompressedPalBytes, err := utils.DecompressLZ77(data[s[i].PalettePtr : s[i].PalettePtr+COMPRESSED_PALETTE_SIZE])
		if err != nil {
			return fmt.Errorf("ERROR DECOMPRESSING POKEMON PALETTE FOR %d: %w", i, err)
		}
		pal := utils.ParsePaletteBytes(decompressedPalBytes)
		if err != nil {
			return fmt.Errorf("MISSING POKEMON PALETTE FOR %d: %w", i, err)
		}
		frontPicBytesCompressed := data[s[i].FrontPicPtr : s[i].FrontPicPtr+4096]
		frontPicBytes, err := utils.DecompressLZ77(frontPicBytesCompressed)
		if err != nil {
			return fmt.Errorf("ERROR DECOMPRESSING POKEMON FRONT PIC FOR %d: %w", i, err)
		}
		err = utils.Save4bppImageBytes(frontPicBytes, "build/images/pokemon/sprites/"+fmt.Sprint(i), pal, 64, 64, true)
		if err != nil {
			return fmt.Errorf("ERROR SAVING POKEMON FRONT PIC FOR %d: %w", i, err)
		}
	}
	return nil
}

func SaveSpeciesData(filepath string, s []*gba.SpeciesData) error {
	file, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("ERROR CREATING FILE: %w", err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	json.NewEncoder(writer).Encode(utils.MapSlice(s,
		func(mon *gba.SpeciesData, i int) JSON {
			return JSON{
				"id":          i,
				"species":     getSpeciesIdentifier(mon, uint16(i+1)),
				"speciesName": mon.SpeciesName,
				"stats": jsonconvert.MarshalSlice([]uint8{
					mon.BaseHP,
					mon.BaseAttack,
					mon.BaseDefense,
					mon.BaseSpAttack,
					mon.BaseSpDefense,
					mon.BaseSpeed,
				}),
				"bst":        mon.Bst,
				"generation": mon.Generation,
				"types": jsonconvert.MarshalSlice(
					utils.PruneDuplicates(
						utils.MapSlice(mon.Types[:],
							func(t uint8, i int) string {
								return TYPES[t]
							}),
					),
				),
				"natDexNum": mon.NatDexNum,
				"abilities": jsonconvert.MarshalSlice(
					utils.PruneDuplicates(
						utils.FilterEmpty(
							utils.MapSlice(mon.Abilities[:2], func(a uint16, i int) string {
								if a == 0 {
									return ""
								}
								return strings.ToUpper(utils.ToSnakeCase(gba.Abilities[a].Name))
							}),
						))),
				"bodyColor":    mon.BodyColor,
				"catchRate":    mon.CatchRate,
				"categoryName": mon.CategoryName,
				"description":  mon.Description,
				"eggCycles":    mon.EggCycles,
				"eggGroups": jsonconvert.MarshalSlice(
					utils.PruneDuplicates(utils.MapSlice(mon.EggGroups[:], func(g uint8, i int) string {
						return EGG_GROUPS[g]
					}),
					)),
				"formChangeTable": jsonconvert.MarshalSlice(
					utils.MapSlice(mon.FormChangeTable, func(change *gba.FormChange, i int) JSON {
						return JSON{
							"form":   change.TargetSpecies,
							"method": change.Method,
							"params": jsonconvert.MarshalSlice([]uint16{change.Param1, change.Param2, change.Param3}),
						}
					},
					)),
				"evYield": jsonconvert.MarshalSlice([]uint8{
					mon.EvYieldHP,
					mon.EvYieldAttack,
					mon.EvYieldDefense,
					mon.EvYieldSpeed,
					mon.EvYieldSpAttack,
					mon.EvYieldSpDefense,
				}),
				"expYield":    mon.ExpYield,
				"genderRatio": mon.GenderRatio,
				"growthRate":  mon.GrowthRate,
				"height":      mon.Height,
				"weight":      mon.Weight,
				"flags":       jsonconvert.MarshalSlice(getFlagArray(mon)),
				"itemCommon":  mon.ItemCommon,
				"itemRare":    mon.ItemRare,
			}
		}))
	err = writer.Flush()
	if err != nil {
		return fmt.Errorf("ERROR WRITING TO FILE: %w", err)
	}
	return nil
}

func getFlagArray(mon *gba.SpeciesData) []string {
	flags := []string{}

	flagMap := map[bool]string{
		mon.IsLegendary:       "LEGENDARY",
		mon.IsMythical:        "MYTHICAL",
		mon.IsUltraBeast:      "ULTRABEAST",
		mon.IsParadox:         "PARADOX",
		mon.IsTotem:           "TOTEM",
		mon.IsMegaEvolution:   "MEGAEVOLUTION",
		mon.IsPrimalReversion: "PRIMAL",
		mon.IsUltraBurst:      "ULTRABURST",
		mon.IsGigantamax:      "GIGANTAMAX",
		mon.IsTeraForm:        "TERAFORM",
		mon.IsAlolanForm:      "ALOLAN",
		mon.IsGalarianForm:    "GALARIAN",
		mon.IsHisuianForm:     "HISUIAN",
		mon.IsPaldeanForm:     "PALDEAN",
	}

	for condition, flag := range flagMap {
		if condition {
			flags = append(flags, flag)
		}
	}

	return flags
}

func getSpeciesIdentifier(mon *gba.SpeciesData, index uint16) string {
	speciesName := strings.ReplaceAll(
		strings.ReplaceAll(
			strings.ToUpper(utils.ToSnakeCase(mon.SpeciesName)),
			" ", "_"),
		"'", "_")
	flagMap := map[bool]string{
		mon.IsTotem:           "TOTEM",
		mon.IsMegaEvolution:   "MEGA",
		mon.IsPrimalReversion: "PRIMAL",
		mon.IsUltraBurst:      "ULTRA",
		mon.IsGigantamax:      "GIGANTAMAX",
		mon.IsTeraForm:        "TERA",
		mon.IsAlolanForm:      "ALOLAN",
		mon.IsGalarianForm:    "GALARIAN",
		mon.IsHisuianForm:     "HISUIAN",
		mon.IsPaldeanForm:     "PALDEAN",
	}
	for condition, flag := range flagMap {
		if condition {
			return speciesName + "_" + flag
		}
	}
	return speciesName
}

package main

import (
	jsonconvert "acheron-save-parser/export/json"
	"acheron-save-parser/gba"
	"acheron-save-parser/utils"
	"bufio"
	"encoding/json"
	"fmt"
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

func SaveSpeciesSprites(data []byte, s []*gba.SpeciesData) error {
	for i := 0; i < len(s); i++ {
		if s[i].FrontPicPtr == gba.BAD_POINTER {
			continue
		}
		decompressedPalBytes, err := utils.DecompressLZ77(data[s[i].PalettePtr : s[i].PalettePtr+40])
		if err != nil {
			return fmt.Errorf("ERROR DECOMPRESSING PALETTE FOR %d: %w", i, err)
		}
		pal := utils.ParsePaletteBytes(decompressedPalBytes)
		if err != nil {
			return fmt.Errorf("MISSING PALETTE FOR %d: %w", i, err)
		}
		frontPicBytes := data[s[i].FrontPicPtr : s[i].FrontPicPtr+4096]
		err = utils.Save4bppImageBytes(frontPicBytes, "build/images/sprites/"+fmt.Sprint(i), pal, 64, 64, true)
		if err != nil {
			return fmt.Errorf("ERROR SAVING FRONT PIC FOR %d: %w", i, err)
		}
	}
	return nil
}

func SaveSpeciesData(filename string, s []*gba.SpeciesData) {
	if !strings.HasSuffix(filename, ".json") {
		filename += ".json"
	}
	file, err := os.Create("build/" + filename)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
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
		fmt.Println("Error writing to file:", err)
	}
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
	rename := SpeciesRenames[index]
	if rename != "" {
		return rename
	}
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

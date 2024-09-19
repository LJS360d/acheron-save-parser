package main

import (
	jsonconvert "acheron-save-parser/export/json"
	"acheron-save-parser/gba"
	"acheron-save-parser/sav"
	"acheron-save-parser/utils"
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

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

func main() {
	savFile := flag.String("s", "", "Path to the save file (.sav)")
	gbaFile := flag.String("g", "", "Path to the GBA ROM file (.gba)")

	// Support for --sav and --gba flags
	flag.StringVar(savFile, "sav", "", "Path to the save file (.sav)")
	flag.StringVar(gbaFile, "gba", "", "Path to the GBA ROM file (.gba)")

	flag.Parse()

	if *savFile == "" || *gbaFile == "" {
		log.Fatal("Both -s/-sav and -g/-gba flags are required.")
	}

	fmt.Printf("Save file path: %s\n", *savFile)
	fmt.Printf("GBA file path: %s\n", *gbaFile)

	gbaBytes, err := os.ReadFile(*gbaFile)
	if err != nil {
		log.Fatal(err)
	}
	savBytes, err := os.ReadFile(*savFile)
	if err != nil {
		log.Fatal(err)
	}
	gba.ParseGbaBytes(gbaBytes)
	sav.ParseSavBytes(savBytes)
	// writeSpeciesData("species.json", gba.Species[1:])
	/*
		directories, err := getDirectories("../game/graphics/pokemon")
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		ReverseSpeciesRenames := make(map[string]uint16)
		for key, value := range SpeciesRenames {
			ReverseSpeciesRenames[value] = key
		}

		for _, mon := range directories {
			if ReverseSpeciesRenames[strings.ToUpper(mon)] == 0 {
				fmt.Println("Missing ID for ", mon)
				continue
			}
			normalPalPath := "../game/graphics/pokemon/" + strings.ReplaceAll(mon, "_", "/") + "/normal.pal"
			uscorePalPath := "../game/graphics/pokemon/" + mon + "/normal.pal"
			// shinyPalPath := "../game/graphics/pokemon/" + strings.ReplaceAll(mon, "_", "/") + "/shiny.pal"
			utils.CopyFile(normalPalPath, "assets/palettes/normal/"+fmt.Sprint(ReverseSpeciesRenames[strings.ToUpper(mon)])+".pal")
			utils.CopyFile(uscorePalPath, "assets/palettes/normal/"+fmt.Sprint(ReverseSpeciesRenames[strings.ToUpper(mon)])+".pal")
		} */
}

func getDirectories(root string) ([]string, error) {
	var dirs []string

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// Check if it's a directory (ignore the root directory itself)
		if info.IsDir() && path != root {
			path = filepath.ToSlash(path)
			// Remove the rootDir from the path
			relativePath := strings.TrimPrefix(path+"/", root)

			// Replace the OS-specific separator with "_"
			relativePath = strings.ReplaceAll(relativePath, "/", "_")

			// Ensure we don't have a leading separator (like "_subdir1")
			relativePath = strings.TrimSuffix(strings.TrimPrefix(relativePath, "_"), "_")

			dirs = append(dirs, relativePath)
		}
		return nil
	})

	return dirs, err
}

type JSON = map[string]interface{}

func writeSpeciesData(filename string, species []*gba.SpeciesData) {
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
	json.NewEncoder(writer).Encode(utils.MapSlice(species,
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

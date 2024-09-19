package main

import (
	"acheron-save-parser/gba"
	"acheron-save-parser/sav"
	"flag"
	"fmt"
	"log"
	"os"
	"slices"
	"strings"
)

func main() {
	savFile := flag.String("s", "", "Path to the save file (.sav)")
	gbaFile := flag.String("g", "", "Path to the GBA ROM file (.gba)")
	outputs := flag.String("o", "", "Comma-separated list of outputs to generate (e.g., species,moves,learnsets,items,sprites)")

	// Support for --sav and --gba flags
	flag.StringVar(savFile, "sav", "", "Path to the save file (.sav)")
	flag.StringVar(gbaFile, "gba", "", "Path to the GBA ROM file (.gba)")
	flag.StringVar(outputs, "output", "", "Comma-separated list of outputs to generate (e.g., species,moves,learnsets,items,sprites)")

	flag.Parse()

	if *savFile == "" || *gbaFile == "" {
		log.Fatal("Both -s/-sav and -g/-gba flags are required.")
	}

	fmt.Printf("Save file path: %s\n", *savFile)
	fmt.Printf("GBA file path: %s\n", *gbaFile)
	fmt.Printf("Outputs: %s\n", *outputs)

	gbaBytes, err := os.ReadFile(*gbaFile)
	if err != nil {
		log.Fatal(err)
	}
	savBytes, err := os.ReadFile(*savFile)
	if err != nil {
		log.Fatal(err)
	}
	g := gba.ParseGbaBytes(gbaBytes)
	sav.ParseSavBytes(savBytes)

	selectedOutputs := strings.Split(*outputs, ",")

	if slices.Contains(selectedOutputs, "species") {
		SaveSpeciesData("build/species_rhh.json", gba.Species[1:])
	}
	if slices.Contains(selectedOutputs, "sprites") {
		SaveSpeciesSprites(gbaBytes, gba.Species[1:])
		SaveSpeciesIcons(gbaBytes, gba.Species[1:], g.IconPalettesTablePtr)
	}
}

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
	"sync"
)

var (
	JSON_BUILDS_PREFIX = ""
)

func main() {
	savFile := flag.String("s", "", "Path to the save file (.sav)")
	gbaFile := flag.String("g", "", "Path to the GBA ROM file (.gba)")
	outputs := flag.String("o", "", "Comma-separated list of outputs to generate (e.g., species,evolutions,moves,learnsets,items,sprites)")
	jsonBuildsPrefix := flag.String("jbp", "", "Prefix to use for JSON builds (generated files under build will have this prefix)")

	flag.StringVar(savFile, "sav", "", "Path to the save file (.sav)")
	flag.StringVar(gbaFile, "gba", "", "Path to the GBA ROM file (.gba)")
	flag.StringVar(outputs, "output", "", "Comma-separated list of outputs to generate (e.g., species,evolutions,moves,learnsets,items,sprites)")
	flag.StringVar(jsonBuildsPrefix, "jsonBuildsPrefix", "", "Prefix to use for JSON builds (generated files under build will have this prefix)")

	flag.Parse()
	fmt.Printf("Save file path: %s\n", *savFile)
	fmt.Printf("GBA file path: %s\n", *gbaFile)
	fmt.Printf("Outputs: %s\n", *outputs)

	if *gbaFile == "" {
		log.Fatal("-g/-gba flag is required.")
	}

	if *jsonBuildsPrefix != "" {
		JSON_BUILDS_PREFIX = *jsonBuildsPrefix
	}

	gbaBytes, err := os.ReadFile(*gbaFile)
	if err != nil {
		log.Fatal(err)
	}
	g := gba.ParseGbaBytes(gbaBytes)
	if *savFile != "" {
		savBytes, err := os.ReadFile(*savFile)
		if err != nil {
			log.Fatal(err)
		}
		sav.ParseSavBytes(savBytes)
	}

	selectedOutputs := strings.Split(*outputs, ",")

	var wg sync.WaitGroup

	if slices.Contains(selectedOutputs, "evolutions") {
		buildTask(&wg, "Evolutions data", func() error {
			return SaveEvolutionsData("build/"+JSON_BUILDS_PREFIX+"evolutions.json", gba.Species)
		})
	}

	if slices.Contains(selectedOutputs, "items") {
		buildTask(&wg, "Items data", func() error {
			return SaveItemsData("build/"+JSON_BUILDS_PREFIX+"items.json", gba.Items[1:])
		})
	}

	if slices.Contains(selectedOutputs, "moves") {
		buildTask(&wg, "Moves data", func() error {
			return SaveMovesData("build/"+JSON_BUILDS_PREFIX+"moves.json", gba.Moves[1:])
		})
	}

	if slices.Contains(selectedOutputs, "species") {
		buildTask(&wg, "Species data", func() error {
			return SaveSpeciesData("build/"+JSON_BUILDS_PREFIX+"species.json", gba.Species[1:])
		})
	}

	if slices.Contains(selectedOutputs, "sprites") {
		buildTask(&wg, "Pokemon sprites", func() error {
			return SaveSpeciesSprites(gbaBytes, gba.Species[1:])
		})

		buildTask(&wg, "Pokemon icons", func() error {
			return SaveSpeciesIcons(gbaBytes, gba.Species[1:], g.IconPalettesTablePtr)
		})

		buildTask(&wg, "Item icons", func() error {
			return SaveItemsIcons(gbaBytes, gba.Items[1:])
		})
	}

	wg.Wait()
}

func buildTask(wg *sync.WaitGroup, taskName string, taskFunc func() error) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Printf("Saving %s...", taskName)
		if err := taskFunc(); err != nil {
			log.Fatal(err)
			return
		}
		log.Printf("Saved %s!", taskName)
	}()
}

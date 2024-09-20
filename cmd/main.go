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

func main() {
	savFile := flag.String("s", "", "Path to the save file (.sav)")
	gbaFile := flag.String("g", "", "Path to the GBA ROM file (.gba)")
	outputs := flag.String("o", "", "Comma-separated list of outputs to generate (e.g., species,moves,learnsets,items,sprites)")

	flag.StringVar(savFile, "sav", "", "Path to the save file (.sav)")
	flag.StringVar(gbaFile, "gba", "", "Path to the GBA ROM file (.gba)")
	flag.StringVar(outputs, "output", "", "Comma-separated list of outputs to generate (e.g., species,moves,learnsets,items,sprites)")

	flag.Parse()
	fmt.Printf("Save file path: %s\n", *savFile)
	fmt.Printf("GBA file path: %s\n", *gbaFile)
	fmt.Printf("Outputs: %s\n", *outputs)

	if *gbaFile == "" {
		log.Fatal("-g/-gba flag is required.")
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

	if slices.Contains(selectedOutputs, "species") {
		wg.Add(1)
		go func() {
			defer wg.Done()
			log.Println("Saving species data...")
			err := SaveSpeciesData("build/species_rhh.json", gba.Species[1:])
			if err != nil {
				log.Fatal(err)
				return
			}
			log.Println("Saved species data!")
		}()
	}

	if slices.Contains(selectedOutputs, "sprites") {
		wg.Add(1)
		go func() {
			defer wg.Done()
			log.Println("Saving Pokemon sprites...")
			err := SaveSpeciesSprites(gbaBytes, gba.Species[1:])
			if err != nil {
				log.Fatal(err)
				return
			}
			log.Println("Saved Pokemon sprites!")
		}()
		wg.Add(1)
		go func() {
			defer wg.Done()
			log.Println("Saving Pokemon icons...")
			err := SaveSpeciesIcons(gbaBytes, gba.Species[1:], g.IconPalettesTablePtr)
			if err != nil {
				log.Fatal(err)
				return
			}
			log.Println("Saved Pokemon icons!")
		}()
		wg.Add(1)
		go func() {
			defer wg.Done()
			log.Println("Saving Item icons...")
			err := SaveItemsIcons(gbaBytes, gba.Items[1:])
			if err != nil {
				log.Fatal(err)
				return
			}
			log.Println("Saved Item icons!")
		}()
	}

	wg.Wait()
}

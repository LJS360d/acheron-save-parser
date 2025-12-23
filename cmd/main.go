package main

import (
	"flag"
	"log"
	"os"
	"rom-parser/gba"
	"rom-parser/memmap"
	"rom-parser/sav"
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
	mapFile := flag.String("m", "", "Path to the build memory map file (.map)")
	outputs := flag.String("o", "", "Comma-separated list of outputs to generate (e.g., species,evolutions,moves,learnsets,items,sprites)")
	jsonBuildsPrefix := flag.String("jbp", "", "Prefix to use for JSON builds (generated files under build will have this prefix)")
	versionExtra := flag.String("vextra", "", "For loading extra GBAConfig for specific builds (e.g., acheron-emerald)")

	flag.StringVar(savFile, "sav", "", "Path to the save file (.sav)")
	flag.StringVar(gbaFile, "gba", "", "Path to the GBA ROM file (.gba)")
	flag.StringVar(mapFile, "map", "", "Path to the build memory map file (.map)")

	flag.StringVar(outputs, "output", "", "Comma-separated list of outputs to generate (e.g., species,evolutions,moves,learnsets,items,sprites)")
	flag.StringVar(jsonBuildsPrefix, "jsonBuildsPrefix", "", "Prefix to use for JSON builds (generated files under build will have this prefix)")

	flag.Parse()
	log.Printf("Save file path: %s\n", *savFile)
	log.Printf("GBA file path: %s\n", *gbaFile)
	log.Printf("Map file path: %s\n", *mapFile)
	log.Printf("Outputs: %s\n", *outputs)

	if *gbaFile == "" {
		log.Fatal("-g/-gba flag is required.")
	}

	if *mapFile != "" {
		regions, err := memmap.ParseMemoryMap(*mapFile)
		if err != nil {
			log.Fatal(err)
		}

		// Create a map to associate symbols with their corresponding offset pointers
		symbolOffsets := map[string]*int{
			"gAbilitiesInfo":    &gba.Config.AbilitiesOffset,
			"gItemsInfo":        &gba.Config.ItemsOffset,
			"gMovesInfo":        &gba.Config.MovesOffset,
			"gNaturesInfo":      &gba.Config.NaturesOffset,
			"gSpeciesInfo":      &gba.Config.SpeciesOffset,
			"gWildMonHeaders":   &gba.Config.WildEncountersOffset,
			"gMapGroups":        &gba.Config.MapGroupsOffset,
			"gRegionMapEntries": &gba.Config.RegionLocationsOffset,
			"gTrainers":         &gba.Config.TrainersOffset,
			"gTrainerSprites":   &gba.Config.TrainerSpritesOffset,
		}

		// Iterate through regions and update the offsets using the map
		for _, region := range regions {
			if offsetPtr, ok := symbolOffsets[region.Symbol]; ok {
				*offsetPtr = int(region.Address) - gba.POINTER_OFFSET
			}
		}
	}
	if *jsonBuildsPrefix != "" {
		JSON_BUILDS_PREFIX = *jsonBuildsPrefix
	}

	gbaBytes, err := os.ReadFile(*gbaFile)
	if err != nil {
		log.Fatal(err)
	}

	g := gba.LoadGbaData(gbaBytes, *versionExtra)
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

	if slices.Contains(selectedOutputs, "trainers") {
		buildTask(&wg, "Trainers data", func() error {
			return SaveJsonEncodable("build/"+JSON_BUILDS_PREFIX+"trainers.json", gba.Trainers)
		})
	}

	if slices.Contains(selectedOutputs, "items") {
		buildTask(&wg, "Items data", func() error {
			return SaveJsonEncodable("build/"+JSON_BUILDS_PREFIX+"items.json", gba.Items[1:])
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

	if slices.Contains(selectedOutputs, "learnsets") {
		buildTask(&wg, "Learnsets data", func() error {
			return SaveLearnsetsData(gbaBytes, "build/"+JSON_BUILDS_PREFIX+"learnsets.json", gba.Species[1:])
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

		buildTask(&wg, "Trainer sprites", func() error {
			return SaveTrainerSprites(gbaBytes, gba.ParseTrainerSpritesBytes(gba.Config.TrainerSpritesOffset, gba.Config.TrainerSpritesCount))
		})
	}

	if slices.Contains(selectedOutputs, "encounters") {
		buildTask(&wg, "Wild encounters", func() error {
			return SaveWildEncountersData("build/"+JSON_BUILDS_PREFIX+"wild_encounters.json", gba.WildEncounters)
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

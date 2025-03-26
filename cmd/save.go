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
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
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
)

const (
	ICON_PALETTES_COUNT     = 6
	PALETTE_SIZE            = 32
	COMPRESSED_PALETTE_SIZE = 40 // yes compressed is bigger than uncompressed, gamefreak probably had a reason to compress otherwise they just had the big stupid
)

func SaveItemsIcons(data []byte, items []*gba.ItemData) error {
	for i := 0; i < len(items); i++ {
		if items[i].IconPalettePtr == gba.NULL_POINTER {
			log.Println("[WARN] MISSING ITEM PALETTE FOR", i, items[i].Name)
			continue
		}
		if items[i].IconPicPtr == gba.NULL_POINTER {
			log.Println("[WARN] MISSING ITEM PIC FOR", i, items[i].Name)
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
		if s[i].IconSpritePtr == gba.NULL_POINTER {
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
		if s[i].FrontPicPtr == gba.NULL_POINTER {
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

// ------------------------------------------------------------

func SaveSpeciesData(filepath string, s []*gba.SpeciesData) error {
	getSpeciesFlags := func(mon *gba.SpeciesData) []string {
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
	getSpeciesIdentifier := func(mon *gba.SpeciesData, index uint16) string {
		speciesName := strings.ReplaceAll(
			strings.ReplaceAll(
				strings.ToUpper(utils.ToSnakeCase(mon.SpeciesName)),
				" ", "_"),
			"'", "_")
		flagMap := map[bool]string{
			mon.IsAlolanForm:      "ALOLAN",
			mon.IsGalarianForm:    "GALARIAN",
			mon.IsHisuianForm:     "HISUIAN",
			mon.IsPaldeanForm:     "PALDEAN",
			mon.IsTotem:           "TOTEM",
			mon.IsMegaEvolution:   "MEGA",
			mon.IsPrimalReversion: "PRIMAL",
			mon.IsUltraBurst:      "ULTRA",
			mon.IsGigantamax:      "GIGANTAMAX",
			mon.IsTeraForm:        "TERA",
		}
		for condition, flag := range flagMap {
			if condition {
				speciesName += "_" + flag
			}
		}
		return speciesName
	}
	return SaveJsonEncodable(filepath, utils.MapSlice(s,
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
				"eggGroups":    jsonconvert.MarshalSlice(utils.PruneDuplicates(mon.EggGroups[:])),
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
					mon.EvYieldSpAttack,
					mon.EvYieldSpDefense,
					mon.EvYieldSpeed,
				}),
				"expYield":    mon.ExpYield,
				"genderRatio": mon.GenderRatio,
				"growthRate":  mon.GrowthRate,
				"height":      mon.Height,
				"weight":      mon.Weight,
				"flags":       jsonconvert.MarshalSlice(getSpeciesFlags(mon)),
				"itemCommon":  mon.ItemCommon,
				"itemRare":    mon.ItemRare,
			}
		}))
}

func SaveMovesData(filepath string, m []*gba.MoveData) error {
	getMoveFlags := func(move *gba.MoveData) []string {
		flags := []string{}

		flagMap := map[bool]string{
			move.AlwaysCriticalHit:                 "ALWAYS_CRITICAL",
			move.AssistBanned:                      "ASSIST_BANNED",
			move.BallisticMove:                     "BALLISTIC",
			move.BitingMove:                        "BITING",
			move.CantUseTwice:                      "CANT_USE_TWICE",
			move.MakesContact:                      "CONTACT",
			move.CopycatBanned:                     "COPYCAT_BANNED",
			move.DamagesAirborne:                   "DAMAGES_AIRBORNE",
			move.DamagesUnderground:                "DAMAGES_UNDERGROUND",
			move.DamagesUnderwater:                 "DAMAGES_UNDERWATER",
			move.AirborneDoubleDamage:              "DOUBLE_DAMAGE_AIRBORNE",
			move.EncoreBanned:                      "ENCORE_BANNED",
			move.ForcePressure:                     "FORCE_PRESSURE",
			move.GravityBanned:                     "GRAVITY_BANNED",
			move.HealingMove:                       "HEALING",
			move.IgnoresTargetDefenseEvasionStages: "IGNORES_DEFENSE_EVASION",
			move.IgnoresKingsRock:                  "IGNORES_KINGS_ROCK",
			move.IgnoresProtect:                    "IGNORES_PROTECT",
			move.IgnoresSubstitute:                 "IGNORES_SUBSTITUTE",
			move.IgnoresTargetAbility:              "IGNORES_TARGET_ABILITY",
			move.IgnoreTypeIfFlyingAndUngrounded:   "IGNORES_TYPE_IF_FLYING_UNGROUNDED",
			move.InstructBanned:                    "INSTRUCT_BANNED",
			move.MagicCoatAffected:                 "MAGIC_COAT_AFFECTED",
			move.MeFirstBanned:                     "ME_FIRST_BANNED",
			move.MetronomeBanned:                   "METRONOME_BANNED",
			move.MimicBanned:                       "MIMIC_BANNED",
			move.MinimizeDoubleDamage:              "MINIMIZE_DOUBLE_DAMAGE",
			move.MirrorMoveBanned:                  "MIRROR_MOVE_BANNED",
			move.ParentalBondBanned:                "PARENTAL_BOND_BANNED",
			move.PowderMove:                        "POWDER",
			move.PulseMove:                         "PULSE",
			move.SketchBanned:                      "SKETCH_BANNED",
			move.SkyBattleBanned:                   "SKY_BATTLE_BANNED",
			move.SleepTalkBanned:                   "SLEEP_TALK_BANNED",
			move.SlicingMove:                       "SLICING",
			move.SnatchAffected:                    "SNATCH_AFFECTED",
			move.SoundMove:                         "SOUND",
			move.ThawsUser:                         "THAWS_USER",
		}

		for condition, flag := range flagMap {
			if condition {
				flags = append(flags, flag)
			}
		}

		return flags
	}
	return SaveJsonEncodable(filepath, utils.MapSlice(m,
		func(move *gba.MoveData, i int) JSON {
			additionalEffects := make([]JSON, move.NumAdditionalEffects)
			for j := 0; j < int(move.NumAdditionalEffects); j++ {
				additionalEffects[j] = JSON{
					"moveEffect":              move.AdditionalEffects[j].MoveEffect,
					"chance":                  move.AdditionalEffects[j].Chance,
					"self":                    move.AdditionalEffects[j].Self,
					"onChargeTurnOnly":        move.AdditionalEffects[j].OnChargeTurnOnly,
					"onlyIfTargetRaisedStats": move.AdditionalEffects[j].OnlyIfTargetRaisedStats,
				}
			}
			return JSON{
				"id":                i,
				"name":              move.Name,
				"description":       move.Description,
				"type":              move.Type,
				"category":          move.Category,
				"pp":                move.Pp,
				"power":             move.Power,
				"accuracy":          move.Accuracy,
				"effect":            move.Effect,
				"target":            move.Target,
				"priority":          move.Priority,
				"recoil":            move.Recoil,
				"criticalHitStage":  move.CriticalHitStage,
				"additionalEffects": jsonconvert.MarshalSlice(additionalEffects),
				"flags":             jsonconvert.MarshalSlice(getMoveFlags(move)),
			}
		}))
}

func SaveItemsData(filepath string, items []*gba.ItemData) error {
	return SaveJsonEncodable(filepath, items)
}

func SaveLearnsetsData(data []byte, filepath string, s []*gba.SpeciesData) error {
	return SaveJsonEncodable(filepath, utils.MapSlice(s,
		func(mon *gba.SpeciesData, i int) JSON {
			levelUpLearnset := parseLevelUpLearnset(data, mon.LevelUpLearnsetPtr)
			return JSON{
				"species": i,
				"levelUpLearnset": jsonconvert.MarshalSlice(
					utils.MapSlice(levelUpLearnset, func(move *LevelUpMove, i int) JSON {
						return JSON{
							"move":  move.Move,
							"level": move.Level,
						}
					})),
				"teachableLearnset": jsonconvert.MarshalSlice(parseLearnset(data, mon.TeachableLearnsetPtr)),
				"eggMoveLearnset":   jsonconvert.MarshalSlice(parseLearnset(data, mon.EggMoveLearnsetPtr)),
			}
		}))
}

func SaveEvolutionsData(filepath string, s []*gba.SpeciesData) error {
	var buildTree func(current uint16, tree *EvolutionTree, paths []EvolutionPath)
	buildTree = func(current uint16, tree *EvolutionTree, paths []EvolutionPath) {
		for _, path := range paths {
			if path.From == current {
				tree.Evolutions = append(tree.Evolutions, path)
				for _, outcome := range path.To {
					buildTree(outcome.Species, tree, paths)
				}
			}
		}
	}
	buildEvolutionTrees := func(speciesData []*gba.SpeciesData) []EvolutionTree {
		// Step 1: Collect all species and paths, grouped by the "from" species
		pathMap := make(map[uint16]*EvolutionPath)
		speciesSet := make(map[uint16]struct{})

		for familyID, species := range speciesData {
			familyID16 := uint16(familyID)
			for _, evo := range species.Evolutions {
				outcome := EvolutionOutcome{
					Species: evo.TargetSpecies,
					Methods: []EvolutionMethod{
						{
							Method: evo.Method,
							Clause: evo.Param,
						},
					},
				}

				// Group evolutions by 'from' species
				if path, exists := pathMap[familyID16]; exists {
					path.To = append(path.To, outcome)
				} else {
					pathMap[familyID16] = &EvolutionPath{
						From: familyID16,
						To:   []EvolutionOutcome{outcome},
					}
				}

				speciesSet[familyID16] = struct{}{}
				speciesSet[evo.TargetSpecies] = struct{}{}
			}
		}

		// Convert the map to a slice of EvolutionPath
		paths := make([]EvolutionPath, 0, len(pathMap))
		for _, path := range pathMap {
			paths = append(paths, *path)
		}

		// Step 2: Identify root species (those that don't appear as 'To')
		potentialRoots := make(map[uint16]struct{})
		for species := range speciesSet {
			potentialRoots[species] = struct{}{}
		}
		for _, path := range paths {
			for _, outcome := range path.To {
				delete(potentialRoots, outcome.Species)
			}
		}

		// Step 3: Build trees from roots
		var trees []EvolutionTree
		for root := range potentialRoots {
			tree := EvolutionTree{
				Family:     root,
				Evolutions: []EvolutionPath{},
			}
			buildTree(root, &tree, paths)
			trees = append(trees, tree)
		}

		return trees
	}
	return SaveJsonEncodable(filepath, buildEvolutionTrees(s))
}

func SaveWildEncountersData(filepath string, w []*gba.WildPokemonHeader) error {
	return SaveJsonEncodable(filepath, w)
}

// ------------------------------------------------------------

type EvolutionMethod struct {
	Method uint16 `json:"method"`
	Clause uint16 `json:"clause"`
}

type EvolutionOutcome struct {
	Species uint16            `json:"species"`
	Methods []EvolutionMethod `json:"methods"`
}

type EvolutionPath struct {
	From uint16             `json:"from"`
	To   []EvolutionOutcome `json:"to"`
}

type EvolutionTree struct {
	Family     uint16          `json:"family"`
	Evolutions []EvolutionPath `json:"evolutions"`
}

type LevelUpMove struct {
	Move  uint16
	Level uint16
}

func parseLevelUpLearnset(data []byte, offset uint32) []*LevelUpMove {
	learnset := make([]*LevelUpMove, 0)
	if offset == gba.NULL_POINTER || offset >= uint32(len(data)) {
		return learnset
	}
	const LEVEL_UP_MOVE_SIZE = 4
	for i := offset; i+LEVEL_UP_MOVE_SIZE < uint32(len(data)); i += LEVEL_UP_MOVE_SIZE {
		move := &LevelUpMove{
			Move:  binary.LittleEndian.Uint16(data[i:]),
			Level: binary.LittleEndian.Uint16(data[i+2:]),
		}
		if move.Move == 0xFFFF && move.Level == 0 {
			break
		}
		learnset = append(learnset, move)
	}
	return learnset
}

func parseLearnset(data []byte, offset uint32) []uint16 {
	learnset := make([]uint16, 0)
	if offset == gba.NULL_POINTER || offset >= uint32(len(data)) {
		return learnset
	}
	for i := offset; i+2 < uint32(len(data)); i += 2 {
		move := binary.LittleEndian.Uint16(data[i+2:])
		if move == 0xFFFF {
			break
		}
		learnset = append(learnset, move)
	}
	return learnset
}

var dirMutexes sync.Map

func SaveJsonEncodable(path string, data any) error {
	dir := filepath.Dir(path)

	// Get or create a mutex for the directory
	mutex, _ := dirMutexes.LoadOrStore(dir, &sync.Mutex{})
	mutex.(*sync.Mutex).Lock()         // Acquire the lock
	defer mutex.(*sync.Mutex).Unlock() // Release the lock

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("ERROR CREATING DIRECTORY: %w", err)
		}
	}

	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("ERROR CREATING FILE: %w", err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	encoder := json.NewEncoder(writer)
	encoder.SetIndent("", "  ") // for pretty printing
	if err := encoder.Encode(data); err != nil {
		return fmt.Errorf("ERROR ENCODING JSON: %w", err)
	}

	if err := writer.Flush(); err != nil {
		return fmt.Errorf("ERROR WRITING TO FILE: %w", err)
	}

	return nil
}

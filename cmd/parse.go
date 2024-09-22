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
	err = writer.Flush()
	if err != nil {
		return fmt.Errorf("ERROR WRITING TO FILE: %w", err)
	}
	return nil
}

func SaveMovesData(filepath string, m []*gba.MoveData) error {
	file, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("ERROR CREATING FILE: %w", err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	json.NewEncoder(writer).Encode(utils.MapSlice(m,
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
	err = writer.Flush()
	if err != nil {
		return fmt.Errorf("ERROR WRITING TO FILE: %w", err)
	}
	return nil
}

func SaveItemsData(filepath string, items []*gba.ItemData) error {
	file, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("ERROR CREATING FILE: %w", err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	json.NewEncoder(writer).Encode(utils.MapSlice(items,
		func(item *gba.ItemData, i int) JSON {
			return JSON{
				"id":              i,
				"name":            item.Name,
				"description":     item.Description,
				"category":        item.Pocket,
				"price":           item.Price,
				"secondaryId":     item.SecondaryId,
				"flingPower":      item.FlingPower,
				"holdEffectParam": item.HoldEffectParam,
				"holdEffect":      item.HoldEffect,
				"battleUsage":     item.BattleUsage,
				"importance":      item.Importance,
				"type":            item.Type,
			}
		}))
	err = writer.Flush()
	if err != nil {
		return fmt.Errorf("ERROR WRITING TO FILE: %w", err)
	}
	return nil
}

func SaveLearnsetsData(data []byte, filepath string, s []*gba.SpeciesData) error {
	file, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("ERROR CREATING FILE: %w", err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	json.NewEncoder(writer).Encode(utils.MapSlice(s,
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
	err = writer.Flush()
	if err != nil {
		return fmt.Errorf("ERROR WRITING TO FILE: %w", err)
	}
	return nil
}

func getSpeciesFlags(mon *gba.SpeciesData) []string {
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

func getMoveFlags(move *gba.MoveData) []string {
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

func getSpeciesIdentifier(mon *gba.SpeciesData, index uint16) string {
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

func SaveEvolutionsData(filepath string, s []*gba.SpeciesData) error {
	file, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("ERROR CREATING FILE: %w", err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	t := BuildEvolutionTrees(s)
	json.NewEncoder(writer).Encode(utils.MapSlice(t, func(e EvolutionTree, i int) JSON {
		return JSON{
			"family": e.Family,
			"evolutions": utils.MapSlice(e.Evolutions, func(ep EvolutionPath, i int) JSON {
				return JSON{
					"from": ep.From,
					"to": utils.MapSlice(ep.To, func(outcome EvolutionOutcome, i int) JSON {
						return JSON{
							"species": outcome.Species,
							"methods": utils.MapSlice(outcome.Methods, func(method EvolutionMethod, i int) JSON {
								return JSON{
									"method": method.Method,
									"clause": method.Clause,
								}
							}),
						}
					}),
				}
			}),
		}
	}))

	err = writer.Flush()
	if err != nil {
		return fmt.Errorf("ERROR WRITING TO FILE: %w", err)
	}
	return nil
}

type EvolutionMethod struct {
	Method uint16
	Clause uint16
}

type EvolutionOutcome struct {
	Species uint16
	Methods []EvolutionMethod
}

type EvolutionPath struct {
	From uint16
	To   []EvolutionOutcome
}

type EvolutionTree struct {
	Family     uint16
	Evolutions []EvolutionPath
}

func BuildEvolutionTrees(speciesData []*gba.SpeciesData) []EvolutionTree {
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

func buildTree(current uint16, tree *EvolutionTree, paths []EvolutionPath) {
	for _, path := range paths {
		if path.From == current {
			tree.Evolutions = append(tree.Evolutions, path)
			for _, outcome := range path.To {
				buildTree(outcome.Species, tree, paths)
			}
		}
	}
}

type LevelUpMove struct {
	Move  uint16
	Level uint16
}

func parseLevelUpLearnset(data []byte, offset uint32) []*LevelUpMove {
	learnset := make([]*LevelUpMove, 0)
	if offset == gba.BAD_POINTER || offset >= uint32(len(data)) {
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
	if offset == gba.BAD_POINTER || offset >= uint32(len(data)) {
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

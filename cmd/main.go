package main

import (
	"acheron-save-parser/gba"
	"acheron-save-parser/utils"
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"reflect"
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
	/* // To Quickly analyze sections with the GameFreak text encoding
	file, err := os.Create("output.txt")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	n := 128
	chunkSize := 1024
	section := gbaBytes[0x08690498-gba.POINTER_OFFSET : 0x08690498-gba.POINTER_OFFSET+28000]
	for start := 0; start < len(section); start += chunkSize {
		end := min(start+chunkSize, len(section)) // Calculate chunk boundaries
		binStr := utils.DecodeGFStringParallel(section[start:end], 32)
		for i := 0; i < len(binStr); i += n {
			chunk := binStr[i:min(i+n, len(binStr))]
			writer.WriteString(chunk + "\n")
			writer.Flush()
		}
	}
	err = writer.Flush()
	if err != nil {
		fmt.Println("Error writing to file:", err)
	} */
	gba.ParseGbaBytes(gbaBytes)
	file, err := os.Create("species.json")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	json.NewEncoder(writer).Encode(utils.MapSlice(gba.Species[1:], func(mon *gba.SpeciesData, i int) JSON {
		return JSON{
			"species": strings.ToUpper(mon.SpeciesName),
			"stats": processSlice(reflect.ValueOf([]uint8{
				mon.BaseHP,
				mon.BaseAttack,
				mon.BaseDefense,
				mon.BaseSpeed,
				mon.BaseSpAttack,
				mon.BaseSpDefense,
			})),
			"bst":        mon.Bst,
			"generation": mon.Generation,
			"types": processSlice(reflect.ValueOf(
				utils.PruneDuplicates(utils.MapSlice(mon.Types[:], func(t uint8, i int) string {
					return TYPES[t]
				})),
			)),
			"natDexNum": mon.NatDexNum,
			"abilities": processSlice(reflect.ValueOf(
				utils.PruneDuplicates(
					utils.FilterEmpty(
						utils.MapSlice(mon.Abilities[:2], func(a uint16, i int) string {
							if a == 0 {
								return ""
							}
							return gba.Abilities[a].Name
						}),
					),
				),
			)),
			// "bodyColor": mon.BodyColor,
			"catchRate":    mon.CatchRate,
			"categoryName": mon.CategoryName,
			"description":  mon.Description,
			"eggCycles":    mon.EggCycles,
			"eggGroups": processSlice(reflect.ValueOf(
				utils.PruneDuplicates(utils.MapSlice(mon.EggGroups[:], func(g uint8, i int) string {
					return EGG_GROUPS[g]
				})),
			)),
			// TODO formChangeTable
			"evYield": processSlice(reflect.ValueOf([]uint8{
				mon.EvYieldHP,
				mon.EvYieldAttack,
				mon.EvYieldDefense,
				mon.EvYieldSpeed,
				mon.EvYieldSpAttack,
				mon.EvYieldSpDefense,
			})),
			"expYield":    mon.ExpYield,
			"genderRatio": mon.GenderRatio,
			"growthRate":  mon.GrowthRate,
			"height":      mon.Height,
			"weight":      mon.Weight,
			"flags":       processSlice(reflect.ValueOf([]string{})),
			"itemCommon":  mon.ItemCommon,
			"itemRare":    mon.ItemRare,
		}
	}))
	err = writer.Flush()
	if err != nil {
		fmt.Println("Error writing to file:", err)
	}
}

type JSON = map[string]interface{}

// converts all exported members of a struct to a map with lowercase field names
func ToJsonMarshal(v interface{}) (JSON, error) {
	val := reflect.ValueOf(v)
	typ := reflect.TypeOf(v)

	// Handle pointers by dereferencing
	if typ.Kind() == reflect.Ptr {
		if val.IsNil() {
			return nil, fmt.Errorf("nil pointer received")
		}
		val = val.Elem()
		typ = typ.Elem()
	}

	// Check if the kind of value is a supported type
	if val.Kind() != reflect.Struct && val.Kind() != reflect.Slice && val.Kind() != reflect.Map {
		return nil, fmt.Errorf("input must be a struct, slice, or map, received: '%v'", val.Kind())
	}

	data := marshalWithLowercaseNames(val)
	return data, nil
}

// Helper function to recursively marshal a struct to a map with lowercase field names
func marshalWithLowercaseNames(v reflect.Value) JSON {
	result := make(JSON)

	if v.Kind() == reflect.Struct {
		typ := v.Type()
		for i := 0; i < v.NumField(); i++ {
			field := typ.Field(i)
			fieldValue := v.Field(i)

			// Ensure the field is exported and is not a method
			if field.PkgPath == "" && fieldValue.Kind() != reflect.Func {
				fieldName := field.Name
				lowercaseFieldName := strings.ToLower(fieldName[:1]) + fieldName[1:]

				if fieldValue.Kind() == reflect.Struct {
					result[lowercaseFieldName] = marshalWithLowercaseNames(fieldValue)
				} else if fieldValue.Kind() == reflect.Slice || fieldValue.Kind() == reflect.Array {
					result[lowercaseFieldName] = processSlice(fieldValue)
				} else if fieldValue.Kind() == reflect.Map {
					result[lowercaseFieldName] = processMap(fieldValue)
				} else {
					result[lowercaseFieldName] = fieldValue.Interface()
				}
			}
		}
	}
	return result
}

// Helper function to process slices
func processSlice(v reflect.Value) interface{} {
	var result []interface{}
	for i := 0; i < v.Len(); i++ {
		elem := v.Index(i)
		if elem.Kind() == reflect.Struct {
			result = append(result, marshalWithLowercaseNames(elem))
		} else if elem.Kind() == reflect.Slice || elem.Kind() == reflect.Array {
			result = append(result, processSlice(elem))
		} else if elem.Kind() == reflect.Map {
			result = append(result, processMap(elem))
		} else {
			result = append(result, elem.Interface())
		}
	}
	return result
}

// Helper function to process maps
func processMap(v reflect.Value) map[string]interface{} {
	result := make(map[string]interface{})
	for _, key := range v.MapKeys() {
		val := v.MapIndex(key)
		if val.Kind() == reflect.Struct {
			result[key.String()] = marshalWithLowercaseNames(val)
		} else if val.Kind() == reflect.Slice || val.Kind() == reflect.Array {
			result[key.String()] = processSlice(val)
		} else if val.Kind() == reflect.Map {
			result[key.String()] = processMap(val)
		} else {
			result[key.String()] = val.Interface()
		}
	}
	return result
}

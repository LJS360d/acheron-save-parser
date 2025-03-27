package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"
)

type Config struct {
	NewPath      string
	OldPath      string
	ObjectID     string
	Outputs      []string
	DeleteInputs bool
	Sort         bool
}

type JSONObject map[string]interface{}

func main() {
	config := parseFlags()

	newArr := readJSON(config.NewPath)
	oldArr := readJSON(config.OldPath)

	modifyAndCompare(newArr, oldArr, config.ObjectID)

	if config.Sort {
		sort.Slice(newArr, func(i, j int) bool {
			idA, okA := newArr[i][config.ObjectID]
			idB, okB := newArr[j][config.ObjectID]

			if !okA || !okB {
				return false
			}
			switch idA.(type) {
			case float64:
				if floatIdA, ok := idA.(float64); ok {
					if floatIdB, ok := idB.(float64); ok {
						return floatIdA < floatIdB
					}
				}
			case string:
				if stringIdA, ok := idA.(string); ok {
					if stringIdB, ok := idB.(string); ok {
						return stringIdA < stringIdB
					}
				}
			}
			return false
		})
	}
	writeOutputs(newArr, config.Outputs)

	if config.DeleteInputs {
		deleteInputs(config.NewPath, config.OldPath)
	}
}

func parseFlags() Config {
	var config Config

	flag.StringVar(&config.NewPath, "newPath", "", "New JSON file path")
	flag.StringVar(&config.NewPath, "new", "", "Alias for -newPath")
	flag.StringVar(&config.OldPath, "oldPath", "", "Old JSON file path")
	flag.StringVar(&config.OldPath, "old", "", "Alias for -oldPath")
	flag.StringVar(&config.ObjectID, "objectId", "id", "Object ID field name")
	flag.StringVar(&config.ObjectID, "id", "id", "Alias for -objectId")
	outputsStr := flag.String("outputs", "", "Comma-separated list of output file paths")
	flag.StringVar(outputsStr, "o", "", "Alias for -outputs")
	flag.BoolVar(&config.DeleteInputs, "deleteInputs", false, "Delete input files")
	flag.BoolVar(&config.DeleteInputs, "delete", false, "Alias for -deleteInputs")
	flag.BoolVar(&config.DeleteInputs, "d", false, "Alias for -deleteInputs")
	flag.BoolVar(&config.Sort, "sortOutputs", false, "Sort output by object ID")
	flag.BoolVar(&config.Sort, "sort", false, "Alias for -sortOutputs")
	flag.BoolVar(&config.Sort, "s", false, "Alias for -sortOutputs")

	flag.Parse()

	if config.NewPath == "" {
		config.NewPath = flag.Lookup("new").Value.String()
	}
	if config.OldPath == "" {
		config.OldPath = flag.Lookup("old").Value.String()
	}
	if *outputsStr == "" {
		*outputsStr = flag.Lookup("o").Value.String()
	}

	if config.NewPath == "" || config.OldPath == "" || *outputsStr == "" {
		fmt.Println("Missing required arguments: -new/-newPath, -old/-oldPath, and -o/-outputs are required.")
		os.Exit(1)
	}

	config.Outputs = strings.Split(*outputsStr, ",")
	return config
}

func readJSON(filePath string) []JSONObject {
	data, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("Error reading file:", err)
		os.Exit(1)
	}

	var result []JSONObject
	err = json.Unmarshal(data, &result)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		os.Exit(1)
	}

	return result
}

func modifyAndCompare(newArr, oldArr []JSONObject, objectID string) {
	for i := range newArr {
		oldObj := findObjectByID(oldArr, newArr[i][objectID], objectID)
		newArr[i]["old"] = oldObj
	}
}

func findObjectByID(arr []JSONObject, id interface{}, objectID string) JSONObject {
	for _, obj := range arr {
		if obj[objectID] == id {
			return obj
		}
	}
	return nil
}

func writeOutputs(arr []JSONObject, outputPaths []string) {
	for _, outputPath := range outputPaths {
		jsonData, err := json.MarshalIndent(arr, "", "  ")
		if err != nil {
			fmt.Println("Error marshalling JSON:", err)
			os.Exit(1)
		}

		err = os.WriteFile(outputPath, jsonData, 0644)
		if err != nil {
			fmt.Println("Error writing file:", err)
			os.Exit(1)
		}
	}

	fmt.Println("Comparison completed and files written to:", strings.Join(outputPaths, ", "))
}

func deleteInputs(newPath, oldPath string) {
	err := os.Remove(newPath)
	if err != nil {
		fmt.Println("Error deleting file:", err)
	}

	err = os.Remove(oldPath)
	if err != nil {
		fmt.Println("Error deleting file:", err)
	}

	fmt.Println("Deleted input files.")
}

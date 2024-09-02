package main

import (
	"acheron-save-parser/gba"
	"flag"
	"fmt"
	"log"
	"os"
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
	gbaData := gba.ParseGbaBytes(gbaBytes)
	fmt.Println(gbaData)
}

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
	/* 	// To quickly analyze sections with the GameFreak text encoding
	   	file, err := os.Create("output.txt")
	   	if err != nil {
	   		fmt.Println("Error creating file:", err)
	   		return
	   	}
	   	defer file.Close()
	   	binStr := sav.ReadStringParallel(gbaBytes[0xc880a4:0xc880a4+8000], 32)
	   	writer := bufio.NewWriter(file)
	   	// Write the string with newlines every n (size of the struct to analyze) characters
	   	n := 28
	   	for i := 0; i < len(binStr); i += n {
	   		chunk := binStr[i:min(i+n, len(binStr))]
	   		writer.WriteString(chunk + "\n")
	   	}
	   	err = writer.Flush()
	   	if err != nil {
	   		fmt.Println("Error writing to file:", err)
	   	} */
	gbaData := gba.ParseGbaBytes(gbaBytes)
	fmt.Println(gbaData)
}

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
	gbaData := gba.ParseGbaBytes(gbaBytes)
	fmt.Println(gbaData)
}

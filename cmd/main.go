package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s [save_file]\n", os.Args[0])
		os.Exit(1)
	}

	if _, err := os.Stat(os.Args[1]); os.IsNotExist(err) {
		fmt.Printf("ERROR: File %s not found\n", os.Args[1])
		os.Exit(1)
	}

	_, err := NewGen3Save(os.Args[1])
	if err != nil {
		fmt.Printf("Error loading save file: %v\n", err)
		os.Exit(1)
	}

	/* fmt.Printf("%s (%s)\n", sf.name, sf.gender)
	fmt.Println(sf.teamcount)
	for _, pkm := range sf.team {
		fmt.Printf("%s Level %d\n", pkm.name, pkm.level)
	} */
}

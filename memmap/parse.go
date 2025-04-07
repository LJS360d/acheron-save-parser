package memmap

import (
	"bufio"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type MemoryRegion struct {
	Section string
	Address uint32
	Size    uint32
	Object  string
	// Can be empty
	Symbol string
}

func ParseMemoryMap(filePath string) ([]MemoryRegion, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	regions := []MemoryRegion{}
	hexRegex := regexp.MustCompile(`(0x[0-9a-fA-F]+)`) // Unified hex regex

	var currentSection string
	var currentObject string

	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)

		if strings.HasPrefix(line, ".") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				currentSection = parts[0]
				addressStr := ""
				sizeStr := ""

				if hexRegex.MatchString(parts[1]) {
					addressStr = parts[1]
					if len(parts) > 2 && hexRegex.MatchString(parts[2]) {
						sizeStr = parts[2]
						if len(parts) > 3 {
							currentObject = strings.Join(parts[3:], " ")
						}
					} else {
						currentObject = strings.Join(parts[2:], " ")
					}
				} else {
					currentObject = strings.Join(parts[1:], " ")
				}

				address, _ := parseHex(addressStr)
				size, _ := parseHex(sizeStr)

				region := MemoryRegion{
					Section: currentSection,
					Address: address,
					Size:    size,
					Object:  currentObject,
				}
				regions = append(regions, region)
			}
		} else if hexRegex.MatchString(line) {
			// Handle lines with just address and symbol
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				addressStr := parts[0]
				symbol := strings.Join(parts[1:], " ")
				address, _ := parseHex(addressStr)

				region := MemoryRegion{
					Section: currentSection,
					Address: address,
					Size:    0, // Size might not be available here
					Object:  currentObject,
					Symbol:  symbol,
				}
				regions = append(regions, region)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return regions, nil
}

func parseHex(hexStr string) (uint32, error) {
	if hexStr == "" {
		return 0, nil
	}
	value, err := strconv.ParseUint(strings.TrimPrefix(hexStr, "0x"), 16, 32)
	return uint32(value), err
}

package main

import (
	"acheron-save-parser/gba"
	"acheron-save-parser/memmap"
	"fmt"
	"slices"
	"strings"
	"text/template"
)

// VersionDataSymbols defines the symbols we are interested in.
var VersionDataSymbols = []string{
	"gAbilitiesInfo",
	"gItemsInfo",
	"gMovesInfo",
	"gNaturesInfo",
	"gSpeciesInfo",
	"gWildMonHeaders",
	"gMapGroups",
	"gRegionMapEntries",
}

func main() {
	regions, err := memmap.ParseMemoryMap("test/memmap/acheron-emerald1.11.1.map")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	configData := gba.GetGbaConfig("")

	tmplString := `GbaConfig{
	AbilityInfoSize:           utils.AlignPointer({{ .AbilityInfoSize }}),
	ItemInfoSize:              utils.AlignPointer({{ .ItemInfoSize }}),
	MoveInfoSize:              utils.AlignPointer({{ .MoveInfoSize }}),
	MoveAdditionalEffectsSize: utils.AlignPointer({{ .MoveAdditionalEffectsSize }}),
	NatureInfoSize:            utils.AlignPointer({{ .NatureInfoSize }}),
	NaturesOffset:             {{ .gNaturesInfo }} - POINTER_OFFSET,
	NaturesCount:              {{ .NaturesCount }},
	PokemonNameLength:         {{ .PokemonNameLength }},
	SpeciesInfoSize:           utils.AlignPointer({{ .SpeciesInfoSize }}),
	WildEncountersOffset:      {{ .gWildMonHeaders }} - POINTER_OFFSET,
	WildEncounterSize:         utils.AlignPointer({{ .WildEncounterSize }}),
	WildPokemonSize:           utils.AlignPointer({{ .WildPokemonSize }}),
	MapGroupsOffset:           {{ .gMapGroups }} - POINTER_OFFSET,
	MapGroupSize:              utils.AlignPointer({{ .MapGroupSize }}),
	MapGroupCounts:            []int{ {{ range $i, $v := .MapGroupCounts }}{{if $i}}, {{end}}{{ $v }}{{ end }} },
	RegionLocationsOffset:     {{ .gRegionMapEntries }} - POINTER_OFFSET,
	RegionLocationSize:        utils.AlignPointer({{ .RegionLocationSize }}),
	RegionLocationsCount:      {{ .RegionLocationsCount }},
}`

	// Create a map to store the symbol addresses.
	symbolAddresses := make(map[string]any)
	for _, region := range regions {
		if slices.Contains(VersionDataSymbols, region.Symbol) {
			symbolAddresses[region.Symbol] = strings.ToLower(fmt.Sprintf("0x%08X", region.Address))
		}
	}
	// Update the symbolAddresses with the configData.
	symbolAddresses["AbilityInfoSize"] = configData.AbilityInfoSize
	symbolAddresses["ItemInfoSize"] = configData.ItemInfoSize
	symbolAddresses["MoveInfoSize"] = configData.MoveInfoSize
	symbolAddresses["MoveAdditionalEffectsSize"] = configData.MoveAdditionalEffectsSize
	symbolAddresses["NatureInfoSize"] = configData.NatureInfoSize
	symbolAddresses["NaturesCount"] = configData.NaturesCount
	symbolAddresses["PokemonNameLength"] = configData.PokemonNameLength
	symbolAddresses["SpeciesInfoSize"] = configData.SpeciesInfoSize
	symbolAddresses["WildEncounterSize"] = configData.WildEncounterSize
	symbolAddresses["WildPokemonSize"] = configData.WildPokemonSize
	symbolAddresses["MapGroupSize"] = configData.MapGroupSize
	symbolAddresses["MapGroupCounts"] = configData.MapGroupCounts
	symbolAddresses["RegionLocationSize"] = configData.RegionLocationSize
	symbolAddresses["RegionLocationsCount"] = configData.RegionLocationsCount
	// Parse the template.
	tmpl, err := template.New("gbaConfig").Parse(tmplString)
	if err != nil {
		fmt.Println("Error parsing template:", err)
		return
	}

	// Execute the template with the config data.
	var result strings.Builder
	err = tmpl.Execute(&result, symbolAddresses)
	if err != nil {
		fmt.Println("Error executing template:", err)
		return
	}

	// Print the resulting string.
	fmt.Println(result.String())
}

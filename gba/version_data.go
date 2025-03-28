package gba

import (
	"regexp"
	"strings"
)

type GbaConfig struct {
	Match string
	// abilities
	AbilityInfoSize int
	// items
	ItemInfoSize int
	// moves
	MoveInfoSize              int
	MoveAdditionalEffectsSize int
	// natures
	NatureInfoSize int
	NaturesOffset  int
	NaturesCount   int
	// species
	PokemonNameLength int
	SpeciesInfoSize   int
	// wild encounters
	WildEncounterSize    int
	WildPokemonSize      int
	WildEncountersOffset int
	WildEncounterCount   int
}

var (
	LatestVersion = "1.11.1"
	versionConfig = map[string]GbaConfig{
		LatestVersion: GbaConfig{
			AbilityInfoSize:           28, // ABILITY_INFO_SIZE_LENGTH16
			ItemInfoSize:              80,
			MoveInfoSize:              48,
			MoveAdditionalEffectsSize: 4,
			NatureInfoSize:            20,
			NaturesOffset:             0x0869797c - POINTER_OFFSET, // from .map
			NaturesCount:              25,
			PokemonNameLength:         13,
			SpeciesInfoSize:           260,
			WildEncounterSize:         20,
			WildPokemonSize:           4,
			WildEncountersOffset:      0x08edd49c - POINTER_OFFSET, // from .map
			WildEncounterCount:        130,                         // TODO
		},
		"1.10.x": GbaConfig{},
		"1.9.x": GbaConfig{
			AbilityInfoSize:           28, // ABILITY_INFO_SIZE_LENGTH16
			ItemInfoSize:              80,
			MoveInfoSize:              52,
			MoveAdditionalEffectsSize: 4,
			NatureInfoSize:            20,
			NaturesOffset:             0x08690498 - POINTER_OFFSET, // from .map
			NaturesCount:              25,
			PokemonNameLength:         13,
			SpeciesInfoSize:           216,
			WildEncounterSize:         20,
			WildPokemonSize:           4,
			WildEncountersOffset:      0x08de747c - POINTER_OFFSET, // from .map
			WildEncounterCount:        130,
		},
	}
)

// GetGbaConfig returns the GBA config for the given version
// If the version is not found, it returns the latest version config
// VersionConfig is a map of version strings to their corresponding GBA config
// 1.x.x is valid syntax, it will return the latest 1.x.x version config
// or the latest overall if a 1.x.x pattern is not matched
func GetGbaConfig(version string) GbaConfig {
	if config, ok := versionConfig[version]; ok {
		config.Match = version
		return config
	}

	if version == "" {
		conf := versionConfig[LatestVersion]
		conf.Match = LatestVersion
		return conf
	}

	// Pattern matching for x.x.x or 1.x.x
	parts := strings.Split(version, ".")
	if len(parts) >= 2 {
		var latestMatch string
		for v := range versionConfig {

			pattern := strings.ReplaceAll(v, "x", `\d+`)
			re := regexp.MustCompile(pattern)
			if re.MatchString(version) {
				if latestMatch == "" || v > latestMatch {
					latestMatch = v
				}
			}
		}

		if latestMatch != "" {
			conf := versionConfig[latestMatch]
			conf.Match = latestMatch
			return conf
		}
	}

	// Return latest version config if no match found
	conf := versionConfig[LatestVersion]
	conf.Match = LatestVersion
	return conf
}

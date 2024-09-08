package gba

import (
	"acheron-save-parser/utils"
	"encoding/binary"
)

type SpeciesData struct {
	BaseHP           uint8
	BaseAttack       uint8
	BaseDefense      uint8
	BaseSpeed        uint8
	BaseSpAttack     uint8
	BaseSpDefense    uint8
	Bst              int
	Generation       int
	Types            [2]uint8
	CatchRate        uint8
	ForceTeraType    uint8
	ExpYield         uint16
	EvYieldHP        uint8 // 2 bits
	EvYieldAttack    uint8 // 2 bits
	EvYieldDefense   uint8 // 2 bits
	EvYieldSpeed     uint8 // 2 bits
	EvYieldSpAttack  uint8 // 2 bits
	EvYieldSpDefense uint8 // 2 bits
	// Padding2                 uint8 // last 4 bits
	ItemCommon         uint16
	ItemRare           uint16
	GenderRatio        uint8
	EggCycles          uint8
	Friendship         uint8
	GrowthRate         uint8 // experience group
	EggGroups          [2]uint8
	Abilities          [NUM_ABILITY_SLOTS]uint16
	SafariZoneFleeRate uint8
	CategoryName       string // [13]uint8
	SpeciesName        string // [POKEMON_NAME_LENGTH + 1]uint8
	CryID              uint16
	NatDexNum          uint16
	Height             uint16 // in decimeters
	Weight             uint16 // in hectograms
	PokemonScale       uint16
	PokemonOffset      uint16
	TrainerScale       uint16
	TrainerOffset      uint16
	// 2 bytes of pointer boundary padding
	descriptionPtr        uint32
	Description           string
	BodyColor             uint8
	NoFlip                bool
	FrontAnimDelay        uint8
	FrontAnimID           uint8
	BackAnimID            uint8
	frontAnimFramesPtr    uint32
	frontPicPtr           uint32
	frontPicFemalePtr     uint32
	backPicPtr            uint32
	backPicFemalePtr      uint32
	palettePtr            uint32
	paletteFemalePtr      uint32
	shinyPalettePtr       uint32
	shinyPaletteFemalePtr uint32
	iconSpritePtr         uint32
	iconSpriteFemalePtr   uint32
	footprintPtr          uint32
	FrontPicSize          uint8
	FrontPicSizeFemale    uint8
	FrontPicYOffset       uint8
	BackPicSize           uint8
	BackPicSizeFemale     uint8
	BackPicYOffset        uint8
	IconPalIndex          uint8 // 3 bits
	IconPalIndexFemale    uint8 // 3 bits
	// Padding3           2 bits of padding
	EnemyMonElevation uint8
	IsLegendary       bool
	IsMythical        bool
	IsUltraBeast      bool
	IsParadox         bool
	IsTotem           bool
	IsMegaEvolution   bool
	IsPrimalReversion bool
	IsUltraBurst      bool
	IsGigantamax      bool
	IsTeraForm        bool
	IsAlolanForm      bool
	IsGalarianForm    bool
	IsHisuianForm     bool
	IsPaldeanForm     bool
	CannotBeTraded    bool
	AllPerfectIVs     bool
	DexForceRequired  bool
	TMIlliterate      bool
	IsFrontierBanned  bool
	// Padding4                 14 bits of padding
	levelUpLearnsetPtr       uint32
	teachableLearnsetPtr     uint32
	eggMoveLearnsetPtr       uint32
	evolutionsPtr            uint32
	formSpeciesIDTablePtr    uint32
	formChangeTablePtr       uint32
	OverworldData            [32]uint8 // TODO
	overworldPalettePtr      uint32
	overworldShinyPalettePtr uint32
}

const (
	NUM_ABILITY_SLOTS   = 3
	POKEMON_NAME_LENGTH = 12
	SPECIES_INFO_SIZE   = 216 // (08/09/2024) on latest commit in rrh/upcoming 212 is the new size (because pointer boundaries)
)

func ParseSpeciesInfoBytes(data []byte, offset int, count int) []*SpeciesData {
	species := make([]*SpeciesData, count)
	for i := 0; i < count; i++ {
		s := &SpeciesData{}
		s.new(data[offset+i*SPECIES_INFO_SIZE : offset+i*SPECIES_INFO_SIZE+SPECIES_INFO_SIZE])
		species[i] = s
		if s.descriptionPtr != BAD_POINTER {
			s.Description = utils.DecodePointerString(data, s.descriptionPtr)
		}
		s.Bst = int(s.BaseHP + s.BaseAttack + s.BaseDefense + s.BaseSpeed + s.BaseSpAttack + s.BaseSpDefense)
		s.Generation = getGenerationByDexNumber(int(s.NatDexNum))
	}
	return species
}

func getGenerationByDexNumber(dexNumber int) int {
	// start: 1, end: 151, generation: 1
	var gens []struct{ start, end, generation int } = []struct{ start, end, generation int }{
		{start: 1, end: 151, generation: 1},
		{start: 152, end: 251, generation: 2},
		{start: 252, end: 386, generation: 3},
		{start: 387, end: 493, generation: 4},
		{start: 494, end: 649, generation: 5},
		{start: 650, end: 721, generation: 6},
		{start: 722, end: 809, generation: 7},
		{start: 810, end: 905, generation: 8},
		{start: 906, end: 1025, generation: 9},
	}

	for _, gen := range gens {
		if gen.start <= dexNumber && dexNumber <= gen.end {
			return gen.generation
		}
	}
	return 0
}

func (s *SpeciesData) new(section []byte /* 216 bytes */) {
	s.BaseHP = section[0x0]
	s.BaseAttack = section[0x1]
	s.BaseDefense = section[0x2]
	s.BaseSpeed = section[0x3]
	s.BaseSpAttack = section[0x4]
	s.BaseSpDefense = section[0x5]
	s.Types[0] = section[0x6]
	s.Types[1] = section[0x7]
	s.CatchRate = section[0x8]
	s.ForceTeraType = section[0x9]
	s.ExpYield = binary.LittleEndian.Uint16(section[0xA:0xC])
	// first 2 bits
	s.EvYieldHP = section[0xC] & 0x3
	s.EvYieldAttack = section[0xC] & 0xC
	s.EvYieldDefense = section[0xC] & 0x30
	s.EvYieldSpeed = section[0xC] & 0xC0
	// first 2 bits
	s.EvYieldSpAttack = section[0xD] & 0x3
	s.EvYieldSpDefense = section[0xD] & 0xC
	// Padding2 last 4 bits
	s.ItemCommon = binary.LittleEndian.Uint16(section[0xE:0x10])
	s.ItemRare = binary.LittleEndian.Uint16(section[0x10:0x12])
	s.GenderRatio = section[0x12]
	s.EggCycles = section[0x13]
	s.Friendship = section[0x14]
	s.GrowthRate = section[0x15]
	s.EggGroups[0] = section[0x16]
	s.EggGroups[1] = section[0x17]
	s.Abilities[0] = binary.LittleEndian.Uint16(section[0x18:0x1A])
	s.Abilities[1] = binary.LittleEndian.Uint16(section[0x1A:0x1C])
	s.Abilities[2] = binary.LittleEndian.Uint16(section[0x1C:0x1E])
	s.SafariZoneFleeRate = section[0x1E]
	s.CategoryName = utils.DecodeGFString(section[0x1F:0x2C])
	s.SpeciesName = utils.DecodeGFString(section[0x2C:0x3A])
	s.CryID = binary.LittleEndian.Uint16(section[0x3A:0x3C])
	s.NatDexNum = binary.LittleEndian.Uint16(section[0x3C:0x3E])
	s.Height = binary.LittleEndian.Uint16(section[0x3E:0x40])
	s.Weight = binary.LittleEndian.Uint16(section[0x40:0x42])
	s.PokemonScale = binary.LittleEndian.Uint16(section[0x42:0x44])
	s.PokemonOffset = binary.LittleEndian.Uint16(section[0x44:0x46])
	s.TrainerScale = binary.LittleEndian.Uint16(section[0x46:0x48])
	s.TrainerOffset = binary.LittleEndian.Uint16(section[0x48:0x4A])
	// 2 bytes of padding pointer boundary padding
	s.descriptionPtr = binary.LittleEndian.Uint32(section[0x4C:0x50]) - POINTER_OFFSET
	// first 7 bits
	s.BodyColor = section[0x50] & 0x7F
	// last bit
	s.NoFlip = section[0x50]&0x80 == 1
	s.FrontAnimDelay = section[0x51]
	s.FrontAnimID = section[0x52]
	s.BackAnimID = section[0x53]
	s.frontAnimFramesPtr = binary.LittleEndian.Uint32(section[0x54:0x58]) - POINTER_OFFSET
	s.frontPicPtr = binary.LittleEndian.Uint32(section[0x58:0x5C]) - POINTER_OFFSET
	s.frontPicFemalePtr = binary.LittleEndian.Uint32(section[0x5C:0x60]) - POINTER_OFFSET
	s.backPicPtr = binary.LittleEndian.Uint32(section[0x60:0x64]) - POINTER_OFFSET
	s.backPicFemalePtr = binary.LittleEndian.Uint32(section[0x64:0x68]) - POINTER_OFFSET
	s.palettePtr = binary.LittleEndian.Uint32(section[0x68:0x6C]) - POINTER_OFFSET
	s.paletteFemalePtr = binary.LittleEndian.Uint32(section[0x6C:0x70]) - POINTER_OFFSET
	s.shinyPalettePtr = binary.LittleEndian.Uint32(section[0x70:0x74]) - POINTER_OFFSET
	s.shinyPaletteFemalePtr = binary.LittleEndian.Uint32(section[0x74:0x78]) - POINTER_OFFSET
	s.iconSpritePtr = binary.LittleEndian.Uint32(section[0x78:0x7C]) - POINTER_OFFSET
	s.iconSpriteFemalePtr = binary.LittleEndian.Uint32(section[0x7C:0x80]) - POINTER_OFFSET
	s.footprintPtr = binary.LittleEndian.Uint32(section[0x80:0x84]) - POINTER_OFFSET
	s.FrontPicSize = section[0x84]
	s.FrontPicSizeFemale = section[0x85]
	s.FrontPicYOffset = section[0x86]
	s.BackPicSize = section[0x87]
	s.BackPicSizeFemale = section[0x88]
	s.BackPicYOffset = section[0x89]
	// first 3 bits
	s.IconPalIndex = section[0x8A] & 0x7
	// next 3 bits
	s.IconPalIndexFemale = section[0x8A] & 0x38
	// last 2 bits are padding
	s.EnemyMonElevation = section[0x8B]
	s.IsLegendary = section[0x8C]&0x1 == 1
	s.IsMythical = section[0x8C]&0x2 == 2
	s.IsUltraBeast = section[0x8C]&0x4 == 4
	s.IsParadox = section[0x8C]&0x8 == 8
	s.IsTotem = section[0x8C]&0x10 == 0x10
	s.IsMegaEvolution = section[0x8C]&0x20 == 0x20
	s.IsPrimalReversion = section[0x8C]&0x40 == 0x40
	s.IsUltraBurst = section[0x8C]&0x80 == 0x80
	s.IsGigantamax = section[0x8D]&0x1 == 1
	s.IsTeraForm = section[0x8D]&0x2 == 2
	s.IsAlolanForm = section[0x8D]&0x4 == 4
	s.IsGalarianForm = section[0x8D]&0x8 == 8
	s.IsHisuianForm = section[0x8D]&0x10 == 0x10
	s.IsPaldeanForm = section[0x8D]&0x20 == 0x20
	s.CannotBeTraded = section[0x8D]&0x40 == 0x40
	s.AllPerfectIVs = section[0x8D]&0x80 == 0x80
	s.DexForceRequired = section[0x8E]&0x1 == 1
	s.TMIlliterate = section[0x8E]&0x2 == 2
	s.IsFrontierBanned = section[0x8E]&0x4 == 4
	// Padding4 14 bits of padding, it also takes the first bit of the next byte
	// 4 bytes of pointer boundary padding
	s.levelUpLearnsetPtr = binary.LittleEndian.Uint32(section[0x94:0x98]) - POINTER_OFFSET
	s.teachableLearnsetPtr = binary.LittleEndian.Uint32(section[0x98:0x9C]) - POINTER_OFFSET
	s.eggMoveLearnsetPtr = binary.LittleEndian.Uint32(section[0x9C:0xA0]) - POINTER_OFFSET
	s.evolutionsPtr = binary.LittleEndian.Uint32(section[0xA0:0xA4]) - POINTER_OFFSET
	s.formSpeciesIDTablePtr = binary.LittleEndian.Uint32(section[0xA4:0xA8]) - POINTER_OFFSET
	s.formChangeTablePtr = binary.LittleEndian.Uint32(section[0xA8:0xAC]) - POINTER_OFFSET
	// TODO
	// s.OverworldData = section[0xAC:0xB0]
	s.overworldPalettePtr = binary.LittleEndian.Uint32(section[0xB0:0xB4]) - POINTER_OFFSET
	s.overworldShinyPalettePtr = binary.LittleEndian.Uint32(section[0xB4:0xB8]) - POINTER_OFFSET
}

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
	FrontPicPtr           uint32
	frontPicFemalePtr     uint32
	backPicPtr            uint32
	backPicFemalePtr      uint32
	PalettePtr            uint32
	paletteFemalePtr      uint32
	ShinyPalettePtr       uint32
	shinyPaletteFemalePtr uint32
	IconSpritePtr         uint32
	IconSpriteFemalePtr   uint32
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
	LevelUpLearnset          []*LevelUpMove
	teachableLearnsetPtr     uint32
	TeachableLearnset        []uint16
	eggMoveLearnsetPtr       uint32
	EggMoveLearnset          []uint16
	evolutionsPtr            uint32
	Evolutions               []*Evolution
	formSpeciesIdTablePtr    uint32
	FormSpeciesIdTable       []uint16
	formChangeTablePtr       uint32
	FormChangeTable          []*FormChange
	OverworldData            [32]uint8 // TODO
	overworldPalettePtr      uint32
	overworldShinyPalettePtr uint32
}

const (
	NUM_ABILITY_SLOTS   = 3
	POKEMON_NAME_LENGTH = 12
	SPECIES_INFO_SIZE   = 212 // (08/09/2024) on latest commit in rrh/upcoming 212 is the new size (because pointer boundaries)
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
		s.Bst = int(s.BaseHP) + int(s.BaseAttack) + int(s.BaseDefense) + int(s.BaseSpeed) + int(s.BaseSpAttack) + int(s.BaseSpDefense)
		s.Generation = getGenerationByDexNumber(int(s.NatDexNum))
		s.FormSpeciesIdTable = parseFormSpeciesIdTable(data, s.formSpeciesIdTablePtr)
		s.FormChangeTable = parseFormChangeTable(data, s.formChangeTablePtr)
		s.Evolutions = parseEvolutions(data, s.evolutionsPtr)
		s.LevelUpLearnset = parseLevelUpLearnset(data, s.levelUpLearnsetPtr)
		s.TeachableLearnset = parseLearnset(data, s.teachableLearnsetPtr)
		s.EggMoveLearnset = parseLearnset(data, s.eggMoveLearnsetPtr)
	}
	return species
}

func parseFormSpeciesIdTable(data []byte, offset uint32) []uint16 {
	table := make([]uint16, 0)
	if offset == BAD_POINTER || offset >= uint32(len(data)) {
		return table
	}
	for i := offset; i+1 < uint32(len(data)); i += 2 {
		value := binary.LittleEndian.Uint16(data[i:])
		if value == 0xFFFF {
			break
		}
		table = append(table, value)
	}
	return table
}

type FormChange struct {
	Method        uint16
	TargetSpecies uint16
	Param1        uint16
	Param2        uint16
	Param3        uint16
}

func parseFormChangeTable(data []byte, offset uint32) []*FormChange {
	table := make([]*FormChange, 0)
	if offset == BAD_POINTER || offset >= uint32(len(data)) {
		return table
	}
	for i := offset; i+12 < uint32(len(data)); i += 12 {
		change := &FormChange{
			Method:        binary.LittleEndian.Uint16(data[i:]),
			TargetSpecies: binary.LittleEndian.Uint16(data[i+2:]),
			Param1:        binary.LittleEndian.Uint16(data[i+4:]),
			Param2:        binary.LittleEndian.Uint16(data[i+6:]),
			Param3:        binary.LittleEndian.Uint16(data[i+8:]),
		}
		if change.Method == 0x0000 {
			break
		}
		table = append(table, change)
	}
	return table
}

type Evolution struct {
	Method        uint16
	Param         uint16
	TargetSpecies uint16
}

func parseEvolutions(data []byte, offset uint32) []*Evolution {
	evolutions := make([]*Evolution, 0)
	if offset == BAD_POINTER || offset >= uint32(len(data)) {
		return evolutions
	}
	const EVO_SIZE = 8
	for i := offset; i+EVO_SIZE < uint32(len(data)); i += EVO_SIZE {
		evolution := &Evolution{
			Method:        binary.LittleEndian.Uint16(data[i:]),
			Param:         binary.LittleEndian.Uint16(data[i+2:]),
			TargetSpecies: binary.LittleEndian.Uint16(data[i+4:]),
		}
		evolutions = append(evolutions, evolution)
		if binary.LittleEndian.Uint16(data[i+EVO_SIZE:]) == 0xFFFF {
			break
		}
	}
	return evolutions
}

type LevelUpMove struct {
	Move  uint16
	Level uint16
}

func parseLevelUpLearnset(data []byte, offset uint32) []*LevelUpMove {
	learnset := make([]*LevelUpMove, 0)
	if offset == BAD_POINTER || offset >= uint32(len(data)) {
		return learnset
	}
	const LEVEL_UP_MOVE_SIZE = 4
	for i := offset; i+LEVEL_UP_MOVE_SIZE < uint32(len(data)); i += LEVEL_UP_MOVE_SIZE {
		move := &LevelUpMove{
			Move:  binary.LittleEndian.Uint16(data[i:]),
			Level: binary.LittleEndian.Uint16(data[i+2:]),
		}
		if move.Move == 0xFFFF && move.Level == 0 {
			break
		}
		learnset = append(learnset, move)
	}
	return learnset
}

func parseLearnset(data []byte, offset uint32) []uint16 {
	learnset := make([]uint16, 0)
	if offset == BAD_POINTER || offset >= uint32(len(data)) {
		return learnset
	}
	for i := offset; i+2 < uint32(len(data)); i += 2 {
		move := binary.LittleEndian.Uint16(data[i+2:])
		if move == 0xFFFF {
			break
		}
		learnset = append(learnset, move)
	}
	return learnset
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
	s.EvYieldHP = section[0xC] & 0b00000011
	s.EvYieldAttack = section[0xC] & 0b00001100 >> 2
	s.EvYieldDefense = section[0xC] & 0b00110000 >> 4
	s.EvYieldSpeed = section[0xC] & 0b11000000 >> 6
	// first 2 bits
	s.EvYieldSpAttack = section[0xD] & 0b00000011
	s.EvYieldSpDefense = section[0xD] & 0b00001100 >> 2
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
	s.FrontPicPtr = binary.LittleEndian.Uint32(section[0x58:0x5C]) - POINTER_OFFSET
	s.frontPicFemalePtr = binary.LittleEndian.Uint32(section[0x5C:0x60]) - POINTER_OFFSET
	s.backPicPtr = binary.LittleEndian.Uint32(section[0x60:0x64]) - POINTER_OFFSET
	s.backPicFemalePtr = binary.LittleEndian.Uint32(section[0x64:0x68]) - POINTER_OFFSET
	s.PalettePtr = binary.LittleEndian.Uint32(section[0x68:0x6C]) - POINTER_OFFSET
	s.paletteFemalePtr = binary.LittleEndian.Uint32(section[0x6C:0x70]) - POINTER_OFFSET
	s.ShinyPalettePtr = binary.LittleEndian.Uint32(section[0x70:0x74]) - POINTER_OFFSET
	s.shinyPaletteFemalePtr = binary.LittleEndian.Uint32(section[0x74:0x78]) - POINTER_OFFSET
	s.IconSpritePtr = binary.LittleEndian.Uint32(section[0x78:0x7C]) - POINTER_OFFSET
	s.IconSpriteFemalePtr = binary.LittleEndian.Uint32(section[0x7C:0x80]) - POINTER_OFFSET
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
	// 4 bytes of pointer boundary padding but only on builds with 216 size
	pointerPadding := SPECIES_INFO_SIZE - 212
	ptrsOffset := 144 + pointerPadding
	s.levelUpLearnsetPtr = binary.LittleEndian.Uint32(section[ptrsOffset:ptrsOffset+4]) - POINTER_OFFSET
	ptrsOffset += 4
	s.teachableLearnsetPtr = binary.LittleEndian.Uint32(section[ptrsOffset:ptrsOffset+4]) - POINTER_OFFSET
	ptrsOffset += 4
	s.eggMoveLearnsetPtr = binary.LittleEndian.Uint32(section[ptrsOffset:ptrsOffset+4]) - POINTER_OFFSET
	ptrsOffset += 4
	s.evolutionsPtr = binary.LittleEndian.Uint32(section[ptrsOffset:ptrsOffset+4]) - POINTER_OFFSET
	ptrsOffset += 4
	s.formSpeciesIdTablePtr = binary.LittleEndian.Uint32(section[ptrsOffset:ptrsOffset+4]) - POINTER_OFFSET
	ptrsOffset += 4
	s.formChangeTablePtr = binary.LittleEndian.Uint32(section[ptrsOffset:ptrsOffset+4]) - POINTER_OFFSET
	ptrsOffset += 4
	// TODO
	// s.OverworldData = section[0xAC:0xB0]
	s.overworldPalettePtr = binary.LittleEndian.Uint32(section[0xB0:0xB4]) - POINTER_OFFSET
	s.overworldShinyPalettePtr = binary.LittleEndian.Uint32(section[0xB4:0xB8]) - POINTER_OFFSET
}

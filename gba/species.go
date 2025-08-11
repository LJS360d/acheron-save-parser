package gba

import (
	"acheron-save-parser/utils"
	"encoding/binary"
)

const (
	NUM_ABILITY_SLOTS = 3
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
	Description           string
	BodyColor             uint8
	NoFlip                bool
	FrontAnimDelay        uint8
	FrontAnimID           uint8
	BackAnimID            uint8
	FrontAnimFramesPtr    uint32
	FrontPicPtr           uint32
	FrontPicFemalePtr     uint32
	BackPicPtr            uint32
	BackPicFemalePtr      uint32
	PalettePtr            uint32
	PaletteFemalePtr      uint32
	ShinyPalettePtr       uint32
	ShinyPaletteFemalePtr uint32
	IconSpritePtr         uint32
	IconSpriteFemalePtr   uint32
	FootprintPtr          uint32
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
	// Padding4                 14/11 bits of padding
	// shadow settings (1.11.x+)
	EnemyShadowXOffset  int8
	EnemyShadowYOffset  int8
	EnemyShadowSize     uint8 // 3 bits
	SuppressEnemyShadow bool
	// padding 12 bits
	LevelUpLearnsetPtr             uint32
	TeachableLearnsetPtr           uint32
	EggMoveLearnsetPtr             uint32
	Evolutions                     []*Evolution
	FormSpeciesIdTable             []uint16
	FormChangeTable                []*FormChange
	OverworldData                  [32]uint8 // TODO
	OverworldDataFemale            [32]uint8 // TODO
	overworldPalettePtr            uint32
	overworldShinyPalettePtr       uint32
	overworldPaletteFemalePtr      uint32
	overworldShinyPaletteFemalePtr uint32
}

func ParseSpeciesInfoBytes(offset int, count int) []*SpeciesData {
	species := make([]*SpeciesData, count)
	for i := 0; i < count; i++ {
		s := &SpeciesData{}
		s.loadFromDataSection(Data[offset+i*Config.SpeciesInfoSize : offset+i*Config.SpeciesInfoSize+Config.SpeciesInfoSize])
		species[i] = s
	}
	return species
}

func parseFormSpeciesIdTable(Data []byte, offset uint32) []uint16 {
	table := make([]uint16, 0)
	if offset == NULL_POINTER || offset >= uint32(len(Data)) {
		return table
	}
	for i := offset; i+1 < uint32(len(Data)); i += 2 {
		value := binary.LittleEndian.Uint16(Data[i:])
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

func parseFormChangeTable(Data []byte, offset uint32) []*FormChange {
	table := make([]*FormChange, 0)
	if offset == NULL_POINTER || offset >= uint32(len(Data)) {
		return table
	}
	for i := offset; i+12 < uint32(len(Data)); i += 12 {
		change := &FormChange{
			Method:        binary.LittleEndian.Uint16(Data[i:]),
			TargetSpecies: binary.LittleEndian.Uint16(Data[i+2:]),
			Param1:        binary.LittleEndian.Uint16(Data[i+4:]),
			Param2:        binary.LittleEndian.Uint16(Data[i+6:]),
			Param3:        binary.LittleEndian.Uint16(Data[i+8:]),
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

func parseEvolutions(Data []byte, offset uint32) []*Evolution {
	evolutions := make([]*Evolution, 0)
	if offset == NULL_POINTER || offset >= uint32(len(Data)) {
		return evolutions
	}
	const EVO_SIZE = 8
	for i := offset; i+EVO_SIZE < uint32(len(Data)); i += EVO_SIZE {
		evolution := &Evolution{
			Method:        binary.LittleEndian.Uint16(Data[i:]),
			Param:         binary.LittleEndian.Uint16(Data[i+2:]),
			TargetSpecies: binary.LittleEndian.Uint16(Data[i+4:]),
		}
		evolutions = append(evolutions, evolution)
		if binary.LittleEndian.Uint16(Data[i+EVO_SIZE:]) == 0xFFFF {
			break
		}
	}
	return evolutions
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

func (s *SpeciesData) loadFromDataSection(section []byte /* 260 bytes */) {
	s.BaseHP = section[0]
	s.BaseAttack = section[1]
	s.BaseDefense = section[2]
	s.BaseSpeed = section[3]
	s.BaseSpAttack = section[4]
	s.BaseSpDefense = section[5]
	// added util
	s.Bst = int(s.BaseHP) + int(s.BaseAttack) + int(s.BaseDefense) + int(s.BaseSpeed) + int(s.BaseSpAttack) + int(s.BaseSpDefense)
	s.Types[0] = section[6]
	s.Types[1] = section[7]
	s.CatchRate = section[8]
	s.ForceTeraType = section[9]
	s.ExpYield = binary.LittleEndian.Uint16(section[10:12])
	// first 2 bits
	s.EvYieldHP = section[12] & 0b00000011
	s.EvYieldAttack = section[12] & 0b00001100 >> 2
	s.EvYieldDefense = section[12] & 0b00110000 >> 4
	s.EvYieldSpeed = section[12] & 0b11000000 >> 6
	// first 2 bits
	s.EvYieldSpAttack = section[13] & 0b00000011
	s.EvYieldSpDefense = section[13] & 0b00001100 >> 2
	// Padding2 last 4 bits of 13th byte
	s.ItemCommon = binary.LittleEndian.Uint16(section[14:16])
	s.ItemRare = binary.LittleEndian.Uint16(section[16:18])
	s.GenderRatio = section[18]
	s.EggCycles = section[19]
	s.Friendship = section[20]
	s.GrowthRate = section[21]
	s.EggGroups[0] = section[22]
	s.EggGroups[1] = section[23]
	s.Abilities[0] = binary.LittleEndian.Uint16(section[24:])
	s.Abilities[1] = binary.LittleEndian.Uint16(section[26:])
	s.Abilities[2] = binary.LittleEndian.Uint16(section[28:])
	s.SafariZoneFleeRate = section[30]
	s.CategoryName = utils.DecodeGFString(section[31:44]) // 13 chars max
	// 13 chars max, technically build dependant
	s.SpeciesName = utils.DecodeGFString(section[44 : 44+Config.PokemonNameLength])
	var i = 45 + Config.PokemonNameLength
	s.CryID = binary.LittleEndian.Uint16(section[i:])
	s.NatDexNum = binary.LittleEndian.Uint16(section[i+2:])
	// added util
	s.Generation = getGenerationByDexNumber(int(s.NatDexNum))
	s.Height = binary.LittleEndian.Uint16(section[i+4:])
	s.Weight = binary.LittleEndian.Uint16(section[i+6:])
	// 2 bytes of padding ??
	s.PokemonScale = binary.LittleEndian.Uint16(section[i+10:])
	s.PokemonOffset = binary.LittleEndian.Uint16(section[i+12:])
	s.TrainerScale = binary.LittleEndian.Uint16(section[i+14:])
	s.TrainerOffset = binary.LittleEndian.Uint16(section[i+16:])
	descriptionPtr := binary.LittleEndian.Uint32(section[i+18:]) - POINTER_OFFSET
	if descriptionPtr != NULL_POINTER && descriptionPtr < uint32(len(Data)) {
		s.Description = utils.DecodePointerString(Data, descriptionPtr)
	}
	// first 7 bits, big endian
	s.BodyColor = section[i+22] & 0x7F
	// last bit
	s.NoFlip = section[i+22]&0x80 == 1
	s.FrontAnimDelay = section[i+23]
	s.FrontAnimID = section[i+24]
	s.BackAnimID = section[i+25]
	s.FrontAnimFramesPtr = binary.LittleEndian.Uint32(section[i+26:]) - POINTER_OFFSET
	s.FrontPicPtr = binary.LittleEndian.Uint32(section[i+30:]) - POINTER_OFFSET
	s.BackPicPtr = binary.LittleEndian.Uint32(section[i+34:]) - POINTER_OFFSET
	s.PalettePtr = binary.LittleEndian.Uint32(section[i+38:]) - POINTER_OFFSET
	s.ShinyPalettePtr = binary.LittleEndian.Uint32(section[i+42:]) - POINTER_OFFSET
	s.IconSpritePtr = binary.LittleEndian.Uint32(section[i+46:]) - POINTER_OFFSET
	// female graphics
	s.FrontPicFemalePtr = binary.LittleEndian.Uint32(section[i+50:]) - POINTER_OFFSET
	s.BackPicFemalePtr = binary.LittleEndian.Uint32(section[i+54:]) - POINTER_OFFSET
	s.PaletteFemalePtr = binary.LittleEndian.Uint32(section[i+58:]) - POINTER_OFFSET
	s.ShinyPaletteFemalePtr = binary.LittleEndian.Uint32(section[i+62:]) - POINTER_OFFSET
	s.IconSpriteFemalePtr = binary.LittleEndian.Uint32(section[i+66:]) - POINTER_OFFSET
	s.FootprintPtr = binary.LittleEndian.Uint32(section[i+68:]) - POINTER_OFFSET
	// 2 bytes of padding
	s.FrontPicSize = section[i+74]
	s.FrontPicYOffset = section[i+75]
	s.BackPicSize = section[i+76]
	s.BackPicYOffset = section[i+77]
	// female graphics
	s.FrontPicSizeFemale = section[i+78]
	s.BackPicSizeFemale = section[i+79]
	// first 3 bits
	s.IconPalIndex = section[i+80] & 0x7
	// next 3 bits
	s.IconPalIndexFemale = section[i+80] & 0x38
	// last 2 bits are padding
	s.EnemyMonElevation = section[i+81]
	s.IsLegendary = section[i+82]&0x1 == 1
	s.IsMythical = section[i+82]&0x2 == 2
	s.IsUltraBeast = section[i+82]&0x4 == 4
	s.IsParadox = section[i+82]&0x8 == 8
	s.IsTotem = section[i+82]&0x10 == 0x10
	s.IsMegaEvolution = section[i+82]&0x20 == 0x20
	s.IsPrimalReversion = section[i+82]&0x40 == 0x40
	s.IsUltraBurst = section[i+82]&0x80 == 0x80
	s.IsGigantamax = section[i+83]&0x1 == 1
	s.IsTeraForm = section[i+83]&0x2 == 2
	s.IsAlolanForm = section[i+83]&0x4 == 4
	s.IsGalarianForm = section[i+83]&0x8 == 8
	s.IsHisuianForm = section[i+83]&0x10 == 0x10
	s.IsPaldeanForm = section[i+83]&0x20 == 0x20
	s.CannotBeTraded = section[i+83]&0x40 == 0x40
	s.AllPerfectIVs = section[i+83]&0x80 == 0x80
	s.DexForceRequired = section[i+84]&0x1 == 1
	s.TMIlliterate = section[i+84]&0x2 == 2
	s.IsFrontierBanned = section[i+84]&0x4 == 4
	// Padding4 11 bits of padding
	s.EnemyShadowXOffset = int8(section[i+86])
	s.EnemyShadowYOffset = int8(section[i+87])
	// first 3 bits
	s.EnemyShadowSize = section[i+88] & 0x7
	// 4th bit
	s.SuppressEnemyShadow = section[i+88]&0x8 == 0x8
	// 12 bits of padding, so i+89 is unused
	// 2 bytes of compiler added padding for pointer alignment

	ptrsOffset := i + 90
	s.LevelUpLearnsetPtr = binary.LittleEndian.Uint32(section[ptrsOffset:ptrsOffset+4]) - POINTER_OFFSET
	ptrsOffset += 4
	s.TeachableLearnsetPtr = binary.LittleEndian.Uint32(section[ptrsOffset:ptrsOffset+4]) - POINTER_OFFSET
	ptrsOffset += 4
	s.EggMoveLearnsetPtr = binary.LittleEndian.Uint32(section[ptrsOffset:ptrsOffset+4]) - POINTER_OFFSET
	ptrsOffset += 4
	evolutionsPtr := binary.LittleEndian.Uint32(section[ptrsOffset:ptrsOffset+4]) - POINTER_OFFSET
	s.Evolutions = parseEvolutions(Data, evolutionsPtr)
	ptrsOffset += 4
	formSpeciesIdTablePtr := binary.LittleEndian.Uint32(section[ptrsOffset:ptrsOffset+4]) - POINTER_OFFSET
	s.FormSpeciesIdTable = parseFormSpeciesIdTable(Data, formSpeciesIdTablePtr)
	ptrsOffset += 4
	formChangeTablePtr := binary.LittleEndian.Uint32(section[ptrsOffset:ptrsOffset+4]) - POINTER_OFFSET
	s.FormChangeTable = parseFormChangeTable(Data, formChangeTablePtr)
	ptrsOffset += 4

	for i, v := range section[180:212] {
		s.OverworldData[i] = uint8(v)
	}
	for i, v := range section[212:244] {
		s.OverworldDataFemale[i] = uint8(v)
	}
	s.overworldPalettePtr = binary.LittleEndian.Uint32(section[244:]) - POINTER_OFFSET
	s.overworldShinyPalettePtr = binary.LittleEndian.Uint32(section[248:]) - POINTER_OFFSET
	s.overworldPaletteFemalePtr = binary.LittleEndian.Uint32(section[252:]) - POINTER_OFFSET
	s.overworldShinyPaletteFemalePtr = binary.LittleEndian.Uint32(section[256:]) - POINTER_OFFSET
}

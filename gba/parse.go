package gba

import (
	"encoding/binary"
	"fmt"
)

type GbaData struct {
	// --- GF Header ---
	RomEntryPoint uint32
	NintendoLogo  []byte // 156 bytes
	GameTitle     string // 12
	GameCode      string // 4
	MakerCode     string // 2
	UnitCode      uint8
	DeviceType    uint8
	// Reserved1    [7]byte
	SoftwareVersion uint8
	Checksum        uint8
	// Reserved2    [2]byte
	RamEntryPoint uint32
	BootMode      uint8 // should be 0x00
	SlaveID       uint8 // should be 0x00
	// unused26 [26]byte
	JoyBusEntryPoint uint32
	// --- GF .text.consts ---
	Language             uint32
	Version              uint32
	GameName             string // 32
	MonFrontPicsPtr      uint64
	MonBackPicsPtr       uint64
	MonNormalPalettesPtr uint64
	MonShinyPalettesPtr  uint64
	MonIconsPtr          uint64
	MonIconPaletteIdsPtr uint64
	PaletteTablesPtr     uint64
	MonSpeciesNamesPtr   uint64
	MoveNamesPtr         uint64
	DecorationsPtr       uint64
	FlagsOffset          uint64
	VarsOffset           uint64
	PokedexOffset        uint64
	Seen1Offset          uint64 // seen1 and seen2 are the same ptr
	Seen2Offset          uint64
	PokedexVar           uint8  // 0x46 (70)
	PokedexFlag          uint16 // 0x8E4 (2276)
	MysteryEventFlag     uint16 // 0x8AC (2220)
	PokedexCount         uint16 // 1025
	PlayerNameLength     uint8  // 7
	TrainerNameLength    uint8  // 10
	PokemonNameLength1   uint8  // 12
	PokemonNameLength2   uint8  // 12
	// --- RHH .text.consts ---
	RhhHeader      string // RHHEXP
	MajorVersion   uint8
	MinorVersion   uint8
	PatchVersion   uint8
	TaggedVersion  uint8
	MovesCount     uint16 // 848
	SpeciesCount   uint16 // 1524
	AbilitiesCount uint16 // 311
	AbilitiesPtr   uint32
	ItemsCount     uint16 // 831
	ItemNameLength uint8  // 20
}

func ParseGbaBytes(data []byte /* 33'554'432 Bytes */) *GbaData {
	fmt.Println(data[0xB2])
	g := &GbaData{
		// --- GF Header ---
		RomEntryPoint: binary.LittleEndian.Uint32(data[0x00:0x04]),
		NintendoLogo:  data[0x04:0xA0],
		GameTitle:     string(data[0xA0:0xAC]),
		GameCode:      string(data[0xAC:0xB0]),
		MakerCode:     string(data[0xB0:0xB2]),
		// data[0xB2] 1 byte of fixed value (96)
		UnitCode:   data[0xB3],
		DeviceType: data[0xB4],
		// Reserved1: data[0xB5:0xBC], // 7 bytes of reserved unused data
		SoftwareVersion: data[0xBC],
		Checksum:        data[0xBD],
		// Reserved2: data[0xBE:0xC0], // 2 bytes of reserved unused data
		RamEntryPoint: binary.LittleEndian.Uint32(data[0xC0:0xC4]),
		BootMode:      data[0xC4],
		SlaveID:       data[0xC5],
		// unused26: data[0xC6:0xD2], // 26 bytes of reserved unused data
		JoyBusEntryPoint: binary.LittleEndian.Uint32(data[0xD2:0xD6]),
		// gap between 214 - 256
		// --- GF .text.consts ---
		Version:  binary.LittleEndian.Uint32(data[0x100:0x104]),
		Language: binary.LittleEndian.Uint32(data[0x104:0x108]),
		GameName: string(data[0x108:0x11F]),
		// pointer madness
		MonFrontPicsPtr:      binary.LittleEndian.Uint64(data[0x11F:0x127]),
		MonBackPicsPtr:       binary.LittleEndian.Uint64(data[0x127:0x133]),
		MonNormalPalettesPtr: binary.LittleEndian.Uint64(data[0x133:0x13B]),
		MonShinyPalettesPtr:  binary.LittleEndian.Uint64(data[0x13B:0x143]),
		MonIconsPtr:          binary.LittleEndian.Uint64(data[0x143:0x14B]),
		MonIconPaletteIdsPtr: binary.LittleEndian.Uint64(data[0x14B:0x153]),
		PaletteTablesPtr:     binary.LittleEndian.Uint64(data[0x153:0x15B]),
		MonSpeciesNamesPtr:   binary.LittleEndian.Uint64(data[0x15B:0x163]),
		MoveNamesPtr:         binary.LittleEndian.Uint64(data[0x163:0x16B]),
		DecorationsPtr:       binary.LittleEndian.Uint64(data[0x16B:0x173]),
		FlagsOffset:          binary.LittleEndian.Uint64(data[0x173:0x17B]),
		VarsOffset:           binary.LittleEndian.Uint64(data[0x17B:0x183]),
		PokedexOffset:        binary.LittleEndian.Uint64(data[0x183:0x18B]),
		Seen1Offset:          binary.LittleEndian.Uint64(data[0x18B:0x193]), // seen1 and seen2 are the same ptr
		Seen2Offset:          binary.LittleEndian.Uint64(data[0x193:0x19B]),
		// end of pointer madness
		PokedexVar:         data[0x19B],
		PokedexFlag:        binary.LittleEndian.Uint16(data[0x19C:0x19E]),
		MysteryEventFlag:   binary.LittleEndian.Uint16(data[0x19D:0x19F]),
		PokedexCount:       binary.LittleEndian.Uint16(data[0x19D:0x19F]),
		PlayerNameLength:   data[0x200],
		TrainerNameLength:  data[0x201],
		PokemonNameLength1: data[0x202],
		PokemonNameLength2: data[0x203],
		// --- RHH .text.consts ---
		RhhHeader:      string(data[0x204:0x20A]),
		MajorVersion:   data[0x20A],
		MinorVersion:   data[0x20B],
		PatchVersion:   data[0x20C],
		TaggedVersion:  data[0x20D],
		MovesCount:     binary.LittleEndian.Uint16(data[0x20E:0x210]),
		SpeciesCount:   binary.LittleEndian.Uint16(data[0x210:0x212]),
		AbilitiesCount: binary.LittleEndian.Uint16(data[0x212:0x214]),
		AbilitiesPtr:   binary.LittleEndian.Uint32(data[0x214:0x218]),
		ItemsCount:     binary.LittleEndian.Uint16(data[0x218:0x21A]),
		ItemNameLength: data[0x21A],
	}
	/* abilities := ParseAbilitiesBytes(data, int(g.AbilitiesOffset), int(g.AbilitiesCount))
	fmt.Println(abilities) */
	return g
}

func ParseAbilitiesBytes(data []byte /* 7464 Bytes */, offset int, count int) []*AbilityData {
	abilities := make([]*AbilityData, count)
	for i := 0; i < count; i++ {
		a := &AbilityData{}
		a.new(data[offset+i*32 : offset+i*32+32])
		abilities = append(abilities, a)
	}
	return abilities
}

type AbilityData struct {
	Name              string // 20 bytes
	DescriptionPtr    uint32
	aiRating          int8
	cantBeCopied      bool
	cantBeSwapped     bool
	cantBeTraced      bool
	cantBeSuppressed  bool
	cantBeOverwritten bool
	breakable         bool
	failsOnImposter   bool
}

func (a *AbilityData) new(section []byte) {
	a.Name = string(section[0:20])
	a.DescriptionPtr = binary.LittleEndian.Uint32(section[20:24])
	a.aiRating = int8(section[24])
	a.cantBeCopied = section[25]&0x1 == 1
	a.cantBeSwapped = section[25]&0x2 == 1
	a.cantBeTraced = section[25]&0x4 == 1
	a.cantBeSuppressed = section[25]&0x8 == 1
	a.cantBeOverwritten = section[25]&0x10 == 1
	a.breakable = section[25]&0x20 == 1
	a.failsOnImposter = section[25]&0x40 == 1
}

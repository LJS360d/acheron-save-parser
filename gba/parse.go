package gba

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"strconv"
	"strings"
)

const (
	POINTER_OFFSET = 0x08000000
	BAD_POINTER    = 0x0f8000000
)

var (
	Abilities []*AbilityData
	Species   []*SpeciesData
	Items     []*ItemData
	Natures   []*NatureData
	Moves     []*MoveData
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
	Language uint32
	Version  uint32
	GameName string // 32
	// MonFrontPicsPtr      uint32
	// MonBackPicsPtr       uint32
	// MonNormalPalettesPtr uint32
	// MonShinyPalettesPtr  uint32
	// MonIconsPtr          uint32
	// MonIconPaletteIdsPtr uint32
	PaletteTablesPtr uint32
	// MonSpeciesNamesPtr uint32
	// MoveNamesPtr       uint32
	DecorationsPtr     uint32
	FlagsOffset        uint32
	VarsOffset         uint32
	PokedexOffset      uint32
	Seen1Offset        uint32 // seen1 and seen2 are the same ptr
	Seen2Offset        uint32
	PokedexVar         uint32 // 0x46 (70)
	PokedexFlag        uint32 // 0x8E40 (2276)
	MysteryEventFlag   uint32 // 0x8AC (2220)
	PokedexCount       uint32 // 1025
	PlayerNameLength   uint8  // 7
	TrainerNameLength  uint8  // 10
	PokemonNameLength1 uint8  // 12
	PokemonNameLength2 uint8  // 12
	// 12 bytes of unknown use
	// 3 bytes of padding
	SaveBlock2Size           uint32 // 0xF2C (3884)
	SaveBlock1Size           uint32 // 0x3CD0 (15568)
	PartyCountOffset         uint32
	PartyOffset              uint32
	WarpFlagsOffset          uint32
	TrainerIdOffset          uint32
	PlayerNameOffset         uint32
	PlayerGenderOffset       uint32
	FrontierStatusOffset     uint32
	FrontierStatusOffset2    uint32
	ExternalEventFlagsOffset uint32
	ExternalEventDataOffset  uint32
	// unk18 uint8
	SpeciesInfoPtr uint32
	// AbilityNamesPtr        uint32
	// AbilityDescriptionsPtr uint64
	ItemsPtr           uint32
	MovesPtr           uint32
	BallGfxPtr         uint32
	BallPallettesPtr   uint32
	GcnLinkFlagsOffset uint32
	GameClearFlag      uint32
	RibbonFlag         uint32
	BagItemsCount      uint8
	BagKeyItemsCount   uint8
	BagPokeballsCount  uint8
	BagTMHMsCount      uint8
	BagBerriesCount    uint8
	PcItemsCount       uint8
	PcItemsOffset      uint32
	GiftRibbonsOffset  uint32
	EnigmaBerryOffset  uint32
	EnigmaBerrySize    uint32
	MoveDescription    uint8 // always NULL ((void *)0)
	// unknown20 			uint32
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
		JoyBusEntryPoint: binary.LittleEndian.Uint32(data[0xE0:0xE4]),
		// gap between 214 - 256
		// --- GF .text.consts ---
		Version:  binary.LittleEndian.Uint32(data[0x100:0x104]),
		Language: binary.LittleEndian.Uint32(data[0x104:0x108]),
		GameName: string(data[0x108:0x11F]),
		/* MonFrontPicsPtr:      binary.LittleEndian.Uint64(data[0x11F:0x127]),
		MonBackPicsPtr:       binary.LittleEndian.Uint64(data[0x127:0x133]),
		MonNormalPalettesPtr: binary.LittleEndian.Uint64(data[0x133:0x13B]),
		MonShinyPalettesPtr:  binary.LittleEndian.Uint64(data[0x13B:0x143]),
		MonIconsPtr:          binary.LittleEndian.Uint64(data[0x143:0x14B]),
		MonIconPaletteIdsPtr: binary.LittleEndian.Uint32(data[0x140:0x145]), */
		PaletteTablesPtr: binary.LittleEndian.Uint32(data[0x140:0x145]),
		// MonSpeciesNamesPtr: binary.LittleEndian.Uint32(data[0x145:0x14A]),
		// MoveNamesPtr:       binary.LittleEndian.Uint32(data[0x163:0x16B]),
		DecorationsPtr:     binary.LittleEndian.Uint32(data[0x14C:0x150]) - POINTER_OFFSET,
		FlagsOffset:        binary.LittleEndian.Uint32(data[0x150:0x154]), /* - POINTER_OFFSET */
		VarsOffset:         binary.LittleEndian.Uint32(data[0x154:0x158]), /* - POINTER_OFFSET */
		PokedexOffset:      binary.LittleEndian.Uint32(data[0x158:0x15C]), /* - POINTER_OFFSET */
		Seen1Offset:        binary.LittleEndian.Uint32(data[0x15C:0x160]), /* - POINTER_OFFSET */ // seen1 and seen2 are the same ptr
		Seen2Offset:        binary.LittleEndian.Uint32(data[0x160:0x164]), /* - POINTER_OFFSET */
		PokedexVar:         binary.LittleEndian.Uint32(data[0x164:0x168]),
		PokedexFlag:        binary.LittleEndian.Uint32(data[0x168:0x16C]),
		MysteryEventFlag:   binary.LittleEndian.Uint32(data[0x16C:0x170]),
		PokedexCount:       binary.LittleEndian.Uint32(data[0x170:0x174]),
		PlayerNameLength:   data[0x174],
		TrainerNameLength:  data[0x175],
		PokemonNameLength1: data[0x176],
		PokemonNameLength2: data[0x177],
		// ---
		// 12 bytes of unknown use
		// 3 bytes of padding
		SaveBlock2Size:           binary.LittleEndian.Uint32(data[0x188:0x18C]),
		SaveBlock1Size:           binary.LittleEndian.Uint32(data[0x18C:0x190]),
		PartyCountOffset:         binary.LittleEndian.Uint32(data[0x190:0x194]), /* - POINTER_OFFSET */
		PartyOffset:              binary.LittleEndian.Uint32(data[0x194:0x198]), /* - POINTER_OFFSET */
		WarpFlagsOffset:          binary.LittleEndian.Uint32(data[0x198:0x19C]), /* - POINTER_OFFSET */
		TrainerIdOffset:          binary.LittleEndian.Uint32(data[0x19C:0x1A0]), /* - POINTER_OFFSET */
		PlayerNameOffset:         binary.LittleEndian.Uint32(data[0x1A0:0x1A4]), /* - POINTER_OFFSET */
		PlayerGenderOffset:       binary.LittleEndian.Uint32(data[0x1A4:0x1A8]), /* - POINTER_OFFSET */
		FrontierStatusOffset:     binary.LittleEndian.Uint32(data[0x1A8:0x1AC]), /* - POINTER_OFFSET */
		FrontierStatusOffset2:    binary.LittleEndian.Uint32(data[0x1AC:0x1B0]), /* - POINTER_OFFSET */
		ExternalEventFlagsOffset: binary.LittleEndian.Uint32(data[0x1B0:0x1B4]), /* - POINTER_OFFSET */
		ExternalEventDataOffset:  binary.LittleEndian.Uint32(data[0x1B4:0x1B8]), /* - POINTER_OFFSET */
		// unk18: data[0x1B8:0x1BC],
		SpeciesInfoPtr: binary.LittleEndian.Uint32(data[0x1BC:0x1C0]) - POINTER_OFFSET,
		// some padding that idk and does not matter
		// AbilityNamesPtr:        binary.LittleEndian.Uint32(data[0x1BC:0x1C0]) - POINTER_OFFSET,
		// AbilityDescriptionsPtr: binary.LittleEndian.Uint64(data[0x1C0:0x1C9]) - POINTER_OFFSET,
		ItemsPtr:           binary.LittleEndian.Uint32(data[0x1C8:0x1CC]) - POINTER_OFFSET,
		MovesPtr:           binary.LittleEndian.Uint32(data[0x1CC:0x1D0]) - POINTER_OFFSET,
		BallGfxPtr:         binary.LittleEndian.Uint32(data[0x1D0:0x1D4]) - POINTER_OFFSET,
		BallPallettesPtr:   binary.LittleEndian.Uint32(data[0x1D4:0x1D8]) - POINTER_OFFSET,
		GcnLinkFlagsOffset: binary.LittleEndian.Uint32(data[0x1D8:0x1DC]), /* - POINTER_OFFSET */
		//
		GameClearFlag:     binary.LittleEndian.Uint32(data[0x1DC:0x1E0]),
		RibbonFlag:        binary.LittleEndian.Uint32(data[0x1E0:0x1E4]),
		BagItemsCount:     data[0x1E4],
		BagKeyItemsCount:  data[0x1E5],
		BagPokeballsCount: data[0x1E6],
		BagTMHMsCount:     data[0x1E7],
		BagBerriesCount:   data[0x1E8],
		PcItemsCount:      data[0x1E9],
		PcItemsOffset:     binary.LittleEndian.Uint32(data[0x1E9:0x1ED]) - POINTER_OFFSET,
		GiftRibbonsOffset: binary.LittleEndian.Uint32(data[0x1ED:0x1F1]) - POINTER_OFFSET,
		EnigmaBerryOffset: binary.LittleEndian.Uint32(data[0x1F1:0x1F5]) - POINTER_OFFSET,
		EnigmaBerrySize:   binary.LittleEndian.Uint32(data[0x1F5:0x1F9]) - POINTER_OFFSET,
		MoveDescription:   data[0x1F9],
		// unknown20: binary.LittleEndian.Uint32(data[0x1FD:0x204]), // 0xFFFFFFFF
		// X bytes of padding
		// --- RHH .text.consts ---
		RhhHeader:      string(data[0x204:0x20A]),
		MajorVersion:   data[0x20A],
		MinorVersion:   data[0x20B],
		PatchVersion:   data[0x20C],
		TaggedVersion:  data[0x20D],
		MovesCount:     binary.LittleEndian.Uint16(data[0x20E:0x210]),
		SpeciesCount:   binary.LittleEndian.Uint16(data[0x210:0x212]),
		AbilitiesCount: binary.LittleEndian.Uint16(data[0x212:0x214]),
		AbilitiesPtr:   binary.LittleEndian.Uint32(data[0x214:0x218]) - POINTER_OFFSET,
		ItemsCount:     binary.LittleEndian.Uint16(data[0x218:0x21A]),
		ItemNameLength: data[0x21A],
	}
	a := ParseAbilitiesBytes(data, int(g.AbilitiesPtr), int(g.AbilitiesCount))
	Abilities = a
	s := ParseSpeciesInfoBytes(data, int(g.SpeciesInfoPtr), int(g.SpeciesCount))
	Species = s
	pal, _ := LoadPalette("images/ivysaur.pal")
	Save4bppImageBytes(data[s[2].frontPicPtr:s[2].frontPicPtr+2048], "images/ivysaur", pal, 64, 64)
	i := ParseItemsInfoBytes(data, int(g.ItemsPtr), int(g.ItemsCount))
	Items = i
	m := ParseMovesInfoBytes(data, int(g.MovesPtr), int(g.MovesCount))
	Moves = m
	NaturesPtr := 0x08690498 - POINTER_OFFSET // (08/09/2024) on latest commit in rrh/upcoming 0x0869797c is the new offset
	NaturesCount := 25
	n := ParseNaturesInfoBytes(data, int(NaturesPtr), int(NaturesCount))
	Natures = n
	return g
}

func LoadPalette(filename string) ([]color.Color, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var palette color.Palette
	scanner := bufio.NewScanner(file)

	lineNumber := 0
	for scanner.Scan() {
		line := scanner.Text()
		lineNumber++

		// Skip the first three header lines (JASC-PAL, version, number of colors).
		if lineNumber <= 3 {
			continue
		}

		// Split the RGB components.
		parts := strings.Split(line, " ")
		if len(parts) != 3 {
			continue
		}

		// Parse the R, G, B values.
		r, _ := strconv.Atoi(parts[0])
		g, _ := strconv.Atoi(parts[1])
		b, _ := strconv.Atoi(parts[2])

		// Append the color to the palette.
		palette = append(palette, color.RGBA{uint8(r), uint8(g), uint8(b), 255})
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return palette, nil
}

func DecompressLZ77(input []byte) ([]byte, error) {
	if len(input) < 4 || input[0] != 0x10 {
		return nil, fmt.Errorf("invalid LZ77 compressed data")
	}

	var decompressed bytes.Buffer
	srcOffset := 4
	length := int(binary.LittleEndian.Uint32(input) >> 8)

	for srcOffset < len(input) {
		// Read the flag byte (8 blocks).
		flag := input[srcOffset]
		srcOffset++

		for i := 0; i < 8; i++ {
			if srcOffset >= len(input) || decompressed.Len() >= length {
				break
			}

			// If bit is 0, it's raw byte data.
			if (flag & 0x80) == 0 {
				decompressed.WriteByte(input[srcOffset])
				srcOffset++
			} else {
				// Else, it's a reference to earlier data (compressed).
				if srcOffset+1 >= len(input) {
					return nil, fmt.Errorf("unexpected end of input")
				}

				byte1 := input[srcOffset]
				byte2 := input[srcOffset+1]
				srcOffset += 2

				disp := int(byte1&0xF)<<8 | int(byte2)
				disp += 1
				copyLen := int(byte1>>4) + 3

				// Copy the bytes from earlier in the decompressed buffer.
				for j := 0; j < copyLen; j++ {
					pos := decompressed.Len() - disp
					if pos < 0 {
						return nil, fmt.Errorf("invalid reference position")
					}
					decompressed.WriteByte(decompressed.Bytes()[pos])
				}
			}
			flag <<= 1
		}
	}

	return decompressed.Bytes(), nil
}

func Save4bppImageBytes(bytes []byte, savename string, palette []color.Color, width, height int) {
	img := image.NewPaletted(image.Rect(0, 0, width, height), palette)

	decompressed, err := DecompressLZ77(bytes)
	if err != nil {
		panic(err)
	}

	tilesX := width / 8  // Number of tiles horizontally
	tilesY := height / 8 // Number of tiles vertically

	// Iterate over each 8x8 tile in the image.
	for tileY := 0; tileY < tilesY; tileY++ {
		for tileX := 0; tileX < tilesX; tileX++ {
			// Calculate the offset to the current tile in the decompressed byte array
			tileOffset := (tileX + tileY*tilesX) * 32 // 32 bytes per tile (8x8 pixels, 4bpp)

			// Iterate over each row in the tile (8 rows per tile)
			for row := 0; row < 8; row++ {
				// Each row contains 4 bytes (2 pixels per byte, 8 pixels total)
				for col := 0; col < 8; col += 2 {
					// Compute the index for the current byte in the decompressed data
					byteIndex := tileOffset + (row * 4) + (col / 2)
					byteVal := decompressed[byteIndex]

					// Extract the high and low nibbles (two pixels from the byte)
					highNibble := (byteVal >> 4) & 0xF // First pixel (left)
					lowNibble := byteVal & 0xF         // Second pixel (right)

					// Calculate the position of the pixel in the final image
					pixelX := tileX*8 + col // X coordinate in the full image
					pixelY := tileY*8 + row // Y coordinate in the full image

					// Set the high nibble pixel (left)
					img.SetColorIndex(pixelX, pixelY, uint8(highNibble))

					// Set the low nibble pixel (right)
					img.SetColorIndex(pixelX+1, pixelY, uint8(lowNibble))
				}
			}
		}
	}

	// Save the image as a PNG file
	file, err := os.Create(savename + ".png")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	err = png.Encode(file, img)
	if err != nil {
		panic(err)
	}
}

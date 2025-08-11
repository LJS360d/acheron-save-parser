package gba

import (
	"acheron-save-parser/utils"
	"encoding/binary"
)

// MAX_TRAINER_ITEMS is the maximum number of items a trainer can have.
const MAX_TRAINER_ITEMS = 4

// TRAINER_NAME_LENGTH is the length of a trainer's name.
const TRAINER_NAME_LENGTH = 10

// Trainer represents an NPC trainer.
type Trainer struct {
	AiFlags              uint32       `json:"aiFlags"`
	Party                []TrainerMon `json:"party"` // built as just *TrainerMon
	Items                [4]uint16    `json:"items"`
	TrainerClass         uint8        `json:"trainerClass"`
	EncounterMusicGender uint8        `json:"encounterMusicGender"` // last bit is gender
	TrainerPic           uint8        `json:"trainerPic"`
	TrainerName          string       `json:"trainerName"`  // built as [11]uint8
	DoubleBattle         bool         `json:"doubleBattle"` // 1 bit

	// Padding              bool   // 1 bit

	StartingStatus uint8 `json:"startingStatus"` // 6 bits
	MugshotColor   uint8 `json:"mugshotColor"`
	PartySize      uint8 `json:"partySize"`
	PoolSize       uint8 `json:"poolSize"`
	PoolRuleIndex  uint8 `json:"poolRuleIndex"`
	PoolPickIndex  uint8 `json:"poolPickIndex"`
	PoolPruneIndex uint8 `json:"poolPruneIndex"`
}

// TrainerMon represents a single Pokémon in a trainer's party.
type TrainerMon struct {
	Nickname         string    `json:"nickname"` // originally *uint8, pointer to a string
	Ev               [6]uint32 `json:"ev"`       // originally *uint8, pointer to an array: [6]uint8, [hp, atk, def, speed, spatk, spdef]
	Iv               [6]uint32 `json:"iv"`       // originally uint32, packed ivs, 6 bits each, [hp, atk, def, speed, spatk, spdef]
	Moves            []uint16  `json:"moves"`    // [4]uint16
	Species          uint16    `json:"species"`
	HeldItem         uint16    `json:"heldItem"`
	Ability          uint16    `json:"ability"`
	Lvl              uint8     `json:"lvl"`
	Ball             uint8     `json:"ball"`
	Friendship       uint8     `json:"friendship"`
	Nature           uint8     `json:"nature"`           // 5 bits
	Gender           uint8     `json:"gender"`           // 2 bits
	IsShiny          bool      `json:"isShiny"`          // 1 bit
	TeraType         uint8     `json:"teraType"`         // 5 bits
	GigantamaxFactor bool      `json:"gigantamaxFactor"` // 1 bit
	ShouldUseDynamax bool      `json:"shouldUseDynamax"` // 1 bit

	// Padding1         uint8 // 1 bit

	DynamaxLevel uint8 `json:"dynamaxLevel"` // 4 bits

	// Padding2         uint8 // 4 bits

	Tags uint32 `json:"tags"`
}

// TrainerSprite represents a sprite for a trainer.
type TrainerSprite struct {
	YOffset         uint8
	FrontPic        CompressedSpriteSheet
	Palette         CompressedSpritePalette
	AnimCmdPtr      uint32 // [][]AnimCmd
	MugshotCoords   Coords16
	MugshotRotation int16
}

// CompressedSpriteSheet holds compressed pixel data for a sprite.
type CompressedSpriteSheet struct {
	DataPtr uint32
	Size    uint16
	Tag     uint16
}

// CompressedSpritePalette holds compressed palette data.
type CompressedSpritePalette struct {
	DataPtr uint32
	Tag     uint16
}

// Coords16 represents a 16-bit coordinate pair.
type Coords16 struct {
	X int16
	Y int16
}

// AnimCmd is a union-like struct for animation commands.
type AnimCmd struct {
	Type  int16
	Frame *AnimFrameCmd
	Loop  *AnimLoopCmd
	Jump  *AnimJumpCmd
}

// AnimFrameCmd represents a single frame in an animation.
type AnimFrameCmd struct {
	ImageValue int16
	FrameDelay int16
	XOffset    int16
	YOffset    int16
	HFlip      int16
	VFlip      int16
}

// AnimLoopCmd represents a loop command in an animation script.
type AnimLoopCmd struct {
	Type     int16
	JumpDest int16
}

// AnimJumpCmd represents a jump command in an animation script.
type AnimJumpCmd struct {
	Type     int16
	JumpDest int16
}

func ParseTrainersBytes(offset int, count int) []*Trainer {
	trainers := make([]*Trainer, 0)
	difficultyCount := count / 3
	// only parses DIFFICULTY_MEDIUM
	for i := difficultyCount; i < count-difficultyCount; i++ {
		t := &Trainer{}
		t.loadFromDataSection(Data[offset+i*Config.TrainerStructSize : offset+i*Config.TrainerStructSize+Config.TrainerStructSize])
		trainers = append(trainers, t)
	}
	return trainers
}

func (t *Trainer) loadFromDataSection(section []byte) {
	// 0x00: u32 aiFlags
	t.AiFlags = binary.LittleEndian.Uint32(section[0:4])

	// 0x04: const struct TrainerMon *party
	partyOffset := binary.LittleEndian.Uint32(section[4:8]) - POINTER_OFFSET

	// 0x08: u16 items[MAX_TRAINER_ITEMS]
	for i := 0; i < MAX_TRAINER_ITEMS; i++ {
		t.Items[i] = binary.LittleEndian.Uint16(section[8+2*i : 10+2*i])
	}

	// 0x10: u8 trainerClass
	t.TrainerClass = section[16]
	// 0x11: u8 encounterMusic_gender
	t.EncounterMusicGender = section[17]
	// 0x12: u8 trainerPic
	t.TrainerPic = section[18]
	// 0x13: u8 trainerName[TRAINER_NAME_LENGTH + 1]
	t.TrainerName = utils.DecodeGFString(section[19 : 19+TRAINER_NAME_LENGTH+1])

	// 0x1E: Bit fields (doubleBattle, padding, startingStatus)
	bitFieldByte := section[30]
	t.DoubleBattle = (bitFieldByte & 0x01) != 0
	// t.Padding = (bitFieldByte >> 1 & 0x01) != 0
	t.StartingStatus = (bitFieldByte >> 2) & 0x3F // 6 bits

	// 0x1F: u8 mugshotColor
	t.MugshotColor = section[31]

	// 0x20-0x24: Remaining u8 fields
	t.PartySize = section[32]
	t.PoolSize = section[33]
	t.PoolRuleIndex = section[34]
	t.PoolPickIndex = section[35]
	t.PoolPruneIndex = section[36]

	// load party
	if partyOffset != 0 {
		for i := 0; i < int(t.PartySize); i++ {
			mon := TrainerMon{}
			mon.loadFromDataSection(Data[partyOffset+uint32(i*Config.TrainerMonStructSize) : partyOffset+uint32(i*Config.TrainerMonStructSize+Config.TrainerMonStructSize)])
			t.Party = append(t.Party, mon)
		}
	}
}

func (t *TrainerMon) loadFromDataSection(section []byte) {
	t.Nickname = utils.DecodePointerString(section, binary.LittleEndian.Uint32(section[0:4]))
	evsPtr := binary.LittleEndian.Uint32(section[4:8]) - POINTER_OFFSET
	if evsPtr != 0 && evsPtr != NULL_POINTER {
		for i := 0; i < 6; i++ {
			t.Ev[i] = binary.LittleEndian.Uint32(Data[evsPtr+uint32(i*4) : evsPtr+uint32(i*4+4)])
		}
	}

	t.Iv = UnpackIVs(binary.LittleEndian.Uint32(section[8:12]))
	for i := 0; i < 4; i++ {
		move := binary.LittleEndian.Uint16(section[12+2*i : 14+2*i])
		t.Moves = append(t.Moves, move)
	}
	t.Species = binary.LittleEndian.Uint16(section[20:22])
	t.HeldItem = binary.LittleEndian.Uint16(section[22:24])
	t.Ability = binary.LittleEndian.Uint16(section[24:26])

	t.Lvl = section[26]
	t.Ball = section[27]
	t.Friendship = section[28]
	// 5 bits
	t.Nature = section[29] & 0b00011111
	// 2 bits
	t.Gender = section[29] & 0b00110000 >> 4
	// 1 bit
	t.IsShiny = section[29]&0b01000000 != 0
	// 5 bits
	t.TeraType = section[29] & 0b10011110
	// 3 bits
	t.GigantamaxFactor = section[29]&0b11000000>>6 != 0
	// 1 bit
	t.ShouldUseDynamax = section[29]&0b10000000 != 0

	// bitFieldByte := section[26]
	// t.Padding = (section[26] >> 1 & 0x01) != 0

	// 4 bits
	t.DynamaxLevel = section[30] & 0b00001111

	// 4 bits of padding

	t.Tags = binary.LittleEndian.Uint32(section[31:35])
}

func ParseTrainerSpritesBytes(offset int, count int) []*TrainerSprite {
	trainerSprites := make([]*TrainerSprite, 0)
	for i := 0; i < count; i++ {
		trainerSprites = append(trainerSprites, &TrainerSprite{})
		trainerSprites[i].loadFromDataSection(Data[offset+i*Config.TrainerSpriteStructSize : offset+i*Config.TrainerSpriteStructSize+Config.TrainerSpriteStructSize])
	}
	return trainerSprites
}

func (t *TrainerSprite) loadFromDataSection(section []byte /* 25 + 3 + 4 bytes */) {
	t.YOffset = uint8(section[0])
	// 3 bytes of padding
	frontPicDataPtr := binary.LittleEndian.Uint32(section[4:8]) - POINTER_OFFSET
	t.FrontPic = CompressedSpriteSheet{
		DataPtr: frontPicDataPtr,
		Size:    binary.LittleEndian.Uint16(section[8:10]),
		Tag:     binary.LittleEndian.Uint16(section[10:12]),
	}
	paletteDataPtr := binary.LittleEndian.Uint32(section[12:16]) - POINTER_OFFSET
	t.Palette = CompressedSpritePalette{
		DataPtr: paletteDataPtr,
		Tag:     binary.LittleEndian.Uint16(section[16:18]),
	}
	t.AnimCmdPtr = binary.LittleEndian.Uint32(section[18:22]) - POINTER_OFFSET
	t.MugshotCoords = Coords16{
		X: int16(binary.LittleEndian.Uint16(section[22:24])),
		Y: int16(binary.LittleEndian.Uint16(section[24:26])),
	}
	t.MugshotRotation = int16(binary.LittleEndian.Uint16(section[26:28]))
	// 4 more bytes of padding (i have no fucking clue why the compiler aligns with 4 more here)
}

/*  [hp, atk, def, speed, spatk, spdef] */
func UnpackIVs(packedIVs uint32) [6]uint32 {
	const mask uint32 = 0x1F // 5 bits set to 1
	ivs := [6]uint32{}
	ivs[0] = packedIVs & mask
	ivs[1] = (packedIVs >> 5) & mask
	ivs[2] = (packedIVs >> 10) & mask
	ivs[3] = (packedIVs >> 15) & mask
	ivs[4] = (packedIVs >> 20) & mask
	ivs[5] = (packedIVs >> 25) & mask

	return ivs
}

package gba

import (
	"acheron-save-parser/utils"
	"encoding/binary"
)

// MAX_TRAINER_ITEMS is the maximum number of items a trainer can have.
const MAX_TRAINER_ITEMS = 4

// TRAINER_NAME_LENGTH is the length of a trainer's name.
const TRAINER_NAME_LENGTH = 10

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

package sav

import (
	"encoding/binary"
)

type Trainer struct {
	name            uint64
	gender          uint8
	SaveWarpFlags   uint8
	Id              uint32
	PublicId        uint16
	PrivateId       uint16
	PlaytimeHours   uint16
	PlaytimeMinutes uint8
	PlaytimeSeconds uint8
	PlaytimeVBlanks uint8
	Options         []byte // 3 bytes
	// Pokedex here (120 bytes)
	// filler_90 		uint64
	SecurityKey uint32
}

func (t *Trainer) new(section []byte) {
	t.name = binary.LittleEndian.Uint64(section[0:8])
	t.gender = section[8]
	t.SaveWarpFlags = section[9]
	t.Id = binary.LittleEndian.Uint32(section[12:16])
	t.PublicId = binary.LittleEndian.Uint16(section[12:14])
	t.PrivateId = binary.LittleEndian.Uint16(section[14:16])
	t.PlaytimeHours = binary.LittleEndian.Uint16(section[16:18])
	t.PlaytimeMinutes = section[18]
	t.PlaytimeSeconds = section[19]
	t.PlaytimeVBlanks = section[20]
	t.Options = section[20:23]
	// Skip 0x18:0x90 (used by Pokedex)
	t.SecurityKey = binary.LittleEndian.Uint32(section[0xAC : 0xAC+4])
}

func (t *Trainer) Gender() string {
	if t.gender == 0 {
		return "boy"
	}
	return "girl"
}

func (t *Trainer) Name() string {
	nameBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(nameBytes, t.name)
	return ReadString(nameBytes)
}

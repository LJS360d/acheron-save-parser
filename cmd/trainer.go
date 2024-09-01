package main

import (
	"encoding/binary"
	"syscall/js"
)

type Trainer struct {
	name            uint64
	gender          uint8
	saveWarpFlags   uint8
	id              uint32
	publicId        uint16
	privateId       uint16
	playtimeHours   uint16
	playtimeMinutes uint8
	playtimeSeconds uint8
	playTimeVBlanks uint8
	options         []byte // 3 bytes
	// Pokedex here (120 bytes)
	// filler_90 		uint64
	securityKey uint32
}

func (t *Trainer) new(section []byte) {
	t.name = binary.LittleEndian.Uint64(section[0:8])
	t.gender = section[8]
	t.saveWarpFlags = section[9]
	t.id = binary.LittleEndian.Uint32(section[12:16])
	t.publicId = binary.LittleEndian.Uint16(section[12:14])
	t.privateId = binary.LittleEndian.Uint16(section[14:16])
	t.playtimeHours = binary.LittleEndian.Uint16(section[16:18])
	t.playtimeMinutes = section[18]
	t.playtimeSeconds = section[19]
	t.playTimeVBlanks = section[20]
	t.options = section[20:23]
	// skip 0x18:0x90 (used by Pokedex)
	t.securityKey = binary.LittleEndian.Uint32(section[0xAC : 0xAC+4])
}

func (t *Trainer) toJS() js.Value {
	return js.ValueOf(map[string]interface{}{
		"name":   t.Name(),
		"gender": t.Gender(),
		// "saveWarpFlags":   t.saveWarpFlags,
		"id": t.id,
		// "publicId":        t.publicId,
		// "privateId":       t.privateId,
		"playtimeHours":   t.playtimeHours,
		"playtimeMinutes": t.playtimeMinutes,
		"playtimeSeconds": t.playtimeSeconds,
		// "playTimeVBlanks": t.playTimeVBlanks,
		// "options":         t.options,
		// "securityKey":     t.securityKey,
	})
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
	return readString(nameBytes)
}

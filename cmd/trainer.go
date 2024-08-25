package main

import "encoding/binary"

type Trainer struct {
	name   string // 7 bytes
	gender string // 1 byte
	// 1 byte of unused data
	id        uint32
	publicId  uint16
	privateId uint16
	time      int    // 5 bytes
	options   []byte // 3 bytes
	// 4 bytes of GameCode
	// 4 bytes of SecurityKey
}

func (t *Trainer) new(section []byte) {
	t.name = readString(section[0:7])
	if section[8] == 0 {
		t.gender = "boy"
	} else {
		t.gender = "girl"
	}
	t.id = binary.LittleEndian.Uint32(section[12:16])
	t.publicId = binary.LittleEndian.Uint16(section[12:14])
	t.privateId = binary.LittleEndian.Uint16(section[14:16])
	t.time = int(binary.LittleEndian.Uint32(section[16:20]))
	t.options = section[20:23]
}

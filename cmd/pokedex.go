package main

import "encoding/binary"

type Pokedex struct {
	order         uint8
	mode          uint8
	nationalMagic uint8
	// unknown2          uint8
	unownPersonality  uint32
	spindaPersonality uint32
	// unknown3          uint32
	// filler []byte // 104 bytes
}

func (p *Pokedex) new(section []byte) {
	// starts from offset 0x18 of section 0
	section = section[0x18:0x90]
	p.order = section[0]
	p.mode = section[1]
	// must equal 0xDA in order to have National mode
	p.nationalMagic = section[2]
	// p.unknown2 = section[3]
	p.unownPersonality = binary.LittleEndian.Uint32(section[4:8])
	p.spindaPersonality = binary.LittleEndian.Uint32(section[8:12])
	// p.unknown3 = binary.LittleEndian.Uint32(section[12:16])
	// p.filler = section[16:104]
}

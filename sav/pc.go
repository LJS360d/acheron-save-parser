package sav

import (
	"encoding/binary"
)

type PC struct {
	CurrentBox uint32
	Pokemon    [420]Pokemon
	BoxNames   string
}

func (pc *PC) new(section []byte) {
	pc.CurrentBox = binary.LittleEndian.Uint32(section[0x0000:0x0004]) - 1
	for i := 0; i < len(pc.Pokemon); i++ {
		ix := 0x0004 + i*80
		pc.Pokemon[i].newBoxed(section[ix : ix+80])
	}
	pc.BoxNames = ReadString(section[0x8344:0x83C2])
}

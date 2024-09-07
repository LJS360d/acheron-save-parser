package sav

import (
	"acheron-save-parser/utils"
	"encoding/binary"
)

type PC struct {
	CurrentBox uint32
	Pokemon    [420]Pokemon
	BoxNames   [14]string
}

func (pc *PC) new(section []byte) {
	pc.CurrentBox = binary.LittleEndian.Uint32(section[0x0000:0x0004]) - 1
	for i := 0; i < len(pc.Pokemon); i++ {
		ix := 0x0004 + i*80
		pc.Pokemon[i].newBoxed(section[ix : ix+80])
	}
	for i := 0; i < len(pc.BoxNames); i++ {
		nameLength := 9
		ix := 0x8344 + i*nameLength
		pc.BoxNames[i] = utils.DecodeGFString(section[ix : ix+nameLength])
	}
}

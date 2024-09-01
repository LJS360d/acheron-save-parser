package main

import (
	"encoding/binary"
	"syscall/js"
)

type PC struct {
	currentBox uint32
	pokemon    [420]Pokemon
	boxNames   string
}

func (pc *PC) new(section []byte) {
	pc.currentBox = binary.LittleEndian.Uint32(section[0x0000:0x0004]) - 1
	for i := 0; i < len(pc.pokemon); i++ {
		ix := 0x0004 + i*80
		pc.pokemon[i].newBoxed(section[ix : ix+80])
	}
	pc.boxNames = readString(section[0x8344:0x83C2])
}

func (pc *PC) toJS() js.Value {
	jsMons := make([]js.Value, len(pc.pokemon))
	for i, value := range pc.pokemon {
		jsMons[i] = value.ToJS()
	}
	return js.ValueOf(map[string]interface{}{
		"currentBox": pc.currentBox,
		"pokemon":    toJSArray(jsMons),
		"boxNames":   pc.boxNames,
	})
}

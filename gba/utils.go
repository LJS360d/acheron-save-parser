package gba

import (
	"acheron-save-parser/sav"
)

func ParsePointerString(data []byte, offset uint32) string {
	start := int(offset)
	end := start
	for end < len(data) && data[end] != 0xFF /* 0xFF is the \0 of GF text encoding */ {
		end++
	}
	return sav.DecodeGFString(data[start:end])
}

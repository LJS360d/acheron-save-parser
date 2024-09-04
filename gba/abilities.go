package gba

import (
	"acheron-save-parser/sav"
	"encoding/binary"
)

type AbilityData struct {
	Name              string // 20 bytes
	DescriptionPtr    uint32
	Description       string
	aiRating          int8
	cantBeCopied      bool
	cantBeSwapped     bool
	cantBeTraced      bool
	cantBeSuppressed  bool
	cantBeOverwritten bool
	breakable         bool
	failsOnImposter   bool
}

const (
	ABILITY_INFO_SIZE_LENGTH12 = 25
	ABILITY_INFO_SIZE_LENGTH16 = 28
)

var (
	AbilityNameLength12 = false
)

func ParseAbilitiesBytes(data []byte, offset int, count int) []*AbilityData {
	abilities := make([]*AbilityData, count)
	a := &AbilityData{}
	var abilityInfoSize int
	if AbilityNameLength12 {
		abilityInfoSize = ABILITY_INFO_SIZE_LENGTH12
		for i := 0; i < count; i++ {
			a.new_name12(data[offset+i*abilityInfoSize : offset+i*abilityInfoSize+abilityInfoSize])
			abilities[i] = a
			a.Description = ParseAbilityDescription(data, a.DescriptionPtr)
		}
	} else {
		abilityInfoSize = ABILITY_INFO_SIZE_LENGTH16
		for i := 0; i < count; i++ {
			a.new_name16(data[offset+i*abilityInfoSize : offset+i*abilityInfoSize+abilityInfoSize])
			abilities[i] = a
			a.Description = ParseAbilityDescription(data, a.DescriptionPtr)
		}
	}
	return abilities
}

func (a *AbilityData) new_name12(section []byte /* 22 + 3 bytes */) {
	a.Name = sav.DecodeGFString(section[0:16])
	a.DescriptionPtr = binary.LittleEndian.Uint32(section[16:20]) - POINTER_OFFSET
	a.aiRating = int8(section[20])
	a.cantBeCopied = section[21]&0x1 == 1
	a.cantBeSwapped = section[21]&0x2 == 1
	a.cantBeTraced = section[21]&0x4 == 1
	a.cantBeSuppressed = section[21]&0x8 == 1
	a.cantBeOverwritten = section[21]&0x10 == 1
	a.breakable = section[21]&0x20 == 1
	a.failsOnImposter = section[21]&0x40 == 1
	// 3 bytes of padding for the pointer boundary
}

func (a *AbilityData) new_name16(section []byte /* 26 + 2 bytes */) {
	a.Name = sav.DecodeGFString(section[0:20])
	a.DescriptionPtr = binary.LittleEndian.Uint32(section[20:24]) - POINTER_OFFSET
	a.aiRating = int8(section[24])
	a.cantBeCopied = section[25]&0x1 == 0x1
	a.cantBeSwapped = section[25]&0x2 == 0x2
	a.cantBeTraced = section[25]&0x4 == 0x4
	a.cantBeSuppressed = section[25]&0x8 == 0x8
	a.cantBeOverwritten = section[25]&0x10 == 0x10
	a.breakable = section[25]&0x20 == 0x20
	a.failsOnImposter = section[25]&0x40 == 0x40
	// 2 bytes of padding for the pointer boundary
}

func ParseAbilityDescription(data []byte, offset uint32) string {
	start := int(offset)
	end := start
	for end < len(data) && data[end] != 173 /* 173->GF Encoding->"." */ {
		end++
	}
	return sav.DecodeGFString(data[start : end+1])
}

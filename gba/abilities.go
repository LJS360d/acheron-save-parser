package gba

import (
	"encoding/binary"
	"rom-parser/utils"
)

type AbilityData struct {
	Name              string `json:"name"` // 20 bytes
	descriptionPtr    uint32
	Description       string `json:"description"`
	aiRating          int8
	cantBeCopied      bool
	cantBeSwapped     bool
	cantBeTraced      bool
	cantBeSuppressed  bool
	cantBeOverwritten bool
	breakable         bool
	failsOnImposter   bool
}

func ParseAbilitiesBytes(offset int, count int) []*AbilityData {
	abilities := make([]*AbilityData, count)
	abilityInfoSize := Config.AbilityInfoSize
	for i := 0; i < count; i++ {
		a := &AbilityData{}
		if abilityInfoSize == 25 {
			a.loadFromDataSection25(Data[offset+i*abilityInfoSize : offset+i*abilityInfoSize+abilityInfoSize])
		} else {
			a.loadFromDataSection(Data[offset+i*abilityInfoSize : offset+i*abilityInfoSize+abilityInfoSize])
		}
		abilities[i] = a
		a.Description = utils.DecodePointerString(Data, a.descriptionPtr)
	}
	return abilities
}

func (a *AbilityData) loadFromDataSection25(section []byte /* 22 + 3 bytes */) {
	a.Name = utils.DecodeGFString(section[0:16])
	a.descriptionPtr = binary.LittleEndian.Uint32(section[16:20]) - POINTER_OFFSET
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

func (a *AbilityData) loadFromDataSection(section []byte /* 26 + 2 bytes */) {
	a.Name = utils.DecodeGFString(section[0:20])
	a.descriptionPtr = binary.LittleEndian.Uint32(section[20:24]) - POINTER_OFFSET
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

package gba

import (
	"acheron-save-parser/utils"
	"encoding/binary"
)

type NatureData struct {
	namePtr                 uint32
	Name                    string
	StatUp                  uint8
	StatDown                uint8
	BackAnim                uint8
	PokeBlockAnim           [2]uint8
	BattlePalacePercents    [4]uint8
	BattlePalaceFlavorText  uint8
	BattlePalaceSmokescreen uint8
	natureGirlMessagePtr    uint32
	NatureGirlMessage       string
}

func ParseNaturesInfoBytes(offset int, count int) []*NatureData {
	natures := make([]*NatureData, count)
	for i := 0; i < count; i++ {
		n := &NatureData{}
		n.loadFromDataSection(Data[offset+i*Config.NatureInfoSize : offset+i*Config.NatureInfoSize+Config.NatureInfoSize])
		natures[i] = n
		n.Name = utils.DecodePointerString(Data, n.namePtr)
		n.NatureGirlMessage = utils.DecodePointerString(Data, n.natureGirlMessagePtr)
	}
	return natures
}

func (n *NatureData) loadFromDataSection(section []byte /* 20 bytes */) {
	n.namePtr = binary.LittleEndian.Uint32(section[0:4]) - POINTER_OFFSET
	//
	n.StatUp = section[4]
	n.StatDown = section[5]
	n.BackAnim = section[6]
	n.PokeBlockAnim[0] = section[7]
	n.PokeBlockAnim[1] = section[8]
	n.BattlePalacePercents[0] = section[9]
	n.BattlePalacePercents[1] = section[10]
	n.BattlePalacePercents[2] = section[11]
	n.BattlePalacePercents[3] = section[12]
	n.BattlePalaceFlavorText = section[13]
	n.BattlePalaceSmokescreen = section[14]
	// 1 byte of pointer boundary padding
	n.natureGirlMessagePtr = binary.LittleEndian.Uint32(section[16:20]) - POINTER_OFFSET
	//
}

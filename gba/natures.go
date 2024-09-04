package gba

import (
	"encoding/binary"
)

type NatureData struct {
	NamePtr                 uint32
	Name                    string
	StatUp                  uint8
	StatDown                uint8
	BackAnim                uint8
	PokeBlockAnim           [2]uint8
	BattlePalacePercents    [4]uint8
	BattlePalaceFlavorText  uint8
	BattlePalaceSmokescreen uint8
	NatureGirlMessagePtr    uint32
	NatureGirlMessage       string
}

const (
	NATURE_INFO_SIZE = 32
)

func ParseNaturesInfoBytes(data []byte, offset int, count int) []*NatureData {
	natures := make([]*NatureData, count)
	for i := 0; i < count; i++ {
		n := &NatureData{}
		n.new(data[offset+i*NATURE_INFO_SIZE : offset+i*NATURE_INFO_SIZE+NATURE_INFO_SIZE])
		natures[i] = n
	}
	return natures
}

func (n *NatureData) new(section []byte /* 32 bytes */) {
	n.NamePtr = binary.LittleEndian.Uint32(section[0:4]) - POINTER_OFFSET
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
	n.NatureGirlMessagePtr = binary.LittleEndian.Uint32(section[15:19]) - POINTER_OFFSET
	//
}

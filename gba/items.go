package gba

import (
	"acheron-save-parser/utils"
	"encoding/binary"
)

type ItemData struct {
	Price           uint32
	SecondaryId     uint16
	FieldUseFuncPtr uint32 // ItemUseFunc
	DescriptionPtr  uint32
	Description     string
	EffectPtr       uint32
	Name            string // 20 bytes
	PluralName      string // 22 bytes
	HoldEffect      uint8
	HoldEffectParam uint8
	Importance      uint8
	Pocket          uint8
	Type            uint8
	BattleUsage     uint8
	FlingPower      uint8
	IconPicPtr      uint32
	IconPalettePtr  uint32
}

const (
	ITEM_INFO_SIZE = 80
)

func ParseItemsInfoBytes(data []byte, offset int, count int) []*ItemData {
	items := make([]*ItemData, count)
	for i := 0; i < count; i++ {
		n := &ItemData{}
		n.new(data[offset+i*ITEM_INFO_SIZE : offset+i*ITEM_INFO_SIZE+ITEM_INFO_SIZE])
		items[i] = n
		n.Description = utils.DecodePointerString(data, n.DescriptionPtr)
	}
	return items
}

func (i *ItemData) new(section []byte /* 80 bytes */) {
	i.Price = binary.LittleEndian.Uint32(section[0:4])
	i.SecondaryId = binary.LittleEndian.Uint16(section[4:6])
	// 2 bytes of padding for the pointer boundary
	i.FieldUseFuncPtr = binary.LittleEndian.Uint32(section[8:12]) - POINTER_OFFSET
	i.DescriptionPtr = binary.LittleEndian.Uint32(section[12:16]) - POINTER_OFFSET
	i.EffectPtr = binary.LittleEndian.Uint32(section[16:20]) - POINTER_OFFSET
	i.Name = utils.DecodeGFString(section[20:40])
	i.PluralName = utils.DecodeGFString(section[40:62])
	i.HoldEffect = section[62]
	i.HoldEffectParam = section[63]
	i.Importance = section[64]
	i.Pocket = section[65]
	i.Type = section[66]
	i.BattleUsage = section[67]
	i.FlingPower = section[68]
	i.IconPicPtr = binary.LittleEndian.Uint32(section[69:73]) - POINTER_OFFSET
	i.IconPalettePtr = binary.LittleEndian.Uint32(section[73:77]) - POINTER_OFFSET
}

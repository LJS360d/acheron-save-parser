package gba

import (
	"acheron-save-parser/utils"
	"encoding/binary"
)

type ItemData struct {
	// index in the Items array, added for convenience
	Id              uint   `json:"id"`
	Price           uint32 `json:"price"`
	SecondaryId     uint16 `json:"secondaryId"`
	fieldUseFuncPtr uint32 `json:"-"`
	descriptionPtr  uint32 `json:"-"`
	Description     string `json:"description"`
	effectPtr       uint32 `json:"-"`
	Name            string `json:"name"`       // 20 bytes
	PluralName      string `json:"pluralName"` // 22 bytes
	HoldEffect      uint8  `json:"holdEffect"`
	HoldEffectParam uint8  `json:"holdEffectParam"`
	Importance      uint8  `json:"importance"`
	Pocket          uint8  `json:"pocket"`
	Type            uint8  `json:"type"`
	BattleUsage     uint8  `json:"battleUsage"`
	FlingPower      uint8  `json:"flingPower"`
	IconPicPtr      uint32 `json:"-"`
	IconPalettePtr  uint32 `json:"-"`
}

func ParseItemsInfoBytes(offset int, count int) []*ItemData {
	items := make([]*ItemData, count)
	for i := 0; i < count; i++ {
		n := &ItemData{}
		n.loadFromDataSection(Data[offset+i*Config.ItemInfoSize : offset+i*Config.ItemInfoSize+Config.ItemInfoSize])
		n.Id = uint(i)
		items[i] = n
		n.Description = utils.DecodePointerString(Data, n.descriptionPtr)
	}
	return items
}

func (i *ItemData) loadFromDataSection(section []byte /* 80 bytes */) {
	i.Price = binary.LittleEndian.Uint32(section[0:4])
	i.SecondaryId = binary.LittleEndian.Uint16(section[4:6])
	// 2 bytes of padding for the pointer boundary
	i.fieldUseFuncPtr = binary.LittleEndian.Uint32(section[8:12]) - POINTER_OFFSET
	i.descriptionPtr = binary.LittleEndian.Uint32(section[12:16]) - POINTER_OFFSET
	i.effectPtr = binary.LittleEndian.Uint32(section[16:20]) - POINTER_OFFSET
	i.Name = utils.DecodeGFString(section[20:40])
	i.PluralName = utils.DecodeGFString(section[40:62])
	i.HoldEffect = section[62]
	i.HoldEffectParam = section[63]
	i.Importance = section[64]
	i.Pocket = section[65]
	i.Type = section[66]
	i.BattleUsage = section[67]
	i.FlingPower = section[68]
	// more padding
	i.IconPicPtr = binary.LittleEndian.Uint32(section[72:76]) - POINTER_OFFSET
	i.IconPalettePtr = binary.LittleEndian.Uint32(section[76:80]) - POINTER_OFFSET
}

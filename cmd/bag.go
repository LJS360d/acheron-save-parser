package main

import (
	"acheron-save-parser/data"
	"encoding/binary"
)

type Bag struct {
	money          uint32
	coins          uint16
	registeredItem uint16
	pcItems        [50]ItemSlot
	items          [30]ItemSlot
	keyItems       [30]ItemSlot
	pokeBalls      [16]ItemSlot
	tmhm           [64]ItemSlot
	berries        [46]ItemSlot
	pokeblocks     [40]Pokeblock
	// filler1 []byte // 52 bytes
}

func (b *Bag) new(section []byte, securityKey uint32) {
	// last 2 bytes of security key
	var keyLower uint16 = uint16(securityKey & 0xFFFF)
	// bag segment starts from offset 0x490 of section 1
	bagSection := section[0x490:]
	b.money = binary.LittleEndian.Uint32(bagSection[0:4]) ^ securityKey
	b.coins = binary.LittleEndian.Uint16(bagSection[4:6]) ^ keyLower
	b.registeredItem = binary.LittleEndian.Uint16(bagSection[6:8])
	offset := 8
	for i := 0; i < len(b.pcItems); i++ {
		b.pcItems[i].new(bagSection[offset+i*4 : offset+i*4+4])
	}
	offset += len(b.pcItems) * 4
	for i := 0; i < len(b.items); i++ {
		ix := offset + i*4
		b.items[i].newEncrypted(bagSection[ix:offset+i*4+4], keyLower)
	}
	offset += len(b.items) * 4
	for i := 0; i < len(b.keyItems); i++ {
		b.keyItems[i].newEncrypted(bagSection[offset+i*4:offset+i*4+4], keyLower)
	}
	offset += len(b.keyItems) * 4
	for i := 0; i < len(b.pokeBalls); i++ {
		b.pokeBalls[i].newEncrypted(bagSection[offset+i*4:offset+i*4+4], keyLower)
	}
	offset += len(b.pokeBalls) * 4
	for i := 0; i < len(b.tmhm); i++ {
		b.tmhm[i].newEncrypted(bagSection[offset+i*4:offset+i*4+4], keyLower)
	}
	offset += len(b.tmhm) * 4
	for i := 0; i < len(b.berries); i++ {
		b.berries[i].newEncrypted(bagSection[offset+i*4:offset+i*4+4], keyLower)
	}
	offset += len(b.berries) * 4
	for i := 0; i < len(b.pokeblocks); i++ {
		b.pokeblocks[i].new(bagSection[offset+i*7 : offset+i*7+7])
	}
	// offset += len(b.pokeblocks) * 7
	// b.filler1 = bagSection[offset : offset+52]
}

type ItemSlot struct {
	itemId   uint16
	quantity uint16
}

func (i *ItemSlot) new(section []byte) {
	i.itemId = binary.LittleEndian.Uint16(section[0:2])
	i.quantity = binary.LittleEndian.Uint16(section[2:4])
}

func (i *ItemSlot) newEncrypted(section []byte, lowerKey uint16) {
	i.itemId = binary.LittleEndian.Uint16(section[0:2])
	i.quantity = binary.LittleEndian.Uint16(section[2:4]) ^ lowerKey
}

func (i *ItemSlot) Name() string {
	return data.ItemName[i.itemId]
}

type Pokeblock struct {
	color  uint8
	spicy  uint8
	dry    uint8
	sweet  uint8
	bitter uint8
	sour   uint8
	feel   uint8
}

func (p *Pokeblock) new(section []byte) {
	p.color = section[0]
	p.spicy = section[1]
	p.dry = section[2]
	p.sweet = section[3]
	p.bitter = section[4]
	p.sour = section[5]
	p.feel = section[6]
}

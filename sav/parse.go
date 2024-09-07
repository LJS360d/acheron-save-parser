package sav

import (
	"bytes"
	"encoding/binary"
)

const (
	SECTOR_DATA_SIZE = 3968
	SAVE3_CHUNK_SIZE = 116
	FOOTER_SIZE      = 12
	SECTOR_SIZE      = SECTOR_DATA_SIZE + SAVE3_CHUNK_SIZE + FOOTER_SIZE // 4096
	SAVE_SLOT_SIZE   = SECTOR_SIZE * 14                                  // 57344
)

type SavData struct {
	Trainer
	Pokedex
	Team
	Bag
	PC
}

func ParseSavBytes(data []byte /* 57'344 Bytes */) *SavData {
	slot1 := data[0:SAVE_SLOT_SIZE] // 14 sectors
	// Save slot 2
	slot2 := data[SAVE_SLOT_SIZE : SAVE_SLOT_SIZE*2] // 14 sectors
	// Hall of fame
	// hof := data[SAVE_SLOT_SIZE*2 : SAVE_SLOT_SIZE*2+SECTOR_SIZE*2] // 2 sectors
	// Trainer hill
	// th := data[SAVE_SLOT_SIZE*2+SECTOR_SIZE*2 : SAVE_SLOT_SIZE*2+SECTOR_SIZE*3] // 1 sector
	// Recorded battle
	// rb := data[SAVE_SLOT_SIZE*2+SECTOR_SIZE*3 : SAVE_SLOT_SIZE*2+SECTOR_SIZE*4] // 1 sector
	activeSlot := getActiveSaveSlot(slot1, slot2)
	save := &SavData{}
	save.ParseSaveSlot(activeSlot)
	return save
}

func (s *SavData) ParseSaveSlot(saveSlot []byte) {
	sections := make([][]byte, 14)

	for i := 0; i < 14; i++ {
		ix := i * SECTOR_SIZE
		section := saveSlot[ix : ix+SECTOR_SIZE]
		footer := section[4084:]
		id := binary.LittleEndian.Uint16(footer[0:2])
		sections[id] = section[0:3968]
	}

	trainer := Trainer{}
	trainer.new(sections[0])
	s.Trainer = trainer

	pokedex := Pokedex{}
	pokedex.new(sections[0])
	s.Pokedex = pokedex

	// TODO map data - 0x00 to 0x234 of section 1

	team := Team{}
	team.new(sections[1])
	s.Team = team

	bag := Bag{}
	bag.new(sections[1], trainer.SecurityKey)
	s.Bag = bag

	// TODO game state flags - 0x00 to 0x??? of section 2
	// TODO game specific data - 0x00 to 0x??? of section 3 and 4

	pc := PC{}
	pc.new(bytes.Join(sections[5:14], nil))
	s.PC = pc
}

func getSaveIndex(data []byte) int {
	// read the last 4 bytes of the save file
	saveIndexRaw := data[4084:]
	// parse it as a number
	id := binary.LittleEndian.Uint16(saveIndexRaw[0:2])
	return int(id)
}

func getActiveSaveSlot(slot1 []byte, slot2 []byte) []byte {
	if getSaveIndex(slot1) < getSaveIndex(slot2) {
		return slot1
	}
	return slot2
}

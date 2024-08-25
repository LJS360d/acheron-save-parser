package main

import (
	"encoding/binary"
	"os"
)

type SaveData struct {
	Trainer
	Team
}

func NewSaveData(filename string) (*SaveData, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	} // 57344
	// Save slot 1
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
	save := &SaveData{}
	save.processSaveSlot(activeSlot)
	return save, nil
}

func (s *SaveData) processSaveSlot(saveSlot []byte) {
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

	team := Team{}
	team.new(sections[1])
	s.Team = team
}

// TODO bag off section 1
type Bag struct {
}

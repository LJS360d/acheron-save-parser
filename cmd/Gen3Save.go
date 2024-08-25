package main

import (
	"encoding/binary"
	"fmt"
	"os"
)

type Gen3Save struct {
	Trainer
	Team
}

func NewGen3Save(filename string) (*Gen3Save, error) {
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
	save := &Gen3Save{}
	save.processSaveSlot(activeSlot)
	return save, nil
}

func (s *Gen3Save) processSaveSlot(saveSlot []byte) {
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

	fmt.Println(s.Team.size)
}

type Trainer struct {
	name   string // 7 bytes
	gender string // 1 byte
	// 1 byte of unused data
	id        uint32
	publicId  uint16
	privateId uint16
	time      int    // 5 bytes
	options   []byte // 3 bytes
	// 4 bytes of GameCode
	// 4 bytes of SecurityKey
}

func (t *Trainer) new(section []byte) {
	t.name = readString(section[0:7])
	if section[8] == 0 {
		t.gender = "boy"
	} else {
		t.gender = "girl"
	}
	t.id = binary.LittleEndian.Uint32(section[12:16])
	t.publicId = binary.LittleEndian.Uint16(section[12:14])
	t.privateId = binary.LittleEndian.Uint16(section[14:16])
	t.time = int(binary.LittleEndian.Uint32(section[16:20]))
	t.options = section[20:23]
}

// Both off section 1
type Team struct {
	size    int
	pokemon []Pokemon // 600 bytes
}

func (t *Team) new(section []byte) {
	// go at offset 0x234 for size
	t.size = int(binary.LittleEndian.Uint32(section[0x234:0x238]))
	// between 0x238 and 0x490 is the pokemon data
	t.pokemon = make([]Pokemon, t.size)
	for i := 0; i < int(t.size); i++ {
		t.pokemon[i].new(section[0x238+(i*100) : 0x238+((i+1)*100)])
	}
}

// TODO bag off section 1
type Bag struct {
}

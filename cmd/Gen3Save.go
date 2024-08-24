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

// 100 bytes
type Pokemon struct {
	personalityValue     uint32
	otId                 uint32
	nickname             string // 10 bytes
	language             uint8  // 3 bits
	hiddenNatureModifier uint8  // 5 bits
	isBadEgg             bool
	hasSpecies           bool
	isEgg                bool
	// blockBoxRS bool
	daysSinceFormChange uint8 // 3 bits - 7 days.
	// unused_13 bool
	otName           string // 7 bytes
	markings         uint8  // 4 bits
	compressedStatus uint8  // 4 bits
	checksum         uint16
	hpLost           uint16 // 14 bits; // 16383 HP.
	shinyModifier    bool
	// unused_1E        bool

	// 2 bytes of unused data
	PokemonData     // 48 bytes - encrypted
	statusCondition uint32
	level           uint8
	mailId          uint8
	currentHp       uint16
	totalHp         uint16
	attack          uint16
	defense         uint16
	speed           uint16
	specialAttack   uint16
	specialDefense  uint16
}

func (p *Pokemon) new(section []byte) {
	p.personalityValue = binary.LittleEndian.Uint32(section[0:4])
	p.otId = binary.LittleEndian.Uint32(section[4:8])
	p.nickname = readString(section[8:18])
	// first 3 bits of [18]
	p.language = section[18] & 0x7 // should always be 2
	// last 5 bits of [18]
	p.hiddenNatureModifier = section[18] & 0x1F
	// first bit of [19]
	p.isBadEgg = section[19]&0x1 == 1
	// second bit of [19]
	p.hasSpecies = section[19]&0x2 == 1
	// third bit of [19]
	p.isEgg = section[19]&0x4 == 1
	// fourth bit of [19]
	// p.blockBoxRS = section[19] & 0x38
	// fifth, sixth and seventh bit of [19]
	p.daysSinceFormChange = section[19] & 0xC0
	// last bit of [19]
	// p.unused_13 = section[19]&0x80 == 1
	p.otName = readString(section[20:28])
	// first 4 bits of [28]
	p.markings = section[28] & 0xF
	// last 4 bits of [28]
	p.compressedStatus = section[28] & 0xF0
	p.checksum = binary.LittleEndian.Uint16(section[30:32])
	// first 14 bits of [32:34]
	p.hpLost = binary.LittleEndian.Uint16(section[32:34]) & 0x3FFF
	// second to last 1 bit of [32:34]
	p.shinyModifier = section[34]&0x80 == 1
	// last bit of [32:34]
	// p.unused_1E = section[34]&0x40 == 1
	p.PokemonData.new(section[32:80], p.otId, p.personalityValue)
	p.statusCondition = binary.LittleEndian.Uint32(section[80:84])
	p.level = section[84]
	p.mailId = section[85]
	p.currentHp = binary.LittleEndian.Uint16(section[86:88])
	p.totalHp = binary.LittleEndian.Uint16(section[88:90])
	p.attack = binary.LittleEndian.Uint16(section[90:92])
	p.defense = binary.LittleEndian.Uint16(section[92:94])
	p.speed = binary.LittleEndian.Uint16(section[94:96])
	p.specialAttack = binary.LittleEndian.Uint16(section[96:98])
	p.specialDefense = binary.LittleEndian.Uint16(section[98:100])
}

// 48 bytes (by default)
// 12 bytes * 4 sections
type PokemonData struct {
	// Growth section
	species    uint16 // 11 bits
	teraType   uint16 // 5 bits
	item       uint16 // 10 bits
	experience uint32
	ppBonuses  uint8 // TODO interpret
	friendship uint8
	// 2 bytes of unused data

	// Moves section
	move1 uint16
	move2 uint16
	move3 uint16
	move4 uint16
	pp1   uint8
	pp2   uint8
	pp3   uint8
	pp4   uint8

	// EVs & Conditions section
	hpEv      uint8
	atkEv     uint8
	defEv     uint8
	speEv     uint8
	spaEv     uint8
	spdEv     uint8
	coolness  uint8
	beauty    uint8
	cuteness  uint8
	smartness uint8
	toughness uint8
	feel      uint8

	// Misc section
	pokerus     uint8 // TODO interpret
	metLocation uint8
	originsInfo uint16 // TODO interpret
	// IVs, isEgg, whichAbility (modified from vanilla)
	// 4 bit each for IVs
	hpIv    uint8
	atkIv   uint8
	defIv   uint8
	speIv   uint8
	spaIv   uint8
	spdIv   uint8
	isEgg   bool
	ability uint8 // TODO definitely modified from vanilla
	// Ribbons & Obedience
	// TODO ribbons
}

func (p *PokemonData) new(data []byte, otId uint32, personalityValue uint32) {
	orders := []string{"GAEM", "GAME", "GEAM", "GEMA", "GMAE", "GMEA", "AGEM", "AGME", "AEGM", "AEMG", "AMGE", "AMEG", "EGAM", "EGMA", "EAGM", "EAMG", "EMGA", "EMAG", "MGAE", "MGEA", "MAGE", "MAEG", "MEGA", "MEAG"}
	key := otId ^ personalityValue
	orderIndex := personalityValue % 24
	order := orders[orderIndex]
	// 12 bytes per section
	sections := make(map[byte][]byte)
	for i := 0; i < 4; i++ {
		section := order[i]
		sectionData := data[i*12 : (i+1)*12]
		decrypted := decryptSubsection(sectionData, key)
		sections[section] = decrypted
	}
	// TODO set decrypted data to the object
	fmt.Println(sections)
}

func decryptSubsection(data []byte, key uint32) []byte {
	if len(data) != 12 {
		return []byte{}
	}
	a := binary.LittleEndian.Uint32(data[0:4]) ^ key
	b := binary.LittleEndian.Uint32(data[4:8]) ^ key
	c := binary.LittleEndian.Uint32(data[8:12]) ^ key
	result := make([]byte, 12)
	binary.LittleEndian.PutUint32(result[0:4], a)
	binary.LittleEndian.PutUint32(result[4:8], b)
	binary.LittleEndian.PutUint32(result[8:12], c)
	return result
}

// TODO bag off section 1
type Bag struct {
}

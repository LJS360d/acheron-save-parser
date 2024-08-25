package main

import (
	"acheron-save-parser/data"
	"encoding/binary"
	"fmt"
	"strings"
)

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
	compressedStatus uint8  // 4 bits // TODO ignored
	checksum         uint16
	hpLost           uint16 // 14 bits; // 16383 HP.
	shinyModifier    bool
	// unused_1E        bool
	PokemonData            // 48 bytes - encrypted
	statusCondition uint32 // TODO ignored
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
	p.hiddenNatureModifier = (section[18] >> 3) & 0x1F
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

func (p *Pokemon) SpeciesName() string {
	return data.SpeciesName[p.species]
}
func (p *Pokemon) ItemName() string {
	return data.ItemName[p.item]
}
func (p *Pokemon) AbilityName() string {
	return data.AbilityName[data.SpeciesAbility[p.species][p.abilityNum]]
}
func (p *Pokemon) NatureName() string {
	return data.NatureName[uint8(p.personalityValue%25)]
}
func (p *Pokemon) Nickname() string {
	return p.nickname + p.nickname11th + p.nickname12th
}
func (p *Pokemon) Moves() []string {
	var moves []string
	moves = append(moves, data.MoveName[p.move1])
	moves = append(moves, data.MoveName[p.move2])
	moves = append(moves, data.MoveName[p.move3])
	moves = append(moves, data.MoveName[p.move4])
	return moves
}
func (p *Pokemon) toSDExportFormat() string {
	var sb strings.Builder
	// Print Pokémon name and item
	if p.nickname != "" {
		if p.ItemName() != "None" {
			sb.WriteString(fmt.Sprintf("\t%s (%s) @ %s\n", p.Nickname(), p.SpeciesName(), p.ItemName()))
		} else {
			sb.WriteString(fmt.Sprintf("\t%s (%s)\n", p.Nickname(), p.SpeciesName()))
		}
	} else {
		if p.ItemName() != "None" {
			sb.WriteString(fmt.Sprintf("\t%s @ %s\n", p.SpeciesName(), p.ItemName()))
		} else {
			sb.WriteString(fmt.Sprintf("\t%s\n", p.SpeciesName()))
		}
	}
	// Print Ability
	sb.WriteString(fmt.Sprintf("\tAbility: %s\n", p.AbilityName()))
	// Print Level
	sb.WriteString(fmt.Sprintf("\tLevel: %d\n", p.level))
	// Print EVs
	sb.WriteString(fmt.Sprintf("\tEVs: %d HP / %d Atk / %d Def / %d SpA / %d SpD / %d Spe\n",
		p.hpEv, p.atkEv, p.defEv, p.speEv, p.spaEv, p.spdEv))
	// Print Nature
	sb.WriteString(fmt.Sprintf("\t%s Nature\n", p.NatureName()))
	// Print IVs
	sb.WriteString(fmt.Sprintf("\tIVs: %d HP / %d Atk / %d Def / %d SpA / %d SpD / %d Spe\n",
		p.hpIv, p.atkIv, p.defIv, p.speIv, p.spaIv, p.spdIv))
	// Print Moves
	sb.WriteString("\t- ")
	moves := p.Moves()
	sb.WriteString(strings.Join(moves, "\n\t- "))
	return sb.String()
}

// 48 bytes (by default)
// 12 bytes * 4 sections
type PokemonData struct {
	// G section
	species  uint16 // 11 bits
	teraType uint16 // 5 bits
	item     uint16 // 10 bits
	// unused_02  uint16 // 6 bits
	experience   uint32 // 21 bits
	nickname11th string // uint32 // 8 bits
	// unused_04 uint32 // 3 bits
	ppBonuses    uint8
	friendship   uint8
	pokeball     uint16 // 6 bits
	nickname12th string // uint16 // 8 bits
	// unused_0A    uint16 // 2 bits
	// A section
	move1             uint16 // 11 bits
	evolutionTracker1 uint16 // 5 bits
	move2             uint16 // 11 bits
	evolutionTracker2 uint16 // 5 bits
	move3             uint16 // 11 bits
	// unused_04 uint16 // 5 bits
	move4 uint16 // 11 bits
	// unused_06 uint16 // 3 bits
	hyperTrainedHp  uint8 // 1 bit
	hyperTrainedAtk uint8 // 1 bit
	pp1             uint8 // 7 bits
	hyperTrainedDef uint8 // 1 bit
	pp2             uint8 // 7 bits
	hyperTrainedSpe uint8 // 1 bit
	pp3             uint8 // 7 bits
	hyperTrainedSpa uint8 // 1 bit
	pp4             uint8 // 7 bits
	hyperTrainedSpd uint8 // 1 bit
	// E section
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
	sheen     uint8
	// M section
	pokerus     uint8 // TODO interpret
	metLocation uint8
	// --- 16 bits of originsInfo
	metLevel     uint16 // 7 bits
	metGame      uint16 // 4 bits
	dynamaxLevel uint16 // 4 bits
	otGender     string // 1 bit
	// --- 32 bits of IVs, isEgg flag, gMaxFactor flag
	hpIv       uint32 // 5 bits
	atkIv      uint32 // 5 bits
	defIv      uint32 // 5 bits
	speIv      uint32 // 5 bits
	spaIv      uint32 // 5 bits
	spdIv      uint32 // 5 bits
	isEgg      bool
	gMaxFactor bool
	// --- Ribbons & other flags
	coolRibbon     uint16 // 3 bits
	beautyRibbon   uint16 // 3 bits
	cuteRibbon     uint16 // 3 bits
	smartRibbon    uint16 // 3 bits
	toughRibbon    uint16 // 3 bits
	championRibbon bool   // 1 bit
	winningRibbon  bool   // 1 bit
	victoryRibbon  bool   // 1 bit
	artistRibbon   bool   // 1 bit
	effortRibbon   bool   // 1 bit
	marineRibbon   bool   // 1 bit
	landRibbon     bool   // 1 bit
	skyRibbon      bool   // 1 bit
	countryRibbon  bool   // 1 bit
	nationalRibbon bool   // 1 bit
	earthRibbon    bool   // 1 bit
	worldRibbon    bool   // 1 bit
	isShadow       bool   // 1 bit
	// unused_0B      bool // 1 bit
	abilityNum             uint8 // 2 bits
	modernFatefulEncounter bool  // 1 bit;
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
	p.parseGSection(sections['G'])
	p.parseASection(sections['A'])
	p.parseESection(sections['E'])
	p.parseMSection(sections['M'])
}

func (p *PokemonData) parseGSection(section []byte) {
	// first 11 bits
	p.species = binary.LittleEndian.Uint16(section[0:2]) & 0x7FF
	// last 5 bits
	p.teraType = binary.LittleEndian.Uint16(section[0:2]) & 0x1F
	// first 10 bits
	p.item = binary.LittleEndian.Uint16(section[2:4]) & 0x3FF
	// last 6 bits
	// p.unused_02 = binary.LittleEndian.Uint32(section[2:4]) & 0x3F
	// first 21 bits
	p.experience = binary.LittleEndian.Uint32(section[4:8]) & 0x1FFFFF
	// the 8 bits after the 21th bit
	p.nickname11th = readString([]byte{byte((binary.LittleEndian.Uint32(section[4:8]) >> 21) & 0xFF)})
	// last 3 bits
	// p.unused_04 = binary.LittleEndian.Uint32(section[4:8]) >> 29
	p.ppBonuses = section[8]
	p.friendship = section[9]
	// first 6 bits
	p.pokeball = binary.LittleEndian.Uint16(section[10:12]) & 0x3F
	// last 8 bits
	p.nickname12th = readString([]byte{byte((binary.LittleEndian.Uint16(section[10:12]) >> 6) & 0xFF)})
	// last 2 bits
	// p.unused_0A = binary.LittleEndian.Uint16(section[10:12]) >> 14
}

func (p *PokemonData) parseASection(section []byte) {
	// first 11 bits
	p.move1 = binary.LittleEndian.Uint16(section[0:2]) & 0x7FF
	// last 5 bits
	p.evolutionTracker1 = (binary.LittleEndian.Uint16(section[0:2]) >> 11) & 0x1F
	// first 11 bits
	p.move2 = binary.LittleEndian.Uint16(section[2:4]) & 0x7FF
	// last 5 bits
	p.evolutionTracker2 = (binary.LittleEndian.Uint16(section[2:4]) >> 11) & 0x1F
	// first 11 bits
	p.move3 = binary.LittleEndian.Uint16(section[4:6]) & 0x7FF
	// last 5 bits
	// p.unused_04 = (binary.LittleEndian.Uint16(section[4:6]) >> 11) & 0x1F
	// first 11 bits
	p.move4 = binary.LittleEndian.Uint16(section[6:8]) & 0x7FF
	// 4th 5th 6th bits
	// p.unused_06 = (section[7] >> 3) & 0x07
	// 7th bit
	p.hyperTrainedHp = (section[7] >> 3) & 0x01
	// 8th bit
	p.hyperTrainedAtk = (section[7] >> 4) & 0x01
	// first 7 bits
	p.pp1 = section[8] & 0x7F
	// last 1 bit
	p.hyperTrainedDef = (section[8] >> 7) & 0x01
	// first 7 bits
	p.pp2 = section[9] & 0x7F
	// last 1 bit
	p.hyperTrainedSpe = (section[9] >> 7) & 0x01
	// first 7 bits
	p.pp3 = section[10] & 0x7F
	// last 1 bit
	p.hyperTrainedSpa = (section[10] >> 7) & 0x01
	// first 7 bits
	p.pp4 = section[11] & 0x7F
	// last 1 bit
	p.hyperTrainedSpd = (section[11] >> 7) & 0x01
}

func (p *PokemonData) parseESection(section []byte) {
	p.hpEv = section[0]
	p.atkEv = section[1]
	p.defEv = section[2]
	p.speEv = section[3]
	p.spaEv = section[4]
	p.spdEv = section[5]
	p.coolness = section[6]
	p.beauty = section[7]
	p.cuteness = section[8]
	p.smartness = section[9]
	p.toughness = section[10]
	p.sheen = section[11]
}

func (p *PokemonData) parseMSection(section []byte) {
	p.pokerus = section[0]
	p.metLocation = section[1]
	// first 7 bits
	p.metLevel = binary.LittleEndian.Uint16(section[2:4]) & 0x7F
	// 8th 9th 10th 11th bits
	p.metGame = binary.LittleEndian.Uint16(section[2:4]) >> 7
	// 12th 13th 14th 15th bits
	p.dynamaxLevel = binary.LittleEndian.Uint16(section[2:4]) >> 12
	// last bit
	if binary.LittleEndian.Uint16(section[2:4])>>15 == 0 {
		p.otGender = "boy"
	} else {
		p.otGender = "girl"
	}
	p.hpIv = binary.LittleEndian.Uint32(section[4:8]) & 0x1F
	p.atkIv = binary.LittleEndian.Uint32(section[4:8]) >> 5 & 0x1F
	p.defIv = binary.LittleEndian.Uint32(section[4:8]) >> 10 & 0x1F
	p.speIv = binary.LittleEndian.Uint32(section[4:8]) >> 15 & 0x1F
	p.spaIv = binary.LittleEndian.Uint32(section[4:8]) >> 20 & 0x1F
	p.spdIv = binary.LittleEndian.Uint32(section[4:8]) >> 25 & 0x1F
	p.isEgg = section[8]&0x1 == 1
	p.gMaxFactor = section[8]&0x2 == 1
	p.coolRibbon = binary.LittleEndian.Uint16(section[8:10]) & 0x7
	p.beautyRibbon = binary.LittleEndian.Uint16(section[8:10]) >> 3 & 0x7
	p.cuteRibbon = binary.LittleEndian.Uint16(section[8:10]) >> 6 & 0x7
	p.smartRibbon = binary.LittleEndian.Uint16(section[8:10]) >> 9 & 0x7
	p.toughRibbon = binary.LittleEndian.Uint16(section[8:10]) >> 12 & 0x7
	p.championRibbon = section[10]&0x1 == 1
	p.winningRibbon = section[10]&0x2 == 1
	p.victoryRibbon = section[10]&0x4 == 1
	p.artistRibbon = section[10]&0x8 == 1
	p.effortRibbon = section[10]&0x10 == 1
	p.marineRibbon = section[10]&0x20 == 1
	p.landRibbon = section[10]&0x40 == 1
	p.skyRibbon = section[10]&0x80 == 1
	p.countryRibbon = section[11]&0x1 == 1
	p.nationalRibbon = section[11]&0x2 == 1
	p.earthRibbon = section[11]&0x4 == 1
	p.worldRibbon = section[11]&0x8 == 1
	p.isShadow = section[11]&0x10 == 1
	// p.unused_0B = section[11]&0x20 == 1
	p.abilityNum = section[11] >> 6 & 0x3
	p.modernFatefulEncounter = section[11]&0x40 == 1
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

package sav

import (
	"acheron-save-parser/gba"
	"acheron-save-parser/utils"
	"encoding/binary"
	"fmt"
	"math"
	"strings"
)

type Team struct {
	Size    int
	Pokemon []Pokemon // 600 bytes
}

func (t *Team) new(section []byte) {
	// go at offset 0x234 for size
	t.Size = int(binary.LittleEndian.Uint32(section[0x234:0x238]))
	// between 0x238 and 0x490 is the pokemon data
	t.Pokemon = make([]Pokemon, t.Size)
	for i := 0; i < int(t.Size); i++ {
		t.Pokemon[i].new(section[0x238+(i*100) : 0x238+((i+1)*100)])
	}
}

type Pokemon struct {
	PersonalityValue     uint32
	OtId                 uint32
	nickname             string // 10 bytes
	language             uint8  // 3 bits
	HiddenNatureModifier uint8  // 5 bits
	IsBadEgg             bool
	HasSpecies           bool
	IsEgg                bool
	// blockBoxRS bool
	DaysSinceFormChange uint8 // 3 bits - 7 days.
	// unused_13 bool
	OtName           string // 7 bytes
	Markings         uint8  // 4 bits
	CompressedStatus uint8  // 4 bits // TODO ignored
	Checksum         uint16
	HpLost           uint16 // 14 bits; // 16383 HP.
	ShinyModifier    bool
	// unused_1E        bool
	PokemonData            // 48 bytes - encrypted
	statusCondition uint32 // TODO ignored
	level           uint8
	MailId          uint8
	CurrentHp       uint16
	TotalHp         uint16
	Attack          uint16
	Defense         uint16
	Speed           uint16
	SpecialAttack   uint16
	SpecialDefense  uint16
}

func (p *Pokemon) new(section []byte) {
	p.PersonalityValue = binary.LittleEndian.Uint32(section[0:4])
	p.OtId = binary.LittleEndian.Uint32(section[4:8])
	p.nickname = utils.DecodeGFString(section[8:18])
	// first 3 bits of [18]
	p.language = section[18] & 0x7 // should always be 2
	// last 5 bits of [18]
	p.HiddenNatureModifier = (section[18] >> 3) & 0x1F
	// first bit of [19]
	p.IsBadEgg = section[19]&0x1 == 1
	// second bit of [19]
	p.HasSpecies = section[19]&0x2 == 1
	// third bit of [19]
	p.IsEgg = section[19]&0x4 == 1
	// fourth bit of [19]
	// p.blockBoxRS = section[19] & 0x38
	// fifth, sixth and seventh bit of [19]
	p.DaysSinceFormChange = section[19] & 0xC0
	// last bit of [19]
	// p.unused_13 = section[19]&0x80 == 1
	p.OtName = utils.DecodeGFString(section[20:28])
	// first 4 bits of [28]
	p.Markings = section[28] & 0xF
	// last 4 bits of [28]
	p.CompressedStatus = section[28] & 0xF0
	p.Checksum = binary.LittleEndian.Uint16(section[30:32])
	// first 14 bits of [32:34]
	p.HpLost = binary.LittleEndian.Uint16(section[32:34]) & 0x3FFF
	// second to last 1 bit of [32:34]
	p.ShinyModifier = section[34]&0x80 == 1
	// last bit of [32:34]
	// p.unused_1E = section[34]&0x40 == 1
	p.PokemonData.new(section[32:80], p.OtId, p.PersonalityValue)
	p.statusCondition = binary.LittleEndian.Uint32(section[80:84])
	p.level = section[84]
	p.MailId = section[85]
	p.CurrentHp = binary.LittleEndian.Uint16(section[86:88])
	p.TotalHp = binary.LittleEndian.Uint16(section[88:90])
	p.Attack = binary.LittleEndian.Uint16(section[90:92])
	p.Defense = binary.LittleEndian.Uint16(section[92:94])
	p.Speed = binary.LittleEndian.Uint16(section[94:96])
	p.SpecialAttack = binary.LittleEndian.Uint16(section[96:98])
	p.SpecialDefense = binary.LittleEndian.Uint16(section[98:100])
}

func (p *Pokemon) newBoxed(section []byte) {
	p.PersonalityValue = binary.LittleEndian.Uint32(section[0:4])
	p.OtId = binary.LittleEndian.Uint32(section[4:8])
	p.nickname = utils.DecodeGFString(section[8:18])
	// first 3 bits of [18]
	p.language = section[18] & 0x7 // should always be 2
	// last 5 bits of [18]
	p.HiddenNatureModifier = (section[18] >> 3) & 0x1F
	// first bit of [19]
	p.IsBadEgg = section[19]&0x1 == 1
	// second bit of [19]
	p.HasSpecies = section[19]&0x2 == 1
	// third bit of [19]
	p.IsEgg = section[19]&0x4 == 1
	// fourth bit of [19]
	// p.blockBoxRS = section[19] & 0x38
	// fifth, sixth and seventh bit of [19]
	p.DaysSinceFormChange = section[19] & 0xC0
	// last bit of [19]
	// p.unused_13 = section[19]&0x80 == 1
	p.OtName = utils.DecodeGFString(section[20:28])
	// first 4 bits of [28]
	p.Markings = section[28] & 0xF
	// last 4 bits of [28]
	p.CompressedStatus = section[28] & 0xF0
	p.Checksum = binary.LittleEndian.Uint16(section[30:32])
	// first 14 bits of [32:34]
	p.HpLost = binary.LittleEndian.Uint16(section[32:34]) & 0x3FFF
	// second to last 1 bit of [32:34]
	p.ShinyModifier = section[34]&0x80 == 1
	// last bit of [32:34]
	// p.unused_1E = section[34]&0x40 == 1
	p.PokemonData.new(section[32:80], p.OtId, p.PersonalityValue)
	// calculate the level using experience and experienceGroup
	p.level = uint8(calculateLevel(int(p.Experience), p.ExperienceGroup()))
}

func (p *Pokemon) SpeciesName() string {
	return gba.Species[p.Species].SpeciesName
}
func (p *Pokemon) ItemName() string {
	return gba.Items[p.Item].Name
}
func (p *Pokemon) Level() int {
	return int(p.level)
}
func (p *Pokemon) AbilityName() string {
	return gba.Abilities[gba.Species[p.Species].Abilities[p.abilityNum]].Name
}
func (p *Pokemon) NatureName() string {
	return gba.Natures[uint8(p.PersonalityValue%25)].Name
}
func (p *Pokemon) Nickname() string {
	return p.nickname + p.nickname11th + p.nickname12th
}
func (p *Pokemon) ExperienceGroup() int {
	return int(gba.Species[p.Species].GrowthRate)
}
func (p *Pokemon) Moves() []string {
	var moves []string
	if p.move1 != 0 {
		moves = append(moves, gba.Moves[p.move1].Name)
	}
	if p.move2 != 0 {
		moves = append(moves, gba.Moves[p.move2].Name)
	}
	if p.move3 != 0 {
		moves = append(moves, gba.Moves[p.move3].Name)
	}
	if p.move4 != 0 {
		moves = append(moves, gba.Moves[p.move4].Name)
	}
	return moves
}
func (p *Pokemon) SDExportFormat() string {
	var sb strings.Builder

	// Pokémon name and item
	if p.nickname != "" {
		sb.WriteString(fmt.Sprintf("%s (%s)", p.Nickname(), p.SpeciesName()))
	} else {
		sb.WriteString(p.SpeciesName())
	}

	if p.Item != 0 {
		sb.WriteString(fmt.Sprintf(" @ %s", p.ItemName()))
	}
	sb.WriteString("\n")

	// Level, Nature, and Ability
	sb.WriteString(fmt.Sprintf("Level: %d\n", p.level))
	sb.WriteString(fmt.Sprintf("%s Nature\n", p.NatureName()))
	sb.WriteString(fmt.Sprintf("Ability: %s\n", p.AbilityName()))

	// EVs
	evs := []string{}

	if p.HpEv > 0 {
		evs = append(evs, fmt.Sprintf("%d HP", p.HpEv))
	}
	if p.AtkEv > 0 {
		evs = append(evs, fmt.Sprintf("%d Atk", p.AtkEv))
	}
	if p.DefEv > 0 {
		evs = append(evs, fmt.Sprintf("%d Def", p.DefEv))
	}
	if p.SpaEv > 0 {
		evs = append(evs, fmt.Sprintf("%d SpA", p.SpaEv))
	}
	if p.SpdEv > 0 {
		evs = append(evs, fmt.Sprintf("%d SpD", p.SpdEv))
	}
	if p.SpeEv > 0 {
		evs = append(evs, fmt.Sprintf("%d Spe", p.SpeEv))
	}

	if len(evs) > 0 {
		sb.WriteString(fmt.Sprintf("EVs: %s\n", strings.Join(evs, " / ")))
	}

	// IVs
	sb.WriteString(fmt.Sprintf("IVs: %d HP / %d Atk / %d Def / %d SpA / %d SpD / %d Spe\n",
		p.HpIv, p.AtkIv, p.DefIv, p.SpaIv, p.SpdIv, p.SpeIv))

	// Moves
	sb.WriteString(fmt.Sprintf("- %s\n", strings.Join(p.Moves(), "\n- ")))

	return sb.String()
}

func calculateLevel(experience int, experienceGroup int) int {
	for level := 1; level <= 100; level++ {
		expRequired := experienceForLevel(level+1, experienceGroup)
		if expRequired > experience {
			return level
		}
	}
	return 100 // If experience is higher than required for level 100
}

func experienceForLevel(level int, experienceGroup int) int {
	switch experienceGroup {
	case 0: // Medium Fast
		return level * level * level
	case 1: // Erratic
		if level <= 50 {
			return level * level * level * (100 - level) / 50
		} else if level <= 68 {
			return level * level * level * (150 - level) / 100
		} else if level <= 98 {
			return level * level * level * int((1911-level*10)/3) / 500
		} else {
			return level * level * level * (160 - level) / 100
		}
	case 2: // Fluctuating
		if level <= 15 {
			return level * level * level * ((level+1)/3 + 24) / 50
		} else if level <= 36 {
			return level * level * level * (level + 14) / 50
		} else {
			return level * level * level * (level/2 + 32) / 50
		}
	case 3: // Medium Slow
		return int(math.Round(float64(6*level*level*level)/5 - float64(15*level*level) + float64(100*level) - 140))
	case 4: // Fast
		return 4 * level * level * level / 5
	case 5: // Slow
		return 5 * level * level * level / 4
	default:
		return 0 // Invalid group
	}
}

// 48 bytes (by default)
// 12 bytes * 4 sections
type PokemonData struct {
	// G section
	Species  uint16 // 11 bits
	TeraType uint16 // 5 bits
	Item     uint16 // 10 bits
	// unused_02  uint16 // 6 bits
	Experience   uint32 // 21 bits
	nickname11th string // uint32 // 8 bits
	// unused_04 uint32 // 3 bits
	PpBonuses    uint8
	Friendship   uint8
	Pokeball     uint16 // 6 bits
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
	HpEv      uint8
	AtkEv     uint8
	DefEv     uint8
	SpeEv     uint8
	SpaEv     uint8
	SpdEv     uint8
	Coolness  uint8
	Beauty    uint8
	Cuteness  uint8
	Smartness uint8
	Toughness uint8
	Sheen     uint8
	// M section
	pokerus     uint8 // TODO interpret
	MetLocation uint8
	// --- 16 bits of originsInfo
	MetLevel     uint16 // 7 bits
	MetGame      uint16 // 4 bits
	DynamaxLevel uint16 // 4 bits
	OtGender     string // 1 bit
	// --- 32 bits of IVs, isEgg flag, gMaxFactor flag
	HpIv       uint32 // 5 bits
	AtkIv      uint32 // 5 bits
	DefIv      uint32 // 5 bits
	SpeIv      uint32 // 5 bits
	SpaIv      uint32 // 5 bits
	SpdIv      uint32 // 5 bits
	IsEgg      bool
	GMaxFactor bool
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
	modernFatefulEncounter bool  // 1 bit
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
	p.Species = binary.LittleEndian.Uint16(section[0:2]) & 0x7FF
	// last 5 bits
	p.TeraType = binary.LittleEndian.Uint16(section[0:2]) & 0x1F
	// first 10 bits
	p.Item = binary.LittleEndian.Uint16(section[2:4]) & 0x3FF
	// last 6 bits
	// p.unused_02 = binary.LittleEndian.Uint32(section[2:4]) & 0x3F
	// first 21 bits
	p.Experience = binary.LittleEndian.Uint32(section[4:8]) & 0x1FFFFF
	// the 8 bits after the 21th bit
	p.nickname11th = utils.DecodeGFString([]byte{byte((binary.LittleEndian.Uint32(section[4:8]) >> 21) & 0xFF)})
	// last 3 bits
	// p.unused_04 = binary.LittleEndian.Uint32(section[4:8]) >> 29
	p.PpBonuses = section[8]
	p.Friendship = section[9]
	// first 6 bits
	p.Pokeball = binary.LittleEndian.Uint16(section[10:12]) & 0x3F
	// last 8 bits
	p.nickname12th = utils.DecodeGFString([]byte{byte((binary.LittleEndian.Uint16(section[10:12]) >> 6) & 0xFF)})
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
	p.HpEv = section[0]
	p.AtkEv = section[1]
	p.DefEv = section[2]
	p.SpeEv = section[3]
	p.SpaEv = section[4]
	p.SpdEv = section[5]
	p.Coolness = section[6]
	p.Beauty = section[7]
	p.Cuteness = section[8]
	p.Smartness = section[9]
	p.Toughness = section[10]
	p.Sheen = section[11]
}

func (p *PokemonData) parseMSection(section []byte) {
	p.pokerus = section[0]
	p.MetLocation = section[1]
	// first 7 bits
	p.MetLevel = binary.LittleEndian.Uint16(section[2:4]) & 0x7F
	// 8th 9th 10th 11th bits
	p.MetGame = binary.LittleEndian.Uint16(section[2:4]) >> 7
	// 12th 13th 14th 15th bits
	p.DynamaxLevel = binary.LittleEndian.Uint16(section[2:4]) >> 12
	// last bit
	if binary.LittleEndian.Uint16(section[2:4])>>15 == 0 {
		p.OtGender = "boy"
	} else {
		p.OtGender = "girl"
	}
	p.HpIv = binary.LittleEndian.Uint32(section[4:8]) & 0x1F
	p.AtkIv = binary.LittleEndian.Uint32(section[4:8]) >> 5 & 0x1F
	p.DefIv = binary.LittleEndian.Uint32(section[4:8]) >> 10 & 0x1F
	p.SpeIv = binary.LittleEndian.Uint32(section[4:8]) >> 15 & 0x1F
	p.SpaIv = binary.LittleEndian.Uint32(section[4:8]) >> 20 & 0x1F
	p.SpdIv = binary.LittleEndian.Uint32(section[4:8]) >> 25 & 0x1F
	p.IsEgg = section[8]&0x1 == 1
	p.GMaxFactor = section[8]&0x2 == 1
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

package gba

import (
	"acheron-save-parser/utils"
	"encoding/binary"
)

type MoveData struct {
	namePtr        uint32
	Name           string
	descriptionPtr uint32
	Description    string
	Effect         uint16
	// -- same 2 bytes
	Type     uint8  // 5 bits
	Category uint8  // 2 bits
	Power    uint16 // 9 bits
	// --
	// -- same 2 bytes
	Accuracy uint8  // 7 bits
	Target   uint16 // 9 bits
	// --
	Pp    uint8
	ZMove uint32 // zMoveEffect / PowerOverride union
	// -- same 4 bytes
	Priority             int8  // 4 bits
	Recoil               uint8 // 7 bits
	StrikeCount          uint8 // 4 bits
	CriticalHitStage     uint8 // 2 bits
	AlwaysCriticalHit    bool
	NumAdditionalEffects uint8 // 2 bits
	// Flags
	MakesContact                      bool
	IgnoresProtect                    bool
	MagicCoatAffected                 bool
	SnatchAffected                    bool
	IgnoresKingsRock                  bool
	PunchingMove                      bool
	BitingMove                        bool
	PulseMove                         bool
	SoundMove                         bool
	BallisticMove                     bool
	PowderMove                        bool
	DanceMove                         bool
	WindMove                          bool
	SlicingMove                       bool
	HealingMove                       bool
	MinimizeDoubleDamage              bool
	IgnoresTargetAbility              bool
	IgnoresTargetDefenseEvasionStages bool
	DamagesUnderground                bool
	DamagesUnderwater                 bool
	DamagesAirborne                   bool
	AirborneDoubleDamage              bool
	IgnoreTypeIfFlyingAndUngrounded   bool
	ThawsUser                         bool
	IgnoresSubstitute                 bool
	ForcePressure                     bool
	CantUseTwice                      bool

	// Ban flags
	GravityBanned      bool
	MirrorMoveBanned   bool
	MeFirstBanned      bool
	MimicBanned        bool
	MetronomeBanned    bool
	CopycatBanned      bool
	AssistBanned       bool // Matches same moves as copycatBanned + semi-invulnerable moves and Mirror Coat.
	SleepTalkBanned    bool
	InstructBanned     bool
	EncoreBanned       bool
	ParentalBondBanned bool
	SkyBattleBanned    bool
	SketchBanned       bool
	// padding // 5 bits // end of word

	Argument uint32
	// primary/secondary effects
	additionalEffectsPtr uint32
	AdditionalEffects    []*MoveAdditionalEffect
	// contest parameters
	ContestEffect         uint8
	ContestCategory       uint8 // 3 bits
	ContestComboStarterId uint8
	ContestComboMoves     [5]uint8
	battleAnimScriptPtr   uint32
}

const (
	MOVE_INFO_SIZE               = 52
	MOVE_ADDITIONAL_EFFECTS_SIZE = 4
)

func ParseMovesInfoBytes(data []byte, offset int, count int) []*MoveData {
	moves := make([]*MoveData, count)
	for i := 0; i < count; i++ {
		m := &MoveData{}
		m.new(data[offset+i*MOVE_INFO_SIZE : offset+i*MOVE_INFO_SIZE+MOVE_INFO_SIZE])
		moves[i] = m
		m.Name = utils.DecodePointerString(data, m.namePtr)
		m.Description = utils.DecodePointerString(data, m.descriptionPtr)
		if m.additionalEffectsPtr != BAD_POINTER {
			m.AdditionalEffects = ParseMoveAdditionalEffects(data, m.additionalEffectsPtr, m.NumAdditionalEffects)
		}
	}
	return moves
}

func (m *MoveData) new(section []byte /* 52 bytes */) {
	m.namePtr = binary.LittleEndian.Uint32(section[0:4]) - POINTER_OFFSET
	m.descriptionPtr = binary.LittleEndian.Uint32(section[4:8]) - POINTER_OFFSET
	m.Effect = binary.LittleEndian.Uint16(section[8:10])
	// first 5 bits
	m.Type = section[10] & 0x1F
	// next 2 bits
	m.Category = section[10] & 0x60 >> 5
	// last bit of [10] added to next byte
	m.Power = uint16(section[11]<<1 | section[10]&0x80>>7)
	// first 7 bits
	m.Accuracy = section[12] & 0x7F
	// last bit of [12] added to next byte
	m.Target = uint16(section[13]<<1 | section[12]&0x80>>7)
	m.Pp = section[14]
	// [15] is padding
	m.ZMove = binary.LittleEndian.Uint32(section[16:20])
	// first 4 bits
	m.Priority = int8(section[20] & 0x0F)
	// last 4 bits and 3 first bits of next byte
	m.Recoil = uint8((section[20]&0xF0)>>4 | (section[21]&0x7)<<4)
	// next 4 bits (after the first 3) of next byte
	m.StrikeCount = (section[21] & 0x78 >> 3)
	// last bit of [21], first bit of [22]
	m.CriticalHitStage = section[21]&0x80>>7 | section[22]&0x01
	// second bit of [22]
	m.AlwaysCriticalHit = section[22]&0b10>>1 == 1
	// 3rd and 4th bits of [22]
	m.NumAdditionalEffects = section[22] & 0x0C >> 2
	// 5th bit of [22]
	m.MakesContact = section[22]&0b10000>>4 == 1
	// 6th bit of [22]
	m.IgnoresProtect = section[22]&0b100000>>5 == 1
	// 7th bit of [22]
	m.MagicCoatAffected = section[22]&0b1000000>>6 == 1
	// 8th bit of [22]
	m.SnatchAffected = section[22]&0b10000000>>7 == 1
	// first bit of [23]
	m.IgnoresKingsRock = section[23]&0b1 == 1
	// second bit of [23]
	m.PunchingMove = section[23]&0b10>>1 == 1
	// third bit of [23]
	m.BitingMove = section[23]&0b100>>2 == 1
	// fourth bit of [23]
	m.PulseMove = section[23]&0b1000>>3 == 1
	// fifth bit of [23]
	m.SoundMove = section[23]&0b10000>>4 == 1
	// sixth bit of [23]
	m.BallisticMove = section[23]&0b100000>>5 == 1
	// seventh bit of [23]
	m.PowderMove = section[23]&0b1000000>>6 == 1
	// 8th bit of [23]
	m.DanceMove = section[23]&0b10000000>>7 == 1
	// 1st bit of [24]
	m.WindMove = section[24]&0b1 == 1
	// 2nd bit of [24]
	m.SlicingMove = section[24]&0b10>>1 == 1
	// 3rd bit of [24]
	m.HealingMove = section[24]&0b100>>2 == 1
	// 4th bit of [24]
	m.MinimizeDoubleDamage = section[24]&0b1000>>3 == 1
	// 5th bit of [24]
	m.IgnoresTargetAbility = section[24]&0b10000>>4 == 1
	// 6th bit of [24]
	m.IgnoresTargetDefenseEvasionStages = section[24]&0b100000>>5 == 1
	// 7th bit of [24]
	m.DamagesUnderground = section[24]&0b1000000>>6 == 1
	// 8th bit of [24]
	m.DamagesUnderwater = section[24]&0b10000000>>7 == 1
	// 1st bit of [25]
	m.DamagesAirborne = section[25]&0b1 == 1
	// 2nd bit of [25]
	m.AirborneDoubleDamage = section[25]&0b10>>1 == 1
	// 3rd bit of [25]
	m.IgnoreTypeIfFlyingAndUngrounded = section[25]&0b100>>2 == 1
	// 4th bit of [25]
	m.ThawsUser = section[25]&0b1000>>3 == 1
	// 5th bit of [25]
	m.IgnoresSubstitute = section[25]&0b10000>>4 == 1
	// 6th bit of [25]
	m.ForcePressure = section[25]&0b100000>>5 == 1
	// 7th bit of [25]
	m.CantUseTwice = section[25]&0b1000000>>6 == 1
	// 8th bit of [25]
	m.GravityBanned = section[25]&0b10000000>>7 == 1
	// 1st bit of [26]
	m.MirrorMoveBanned = section[26]&0b1 == 1
	// 2nd bit of [26]
	m.MeFirstBanned = section[26]&0b10>>1 == 1
	// 3rd bit of [26]
	m.MimicBanned = section[26]&0b100>>2 == 1
	// 4th bit of [26]
	m.MetronomeBanned = section[26]&0b1000>>3 == 1
	// 5th bit of [26]
	m.CopycatBanned = section[26]&0b10000>>4 == 1
	// 6th bit of [26]
	m.AssistBanned = section[26]&0b100000>>5 == 1
	// 7th bit of [26]
	m.SleepTalkBanned = section[26]&0b1000000>>6 == 1
	// 8th bit of [26]
	m.InstructBanned = section[26]&0b10000000>>7 == 1
	// 1st bit of [27]
	m.EncoreBanned = section[27]&0b1 == 1
	// 2nd bit of [27]
	m.ParentalBondBanned = section[27]&0b10>>1 == 1
	// 3rd bit of [27]
	m.SkyBattleBanned = section[27]&0b100>>2 == 1
	// 4th bit of [27]
	m.SketchBanned = section[27]&0b1000>>3 == 1
	// last 4 bits are unused
	// 4 bytes of padding
	m.Argument = binary.LittleEndian.Uint32(section[32:36])
	m.additionalEffectsPtr = binary.LittleEndian.Uint32(section[36:40]) - POINTER_OFFSET
	m.ContestEffect = section[40]
	// first 3 bits
	m.ContestCategory = section[41] & 0x07
	// last (most significant) 5 bits of [41] + first 3 bits of [42]
	m.ContestComboStarterId = section[41]&0b1111000<<3 | section[42]&0b111
	// rest of [42] is padding
	copy(m.ContestComboMoves[:], section[43:48])
	m.battleAnimScriptPtr = binary.LittleEndian.Uint32(section[48:52]) - POINTER_OFFSET
}

type MoveAdditionalEffect struct {
	MoveEffect              uint16
	Self                    bool
	OnlyIfTargetRaisedStats bool
	OnChargeTurnOnly        bool
	Chance                  uint8
}

func ParseMoveAdditionalEffects(data []byte, offset uint32, count uint8) []*MoveAdditionalEffect {
	effects := make([]*MoveAdditionalEffect, 0)
	for i := 0; i < int(count); i++ {
		ix := offset + uint32(i)*MOVE_ADDITIONAL_EFFECTS_SIZE
		effect := &MoveAdditionalEffect{
			MoveEffect:              binary.LittleEndian.Uint16(data[ix : ix+2]),
			Self:                    data[ix+2]&0x01 != 0,
			OnlyIfTargetRaisedStats: data[ix+2]&0x02 != 0,
			OnChargeTurnOnly:        data[ix+2]&0x04 != 0,
			Chance:                  data[ix+3],
		}
		effects = append(effects, effect)
	}
	return effects
}

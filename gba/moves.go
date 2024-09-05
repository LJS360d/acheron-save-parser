package gba

import (
	"encoding/binary"
	"fmt"
)

type MoveData struct {
	NamePtr        uint32
	Name           string
	DescriptionPtr uint32
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
	ZMove uint8 // zMoveEffect / PowerOverride union
	// -- same 4 bytes
	Priority             int8  // 4 bits
	Recoil               uint8 // 7 bits
	StrikeCount          uint8 // 4 bits
	CriticalHitStage     uint8 // 2 bits
	AlwaysCriticalHit    bool
	NumAdditionalEffects uint8 // 2 bits // limited to 3 - don't want to get too crazy
	// (padding?) 12 bits left to complete this word - continues into flags
	// --
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
	SlicingMove                       bool // end of word
	HealingMove                       bool
	MinimizeDoubleDamage              bool
	IgnoresTargetAbility              bool
	IgnoresTargetDefenseEvasionStages bool
	DamagesUnderground                bool
	DamagesUnderwater                 bool
	DamagesAirborne                   bool
	DamagesAirborneDoubleDamage       bool
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
	AdditionalEffectsPtr uint32
	// contest parameters
	ContestEffect         uint8
	ContestCategory       uint8 // 3 bits
	ContestComboStarterId uint8
	ContestComboMoves     [5]uint8
	BattleAnimScriptPtr   uint32
}

const (
	MOVE_INFO_SIZE = 52
)

func ParseMovesInfoBytes(data []byte, offset int, count int) []*MoveData {
	moves := make([]*MoveData, count)
	for i := 0; i < count; i++ {
		m := &MoveData{}
		m.new(data[offset+i*MOVE_INFO_SIZE : offset+i*MOVE_INFO_SIZE+MOVE_INFO_SIZE])
		moves[i] = m
		m.Name = ParsePointerString(data, m.NamePtr)
		m.Description = ParsePointerString(data, m.DescriptionPtr)
		fmt.Println(m.Name, m.Power)
	}
	return moves
}

func (m *MoveData) new(section []byte /* 52 bytes */) {
	m.NamePtr = binary.LittleEndian.Uint32(section[0:4]) - POINTER_OFFSET
	m.DescriptionPtr = binary.LittleEndian.Uint32(section[4:8]) - POINTER_OFFSET
	m.Effect = binary.LittleEndian.Uint16(section[8:10])
	// first 5 bits
	m.Type = section[10] & 0x1F
	// next 2 bits
	m.Category = section[10] & 0x60 >> 5
	// last bit of [10] & next byte
	// TODO figure out move power because what the fuck.
	// uint16(section[11]<<1 | section[10]&0x01) is the closest to being correct
	// but some moves get a 1 added or subtracted to their power, which should not happen
	// cant figure out a pattern, its definitely tied to the binary representation of the number
	// the power should be the last bit of [10] added to the [11] byte but in any way i tried it does not work
	m.Power = uint16(section[11]<<1 | section[10]&0x01)
	// first 7 bits
	m.Accuracy = section[12] & 0x7F
	// last bit of [12] & next byte
	m.Target = binary.LittleEndian.Uint16([]byte{section[12] & 0x80, section[13]})
	m.Pp = section[14]
	m.ZMove = section[15]
	// first 4 bits
	m.Priority = int8(section[16] & 0xF)
	// last 4 bits of [16] and 3 first bits of [17]
	m.Recoil = (section[16] & 0x0F) | (section[17] & 0xE0)
	// next 4 bits (after the first 3) of [17]
	m.StrikeCount = (section[17] & 0x3E) >> 3
}

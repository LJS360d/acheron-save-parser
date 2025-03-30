package gba

import (
	"acheron-save-parser/utils"
	"encoding/binary"
	"log"
)

type MapHeader struct {
	mapLayoutPtr       uint32
	eventsPtr          uint32
	mapScriptsPtr      uint32
	connectionsPtr     uint32
	Music              uint16
	MapLayoutID        uint16
	RegionMapSectionID uint8
	Cave               uint8 // actually its a boolean
	Weather            uint8
	MapType            uint8
	// Filler18           [2]uint8
	AllowCycling  bool
	AllowEscaping bool
	AllowRunning  bool
	ShowMapName   bool //:5
	// 4 unused bits?
	BattleType uint8
}

type RegionMapLocation struct {
	X       uint8
	Y       uint8
	Width   uint8
	Height  uint8
	namePtr uint32
	Name    string

	Id int // added for convenience
}

func ParseRegionMapGroupsMatrix(startOffset int) [][]*MapHeader {
	matrix := make([][]*MapHeader, 0)
	if startOffset == NULL_POINTER || startOffset >= len(Data) {
		log.Printf("[WARN] Region map groups matrix (gMapGroups) offset [%x] is null or out of bounds", startOffset)
		return matrix
	}
	for i, count := range Config.MapGroupCounts {
		if count == 0 {
			continue
		}
		offset := startOffset + i*4
		entryPtr := binary.LittleEndian.Uint32(Data[offset:]) - POINTER_OFFSET
		if entryPtr == NULL_POINTER || int(entryPtr) >= len(Data) {
			log.Printf("[WARN] Region map groups entry offset [%x] is null or out of bounds", startOffset)
			continue
		}
		groupsPtr := binary.LittleEndian.Uint32(Data[entryPtr:]) - POINTER_OFFSET
		groups := parseRegionMapGroups(int(groupsPtr), count)
		matrix = append(matrix, groups)
	}
	return matrix
}

func parseRegionMapGroups(startOffset int, count int) []*MapHeader {
	groups := make([]*MapHeader, 0)
	if startOffset == NULL_POINTER || startOffset >= len(Data) {
		log.Printf("[WARN] Region map groups offset [%x] is null or out of bounds", startOffset)
		return groups
	}
	for i := 0; i < count; i++ {
		offset := startOffset + i*Config.MapGroupSize
		if offset+Config.MapGroupSize >= len(Data) {
			log.Printf("[WARN] Reached end of data in region map groups parsing at %d", i)
			break
		}
		group := &MapHeader{}
		group.mapLayoutPtr = binary.LittleEndian.Uint32(Data[offset:]) - POINTER_OFFSET
		group.eventsPtr = binary.LittleEndian.Uint32(Data[offset+4:]) - POINTER_OFFSET
		group.mapScriptsPtr = binary.LittleEndian.Uint32(Data[offset+8:]) - POINTER_OFFSET
		group.connectionsPtr = binary.LittleEndian.Uint32(Data[offset+12:]) - POINTER_OFFSET
		group.Music = binary.LittleEndian.Uint16(Data[offset+16:])
		group.MapLayoutID = binary.LittleEndian.Uint16(Data[offset+18:])
		group.RegionMapSectionID = Data[offset+20]
		group.Cave = Data[offset+21]
		group.Weather = Data[offset+22]
		group.MapType = Data[offset+23]
		// group.Filler18 = [2]uint8{Data[offset+24], Data[offset+25]}
		group.AllowCycling = Data[offset+26]&0x1 == 1
		group.AllowEscaping = Data[offset+26]&0x2 == 0x2
		group.AllowRunning = Data[offset+26]&0x4 == 0x4
		group.ShowMapName = Data[offset+26]&0x20 == 0x20
		group.BattleType = Data[offset+27]

		if group.MapLayoutID == 0 {
			break
		}

		groups = append(groups, group)
	}
	return groups
}

func ParseRegionMapLocations(startOffset int) []*RegionMapLocation {
	locations := make([]*RegionMapLocation, 0)
	if startOffset == NULL_POINTER || startOffset >= len(Data) {
		log.Printf("[WARN] Region map locations offset [%x] is null or out of bounds", startOffset)
		return locations
	}
	for i := 0; i < Config.RegionLocationsCount; i++ {
		offset := startOffset + i*Config.RegionLocationSize
		location := &RegionMapLocation{}
		location.X = Data[offset]
		location.Y = Data[offset+1]
		location.Width = Data[offset+2]
		location.Height = Data[offset+3]
		location.namePtr = binary.LittleEndian.Uint32(Data[offset+4:]) - POINTER_OFFSET
		location.Name = utils.DecodePointerString(Data, location.namePtr)

		location.Id = i
		locations = append(locations, location)

	}
	return locations
}

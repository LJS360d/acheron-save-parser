package main

import (
	"encoding/binary"
	"os"
	"strings"
)

type Gen3Save struct {
	name      string
	gender    string
	game      string
	teamcount int
	time      int
	team      []Gen3Pokemon
	boxes     []Gen3Pokemon
}

func readstring(text []byte) string {
	chars := "0123456789!?.-         ,  ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	ret := ""
	for _, i := range text {
		c := int(i) - 161
		if c < 0 || c >= len(chars) {
			ret += " "
		} else {
			ret += string(chars[c])
		}
	}
	return strings.TrimSpace(ret)
}

func getindex(data []byte) []int {
	ret := make([]int, 2)

	for i := 0; i < 14; i++ {
		ix := i * 4096
		section := data[ix : ix+4096]
		footer := section[4084:]
		id := binary.LittleEndian.Uint16(footer[0:2])
		index := binary.LittleEndian.Uint32(footer[4:8])

		if id == 0 {
			ds := []uint16{
				binary.LittleEndian.Uint16(section[14:16]),
				uint16(section[16]),
				uint16(section[17]),
			}
			dt := (int(ds[0]) * 3600) + (int(ds[1]) * 60) + int(ds[2])
			ret[0] = int(index)
			ret[1] = dt
			break
		}
	}

	return ret
}

func NewGen3Save(filename string) (*Gen3Save, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	filea := data[0:57344]
	fileb := data[57344:114688]

	a := getindex(filea)
	b := getindex(fileb)

	var savedata []byte
	var time int

	if a[0] < b[0] {
		savedata = fileb
		time = b[1]
	} else if a[0] > b[0] {
		savedata = filea
		time = a[1]
	} else {
		if a[1] < b[1] {
			savedata = fileb
			time = b[1]
		} else {
			savedata = filea
			time = a[1]
		}
	}

	save := &Gen3Save{time: time}
	save.process(savedata)
	return save, nil
}

func (s *Gen3Save) process(savedata []byte) {
	sections := make([][]byte, 14)
	for i := 0; i < 14; i++ {
		ix := i * 4096
		section := savedata[ix : ix+4096]
		footer := section[4084:]
		id := binary.LittleEndian.Uint16(footer[0:2])
		sections[id] = section[0:3968]
	}

	s.name = readstring(sections[0][0:7])
	s.team = []Gen3Pokemon{}
	s.boxes = []Gen3Pokemon{}

	gamecode := binary.LittleEndian.Uint32(sections[0][172:176])
	s.game = "emerald"
	if gamecode == 0 {
		s.game = "rubysapphire"
	}
	if gamecode == 1 {
		s.game = "fireredleafgreen"
	}

	var teamoffset int
	if gamecode == 1 {
		s.teamcount = int(binary.LittleEndian.Uint32(sections[1][52:56]))
		teamoffset = 56
	} else {
		s.teamcount = int(binary.LittleEndian.Uint32(sections[1][564:568]))
		teamoffset = 568
	}

	gender := sections[0][8]
	s.gender = ""
	if gender == 0 {
		s.gender = "boy"
	}
	if gender == 1 {
		s.gender = "girl"
	}

	dex := make([]byte, 0, 33600)
	for i := 5; i < 14; i++ {
		dex = append(dex, sections[i]...)
	}
	dex = dex[4:33604]

	for i := 0; i < 420; i++ {
		pkm := NewGen3Pokemon(dex[i*80 : (i+1)*80])
		if pkm.species.name != "" {
			s.boxes = append(s.boxes, *pkm)
		}
	}

	for i := 0; i < s.teamcount; i++ {
		ofs := teamoffset + (i * 100)
		pkm := NewGen3Pokemon(sections[1][ofs : ofs+80])
		if pkm.species.name != "" {
			s.team = append(s.team, *pkm)
		}
	}
}

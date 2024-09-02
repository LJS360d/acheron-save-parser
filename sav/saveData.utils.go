package sav

import (
	"encoding/binary"
	"strings"
)

const SECTOR_DATA_SIZE = 3968
const SAVE3_CHUNK_SIZE = 116
const FOOTER_SIZE = 12
const SECTOR_SIZE = SECTOR_DATA_SIZE + SAVE3_CHUNK_SIZE + FOOTER_SIZE // 4096
const SAVE_SLOT_SIZE = SECTOR_SIZE * 14                               // 57344

func ReadString(text []byte) string {
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

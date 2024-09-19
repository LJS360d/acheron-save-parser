package utils

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"strings"
)

func DecodePointerString(data []byte, offset uint32) string {
	start := int(offset)
	end := start
	for end < len(data) && data[end] != 0xFF /* 0xFF is the \0 of GF text encoding */ {
		end++
	}
	return DecodeGFString(data[start:end])
}

func DecodeGFString(text []byte) string {
	chars := "0123456789!?.-     '   ,  ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	// TODO finish special chars map
	specialChars := map[int]string{
		-134: "é",
	}
	ret := ""
	for _, i := range text {
		c := int(i) - 161
		if c < 0 || c >= len(chars) {
			if specialChars[c] != "" {
				ret += specialChars[c]
			} else {
				ret += " "
			}
		} else {
			ret += string(chars[c])
		}
	}
	return strings.TrimSpace(ret)
}

func DecodeGFStringParallel(text []byte, numGoroutines int) string {
	chunkSize := len(text) / numGoroutines
	chunks := make([][]byte, numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		start := i * chunkSize
		end := min(start+chunkSize, len(text))
		chunks[i] = text[start:end]
	}
	results := make(chan string, numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func(chunk []byte) {
			result := DecodeGFString(chunk)
			results <- result
		}(chunks[i])
	}
	var finalResult string
	for i := 0; i < numGoroutines; i++ {
		finalResult += <-results
	}

	return strings.TrimSpace(finalResult)
}

func ParsePaletteBytes(data []byte) []color.Color {
	if len(data)%2 != 0 {
		panic("Palette data size must be even")
	}

	var palette []color.Color
	reader := bytes.NewReader(data)

	for i := 0; i < len(data); i += 2 {
		var col uint16
		binary.Read(reader, binary.LittleEndian, &col)

		// GBA palette is BGR555
		b := uint8((col>>10)&0x1F) << 3
		g := uint8((col>>5)&0x1F) << 3
		r := uint8(col&0x1F) << 3

		// Append the color to the palette
		palette = append(palette, color.RGBA{R: r, G: g, B: b, A: 255})
	}

	return palette
}

func DecompressLZ77(input []byte) ([]byte, error) {
	// Ensure valid LZ77 header (first byte should be 0x10)
	if len(input) < 4 || input[0] != 0x10 {
		return nil, fmt.Errorf("invalid LZ77 compressed data")
	}

	var decompressed bytes.Buffer                         // Will hold the decompressed data
	srcOffset := 4                                        // Start reading after the header
	length := int(binary.LittleEndian.Uint32(input) >> 8) // Get the expected decompressed length

	for srcOffset < len(input) { // Continue until we've processed all the input data
		// Read the flag byte (used to determine if the next data is literal or compressed reference)
		flag := input[srcOffset]
		srcOffset++

		// Each flag byte represents 8 blocks of data (8 bits)
		for i := 0; i < 8; i++ {
			// If we've exceeded the source or decompressed length, break early
			if srcOffset >= len(input) || decompressed.Len() >= length {
				break
			}

			// If the bit is 0, the next byte is a raw byte (literal data)
			if (flag & 0x80) == 0 {
				decompressed.WriteByte(input[srcOffset])
				srcOffset++
			} else {
				// Else, the next 2 bytes represent a back-reference to earlier data (compressed)
				if srcOffset+1 >= len(input) {
					return nil, fmt.Errorf("unexpected end of input")
				}

				// Read the 2 bytes
				byte1 := input[srcOffset]
				byte2 := input[srcOffset+1]
				srcOffset += 2

				// The displacement (distance back to copy from) is stored in the 12 least significant bits.
				disp := (int(byte1&0xF) << 8) | int(byte2)
				disp += 1 // Adjust displacement by 1 because it's always +1 in GBA LZ77

				// The length of the copied sequence is stored in the high nibble of byte1 (4 most significant bits).
				copyLen := int(byte1>>4) + 3 // Minimum length is 3

				// Copy the referenced sequence from the decompressed buffer
				for j := 0; j < copyLen; j++ {
					// Calculate the source position in the already decompressed data
					pos := decompressed.Len() - disp
					if pos < 0 {
						return nil, fmt.Errorf("invalid reference position")
					}
					// Write the copied byte to the decompressed buffer
					decompressed.WriteByte(decompressed.Bytes()[pos])
				}
			}
			flag <<= 1 // Shift the flag byte for the next block
		}
	}

	return decompressed.Bytes(), nil
}

func Save4bppImageBytes(bytes []byte, savename string, palette []color.Color, width, height int, transparentBg bool) error {
	img := image.NewPaletted(image.Rect(0, 0, width, height), palette)

	// Set transparency for the first color in the palette
	if transparentBg && len(img.Palette) > 0 {
		transparentColor := img.Palette[0]
		img.Palette[0] = color.RGBA{
			R: transparentColor.(color.RGBA).R,
			G: transparentColor.(color.RGBA).G,
			B: transparentColor.(color.RGBA).B,
			A: 0, // Set alpha to 0 (transparent)
		}
	}

	// Decompress the LZ77 bytes
	decompressed, err := DecompressLZ77(bytes)
	if err != nil {
		return err
	}

	// Number of tiles horizontally and vertically
	tilesX := width / 8
	tilesY := height / 8

	// Total number of tiles
	totalTiles := tilesX * tilesY

	// Start processing the decompressed data
	for tileIndex := 0; tileIndex < totalTiles; tileIndex++ {
		// Determine the tile's (x, y) position in the grid
		tileX := tileIndex % tilesX
		tileY := tileIndex / tilesX

		// Offset for the current tile in the decompressed data (32 bytes per tile)
		tileOffset := tileIndex * 32

		for row := 0; row < 8; row++ {
			for col := 0; col < 8; col += 2 {
				byteIndex := tileOffset + row*4 + col/2
				byteVal := decompressed[byteIndex]

				lowNibble := byteVal & 0xF
				highNibble := (byteVal >> 4) & 0xF

				pixelX := tileX*8 + col
				pixelY := tileY*8 + row

				img.SetColorIndex(pixelX, pixelY, uint8(lowNibble))
				img.SetColorIndex(pixelX+1, pixelY, uint8(highNibble))
			}
		}
	}

	// Save the resulting image as a PNG file
	file, err := os.Create(savename + ".png")
	if err != nil {
		return err
	}
	defer file.Close()

	err = png.Encode(file, img)
	if err != nil {
		return err
	}
	return nil
}

package utils

import (
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

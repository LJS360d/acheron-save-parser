package gba_test

import (
	"rom-parser/gba"
	"testing"
)

func TestGetGbaConfig(t *testing.T) {
	config112 := gba.GetGbaConfig("1.11.2")
	if config112.Match != gba.LatestVersion {
		t.Errorf("Expected 1.11.2 to return LatestVersion, got %s", config112.Match)
	}

	config109 := gba.GetGbaConfig("1.10.9")
	if config109.Match != "1.10.x" {
		t.Errorf("Expected 1.10.9 to return 1.10.x, got %s", config109.Match)
	}

	config92 := gba.GetGbaConfig("1.9.2")
	if config92.Match != "1.9.x" {
		t.Errorf("Expected 1.9.2 to return 1.9.x, got %s", config92.Match)
	}

	configEmpty := gba.GetGbaConfig("")
	if configEmpty.Match != gba.LatestVersion {
		t.Errorf("Expected empty string to return LatestVersion, got %s", configEmpty.Match)
	}
}

package utils

import (
	"testing"
)

func TestNewUUID(t *testing.T) {
	uuid := NewUUID()
	if len(uuid) == 0 {
		t.Error("Expected non-empty UUID string")
	}
	// Basic format check: UUIDs should have 4 hyphens
	hyphenCount := 0
	for _, c := range uuid {
		if c == '-' {
			hyphenCount++
		}
	}
	if hyphenCount != 4 {
		t.Errorf("Expected 4 hyphens in UUID, got %d", hyphenCount)
	}
}

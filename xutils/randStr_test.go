package xutils

import (
	"strings"
	"testing"
)

func TestRandomString(t *testing.T) {
	tests := []struct {
		name   string
		length int
	}{
		{"zero length", 0},
		{"normal length", 10},
		{"long length", 100},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RandomString(tt.length)

			// Check length
			if len(result) != tt.length {
				t.Errorf("String(%d) length = %d; want %d", tt.length, len(result), tt.length)
			}

			// Check if all characters are valid
			for _, char := range result {
				if !strings.ContainsRune(letterBytes, char) {
					t.Errorf("String(%d) contains invalid character: %c", tt.length, char)
				}
			}
		})
	}
}

func TestRandomStringDistribution(t *testing.T) {
	// Test for reasonable distribution of characters
	length := 1000
	result := RandomString(length)

	// Count occurrences of each character
	counts := make(map[rune]int)
	for _, char := range result {
		counts[char]++
	}

	// Check if each character appears at least once
	// This is a basic check, not a statistical test
	for _, char := range letterBytes {
		if counts[char] == 0 {
			t.Errorf("Character %c did not appear in the random string", char)
		}
	}
}

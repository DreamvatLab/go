package xutils

import (
	"testing"
)

func TestGenerateStringID(t *testing.T) {
	// Test generating string ID
	id := GenerateStringID()
	if id == "" {
		t.Error("GenerateStringID() returned empty string")
	}

	// Test if multiple generated IDs are different
	id2 := GenerateStringID()
	if id == id2 {
		t.Error("GenerateStringID() generated duplicate IDs")
	}
}

func TestGenerateUInt64ID(t *testing.T) {
	// Test generating uint64 ID
	id := GenerateUInt64ID()
	if id == 0 {
		t.Error("GenerateUInt64ID() returned zero")
	}

	// Test if multiple generated IDs are different
	id2 := GenerateUInt64ID()
	if id == id2 {
		t.Error("GenerateUInt64ID() generated duplicate IDs")
	}
}

func TestSonyflakeIDGenerator(t *testing.T) {
	// Create a new generator instance
	generator := NewSonyflakeIDGenerator()

	// Test generating string ID
	stringID := generator.GenerateStringID()
	if stringID == "" {
		t.Error("SonyflakeIDGenerator.GenerateStringID() returned empty string")
	}

	// Test generating uint64 ID
	uint64ID := generator.GenerateUInt64ID()
	if uint64ID == 0 {
		t.Error("SonyflakeIDGenerator.GenerateUInt64ID() returned zero")
	}

	// Test if multiple generated IDs are different
	stringID2 := generator.GenerateStringID()
	if stringID == stringID2 {
		t.Error("SonyflakeIDGenerator.GenerateStringID() generated duplicate IDs")
	}

	uint64ID2 := generator.GenerateUInt64ID()
	if uint64ID == uint64ID2 {
		t.Error("SonyflakeIDGenerator.GenerateUInt64ID() generated duplicate IDs")
	}
}

func TestIDGeneratorInterface(t *testing.T) {
	// Test interface implementation
	var _ IIDGenerator = NewSonyflakeIDGenerator()
}

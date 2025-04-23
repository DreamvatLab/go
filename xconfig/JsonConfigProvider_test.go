package xconfig

import (
	"os"
	"testing"
)

func TestNewJsonConfigProvider(t *testing.T) {
	// Test with default config file
	t.Run("Default config file", func(t *testing.T) {
		// Create a temporary config file
		tempFile, err := os.CreateTemp("", "config*.json")
		if err != nil {
			t.Fatalf("Failed to create temp file: %v", err)
		}
		defer os.Remove(tempFile.Name())

		// Write test data to the temp file
		testData := []byte(`{"test": "value"}`)
		if _, err := tempFile.Write(testData); err != nil {
			t.Fatalf("Failed to write to temp file: %v", err)
		}
		tempFile.Close()

		// Test with the temp file
		provider := NewJsonConfigProvider(tempFile.Name())
		if provider == nil {
			t.Fatal("Expected non-nil provider")
		}

		// Verify the provider works by getting a value
		value := provider.GetString("test")
		if value != "value" {
			t.Errorf("Expected 'value', got '%s'", value)
		}
	})

	// Test with non-existent file
	t.Run("Non-existent file", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected panic for non-existent file")
			}
		}()
		NewJsonConfigProvider("non_existent_file.json")
	})
}

func TestJsonConfigProvider_GetStruct(t *testing.T) {
	// Create a test config file
	tempFile, err := os.CreateTemp("", "config*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	// Write test data with nested structure
	testData := []byte(`{
		"user": {
			"name": "John Doe",
			"age": 30,
			"active": true
		}
	}`)
	if _, err := tempFile.Write(testData); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tempFile.Close()

	provider := NewJsonConfigProvider(tempFile.Name())

	// Define a struct to unmarshal into
	type User struct {
		Name   string `json:"name"`
		Age    int    `json:"age"`
		Active bool   `json:"active"`
	}

	t.Run("Valid struct unmarshaling", func(t *testing.T) {
		var user User
		err := provider.GetStruct("user", &user)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		// Verify the struct fields
		if user.Name != "John Doe" {
			t.Errorf("Expected name to be 'John Doe', got '%s'", user.Name)
		}
		if user.Age != 30 {
			t.Errorf("Expected age to be 30, got %d", user.Age)
		}
		if !user.Active {
			t.Error("Expected active to be true")
		}
	})

	t.Run("Invalid key", func(t *testing.T) {
		var user User
		err := provider.GetStruct("nonexistent", &user)
		if err == nil {
			t.Error("Expected error for nonexistent key")
		}
	})
}

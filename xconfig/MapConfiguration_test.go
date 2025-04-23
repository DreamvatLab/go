package xconfig

import (
	"testing"
)

func TestMapConfiguration_GetMap(t *testing.T) {
	config := MapConfiguration{
		"nested": map[string]interface{}{
			"key": "value",
		},
	}

	t.Run("Valid nested map", func(t *testing.T) {
		nested := config.GetMap("nested")
		if nested == nil {
			t.Fatal("Expected non-nil nested map")
		}
		if nested["key"] != "value" {
			t.Errorf("Expected 'value', got '%v'", nested["key"])
		}
	})

	t.Run("Non-existent nested map", func(t *testing.T) {
		nested := config.GetMap("nonexistent")
		if nested != nil {
			t.Error("Expected nil for nonexistent nested map")
		}
	})
}

func TestMapConfiguration_GetMapSlice(t *testing.T) {
	config := MapConfiguration{
		"items": []interface{}{
			map[string]interface{}{"id": 1, "name": "item1"},
			map[string]interface{}{"id": 2, "name": "item2"},
		},
	}

	t.Run("Valid map slice", func(t *testing.T) {
		items := config.GetMapSlice("items")
		if len(items) != 2 {
			t.Errorf("Expected 2 items, got %d", len(items))
		}
		if items[0]["id"] != 1 {
			t.Errorf("Expected id 1, got %v", items[0]["id"])
		}
		if items[0]["name"] != "item1" {
			t.Errorf("Expected name 'item1', got '%v'", items[0]["name"])
		}
	})

	t.Run("Non-existent map slice", func(t *testing.T) {
		items := config.GetMapSlice("nonexistent")
		if items != nil {
			t.Error("Expected nil for nonexistent map slice")
		}
	})
}

func TestMapConfiguration_GetString(t *testing.T) {
	config := MapConfiguration{
		"string_value": "test",
		"nested": map[string]interface{}{
			"key": "nested_value",
		},
	}

	t.Run("Valid string", func(t *testing.T) {
		value := config.GetString("string_value")
		if value != "test" {
			t.Errorf("Expected 'test', got '%s'", value)
		}
	})

	t.Run("Nested string", func(t *testing.T) {
		value := config.GetString("nested.key")
		if value != "nested_value" {
			t.Errorf("Expected 'nested_value', got '%s'", value)
		}
	})

	t.Run("Non-existent string", func(t *testing.T) {
		value := config.GetString("nonexistent")
		if value != "" {
			t.Errorf("Expected empty string, got '%s'", value)
		}
	})
}

func TestMapConfiguration_GetStringDefault(t *testing.T) {
	config := MapConfiguration{
		"string_value": "test",
	}

	t.Run("Existing string", func(t *testing.T) {
		value := config.GetStringDefault("string_value", "default")
		if value != "test" {
			t.Errorf("Expected 'test', got '%s'", value)
		}
	})

	t.Run("Non-existent string with default", func(t *testing.T) {
		value := config.GetStringDefault("nonexistent", "default")
		if value != "default" {
			t.Errorf("Expected 'default', got '%s'", value)
		}
	})
}

func TestMapConfiguration_GetBool(t *testing.T) {
	config := MapConfiguration{
		"bool_value": true,
		"nested": map[string]interface{}{
			"key": false,
		},
	}

	t.Run("Valid boolean", func(t *testing.T) {
		value := config.GetBool("bool_value")
		if !value {
			t.Error("Expected true")
		}
	})

	t.Run("Nested boolean", func(t *testing.T) {
		value := config.GetBool("nested.key")
		if value {
			t.Error("Expected false")
		}
	})

	t.Run("Non-existent boolean", func(t *testing.T) {
		value := config.GetBool("nonexistent")
		if value {
			t.Error("Expected false for nonexistent boolean")
		}
	})
}

func TestMapConfiguration_GetFloat64(t *testing.T) {
	config := MapConfiguration{
		"float_value": 3.14,
		"nested": map[string]interface{}{
			"key": 42.0,
		},
	}

	t.Run("Valid float", func(t *testing.T) {
		value := config.GetFloat64("float_value")
		if value != 3.14 {
			t.Errorf("Expected 3.14, got %f", value)
		}
	})

	t.Run("Nested float", func(t *testing.T) {
		value := config.GetFloat64("nested.key")
		if value != 42.0 {
			t.Errorf("Expected 42.0, got %f", value)
		}
	})

	t.Run("Non-existent float", func(t *testing.T) {
		value := config.GetFloat64("nonexistent")
		if value != 0 {
			t.Errorf("Expected 0, got %f", value)
		}
	})
}

func TestMapConfiguration_GetInt(t *testing.T) {
	config := MapConfiguration{
		"int_value": 42,
		"nested": map[string]interface{}{
			"key": 100,
		},
	}

	t.Run("Valid integer", func(t *testing.T) {
		value := config.GetInt("int_value")
		if value != 42 {
			t.Errorf("Expected 42, got %d", value)
		}
	})

	t.Run("Nested integer", func(t *testing.T) {
		value := config.GetInt("nested.key")
		if value != 100 {
			t.Errorf("Expected 100, got %d", value)
		}
	})

	t.Run("Non-existent integer", func(t *testing.T) {
		value := config.GetInt("nonexistent")
		if value != 0 {
			t.Errorf("Expected 0, got %d", value)
		}
	})
}

func TestMapConfiguration_GetIntDefault(t *testing.T) {
	config := MapConfiguration{
		"int_value": 42,
	}

	t.Run("Existing integer", func(t *testing.T) {
		value := config.GetIntDefault("int_value", 0)
		if value != 42 {
			t.Errorf("Expected 42, got %d", value)
		}
	})

	t.Run("Non-existent integer with default", func(t *testing.T) {
		value := config.GetIntDefault("nonexistent", 100)
		if value != 100 {
			t.Errorf("Expected 100, got %d", value)
		}
	})
}

func TestMapConfiguration_GetStringSlice(t *testing.T) {
	config := MapConfiguration{
		"string_slice": []interface{}{"item1", "item2", "item3"},
		"nested": map[string]interface{}{
			"key": []interface{}{"nested1", "nested2"},
		},
	}

	t.Run("Valid string slice", func(t *testing.T) {
		slice := config.GetStringSlice("string_slice")
		if len(slice) != 3 {
			t.Errorf("Expected 3 items, got %d", len(slice))
		}
		if slice[0] != "item1" {
			t.Errorf("Expected 'item1', got '%s'", slice[0])
		}
	})

	t.Run("Nested string slice", func(t *testing.T) {
		slice := config.GetStringSlice("nested.key")
		if len(slice) != 2 {
			t.Errorf("Expected 2 items, got %d", len(slice))
		}
		if slice[0] != "nested1" {
			t.Errorf("Expected 'nested1', got '%s'", slice[0])
		}
	})

	t.Run("Non-existent string slice", func(t *testing.T) {
		slice := config.GetStringSlice("nonexistent")
		if len(slice) != 0 {
			t.Error("Expected empty slice for nonexistent string slice")
		}
	})
}

func TestMapConfiguration_GetIntSlice(t *testing.T) {
	config := MapConfiguration{
		"int_slice": []interface{}{1, 2, 3, 4, 5},
		"nested": map[string]interface{}{
			"key": []interface{}{10, 20},
		},
	}

	t.Run("Valid integer slice", func(t *testing.T) {
		slice := config.GetIntSlice("int_slice")
		if len(slice) != 5 {
			t.Errorf("Expected 5 items, got %d", len(slice))
		}
		if slice[0] != 1 {
			t.Errorf("Expected 1, got %d", slice[0])
		}
	})

	t.Run("Nested integer slice", func(t *testing.T) {
		slice := config.GetIntSlice("nested.key")
		if len(slice) != 2 {
			t.Errorf("Expected 2 items, got %d", len(slice))
		}
		if slice[0] != 10 {
			t.Errorf("Expected 10, got %d", slice[0])
		}
	})

	t.Run("Non-existent integer slice", func(t *testing.T) {
		slice := config.GetIntSlice("nonexistent")
		if len(slice) != 0 {
			t.Error("Expected empty slice for nonexistent integer slice")
		}
	})
}

package xconv

import (
	"math"
	"testing"
)

func TestToString(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected string
	}{
		{"string", "hello", "hello"},
		{"int", 123, "123"},
		{"int32", int32(123), "123"},
		{"int64", int64(123), "123"},
		{"float32", float32(123.45), "123.45"},
		{"float64", 123.45, "123.45"},
		{"bool", true, "true"},
		{"unsupported", struct{}{}, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ToString(tt.input)
			if result != tt.expected {
				t.Errorf("ToString(%v) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestToInt(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected int
	}{
		{"int", 123, 123},
		{"int64", int64(123), 123},
		{"int32", int32(123), 123},
		{"float32", float32(123.45), 123},
		{"float64", 123.45, 123},
		{"string", "123", 123},
		{"string_with_comma", "1,234", 1234},
		{"invalid_string", "abc", 0},
		{"unsupported", struct{}{}, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ToInt(tt.input)
			if result != tt.expected {
				t.Errorf("ToInt(%v) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestToIntRound(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected int
	}{
		{"int", 123, 123},
		{"int64", int64(123), 123},
		{"int32", int32(123), 123},
		{"float32_round_up", float32(123.6), 124},
		{"float32_round_down", float32(123.4), 123},
		{"float64_round_up", 123.6, 124},
		{"float64_round_down", 123.4, 123},
		{"string", "123.6", 124},
		{"string_with_comma", "1,234.6", 1235},
		{"invalid_string", "abc", 0},
		{"unsupported", struct{}{}, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ToIntRound(tt.input)
			if result != tt.expected {
				t.Errorf("ToIntRound(%v) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestToInt32(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected int32
	}{
		{"int32", int32(123), 123},
		{"int64", int64(123), 123},
		{"int", 123, 123},
		{"float32", float32(123.45), 123},
		{"float64", 123.45, 123},
		{"string", "123", 123},
		{"string_with_comma", "1,234", 1234},
		{"invalid_string", "abc", 0},
		{"unsupported", struct{}{}, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ToInt32(tt.input)
			if result != tt.expected {
				t.Errorf("ToInt32(%v) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestToInt32Round(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected int32
	}{
		{"int32", int32(123), 123},
		{"int64", int64(123), 123},
		{"int", 123, 123},
		{"float32_round_up", float32(123.6), 124},
		{"float32_round_down", float32(123.4), 123},
		{"float64_round_up", 123.6, 124},
		{"float64_round_down", 123.4, 123},
		{"string", "123.6", 124},
		{"string_with_comma", "1,234.6", 1235},
		{"invalid_string", "abc", 0},
		{"unsupported", struct{}{}, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ToInt32Round(tt.input)
			if result != tt.expected {
				t.Errorf("ToInt32Round(%v) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestToInt64(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected int64
	}{
		{"int64", int64(123), 123},
		{"int32", int32(123), 123},
		{"int", 123, 123},
		{"float32", float32(123.45), 123},
		{"float64", 123.45, 123},
		{"string", "123", 123},
		{"string_with_comma", "1,234", 1234},
		{"invalid_string", "abc", 0},
		{"unsupported", struct{}{}, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ToInt64(tt.input)
			if result != tt.expected {
				t.Errorf("ToInt64(%v) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestToInt64Round(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected int64
	}{
		{"int64", int64(123), 123},
		{"int32", int32(123), 123},
		{"int", 123, 123},
		{"float32_round_up", float32(123.6), 124},
		{"float32_round_down", float32(123.4), 123},
		{"float64_round_up", 123.6, 124},
		{"float64_round_down", 123.4, 123},
		{"string", "123.6", 124},
		{"string_with_comma", "1,234.6", 1235},
		{"invalid_string", "abc", 0},
		{"unsupported", struct{}{}, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ToInt64Round(tt.input)
			if result != tt.expected {
				t.Errorf("ToInt64Round(%v) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestToFloat32(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected float32
	}{
		{"int", 123, 123.0},
		{"int32", int32(123), 123.0},
		{"int64", int64(123), 123.0},
		{"float32", float32(123.45), 123.45},
		{"float64", 123.45, 123.45},
		{"string", "123.45", 123.45},
		{"string_with_comma", "1,234.56", 1234.56},
		{"invalid_string", "abc", 0},
		{"unsupported", struct{}{}, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ToFloat32(tt.input)
			if result != tt.expected {
				t.Errorf("ToFloat32(%v) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestToFloat64(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected float64
	}{
		{"int", 123, 123.0},
		{"int32", int32(123), 123.0},
		{"int64", int64(123), 123.0},
		{"float64", 123.45, 123.45},
		{"float32", float32(123.45), 123.45},
		{"string", "123.45", 123.45},
		{"string_with_comma", "1,234.56", 1234.56},
		{"invalid_string", "abc", 0},
		{"unsupported", struct{}{}, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ToFloat64(tt.input)
			if math.Abs(result-tt.expected) > 1e-5 {
				t.Errorf("ToFloat64(%v) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestToBool(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected bool
	}{
		{"bool_true", true, true},
		{"bool_false", false, false},
		{"string_true", "true", true},
		{"string_false", "false", false},
		{"string_1", "1", true},
		{"string_0", "0", false},
		{"unsupported", struct{}{}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ToBool(tt.input)
			if result != tt.expected {
				t.Errorf("ToBool(%v) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

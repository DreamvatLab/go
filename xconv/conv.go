package xconv

import (
	"math"
	"strconv"
	"strings"
)

// ToString converts any supported type to a string representation.
// Supported types: string, int, int32, int64, float32, float64, bool
// Returns empty string for unsupported types.
func ToString(obj interface{}) string {
	switch v := obj.(type) {
	case string:
		return v
	case int:
		return strconv.Itoa(v)
	case int32:
		return strconv.FormatInt(int64(v), 10)
	case int64:
		return strconv.FormatInt(v, 10)
	case float32:
		return strconv.FormatFloat(float64(v), 'f', 2, 64)
	case float64:
		return strconv.FormatFloat(v, 'f', 2, 64)
	case bool:
		return strconv.FormatBool(v)
	default:
		return ""
	}
}

// ToInt converts any supported type to an int.
// Supported types: int, int64, int32, float32, float64, string
// For string inputs, commas are removed before conversion.
// Returns 0 for unsupported types or conversion errors.
func ToInt(obj interface{}) int {
	switch v := obj.(type) {
	case int:
		return v
	case int64:
		return int(v)
	case int32:
		return int(v)
	case float32:
		return int(v)
	case float64:
		return int(v)
	case string:
		v = strings.ReplaceAll(v, ",", "") // Remove commas
		f, _ := strconv.ParseFloat(v, 64)
		return int(f)
	default:
		return 0
	}
}

// ToIntRound converts any supported type to an int with rounding.
// Similar to ToInt but rounds float values to the nearest integer.
// Supported types: int, int64, int32, float32, float64, string
// Returns 0 for unsupported types or conversion errors.
func ToIntRound(obj interface{}) int {
	switch v := obj.(type) {
	case int:
		return v
	case int64:
		return int(v)
	case int32:
		return int(v)
	case float32:
		return int(math.Round(float64(v)))
	case float64:
		return int(math.Round(v))
	case string:
		v = strings.ReplaceAll(v, ",", "") // Remove commas
		f, _ := strconv.ParseFloat(v, 64)
		return int(math.Round(f))
	default:
		return 0
	}
}

// ToInt32 converts any supported type to an int32.
// Supported types: int32, int64, int, float32, float64, string
// For string inputs, commas are removed before conversion.
// Returns 0 for unsupported types or conversion errors.
func ToInt32(obj interface{}) int32 {
	switch v := obj.(type) {
	case int32:
		return v
	case int64:
		return int32(v)
	case int:
		return int32(v)
	case float32:
		return int32(v)
	case float64:
		return int32(v)
	case string:
		v = strings.ReplaceAll(v, ",", "") // Remove commas
		f, _ := strconv.ParseFloat(v, 64)
		return int32(f)
	default:
		return 0
	}
}

// ToInt32Round converts any supported type to an int32 with rounding.
// Similar to ToInt32 but rounds float values to the nearest integer.
// Supported types: int32, int64, int, float32, float64, string
// Returns 0 for unsupported types or conversion errors.
func ToInt32Round(obj interface{}) int32 {
	switch v := obj.(type) {
	case int32:
		return v
	case int64:
		return int32(v)
	case int:
		return int32(v)
	case float32:
		return int32(math.Round(float64(v)))
	case float64:
		return int32(math.Round(v))
	case string:
		v = strings.ReplaceAll(v, ",", "") // Remove commas
		f, _ := strconv.ParseFloat(v, 64)
		return int32(math.Round(f))
	default:
		return 0
	}
}

// ToInt64 converts any supported type to an int64.
// Supported types: int64, int32, int, float32, float64, string
// For string inputs, commas are removed before conversion.
// Returns 0 for unsupported types or conversion errors.
func ToInt64(obj interface{}) int64 {
	switch v := obj.(type) {
	case int64:
		return v
	case int32:
		return int64(v)
	case int:
		return int64(v)
	case float32:
		return int64(v)
	case float64:
		return int64(v)
	case string:
		v = strings.ReplaceAll(v, ",", "") // Remove commas
		f, _ := strconv.ParseFloat(v, 64)
		return int64(f)
	default:
		return 0
	}
}

// ToInt64Round converts any supported type to an int64 with rounding.
// Similar to ToInt64 but rounds float values to the nearest integer.
// Supported types: int64, int32, int, float32, float64, string
// Returns 0 for unsupported types or conversion errors.
func ToInt64Round(obj interface{}) int64 {
	switch v := obj.(type) {
	case int64:
		return v
	case int32:
		return int64(v)
	case int:
		return int64(v)
	case float32:
		return int64(math.Round(float64(v)))
	case float64:
		return int64(math.Round(v))
	case string:
		v = strings.ReplaceAll(v, ",", "") // Remove commas
		f, _ := strconv.ParseFloat(v, 64)
		return int64(math.Round(f))
	default:
		return 0
	}
}

// ToFloat32 converts any supported type to a float32.
// Supported types: int, int32, int64, float32, float64, string
// For string inputs, commas are removed before conversion.
// Returns 0 for unsupported types or conversion errors.
func ToFloat32(obj interface{}) float32 {
	switch v := obj.(type) {
	case int:
		return float32(v)
	case int32:
		return float32(v)
	case int64:
		return float32(v)
	case float32:
		return v
	case float64:
		return float32(v)
	case string:
		v = strings.ReplaceAll(v, ",", "") // Remove commas
		r, _ := strconv.ParseFloat(v, 32)
		return float32(r)
	default:
		return 0
	}
}

// ToFloat64 converts any supported type to a float64.
// Supported types: int, int32, int64, float64, float32, string
// For string inputs, commas are removed before conversion.
// Returns 0 for unsupported types or conversion errors.
func ToFloat64(obj interface{}) float64 {
	switch v := obj.(type) {
	case int:
		return float64(v)
	case int32:
		return float64(v)
	case int64:
		return float64(v)
	case float64:
		return v
	case float32:
		return float64(v)
	case string:
		v = strings.ReplaceAll(v, ",", "") // Remove commas
		r, _ := strconv.ParseFloat(v, 64)
		return r
	default:
		return 0
	}
}

// ToBool converts any supported type to a boolean.
// Supported types: bool, string
// For string inputs, uses strconv.ParseBool for conversion.
// Returns false for unsupported types or conversion errors.
func ToBool(obj interface{}) bool {
	switch v := obj.(type) {
	case bool:
		return v
	case string:
		r, _ := strconv.ParseBool(v)
		return r
	default:
		return false
	}
}

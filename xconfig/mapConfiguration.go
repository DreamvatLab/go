package xconfig

import (
	"fmt"
	"strings"

	"github.com/Lukiya/go/xconv"
)

// MapConfiguration represents a configuration stored in a map structure
type MapConfiguration map[string]interface{}

// GetMap retrieves a nested map configuration by key
func (x *MapConfiguration) GetMap(key string) (r MapConfiguration) {
	v := getValue(key, *x)
	if v != nil {
		a, ok := v.(interface{})
		if !ok {
			fmt.Printf("convert failed. %t -> interface{}", v)
		} else {
			b := a.(map[string]interface{})
			r = MapConfiguration(b)

			return r
		}
	}
	return nil
}

// GetMapSlice retrieves a slice of map configurations by key
func (x *MapConfiguration) GetMapSlice(key string) []MapConfiguration {
	v := getValue(key, *x)
	if v != nil {
		a, ok := v.([]interface{})
		if !ok {
			fmt.Printf("convert failed. %t -> []interface{}", v)
			return nil
		}
		r := make([]MapConfiguration, 0, len(a))
		for _, b := range a {
			if m, ok := b.(map[string]interface{}); ok {
				r = append(r, MapConfiguration(m))
			}
		}
		return r
	}
	return nil
}

// GetString retrieves a string value by key
func (x *MapConfiguration) GetString(key string) string {
	v := getValue(key, *x)
	if v != nil {
		r, ok := v.(string)
		if !ok {
			fmt.Printf("convert failed. %t -> string", v)
		} else {
			return r
		}
	}
	return ""
}

// GetStringDefault retrieves a string value by key with a default value
func (x *MapConfiguration) GetStringDefault(key, defaultValue string) string {
	r := x.GetString(key)
	if r != "" {
		return r
	}
	return defaultValue
}

// GetBool retrieves a boolean value by key
func (x *MapConfiguration) GetBool(key string) bool {
	v := getValue(key, *x)
	if v != nil {
		r, ok := v.(bool)
		if !ok {
			fmt.Printf("convert failed. %t -> bool", v)
		} else {
			return r
		}
	}
	return false
}

// GetFloat64 retrieves a float64 value by key
func (x *MapConfiguration) GetFloat64(key string) float64 {
	v := getValue(key, *x)
	if v != nil {
		r, ok := v.(float64)
		if !ok {
			fmt.Printf("convert failed. %t -> float64", v)
		} else {
			return r
		}
	}
	return 0
}

// GetInt retrieves an integer value by key
func (x *MapConfiguration) GetInt(key string) int {
	v := getValue(key, *x)
	if v != nil {
		return xconv.ToInt(v)
	}
	return 0
}

// GetIntDefault retrieves an integer value by key with a default value
func (x *MapConfiguration) GetIntDefault(key string, defaultValue int) int {
	r := x.GetInt(key)
	if r != 0 {
		return r
	}
	return defaultValue
}

// GetStringSlice retrieves a slice of strings by key
func (x *MapConfiguration) GetStringSlice(key string) []string {
	v := getValue(key, *x)
	if v != nil {
		slice, ok := v.([]interface{})
		if !ok {
			fmt.Printf("convert failed. %t -> interface slice", v)
		} else {
			var r []string
			for _, e := range slice {
				a, ok := e.(string)
				if ok {
					r = append(r, a)
				}
			}
			return r
		}
	}
	return make([]string, 0)
}

// GetIntSlice retrieves a slice of integers by key
func (x *MapConfiguration) GetIntSlice(key string) []int {
	v := getValue(key, *x)
	if v != nil {
		slice, ok := v.([]interface{})
		if !ok {
			fmt.Printf("convert failed. %t -> interface slice", v)
		} else {
			var r []int
			for _, e := range slice {
				switch val := e.(type) {
				case float64:
					r = append(r, int(val))
				case int:
					r = append(r, val)
				case int64:
					r = append(r, int(val))
				case int32:
					r = append(r, int(val))
				}
			}
			return r
		}
	}
	return make([]int, 0)
}

// getValue retrieves a value from the configuration using a dot-separated key path
func getValue(key string, c MapConfiguration) interface{} {
	keys := strings.Split(key, ".")
	keyCount := len(keys)

	for i := 0; i < keyCount; i++ {
		if i < keyCount-1 {
			c = getNode(keys[i], c)
			if c == nil {
				break
			}
		} else {
			r, ok := c[keys[i]]
			if ok {
				return r
			}
		}
	}

	return nil
}

// getNode retrieves a nested map node by key
func getNode(key string, c map[string]interface{}) map[string]interface{} {
	v, ok := c[key]
	if ok {
		return v.(map[string]interface{})
	}
	return nil
}

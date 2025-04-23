package xconfig

// IConfigProvider defines the interface for configuration providers
type IConfigProvider interface {
	// GetStruct retrieves a configuration section as a struct
	GetStruct(key string, target interface{}) error
	// GetString retrieves a string value from the configuration
	GetString(key string) string
	// GetStringDefault retrieves a string value from the configuration with a default value
	GetStringDefault(key string, defaultValue string) string
	// GetBool retrieves a boolean value from the configuration
	GetBool(key string) bool
	// GetFloat64 retrieves a float64 value from the configuration
	GetFloat64(key string) float64
	// GetInt retrieves an integer value from the configuration
	GetInt(key string) int
	// GetIntDefault retrieves an integer value from the configuration with a default value
	GetIntDefault(key string, defaultValue int) int
	// GetStringSlice retrieves a slice of strings from the configuration
	GetStringSlice(key string) []string
	// GetIntSlice retrieves a slice of integers from the configuration
	GetIntSlice(key string) []int
}

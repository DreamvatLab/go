package xutils

import (
	"fmt"

	"github.com/DreamvatLab/go/xlog"
	"github.com/sony/sonyflake"
)

// Default ID generator instance
var (
	_defaultIDGenerator = NewSonyflakeIDGenerator()
)

// GenerateStringID generates a unique string ID using the default generator
func GenerateStringID() string {
	return _defaultIDGenerator.GenerateStringID()
}

// GenerateUInt64ID generates a unique uint64 ID using the default generator
func GenerateUInt64ID() uint64 {
	return _defaultIDGenerator.GenerateUInt64ID()
}

// IIDGenerator defines the interface for ID generation
type IIDGenerator interface {
	GenerateStringID() string
	GenerateUInt64ID() uint64
}

// SonyflakeIDGenerator implements IIDGenerator using Sonyflake algorithm
type SonyflakeIDGenerator struct {
	generator *sonyflake.Sonyflake
}

// NewSonyflakeIDGenerator creates a new instance of SonyflakeIDGenerator
func NewSonyflakeIDGenerator() IIDGenerator {
	return &SonyflakeIDGenerator{
		generator: sonyflake.NewSonyflake(sonyflake.Settings{}),
	}
}

// GenerateStringID converts the uint64 ID to a hexadecimal string representation
func (x *SonyflakeIDGenerator) GenerateStringID() string {
	return fmt.Sprintf("%x", x.GenerateUInt64ID())
}

// GenerateUInt64ID generates a unique uint64 ID using the Sonyflake algorithm
func (x *SonyflakeIDGenerator) GenerateUInt64ID() uint64 {
	id, err := x.generator.NextID()
	if err != nil {
		xlog.Errorf("flake.NextID() failed with %s\n", err)
	}
	return id
}

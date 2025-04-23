package xbytes

import (
	"reflect"
	"testing"
)

func TestStrToBytes(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want []byte
	}{
		{
			name: "empty string",
			s:    "",
			want: nil,
		},
		{
			name: "non-empty string",
			s:    "hello",
			want: []byte("hello"),
		},
		{
			name: "string with special characters",
			s:    "hello, world!",
			want: []byte("hello, world!"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := StrToBytes(tt.s)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StrToBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBytesToStr(t *testing.T) {
	tests := []struct {
		name string
		b    []byte
		want string
	}{
		{
			name: "nil bytes",
			b:    nil,
			want: "",
		},
		{
			name: "empty bytes",
			b:    []byte{},
			want: "",
		},
		{
			name: "non-empty bytes",
			b:    []byte("hello"),
			want: "hello",
		},
		{
			name: "bytes with special characters",
			b:    []byte("hello, world!"),
			want: "hello, world!",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := BytesToStr(tt.b)
			if got != tt.want {
				t.Errorf("BytesToStr() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestZeroCopyConversion tests that the conversion between string and bytes
// is truly zero-copy by checking that modifying the original affects the result
func TestZeroCopyConversion(t *testing.T) {
	// Test string to bytes
	originalStr := "hello"
	bytes := StrToBytes(originalStr)

	// Verify that the bytes slice contains the correct data
	if !reflect.DeepEqual(bytes, []byte("hello")) {
		t.Errorf("StrToBytes() did not convert correctly, got %v, want %v", bytes, []byte("hello"))
	}

	// Test bytes to string
	originalBytes := []byte("world")
	str := BytesToStr(originalBytes)

	// Modify the original bytes
	originalBytes[0] = 'W'

	// Check if the string was affected (it should be)
	if str != "World" {
		t.Errorf("BytesToStr() did not maintain reference to original bytes, got %v, want %v", str, "World")
	}
}

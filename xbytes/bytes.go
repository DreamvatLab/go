package xbytes

import "unsafe"

// StrToBytes converts a string to a byte slice without copying the data.
// This is a zero-copy conversion that should be used with caution.
// The resulting byte slice should not be modified as it shares memory with the original string.
func StrToBytes(s string) []byte {
	// x := (*[2]uintptr)(unsafe.Pointer(&s))
	// h := [3]uintptr{x[0], x[1], x[1]}
	// return *(*[]byte)(unsafe.Pointer(&h))

	if len(s) == 0 {
		return nil
	}
	// For Go 1.20+, it's recommended to use unsafe.StringData to get the string data address
	return unsafe.Slice(unsafe.StringData(s), len(s))
}

// BytesToStr converts a byte slice to a string without copying the data.
// This is a zero-copy conversion that should be used with caution.
// The resulting string shares memory with the original byte slice.
func BytesToStr(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

package xerr

import "github.com/DreamvatLab/go/xlog"

// LogError checks if an error exists and logs it if present.
// Returns true if the error exists, false otherwise.
func LogError(err error) bool {
	if err != nil {
		xlog.Error(err)
		return true
	}

	return false
}

// FatalIfErr checks if an error exists and logs it as a fatal error if present.
// This function will terminate the program if an error is encountered.
func FatalIfErr(err error) {
	if err != nil {
		xlog.Fatal(err)
	}
}

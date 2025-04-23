package xerr

import "github.com/Lukiya/go/xlog"

// HasError checks if an error exists and logs it if present.
// Returns true if the error exists, false otherwise.
func HasError(err error) bool {
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

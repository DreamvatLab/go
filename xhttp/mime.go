package xhttp

import "regexp"

var (
	// Standard base64 regex that validates:
	// 1. Only allows valid base64 characters (A-Z, a-z, 0-9, +, /)
	// 2. String length must be multiple of 4
	// 3. Proper padding (=) at the end if needed
	_base64Regex = regexp.MustCompile(`^(?:[A-Za-z0-9+/]{4})*(?:[A-Za-z0-9+/]{2}==|[A-Za-z0-9+/]{3}=)?$`)
)

func IsBase64String(str string) bool {
	return _base64Regex.MatchString(str)
}

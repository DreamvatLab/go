package xhttp

import "testing"

func TestIsBase64String(t *testing.T) {
	tests := []struct {
		name string
		str  string
		want bool
	}{
		{
			name: "valid base64 - simple string",
			str:  "SGVsbG8gV29ybGQ=", // "Hello World"
			want: true,
		},
		{
			name: "valid base64 - empty string",
			str:  "",
			want: true,
		},
		{
			name: "valid base64 - with padding ==",
			str:  "YWJj", // "abc"
			want: true,
		},
		{
			name: "valid base64 - with padding =",
			str:  "YWJjZGU=", // "abcde"
			want: true,
		},
		{
			name: "valid base64 - special characters",
			str:  "SGVsbG8gV29ybGQhIQ==", // "Hello World!!"
			want: true,
		},
		{
			name: "invalid base64 - wrong padding",
			str:  "SGVsbG8gV29ybGQ==", // "Hello World" with wrong padding (should be single =)
			want: false,
		},
		{
			name: "invalid base64 - invalid characters",
			str:  "SGVsbG8gV29ybGQ$",
			want: false,
		},
		{
			name: "invalid base64 - wrong length",
			str:  "YWJjZ", // length not multiple of 4
			want: false,
		},
		{
			name: "invalid base64 - multiple padding",
			str:  "YWJj====",
			want: false,
		},
		{
			name: "invalid base64 - padding in middle",
			str:  "YWJj=ZGU=",
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsBase64String(tt.str); got != tt.want {
				t.Errorf("IsBase64String() = %v, want %v", got, tt.want)
			}
		})
	}
}

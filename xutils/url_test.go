package xutils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJointURL(t *testing.T) {
	tests := []struct {
		name     string
		basePath string
		paths    []string
		want     string
		wantErr  bool
	}{
		{
			name:     "simple path join",
			basePath: "http://example.com",
			paths:    []string{"path", "to", "resource"},
			want:     "http://example.com/path/to/resource",
			wantErr:  false,
		},
		{
			name:     "with trailing slash",
			basePath: "http://example.com/",
			paths:    []string{"path", "to", "resource"},
			want:     "http://example.com/path/to/resource",
			wantErr:  false,
		},
		{
			name:     "with double slashes",
			basePath: "http://example.com//",
			paths:    []string{"//path//", "//to//", "//resource//"},
			want:     "http://example.com/path/to/resource",
			wantErr:  false,
		},
		{
			name:     "invalid base URL",
			basePath: "://invalid",
			paths:    []string{"path"},
			want:     "",
			wantErr:  true,
		},
		{
			name:     "empty paths",
			basePath: "http://example.com",
			paths:    []string{},
			want:     "http://example.com",
			wantErr:  false,
		},
		{
			name:     "file protocol",
			basePath: "file:///home/user",
			paths:    []string{"documents", "file.txt"},
			want:     "file:///home/user/documents/file.txt",
			wantErr:  false,
		},
		{
			name:     "ftp protocol",
			basePath: "ftp://ftp.example.com",
			paths:    []string{"downloads", "software"},
			want:     "ftp://ftp.example.com/downloads/software",
			wantErr:  false,
		},
		{
			name:     "relative path",
			basePath: "/home/user",
			paths:    []string{"documents", "file.txt"},
			want:     "/home/user/documents/file.txt",
			wantErr:  false,
		},
		{
			name:     "relative path with out prefix /",
			basePath: "home/user",
			paths:    []string{"documents", "file.txt"},
			want:     "home/user/documents/file.txt",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := JointURL(tt.basePath, tt.paths...)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got.String())
		})
	}
}

func TestJointURLString(t *testing.T) {
	tests := []struct {
		name     string
		basePath string
		paths    []string
		want     string
		wantErr  bool
	}{
		{
			name:     "simple path join",
			basePath: "http://example.com",
			paths:    []string{"path", "to", "resource"},
			want:     "http://example.com/path/to/resource",
			wantErr:  false,
		},
		{
			name:     "with query parameters",
			basePath: "http://example.com?param=value",
			paths:    []string{"path"},
			want:     "http://example.com/path?param=value",
			wantErr:  false,
		},
		{
			name:     "invalid base URL",
			basePath: "://invalid",
			paths:    []string{"path"},
			want:     "",
			wantErr:  true,
		},
		{
			name:     "file protocol with query",
			basePath: "file:///home/user?type=text",
			paths:    []string{"documents", "file.txt"},
			want:     "file:///home/user/documents/file.txt?type=text",
			wantErr:  false,
		},
		{
			name:     "sftp protocol",
			basePath: "sftp://user@example.com",
			paths:    []string{"backup", "data"},
			want:     "sftp://user@example.com/backup/data",
			wantErr:  false,
		},
		{
			name:     "relative path with query",
			basePath: "/home/user?filter=all",
			paths:    []string{"documents", "file.txt"},
			want:     "/home/user/documents/file.txt?filter=all",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := JointURLString(tt.basePath, tt.paths...)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

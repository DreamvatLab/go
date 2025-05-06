package xslice

import (
	"testing"
)

func TestFindItem(t *testing.T) {
	tests := []struct {
		name     string
		slice    []int
		fnFilter func(int) bool
		wantIdx  int
		wantVal  int
	}{
		{
			name:     "empty slice",
			slice:    []int{},
			fnFilter: func(x int) bool { return x > 0 },
			wantIdx:  -1,
			wantVal:  0,
		},
		{
			name:     "found in small slice",
			slice:    []int{1, 2, 3, 4, 5},
			fnFilter: func(x int) bool { return x == 3 },
			wantIdx:  2,
			wantVal:  3,
		},
		{
			name:     "not found in small slice",
			slice:    []int{1, 2, 3, 4, 5},
			fnFilter: func(x int) bool { return x > 10 },
			wantIdx:  -1,
			wantVal:  0,
		},
		{
			name:     "found in large slice",
			slice:    make([]int, 20000),
			fnFilter: func(x int) bool { return x == 100 },
			wantIdx:  100,
			wantVal:  100,
		},
	}

	// Initialize large slice with values
	for i := range tests[3].slice {
		tests[3].slice[i] = i
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotIdx, gotVal := FindItem(tt.slice, tt.fnFilter)
			if gotIdx != tt.wantIdx {
				t.Errorf("FindItem() gotIdx = %v, want %v", gotIdx, tt.wantIdx)
			}
			if gotVal != tt.wantVal {
				t.Errorf("FindItem() gotVal = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestFindItemLinear(t *testing.T) {
	tests := []struct {
		name     string
		slice    []string
		fnFilter func(string) bool
		wantIdx  int
		wantVal  string
	}{
		{
			name:     "empty slice",
			slice:    []string{},
			fnFilter: func(s string) bool { return len(s) > 0 },
			wantIdx:  -1,
			wantVal:  "",
		},
		{
			name:     "found in slice",
			slice:    []string{"a", "b", "c", "d"},
			fnFilter: func(s string) bool { return s == "c" },
			wantIdx:  2,
			wantVal:  "c",
		},
		{
			name:     "not found in slice",
			slice:    []string{"a", "b", "c", "d"},
			fnFilter: func(s string) bool { return s == "z" },
			wantIdx:  -1,
			wantVal:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotIdx, gotVal := FindItemLinear(tt.slice, tt.fnFilter)
			if gotIdx != tt.wantIdx {
				t.Errorf("FindItemLinear() gotIdx = %v, want %v", gotIdx, tt.wantIdx)
			}
			if gotVal != tt.wantVal {
				t.Errorf("FindItemLinear() gotVal = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestFindItemConcurrent(t *testing.T) {
	tests := []struct {
		name     string
		slice    []int
		fnFilter func(int) bool
		wantIdx  int
		wantVal  int
	}{
		{
			name:     "empty slice",
			slice:    []int{},
			fnFilter: func(x int) bool { return x > 0 },
			wantIdx:  -1,
			wantVal:  0,
		},
		{
			name:     "found in slice",
			slice:    make([]int, 20000),
			fnFilter: func(x int) bool { return x == 100 },
			wantIdx:  100,
			wantVal:  100,
		},
		{
			name:     "not found in slice",
			slice:    make([]int, 20000),
			fnFilter: func(x int) bool { return x == 20000 },
			wantIdx:  -1,
			wantVal:  0,
		},
	}

	// Initialize slices with values
	for i := range tests[1].slice {
		tests[1].slice[i] = i
	}
	for i := range tests[2].slice {
		tests[2].slice[i] = i
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotIdx, gotVal := FindItemConcurrent(tt.slice, tt.fnFilter)
			if gotIdx != tt.wantIdx {
				t.Errorf("FindItemConcurrent() gotIdx = %v, want %v", gotIdx, tt.wantIdx)
			}
			if gotVal != tt.wantVal {
				t.Errorf("FindItemConcurrent() gotVal = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

package xslice

import (
	"testing"
)

func TestHasStr(t *testing.T) {
	tests := []struct {
		name   string
		source []string
		in     string
		want   bool
	}{
		{
			name:   "empty slice",
			source: []string{},
			in:     "test",
			want:   false,
		},
		{
			name:   "found in slice",
			source: []string{"a", "b", "c"},
			in:     "b",
			want:   true,
		},
		{
			name:   "not found in slice",
			source: []string{"a", "b", "c"},
			in:     "d",
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HasStr(tt.source, tt.in); got != tt.want {
				t.Errorf("HasStr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHasAnyStr(t *testing.T) {
	tests := []struct {
		name   string
		source []string
		in     []string
		want   bool
	}{
		{
			name:   "empty source",
			source: []string{},
			in:     []string{"a", "b"},
			want:   false,
		},
		{
			name:   "empty in",
			source: []string{"a", "b", "c"},
			in:     []string{},
			want:   false,
		},
		{
			name:   "found any",
			source: []string{"a", "b", "c"},
			in:     []string{"d", "b", "e"},
			want:   true,
		},
		{
			name:   "not found any",
			source: []string{"a", "b", "c"},
			in:     []string{"d", "e", "f"},
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HasAnyStr(tt.source, tt.in); got != tt.want {
				t.Errorf("HasAnyStr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHasAllStr(t *testing.T) {
	tests := []struct {
		name   string
		source []string
		in     []string
		want   bool
	}{
		{
			name:   "empty source",
			source: []string{},
			in:     []string{"a", "b"},
			want:   false,
		},
		{
			name:   "empty in",
			source: []string{"a", "b", "c"},
			in:     []string{},
			want:   true,
		},
		{
			name:   "has all",
			source: []string{"a", "b", "c"},
			in:     []string{"a", "b"},
			want:   true,
		},
		{
			name:   "missing some",
			source: []string{"a", "b", "c"},
			in:     []string{"a", "d"},
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HasAllStr(tt.source, tt.in); got != tt.want {
				t.Errorf("HasAllStr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHasInt(t *testing.T) {
	tests := []struct {
		name   string
		source []int
		in     int
		want   bool
	}{
		{
			name:   "empty slice",
			source: []int{},
			in:     1,
			want:   false,
		},
		{
			name:   "found in slice",
			source: []int{1, 2, 3},
			in:     2,
			want:   true,
		},
		{
			name:   "not found in slice",
			source: []int{1, 2, 3},
			in:     4,
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HasInt(tt.source, tt.in); got != tt.want {
				t.Errorf("HasInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHasAnyInt(t *testing.T) {
	tests := []struct {
		name   string
		source []int
		in     []int
		want   bool
	}{
		{
			name:   "empty source",
			source: []int{},
			in:     []int{1, 2},
			want:   false,
		},
		{
			name:   "empty in",
			source: []int{1, 2, 3},
			in:     []int{},
			want:   false,
		},
		{
			name:   "found any",
			source: []int{1, 2, 3},
			in:     []int{4, 2, 5},
			want:   true,
		},
		{
			name:   "not found any",
			source: []int{1, 2, 3},
			in:     []int{4, 5, 6},
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HasAnyInt(tt.source, tt.in); got != tt.want {
				t.Errorf("HasAnyInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHasAllInt(t *testing.T) {
	tests := []struct {
		name   string
		source []int
		in     []int
		want   bool
	}{
		{
			name:   "empty source",
			source: []int{},
			in:     []int{1, 2},
			want:   false,
		},
		{
			name:   "empty in",
			source: []int{1, 2, 3},
			in:     []int{},
			want:   true,
		},
		{
			name:   "has all",
			source: []int{1, 2, 3},
			in:     []int{1, 2},
			want:   true,
		},
		{
			name:   "missing some",
			source: []int{1, 2, 3},
			in:     []int{1, 4},
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HasAllInt(tt.source, tt.in); got != tt.want {
				t.Errorf("HasAllInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHasFloat64(t *testing.T) {
	tests := []struct {
		name   string
		source []float64
		in     float64
		want   bool
	}{
		{
			name:   "empty slice",
			source: []float64{},
			in:     1.0,
			want:   false,
		},
		{
			name:   "exact match",
			source: []float64{1.0, 2.0, 3.0},
			in:     2.0,
			want:   true,
		},
		{
			name:   "within epsilon",
			source: []float64{1.0, 2.0, 3.0},
			in:     2.0 + 1e-10,
			want:   true,
		},
		{
			name:   "not found",
			source: []float64{1.0, 2.0, 3.0},
			in:     4.0,
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HasFloat64(tt.source, tt.in); got != tt.want {
				t.Errorf("HasFloat64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHasAnyFloat64(t *testing.T) {
	tests := []struct {
		name   string
		source []float64
		in     []float64
		want   bool
	}{
		{
			name:   "empty source",
			source: []float64{},
			in:     []float64{1.0, 2.0},
			want:   false,
		},
		{
			name:   "empty in",
			source: []float64{1.0, 2.0, 3.0},
			in:     []float64{},
			want:   false,
		},
		{
			name:   "found any",
			source: []float64{1.0, 2.0, 3.0},
			in:     []float64{4.0, 2.0, 5.0},
			want:   true,
		},
		{
			name:   "found any within epsilon",
			source: []float64{1.0, 2.0, 3.0},
			in:     []float64{4.0, 2.0 + 1e-10, 5.0},
			want:   true,
		},
		{
			name:   "not found any",
			source: []float64{1.0, 2.0, 3.0},
			in:     []float64{4.0, 5.0, 6.0},
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HasAnyFloat64(tt.source, tt.in); got != tt.want {
				t.Errorf("HasAnyFloat64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHasAllFloat64(t *testing.T) {
	tests := []struct {
		name   string
		source []float64
		in     []float64
		want   bool
	}{
		{
			name:   "empty source",
			source: []float64{},
			in:     []float64{1.0, 2.0},
			want:   false,
		},
		{
			name:   "empty in",
			source: []float64{1.0, 2.0, 3.0},
			in:     []float64{},
			want:   true,
		},
		{
			name:   "has all",
			source: []float64{1.0, 2.0, 3.0},
			in:     []float64{1.0, 2.0},
			want:   true,
		},
		{
			name:   "has all within epsilon",
			source: []float64{1.0, 2.0, 3.0},
			in:     []float64{1.0 + 1e-10, 2.0},
			want:   true,
		},
		{
			name:   "missing some",
			source: []float64{1.0, 2.0, 3.0},
			in:     []float64{1.0, 4.0},
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HasAllFloat64(tt.source, tt.in); got != tt.want {
				t.Errorf("HasAllFloat64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHasFloat32(t *testing.T) {
	tests := []struct {
		name   string
		source []float32
		in     float32
		want   bool
	}{
		{
			name:   "empty slice",
			source: []float32{},
			in:     1.0,
			want:   false,
		},
		{
			name:   "exact match",
			source: []float32{1.0, 2.0, 3.0},
			in:     2.0,
			want:   true,
		},
		{
			name:   "within epsilon",
			source: []float32{1.0, 2.0, 3.0},
			in:     2.0 + 1e-7,
			want:   true,
		},
		{
			name:   "not found",
			source: []float32{1.0, 2.0, 3.0},
			in:     4.0,
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HasFloat32(tt.source, tt.in); got != tt.want {
				t.Errorf("HasFloat32() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHasAnyFloat32(t *testing.T) {
	tests := []struct {
		name   string
		source []float32
		in     []float32
		want   bool
	}{
		{
			name:   "empty source",
			source: []float32{},
			in:     []float32{1.0, 2.0},
			want:   false,
		},
		{
			name:   "empty in",
			source: []float32{1.0, 2.0, 3.0},
			in:     []float32{},
			want:   false,
		},
		{
			name:   "found any",
			source: []float32{1.0, 2.0, 3.0},
			in:     []float32{4.0, 2.0, 5.0},
			want:   true,
		},
		{
			name:   "found any within epsilon",
			source: []float32{1.0, 2.0, 3.0},
			in:     []float32{4.0, 2.0 + 1e-7, 5.0},
			want:   true,
		},
		{
			name:   "not found any",
			source: []float32{1.0, 2.0, 3.0},
			in:     []float32{4.0, 5.0, 6.0},
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HasAnyFloat32(tt.source, tt.in); got != tt.want {
				t.Errorf("HasAnyFloat32() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHasAllFloat32(t *testing.T) {
	tests := []struct {
		name   string
		source []float32
		in     []float32
		want   bool
	}{
		{
			name:   "empty source",
			source: []float32{},
			in:     []float32{1.0, 2.0},
			want:   false,
		},
		{
			name:   "empty in",
			source: []float32{1.0, 2.0, 3.0},
			in:     []float32{},
			want:   true,
		},
		{
			name:   "has all",
			source: []float32{1.0, 2.0, 3.0},
			in:     []float32{1.0, 2.0},
			want:   true,
		},
		{
			name:   "has all within epsilon",
			source: []float32{1.0, 2.0, 3.0},
			in:     []float32{1.0 + 1e-7, 2.0},
			want:   true,
		},
		{
			name:   "missing some",
			source: []float32{1.0, 2.0, 3.0},
			in:     []float32{1.0, 4.0},
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HasAllFloat32(tt.source, tt.in); got != tt.want {
				t.Errorf("HasAllFloat32() = %v, want %v", got, tt.want)
			}
		})
	}
}

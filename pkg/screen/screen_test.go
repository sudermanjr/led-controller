package screen

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_nearestMultiple(t *testing.T) {
	tests := []struct {
		name   string
		input  int
		target int
		up     bool
		want   int
	}{
		{"7 up", 7, 8, true, 8},
		{"15 up", 15, 8, true, 16},
		{"9 down", 9, 8, false, 8},
		{"9 up", 9, 8, true, 16},
		{"up doesn't matter 1", 64, 8, true, 64},
		{"up doesn't matter 2", 64, 8, false, 64},
		{"3 6 up", 3, 6, true, 6},
		{"8 9 up", 8, 9, true, 9},
		{"50 8 down", 50, 8, false, 48},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := nearestMultiple(tt.input, tt.target, tt.up)
			assert.Equal(t, tt.want, got)
		})
	}
}

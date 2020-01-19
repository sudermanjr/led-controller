package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_scaleHomekitBrightness(t *testing.T) {
	tests := []struct {
		name  string
		value int
		want  int
	}{
		{name: "top", value: 100, want: 200},
		{name: "bottom", value: 0, want: 30},
		{name: "middle", value: 50, want: 115},
		{name: "other", value: 18, want: 60},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := scaleHomekitBrightness(tt.value)
			assert.Equal(t, tt.want, got)
		})
	}
}

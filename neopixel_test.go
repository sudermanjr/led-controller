package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_brightnessBounds(t *testing.T) {
	tests := []struct {
		name  string
		value int
		want  int
	}{
		{name: "good", value: 100, want: 100},
		{name: "high", value: 251, want: 200},
		{name: "low", value: 10, want: 30},
	}
	for _, tt := range tests {
		led := &ledArray{}
		t.Run(tt.name, func(t *testing.T) {
			led.brightness = tt.value
			led.checkBrightness()
			assert.Equal(t, tt.want, led.brightness)
		})
	}
}

func Test_stepRamp(t *testing.T) {
	tests := []struct {
		name     string
		start    float64
		stop     float64
		duration float64
		want     []int
	}{
		{
			name:     "basic up",
			start:    0,
			stop:     150,
			duration: 10,
			want:     []int{0, 15, 30, 45, 60, 75, 90, 105, 120, 135},
		},
		{
			name:     "basic down",
			start:    150,
			stop:     0,
			duration: 10,
			want:     []int{150, 135, 120, 105, 90, 75, 60, 45, 30, 15},
		},
		{
			name:     "up a little in a long time",
			start:    30,
			stop:     35,
			duration: 20,
			want:     []int{30, 30, 30, 30, 31, 31, 31, 31, 32, 32, 32, 32, 33, 33, 33, 33, 34, 34, 34, 34},
		},
		{
			name:     "down a little in a long time",
			start:    35,
			stop:     30,
			duration: 20,
			want:     []int{35, 34, 34, 34, 34, 33, 33, 33, 33, 32, 32, 32, 32, 31, 31, 31, 31, 30, 30, 30},
		},
		{
			name:     "down a lot in a short time",
			start:    200,
			stop:     5,
			duration: 5,
			want:     []int{200, 161, 122, 83, 44},
		},
	}
	for _, tt := range tests {
		got := stepRamp(tt.start, tt.stop, tt.duration)
		assert.EqualValues(t, tt.want, got)
	}
}

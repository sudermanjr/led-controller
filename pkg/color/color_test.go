package color

import (
	"fmt"
	"os"
	"testing"

	"github.com/lucasb-eyer/go-colorful"
	"github.com/stretchr/testify/assert"
	"github.com/sudermanjr/led-controller/pkg/utils"
)

var testGradient1 = GradientTable{
	{HexToColor("#9e0142"), 0.0},
	{HexToColor("#d53e4f"), 0.1},
	{HexToColor("#f46d43"), 0.2},
	{HexToColor("#fdae61"), 0.3},
	{HexToColor("#fee090"), 0.4},
	{HexToColor("#ffffbf"), 0.5},
	{HexToColor("#e6f598"), 0.6},
	{HexToColor("#abdda4"), 0.7},
	{HexToColor("#66c2a5"), 0.8},
	{HexToColor("#3288bd"), 0.9},
	{HexToColor("#5e4fa2"), 1.0},
}

var testGradient2 = GradientTable{
	{HexToColor("#4e3cec"), 0.0},
	{HexToColor("#5b3ee8"), 0.1},
	{HexToColor("#6941e5"), 0.2},
	{HexToColor("#7643e1"), 0.3},
	{HexToColor("#8346dd"), 0.4},
	{HexToColor("#9148da"), 0.5},
	{HexToColor("#9e4bd6"), 0.6},
	{HexToColor("#ab4dd2"), 0.7},
	{HexToColor("#b94fcf"), 0.8},
	{HexToColor("#c652cb"), 0.9},
	{HexToColor("#e157c4"), 1.0},
}

func TestParseHex(t *testing.T) {
	tests := []struct {
		name string
		hex  string
		want colorful.Color
	}{
		{"blue", "#0000ff", colorful.Color{R: 0, G: 0, B: 1}},
		{"yellow", "#ffff00", colorful.Color{R: 1, G: 1, B: 0}},
		{"red", "#ff0000", colorful.Color{R: 1, G: 0, B: 0}},
		{"black", "#000000", colorful.Color{R: 0, G: 0, B: 0}},
		{"green", "#00ff00", colorful.Color{R: 0, G: 1, B: 0}},
		{"white", "#ffffff", colorful.Color{R: 1, G: 1, B: 1}},
		{"notacolor", "ff", colorful.Color{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := HexToColor(tt.hex)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestColorToUint32(t *testing.T) {
	tests := []struct {
		name  string
		color colorful.Color
		want  uint32
	}{
		{"white", colorful.Color{R: 1, G: 1, B: 1}, uint32(16777215)},
		{"black", colorful.Color{R: 0, G: 0, B: 0}, uint32(0)},
		{"green", colorful.Color{R: 0, G: 1, B: 0}, uint32(65280)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ToUint32(tt.color)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestGradientTable_GetInterpolatedColor(t *testing.T) {
	tests := []struct {
		name  string
		value float64
		want  colorful.Color
	}{
		{"one", 1.0, colorful.Color{R: 0.3686274518393889, G: 0.30980394385954535, B: 0.635294122225692}},
		{"two", 1.1, colorful.Color{R: 0.3686274509803922, G: 0.30980392156862746, B: 0.6352941176470588}},
		{"three", 1.2, colorful.Color{R: 0.3686274509803922, G: 0.30980392156862746, B: 0.6352941176470588}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := testGradient1.GetInterpolatedColor(tt.value)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestGradientPNG(t *testing.T) {
	tests := []struct {
		name     string
		gradient GradientTable
		h        int
		w        int
		testFile string
	}{
		{"1024x40 1", testGradient1, 1024, 40, "1024x40-gradient-1.png"},
		{"1024x1024 1", testGradient1, 1024, 1024, "1024x1024-gradient-1.png"},
		{"2048x40 2", testGradient2, 2048, 40, "2048x40-gradient-2.png"},
	}
	for _, tt := range tests {
		os.Remove("gradient.png")
		t.Run(tt.name, func(t *testing.T) {
			GradientPNG(tt.gradient, tt.h, tt.w)
			assert.FileExistsf(t, "gradient.png", "gradient.png should exist")
			match := utils.DeepCompareFiles("gradient.png", "testdata/"+tt.testFile)
			assert.Truef(t, match, "the files must match")
		})
		os.Remove("gradient.png")
	}
}

func TestGradientColorList(t *testing.T) {
	tests := []struct {
		name     string
		gradient GradientTable
		length   int
		want     []colorful.Color
	}{
		{"one", testGradient1, 1, []colorful.Color{{R: 0.6196077933795217, G: 0.003922138953572327, B: 0.2588235191354816}}},
		{
			"two",
			testGradient2,
			5,
			[]colorful.Color{
				{R: 0.3058824116295, G: 0.23529411042695136, B: 0.9254902064503917},
				{R: 0.4117647268595018, G: 0.2549019781745598, B: 0.8980392249706234},
				{R: 0.5137254867393809, G: 0.2745098482413161, B: 0.866666674535384},
				{R: 0.6196078211495166, G: 0.29411772176478884, B: 0.8392156924085546},
				{R: 0.72549015896929, G: 0.3098040307126474, B: 0.8117647099046337},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GradientColorList(tt.gradient, tt.length)
			fmt.Println(got)
			assert.Equal(t, tt.want, got)
		})
	}
}

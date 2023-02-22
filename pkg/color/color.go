package color

import (
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"os"
	"strconv"
	"strings"

	"github.com/lucasb-eyer/go-colorful"
)

// ColorMap is a map of named colors to hex values
var ColorMap = map[string]string{
	"warmwhite": "#ffe9c2",
	"blue":      "#0000ff",
	"green":     "#00ff00",
	"yellow":    "#ffaf33",
	"purple":    "#af33ff",
	"red":       "#ff0000",
	"teal":      "#33ffd1",
	"pink":      "#ff08c7",
	"white":     "#ffffff",
	"black":     "#000000", // This basically equates to off.
}

// GradientTable contains the "keypoints" of the colorgradient you want to generate.
// The position of each keypoint has to live in the range [0,1]
type GradientTable []struct {
	Col colorful.Color
	Pos float64
}

// GetInterpolatedColor is the meat of the gradient computation. It returns
// an HCL-blend between the two colors around `t`.
// Note: It relies heavily on the fact that the gradient keypoints are sorted.
func (gt GradientTable) GetInterpolatedColor(t float64) colorful.Color {
	for i := 0; i < len(gt)-1; i++ {
		c1 := gt[i]
		c2 := gt[i+1]
		if c1.Pos <= t && t <= c2.Pos {
			// We are in between c1 and c2. Go blend them!
			t := (t - c1.Pos) / (c2.Pos - c1.Pos)
			return c1.Col.BlendHcl(c2.Col, t).Clamped()
		}
	}

	// Nothing found? Means we're at (or past) the last gradient keypoint.
	return gt[len(gt)-1].Col
}

// HexToColor converts a hex string to a Color
func HexToColor(s string) colorful.Color {
	c, err := colorful.Hex(s)
	if err != nil {
		fmt.Println("error converting hex string to color", "error", err, "hex", s) // TODO - logging better
	}
	return c
}

// ToUint32 converts a color object to a uint32
// for use by the neopixel
func ToUint32(color colorful.Color) uint32 {
	hex := color.Hex()
	hex = strings.Replace(hex, "#", "", -1)
	value, _ := strconv.ParseUint(hex, 16, 32)

	return uint32(value)
}

// GradientPNG generates a gradient PNG as an example
func GradientPNG(gradient GradientTable, h int, w int) error {
	img := image.NewRGBA(image.Rect(0, 0, w, h))

	colorList := GradientColorList(gradient, h)
	for vert, color := range colorList {
		draw.Draw(img, image.Rect(0, vert, w, vert+1), &image.Uniform{color}, image.Point{}, draw.Src)
	}

	outpng, err := os.Create("gradient.png")
	if err != nil {
		return fmt.Errorf("error storing png: %w", err)
	}
	defer outpng.Close()

	err = png.Encode(outpng, img)
	if err != nil {
		return err
	}
	return nil
}

// GradientColorList generates a list of colors for a GradientTable
// length: the number of colors you want
func GradientColorList(gradient GradientTable, length int) []colorful.Color {
	var list []colorful.Color
	for j := 0; j < length; j++ {
		c := gradient.GetInterpolatedColor(float64(j) / float64(length))
		list = append(list, c)
	}
	return list
}

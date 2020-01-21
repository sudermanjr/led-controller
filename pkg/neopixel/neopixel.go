package neopixel

import (
	"time"

	"github.com/lucasb-eyer/go-colorful"
	ws2811 "github.com/rpi-ws281x/rpi-ws281x-go"
	"k8s.io/klog"

	"github.com/sudermanjr/led-controller/pkg/color"
)

type wsEngine interface {
	Init() error
	Render() error
	Wait() error
	Fini()
	Leds(channel int) []uint32
	SetBrightness(channel int, brightness int)
}

// LEDArray is a struct for interacting with LEDs
type LEDArray struct {
	WS            wsEngine
	MaxBrightness int
	MinBrightness int
	Brightness    int
	Color         colorful.Color
	FadeDuration  int
}

var demoGradient = color.GradientTable{
	{color.HexToColor("#9e0142"), 0.0},
	{color.HexToColor("#d53e4f"), 0.1},
	{color.HexToColor("#f46d43"), 0.2},
	{color.HexToColor("#fdae61"), 0.3},
	{color.HexToColor("#fee090"), 0.4},
	{color.HexToColor("#ffffbf"), 0.5},
	{color.HexToColor("#e6f598"), 0.6},
	{color.HexToColor("#abdda4"), 0.7},
	{color.HexToColor("#66c2a5"), 0.8},
	{color.HexToColor("#3288bd"), 0.9},
	{color.HexToColor("#5e4fa2"), 1.0},
}

// NewLEDArray creates a new array and initializes it
func NewLEDArray(minBrightness int, maxBrightness, ledCount int, fadeDuration int) (*LEDArray, error) {
	// Setup the LED lights
	opt := ws2811.DefaultOptions
	opt.Channels[0].Brightness = maxBrightness
	opt.Channels[0].LedCount = ledCount
	dev, err := ws2811.MakeWS2811(&opt)
	if err != nil {
		klog.Error(err)
		return nil, err
	}

	led := &LEDArray{
		WS:            dev,
		Brightness:    minBrightness,
		MinBrightness: minBrightness,
		MaxBrightness: maxBrightness,
		FadeDuration:  fadeDuration,
		Color:         color.HexToColor(color.ColorMap["white"]),
	}

	err = dev.Init()
	if err != nil {
		klog.Error("could not initialize array, did you run as root?")
		klog.Error(err)
		return nil, err
	}

	return led, nil
}

// Display changes all of the LEDs one at a time
// delay: sets the time between each LED coming on
// brightness: sets the brightness for the entire thing
func (led *LEDArray) Display(delay int) error {
	klog.V(6).Infof("setting led array to color: %v, delay: %d, brightness: %d", led.Color, delay, led.Brightness)
	err := led.setBrightness()
	if err != nil {
		return err
	}
	for i := 0; i < len(led.WS.Leds(0)); i++ {
		led.WS.Leds(0)[i] = color.ToUint32(led.Color)
		klog.V(9).Infof("setting led %d to color: %v brightness: %d", i, led.Color, led.Brightness)
		if err := led.WS.Render(); err != nil {
			klog.Error(err)
			return err
		}
		time.Sleep(time.Duration(delay) * time.Millisecond)
	}
	return nil
}

// setBrightness turns the LED array to a brightness value
// and sets the led.brightness value accordingly
// if it goes out of bounds, it will be set to min or max
func (led *LEDArray) setBrightness() error {
	led.checkBrightness()
	klog.V(10).Infof("setting brightness to %d", led.Brightness)
	led.WS.SetBrightness(0, led.Brightness)
	err := led.WS.Render()
	if err != nil {
		return err
	}
	return nil
}

// SetMaxBrightness fades the LED array to maximum brightness
func (led *LEDArray) SetMaxBrightness() error {
	err := led.Fade(led.MaxBrightness)
	if err != nil {
		return err
	}
	return nil
}

// SetMinBrightness fades the LED array to the minimum brightness
func (led *LEDArray) SetMinBrightness() error {
	err := led.Fade(led.MinBrightness)
	if err != nil {
		return err
	}
	return nil
}

// checkBrightness checks to see if the value is
// inside the min/max bounds. If it is out, fix it
func (led *LEDArray) checkBrightness() {
	// Check the bounds
	klog.V(10).Infof("comparing value %d to min: %d, max: %d", led.Brightness, led.MinBrightness, led.MaxBrightness)
	if led.Brightness < led.MinBrightness {
		klog.V(8).Infof("brightness %d below bounds, setting to %d", led.Brightness, led.MinBrightness)
		led.Brightness = led.MinBrightness
		return
	}
	if led.Brightness > led.MaxBrightness {
		klog.V(8).Infof("brightness %d above bounds, setting to %d", led.Brightness, led.MaxBrightness)
		led.Brightness = led.MaxBrightness
		return
	}
	klog.V(8).Infof("not out of bounds. leaving it set as %d", led.Brightness)
}

// Fade goes to a new brightness in the duration specified
func (led *LEDArray) Fade(target int) error {
	klog.V(8).Infof("fading brightness to %d", target)
	klog.V(8).Infof("setting color to %v", led.Color)
	ramp := stepRamp(float64(led.Brightness), float64(target), float64(led.FadeDuration))

	//Set the color on all the LEDs
	for i := 0; i < len(led.WS.Leds(0)); i++ {
		led.WS.Leds(0)[i] = color.ToUint32(led.Color)
	}

	for _, step := range ramp {
		klog.V(10).Infof("processing step: %d", step)
		led.Brightness = step
		err := led.setBrightness()
		if err != nil {
			return err
		}
		time.Sleep(time.Millisecond)
	}
	return nil
}

// stepRamp returns a list of steps in a brightness ramp up
func stepRamp(start float64, stop float64, duration float64) []int {
	slope := (stop - start) / duration
	klog.V(7).Infof("slope of ramp: %f", slope)

	var ramp []int
	for i := 0; i < int(duration); i++ {
		point := start + (slope * float64(i))
		ramp = append(ramp, int(point))
	}
	klog.V(7).Infof("calculated ramp: %v", ramp)
	return ramp
}

//Demo runs a demo of the LED capabilities
func (led *LEDArray) Demo(count int, delay int, gradientLength int) {
	for i := 0; i < (count); i++ {
		for colorName, colorValue := range color.ColorMap {
			klog.Infof("displaying: %s", colorName)
			led.Color = color.HexToColor(colorValue)
			_ = led.Display(delay)
		}
		led.Color = color.HexToColor(color.ColorMap["black"])
		_ = led.Display(0)
		time.Sleep(500 * time.Millisecond)

		// Second part of demo - go through a color gradient really fast.
		klog.V(3).Infof("starting color gradient")
		colorList := color.GradientColorList(demoGradient, gradientLength)
		for _, gradColor := range colorList {
			led.Color = gradColor
			_ = led.Display(0)
			time.Sleep(time.Duration(delay) * time.Nanosecond)
		}
	}
	_ = led.Fade(led.MinBrightness)

}

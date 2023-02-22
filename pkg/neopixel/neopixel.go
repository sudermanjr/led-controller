package neopixel

import (
	"time"

	"github.com/lucasb-eyer/go-colorful"
	ws2811 "github.com/rpi-ws281x/rpi-ws281x-go"
	"go.uber.org/zap"

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

	Logger *zap.SugaredLogger
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
func NewLEDArray(minBrightness int, maxBrightness, ledCount int, fadeDuration int, logger *zap.SugaredLogger) (*LEDArray, error) {
	// Setup the LED lights
	opt := ws2811.DefaultOptions
	opt.Channels[0].Brightness = maxBrightness
	opt.Channels[0].LedCount = ledCount
	dev, err := ws2811.MakeWS2811(&opt)
	if err != nil {
		return nil, err
	}

	led := &LEDArray{
		WS:            dev,
		Brightness:    minBrightness,
		MinBrightness: minBrightness,
		MaxBrightness: maxBrightness,
		FadeDuration:  fadeDuration,
		Color:         color.HexToColor(color.ColorMap["warmwhite"]),
		Logger:        logger,
	}

	err = dev.Init()
	if err != nil {
		led.Logger.Warnf("could not initialize array, did you run as root?")
		return nil, err
	}

	return led, nil
}

// Display changes all of the LEDs one at a time
// delay: sets the time between each LED changing in milliseconds
// brightness: sets the brightness for the entire thing
func (led *LEDArray) Display(delay int) error {
	led.Logger.Debugw("setting led array to color",
		"color", led.Color,
		"delay", delay,
		"brightness", led.Brightness,
	)
	err := led.SetBrightness()
	if err != nil {
		return err
	}
	for i := 0; i < len(led.WS.Leds(0)); i++ {
		led.WS.Leds(0)[i] = color.ToUint32(led.Color)
		led.Logger.Debugw("setting led", "number", i,
			"color", led.Color,
			"brightness", led.Brightness,
		)
		if err := led.WS.Render(); err != nil {
			return err
		}
		time.Sleep(time.Duration(delay) * time.Millisecond)
	}
	return nil
}

// setBrightness turns the LED array to a brightness value
// and sets the led.brightness value accordingly
// if it goes out of bounds, it will be set to min or max
func (led *LEDArray) SetBrightness() error {
	led.checkBrightness()
	led.Logger.Debugw("setting brightness", "value", led.Brightness)
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
	led.Logger.Debugw("comparing value to min/max", "value", led.Brightness, "min", led.MinBrightness, "max", led.MaxBrightness)
	if led.Brightness < led.MinBrightness {
		led.Brightness = led.MinBrightness
		return
	}
	if led.Brightness > led.MaxBrightness {
		led.Brightness = led.MaxBrightness
		return
	}
}

// Fade goes to a new brightness in the duration specified
func (led *LEDArray) Fade(target int) error {
	led.Logger.Debugw("fading led", "target", target, "color", led.Color)
	ramp := stepRamp(float64(led.Brightness), float64(target), float64(led.FadeDuration))
	led.Logger.Debugw("calculated ramp", "ramp", ramp)

	//Set the color on all the LEDs
	for i := 0; i < len(led.WS.Leds(0)); i++ {
		led.WS.Leds(0)[i] = color.ToUint32(led.Color)
	}

	for _, step := range ramp {
		led.Brightness = step
		err := led.SetBrightness()
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

	var ramp []int
	for i := 0; i < int(duration); i++ {
		point := start + (slope * float64(i))
		ramp = append(ramp, int(point))
	}
	return ramp
}

// Demo runs a demo of the LED capabilities
func (led *LEDArray) Demo(count int, delay int, gradientLength int) {
	for i := 0; i < (count); i++ {
		for colorName, colorValue := range color.ColorMap {
			led.Logger.Debugw("diplaying color", "color", colorName)
			led.Color = color.HexToColor(colorValue)

			_ = led.Display(delay)
		}
		led.Color = color.HexToColor(color.ColorMap["black"])
		_ = led.Display(0)
		time.Sleep(500 * time.Millisecond)

		// Second part of demo - go through a color gradient really fast.
		led.Logger.Infow("starting color gradient")
		colorList := color.GradientColorList(demoGradient, gradientLength)
		for _, gradColor := range colorList {
			led.Color = gradColor
			_ = led.Display(delay / 10)
		}
	}
	_ = led.Fade(led.MinBrightness)

}

func (led *LEDArray) FadeToggleOnOff() {
	var err error
	if led.Brightness == led.MinBrightness {
		err = led.SetMaxBrightness()
	} else {
		err = led.SetMinBrightness()
	}
	if err != nil {
		led.Logger.Errorw("could not change brightness from button press", "error", err)
	}
}

func (led *LEDArray) ToggleOnOff() {
	var err error
	led.Color = color.HexToColor(color.ColorMap["warmwhite"])
	if led.Brightness == led.MinBrightness {
		led.Brightness = led.MaxBrightness
		err = led.Display(1)
	} else {
		led.Brightness = led.MinBrightness
		err = led.Display(1)
	}
	if err != nil {
		led.Logger.Errorw("could not change brightness from button press", "error", err)
	}
}

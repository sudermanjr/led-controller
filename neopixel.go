package main

import (
	"time"

	"github.com/lucasb-eyer/go-colorful"
	ws2811 "github.com/rpi-ws281x/rpi-ws281x-go"
	"k8s.io/klog"
)

var colors = map[string]string{
	"blue":   "#0000ff",
	"green":  "#00ff00",
	"yellow": "#ffaf33",
	"purple": "#af33ff",
	"red":    "#ff0000",
	"teal":   "#33ffd1",
	"pink":   "#ff08c7",
	"white":  "#ffffff",
	"black":  "#000000", // This basically equates to off.
}

type wsEngine interface {
	Init() error
	Render() error
	Wait() error
	Fini()
	Leds(channel int) []uint32
	SetBrightness(channel int, brightness int)
}

// ledArray is a struct for interacting with LEDs
type ledArray struct {
	ws         wsEngine
	brightness int
	color      colorful.Color
}

func newledArray() (*ledArray, error) {
	// Setup the LED lights
	opt := ws2811.DefaultOptions
	opt.Channels[0].Brightness = maxBrightness
	opt.Channels[0].LedCount = ledCount
	dev, err := ws2811.MakeWS2811(&opt)
	if err != nil {
		klog.Error(err)
		return nil, err
	}

	led := &ledArray{
		ws: dev,
	}

	err = dev.Init()
	if err != nil {
		klog.Error("could not initialize array, did you run as root?")
		klog.Error(err)
		return nil, err
	}
	// Start with brightness off and color white
	led.brightness = minBrightness
	led.color = HexToColor(colors["white"])

	return led, nil
}

// display changes all of the LEDs one at a time
// delay: sets the time between each LED coming on
// brightness: sets the brightness for the entire thing
func (led *ledArray) display(delay int) error {
	klog.V(6).Infof("setting led array to color: %v, delay: %d, brightness: %d", led.color, delay, led.brightness)
	err := led.setBrightness(led.brightness)
	if err != nil {
		return err
	}
	for i := 0; i < len(led.ws.Leds(0)); i++ {
		led.ws.Leds(0)[i] = ColorToUint32(led.color)
		klog.V(10).Infof("setting led %d", i)
		if err := led.ws.Render(); err != nil {
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
func (led *ledArray) setBrightness(value int) error {
	value = brightnessBounds(value)
	klog.V(8).Infof("setting brightness to %d", value)
	led.ws.SetBrightness(0, value)
	err := led.ws.Render()
	if err != nil {
		return err
	}
	led.brightness = value
	klog.V(10).Infof("storing brightness as %d", led.brightness)
	return nil
}

// brightnessBounds checks to see if the value is
// inside the min/max bounds. If it is out, return
// the appropriate min or max
func brightnessBounds(value int) int {
	// Check the bounds
	klog.V(10).Infof("comparing value %d to min: %d, max: %d", value, minBrightness, maxBrightness)
	if value < minBrightness {
		klog.V(8).Infof("brightness %d below bounds, setting to %d", value, minBrightness)
		return minBrightness
	} else if value > maxBrightness {
		klog.V(8).Infof("brightness %d above bounds, setting to %d", value, maxBrightness)
		return maxBrightness
	}
	klog.V(8).Infof("not out of bounds. returning %d", value)
	return value
}

// fade goes to a new brightness in the duration specified
func (led *ledArray) fade(target int) error {
	klog.V(8).Infof("fading brightness to %d", target)
	klog.V(8).Infof("setting color to %v", led.color)
	ramp := stepRamp(float64(led.brightness), float64(target), float64(fadeDuration))

	//Set the color on all the LEDs
	for i := 0; i < len(led.ws.Leds(0)); i++ {
		led.ws.Leds(0)[i] = ColorToUint32(led.color)
	}

	for _, step := range ramp {
		klog.V(10).Infof("processing step: %d", step)
		err := led.setBrightness(step)
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

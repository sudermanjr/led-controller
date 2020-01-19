package main

import (
	"time"

	ws2811 "github.com/rpi-ws281x/rpi-ws281x-go"
	"k8s.io/klog"
)

var colors = map[string]uint32{
	"blue":   uint32(0x0000ff),
	"green":  uint32(0x00ff00),
	"yellow": uint32(0xffaf33),
	"purple": uint32(0xaf33ff),
	"red":    uint32(0xff0000),
	"teal":   uint32(0x33ffd1),
	"pink":   uint32(0xff08c7),
	"white":  uint32(0xffffff),
}

const off = uint32(0x000000)

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
	ws         wsEngine
	brightness int
	color      uint32
}

func newLEDArray() (*LEDArray, error) {
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
		ws: dev,
	}

	err = dev.Init()
	if err != nil {
		klog.Error(err)
		return nil, err
	}
	// Start off
	led.brightness = 0
	led.color = off
	return led, nil
}

// display changes all of the LEDs one at a time
// delay: sets the time between each LED coming on
// brightness: sets the brightness for the entire thing
func (led *LEDArray) display(color uint32, delay int, brightness int) error {
	led.ws.SetBrightness(0, brightness)
	for i := 0; i < len(led.ws.Leds(0)); i++ {
		led.ws.Leds(0)[i] = color
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
func (led *LEDArray) setBrightness(value int) error {

	value = brightnessBounds(value)
	led.ws.SetBrightness(0, value)
	err := led.ws.Render()
	if err != nil {
		return err
	}
	led.brightness = value
	return nil
}

// brightnessBounds checks to see if the value is
// inside the min/max bounds. If it is out, return
// the appropriate min or max
func brightnessBounds(value int) int {
	// Check the bounds
	klog.V(10).Infof("comparing value %d to min: %d, max: %d", value, minBrightness, maxBrightness)
	if value < minBrightness {
		return minBrightness
	} else if value > maxBrightness {
		return maxBrightness
	}
	return value
}

// fade goes to a new brightness in the duration specified
func (led *LEDArray) fade(color uint32, target int) error {

	ramp := stepRamp(float64(led.brightness), float64(target), float64(fadeDuration))

	//Set the color on all the LEDs
	for i := 0; i < len(led.ws.Leds(0)); i++ {
		led.ws.Leds(0)[i] = color
	}

	for _, step := range ramp {
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

	var ramp []int
	for i := 0; i < int(duration); i++ {
		point := start + (slope * float64(i))
		ramp = append(ramp, int(point))
	}
	return ramp
}

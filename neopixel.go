package main

import (
	"math"
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
	ws wsEngine
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

	cw := &LEDArray{
		ws: dev,
	}

	err = dev.Init()
	if err != nil {
		klog.Error(err)
		return nil, err
	}
	return cw, nil
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

// fade goes to a new brightness in the duration specified
func (led *LEDArray) fade(color uint32, start int, target int, duration int) error {

	// Number of steps of brightness to go per millisecond
	stepSize := math.Abs(float64((target - start) / duration))

	ramp := stepRamp(float64(duration), stepSize)

	//Set the color on all the LEDs
	for i := 0; i < len(led.ws.Leds(0)); i++ {
		led.ws.Leds(0)[i] = color
	}

	//Fade in
	for _, step := range ramp {
		led.ws.SetBrightness(0, step)
		err := led.ws.Render()
		if err != nil {
			return err
		}
		time.Sleep(time.Millisecond)
	}
	return nil
}

// stepRamp returns a list of steps in a brightness ramp up
func stepRamp(steps float64, size float64) []int {
	var ramp []int
	step := 0
	for i := 0; i < int(steps); i++ {
		ramp = append(ramp, step)
		step = step + int(size)
	}
	return ramp
}

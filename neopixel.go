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

func (led *LEDArray) display(color uint32, delay int) error {
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

func (led *LEDArray) fade(color uint32, brightness int, ledDelay int, brightnessDelay int) error {
	for b := 0; b < brightness; b++ {
		led.ws.SetBrightness(0, b)
		time.Sleep(time.Duration(brightnessDelay) * time.Millisecond)

		err := led.display(color, ledDelay)
		if err != nil {
			klog.Error(err)
			return err
		}
	}
	return nil
}

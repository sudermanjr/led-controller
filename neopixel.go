package main

import (
	"time"

	ws2811 "github.com/rpi-ws281x/rpi-ws281x-go"
)

const (
	brightness = 100
	ledCounts  = 12
	sleepTime  = 100
	blue       = uint32(0x0000ff)
	green      = uint32(0x00ff00)
	yellow     = uint32(0xffaf33)
	purple     = uint32(0xaf33ff)
	red        = uint32(0xff0000)
	teal       = uint32(0x33ffd1)
	pink       = uint32(0xff08c7)
	off        = uint32(0x000000)
)

type wsEngine interface {
	Init() error
	Render() error
	Wait() error
	Fini()
	Leds(channel int) []uint32
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

type colorWipe struct {
	ws wsEngine
}

func (cw *colorWipe) setup() error {
	return cw.ws.Init()
}

func (cw *colorWipe) display(color uint32) error {
	for i := 0; i < len(cw.ws.Leds(0)); i++ {
		cw.ws.Leds(0)[i] = color
		if err := cw.ws.Render(); err != nil {
			return err
		}
		time.Sleep(sleepTime * time.Millisecond)
	}
	return nil
}

// Demo runs a demo of the lights
func Demo() {
	opt := ws2811.DefaultOptions
	opt.Channels[0].Brightness = brightness
	opt.Channels[0].LedCount = ledCounts

	dev, err := ws2811.MakeWS2811(&opt)
	checkError(err)

	cw := &colorWipe{
		ws: dev,
	}
	checkError(cw.setup())
	defer dev.Fini()

	_ = cw.display(blue)
	_ = cw.display(green)
	_ = cw.display(yellow)
	_ = cw.display(purple)
	_ = cw.display(red)
	_ = cw.display(teal)
	_ = cw.display(pink)
	_ = cw.display(off)
}

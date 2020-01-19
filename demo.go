package main

import (
	ws2811 "github.com/rpi-ws281x/rpi-ws281x-go"
)

// Demo runs a demo of the lights
func Demo() {
	opt := ws2811.DefaultOptions

	opt.Channels[0].Brightness = brightness
	opt.Channels[0].LedCount = ledCount

	dev, err := ws2811.MakeWS2811(&opt)
	checkError(err)

	cw := &colorWipe{
		ws:    dev,
		delay: demoDelay,
	}
	checkError(cw.setup())
	defer dev.Fini()

	for i := 1; i < demoCount; i++ {
		_ = cw.display(blue)
		_ = cw.display(green)
		_ = cw.display(yellow)
		_ = cw.display(purple)
		_ = cw.display(red)
		_ = cw.display(teal)
		_ = cw.display(pink)
		_ = cw.display(white)
	}

	_ = cw.display(off)
}

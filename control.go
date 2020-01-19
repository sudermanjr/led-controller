package main

import (
	ws2811 "github.com/rpi-ws281x/rpi-ws281x-go"
)

// On turns on the lights
func On(colorName string) {
	opt := ws2811.DefaultOptions

	opt.Channels[0].Brightness = brightness
	opt.Channels[0].LedCount = ledCount

	dev, err := ws2811.MakeWS2811(&opt)
	checkError(err)

	cw := &colorWipe{
		ws: dev,
	}
	checkError(cw.setup())
	defer dev.Fini()

	_ = cw.on(colors[colorName])
}

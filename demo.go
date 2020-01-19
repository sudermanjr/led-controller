package main

import (
	ws2811 "github.com/rpi-ws281x/rpi-ws281x-go"
	"k8s.io/klog"
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

	for i := 1; i < (demoCount + 1); i++ {
		for colorName, color := range colors {
			klog.Infof("displaying: %s", colorName)
			_ = cw.display(color)
		}
	}

	_ = cw.display(off)
}

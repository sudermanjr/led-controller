package main

import (
	"time"
)

const (
	blue   = uint32(0x0000ff)
	green  = uint32(0x00ff00)
	yellow = uint32(0xffaf33)
	purple = uint32(0xaf33ff)
	red    = uint32(0xff0000)
	teal   = uint32(0x33ffd1)
	pink   = uint32(0xff08c7)
	white  = uint32(0xffffff)
	off    = uint32(0x000000)
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
	ws    wsEngine
	delay int
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
		time.Sleep(time.Duration(cw.delay) * time.Millisecond)
	}
	return nil
}

package dashboard

import (
	"time"

	"github.com/stianeikeland/go-rpio/v4"
)

func (a *App) WatchButton() {
	if err := rpio.Open(); err != nil {
		a.Logger.Errorw("could not init gpio", "error", err)
		return
	}
	defer rpio.Close()
	buttonPin := rpio.Pin(a.ButtonPin)

	buttonPin.Input()
	buttonPin.PullUp()
	buttonPin.Detect(rpio.FallEdge)     // enable falling edge event detection
	defer buttonPin.Detect(rpio.NoEdge) // disable edge event detection

	for {
		if buttonPin.EdgeDetected() { // check if event occured
			a.Logger.Debugw("button pressed", "gpio", a.ButtonPin)
		}
		time.Sleep(time.Second / 2)
	}
}

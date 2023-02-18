package dashboard

import (
	"fmt"
	"time"

	"github.com/stianeikeland/go-rpio/v4"
)

func (a *App) WatchButton() {
	if err := rpio.Open(); err != nil {
		fmt.Println(err)
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
			fmt.Println("button pressed")
		}
		time.Sleep(time.Second / 2)
	}
}

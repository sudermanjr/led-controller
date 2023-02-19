package dashboard

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
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
			jsonStr := []byte(`{"button":"pressed"}`)

			res, err := http.Post(fmt.Sprintf("http://localhost:%d/button", a.Port), "application/json", bytes.NewBuffer(jsonStr))
			if err != nil {
				a.Logger.Errorw("error making button request", "error", err)
				continue
			}
			body, _ := io.ReadAll(res.Body)
			a.Logger.Infow("got response from button press", "statusCode", res.StatusCode, "headers", res.Header, "body", body)
			res.Body.Close()
		}
		time.Sleep(time.Second / 2)
	}
}

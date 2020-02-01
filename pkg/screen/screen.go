package screen

import (
	"time"

	lcd1602 "github.com/pimvanhespen/go-pi-lcd1602"
	"github.com/pimvanhespen/go-pi-lcd1602/synchronized"
)

// Demo runs a screen demo
func Demo() {
	lcdi := lcd1602.New(
		25,                    //rs
		24,                    //enable
		[]int{23, 17, 18, 22}, //datapins
		16,                    //lineSize
	)
	lcd := synchronized.NewSynchronizedLCD(lcdi)
	lcd.Initialize()
	lcd.WriteLines("The LCD Screen ", " is working   ")
	time.Sleep(5 * time.Second)
	lcd.Clear()
	lcd.Close()
}

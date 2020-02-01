package screen

import (
	lcd1602 "github.com/pimvanhespen/go-pi-lcd1602"
	"github.com/pimvanhespen/go-pi-lcd1602/animations"
	"github.com/pimvanhespen/go-pi-lcd1602/stringutils"
	"github.com/pimvanhespen/go-pi-lcd1602/synchronized"
	"github.com/sudermanjr/led-controller/pkg/utils"
	"k8s.io/klog"
)

// Display is a screen that you can display info on
type Display struct {
	LCD *synchronized.SynchronizedLCD
}

// NewDisplay returns a *Display
func NewDisplay(rs int, e int, dp []int, ls int) (*Display, error) {
	lcdi := lcd1602.New(rs, e, dp, ls)
	lcd := synchronized.NewSynchronizedLCD(lcdi)
	lcd.Initialize()
	obj := &Display{
		LCD: lcd,
	}
	return obj, nil
}

// Demo runs a screen demo
func (display *Display) Demo() {

	animations := []animations.Animation{
		animations.None(stringutils.Center("lcd demo", 16)),
		animations.GarbleLeftSimple(stringutils.Center("garble left", 16)),
		animations.GarbleRightSimple(stringutils.Center("garble right", 16)),
		animations.SlideInLeft(stringutils.Center("slide in left", 16)),
		animations.SlideInRight(stringutils.Center("slide in right", 16)),
		animations.SlideOutLeft(stringutils.Center("slide out left", 16)),
		animations.SlideOutRight(stringutils.Center("slide out right", 16)),
	}

	for index, animation := range animations {
		line := lcd1602.LINE_1
		if index%2 == 0 {
			line = lcd1602.LINE_2
		}
		wait := display.LCD.Animate(animation, line)
		<-wait
	}
}

// InfoDisplay shows current information
func (display *Display) InfoDisplay() error {
	ip, err := utils.IPAddress()
	if err != nil {
		klog.Error(err)
		return err
	}
	display.LCD.WriteLines("LED Controller ", stringutils.Center(ip, 16))
	return nil
}

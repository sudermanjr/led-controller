package homekit

import (
	"github.com/brutella/hc"
	"github.com/brutella/hc/accessory"
	"github.com/lucasb-eyer/go-colorful"
	"k8s.io/klog"

	"github.com/sudermanjr/led-controller/pkg/neopixel"
	"github.com/sudermanjr/led-controller/pkg/utils"
)

//Start starts the homekit server
func Start(homekitPin string, led *neopixel.LEDArray) {
	// create an accessory
	info := accessory.Info{
		Name:         "LED",
		Manufacturer: "sudermanjr",
		Model:        "led-controller",
	}
	ac := accessory.NewLightbulb(info)

	// configure the ip transport
	config := hc.Config{Pin: homekitPin}
	t, err := hc.NewIPTransport(config, ac.Accessory)
	if err != nil {
		klog.Fatal(err)
	}

	ac.Lightbulb.On.OnValueRemoteUpdate(func(on bool) {
		if on {
			klog.Infof("Switch is on")
			err = led.Fade(led.MaxBrightness)
			if err != nil {
				klog.Error(err)
			}
		} else {
			klog.Infof("Switch is off")
			err = led.Fade(led.MinBrightness)
			if err != nil {
				klog.Error(err)
			}
		}
	})

	ac.Lightbulb.Hue.OnValueRemoteUpdate(func(value float64) {
		klog.Infof("homekit hue set to: %f", value)
		led.Color = modifyHue(led.Color, value)
		err = led.Display(0)
		if err != nil {
			klog.Error(err)
		}
	})

	ac.Lightbulb.Saturation.OnValueRemoteUpdate(func(value float64) {
		klog.Infof("homekit saturation set to %f", value)
		led.Color = modifySaturation(led.Color, value)
		err = led.Display(0)
		if err != nil {
			klog.Error(err)
		}
	})

	ac.Lightbulb.Brightness.OnValueRemoteUpdate(func(value int) {
		klog.Infof("homekit brightness set to: %d", value)
		err = led.Fade(utils.ScaleBrightness(value, led.MinBrightness, led.MaxBrightness))
		if err != nil {
			klog.Error(err)
		}
	})

	hc.OnTermination(func() {
		klog.Info("terminated. turning off lights")
		err = led.Fade(led.MinBrightness)
		if err != nil {
			klog.Error(err)
		}
		<-t.Stop()
	})

	klog.Info("starting homekit server...")
	klog.Infof("max-brightness: %d", led.MaxBrightness)
	klog.Infof("min-brightness: %d", led.MinBrightness)
	klog.Infof("fade-duration %d", led.FadeDuration)

	t.Start()
}

//modifySaturation changes the saturation and returns a new color
func modifySaturation(oldColor colorful.Color, saturation float64) colorful.Color {
	h, s, v := oldColor.Hsv()
	klog.V(8).Infof("old color h: %f, s: %f, v: %f", h, s, v)
	s = saturation * .1 // hc sends this 1-100, but colorful uses 0-1
	newColor := colorful.Hsv(h, s, v)
	klog.V(8).Infof("new color h: %f, s: %f, v: %f", h, s, v)
	return newColor
}

//modifyHue changes the hue and returns a new color
func modifyHue(oldColor colorful.Color, hue float64) colorful.Color {
	h, s, v := oldColor.Hsv()
	klog.V(8).Infof("old color h: %f, s: %f, v: %f", h, s, v)
	h = hue
	newColor := colorful.Hsv(h, s, v)
	klog.V(8).Infof("new color h: %f, s: %f, v: %f", h, s, v)
	return newColor
}

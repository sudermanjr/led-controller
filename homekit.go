package main

import (
	"github.com/brutella/hc"
	"github.com/brutella/hc/accessory"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/spf13/cobra"
	"k8s.io/klog"
)

var homekitPin string

func init() {
	rootCmd.AddCommand(homekitCmd)

	homekitCmd.Flags().StringVar(&homekitPin, "homekit-pin", "29847290", "The pin that homekit will use to authenticate with this device.")
}

var homekitCmd = &cobra.Command{
	Use:   "homekit",
	Short: "Run the lights as a homekit accessory.",
	Long:  `Run the lights as a homekit accessory.`,
	Run: func(cmd *cobra.Command, args []string) {
		startHomekit()
	},
}

func startHomekit() {
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

	led, err := newledArray()
	if err != nil {
		klog.Fatal(err)
	}
	defer led.ws.Fini()

	ac.Lightbulb.On.OnValueRemoteUpdate(func(on bool) {
		if on {
			klog.Infof("Switch is on")
			err = led.fade(maxBrightness)
			if err != nil {
				klog.Error(err)
			}
		} else {
			klog.Infof("Switch is off")
			err = led.fade(minBrightness)
			if err != nil {
				klog.Error(err)
			}
		}
	})

	ac.Lightbulb.Hue.OnValueRemoteUpdate(func(value float64) {
		klog.Infof("homekit hue set to: %f", value)
		led.color = modifyHue(led.color, value)
		err = led.display(0)
		if err != nil {
			klog.Error(err)
		}
	})

	ac.Lightbulb.Saturation.OnValueRemoteUpdate(func(value float64) {
		klog.Infof("homekit saturation set to %f", value)
		led.color = modifySaturation(led.color, value)
		err = led.display(0)
		if err != nil {
			klog.Error(err)
		}
	})

	ac.Lightbulb.Brightness.OnValueRemoteUpdate(func(value int) {
		klog.Infof("homekit brightness set to: %d", value)
		err = led.fade(scaleHomekitBrightness(value))
		if err != nil {
			klog.Error(err)
		}
	})

	hc.OnTermination(func() {
		klog.Info("terminated. turning off lights")
		err = led.fade(minBrightness)
		if err != nil {
			klog.Error(err)
		}
		<-t.Stop()
	})

	klog.Info("starting homekit server...")
	klog.Infof("max-brightness: %d", maxBrightness)
	klog.Infof("min-brightness: %d", minBrightness)
	klog.Infof("fade-duration %d", fadeDuration)

	t.Start()
}

// scaleHomekitBrightness converts a 0-100 homekit brightness
// to the scale of the controller (min - max)
// math isn't as easy as it used to be for me:
// https://stackoverflow.com/questions/5294955/how-to-scale-down-a-range-of-numbers-with-a-known-min-and-max-value
func scaleHomekitBrightness(value int) int {
	min := 0
	max := 100
	a := minBrightness
	b := maxBrightness

	new := ((b-a)*(value-min))/(max-min) + a

	return new
}

//modifySaturation changes the saturation and returns a new color
func modifySaturation(oldColor colorful.Color, saturation float64) colorful.Color {
	h, s, v := oldColor.Hsv()
	klog.V(8).Infof("old color h: %f, s: %f, v: %f", h, s, v)
	s = saturation
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

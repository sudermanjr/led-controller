package main

import (
	"github.com/brutella/hc"
	"github.com/brutella/hc/accessory"
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

	led, err := newLEDArray()
	if err != nil {
		klog.Fatal(err)
	}
	defer led.ws.Fini()

	ac.Lightbulb.On.OnValueRemoteUpdate(func(on bool) {
		if on {
			klog.Infof("Switch is on")
			err = led.fade(colors["white"], 150)
			if err != nil {
				klog.Error(err)
			}
		} else {
			klog.Infof("Switch is off")
			err = led.display(off, 0, 0)
			if err != nil {
				klog.Error(err)
			}
		}
	})

	ac.Lightbulb.Hue.OnValueRemoteUpdate(func(value float64) {
		klog.Infof("homekit hue set to: %f", value)
	})

	ac.Lightbulb.Saturation.OnValueRemoteUpdate(func(value float64) {
		klog.Infof("homekit saturation set to %f", value)
	})

	ac.Lightbulb.Brightness.OnValueRemoteUpdate(func(value int) {
		klog.Infof("homekit brightness set to: %d", value)
		err = led.fade(colors["white"], scaleHomekitBrightness(value))
		if err != nil {
			klog.Error(err)
		}
	})

	hc.OnTermination(func() {
		klog.Info("terminated. turning off lights")
		err = led.display(off, 0, 0)
		if err != nil {
			klog.Error(err)
		}
		<-t.Stop()
	})

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

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
			err = led.fade(colors["white"], 0, 150, 20)
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
		klog.Infof("hue set to: %f", value)
	})

	ac.Lightbulb.Saturation.OnValueRemoteUpdate(func(value float64) {
		klog.Infof("saturation set to %f", value)
	})

	ac.Lightbulb.Brightness.OnValueRemoteUpdate(func(value int) {
		klog.Infof("brightness set to: %d", value)
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

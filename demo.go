package main

import (
	ws2811 "github.com/rpi-ws281x/rpi-ws281x-go"
	"github.com/spf13/cobra"
	"k8s.io/klog"
)

func init() {
	rootCmd.AddCommand(demoCmd)

	demoCmd.Flags().IntVar(&demoDelay, "speed", 100, "The delay in ms of the demo program.")
	demoCmd.Flags().IntVar(&demoCount, "count", 1, "The number of loops to run the demo.")
}

var demoCmd = &cobra.Command{
	Use:   "demo",
	Short: "Run a demo.",
	Long:  `Runs a demo.`,
	Run: func(cmd *cobra.Command, args []string) {

		opt := ws2811.DefaultOptions

		opt.Channels[0].Brightness = brightness
		opt.Channels[0].LedCount = ledCount

		dev, err := ws2811.MakeWS2811(&opt)
		checkError(err)

		cw := &colorWipe{
			ws: dev,
		}
		checkError(cw.setup())
		defer dev.Fini()

		for i := 1; i < (demoCount + 1); i++ {
			for colorName, color := range colors {
				klog.Infof("displaying: %s", colorName)
				_ = cw.display(color, demoDelay)
			}
		}

		_ = cw.display(off, demoDelay)
	},
}

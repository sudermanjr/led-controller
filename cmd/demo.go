package cmd

import (
	"github.com/spf13/cobra"
	"k8s.io/klog"

	"github.com/sudermanjr/led-controller/pkg/neopixel"
)

var (
	demoBrightness     int
	demoDelay          int
	demoCount          int
	demoGradientLength int
)

func init() {
	rootCmd.AddCommand(demoCmd)

	demoCmd.Flags().IntVar(&demoDelay, "delay", 100, "The delay in ms of the demo program.")
	demoCmd.Flags().IntVar(&demoCount, "count", 1, "The number of loops to run the demo.")
	demoCmd.Flags().IntVar(&demoBrightness, "brightness", 150, "The brightness to run the demo at. Must be between min and max.")
	demoCmd.Flags().IntVar(&demoGradientLength, "gradient-count", 2048, "The number of steps in the gradient.")
}

var demoCmd = &cobra.Command{
	Use:   "demo",
	Short: "Run a demo.",
	Long:  `Runs a demo.`,
	Run: func(cmd *cobra.Command, args []string) {

		// Initialize the LEDs
		led, err := neopixel.NewLEDArray(minBrightness, maxBrightness, ledCount, fadeDuration)
		if err != nil {
			klog.Fatal(err)
		}
		defer led.WS.Fini()

		led.Brightness = demoBrightness
		led.Demo(demoCount, demoDelay, demoGradientLength)
	},
}

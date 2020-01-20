package cmd

import (
	"time"

	"github.com/spf13/cobra"
	"k8s.io/klog"

	"github.com/sudermanjr/led-controller/pkg/color"
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

var demoGradient = color.GradientTable{
	{color.HexToColor("#9e0142"), 0.0},
	{color.HexToColor("#d53e4f"), 0.1},
	{color.HexToColor("#f46d43"), 0.2},
	{color.HexToColor("#fdae61"), 0.3},
	{color.HexToColor("#fee090"), 0.4},
	{color.HexToColor("#ffffbf"), 0.5},
	{color.HexToColor("#e6f598"), 0.6},
	{color.HexToColor("#abdda4"), 0.7},
	{color.HexToColor("#66c2a5"), 0.8},
	{color.HexToColor("#3288bd"), 0.9},
	{color.HexToColor("#5e4fa2"), 1.0},
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
		// Loops through our list of pre-defined colors and display them in order.
		for i := 0; i < (demoCount); i++ {
			for colorName, colorValue := range color.ColorMap {
				klog.Infof("displaying: %s", colorName)
				led.Color = color.HexToColor(colorValue)
				_ = led.Display(demoDelay)
			}
			_ = led.Fade(minBrightness)
			time.Sleep(500 * time.Millisecond)

			// Second part of demo - go through a color gradient really fast.
			klog.V(3).Infof("starting color gradient")
			colorList := color.GradientColorList(demoGradient, demoGradientLength)
			for _, gradColor := range colorList {
				led.Color = gradColor
				led.Brightness = demoBrightness
				_ = led.Display(0)
				time.Sleep(time.Duration(demoDelay) * time.Nanosecond)
			}
		}

		_ = led.Fade(minBrightness)
	},
}

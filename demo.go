package main

import (
	"github.com/spf13/cobra"
	"k8s.io/klog"
	"time"
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

var demoGradient = GradientTable{
	{HexToColor("#9e0142"), 0.0},
	{HexToColor("#d53e4f"), 0.1},
	{HexToColor("#f46d43"), 0.2},
	{HexToColor("#fdae61"), 0.3},
	{HexToColor("#fee090"), 0.4},
	{HexToColor("#ffffbf"), 0.5},
	{HexToColor("#e6f598"), 0.6},
	{HexToColor("#abdda4"), 0.7},
	{HexToColor("#66c2a5"), 0.8},
	{HexToColor("#3288bd"), 0.9},
	{HexToColor("#5e4fa2"), 1.0},
}

var demoCmd = &cobra.Command{
	Use:   "demo",
	Short: "Run a demo.",
	Long:  `Runs a demo.`,
	Run: func(cmd *cobra.Command, args []string) {

		// Initialize the LEDs
		led, err := newledArray()
		if err != nil {
			klog.Fatal(err)
		}
		defer led.ws.Fini()

		led.brightness = demoBrightness
		// Loops through our list of pre-defined colors and display them in order.
		for i := 0; i < (demoCount); i++ {
			for colorName, color := range colors {
				klog.Infof("displaying: %s", colorName)
				led.color = HexToColor(color)
				_ = led.display(demoDelay)
			}
			_ = led.fade(minBrightness)
			time.Sleep(500 * time.Millisecond)

			// Second part of demo - go through a color gradient really fast.
			klog.V(3).Infof("starting color gradient")
			colorList := GradientColorList(demoGradient, demoGradientLength)
			for _, gradColor := range colorList {
				led.color = gradColor
				led.brightness = demoBrightness
				_ = led.display(0)
				time.Sleep(time.Duration(demoDelay) * time.Nanosecond)
			}
		}

		_ = led.fade(minBrightness)
	},
}

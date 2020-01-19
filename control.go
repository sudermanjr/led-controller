package main

import (
	"github.com/spf13/cobra"
	"k8s.io/klog"
)

var onBrightness int

func init() {
	rootCmd.AddCommand(onCmd)
	rootCmd.AddCommand(offCmd)
	onCmd.Flags().StringVarP(&colorName, "color", "c", "white", "The color to turn the lights on to.")
	onCmd.Flags().IntVar(&onBrightness, "brightness", 100, "The brightness setting. Range is 1-250.")

}

var onCmd = &cobra.Command{
	Use:   "on",
	Short: "Turn on the lights.",
	Long:  `Turns on the lights to a specific color.`,
	Run: func(cmd *cobra.Command, args []string) {
		led, err := newLEDArray()
		if err != nil {
			klog.Fatal(err)
		}
		defer led.ws.Fini()

		_ = led.fade(colors[colorName], 0, onBrightness, 10)
	},
}

var offCmd = &cobra.Command{
	Use:   "off",
	Short: "Turn off the lights.",
	Long:  `Turns off the lights.`,
	Run: func(cmd *cobra.Command, args []string) {
		led, err := newLEDArray()
		if err != nil {
			klog.Fatal(err)
		}
		defer led.ws.Fini()

		_ = led.display(off, 0, 0)
	},
}

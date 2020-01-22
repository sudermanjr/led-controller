package cmd

import (
	"github.com/spf13/cobra"
	"k8s.io/klog"

	"github.com/sudermanjr/led-controller/pkg/color"
	"github.com/sudermanjr/led-controller/pkg/neopixel"
)

var onBrightness int

func init() {
	rootCmd.AddCommand(onCmd)
	rootCmd.AddCommand(offCmd)
	onCmd.Flags().StringVarP(&colorName, "color", "c", "white", "The color to turn the lights on to.")
	onCmd.Flags().IntVar(&onBrightness, "brightness", 100, "The brightness setting. Range is between the min-brightness and max-brightness.")

}

var onCmd = &cobra.Command{
	Use:   "on",
	Short: "Turn on the lights.",
	Long:  `Turns on the lights to a specific color.`,
	Run: func(cmd *cobra.Command, args []string) {
		led, err := neopixel.NewLEDArray(minBrightness, maxBrightness, ledCount, fadeDuration)
		if err != nil {
			klog.Fatal(err)
		}
		defer led.WS.Fini()
		led.Color = color.HexToColor(color.ColorMap[colorName])
		led.Brightness = onBrightness
		err = led.Display(0)
		if err != nil {
			klog.Fatal(err)
		}
	},
}

var offCmd = &cobra.Command{
	Use:   "off",
	Short: "Turn off the lights.",
	Long:  `Turns off the lights.`,
	Run: func(cmd *cobra.Command, args []string) {
		led, err := neopixel.NewLEDArray(minBrightness, maxBrightness, ledCount, fadeDuration)
		if err != nil {
			klog.Fatal(err)
		}
		defer led.WS.Fini()
		err = led.SetMinBrightness()
		if err != nil {
			klog.Fatal(err)
		}
	},
}

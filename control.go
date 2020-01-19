package main

import (
	ws2811 "github.com/rpi-ws281x/rpi-ws281x-go"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(onCmd)
	rootCmd.AddCommand(offCmd)
	onCmd.Flags().StringVarP(&colorName, "color", "c", "white", "The color to turn the lights on to.")

}

var onCmd = &cobra.Command{
	Use:   "on",
	Short: "Turn on the lights.",
	Long:  `Turns on the lights to a specific color.`,
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

		_ = cw.display(colors[colorName], 0)
	},
}

var offCmd = &cobra.Command{
	Use:   "off",
	Short: "Turn off the lights.",
	Long:  `Turns off the lights.`,
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

		_ = cw.display(off, 0)
	},
}

package cmd

import (
	"github.com/spf13/cobra"
	"k8s.io/klog"

	"github.com/sudermanjr/led-controller/pkg/homekit"
	"github.com/sudermanjr/led-controller/pkg/neopixel"
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

		led, err := neopixel.NewLEDArray(minBrightness, maxBrightness, ledCount, fadeDuration)
		if err != nil {
			klog.Fatal(err)
		}
		defer led.WS.Fini()

		homekit.Start(homekitPin, led)
	},
}

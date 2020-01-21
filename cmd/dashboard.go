package cmd

import (
	"github.com/spf13/cobra"

	"k8s.io/klog"

	"github.com/sudermanjr/led-controller/pkg/dashboard"
	"github.com/sudermanjr/led-controller/pkg/neopixel"
)

var (
	serverPort int
)

func init() {
	rootCmd.AddCommand(dashboardCmd)
	dashboardCmd.PersistentFlags().IntVarP(&serverPort, "port", "p", 8080, "The port to serve the dashboard on.")
}

var dashboardCmd = &cobra.Command{
	Use:   "dashboard",
	Short: "Run a dashboard",
	Long:  `Run a dashboard`,
	Run: func(cmd *cobra.Command, args []string) {
		// Initialize the LEDs
		led, err := neopixel.NewLEDArray(minBrightness, maxBrightness, ledCount, fadeDuration)
		if err != nil {
			klog.Fatal(err)
		}

		app := dashboard.App{
			Array: led,
			Port:  serverPort,
		}
		app.Initialize()
		app.Run()
	},
}

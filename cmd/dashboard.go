package cmd

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"

	"k8s.io/klog"

	"github.com/sudermanjr/led-controller/pkg/dashboard"
	"github.com/sudermanjr/led-controller/pkg/homekit"
	"github.com/sudermanjr/led-controller/pkg/neopixel"
	"github.com/sudermanjr/led-controller/pkg/screen"
)

var (
	serverPort int
	homekitPin string
)

func init() {
	rootCmd.AddCommand(dashboardCmd)
	dashboardCmd.PersistentFlags().IntVarP(&serverPort, "port", "p", 8080, "The port to serve the dashboard on.")
	dashboardCmd.PersistentFlags().StringVar(&homekitPin, "homekit-pin", "29847290", "The pin that homekit will use to authenticate with this device.")
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

		// Initialize the LCD display
		display, err := screen.NewDisplay(displayRSPin, displayEPin, displayDataPins, lineSize)
		if err != nil {
			klog.Fatal(err)
		}
		defer display.LCD.Close()

		app := dashboard.App{
			Array:  led,
			Port:   serverPort,
			Screen: display,
		}
		app.Initialize()

		go homekit.Start(homekitPin, led)
		go app.Run()

		// create a channel to respond to signals
		signals := make(chan os.Signal, 1)
		defer close(signals)

		signal.Notify(signals, syscall.SIGTERM)
		signal.Notify(signals, syscall.SIGINT)
		s := <-signals
		//stop <- true
		klog.Infof("Exiting, got signal: %v", s)
	},
}

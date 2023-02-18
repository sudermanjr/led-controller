package cmd

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"

	"github.com/sudermanjr/led-controller/pkg/dashboard"
	"github.com/sudermanjr/led-controller/pkg/homekit"
	"github.com/sudermanjr/led-controller/pkg/screen"
)

var (
	homekitPin     string
	screenAttached bool
	app            = dashboard.App{}
)

func init() {
	rootCmd.AddCommand(dashboardCmd)
	dashboardCmd.PersistentFlags().IntVarP(&app.Port, "port", "p", 8080, "The port to serve the dashboard on.")
	dashboardCmd.PersistentFlags().StringVar(&homekitPin, "homekit-pin", "29847290", "The pin that homekit will use to authenticate with this device.")
	dashboardCmd.PersistentFlags().BoolVar(&screenAttached, "screen", false, "Set to true if you have a screen attached.")
}

var dashboardCmd = &cobra.Command{
	Use:   "dashboard",
	Short: "Run a dashboard",
	Long:  `Run a dashboard`,
	Run: func(cmd *cobra.Command, args []string) {
		if screenAttached {
			display, err := screen.NewDisplay()
			if err != nil {
				app.Logger.Fatalw("failed to initialize screen", "error", err)
			}
			app.Screen = display
		}

		app.Initialize()
		go homekit.Start(homekitPin, app.Array)
		go app.Run()

		// create a channel to respond to signals
		signals := make(chan os.Signal, 1)
		defer close(signals)

		signal.Notify(signals, syscall.SIGTERM)
		signal.Notify(signals, syscall.SIGINT)
		s := <-signals
		//stop <- true
		app.Logger.Infow("got signal, exiting", "signal", s)
	},
}

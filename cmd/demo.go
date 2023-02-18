package cmd

import (
	"github.com/spf13/cobra"

	"github.com/sudermanjr/led-controller/pkg/screen"
)

var (
	demoBrightness     int
	demoDelay          int
	demoCount          int
	demoGradientLength int
)

func init() {
	rootCmd.AddCommand(demoCmd)
	rootCmd.AddCommand(screenDemoCmd)

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

		defer app.Array.WS.Fini()

		app.Array.Brightness = demoBrightness
		app.Array.Demo(demoCount, demoDelay, demoGradientLength)
	},
}

var screenDemoCmd = &cobra.Command{
	Use:   "screen-demo",
	Short: "Run a demo of the lcd screen.",
	Long:  `Runs a demo of the lcd screen.`,
	Run: func(cmd *cobra.Command, args []string) {

		display, err := screen.NewDisplay(app.Logger)
		if err != nil {
			app.Logger.Fatalw("failed to initialize display", "error", err)
		}
		err = display.Demo()
		if err != nil {
			app.Logger.Fatalw("failed to run display demo", "error", err)
		}
	},
}

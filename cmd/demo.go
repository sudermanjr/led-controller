package cmd

import (
	"fmt"
	"os"

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

	demoCmd.Flags().IntVar(&demoDelay, "delay", 100, "The delay in ms of the demo program.")
	demoCmd.Flags().IntVar(&demoCount, "count", 1, "The number of loops to run the demo.")
	demoCmd.Flags().IntVar(&demoBrightness, "brightness", 150, "The brightness to run the demo at. Must be between min and max.")
	demoCmd.Flags().IntVar(&demoGradientLength, "gradient-count", 2048, "The number of steps in the gradient.")

	demoCmd.AddCommand(screenDemoCmd)
	demoCmd.AddCommand(ledDemoCmd)
}

var demoCmd = &cobra.Command{
	Use:   "demo",
	Short: "Commands for running demos.",
	Long:  `Commands for running demos.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println("You must pass a sub-command. Run with --help for help.")
			os.Exit(1)
		}
	},
}

var ledDemoCmd = &cobra.Command{
	Use:   "led",
	Short: "Run a demo of the LED array",
	Run: func(cmd *cobra.Command, args []string) {
		defer app.Array.WS.Fini()
		app.Array.Brightness = demoBrightness
		app.Array.Demo(demoCount, demoDelay, demoGradientLength)
	},
}

var screenDemoCmd = &cobra.Command{
	Use:   "screen",
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

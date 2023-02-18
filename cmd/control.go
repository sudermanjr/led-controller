package cmd

import (
	"github.com/spf13/cobra"

	"github.com/sudermanjr/led-controller/pkg/color"
	"github.com/sudermanjr/led-controller/pkg/screen"
)

var (
	onBrightness   int
	displayTextX   int
	displayTextY   int
	displayMessage string
	scrollMessage  bool
)

func init() {
	rootCmd.AddCommand(onCmd)
	rootCmd.AddCommand(offCmd)
	onCmd.Flags().StringVarP(&colorName, "color", "c", "white", "The color to turn the lights on to.")
	onCmd.Flags().IntVar(&onBrightness, "brightness", 100, "The brightness setting. Range is between the min-brightness and max-brightness.")

	rootCmd.AddCommand(displayText)
	displayText.Flags().StringVarP(&displayMessage, "message", "m", "LED-Controller", "The text to display.")
	displayText.Flags().IntVarP(&displayTextX, "x-coordinate", "x", 0, "The x-coordinate of the text")
	displayText.Flags().IntVarP(&displayTextY, "y-coordinate", "y", 0, "The y-coordinate of the text")
	displayText.Flags().BoolVar(&scrollMessage, "scroll", false, "If true, the message will scroll.")
}

var onCmd = &cobra.Command{
	Use:   "on",
	Short: "Turn on the lights.",
	Long:  `Turns on the lights to a specific color.`,
	Run: func(cmd *cobra.Command, args []string) {
		defer app.Array.WS.Fini()
		app.Array.Color = color.HexToColor(color.ColorMap[colorName])
		app.Array.Brightness = onBrightness
		err := app.Array.Display(0)
		if err != nil {
			app.Logger.Fatalw("failed to turn on lights", "error", err)
		}
	},
}

var offCmd = &cobra.Command{
	Use:   "off",
	Short: "Turn off the lights.",
	Long:  `Turns off the lights.`,
	Run: func(cmd *cobra.Command, args []string) {
		defer app.Array.WS.Fini()
		err := app.Array.SetMinBrightness()
		if err != nil {
			app.Logger.Fatalw("failed to set brightness", "error", err)
		}
	},
}

var displayText = &cobra.Command{
	Use:   "display-text",
	Short: "Display text.",
	Long:  `Displays text at some coordinates on the LCD screen.`,
	Run: func(cmd *cobra.Command, args []string) {
		display, err := screen.NewDisplay(app.Logger)
		if err != nil {
			app.Logger.Fatalw("could not init display", "error", err)
		}

		if scrollMessage {
			err := display.ScrollText(displayTextX, displayTextY, displayMessage)
			if err != nil {
				app.Logger.Fatalw("could not display scrolling message", "error", err, "scrollMessage", scrollMessage)
			}
		} else {
			err := display.Text(displayTextX, displayTextY, displayMessage)
			if err != nil {
				app.Logger.Fatalw("could not display text", "error", err, "displayMessage", displayMessage)
			}
		}
	},
}

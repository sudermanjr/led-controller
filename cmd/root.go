package cmd

import (
	"flag"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"k8s.io/klog"
)

var (
	version         = "development"
	commit          = "n/a"
	ledCount        int
	maxBrightness   int
	minBrightness   int
	fadeDuration    int
	colorName       string
	displayRSPin    int
	displayEPin     int
	lineSize        int
	displayDataPins []int
)

func init() {
	// Flags
	rootCmd.PersistentFlags().IntVarP(&ledCount, "led-count", "l", 12, "The number of LEDs in the array.")
	rootCmd.PersistentFlags().IntVar(&maxBrightness, "max-brightness", 200, "The maximum brightness that will work within the 0-250 range.")
	rootCmd.PersistentFlags().IntVar(&minBrightness, "min-brightness", 25, "The minimum brightness that will work within the 0-250 range.")
	rootCmd.PersistentFlags().IntVarP(&fadeDuration, "fade-duration", "f", 100, "The duration of fade-ins and fade-outs in ms.")

	// LCD Screen Flags
	rootCmd.PersistentFlags().IntVar(&displayRSPin, "rs-pin", 25, "The GPIO number connected to the RS Pin on the LCD display.")
	rootCmd.PersistentFlags().IntVar(&displayEPin, "e-pin", 24, "The GPIO number connected to the E pin on the LCD display.")
	rootCmd.PersistentFlags().IntVar(&lineSize, "line-size", 16, "The line size of the LCD display.")
	rootCmd.PersistentFlags().IntSliceVar(&displayDataPins, "data-pins", []int{23, 17, 18, 22}, "The data pins connected to the LCD")

	//Commands
	rootCmd.AddCommand(versionCmd)

	klog.InitFlags(nil)
	pflag.CommandLine.AddGoFlag(flag.CommandLine.Lookup("v"))

	environmentVariables := map[string]string{
		"LED_COUNT":      "led-count",
		"MAX_BRIGHTNESS": "max-brightness",
		"MIN_BRIGHTNESS": "min-brightness",
		"FADE_DURATION":  "fade-duration",
	}

	for env, flag := range environmentVariables {
		flag := rootCmd.PersistentFlags().Lookup(flag)
		flag.Usage = fmt.Sprintf("%v [%v]", flag.Usage, env)
		if value := os.Getenv(env); value != "" {
			err := flag.Value.Set(value)
			if err != nil {
				klog.Errorf("Error setting flag %v to %s from environment variable %s", flag, value, env)
			}
		}
	}
}

var rootCmd = &cobra.Command{
	Use:   "led-controller",
	Short: "led-controller",
	Long:  `A cli for running neopixels`,
	Run: func(cmd *cobra.Command, args []string) {
		klog.Error("You must specify a sub-command.")
		err := cmd.Help()
		if err != nil {
			klog.Error(err)
		}
		os.Exit(1)
	},
}

// Execute the stuff
func Execute(VERSION string, COMMIT string) {
	version = VERSION
	commit = COMMIT
	if err := rootCmd.Execute(); err != nil {
		klog.Error(err)
		os.Exit(1)
	}
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Prints the current version of the tool.",
	Long:  `Prints the current version.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Version:" + version + " Commit:" + commit)
	},
}

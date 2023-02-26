package cmd

import (
	"flag"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/sudermanjr/led-controller/pkg/neopixel"
	"go.uber.org/zap"
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
	rootCmd.PersistentFlags().IntVar(&maxBrightness, "max-brightness", 250, "The maximum brightness that will work within the 0-250 range.")
	rootCmd.PersistentFlags().IntVar(&minBrightness, "min-brightness", 25, "The minimum brightness that will work within the 0-250 range.")
	rootCmd.PersistentFlags().IntVarP(&fadeDuration, "fade-duration", "f", 100, "The duration of fade-ins and fade-outs in ms.")

	// LCD Screen Flags
	rootCmd.PersistentFlags().IntVar(&displayRSPin, "rs-pin", 25, "The GPIO number connected to the RS Pin on the LCD display.")
	rootCmd.PersistentFlags().IntVar(&displayEPin, "e-pin", 24, "The GPIO number connected to the E pin on the LCD display.")
	rootCmd.PersistentFlags().IntVar(&lineSize, "line-size", 16, "The line size of the LCD display.")
	rootCmd.PersistentFlags().IntSliceVar(&displayDataPins, "data-pins", []int{23, 17, 18, 22}, "The data pins connected to the LCD")

	//Commands
	rootCmd.AddCommand(versionCmd)

	// initialize logging flags
	logLevel = zap.LevelFlag("v", zap.InfoLevel, "log level")
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
				app.Logger.Errorw("Error setting flag from environment variable", "flag", flag, "value", value, "envVar", env)
			}
		}
	}
}

var rootCmd = &cobra.Command{
	Use:   "led-controller",
	Short: "led-controller",
	Long:  `A cli for running neopixels`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		logConfig := &zap.Config{
			Encoding:         "json",
			EncoderConfig:    encoderConfig,
			OutputPaths:      []string{"stdout"},
			ErrorOutputPaths: []string{"stderr"},
			Level:            zap.NewAtomicLevelAt(*logLevel),
		}

		l, err := logConfig.Build(zap.AddStacktrace(zap.DPanicLevel))
		if err != nil {
			return err
		}
		app.Logger = l.Sugar()

		app.Array, err = neopixel.NewLEDArray(minBrightness, maxBrightness, ledCount, fadeDuration, app.Logger)
		if err != nil {
			return err
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		app.Logger.Errorw("you must specify a sub-command")
		err := cmd.Help()
		if err != nil {
			app.Logger.Errorw("error displaying help", "error", err)
		}
		os.Exit(1)
	},
}

// Execute the stuff
func Execute(VERSION string, COMMIT string) {
	version = VERSION
	commit = COMMIT
	if err := rootCmd.Execute(); err != nil {
		app.Logger.Fatalw("failed to execute cobra cmd", "error", err)
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

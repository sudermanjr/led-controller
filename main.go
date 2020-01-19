package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"k8s.io/klog"
)

func main() {
	Execute(version, commit)
}

var (
	version       = "development"
	commit        = "n/a"
	ledCount      int
	maxBrightness int
	minBrightness int
	fadeDuration  int
	colorName     string
)

func init() {
	// Flags
	rootCmd.PersistentFlags().IntVarP(&ledCount, "led-count", "l", 12, "The number of LEDs in the array.")
	rootCmd.PersistentFlags().IntVar(&maxBrightness, "max-brightness", 200, "The maximum brightness that will work within the 0-250 range.")
	rootCmd.PersistentFlags().IntVar(&minBrightness, "min-brightness", 30, "The minimum brightness that will work within the 0-250 range.")

	rootCmd.PersistentFlags().IntVarP(&fadeDuration, "fade-duration", "f", 30, "The duration of fade-ins and fade-outs in ms.")

	//Commands
	rootCmd.AddCommand(versionCmd)

	klog.InitFlags(nil)
	flag.Parse()
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)

	environmentVariables := map[string]string{
		"LED_COUNT": "led-count",
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

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
	version    = "development"
	commit     = "n/a"
	ledCount   int
	brightness int
	demoDelay  int
	demoCount  int
)

func init() {
	// Flags
	rootCmd.PersistentFlags().IntVar(&ledCount, "led-count", 12, "The number of LEDs in the array.")
	rootCmd.PersistentFlags().IntVar(&brightness, "brightness", 100, "The brightnes to run the LEDs at.")

	//Commands
	rootCmd.AddCommand(demo)
	rootCmd.AddCommand(versionCmd)

	demo.Flags().IntVar(&demoDelay, "speed", 200, "The delay in ms of the demo program.")
	demo.Flags().IntVar(&demoCount, "count", 2, "The number of loops to run the demo.")

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

var demo = &cobra.Command{
	Use:   "demo",
	Short: "Run a demo.",
	Long:  `Runs a demo.`,
	Run: func(cmd *cobra.Command, args []string) {
		Demo()
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

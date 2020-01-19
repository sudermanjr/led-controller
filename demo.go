package main

import (
	"github.com/spf13/cobra"
	"k8s.io/klog"
)

func init() {
	rootCmd.AddCommand(demoCmd)

	demoCmd.Flags().IntVar(&demoDelay, "speed", 100, "The delay in ms of the demo program.")
	demoCmd.Flags().IntVar(&demoCount, "count", 1, "The number of loops to run the demo.")
}

var demoCmd = &cobra.Command{
	Use:   "demo",
	Short: "Run a demo.",
	Long:  `Runs a demo.`,
	Run: func(cmd *cobra.Command, args []string) {

		led, err := newLEDArray()
		if err != nil {
			klog.Fatal(err)
		}
		defer led.ws.Fini()

		for i := 1; i < (demoCount + 1); i++ {
			for colorName, color := range colors {
				klog.Infof("displaying: %s", colorName)
				_ = led.display(color, demoDelay, 150)
			}
		}

		_ = led.display(off, 0, 0)
	},
}

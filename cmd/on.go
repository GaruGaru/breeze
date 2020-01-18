package cmd

import (
	"breeze/fan"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(onCmd)
}

var onCmd = &cobra.Command{
	Use:   "on",
	Short: "turn on specified gpio pin",
	Run: func(cmd *cobra.Command, args []string) {
		fanController, err := fan.NewGpio(gpioPin)
		defer func() {
			if err := fanController.Close(); err != nil {
				panic(err)
			}
		}()
		if err != nil {
			panic(err)
		}

		fanController.On()
	},
}

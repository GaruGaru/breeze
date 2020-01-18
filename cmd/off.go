package cmd

import (
	"breeze/fan"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(offCmd)
}

var offCmd = &cobra.Command{
	Use:   "off",
	Short: "turn off specified gpio pin",
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

		fanController.Off()
	},
}

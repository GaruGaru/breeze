package cmd

import (
	"breeze/fan"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
	"syscall"
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

		if err := fanController.Off(); err != nil {
			panic(err)
		}

		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		defer signal.Stop(quit)

		for {
			select {
			default:
				if err := fanController.Off(); err != nil {
					panic(err)
				}
			case <-quit:
				if err := fanController.Close(); err != nil {
					panic(err)
				}
			}
		}

	},
}

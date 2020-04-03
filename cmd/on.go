package cmd

import (
	"breeze/fan"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
	"syscall"
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

		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		defer signal.Stop(quit)

		for {
			select {
			default:
				if err := fanController.On(); err != nil {
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

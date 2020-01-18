package cmd

import (
	"breeze/controller"
	"breeze/fan"
	"breeze/sensor"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"time"
)

var (
	targetTemperature          float64
	cooldownTemperaturePercent float64
	checkLatency               time.Duration
)

func init() {
	controllerCmd.Flags().DurationVarP(&checkLatency, "check-every", "d", 5*time.Second, "controller loop delay")
	controllerCmd.Flags().Float64VarP(&targetTemperature, "target-temperature", "t", 65, "target temperature in celsius")
	controllerCmd.Flags().Float64VarP(&cooldownTemperaturePercent, "temperature-cooldown-percent", "c", 15, "temperature percentage for cool down (relative to target-temperature)")
	rootCmd.AddCommand(controllerCmd)
}

var controllerCmd = &cobra.Command{
	Use:   "controller",
	Short: "run breeze smart fan controller based on temperature",
	Run: func(cmd *cobra.Command, args []string) {
		logrus.WithFields(logrus.Fields{
			"pin":              gpioPin,
			"target":           targetTemperature,
			"cooldown_percent": cooldownTemperaturePercent,
		}).Info("breeze controller started")

		fanController, err := fan.NewGpio(gpioPin)
		defer func() {
			if err := fanController.Close(); err != nil {
				panic(err)
			}
		}()
		if err != nil {
			panic(err)
		}

		tempSensor := sensor.Cpu{}

		contr := controller.New(controller.Config{
			Delay:           checkLatency,
			Threshold:       targetTemperature,
			CoolDownPercent: cooldownTemperaturePercent,
		})

		if err := contr.Run(fanController, tempSensor); err != nil {
			panic(err)
		}
	},
}

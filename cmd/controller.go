package cmd

import (
	"breeze/controller"
	"breeze/fan"
	"breeze/metrics"
	"breeze/sensor"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"time"
)

var (
	targetTemperature          float64
	cooldownTemperaturePercent float64
	checkLatency               time.Duration
	metricsEnabled             bool
	metricsServerPort          int
	metricsServerAddr          string
)

func init() {
	controllerCmd.Flags().DurationVarP(&checkLatency, "check-every", "d", 5*time.Second, "controller loop delay")
	controllerCmd.Flags().Float64VarP(&targetTemperature, "target-temperature", "t", 65, "target temperature in celsius")
	controllerCmd.Flags().Float64VarP(&cooldownTemperaturePercent, "temperature-cooldown-percent", "c", 15, "temperature percentage for cool down (relative to target-temperature)")
	controllerCmd.Flags().BoolVar(&metricsEnabled, "metrics", false, "expose metrics endpoint for prometheus")
	controllerCmd.Flags().StringVar(&metricsServerAddr, "metrics-addr", "0.0.0.0", "metrics server bind address")
	controllerCmd.Flags().IntVar(&metricsServerPort, "metrics-port", 9999, "metrics server port ")

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

		if metricsEnabled {
			go func() {
				metricsServer := metrics.New(tempSensor)
				if err := metricsServer.Run(metricsServerAddr, metricsServerPort); err != nil {
					logrus.Error("error running metrics server", err)
				}
			}()
		}

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

package cmd

import (
	"breeze/controller"
	"breeze/fan"
	"breeze/metrics"
	"breeze/sensor"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	targetTemperature          float64
	cooldownTemperaturePercent float64
	checkLatency               time.Duration
	metricsEnabled             bool
	metricsServerPort          int
	metricsServerAddr          string
	nodeName                   string
	quietHoursRaw              string
)

func init() {
	controllerCmd.Flags().DurationVarP(&checkLatency, "check-every", "d", 5*time.Second, "controller loop delay")
	controllerCmd.Flags().Float64VarP(&targetTemperature, "target-temperature", "t", 65, "target temperature in celsius")
	controllerCmd.Flags().Float64VarP(&cooldownTemperaturePercent, "temperature-cooldown-percent", "c", 15, "temperature percentage for cool down (relative to target-temperature)")
	controllerCmd.Flags().BoolVar(&metricsEnabled, "metrics", false, "expose metrics endpoint for prometheus")
	controllerCmd.Flags().StringVar(&metricsServerAddr, "metrics-addr", "0.0.0.0", "metrics server bind address")
	controllerCmd.Flags().IntVar(&metricsServerPort, "metrics-port", 9999, "metrics server port ")
	controllerCmd.Flags().StringVar(&nodeName, "node-name", envOrDefault("NODE_NAME", "notset"), "metrics node name label")
	controllerCmd.Flags().StringVar(&quietHoursRaw, "quiet-hours", "", "quiet hours in 24hr format eg: 22-7")

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
				metricsServer := metrics.New(nodeName, tempSensor)
				if err := metricsServer.Run(metricsServerAddr, metricsServerPort); err != nil {
					logrus.Error("error running metrics server", err)
				}
			}()
		}

		quietHours, err := parseQuietHours(quietHoursRaw)
		if err != nil {
			log.Fatal(err)
		}

		contr := controller.New(controller.Config{
			Threshold:       targetTemperature,
			CoolDownPercent: cooldownTemperaturePercent,
			Delay:           checkLatency,
			QuietHours:      quietHours,
		})

		if err := contr.Run(fanController, tempSensor); err != nil {
			panic(err)
		}
	},
}

func parseQuietHours(raw string) ([2]int, error) {
	if len(raw) == 0 {
		return [2]int{}, nil
	}
	rawParts := strings.Split(raw, "-")
	if len(rawParts) != 2 {
		return [2]int{}, fmt.Errorf("invalid quiet hours value: %s", raw)
	}

	valBegin, err := strconv.Atoi(rawParts[0])
	if err != nil {
		return [2]int{}, err
	}

	if valBegin < 0 || valBegin > 23 {
		return [2]int{}, fmt.Errorf("invalid quiet hours value: %d must be in [0-23]", valBegin)
	}

	valEnd, err := strconv.Atoi(rawParts[0])
	if err != nil {
		return [2]int{}, err
	}

	if valEnd < 0 || valEnd > 23 {
		return [2]int{}, fmt.Errorf("invalid quiet hours value: %d must be in [0-23]", valEnd)
	}

	return [2]int{valBegin, valEnd}, nil
}

func envOrDefault(key string, def string) string {
	val, found := os.LookupEnv(key)
	if found {
		return val
	}
	return def
}

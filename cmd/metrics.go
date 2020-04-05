package cmd

import (
	"breeze/metrics"
	"breeze/sensor"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	metricsCmd.Flags().StringVar(&metricsServerAddr, "metrics-addr", "0.0.0.0", "metrics server bind address")
	metricsCmd.Flags().IntVar(&metricsServerPort, "metrics-port", 9999, "metrics server port ")
	rootCmd.AddCommand(metricsCmd)
}

var metricsCmd = &cobra.Command{
	Use:   "metrics",
	Short: "run metrics server",
	Run: func(cmd *cobra.Command, args []string) {
		tempSensor := sensor.Cpu{}
		metricsServer := metrics.New(tempSensor)
		if err := metricsServer.Run(metricsServerAddr, metricsServerPort); err != nil {
			logrus.Error("error running metrics server", err)
		}
	},
}

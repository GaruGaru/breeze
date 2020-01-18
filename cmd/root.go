package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var (
	gpioPin int
)

func init() {
	rootCmd.PersistentFlags().IntVarP(&gpioPin, "gpio-pin", "p", 14, "gpio fan controller pin number")
}

var rootCmd = &cobra.Command{
	Use:   "breeze",
	Short: "Breeze is smart a gpio fan controller",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

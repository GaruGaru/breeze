package main

import (
	"breeze/controller"
	"breeze/fan"
	"breeze/sensor"
	"time"
)

func main() {

	// TODO configuration from flag/env
	fanController, err := fan.NewGpio(14)
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
		Delay: 3 * time.Second,
		Threshold: 65,
	})

	if err := contr.Run(fanController, tempSensor); err != nil {
		panic(err)
	}
}

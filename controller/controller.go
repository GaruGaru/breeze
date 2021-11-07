package controller

import (
	"breeze/fan"
	"breeze/sensor"
	log "github.com/sirupsen/logrus"
	"time"
)

const (
	DefaultDelay              = 10 * time.Second
	DefaultThreshold          = 65.0
	DefaultCoolDownPercentage = 15.0
)

type Config struct {
	Threshold       float64
	CoolDownPercent float64
	Delay           time.Duration
	QuietHours      [2]int
}

type status struct {
	targetTemperature *float64
	coolingDown       bool
	beginCoolDown     *time.Time
}

type Controller struct {
	config Config
	status status
}

func New(config Config) *Controller {
	if config.Delay == 0 {
		config.Delay = DefaultDelay
	}

	if config.Threshold == 0 {
		config.Threshold = DefaultThreshold
	}

	if config.CoolDownPercent == 0 {
		config.CoolDownPercent = DefaultCoolDownPercentage
	}
	return &Controller{config: config}
}

func (c *Controller) Run(fan fan.Controller, sensor sensor.Thermal) error {
	for {
		if c.isInQuietHours(time.Now()) {
			log.Info("trying to run but was in quiet hours")
			time.Sleep(c.config.Delay)
			continue
		}

		temp, err := sensor.Read()
		if err != nil {
			return err
		}

		log.Debugf("temperature %d°", int(temp))

		if !c.status.coolingDown && temp > c.config.Threshold {
			now := time.Now()
			c.status.coolingDown = true
			c.status.beginCoolDown = &now
			c.status.targetTemperature = float64Ptr(c.config.Threshold - percent(c.config.CoolDownPercent, c.config.Threshold))

			log.Infof("temperature threshold reached %f°, cooling down until %f°", temp, *c.status.targetTemperature)
		}

		if c.status.coolingDown && temp <= *c.status.targetTemperature {
			coolingTime := time.Since(*c.status.beginCoolDown)
			log.Infof("target temperature (%f°) reached in %d seconds", *c.status.targetTemperature, int(coolingTime.Seconds()))

			c.status.coolingDown = false
			c.status.beginCoolDown = nil
			c.status.targetTemperature = nil
		}

		if c.status.coolingDown {
			if err := fan.On(); err != nil {
				log.Errorf("error turning on fan: %s", err.Error())
			}
		} else {
			if err := fan.Off(); err != nil {
				log.Errorf("error turning off fan: %s", err.Error())
			}
		}

		time.Sleep(c.config.Delay)
	}
}

func (c *Controller) isInQuietHours(t time.Time) bool {
	if c.config.QuietHours == [2]int{} {
		return false
	}
	hr := t.Hour()

	//cross day
	if c.config.QuietHours[0] > c.config.QuietHours[1] {
		return hr >= c.config.QuietHours[0] || hr <= c.config.QuietHours[1]
	}

	return hr >= c.config.QuietHours[0] && hr <= c.config.QuietHours[1]
}

func float64Ptr(val float64) *float64 {
	return &val
}

func percent(p, val float64) float64 {
	return (p * val) / 100
}

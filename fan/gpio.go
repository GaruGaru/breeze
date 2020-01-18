package fan

import "github.com/warthog618/gpio"

type Gpio struct {
	pinNum int
	pin    *gpio.Pin
}

func NewGpio(pinNum int) (*Gpio, error) {
	controller := &Gpio{pinNum: pinNum}
	if err := controller.Init(); err != nil {
		return nil, err
	}
	return controller, nil
}

func (g *Gpio) Init() error {
	if err := gpio.Open(); err != nil {
		return err
	}

	g.pin = gpio.NewPin(g.pinNum)
	return nil
}

func (g *Gpio) On() {
	if g.pin == nil {
		panic("gpio controller not initialized")
	}
	g.pin.High()
}

func (g *Gpio) Off() {
	if g.pin == nil {
		panic("gpio controller not initialized")
	}
	g.pin.Low()
}

func (g *Gpio) Close() error {
	return gpio.Close()
}

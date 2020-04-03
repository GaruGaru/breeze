package fan

import (
	"github.com/hashicorp/go-multierror"
	"github.com/warthog618/gpiod"
)

const (
	consumerName = "breeze"
	gpioChip = "gpiochip0"
)
type Gpio struct {
	pinNum int
	chip   *gpiod.Chip
	line   *gpiod.Line
}

func NewGpio(pinNum int) (*Gpio, error) {
	chip, err := gpiod.NewChip(gpioChip, gpiod.WithConsumer(consumerName))
	if err != nil {
		return nil, err
	}

	line, err := chip.RequestLine(pinNum, gpiod.AsOutput(1))
	if err != nil {
		return nil, err
	}

	return &Gpio{
		pinNum: pinNum,
		chip:   chip,
		line:   line,
	}, nil
}

func (g *Gpio) On() error {
	return g.line.SetValue(1)
}

func (g *Gpio) Off() error {
	return g.line.SetValue(0)
}

func (g *Gpio) Close() error {
	var errors error
	if err := g.line.Close(); err != nil {
		errors = multierror.Append(errors, err)
	}
	if err := g.chip.Close(); err != nil {
		errors = multierror.Append(errors, err)
	}
	return errors
}

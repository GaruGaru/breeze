package sensor

import (
	"bytes"
	"io/ioutil"
	"strconv"
)

const cpuThermalPath = "/sys/class/thermal/thermal_zone0/temp"

type Cpu struct{}

func (c Cpu) Read() (float64, error) {
	data, err := ioutil.ReadFile(cpuThermalPath)
	if err != nil {
		return 0, err
	}

	data = bytes.ReplaceAll(data, []byte("\n"), []byte{})

	value, err := strconv.ParseInt(string(data), 10, 32)
	if err != nil {
		return 0, err
	}

	return float64(value) / 1000, nil
}

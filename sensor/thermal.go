package sensor

type Thermal interface {
	Read() (float64, error)
}

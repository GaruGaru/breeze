package fan

type Controller interface {
	Init() error
	On()
	Off()
	Close() error
}

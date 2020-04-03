package fan

type Controller interface {
	On() error
	Off() error
	Close() error
}

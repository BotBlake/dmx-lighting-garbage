package dmx

type Driver interface {
	Open() error
	Close() error
	SendFrame(universe int, frame []byte) error
	Info() string
}

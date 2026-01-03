package dmx

const (
	Channels    = 512
	StartCode   = 0x00
	FrameLength = Channels + 1
)

type Universe struct {
	ID    int
	Frame [FrameLength]byte
}

func NewUniverse(id int) *Universe {
	u := &Universe{ID: id}
	u.Frame[0] = StartCode
	return u
}

package dmx

type CommandType int

const (
	SetChannel CommandType = iota
	ClearUniverse
)

type Command struct {
	Type       CommandType
	UniverseID int
	Channel    int
	Value      byte
}

package device

import (
	"dmx-lighting-garbage/internal/dmx"
	"dmx-lighting-garbage/internal/util"
	"strings"
)

type Instance struct {
	Name     string
	Profile  *DDFDevice
	Universe *dmx.Universe
	Address  int // 1-based start address
}

// Resolve absolute DMX channel
func (i *Instance) dmx(ch int) int {
	return i.Address + ch
}

func (i *Instance) SetBrightness(percent int) {
	if i.Profile.Functions.Dimmer == nil {
		return
	}

	ch := i.Profile.Functions.Dimmer.DMXChannel
	i.Universe.Frame[i.dmx(ch)] = util.PercentToDMX(percent)
}

func (i *Instance) whiteChannel() *int {
	for _, r := range i.Profile.Functions.RawList {
		if strings.EqualFold(r.Name, "white") {
			return &r.DMXChannel
		}
	}
	return nil
}

func (i *Instance) SetColor(hex string) {
	rgb := i.Profile.Functions.RGB
	if rgb == nil {
		return
	}

	r, g, b, err := util.HexToRGB(hex)
	if err != nil {
		return
	}

	white := i.whiteChannel()

	// Pure white → prefer white channel
	if white != nil && r == 255 && g == 255 && b == 255 {
		i.Universe.Frame[i.dmx(*white)] = 255
		i.Universe.Frame[i.dmx(rgb.Red.DMXChannel)] = 0
		i.Universe.Frame[i.dmx(rgb.Green.DMXChannel)] = 0
		i.Universe.Frame[i.dmx(rgb.Blue.DMXChannel)] = 0
		return
	}

	i.Universe.Frame[i.dmx(rgb.Red.DMXChannel)] = r
	i.Universe.Frame[i.dmx(rgb.Green.DMXChannel)] = g
	i.Universe.Frame[i.dmx(rgb.Blue.DMXChannel)] = b

	if white != nil {
		i.Universe.Frame[i.dmx(*white)] = 0
	}
}

func (i *Instance) Clear() {
	if rgb := i.Profile.Functions.RGB; rgb != nil {
		i.Universe.Frame[i.dmx(rgb.Red.DMXChannel)] = 0
		i.Universe.Frame[i.dmx(rgb.Green.DMXChannel)] = 0
		i.Universe.Frame[i.dmx(rgb.Blue.DMXChannel)] = 0
	}

	if w := i.whiteChannel(); w != nil {
		i.Universe.Frame[i.dmx(*w)] = 0
	}

	if d := i.Profile.Functions.Dimmer; d != nil {
		i.Universe.Frame[i.dmx(d.DMXChannel)] = 0
	}
}

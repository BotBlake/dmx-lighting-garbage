package device

import (
	"dmx-lighting-garbage/internal/dmx"
	"dmx-lighting-garbage/internal/util"
	"strings"
)

type Instance struct {
	Name       string
	Profile    *DDFDevice
	UniverseID int
	Address    int // 1-based start address

	X float32
	Y float32
	Z float32
}

// Resolve absolute DMX channel
func (i *Instance) dmx(ch int) int {
	return i.Address + ch
}

func (i *Instance) SetBrightness(percent int) []dmx.Command {
	ch := i.Profile.Functions.Dimmer.DMXChannel
	return []dmx.Command{
		{
			Type:       dmx.SetChannel,
			UniverseID: i.UniverseID,
			Channel:    i.dmx(ch),
			Value:      util.PercentToDMX(percent),
		},
	}
}

func (i *Instance) whiteChannel() *int {
	for _, r := range i.Profile.Functions.RawList {
		if strings.EqualFold(r.Name, "white") {
			return &r.DMXChannel
		}
	}
	return nil
}

func (i *Instance) SetColor(hex string) []dmx.Command {
	rgb := i.Profile.Functions.RGB
	if rgb == nil {
		return nil
	}

	r, g, b, err := util.HexToRGB(hex)
	if err != nil {
		return nil
	}

	cmds := []dmx.Command{}
	white := i.whiteChannel()

	// Pure white → prefer white channel
	if white != nil && r == 255 && g == 255 && b == 255 {
		cmds = append(cmds,
			dmx.Command{
				Type:       dmx.SetChannel,
				UniverseID: i.UniverseID,
				Channel:    i.dmx(*white),
				Value:      255,
			},
			dmx.Command{
				Type:       dmx.SetChannel,
				UniverseID: i.UniverseID,
				Channel:    i.dmx(rgb.Red.DMXChannel),
				Value:      0,
			},
			dmx.Command{
				Type:       dmx.SetChannel,
				UniverseID: i.UniverseID,
				Channel:    i.dmx(rgb.Green.DMXChannel),
				Value:      0,
			},
			dmx.Command{
				Type:       dmx.SetChannel,
				UniverseID: i.UniverseID,
				Channel:    i.dmx(rgb.Blue.DMXChannel),
				Value:      0,
			},
		)
		return cmds
	}

	// Normal RGB output
	cmds = append(cmds,
		dmx.Command{
			Type:       dmx.SetChannel,
			UniverseID: i.UniverseID,
			Channel:    i.dmx(rgb.Red.DMXChannel),
			Value:      r,
		},
		dmx.Command{
			Type:       dmx.SetChannel,
			UniverseID: i.UniverseID,
			Channel:    i.dmx(rgb.Green.DMXChannel),
			Value:      g,
		},
		dmx.Command{
			Type:       dmx.SetChannel,
			UniverseID: i.UniverseID,
			Channel:    i.dmx(rgb.Blue.DMXChannel),
			Value:      b,
		},
	)

	// Ensure white channel is off if present
	if white != nil {
		cmds = append(cmds,
			dmx.Command{
				Type:       dmx.SetChannel,
				UniverseID: i.UniverseID,
				Channel:    i.dmx(*white),
				Value:      0,
			},
		)
	}

	return cmds
}

func (i *Instance) Clear() []dmx.Command {
	cmds := []dmx.Command{}

	// Clear RGB channels
	if rgb := i.Profile.Functions.RGB; rgb != nil {
		cmds = append(cmds,
			dmx.Command{
				Type:       dmx.SetChannel,
				UniverseID: i.UniverseID,
				Channel:    i.dmx(rgb.Red.DMXChannel),
				Value:      0,
			},
			dmx.Command{
				Type:       dmx.SetChannel,
				UniverseID: i.UniverseID,
				Channel:    i.dmx(rgb.Green.DMXChannel),
				Value:      0,
			},
			dmx.Command{
				Type:       dmx.SetChannel,
				UniverseID: i.UniverseID,
				Channel:    i.dmx(rgb.Blue.DMXChannel),
				Value:      0,
			},
		)
	}

	// Clear white channel
	if w := i.whiteChannel(); w != nil {
		cmds = append(cmds,
			dmx.Command{
				Type:       dmx.SetChannel,
				UniverseID: i.UniverseID,
				Channel:    i.dmx(*w),
				Value:      0,
			},
		)
	}

	// Clear dimmer
	if d := i.Profile.Functions.Dimmer; d != nil {
		cmds = append(cmds,
			dmx.Command{
				Type:       dmx.SetChannel,
				UniverseID: i.UniverseID,
				Channel:    i.dmx(d.DMXChannel),
				Value:      0,
			},
		)
	}

	return cmds
}

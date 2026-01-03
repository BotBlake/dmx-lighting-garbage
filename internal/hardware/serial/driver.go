package serialhw

import (
	"time"

	"dmx-lighting-garbage/internal/dmx"

	"go.bug.st/serial"
)

var _ dmx.Driver = (*Driver)(nil)

type breakCapable interface {
	Break(time.Duration) error
}

type Driver struct {
	portName string
	port     serial.Port
	breaker  breakCapable
}

func New(port string) *Driver {
	return &Driver{portName: port}
}

func (d *Driver) Open() error {
	mode := &serial.Mode{
		BaudRate: 250000,
		DataBits: 8,
		Parity:   serial.NoParity,
		StopBits: serial.TwoStopBits,
	}

	p, err := serial.Open(d.portName, mode)
	if err != nil {
		return err
	}

	d.port = p
	if b, ok := p.(breakCapable); ok {
		d.breaker = b
		print("Using native BREAK support\n")
	}
	return nil
}

func (d *Driver) Close() error {
	if d.port != nil {
		return d.port.Close()
	}
	return nil
}

func (d *Driver) SendFrame(_ int, frame []byte) error {
	if d.breaker != nil {
		// Use native BREAK support if available
		d.breaker.Break(120 * time.Microsecond)
	} else {
		// Fallback DMX BREAK via baud rate drop (HORRIBLE) ---

		// 1) Switch to slow baud to force framing error
		d.port.SetMode(&serial.Mode{
			BaudRate: 9600,
			DataBits: 8,
			Parity:   serial.NoParity,
			StopBits: serial.OneStopBit,
		})
		// 2) Send a zero byte -> line held low long enough
		d.port.Write([]byte{0x00})

		// 3) Wait ≥ 88µs (DMX spec)
		time.Sleep(120 * time.Microsecond)
		// 4) Restore DMX UART settings
		d.port.SetMode(&serial.Mode{
			BaudRate: 250000,
			DataBits: 8,
			Parity:   serial.NoParity,
			StopBits: serial.TwoStopBits,
		})
	}

	// --- MAB ---
	time.Sleep(12 * time.Microsecond)

	// --- DMX FRAME ---
	_, err := d.port.Write(frame)
	return err
}

func (d *Driver) Info() string {
	return "Serial DMX (" + d.portName + ")"
}

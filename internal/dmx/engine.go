package dmx

import (
	"time"
)

type Engine struct {
	bindings []Binding
	rate     time.Duration
	stop     chan struct{}
}

func NewEngine(bindings []Binding, fps int) *Engine {
	return &Engine{
		bindings: bindings,
		rate:     time.Second / time.Duration(fps),
		stop:     make(chan struct{}),
	}
}

func (e *Engine) Run(universes map[int]*Universe) {
	ticker := time.NewTicker(e.rate)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			for _, b := range e.bindings {
				u, ok := universes[b.Universe]
				if !ok {
					continue
				}
				_ = b.Driver.SendFrame(b.Universe, u.Frame[:])
			}
		case <-e.stop:
			return
		}
	}
}

func (e *Engine) Stop() {
	close(e.stop)
}

package dmx

import (
	"fmt"
	"sync"
	"time"
)

type Engine struct {
	mu sync.RWMutex

	fps  int
	rate time.Duration

	bindings  map[int]Binding
	universes map[int]*Universe

	commands chan Command

	stop    chan struct{}
	running bool
}

func NewEngine(fps int) *Engine {
	return &Engine{
		fps:       fps,
		rate:      time.Second / time.Duration(fps),
		bindings:  make(map[int]Binding),
		universes: make(map[int]*Universe),
		commands:  make(chan Command, 1024), // buffered
		stop:      make(chan struct{}),
	}
}

func (e *Engine) AddUniverse(id int, u *Universe) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if _, exists := e.universes[id]; exists {
		return fmt.Errorf("universe %d already exists", id)
	}

	e.universes[id] = u
	return nil
}

func (e *Engine) AddBinding(b Binding) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if _, exists := e.bindings[b.Universe.ID]; exists {
		return fmt.Errorf("binding already exists for universe %d", b.Universe.ID)
	}

	e.bindings[b.Universe.ID] = b
	return nil
}

func (e *Engine) Run() {
	ticker := time.NewTicker(e.rate)
	defer ticker.Stop()

	for {
		select {
		case cmd := <-e.commands:
			e.apply(cmd)

		case <-ticker.C:
			e.tick()

		case <-e.stop:
			return
		}
	}
}

func (e *Engine) tick() {
	e.mu.RLock()
	defer e.mu.RUnlock()

	for uid, binding := range e.bindings {
		u, ok := e.universes[uid]
		if !ok {
			panic(fmt.Sprintf(
				"engine error: universe %d has binding but no universe",
				uid,
			))
		}

		_ = binding.Driver.SendFrame(uid, u.Frame[:])
	}
}

func (e *Engine) apply(cmd Command) {
	e.mu.Lock()
	defer e.mu.Unlock()

	switch cmd.Type {
	case SetChannel:
		u, ok := e.universes[cmd.UniverseID]
		if !ok {
			panic("set channel on missing universe")
		}
		u.Frame[cmd.Channel] = cmd.Value

	case ClearUniverse:
		u, ok := e.universes[cmd.UniverseID]
		if !ok {
			panic("clear missing universe")
		}
		for i := range u.Frame {
			u.Frame[i] = 0
		}
	}
}

func (e *Engine) Queue(cmds []Command) {
	for _, cmd := range cmds {
		e.commands <- cmd
	}
}

func (e *Engine) Stop() {
	close(e.stop)
}

func (e *Engine) SetChannel(universeID int, channel int, value byte) {
	e.mu.Lock()
	defer e.mu.Unlock()
	if u, ok := e.universes[universeID]; ok {
		u.Frame[channel] = value
	}
}

func (e *Engine) GetOutputs() map[int][]byte {
	e.mu.RLock()
	defer e.mu.RUnlock()
	outputs := make(map[int][]byte)
	for uid, u := range e.universes {
		frameCopy := make([]byte, len(u.Frame))
		copy(frameCopy, u.Frame[:])
		outputs[uid] = frameCopy
	}
	return outputs
}

func (e *Engine) Clear() {
	e.mu.Lock()
	defer e.mu.Unlock()
	for _, u := range e.universes {
		for i := range u.Frame {
			u.Frame[i] = 0
		}
	}
}

package event

import "sync"

type Event string

type Emitter struct {
	channels   map[Event]chan any
	mu         sync.RWMutex
	bufferSize int
}

func NewEmitter() *Emitter {
	return &Emitter{
		channels:   make(map[Event]chan any),
		bufferSize: 10, //nolint:mnd // default buffer size
	}
}

func (e *Emitter) Emit(event Event, data any) {
	e.mu.RLock()
	channel := e.channels[event]
	e.mu.RUnlock()

	select {
	case channel <- data:
	default:
	}
}

func (e *Emitter) On(event Event) <-chan any {
	ch := make(chan any, e.bufferSize)

	e.mu.Lock()
	e.channels[event] = ch
	e.mu.Unlock()

	return ch
}

func (e *Emitter) Close() {
	e.mu.Lock()
	defer e.mu.Unlock()

	for event, channel := range e.channels {
		close(channel)
		delete(e.channels, event)
	}
}

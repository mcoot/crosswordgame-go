package events

import "sync"

// TODO: Test this in the abstract
// TODO: Realistically, can't use this without
// 		 rearchitecting game logic to be in a separate routine, at the least

type Kind string

type Event interface {
	Kind() Kind
}

type TypedEvent[T any] struct {
	kind Kind
	Data T
}

func (e TypedEvent[T]) Kind() Kind {
	return e.kind
}

type EventBus struct {
	subscribers map[Kind][]chan Event
	mutex       *sync.Mutex
}

func NewEventBus() *EventBus {
	return &EventBus{
		subscribers: make(map[Kind][]chan Event),
		mutex:       &sync.Mutex{},
	}
}

func (b *EventBus) Subscribe(kind Kind) chan Event {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	ch := make(chan Event)
	b.subscribers[kind] = append(b.subscribers[kind], ch)
	return ch
}

func (b *EventBus) Unsubscribe(kind Kind, ch chan Event) {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	subs := b.subscribers[kind]
	for i, sub := range subs {
		if sub == ch {
			b.subscribers[kind] = append(subs[:i], subs[i+1:]...)
			close(ch)
			return
		}
	}
}

func (b *EventBus) Publish(event Event) {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	for _, ch := range b.subscribers[event.Kind()] {
		ch <- event
	}
}

type EventConsumer struct {
	bus   *EventBus
	kinds []Kind
}

func NewEventConsumer(bus *EventBus, kinds ...Kind) *EventConsumer {
	return &EventConsumer{
		bus:   bus,
		kinds: kinds,
	}
}

// EventHandler takes an event, and returns true if processing should continue,
// or false if the consumer should unsubscribe
type EventHandler func(e Event) bool

// Run subscribes to the consumer's event kinds, and blocks while processing events for them
// terminates when the handler for all event kinds returns false
func (c *EventConsumer) Run(f EventHandler) {
	wg := sync.WaitGroup{}
	for _, kind := range c.kinds {
		ch := c.bus.Subscribe(kind)
		wg.Add(1)
		go func() {
			defer wg.Done()
			for event := range ch {
				if !f(event) {
					c.bus.Unsubscribe(kind, ch)
					return
				}
			}
		}()
	}
	wg.Wait()
}

type TypedEventHandler[T any] func(e TypedEvent[T]) bool

func (h TypedEventHandler[T]) ToEventHandler() EventHandler {
	return func(e Event) bool {
		if typed, ok := e.(TypedEvent[T]); ok {
			return h(typed)
		}
		return true
	}
}

func MuxEventHandlers(handlers map[Kind]EventHandler) EventHandler {
	return func(e Event) bool {
		if handler, ok := handlers[e.Kind()]; ok {
			return handler(e)
		}
		return true
	}
}

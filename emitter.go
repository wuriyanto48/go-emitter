package emitter

import (
	"context"
	"fmt"
	"sync"
)

// EventFunc the event function
type EventFunc[T any] func(T, *sync.WaitGroup)

// Emitter structure
type Emitter[T any] struct {
	events    map[string]([]chan T)
	wgH       *sync.WaitGroup
	wgE       *sync.WaitGroup
	mx        *sync.Mutex
	ctx       context.Context
}

// NewEmitter the emitter constructor
func NewEmitter[T any](ctx context.Context) *Emitter[T] {
	events := make(map[string]([]chan T))
	return &Emitter[T]{
		events:    events,
		wgH:       &sync.WaitGroup{},
		wgE:       &sync.WaitGroup{},
		mx:        &sync.Mutex{},
		ctx:       ctx,
	}
}

func (e *Emitter[T]) Emit(event string, eventArg T) {
	e.wgE.Add(1)

	go func() {
		e.mx.Lock()
		defer func() { e.mx.Unlock() }()
		defer func() { e.wgE.Done() }()

		if _, ok := e.events[event]; ok {
			var wg sync.WaitGroup
			for _, argC := range e.events[event] {
				wg.Add(1)
				go func(aC chan T, wg *sync.WaitGroup) {
					defer func() { wg.Done() }()
					
					select {
					case aC <- eventArg:
					case <-e.ctx.Done():
						return
					}
				}(argC, &wg)
			}

			wg.Wait()
		}
	}()
}

func (e *Emitter[T]) On(event string, handler EventFunc[T]) {

	e.mx.Lock()
	argC := make(chan T, 1)
	if _, ok := e.events[event]; ok {
		e.events[event] = append(e.events[event], argC)
	} else {
		e.events[event] = []chan T{argC}
	}
	e.mx.Unlock()

	go func() {
		for {
			// The try-receive operation here is to
			// try to exit the worker goroutine as
			// early as possible. Try-receive
			// optimized by the standard Go
			// compiler, so they are very efficient.
			select {
			case <-e.ctx.Done():
				fmt.Printf("handler canceled\n")
				return
			default: // add default case to avoid blocking from this try-receive operation
			}

			select {
			case arg := <-argC:
				e.wgH.Add(1)
				go handler(arg, e.wgH)
			case <-e.ctx.Done():
				fmt.Printf("handler canceled\n")
				return
			}
		}

	}()
}

func (e *Emitter[T]) Wait() {
	// wait all emitter
	e.wgE.Wait()

	// wait all handler
	e.wgH.Wait()

}

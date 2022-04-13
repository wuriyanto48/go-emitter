package emitter

import (
	"context"
	"testing"
	"sync"
)

type Arg struct {
	X uint
	Y uint
}

type Err struct {}

func (e *Err) Error() string {
	return "err"
}

func TestEmitter(t *testing.T) {

	operation := func() error {
		return &Err{}
	}

	em := NewEmitter[error](context.Background())

	em.On("error", func(err error, wg *sync.WaitGroup) {
		_, ok := err.(*Err)
		if !ok {
			t.Error("err is not *Err")
		}

		
		wg.Done()
	})

	em.Emit("error", operation())

	em.Wait()

	// -----------------------------------------------

	emCalc := NewEmitter[Arg](context.Background())

	emCalc.On("sum", func(arg Arg, wg *sync.WaitGroup) {
		res := arg.X + arg.Y
		if res != 10 {
			t.Error("res is not equal to 10")
		}

		wg.Done()
	})

	emCalc.Emit("sum", Arg{X: 5, Y: 5})

	emCalc.Wait()
}
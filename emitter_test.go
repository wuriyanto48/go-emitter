package emitter

import (
	"context"
	"testing"
	"sync"
	"time"
	"log"
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

}

func TestEmitterWithDelay(t *testing.T) {

	ctx, cancel := context.WithCancel(context.Background())
	// ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*800)
	defer func() { cancel() }()

	emCalc := NewEmitter[Arg](ctx)

	emCalc.On("sum", func(arg Arg, wg *sync.WaitGroup) {
		// simulate heavy work
		time.Sleep(3 * time.Second)

		res := arg.X + arg.Y
		if res != 11 {
			t.Error("res is not equal to 10")
		}
		log.Printf(" res : %d\n", res)

		wg.Done()
	})

	emCalc.Emit("sum", Arg{X: 5, Y: 5})

	emCalc.Wait()

}
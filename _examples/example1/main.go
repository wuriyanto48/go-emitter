package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/wuriyanto48/go-emitter"
)

type Data struct {
	X float64
	Y float64
}

func main() {
	start := time.Now()

	ctx, cancel := context.WithCancel(context.Background())
	// ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*800)
	defer func() { cancel() }()

	em1 := emitter.NewEmitter[Data](ctx)
	em2 := emitter.NewEmitter[string](ctx)

	// em1.On("succeed", func(arg Data) {
	// 	z := arg.X * arg.Y
	// 	// time.Sleep(2 * time.Second)
	// 	fmt.Println(z)
	// })

	em1.On("succeed", func(arg Data, wg *sync.WaitGroup) {
		z := arg.X * arg.Y
		time.Sleep(2 * time.Second)
		fmt.Println("-- ", z)
		wg.Done()
	})

	em2.On("progress", func(arg string, wg *sync.WaitGroup) {
		time.Sleep(2 * time.Second)
		fmt.Println("processed: ", arg)
		wg.Done()
	})

	em1.Emit("succeed", Data{X: 10, Y: 5})
	em1.Emit("succeed", Data{X: 100, Y: 5})
	em1.Emit("succeed", Data{X: 2, Y: 2})
	em1.Emit("succeed", Data{X: 10, Y: 10})

	em2.Emit("progress", "run....")
	em2.Emit("progress", "ron....")
	em2.Emit("progress", "ren....")
	// time.Sleep(2 * time.Second)

	em1.Wait()
	em2.Wait()

	elapsed := time.Since(start)
	log.Printf("took %s", elapsed)

}

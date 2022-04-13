## Go Emitter

It's similar to Node's Event Emitter(https://nodejs.org/api/events.html)

#
~ it uses generic features, so Go version 1.18 is the minimum requirement
#

### Usage

```go
type Data struct {
	X float64
	Y float64
}

func main() {
	start := time.Now()

	ctx, cancel := context.WithCancel(context.Background())
	defer func() { cancel() }()

	em1 := emitter.NewEmitter[Data](ctx)
	em2 := emitter.NewEmitter[string](ctx)

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

	em1.Wait()
	em2.Wait()

	elapsed := time.Since(start)
	log.Printf("took %s", elapsed)

}
```
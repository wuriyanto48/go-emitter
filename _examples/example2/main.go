package main

import (
	"bufio"
	"os"
	"strings"
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/wuriyanto48/go-emitter"
)

func main() {
	start := time.Now()

	ctx, cancel := context.WithCancel(context.Background())
	// ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*800)
	defer func() { cancel() }()

	em2 := emitter.NewEmitter[string](ctx)

	em2.On("progress", func(arg string, wg *sync.WaitGroup) {
		time.Sleep(2 * time.Second)
		fmt.Println("processed: ", arg)
		wg.Done()
	})

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("input: ")

		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		input = strings.TrimSpace(input)

		if input == "exit" {
			break
		}

		// fmt.Println("you typed: ", input)
		em2.Emit("progress", input)
	}

	em2.Wait()

	elapsed := time.Since(start)
	log.Printf("took %s", elapsed)

}

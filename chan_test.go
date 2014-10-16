package buffer

import (
	"fmt"
	"sync"
	"testing"
)

func TestMemChan(t *testing.T) {
	input := make(chan []byte)
	output := make(chan []byte)
	go MemChan(input, output)

	group := new(sync.WaitGroup)
	group.Add(1)

	go func() {
		for i := 0; i < 10; i++ {
			input <- []byte("test")
		}

		go func() {
			defer group.Done()
			for out := range output {
				fmt.Println("Output:", string(out))
			}
		}()

		close(input)

	}()

	group.Wait()
}

package buffer

import (
	"fmt"
	"sync"
	"testing"
)

func TestPipe(t *testing.T) {
	r, w := Pipe()

	group := new(sync.WaitGroup)
	group.Add(1)

	go func(){
		for i := 0; i < 10; i++ {
			w.Write([]byte(fmt.Sprintf("%d", i)))
		}
		w.Close()
	}()
	

	go func(){
		defer group.Done()
		data := make([]byte, 32*1024)
		for {
			if n, err := r.Read(data); err == nil {
				fmt.Print(string(data[:n]))
			} else {
				fmt.Println("")
				return
			}
		}
	}()

	group.Wait()
}

func TestMemChan(t *testing.T) {
	input := make(chan []byte)
	output := make(chan []byte)
	go MemChan(input, output)

	group := new(sync.WaitGroup)
	group.Add(1)

	go func() {
		for i := 0; i < 10; i++ {
			input <- []byte(fmt.Sprintf("%d", i))
		}

		go func() {
			defer group.Done()
			for out := range output {
				fmt.Print(string(out))
			}
			fmt.Println("")
		}()

		close(input)

	}()

	group.Wait()
}

func TestChan(t *testing.T) {
	input := make(chan []byte)
	output := make(chan []byte)
	go Chan(input, output)

	group := new(sync.WaitGroup)
	group.Add(1)

	go func() {
		for i := 0; i < 10; i++ {
			input <- []byte(fmt.Sprintf("%d", i))
		}

		go func() {
			defer group.Done()
			for out := range output {
				fmt.Print(string(out))
			}
			fmt.Println("")
		}()

		close(input)

	}()

	group.Wait()
}

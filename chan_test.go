package buffer

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"sync"
	"testing"
)

func TestPipe(t *testing.T) {
	r, w := Pipe()

	go func() {
		for i := 0; i < 10; i++ {
			w.Write([]byte(fmt.Sprintf("%d", i)))
		}
		w.Close()
	}()

	data, _ := ioutil.ReadAll(r)
	if !bytes.Equal(data, []byte("0123456789")) {
		t.Error("Not equal!")
	}
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
			s := ""
			for out := range output {
				s += string(out)
			}
			if s != "0123456789" {
				t.Error("Not equal!")
			}
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
			s := ""
			for out := range output {
				s += string(out)
			}
			if s != "0123456789" {
				t.Error("Not equal!")
			}
		}()

		close(input)

	}()

	group.Wait()
}

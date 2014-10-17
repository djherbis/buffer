package buffer

import (
	"io"
)

func Pipe() (io.ReadCloser, io.WriteCloser) {
	inReader, inWriter := io.Pipe()
	outReader, outWriter := io.Pipe()

	input := make(chan []byte)
	output := make(chan []byte)

	go Chan(input, output)
	go inFeed(inReader, input)
	go outFeed(outWriter, output)

	return outReader, inWriter
}

func inFeed(r io.ReadCloser, in chan<- []byte) {
	for {
		data := make([]byte, 32*1024)
		if n, err := r.Read(data); err == nil {
			in <- data[:n]
		} else {
			close(in)
			return
		}
	}
}

func outFeed(w io.WriteCloser, out <-chan []byte) {
	for output := range out {
		for len(output) > 0 {
			n, _ := w.Write(output)
			output = output[n:]
		}

	}
	w.Close()
}

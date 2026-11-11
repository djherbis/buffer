package buffer_test

import (
	"bytes"
	"io"
	"io/ioutil"
	"testing"

	"github.com/djherbis/buffer"
)

func TestMultiReadToEnd(t *testing.T) {
	payload := bytes.Repeat([]byte("buffered data\n"), 400)
	factories := []struct {
		name string
		new  func() buffer.Buffer
	}{
		{"zero memory", func() buffer.Buffer {
			return buffer.NewMulti(buffer.New(0), buffer.New(int64(len(payload))))
		}},
		{"zero memory at", func() buffer.Buffer {
			return buffer.NewMultiAt(buffer.New(0), buffer.New(int64(len(payload))))
		}},
		{"nested zero memory", func() buffer.Buffer {
			return buffer.NewMulti(buffer.New(0), buffer.New(0), buffer.New(int64(len(payload))))
		}},
		{"file only", func() buffer.Buffer {
			return buffer.NewUnboundedBuffer(0, 1024)
		}},
		{"memory and memory", func() buffer.Buffer {
			return buffer.NewMulti(buffer.New(64), buffer.New(int64(len(payload))))
		}},
		{"memory and files", func() buffer.Buffer {
			return buffer.NewUnboundedBuffer(64, 1024)
		}},
	}
	readers := []struct {
		name string
		read func(io.Reader) ([]byte, error)
	}{
		{"ReadAll", ioutil.ReadAll},
		{"Copy", func(r io.Reader) ([]byte, error) {
			var dst bytes.Buffer
			_, err := io.Copy(&dst, r)
			return dst.Bytes(), err
		}},
	}
	for _, factory := range factories {
		for _, reader := range readers {
			t.Run(factory.name+"/"+reader.name, func(t *testing.T) {
				buf := factory.new()
				defer buf.Reset()
				for cycle := 0; cycle < 2; cycle++ {
					if n, err := buf.Write(payload); n != len(payload) || err != nil {
						t.Fatalf("Write = (%d, %v), want (%d, nil)", n, err, len(payload))
					}
					got, err := reader.read(buf)
					if err != nil {
						t.Fatal(err)
					}
					if !bytes.Equal(got, payload) {
						t.Fatalf("cycle %d: read %d bytes, want %d; %d bytes remain buffered", cycle, len(got), len(payload), buf.Len())
					}
					if buf.Len() != 0 {
						t.Fatalf("Len = %d after draining", buf.Len())
					}
					if n, err := buf.Read(make([]byte, 1)); n != 0 || err != io.EOF {
						t.Fatalf("empty Read = (%d, %v), want (0, EOF)", n, err)
					}
				}
			})
		}
	}
}

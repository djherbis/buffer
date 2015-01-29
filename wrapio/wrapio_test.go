package wrapio

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"
)

func TestWrap(t *testing.T) {
	if file, err := ioutil.TempFile("", "wrap.test"); err == nil {
		w := NewWrapWriter(file, 0, 3)
		w.Write([]byte("abcdef"))

		r := NewWrapReader(file, 0, 2)
		data := make([]byte, 6)
		r.Read(data)
		if !bytes.Equal(data, []byte("dedede")) {
			t.Error("Wrapper error!")
		}
		file.Close()
		os.Remove(file.Name())
	}
}

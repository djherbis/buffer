package buffer

import (
	"bytes"
	"crypto/rand"
	"io"
	"io/ioutil"
	"os"
	"testing"
)

func TestSpill(t *testing.T) {
	buf := NewMulti(New(5), NewDiscard())
	buf.Write([]byte("Hello World"))
	data := make([]byte, 12)
	n, _ := buf.Read(data)
	if !bytes.Equal(data[:n], []byte("Hello")) {
		t.Error("ReadAt Failed. " + string(data[:n]))
	}
}

func TestReadAt(t *testing.T) {
	buf := NewMulti(NewFile(2), NewFile(2), NewFile(2), NewFile(2), NewFile(2), NewFile(2))
	buf.Write([]byte("Hello World"))
	data := make([]byte, 10)
	n, _ := buf.ReadAt(data, 6)
	if !bytes.Equal(data[:n], []byte("World")) {
		t.Error("ReadAt Failed. " + string(data[:n]))
	}
	buf.Reset()
}

func TestFF(t *testing.T) {
	buf := NewMulti(New(2), New(2), New(2), New(2), New(2), New(2), New(2))
	buf.Write([]byte("Hello"))
	buf.FastForward(3)
	data := make([]byte, 4)
	n, _ := buf.Read(data)
	if !bytes.Equal(data[:n], []byte("lo")) {
		t.Error("FastForward error! " + string(data))
	}
}

func TestWrap(t *testing.T) {
	if file, err := ioutil.TempFile("D:\\Downloads\\temp", "wrap.test"); err == nil {
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

func TestFile(t *testing.T) {
	buf := NewFile(1024)
	checkCap(t, buf, 1024)
	runPerfectSeries(t, buf)
	buf.Reset()

	buf = NewFile(3)
	buf.Write([]byte("abc"))
	buf.Read(make([]byte, 1))
	buf.Write([]byte("a"))
	d, _ := ioutil.ReadAll(buf)
	if !bytes.Equal(d, []byte("bca")) {
		t.Error("back and forth error!")
	}
}

func TestMem(t *testing.T) {
	buf := New(1024)
	checkCap(t, buf, 1024)
	runPerfectSeries(t, buf)
	buf.Reset()
}

func TestFilePartition(t *testing.T) {
	buf := NewPartition(func() Buffer { return NewFile(1024) })
	checkCap(t, buf, MaxCap())
	runPerfectSeries(t, buf)
	buf.Reset()
}

func TestMulti(t *testing.T) {
	buf := NewMulti(New(5), New(5), NewFile(500), NewPartition(func() Buffer { return New(1024) }))
	checkCap(t, buf, MaxCap())
	runPerfectSeries(t, buf)
	isPerfectMatch(t, buf, 1024*1024)
	buf.Reset()
}

func runPerfectSeries(t *testing.T, buf Buffer) {
	checkEmpty(t, buf)
	simple(t, buf)

	max := LimitAlloc(buf.Cap())
	isPerfectMatch(t, buf, 0)
	for i := int64(1); i < max; i *= 2 {
		isPerfectMatch(t, buf, i)
	}
	isPerfectMatch(t, buf, max)
}

func simple(t *testing.T, buf Buffer) {
	buf.Write([]byte("hello world"))
	data, err := ioutil.ReadAll(buf)
	if err != nil {
		t.Error(err.Error())
	}
	if !bytes.Equal([]byte("hello world"), data) {
		t.Error("Hello world failed.")
	}

	buf.Write([]byte("hello world"))
	data = make([]byte, 3)
	buf.Read(data)
	buf.Write([]byte(" yolo"))
	data, err = ioutil.ReadAll(buf)
	if err != nil {
		t.Error(err.Error())
	}
	if !bytes.Equal([]byte("lo world yolo"), data) {
		t.Error("Buffer crossing error :(", string(data))
	}
}

func buildOutputs(t *testing.T, buf Buffer, size int64) (wrote []byte, read []byte) {
	r := io.LimitReader(rand.Reader, size)
	tee := io.TeeReader(r, buf)

	wrote, _ = ioutil.ReadAll(tee)
	read, _ = ioutil.ReadAll(buf)

	return wrote, read
}

func isPerfectMatch(t *testing.T, buf Buffer, size int64) {
	wrote, read := buildOutputs(t, buf, size)
	if !bytes.Equal(wrote, read) {
		t.Error("Buffer should have matched")
	}
}

func checkEmpty(t *testing.T, buf Buffer) {
	if !Empty(buf) {
		t.Error("Buffer should start empty!")
	}
}

func checkCap(t *testing.T, buf Buffer, correctCap int64) {
	if buf.Cap() != correctCap {
		t.Error("Buffer cap is incorrect", buf.Cap(), correctCap)
	}
}

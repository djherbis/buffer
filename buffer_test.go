package buffer

import (
	"bytes"
	"crypto/rand"
	"encoding/gob"
	"io"
	"io/ioutil"
	"os"
	"testing"

	"github.com/djherbis/buffer/wrapio"
)

func BenchmarkMemory(b *testing.B) {
	buf := New(32 * 1024)
	for i := 0; i < b.N; i++ {
		io.Copy(buf, io.LimitReader(rand.Reader, 32*1024))
		io.Copy(ioutil.Discard, buf)
	}
}

func TestWriteAt(t *testing.T) {
	var b BufferAt

	b = New(5)
	BufferAtTester(t, b)

	file, err := ioutil.TempFile("", "buffer")
	if err != nil {
		t.Error(err.Error())
	}
	defer os.Remove(file.Name())
	defer file.Close()

	b = NewFile(5, file)
	BufferAtTester(t, b)

}

func BufferAtTester(t *testing.T, b BufferAt) {
	b.WriteAt([]byte("abc"), 0)
	Compare(t, b, "abc")

	b.WriteAt([]byte("abc"), 1)
	Compare(t, b, "aabc")

	b.WriteAt([]byte("abc"), 2)
	Compare(t, b, "aaabc")

	b.WriteAt([]byte("abc"), 3)
	Compare(t, b, "aaaab")

	b.Read(make([]byte, 2))

	Compare(t, b, "aab")
	b.Reset()
}

func Compare(t *testing.T, b BufferAt, s string) {
	data := make([]byte, b.Len())
	n, _ := b.ReadAt(data, 0)
	if string(data[:n]) != s {
		t.Error("Mismatch:", string(data[:n]), s)
	}
}

// TODO Rebuild Ring Tests

// TODO Rebuild Spill Tests

/**/
func TestGob(t *testing.T) {
	str := "HelloWorld"

	file, err := ioutil.TempFile(os.TempDir(), "buffer")
	if err != nil {
		t.Error(err.Error())
	}
	defer os.Remove(file.Name())
	defer file.Close()

	buf := NewFile(20, file)
	buf.Write([]byte(str))
	b := bytes.NewBuffer(nil)
	var test Buffer = buf
	if err := gob.NewEncoder(b).Encode(&test); err != nil {
		t.Error(err.Error())
		return
	}
	var buffer Buffer
	if err := gob.NewDecoder(b).Decode(&buffer); err != nil {
		t.Error(err.Error())
		return
	}
	data := make([]byte, 10)
	if _, err := buffer.Read(data); err != nil && err != io.EOF {
		t.Error(err.Error())
	}
	if !bytes.Equal(data, []byte(str)) {
		t.Error("Gob Recover Failed... " + string(data))
	}
	buf.Reset()
}

/**/

func TestSpill(t *testing.T) {
	buf := NewMulti(New(5), NewDiscard())
	buf.Write([]byte("Hello World"))
	data := make([]byte, 12)
	n, _ := buf.Read(data)
	if !bytes.Equal(data[:n], []byte("Hello")) {
		t.Error("ReadAt Failed. " + string(data[:n]))
	}
}

func TestWrap(t *testing.T) {
	if file, err := ioutil.TempFile("", "wrap.test"); err == nil {
		w := wrapio.NewWrapWriter(file, 0, 3)
		w.Write([]byte("abcdef"))

		r := wrapio.NewWrapReader(file, 0, 2)
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
	file, err := ioutil.TempFile("", "buffer")
	if err != nil {
		t.Error(err.Error())
	}
	defer os.Remove(file.Name())
	defer file.Close()

	buf := NewFile(1024, file)
	checkCap(t, buf, 1024)
	runPerfectSeries(t, buf)
	buf.Reset()

	buf = NewFile(3, file)
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
	buf := NewPartition(func() Buffer { return New(1024) })
	checkCap(t, buf, MAXINT64)
	runPerfectSeries(t, buf)
	buf.Reset()
}

func TestMulti(t *testing.T) {
	file, err := ioutil.TempFile("", "buffer")
	if err != nil {
		t.Error(err.Error())
	}
	defer os.Remove(file.Name())
	defer file.Close()

	buf := NewMulti(New(5), New(5), NewFile(500, file), NewPartition(func() Buffer { return New(1024) }))
	checkCap(t, buf, MAXINT64)
	runPerfectSeries(t, buf)
	isPerfectMatch(t, buf, 1024*1024)
	buf.Reset()
}

func runPerfectSeries(t *testing.T, buf Buffer) {
	checkEmpty(t, buf)
	simple(t, buf)

	max := int64(1024)
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

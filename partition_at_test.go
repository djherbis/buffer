package buffer

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestPartitionAtWriteAtCount(t *testing.T) {
	for _, tt := range []struct {
		name    string
		initial string
		data    string
		off     int64
		want    string
	}{
		{"single buffer", "", "ab", 0, "ab"},
		{"two buffers", "", "abcd", 0, "abcd"},
		{"three buffers", "", "abcdefg", 0, "abcdefg"},
		{"overwrite across buffers", "abcdefghi", "WXYZ", 2, "abWXYZghi"},
		{"extend across buffers", "abc", "WXYZ", 2, "abWXYZ"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			buf := NewPartitionAt(NewMemPoolAt(3))
			defer buf.Reset()
			if _, err := buf.Write([]byte(tt.initial)); err != nil {
				t.Fatal(err)
			}
			n, err := buf.WriteAt([]byte(tt.data), tt.off)
			if n != len(tt.data) || err != nil {
				t.Errorf("WriteAt = (%d, %v), want (%d, nil)", n, err, len(tt.data))
			}
			got, err := ioutil.ReadAll(buf)
			if err != nil || string(got) != tt.want {
				t.Errorf("ReadAll = (%q, %v), want (%q, nil)", got, err, tt.want)
			}
		})
	}
}

func TestPartitionAtWriteAtCountOnGrowError(t *testing.T) {
	dir, err := ioutil.TempDir("", "buffer-pool")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	buf := NewPartitionAt(NewFilePoolAt(3, filepath.Join(dir, "missing")), New(3))
	defer buf.Reset()
	if _, err := buf.Write([]byte("x")); err != nil {
		t.Fatal(err)
	}
	n, err := buf.WriteAt([]byte("abcd"), 0)
	if n != 3 || !os.IsNotExist(err) {
		t.Errorf("WriteAt = (%d, %v), want (3, file-not-found error)", n, err)
	}
	got, err := ioutil.ReadAll(buf)
	if err != nil || string(got) != "abc" {
		t.Errorf("ReadAll = (%q, %v), want (%q, nil)", got, err, "abc")
	}
}

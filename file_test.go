package buffer

import (
	"bytes"
	"encoding/gob"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestGobFileWhitespacePath(t *testing.T) {
	dir, err := ioutil.TempDir("", "buffer-gob")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	tests := []struct {
		name string
		path string
	}{
		{"plain", "file"},
		{"space", "file with spaces"},
		{"directory", filepath.Join("directory with spaces", "file")},
		{"tab", "file\twith\ttabs"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "tab" && runtime.GOOS == "windows" {
				t.Skip("Windows does not allow tabs in file names")
			}
			filename := filepath.Join(dir, tt.path)
			if err := os.MkdirAll(filepath.Dir(filename), 0700); err != nil {
				t.Fatal(err)
			}
			file, err := os.Create(filename)
			if err != nil {
				t.Fatal(err)
			}
			defer file.Close()
			buf := NewFile(10, file)
			if n, err := buf.Write([]byte("abcdef")); n != 6 || err != nil {
				t.Fatalf("initial write = %d, %v; want 6, nil", n, err)
			}
			if n, err := buf.Read(make([]byte, 2)); n != 2 || err != nil {
				t.Fatalf("initial read = %d, %v; want 2, nil", n, err)
			}
			var encoded bytes.Buffer
			if err := gob.NewEncoder(&encoded).Encode(&buf); err != nil {
				t.Fatal(err)
			}
			var restored BufferAt
			if err := gob.NewDecoder(&encoded).Decode(&restored); err != nil {
				t.Fatal(err)
			}
			defer restored.(*fileBuffer).file.Close()
			if restored.Cap() != 10 || restored.Len() != 4 {
				t.Fatalf("restored cap, len = %d, %d; want 10, 4", restored.Cap(), restored.Len())
			}
			if n, err := restored.Write([]byte("gh")); n != 2 || err != nil {
				t.Fatalf("restored write = %d, %v; want 2, nil", n, err)
			}
			data, err := ioutil.ReadAll(restored)
			if err != nil {
				t.Fatal(err)
			}
			if string(data) != "cdefgh" {
				t.Fatalf("restored contents = %q; want %q", data, "cdefgh")
			}
		})
	}
}

func TestFileUnmarshalMissingFilename(t *testing.T) {
	for _, data := range []string{"", "filename", "\n10 0 0\n"} {
		var buf fileBuffer
		if err := buf.UnmarshalBinary([]byte(data)); err == nil {
			t.Errorf("UnmarshalBinary(%q) succeeded without a complete filename", data)
		}
	}
}

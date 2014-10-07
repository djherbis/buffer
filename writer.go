package buffer

import(
  "io"
)

type Writer struct {
  buf Buffer
  io.Writer
  buffering bool
}

func NewWriter(w io.Writer, b Buffer) *Writer {
  return &Writer{
    Writer: w,
    buf: b,
  }
}

func (b *Writer) Available() int64 {
  return Gap(b.buf)
} 

func (b *Writer) Buffered() int64 {
  return b.buf.Len()
}

func (b *Writer) Write(p []byte) (n int, err error) {
  for len(p) > 0 {
    if b.buffering {
      m, err := b.buf.Write(p)
      n += m
      p = p[m:]
      if err != nil {
        return n, err
      }
      b.Flush()
    } else {
      m, er := b.Writer.Write(p)
      n += m
      p = p[m:]
      if er != nil {
        b.buffering = true
      }
    }
  }
  return n, nil
}

func (b *Writer) Flush() error {
  for !Empty(b.buf) {
    if _, err := io.Copy(b.Writer, b.buf); err != nil {
      return err
    }
  }
  b.buffering = false
  return nil
}

func (b *Writer) Close() (err error) {
  
  if err = b.Flush(); err != nil {
    return err
  }

  if closer, ok := b.Writer.(io.Closer); ok {
    if err = closer.Close(); err != nil {
      return err
    }
  }

  return nil
}
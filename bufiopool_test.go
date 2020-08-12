package bp

import(
  "testing"
  "strings"
  "bytes"
  "bufio"
)

func TestBufioPoolBufSize(t *testing.T) {
  t.Run("reader/default/get", func(tt *testing.T) {
    p := NewBufioReaderPool(10)

    d1 := p.Get(strings.NewReader("hello"))
    d2 := p.Get(strings.NewReader("helloworld"))
    if d1.Size() != defaultBufioSize {
      tt.Errorf("expect size=%d actual:%d", defaultBufioSize, d1.Size())
    }
    if d2.Size() != defaultBufioSize {
      tt.Errorf("expect size=%d actual:%d", defaultBufioSize, d2.Size())
    }
  })
  t.Run("writer/default/get", func(tt *testing.T) {
    p := NewBufioWriterPool(10)

    d1 := p.Get(bytes.NewBuffer(nil))
    d2 := p.Get(bytes.NewBuffer(make([]byte, 0, 1024)))
    if d1.Size() != defaultBufioSize {
      tt.Errorf("expect size=%d actual:%d", defaultBufioSize, d1.Size())
    }
    if d2.Size() != defaultBufioSize {
      tt.Errorf("expect size=%d actual:%d", defaultBufioSize, d2.Size())
    }
  })
  t.Run("reader/default/putsmall", func(tt *testing.T) {
    p := NewBufioReaderPool(10)

    p.Put(bufio.NewReaderSize(nil, 100))
    d1 := p.Get(strings.NewReader("hello"))
    if d1.Size() != defaultBufioSize {
      tt.Errorf("small buffer rearrenged default buffer: %d", d1.Size())
    }
  })
  t.Run("reader/default/putlarge", func(tt *testing.T) {
    p := NewBufioReaderPool(10)

    size := 8 * 1024
    p.Put(bufio.NewReaderSize(nil, size))
    d1 := p.Get(strings.NewReader("hello"))
    if d1.Size() != size {
      tt.Errorf("no strict check large buffer by default: %d", d1.Size())
    }
  })
  t.Run("writer/default/putsmall", func(tt *testing.T) {
    p := NewBufioWriterPool(10)

    p.Put(bufio.NewWriterSize(nil, 100))
    d1 := p.Get(bytes.NewBuffer(nil))
    if d1.Size() != defaultBufioSize {
      tt.Errorf("small buffer rearrenged default buffer: %d", d1.Size())
    }
  })
  t.Run("writer/default/putlarge", func(tt *testing.T) {
    p := NewBufioWriterPool(10)

    size := 8 * 1024
    p.Put(bufio.NewWriterSize(nil, size))
    d1 := p.Get(bytes.NewBuffer(nil))
    if d1.Size() != size {
      tt.Errorf("no strict check large buffer by default: %d", d1.Size())
    }
  })

  t.Run("reader/8kb/get", func(tt *testing.T) {
    bufSize := 8 * 1024
    p := NewBufioReaderSizePool(10, bufSize)

    d1 := p.Get(strings.NewReader("hello"))
    d2 := p.Get(strings.NewReader("helloworld"))
    if d1.Size() != bufSize {
      tt.Errorf("expect size=%d actual:%d", defaultBufioSize, d1.Size())
    }
    if d2.Size() != bufSize {
      tt.Errorf("expect size=%d actual:%d", defaultBufioSize, d2.Size())
    }
  })
  t.Run("writer/8kb/get", func(tt *testing.T) {
    bufSize := 8 * 1024
    p := NewBufioWriterSizePool(10, bufSize)

    d1 := p.Get(bytes.NewBuffer(nil))
    d2 := p.Get(bytes.NewBuffer(make([]byte, 0, 1024)))
    if d1.Size() != bufSize {
      tt.Errorf("expect size=%d actual:%d", defaultBufioSize, d1.Size())
    }
    if d2.Size() != bufSize {
      tt.Errorf("expect size=%d actual:%d", defaultBufioSize, d2.Size())
    }
  })
  t.Run("reader/8kb/putsmall", func(tt *testing.T) {
    bufSize := 8 * 1024
    p := NewBufioReaderSizePool(10, bufSize)

    p.Put(bufio.NewReaderSize(nil, 4 * 1024))
    d1 := p.Get(strings.NewReader(""))
    if d1.Size() != bufSize {
      tt.Errorf("small buffer rearrenged sized buffer: %d", d1.Size())
    }
  })
  t.Run("writer/8kb/putsmall", func(tt *testing.T) {
    bufSize := 8 * 1024
    p := NewBufioWriterSizePool(10, bufSize)

    p.Put(bufio.NewWriterSize(nil, 4 * 1024))
    d1 := p.Get(bytes.NewBuffer(nil))
    if d1.Size() != bufSize {
      tt.Errorf("small buffer rearrenged sized buffer: %d", d1.Size())
    }
  })
  t.Run("reader/8kb/putlarge", func(tt *testing.T) {
    bufSize := 8 * 1024
    p := NewBufioReaderSizePool(10, bufSize)

    p.Put(bufio.NewReaderSize(nil, 16 * 1024))
    d1 := p.Get(strings.NewReader(""))
    if d1.Size() != bufSize {
      tt.Errorf("strict same size buffer by sized pool: %d", d1.Size())
    }
  })
  t.Run("writer/8kb/putlarge", func(tt *testing.T) {
    bufSize := 8 * 1024
    p := NewBufioWriterSizePool(10, bufSize)

    p.Put(bufio.NewWriterSize(nil, 16 * 1024))
    d1 := p.Get(bytes.NewBuffer(nil))
    if d1.Size() != bufSize {
      tt.Errorf("strict same size buffer by sized pool: %d", d1.Size())
    }
  })
  t.Run("reader/8kb/putsamesize", func(tt *testing.T) {
    bufSize := 8 * 1024
    p := NewBufioReaderSizePool(10, bufSize)

    p.Put(bufio.NewReaderSize(nil, bufSize))
    d1 := p.Get(strings.NewReader(""))
    if d1.Size() != bufSize {
      tt.Errorf("same size ok: %d", d1.Size())
    }
  })
  t.Run("writer/8kb/putsamesize", func(tt *testing.T) {
    bufSize := 8 * 1024
    p := NewBufioWriterSizePool(10, bufSize)

    p.Put(bufio.NewWriterSize(nil, bufSize))
    d1 := p.Get(bytes.NewBuffer(nil))
    if d1.Size() != bufSize {
      tt.Errorf("same size ok: %d", d1.Size())
    }
  })
}

func TestBufioPoolLenCap(t *testing.T) {
  t.Run("reader/default/getput", func(tt *testing.T) {
    p := NewBufioReaderPool(10)
    if 0 != p.Len() {
      tt.Errorf("initial len 0")
    }
    if 10 != p.Cap() {
      tt.Errorf("initial cap 10")
    }

    data := p.Get(strings.NewReader(""))
    if 0 != p.Len() {
      tt.Errorf("initial len 0")
    }
    if 10 != p.Cap() {
      tt.Errorf("initial cap 10")
    }
    p.Put(data)

    if 1 != p.Len() {
      tt.Errorf("put one")
    }
    if 10 != p.Cap() {
      tt.Errorf("initial cap 10")
    }

    d1 := p.Get(strings.NewReader(""))
    if 0 != p.Len() {
      tt.Errorf("aquire pool")
    }
    p.Put(d1)
    if 1 != p.Len() {
      tt.Errorf("release pool")
    }
  })
  t.Run("reader/2kb/getput", func(tt *testing.T) {
    bufSize := 2 * 1024
    p := NewBufioReaderSizePool(10, bufSize)
    if 0 != p.Len() {
      tt.Errorf("initial len 0")
    }
    if 10 != p.Cap() {
      tt.Errorf("initial cap 10")
    }

    data := p.Get(strings.NewReader(""))
    if 0 != p.Len() {
      tt.Errorf("initial len 0")
    }
    if 10 != p.Cap() {
      tt.Errorf("initial cap 10")
    }
    p.Put(data)

    if 1 != p.Len() {
      tt.Errorf("put one")
    }
    if 10 != p.Cap() {
      tt.Errorf("initial cap 10")
    }

    d1 := p.Get(strings.NewReader(""))
    if 0 != p.Len() {
      tt.Errorf("aquire pool")
    }
    p.Put(d1)
    if 1 != p.Len() {
      tt.Errorf("release pool")
    }
  })
  t.Run("reader/default/maxcap", func(tt *testing.T) {
    p := NewBufioReaderPool(10)
    s := make([]*bufio.Reader, 0)
    for i := 0; i < 10; i += 1 {
      d := p.Get(strings.NewReader(""))
      s = append(s, d)
    }
    for _, d := range s {
      p.Put(d)
    }
    if 10 != p.Len() {
      tt.Errorf("fill-ed pool: %d", p.Len())
    }
    if 10 != p.Cap() {
      tt.Errorf("max capacity = 10")
    }
    p.Put(bufio.NewReaderSize(nil, defaultBufioSize))
    p.Put(bufio.NewReaderSize(nil, defaultBufioSize))

    if 10 != p.Len() {
      tt.Errorf("fixed size pool")
    }
    if 10 != p.Cap() {
      tt.Errorf("max capacity = 10")
    }
  })
  t.Run("reader/2kb/maxcap", func(tt *testing.T) {
    bufSize := 2 * 1024
    p := NewBufioReaderSizePool(10, bufSize)
    s := make([]*bufio.Reader, 0)
    for i := 0; i < 10; i += 1 {
      d := p.Get(strings.NewReader(""))
      s = append(s, d)
    }
    for _, d := range s {
      p.Put(d)
    }
    if 10 != p.Len() {
      tt.Errorf("fill-ed pool: %d", p.Len())
    }
    if 10 != p.Cap() {
      tt.Errorf("max capacity = 10")
    }
    p.Put(bufio.NewReaderSize(nil, bufSize))
    p.Put(bufio.NewReaderSize(nil, bufSize))

    if 10 != p.Len() {
      tt.Errorf("fixed size pool")
    }
    if 10 != p.Cap() {
      tt.Errorf("max capacity = 10")
    }
  })

  t.Run("writer/default/getput", func(tt *testing.T) {
    p := NewBufioWriterPool(10)
    if 0 != p.Len() {
      tt.Errorf("initial len 0")
    }
    if 10 != p.Cap() {
      tt.Errorf("initial cap 10")
    }

    data := p.Get(bytes.NewBuffer(nil))
    if 0 != p.Len() {
      tt.Errorf("initial len 0")
    }
    if 10 != p.Cap() {
      tt.Errorf("initial cap 10")
    }
    p.Put(data)

    if 1 != p.Len() {
      tt.Errorf("put one")
    }
    if 10 != p.Cap() {
      tt.Errorf("initial cap 10")
    }

    d1 := p.Get(bytes.NewBuffer(nil))
    if 0 != p.Len() {
      tt.Errorf("aquire pool")
    }
    p.Put(d1)
    if 1 != p.Len() {
      tt.Errorf("release pool")
    }
  })
  t.Run("writer/2kb/getput", func(tt *testing.T) {
    bufSize := 2 * 1024
    p := NewBufioWriterSizePool(10, bufSize)
    if 0 != p.Len() {
      tt.Errorf("initial len 0")
    }
    if 10 != p.Cap() {
      tt.Errorf("initial cap 10")
    }

    data := p.Get(bytes.NewBuffer(nil))
    if 0 != p.Len() {
      tt.Errorf("initial len 0")
    }
    if 10 != p.Cap() {
      tt.Errorf("initial cap 10")
    }
    p.Put(data)

    if 1 != p.Len() {
      tt.Errorf("put one")
    }
    if 10 != p.Cap() {
      tt.Errorf("initial cap 10")
    }

    d1 := p.Get(bytes.NewBuffer(nil))
    if 0 != p.Len() {
      tt.Errorf("aquire pool")
    }
    p.Put(d1)
    if 1 != p.Len() {
      tt.Errorf("release pool")
    }
  })
  t.Run("writer/default/maxcap", func(tt *testing.T) {
    p := NewBufioWriterPool(10)
    s := make([]*bufio.Writer, 0)
    for i := 0; i < 10; i += 1 {
      d := p.Get(bytes.NewBuffer(nil))
      s = append(s, d)
    }
    for _, d := range s {
      p.Put(d)
    }
    if 10 != p.Len() {
      tt.Errorf("fill-ed pool: %d", p.Len())
    }
    if 10 != p.Cap() {
      tt.Errorf("max capacity = 10")
    }
    p.Put(bufio.NewWriterSize(nil, defaultBufioSize))
    p.Put(bufio.NewWriterSize(nil, defaultBufioSize))

    if 10 != p.Len() {
      tt.Errorf("fixed size pool")
    }
    if 10 != p.Cap() {
      tt.Errorf("max capacity = 10")
    }
  })
  t.Run("writer/2kb/maxcap", func(tt *testing.T) {
    bufSize := 2 * 1024
    p := NewBufioWriterSizePool(10, bufSize)
    s := make([]*bufio.Writer, 0)
    for i := 0; i < 10; i += 1 {
      d := p.Get(bytes.NewBuffer(nil))
      s = append(s, d)
    }
    for _, d := range s {
      p.Put(d)
    }
    if 10 != p.Len() {
      tt.Errorf("fill-ed pool: %d", p.Len())
    }
    if 10 != p.Cap() {
      tt.Errorf("max capacity = 10")
    }
    p.Put(bufio.NewWriterSize(nil, bufSize))
    p.Put(bufio.NewWriterSize(nil, bufSize))

    if 10 != p.Len() {
      tt.Errorf("fixed size pool")
    }
    if 10 != p.Cap() {
      tt.Errorf("max capacity = 10")
    }
  })
}

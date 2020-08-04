package bp

import(
  "testing"
  "bytes"
)

func TestBufferPoolBufSize(t *testing.T) {
  bufSize := 8
  t.Run("getsamecap", func(tt *testing.T) {
    p := NewBufferPool(10, bufSize)

    d1 := p.Get()
    d2 := p.Get()
    if d1.Cap() != bufSize {
      tt.Errorf("buf size = %d", bufSize)
    }
    if d2.Cap() != bufSize {
      tt.Errorf("buf size = %d", bufSize)
    }
  })
  t.Run("getput/samecap", func(tt *testing.T) {
    p := NewBufferPool(10, bufSize)
    d1 := p.Get()
    p.Put(d1)
    d2 := p.Get()
    p.Put(d2)

    d3 := p.Get()
    d4 := p.Get()
    if d3.Cap() != bufSize {
      tt.Errorf("buf size = %d", bufSize)
    }
    if d4.Cap() != bufSize {
      tt.Errorf("buf size = %d", bufSize)
    }
  })
  t.Run("getput/largecap", func(tt *testing.T) {
    p := NewBufferPool(10, bufSize)
    p.Put(bytes.NewBuffer(make([]byte, 0, 123)))

    d1 := p.Get()
    if d1.Cap() == bufSize {
      tt.Errorf("manually set buffer capacity be different")
    }
    if d1.Len() != 0 {
      tt.Errorf("manually set buffer be Reset")
    }
  })
}

func TestBufferPoolDiscard(t *testing.T) {
  bufSize := 8
  t.Run("freecap/samesize", func(tt *testing.T) {
    p := NewBufferPool(10, bufSize)

    d1 := bytes.NewBuffer(make([]byte, 0, bufSize))
    if p.Put(d1) != true {
      tt.Errorf("freecap %d", p.Cap())
    }
  })
  t.Run("fullcap/samesize", func(tt *testing.T) {
    p := NewBufferPool(2, bufSize)
    p.Put(bytes.NewBuffer(make([]byte, 0, bufSize)))
    p.Put(bytes.NewBuffer(make([]byte, 0, bufSize)))

    d1 := bytes.NewBuffer(make([]byte, 0, bufSize))
    if p.Put(d1) {
      tt.Errorf("fulled capacity %d", p.Cap())
    }
  })
  t.Run("freecap/largesize", func(tt *testing.T) {
    p := NewBufferPool(10, bufSize)
    p.Put(bytes.NewBuffer(make([]byte, 0, bufSize)))

    if p.Put(bytes.NewBuffer(make([]byte, 0, 100))) != true {
      tt.Errorf("put ok")
    }
    if p.Put(bytes.NewBuffer(make([]byte, 100))) != true {
      tt.Errorf("put ok")
    }
  })
  t.Run("freecap/smallsize", func(tt *testing.T) {
    p := NewBufferPool(10, bufSize)
    p.Put(bytes.NewBuffer(make([]byte, 0, bufSize)))

    if p.Put(bytes.NewBuffer(nil)) != true {
      tt.Errorf("put ok // call bytes.Grow")
    }
    if p.Put(bytes.NewBuffer(make([]byte, 0, 1))) != true {
      tt.Errorf("put ok // call bytes.Grow")
    }
    if p.Put(bytes.NewBuffer(make([]byte, 1))) != true {
      tt.Errorf("put ok // call bytes.Grow")
    }
  })
}

func TestBufferPoolLenCap(t *testing.T) {
  t.Run("getput", func(tt *testing.T) {
    bufSize := 1
    p := NewBufferPool(10, bufSize)
    if 0 != p.Len() {
      tt.Errorf("initial len 0")
    }
    if 10 != p.Cap() {
      tt.Errorf("initial cap 10")
    }

    data := p.Get()
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

    d1 := p.Get()
    if 0 != p.Len() {
      tt.Errorf("aquire pool")
    }
    p.Put(d1)
    if 1 != p.Len() {
      tt.Errorf("release pool")
    }
  })
  t.Run("maxcap", func(tt *testing.T) {
    bufSize := 1
    p := NewBufferPool(10, bufSize)
    s := make([]*bytes.Buffer, 0)

    for i := 0; i < 10; i += 1 {
      d := p.Get()
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

    d1 := bytes.NewBuffer(make([]byte, 0, bufSize))
    d2 := bytes.NewBuffer(make([]byte, 0, bufSize))
    p.Put(d1)
    p.Put(d2)

    if 10 != p.Len() {
      tt.Errorf("fixed size pool")
    }
    if 10 != p.Cap() {
      tt.Errorf("max capacity = 10")
    }
  })
}

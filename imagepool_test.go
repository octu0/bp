package bp

import(
  "testing"
  "image"
)

func TestImageRGBAPoolBufSize(t *testing.T) {
  t.Run("getsamecap", func(tt *testing.T) {
    rect := image.Rect(0, 0, 100, 100)
    pool := NewImageRGBAPool(10, rect)
    img  := image.NewRGBA(rect)

    d1 := pool.GetRef()
    d2 := pool.GetRef()
    if cap(d1.pix) != cap(img.Pix) {
      tt.Errorf("buf size = %d", cap(img.Pix))
    }
    if cap(d2.pix) != cap(img.Pix) {
      tt.Errorf("buf size = %d", cap(img.Pix))
    }
  })
  t.Run("getput/smallcap", func(tt *testing.T) {
    rect := image.Rect(0, 0, 100, 100)
    pool := NewImageRGBAPool(10, rect)
    img  := image.NewRGBA(rect)

    r1 := image.Rect(0, 0, 64, 36)
    i1 := image.NewRGBA(r1)
    if pool.Put(i1.Pix) {
      tt.Errorf("discard small pix")
    }

    d1 := pool.GetRef()
    if cap(d1.pix) != cap(img.Pix) {
      tt.Errorf("discard small pix = %d", cap(d1.pix))
    }
    if len(d1.pix) != len(img.Pix) {
      tt.Errorf("discard small pix")
    }
  })
  t.Run("getput/largecap", func(tt *testing.T) {
    rect := image.Rect(0, 0, 100, 100)
    pool := NewImageRGBAPool(10, rect)
    img  := image.NewRGBA(rect)

    r1 := image.Rect(0, 0, 640, 360)
    i1 := image.NewRGBA(r1)
    if pool.Put(i1.Pix) != true {
      tt.Errorf("allow large pix")
    }

    d1 := pool.GetRef()
    if cap(d1.pix) == cap(img.Pix) {
      tt.Errorf("large pix ok")
    }
    if len(d1.pix) != len(img.Pix) {
      tt.Errorf("same len")
    }
  })
}

func TestImageYCbCrPoolBufSize(t *testing.T) {
  t.Run("getsamecap", func(tt *testing.T) {
    rect := image.Rect(0, 0, 100, 100)
    i420 := image.YCbCrSubsampleRatio420
    pool := NewImageYCbCrPool(10, rect, i420)
    img  := image.NewYCbCr(rect, i420)

    c  := cap(img.Y) + cap(img.Cb) + cap(img.Cr)
    d1 := pool.GetRef()
    d2 := pool.GetRef()
    if cap(d1.pix) != c {
      tt.Errorf("buf size = %d", c)
    }
    if cap(d2.pix) != c {
      tt.Errorf("buf size = %d", c)
    }
  })
  t.Run("getput/smallcap", func(tt *testing.T) {
    rect := image.Rect(0, 0, 100, 100)
    i420 := image.YCbCrSubsampleRatio420
    pool := NewImageYCbCrPool(10, rect, i420)
    img  := image.NewYCbCr(rect, i420)

    r1 := image.Rect(0, 0, 64, 36)
    i1 := image.NewYCbCr(r1, i420)
    c1 := cap(i1.Y) + cap(i1.Cb) + cap(i1.Cr)
    v1 := make([]byte, c1)
    if pool.Put(v1) {
      tt.Errorf("discard small pix")
    }

    c  := cap(img.Y) + cap(img.Cb) + cap(img.Cr)
    l  := len(img.Y) + len(img.Cb) + len(img.Cr)
    d1 := pool.GetRef()
    if cap(d1.pix) != c {
      tt.Errorf("discard small pix = %d", cap(d1.pix))
    }
    if len(d1.pix) != l {
      tt.Errorf("discard small pix")
    }
  })
  t.Run("getput/largecap", func(tt *testing.T) {
    rect := image.Rect(0, 0, 100, 100)
    i420 := image.YCbCrSubsampleRatio420
    pool := NewImageYCbCrPool(10, rect, i420)
    img  := image.NewYCbCr(rect, i420)

    r1 := image.Rect(0, 0, 640, 360)
    i1 := image.NewYCbCr(r1, i420)
    c1 := cap(i1.Y) + cap(i1.Cb) + cap(i1.Cr)
    v1 := make([]byte, c1)
    if pool.Put(v1) != true {
      tt.Errorf("allow large pix")
    }

    c  := cap(img.Y) + cap(img.Cb) + cap(img.Cr)
    l  := len(img.Y) + len(img.Cb) + len(img.Cr)
    d1 := pool.GetRef()
    if cap(d1.pix) == c {
      if cap(d1.pix) != c1 {
        tt.Errorf("large pix ok")
      }
    }
    if len(d1.pix) != l {
      tt.Errorf("large pix ok")
    }
  })
}

func TestImageRGBAPoolCapLen(t *testing.T) {
  t.Run("getput", func(tt *testing.T) {
    r := image.Rect(0, 0, 16, 9)
    p := NewImageRGBAPool(10, r)
    if 0 != p.Len() {
      tt.Errorf("initial len 0")
    }
    if 10 != p.Cap() {
      tt.Errorf("initial cap 10")
    }

    data := p.GetRef()
    if 0 != p.Len() {
      tt.Errorf("initial len 0")
    }
    if 10 != p.Cap() {
      tt.Errorf("initial cap 10")
    }
    p.Put(data.pix)

    if 1 != p.Len() {
      tt.Errorf("put one")
    }
    if 10 != p.Cap() {
      tt.Errorf("initial cap 10")
    }

    d1 := p.GetRef()
    if 0 != p.Len() {
      tt.Errorf("aquire pool")
    }
    p.Put(d1.pix)
    if 1 != p.Len() {
      tt.Errorf("release pool")
    }
  })
  t.Run("maxcap", func(tt *testing.T) {
    r := image.Rect(0, 0, 16, 9)
    p := NewImageRGBAPool(10, r)
    s := make([]*ImageRGBARef, 0)
    for i := 0; i < 10; i += 1 {
      r := p.GetRef()
      s = append(s, r)
    }
    for _, r := range s {
      p.Put(r.pix)
    }

    if 10 != p.Len() {
      tt.Errorf("fill-ed pool: %d", p.Len())
    }
    if 10 != p.Cap() {
      tt.Errorf("max capacity = 10")
    }

    i1 := image.NewRGBA(r)
    d1 := newImageRGBARef(i1.Pix, i1, p)
    i2 := image.NewRGBA(r)
    d2 := newImageRGBARef(i2.Pix, i2, p)
    p.Put(d1.pix)
    p.Put(d2.pix)
    if 10 != p.Len() {
      tt.Errorf("fixed size pool")
    }
    if 10 != p.Cap() {
      tt.Errorf("max capacity = 10")
    }
  })
}

func TestImageYCbCrPoolCapLen(t *testing.T) {
  t.Run("getput", func(tt *testing.T) {
    r := image.Rect(0, 0, 16, 9)
    i420 := image.YCbCrSubsampleRatio420
    p := NewImageYCbCrPool(10, r, i420)
    if 0 != p.Len() {
      tt.Errorf("initial len 0")
    }
    if 10 != p.Cap() {
      tt.Errorf("initial cap 10")
    }

    data := p.GetRef()
    if 0 != p.Len() {
      tt.Errorf("initial len 0")
    }
    if 10 != p.Cap() {
      tt.Errorf("initial cap 10")
    }
    p.Put(data.pix)

    if 1 != p.Len() {
      tt.Errorf("put one")
    }
    if 10 != p.Cap() {
      tt.Errorf("initial cap 10")
    }

    d1 := p.GetRef()
    if 0 != p.Len() {
      tt.Errorf("aquire pool")
    }
    p.Put(d1.pix)
    if 1 != p.Len() {
      tt.Errorf("release pool")
    }
  })
  t.Run("maxcap", func(tt *testing.T) {
    r := image.Rect(0, 0, 16, 9)
    i420 := image.YCbCrSubsampleRatio420
    p := NewImageYCbCrPool(10, r, i420)
    s := make([]*ImageYCbCrRef, 0)
    for i := 0; i < 10; i += 1 {
      r := p.GetRef()
      s = append(s, r)
    }
    for _, r := range s {
      p.Put(r.pix)
    }

    if 10 != p.Len() {
      tt.Errorf("fill-ed pool: %d", p.Len())
    }
    if 10 != p.Cap() {
      tt.Errorf("max capacity = 10")
    }

    i1 := image.NewYCbCr(r, i420)
    v1 := make([]byte, cap(i1.Y) + cap(i1.Cb) + cap(i1.Cr)) // non ref
    d1 := newImageYCbCrRef(v1, i1, p)
    i2 := image.NewYCbCr(r, i420)
    v2 := make([]byte, cap(i2.Y) + cap(i2.Cb) + cap(i2.Cr)) // non ref
    d2 := newImageYCbCrRef(v2, i2, p)
    p.Put(d1.pix)
    p.Put(d2.pix)
    if 10 != p.Len() {
      tt.Errorf("fixed size pool")
    }
    if 10 != p.Cap() {
      tt.Errorf("max capacity = 10")
    }
  })
}

func TestYCbCrPoolPanic(t *testing.T) {
  t.Run("no panic i420", func(tt *testing.T) {
    defer func(){
      r := recover()
      if r != nil {
        tt.Errorf("no panic i420")
      }
    }()

    rect := image.Rect(0, 0, 16, 9)
    i420 := image.YCbCrSubsampleRatio420
    _ = NewImageYCbCrPool(10, rect, i420)
  })
  t.Run("panic sample != i420", func(tt *testing.T) {
    defer func(){
      r := recover()
      if r == nil {
        tt.Errorf("must panic")
      }
      s, ok := r.(string)
      if ok != true {
        tt.Errorf("panic string")
      }
      if s != notyetSupportedSampleRate {
        tt.Errorf("other panic")
      }
    }()

    rect := image.Rect(0, 0, 16, 9)
    i410 := image.YCbCrSubsampleRatio410
    _ = NewImageYCbCrPool(10, rect, i410)
  })
}

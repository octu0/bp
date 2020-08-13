package bp

import(
  "testing"
  "image"
)

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

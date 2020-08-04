package bp

import(
  "image"
)

type ImageRGBAPool struct {
  pool   chan []uint8
  rect   image.Rectangle
  width  int
  height int
  stride int
  length int
  ch     CalibrateHandler
}

func NewImageRGBAPool(poolSize int, rect image.Rectangle, funcs ...optionFunc) *ImageRGBAPool {
  opt := new(option)
  for _, fn := range funcs {
    fn(opt)
  }

  b := &ImageRGBAPool{
    pool:   make(chan []uint8, poolSize),
    rect:   rect,
    width:  rect.Dx(),
    height: rect.Dy(),
    stride: rect.Dx() * 4,
    length: rect.Dx() * rect.Dy() * 4,
    ch:     opt.calibrator,
  }

  if opt.preload {
    b.calibrate()
  }

  return b
}

func (b *ImageRGBAPool) calibrate() {
  if b.ch != nil {
    b.ch.CalibrateImageRGBAPool(b)
  }
}

func (b *ImageRGBAPool) GetRef() *ImageRGBARef {
  var pix []uint8
  select {
  case pix = <-b.pool:
    // reuse exists pool
  default:
    // create []uint8
    pix = make([]uint8, b.length)
    b.calibrate()
  }
  ref := newImageRGBARef(pix, &image.RGBA{
    Pix:    pix,
    Stride: b.stride,
    Rect:   b.rect,
  }, b)
  ref.setFinalizer()
  return ref
}

func (b *ImageRGBAPool) Put(pix []uint8) bool {
  if cap(pix) < b.length {
    // discard small buffer
    return false
  }

  select {
  case b.pool <- pix[: b.length]:
    // free capacity
    return true
  default:
    // full capacity, discard it
    return false
  }
}

func (b *ImageRGBAPool) Len() int {
  return len(b.pool)
}

func (b *ImageRGBAPool) Cap() int {
  return cap(b.pool)
}


type ImageYCbCrPool struct {
  pool     chan []uint8
  rect     image.Rectangle
  sample   image.YCbCrSubsampleRatio
  yIdx     int
  uIdx     int
  vIdx     int
  strideY  int
  strideUV int
  length   int
  ch       CalibrateHandler
}

func yuvSize(rect image.Rectangle, sample image.YCbCrSubsampleRatio) (int, int) {
  w, h := rect.Dx(), rect.Dy()
  if sample == image.YCbCrSubsampleRatio420 {
    y  := ((rect.Max.X + 1) / 2) - (rect.Min.X / 2)
    uv := ((rect.Max.Y + 1) / 2) - (rect.Min.Y / 2)
    return y, uv
  }
  // 4:4:4
  return w, h
}

func NewImageYCbCrPool(poolSize int, rect image.Rectangle, sample image.YCbCrSubsampleRatio, funcs ...optionFunc) *ImageYCbCrPool {
  opt := new(option)
  for _, fn := range funcs {
    fn(opt)
  }

  w, h  := rect.Dx(), rect.Dy()
  y, uv := yuvSize(rect, sample)

  i0 := (w * h) + (0 * y * uv)
  i1 := (w * h) + (1 * y * uv)
  i2 := (w * h) + (2 * y * uv)
  b  := &ImageYCbCrPool{
    pool:     make(chan []uint8, poolSize),
    rect:     rect,
    sample:   sample,
    yIdx:     i0,
    uIdx:     i1,
    vIdx:     i2,
    strideY:  y,
    strideUV: uv,
    length:   i2,
    ch:       opt.calibrator,
  }

  if opt.preload {
    b.calibrate()
  }

  return b
}

func (b *ImageYCbCrPool) calibrate() {
  if b.ch != nil {
    b.ch.CalibrateImageYCbCrPool(b)
  }
}

func (b *ImageYCbCrPool) GetRef() *ImageYCbCrRef {
  var pix []uint8
  select {
  case pix = <-b.pool:
    // reuse exists pool
  default:
    // create []uint8
    pix = make([]uint8, b.length)
    b.calibrate()
  }
  ref := newImageYCbCrRef(pix, &image.YCbCr{
    Y:       pix[0      : b.yIdx : b.yIdx],
    Cb:      pix[b.yIdx : b.uIdx : b.uIdx],
    Cr:      pix[b.uIdx : b.vIdx : b.vIdx],
    YStride: b.strideY,
    CStride: b.strideUV,
    Rect:    b.rect,
    SubsampleRatio: b.sample,
  }, b)
  ref.setFinalizer()
  return ref
}

func (b *ImageYCbCrPool) Put(pix []uint8) bool {
  if cap(pix) < b.length {
    // discard small buffer
    return false
  }

  select {
  case b.pool <- pix[: b.length]:
    // free capacity
    return true
  default:
    // full capacity, discard it
    return false
  }
}

func (b *ImageYCbCrPool) Len() int {
  return len(b.pool)
}

func (b *ImageYCbCrPool) Cap() int {
  return cap(b.pool)
}

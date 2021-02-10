package bp

import (
	"image"
)

const (
	notyetSupportedSampleRate string = "not yet supported sample rate"
)

type ImageRGBAPool struct {
	pool   chan []uint8
	rect   image.Rectangle
	width  int
	height int
	stride int
	length int
}

func NewImageRGBAPool(poolSize int, rect image.Rectangle, funcs ...optionFunc) *ImageRGBAPool {
	opt := newOption()
	for _, fn := range funcs {
		fn(opt)
	}

	b := &ImageRGBAPool{
		pool: make(chan []uint8, poolSize),
		// other field initialize to b.init(rect, sample)
	}
	b.init(rect)

	if opt.preload {
		b.preload(opt.preloadRate)
	}

	return b
}

func (b *ImageRGBAPool) init(rect image.Rectangle) {
	b.rect = rect
	b.width = rect.Dx()
	b.height = rect.Dy()
	b.stride = rect.Dx() * 4
	b.length = rect.Dx() * rect.Dy() * 4
}

func (b *ImageRGBAPool) createImageRGBARef(pix []uint8, pool *ImageRGBAPool) *ImageRGBARef {
	ref := newImageRGBARef(pix, &image.RGBA{
		Pix:    pix,
		Stride: b.stride,
		Rect:   b.rect,
	}, pool)
	ref.setFinalizer()
	return ref
}

func (b *ImageRGBAPool) GetRef() *ImageRGBARef {
	var pix []uint8
	select {
	case pix = <-b.pool:
		// reuse exists pool
	default:
		// create []uint8
		pix = make([]uint8, b.length)
	}
	return b.createImageRGBARef(pix, b)
}

func (b *ImageRGBAPool) preload(rate float64) {
	if 0 < cap(b.pool) {
		preloadSize := int(float64(cap(b.pool)) * rate)
		for i := 0; i < preloadSize; i += 1 {
			b.Put(make([]uint8, b.length))
		}
	}
}

func (b *ImageRGBAPool) Put(pix []uint8) bool {
	if cap(pix) < b.length {
		// discard small buffer
		return false
	}

	select {
	case b.pool <- pix[:b.length]:
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

type ImageNRGBAPool struct {
	ImageRGBAPool
}

func NewImageNRGBAPool(poolSize int, rect image.Rectangle, funcs ...optionFunc) *ImageNRGBAPool {
	opt := newOption()
	for _, fn := range funcs {
		fn(opt)
	}

	b := new(ImageNRGBAPool)
	b.pool = make(chan []uint8, poolSize)
	b.init(rect)

	if opt.preload {
		b.preload(opt.preloadRate)
	}
	return b
}

func (b *ImageNRGBAPool) createImageNRGBARef(pix []uint8, pool *ImageNRGBAPool) *ImageNRGBARef {
	ref := newImageNRGBARef(pix, &image.NRGBA{
		Pix:    pix,
		Stride: b.stride,
		Rect:   b.rect,
	}, pool)
	ref.setFinalizer()
	return ref
}

func (b *ImageNRGBAPool) GetRef() *ImageNRGBARef {
	var pix []uint8
	select {
	case pix = <-b.pool:
		// reuse exists pool
	default:
		// create []uint8
		pix = make([]uint8, b.length)
	}
	return b.createImageNRGBARef(pix, b)
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
}

func yuvSize(rect image.Rectangle, sample image.YCbCrSubsampleRatio) (int, int) {
	w, h := rect.Dx(), rect.Dy()
	if sample == image.YCbCrSubsampleRatio420 {
		y := ((rect.Max.X + 1) / 2) - (rect.Min.X / 2)
		uv := ((rect.Max.Y + 1) / 2) - (rect.Min.Y / 2)
		return y, uv
	}
	// 4:4:4
	return w, h
}

func NewImageYCbCrPool(poolSize int, rect image.Rectangle, sample image.YCbCrSubsampleRatio, funcs ...optionFunc) *ImageYCbCrPool {
	opt := newOption()
	for _, fn := range funcs {
		fn(opt)
	}

	if sample != image.YCbCrSubsampleRatio420 {
		panic(notyetSupportedSampleRate)
	}
	b := &ImageYCbCrPool{
		pool: make(chan []uint8, poolSize),
		// other field initialize to b.init(rect, sample)
	}
	b.init(rect, sample)

	if opt.preload {
		b.preload(opt.preloadRate)
	}
	return b
}
func (b *ImageYCbCrPool) init(rect image.Rectangle, sample image.YCbCrSubsampleRatio) {
	w, h := rect.Dx(), rect.Dy()
	y, uv := yuvSize(rect, sample)

	i0 := (w * h) + (0 * y * uv)
	i1 := (w * h) + (1 * y * uv)
	i2 := (w * h) + (2 * y * uv)

	b.rect = rect
	b.sample = sample
	b.yIdx = i0
	b.uIdx = i1
	b.vIdx = i2
	b.strideY = y
	b.strideUV = uv
	b.length = i2
}

func (b *ImageYCbCrPool) createImageYCbCrRef(pix []uint8, pool *ImageYCbCrPool) *ImageYCbCrRef {
	ref := newImageYCbCrRef(pix, &image.YCbCr{
		Y:              pix[0:b.yIdx:b.yIdx],
		Cb:             pix[b.yIdx:b.uIdx:b.uIdx],
		Cr:             pix[b.uIdx:b.vIdx:b.vIdx],
		YStride:        b.strideY,
		CStride:        b.strideUV,
		Rect:           b.rect,
		SubsampleRatio: b.sample,
	}, pool)
	ref.setFinalizer()
	return ref
}

func (b *ImageYCbCrPool) GetRef() *ImageYCbCrRef {
	var pix []uint8
	select {
	case pix = <-b.pool:
		// reuse exists pool
	default:
		// create []uint8
		pix = make([]uint8, b.length)
	}
	return b.createImageYCbCrRef(pix, b)
}

func (b *ImageYCbCrPool) preload(rate float64) {
	if 0 < cap(b.pool) {
		preloadSize := int(float64(cap(b.pool)) * rate)
		for i := 0; i < preloadSize; i += 1 {
			b.Put(make([]uint8, b.length))
		}
	}
}

func (b *ImageYCbCrPool) Put(pix []uint8) bool {
	if cap(pix) < b.length {
		// discard small buffer
		return false
	}

	select {
	case b.pool <- pix[:b.length]:
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

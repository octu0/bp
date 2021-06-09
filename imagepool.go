package bp

import (
	"image"
)

const (
	notyetSupportedSampleRate string = "not yet supported sample rate"
)

type ImageRGBAPool struct {
	pool   chan []byte
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
		pool: make(chan []byte, poolSize),
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
	b.stride = imageRGBAStride(rect)
	b.length = rect.Dx() * rect.Dy() * 4
}

func (b *ImageRGBAPool) createImageRGBARef(pix []byte, pool *ImageRGBAPool) *ImageRGBARef {
	ref := newImageRGBARef(pix, &image.RGBA{
		Pix:    pix,
		Stride: b.stride,
		Rect:   b.rect,
	}, pool)
	ref.setFinalizer()
	return ref
}

func (b *ImageRGBAPool) GetRef() *ImageRGBARef {
	var pix []byte
	select {
	case pix = <-b.pool:
		// reuse exists pool
	default:
		// create []byte
		pix = make([]byte, b.length)
	}
	return b.createImageRGBARef(pix, b)
}

func (b *ImageRGBAPool) preload(rate float64) {
	if 0 < cap(b.pool) {
		preloadSize := int(float64(cap(b.pool)) * rate)
		for i := 0; i < preloadSize; i += 1 {
			b.Put(make([]byte, b.length))
		}
	}
}

func (b *ImageRGBAPool) Put(pix []byte) bool {
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
	b.pool = make(chan []byte, poolSize)
	b.init(rect)

	if opt.preload {
		b.preload(opt.preloadRate)
	}
	return b
}

func (b *ImageNRGBAPool) createImageNRGBARef(pix []byte, pool *ImageNRGBAPool) *ImageNRGBARef {
	ref := newImageNRGBARef(pix, &image.NRGBA{
		Pix:    pix,
		Stride: b.stride,
		Rect:   b.rect,
	}, pool)
	ref.setFinalizer()
	return ref
}

func (b *ImageNRGBAPool) GetRef() *ImageNRGBARef {
	var pix []byte
	select {
	case pix = <-b.pool:
		// reuse exists pool
	default:
		// create []byte
		pix = make([]byte, b.length)
	}
	return b.createImageNRGBARef(pix, b)
}

type ImageYCbCrPool struct {
	pool     chan []byte
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
		cw := ((rect.Max.X + 1) / 2) - (rect.Min.X / 2)
		ch := ((rect.Max.Y + 1) / 2) - (rect.Min.Y / 2)
		return cw, ch
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
		pool: make(chan []byte, poolSize),
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
	cw, ch := yuvSize(rect, sample)

	i0 := (w * h) + (0 * cw * ch)
	i1 := (w * h) + (1 * cw * ch)
	i2 := (w * h) + (2 * cw * ch)

	b.rect = rect
	b.sample = sample
	b.yIdx = i0
	b.uIdx = i1
	b.vIdx = i2
	b.strideY = w
	b.strideUV = cw
	b.length = i2
}

func (b *ImageYCbCrPool) YStride(stride int) {
	b.strideY = stride
}

func (b *ImageYCbCrPool) UVStride(stride int) {
	b.strideUV = stride
}

func (b *ImageYCbCrPool) createImageYCbCrRef(pix []byte, pool *ImageYCbCrPool) *ImageYCbCrRef {
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
	var pix []byte
	select {
	case pix = <-b.pool:
		// reuse exists pool
	default:
		// create []byte
		pix = make([]byte, b.length)
	}
	return b.createImageYCbCrRef(pix, b)
}

func (b *ImageYCbCrPool) preload(rate float64) {
	if 0 < cap(b.pool) {
		preloadSize := int(float64(cap(b.pool)) * rate)
		for i := 0; i < preloadSize; i += 1 {
			b.Put(make([]byte, b.length))
		}
	}
}

func (b *ImageYCbCrPool) Put(pix []byte) bool {
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

func imageRGBAStride(rect image.Rectangle) int {
	return rect.Dx() * 4
}

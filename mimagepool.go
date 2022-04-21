package bp

import (
	"image"
	"sort"
)

type MultiImageRGBAPool struct {
	tuples []imagepoolTuple
	pools  []*ImageRGBAPool
}

func (b *MultiImageRGBAPool) find(r image.Rectangle) (*ImageRGBAPool, bool) {
	if r.Empty() {
		return nil, false
	}

	for i, t := range b.tuples {
		if rectIn(t.rect, r) {
			return b.pools[i], true
		}
	}
	return nil, false
}

func (b *MultiImageRGBAPool) GetRef(r image.Rectangle) *ImageRGBARef {
	if pool, ok := b.find(r); ok {
		ref := pool.GetRef()
		b.adjust(ref, r)
		return ref
	}

	pool := &ImageRGBAPool{}
	pool.init(r)

	pix := make([]uint8, pool.length)
	ref := pool.createImageRGBARef(pix, b.pools[len(b.pools)-1])
	b.adjust(ref, r)
	return ref
}

func (b *MultiImageRGBAPool) Put(pix []uint8, r image.Rectangle) bool {
	if pool, ok := b.find(r); ok {
		return pool.Put(pix)
	}
	// discard
	return false
}

func (b *MultiImageRGBAPool) adjust(ref *ImageRGBARef, r image.Rectangle) {
	ref.Img.Rect = r
	ref.Img.Stride = imageRGBAStride(r)
}

type MultiImageNRGBAPool struct {
	tuples []imagepoolTuple
	pools  []*ImageNRGBAPool
}

func NewMultiImageNRGBAPool(funcs ...multiImageBufferPoolOptionFunc) *MultiImageNRGBAPool {
	mOpt := newMultiImageBufferPoolOption()
	for _, fn := range funcs {
		fn(mOpt)
	}

	tuples := uniqImagepoolTuple(mOpt.tuples)
	sortTuples(tuples)

	pools := make([]*ImageNRGBAPool, len(tuples))
	for i, t := range tuples {
		pools[i] = NewImageNRGBAPool(t.poolSize, t.rect, mOpt.poolFuncs...)
	}
	return &MultiImageNRGBAPool{
		tuples: tuples,
		pools:  pools,
	}
}

func (b *MultiImageNRGBAPool) find(r image.Rectangle) (*ImageNRGBAPool, bool) {
	if r.Empty() {
		return nil, false
	}

	for i, t := range b.tuples {
		if rectIn(t.rect, r) {
			return b.pools[i], true
		}
	}
	return nil, false
}

func (b *MultiImageNRGBAPool) GetRef(r image.Rectangle) *ImageNRGBARef {
	if pool, ok := b.find(r); ok {
		ref := pool.GetRef()
		b.adjust(ref, r)
		return ref
	}

	pool := &ImageNRGBAPool{}
	pool.init(r)

	pix := make([]uint8, pool.length)
	ref := pool.createImageNRGBARef(pix, b.pools[len(b.pools)-1])
	b.adjust(ref, r)
	return ref
}

func (b *MultiImageNRGBAPool) Put(pix []uint8, r image.Rectangle) bool {
	if pool, ok := b.find(r); ok {
		return pool.Put(pix)
	}
	// discard
	return false
}

func (b *MultiImageNRGBAPool) adjust(ref *ImageNRGBARef, r image.Rectangle) {
	ref.Img.Rect = r
	ref.Img.Stride = imageRGBAStride(r)
}

type MultiImageYCbCrPool struct {
	tuples []imagepoolTuple
	pools  []*ImageYCbCrPool
	sample image.YCbCrSubsampleRatio
}

func NewMultiImageYCbCrPool(sample image.YCbCrSubsampleRatio, funcs ...multiImageBufferPoolOptionFunc) *MultiImageYCbCrPool {
	mOpt := newMultiImageBufferPoolOption()
	for _, fn := range funcs {
		fn(mOpt)
	}

	tuples := uniqImagepoolTuple(mOpt.tuples)
	sortTuples(tuples)

	pools := make([]*ImageYCbCrPool, len(tuples))
	for i, t := range tuples {
		pools[i] = NewImageYCbCrPool(t.poolSize, t.rect, sample, mOpt.poolFuncs...)
	}
	return &MultiImageYCbCrPool{
		tuples: tuples,
		pools:  pools,
		sample: sample,
	}
}

func (b *MultiImageYCbCrPool) find(r image.Rectangle) (*ImageYCbCrPool, bool) {
	for i, t := range b.tuples {
		if rectIn(t.rect, r) {
			return b.pools[i], true
		}
	}
	return nil, false
}

func (b *MultiImageYCbCrPool) GetRef(r image.Rectangle) *ImageYCbCrRef {
	if pool, ok := b.find(r); ok {
		ref := pool.GetRef()
		b.adjust(ref, r)
		return ref
	}

	pool := &ImageYCbCrPool{}
	pool.init(r, b.sample)

	pix := make([]uint8, pool.length)
	ref := pool.createImageYCbCrRef(pix, b.pools[len(b.pools)-1])
	b.adjust(ref, r)
	return ref
}

func (b *MultiImageYCbCrPool) Put(pix []uint8, r image.Rectangle) bool {
	if pool, ok := b.find(r); ok {
		return pool.Put(pix)
	}
	// discard
	return false
}

func (b *MultiImageYCbCrPool) adjust(ref *ImageYCbCrRef, r image.Rectangle) {
	w, h := r.Dx(), r.Dy()
	cw, ch := yuvSize(r, b.sample)

	i0 := (w * h) + (0 * cw * ch)
	i1 := (w * h) + (1 * cw * ch)
	i2 := (w * h) + (2 * cw * ch)

	ref.Img.Y = ref.pix[0:i0:i0]
	ref.Img.Cb = ref.pix[i0:i1:i1]
	ref.Img.Cr = ref.pix[i1:i2:i2]

	ref.Img.Rect = r
	ref.Img.YStride = w
	ref.Img.CStride = cw
}

type multiImageBufferPoolOptionFunc func(*multiImageBufferPoolOption)

type multiImageBufferPoolOption struct {
	tuples    []imagepoolTuple
	poolFuncs []optionFunc
}

type imagepoolTuple struct {
	poolSize int
	rect     image.Rectangle
}

func newMultiImageBufferPoolOption() *multiImageBufferPoolOption {
	return &multiImageBufferPoolOption{
		tuples:    make([]imagepoolTuple, 0),
		poolFuncs: make([]optionFunc, 0),
	}
}

func MultiImagePoolSize(poolSize int, rect image.Rectangle) multiImageBufferPoolOptionFunc {
	return func(opt *multiImageBufferPoolOption) {
		opt.tuples = append(opt.tuples, imagepoolTuple{poolSize, rect})
	}
}

func MultiImagePoolOption(funcs ...optionFunc) multiImageBufferPoolOptionFunc {
	return func(opt *multiImageBufferPoolOption) {
		opt.poolFuncs = append(opt.poolFuncs, funcs...)
	}
}

func uniqImagepoolTuple(tuples []imagepoolTuple) []imagepoolTuple {
	uniq := make(map[string]imagepoolTuple)
	for _, t := range tuples {
		if _, ok := uniq[t.rect.String()]; ok {
			continue
		}
		uniq[t.rect.String()] = t
	}
	uniqTuples := make([]imagepoolTuple, 0, len(uniq))
	for _, t := range uniq {
		uniqTuples = append(uniqTuples, imagepoolTuple{t.poolSize, t.rect})
	}
	return uniqTuples
}

func sortTuples(tuples []imagepoolTuple) {
	sort.Slice(tuples, func(a, b int) bool {
		if tuples[a].rect.Dx() == tuples[b].rect.Dx() {
			return tuples[a].rect.Dy() < tuples[b].rect.Dy()
		}
		return tuples[a].rect.Dx() < tuples[b].rect.Dx()
	})
}

func NewMultiImageRGBAPool(funcs ...multiImageBufferPoolOptionFunc) *MultiImageRGBAPool {
	mOpt := newMultiImageBufferPoolOption()
	for _, fn := range funcs {
		fn(mOpt)
	}

	tuples := uniqImagepoolTuple(mOpt.tuples)
	sortTuples(tuples)

	pools := make([]*ImageRGBAPool, len(tuples))
	for i, t := range tuples {
		pools[i] = NewImageRGBAPool(t.poolSize, t.rect, mOpt.poolFuncs...)
	}
	return &MultiImageRGBAPool{
		tuples: tuples,
		pools:  pools,
	}
}

func rectIn(src, tgt image.Rectangle) bool {
	if tgt.Dx() <= src.Dx() && tgt.Dy() <= src.Dy() {
		return true
	}
	return false
}

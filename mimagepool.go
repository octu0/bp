package bp

import (
	"image"
	"sort"
)

type multiImageBufferPoolOptionFunc func(*multiImageBufferPoolOption)

type multiImageBufferPoolOption struct {
	tuples    []imagepoolTuple
	poolFuncs []optionFunc
}

func newMultiImageBufferPoolOption() *multiImageBufferPoolOption {
	return &multiImageBufferPoolOption{
		tuples:    make([]imagepoolTuple, 0),
		poolFuncs: make([]optionFunc, 0),
	}
}

type imagepoolTuple struct {
	poolSize int
	rect     image.Rectangle
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

func rectIn(src, tgt image.Rectangle) bool {
	if tgt.Dx() <= src.Dx() && tgt.Dy() <= src.Dy() {
		return true
	}
	return false
}

type MultiImageRGBAPool struct {
	tuples []imagepoolTuple
	pools  []*ImageRGBAPool
}

func NewMultiImageRGBAPool(funcs ...multiImageBufferPoolOptionFunc) *MultiImageRGBAPool {
	mOpt := newMultiImageBufferPoolOption()
	for _, fn := range funcs {
		fn(mOpt)
	}

	tuples := uniqImagepoolTuple(mOpt.tuples)
	poolFuncs := mOpt.poolFuncs

	sortTuples(tuples)
	pools := make([]*ImageRGBAPool, len(tuples))
	for i, t := range tuples {
		pools[i] = NewImageRGBAPool(t.poolSize, t.rect, poolFuncs...)
	}
	return &MultiImageRGBAPool{
		tuples: tuples,
		pools:  pools,
	}
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
		return pool.GetRef()
	}

	pool := &ImageRGBAPool{}
	pool.init(r)

	pix := make([]uint8, pool.length)
	return pool.createImageRGBARef(pix, b.pools[len(b.pools)-1])
}

func (b *MultiImageRGBAPool) Put(pix []uint8, r image.Rectangle) bool {
	if pool, ok := b.find(r); ok {
		return pool.Put(pix)
	}
	// discard
	return false
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
	poolFuncs := mOpt.poolFuncs

	sortTuples(tuples)
	pools := make([]*ImageYCbCrPool, len(tuples))
	for i, t := range tuples {
		pools[i] = NewImageYCbCrPool(t.poolSize, t.rect, sample, poolFuncs...)
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
		return pool.GetRef()
	}

	pool := &ImageYCbCrPool{}
	pool.init(r, b.sample)

	pix := make([]uint8, pool.length)
	return pool.createImageYCbCrRef(pix, b.pools[len(b.pools)-1])
}

func (b *MultiImageYCbCrPool) Put(pix []uint8, r image.Rectangle) bool {
	if pool, ok := b.find(r); ok {
		return pool.Put(pix)
	}
	// discard
	return false
}

package bp

import (
	"bytes"
	"sort"
)

type multiBufferPoolOptionFunc func(*multiBufferPoolOption)

type multiBufferPoolOption struct {
	tuples    []bufferpoolTuple
	poolFuncs []optionFunc
}

func newMultiBufferPoolOption() *multiBufferPoolOption {
	return &multiBufferPoolOption{
		tuples:    make([]bufferpoolTuple, 0),
		poolFuncs: make([]optionFunc, 0),
	}
}

type bufferpoolTuple struct {
	poolSize, bufSize int
}

func MultiBufferPoolSize(poolSize int, bufSize int) multiBufferPoolOptionFunc {
	return func(opt *multiBufferPoolOption) {
		opt.tuples = append(opt.tuples, bufferpoolTuple{poolSize, bufSize})
	}
}

func MultiBufferPoolOption(funcs ...optionFunc) multiBufferPoolOptionFunc {
	return func(opt *multiBufferPoolOption) {
		opt.poolFuncs = append(opt.poolFuncs, funcs...)
	}
}

func uniqBufferpoolTuple(tuples []bufferpoolTuple) []bufferpoolTuple {
	uniq := make(map[int]bufferpoolTuple)
	for _, t := range tuples {
		if _, ok := uniq[t.bufSize]; ok {
			continue
		}
		uniq[t.bufSize] = t
	}
	uniqTuples := make([]bufferpoolTuple, 0, len(uniq))
	for _, t := range uniq {
		uniqTuples = append(uniqTuples, bufferpoolTuple{t.poolSize, t.bufSize})
	}
	return uniqTuples
}

type MultiBufferPool struct {
	tuples []bufferpoolTuple
	pools  []*BufferPool
}

func NewMultiBufferPool(funcs ...multiBufferPoolOptionFunc) *MultiBufferPool {
	mOpt := newMultiBufferPoolOption()
	for _, fn := range funcs {
		fn(mOpt)
	}

	tuples := uniqBufferpoolTuple(mOpt.tuples)
	poolFuncs := mOpt.poolFuncs

	sort.Slice(tuples, func(a, b int) bool {
		return tuples[a].bufSize < tuples[b].bufSize
	})
	pools := make([]*BufferPool, len(tuples))
	for i, t := range tuples {
		pools[i] = NewBufferPool(t.poolSize, t.bufSize, poolFuncs...)
	}
	return &MultiBufferPool{
		tuples: tuples,
		pools:  pools,
	}
}

func (b *MultiBufferPool) find(size int) (*BufferPool, bool) {
	for i, t := range b.tuples {
		if size <= t.bufSize {
			return b.pools[i], true
		}
	}
	return nil, false
}

func (b *MultiBufferPool) GetRef(size int) *BufferRef {
	if pool, ok := b.find(size); ok {
		data := pool.Get()
		ref := newBufferRef(data, pool)
		ref.setFinalizer()
		return ref
	}

	data := bytes.NewBuffer(make([]byte, 0, size))
	ref := newBufferRef(data, b.pools[len(b.pools)-1])
	ref.setFinalizer()
	return ref
}

func (b *MultiBufferPool) Get(size int) *bytes.Buffer {
	if pool, ok := b.find(size); ok {
		return pool.Get()
	}
	return bytes.NewBuffer(make([]byte, 0, size))
}

func (b *MultiBufferPool) Put(data *bytes.Buffer) bool {
	if pool, ok := b.find(data.Cap()); ok {
		return pool.Put(data)
	}
	// discard
	return false
}

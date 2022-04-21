//go:build aix || darwin || dragonfly || freebsd || linux || netbsd || openbsd || solaris
// +build aix darwin dragonfly freebsd linux netbsd openbsd solaris

package bp

import (
	"sort"
)

type MultiMmapBytePool struct {
	pools []*MmapBytePool
}

func (b *MultiMmapBytePool) find(size int) (*MmapBytePool, bool) {
	for i, p := range b.pools {
		if size <= p.alignSize {
			return b.pools[i], true
		}
	}
	return nil, false
}

func (b *MultiMmapBytePool) GetRef(size int) *ByteRef {
	if pool, ok := b.find(size); ok {
		data := pool.Get()
		ref := newByteRef(data[:size], pool)
		ref.setFinalizer()
		return ref
	}

	// fallback
	data := make([]byte, size)
	ref := newByteRef(data, b.pools[len(b.pools)-1])
	ref.setFinalizer()
	return ref
}

func (b *MultiMmapBytePool) Get(size int) []byte {
	if pool, ok := b.find(size); ok {
		data := pool.Get()
		return data[:size]
	}
	return make([]byte, size) // fallback
}

func (b *MultiMmapBytePool) Put(data []byte) bool {
	if pool, ok := b.find(cap(data)); ok {
		return pool.Put(data)
	}
	// discard
	return false
}

type multiMmapBytePoolOptionFunc func(*multiMmapBytePoolOption)

type multiMmapBytePoolOption struct {
	tuples    []mmapBytepoolTuple
	poolFuncs []optionFunc
}

type mmapBytepoolTuple struct {
	poolSize, bufSize int
}

func newMultiMmapBytePoolOption() *multiMmapBytePoolOption {
	return &multiMmapBytePoolOption{
		tuples:    make([]mmapBytepoolTuple, 0),
		poolFuncs: make([]optionFunc, 0),
	}
}

func MultiMmapBytePoolSize(poolSize int, bufSize int) multiMmapBytePoolOptionFunc {
	return func(opt *multiMmapBytePoolOption) {
		opt.tuples = append(opt.tuples, mmapBytepoolTuple{poolSize, bufSize})
	}
}

func MultiMmapBytePoolOption(funcs ...optionFunc) multiMmapBytePoolOptionFunc {
	return func(opt *multiMmapBytePoolOption) {
		opt.poolFuncs = append(opt.poolFuncs, funcs...)
	}
}

func uniqMmapBytepoolTuple(tuples []mmapBytepoolTuple) []mmapBytepoolTuple {
	uniq := make(map[int]mmapBytepoolTuple)
	for _, t := range tuples {
		if _, ok := uniq[t.bufSize]; ok {
			continue
		}
		uniq[t.bufSize] = t
	}
	uniqTuples := make([]mmapBytepoolTuple, 0, len(uniq))
	for _, t := range uniq {
		uniqTuples = append(uniqTuples, mmapBytepoolTuple{t.poolSize, t.bufSize})
	}
	return uniqTuples
}

func NewMultiMmapBytePool(funcs ...multiMmapBytePoolOptionFunc) *MultiMmapBytePool {
	mOpt := newMultiMmapBytePoolOption()
	for _, fn := range funcs {
		fn(mOpt)
	}

	tuples := uniqMmapBytepoolTuple(mOpt.tuples)
	poolFuncs := mOpt.poolFuncs

	pools := make([]*MmapBytePool, len(tuples))
	for i, t := range tuples {
		pools[i] = NewMmapBytePool(t.poolSize, t.bufSize, poolFuncs...)
	}
	sort.Slice(pools, func(a, b int) bool {
		return pools[a].alignSize < pools[b].alignSize
	})
	return &MultiMmapBytePool{
		pools: pools,
	}
}

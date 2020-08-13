package bp

import(
  "sort"
)

type multiBytePoolOptionFunc func(*multiBytePoolOption)

type multiBytePoolOption struct {
  tuples    []bytepoolTuple
  poolFuncs []optionFunc
}

func newMultiBytePoolOption() *multiBytePoolOption {
  return &multiBytePoolOption{
    tuples:    make([]bytepoolTuple, 0),
    poolFuncs: make([]optionFunc, 0),
  }
}

type bytepoolTuple struct {
  poolSize, bufSize int
}

func MultiBytePoolSize(poolSize int, bufSize int) multiBytePoolOptionFunc {
  return func(opt *multiBytePoolOption) {
    opt.tuples = append(opt.tuples, bytepoolTuple{poolSize, bufSize})
  }
}

func MultiBytePoolOption(funcs ...optionFunc) multiBytePoolOptionFunc {
  return func(opt *multiBytePoolOption) {
    opt.poolFuncs = append(opt.poolFuncs, funcs...)
  }
}

func uniqBytepoolTuple(tuples []bytepoolTuple) []bytepoolTuple {
  uniq := make(map[int]bytepoolTuple)
  for _, t := range tuples {
    if _, ok := uniq[t.bufSize]; ok {
      continue
    }
    uniq[t.bufSize] = t
  }
  uniqTuples := make([]bytepoolTuple, 0, len(uniq))
  for _, t := range uniq {
    uniqTuples = append(uniqTuples, bytepoolTuple{t.poolSize, t.bufSize})
  }
  return uniqTuples
}

type MultiBytePool struct {
  tuples []bytepoolTuple
  pools  []*BytePool
}

func NewMultiBytePool(funcs ...multiBytePoolOptionFunc) *MultiBytePool {
  mOpt := newMultiBytePoolOption()
  for _, fn := range funcs {
    fn(mOpt)
  }

  tuples    := uniqBytepoolTuple(mOpt.tuples)
  poolFuncs := mOpt.poolFuncs

  sort.Slice(tuples, func(a, b int) bool {
    return tuples[a].bufSize < tuples[b].bufSize
  })
  pools    := make([]*BytePool, len(tuples))
  for i, t := range tuples {
    pools[i] = NewBytePool(t.poolSize, t.bufSize, poolFuncs...)
  }
  return &MultiBytePool{
    tuples: tuples,
    pools:  pools,
  }
}

func (b *MultiBytePool) find(size int) (*BytePool, bool) {
  for i, t := range b.tuples {
    if size <= t.bufSize {
      return b.pools[i], true
    }
  }
  return nil, false
}

func (b *MultiBytePool) GetRef(size int) *ByteRef {
  if pool, ok := b.find(size); ok {
    data := pool.Get()
    ref  := newByteRef(data[:size], pool)
    ref.setFinalizer()
    return ref
  }

  data := make([]byte, size)
  ref  := newByteRef(data, b.pools[len(b.pools) - 1])
  ref.setFinalizer()
  return ref
}

func (b *MultiBytePool) Get(size int) []byte {
  if pool, ok := b.find(size); ok {
    data := pool.Get()
    return data[:size]
  }
  return make([]byte, size)
}

func (b *MultiBytePool) Put(data []byte) bool {
  if pool, ok := b.find(cap(data)); ok {
    return pool.Put(data)
  }
  // discard
  return false
}

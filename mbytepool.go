package bp

import(
  "sort"
)

type MultiBytePool struct {
  tuples []bytepoolTuple
  pools  []*BytePool
}

type bytepoolTuple struct {
  poolSize, bufSize int
}

func MultiBytePoolSize(poolSize int, bufSize int) bytepoolTuple {
  return bytepoolTuple{poolSize, bufSize}
}

func NewMultiBytePool(tuples []bytepoolTuple, funcs ...optionFunc) *MultiBytePool {
  sort.Slice(tuples, func(a, b int) bool {
    return tuples[a].bufSize < tuples[b].bufSize
  })
  pools    := make([]*BytePool, len(tuples))
  for i, t := range tuples {
    pools[i] = NewBytePool(t.poolSize, t.bufSize, funcs...)
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
    ref := newByteRef(data[:size], pool)
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

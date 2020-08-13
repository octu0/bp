package bp

type BytePool struct {
  pool    chan []byte
  bufSize int
}

func NewBytePool(poolSize int, bufSize int, funcs ...optionFunc) *BytePool {
  opt := newOption()
  for _, fn := range funcs {
    fn(opt)
  }

  b := &BytePool{
    pool:    make(chan []byte, poolSize),
    bufSize: bufSize,
  }

  if opt.preload {
    b.preload(opt.preloadRate)
  }

  return b
}

func (b *BytePool) GetRef() *ByteRef {
  data := b.Get()

  ref := newByteRef(data, b)
  ref.setFinalizer()
  return ref
}

func (b *BytePool) preload(rate float64) {
  if 0 < cap(b.pool) {
    preloadSize := int(float64(cap(b.pool)) * rate)
    for i := 0; i < preloadSize; i += 1 {
      b.Put(make([]byte, b.bufSize))
    }
  }
}

func (b *BytePool) Get() []byte {
  var data []byte
  select {
  case data = <-b.pool:
    // reuse exists pool
  default:
    // create []byte
    data = make([]byte, b.bufSize)
  }
  return data
}

func (b *BytePool) Put(data []byte) bool {
  if cap(data) < b.bufSize {
    // discard small buffer
    return false
  }

  select {
  case b.pool <- data[: b.bufSize]:
    // free capacity
    return true
  default:
    // full capacity, discard it
    return false
  }
}

func (b *BytePool) Len() int {
  return len(b.pool)
}

func (b *BytePool) Cap() int {
  return cap(b.pool)
}

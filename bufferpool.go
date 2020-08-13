package bp

import(
  "bytes"
)

type BufferPool struct {
  pool       chan *bytes.Buffer
  bufSize    int
  maxBufSize int
}

func NewBufferPool(poolSize int, bufSize int, funcs ...optionFunc) *BufferPool {
  opt := newOption()
  for _, fn := range funcs {
    fn(opt)
  }

  b := &BufferPool{
    pool:       make(chan *bytes.Buffer, poolSize),
    bufSize:    bufSize,
    maxBufSize: int(opt.maxBufSizeFactor * float64(bufSize)),
  }
  if b.maxBufSize < 1 {
    b.maxBufSize = bufSize
  }

  if opt.preload {
    b.preload(opt.preloadRate)
  }

  return b
}

func (b *BufferPool) GetRef() *BufferRef {
  data := b.Get()

  ref := newBufferRef(data, b)
  ref.setFinalizer()
  return ref
}

func (b *BufferPool) preload(rate float64) {
  if 0 < cap(b.pool) {
    preloadSize := int(float64(cap(b.pool)) * rate)
    for i := 0; i < preloadSize; i += 1 {
      b.Put(bytes.NewBuffer(make([]byte, 0, b.bufSize)))
    }
  }
}

func (b *BufferPool) Get() *bytes.Buffer {
  var data *bytes.Buffer
  select {
  case data = <-b.pool:
    // reuse exists pool
  default:
    // create *bytes.Buffer w/ []byte
    data = bytes.NewBuffer(make([]byte, 0, b.bufSize))
  }
  return data
}

func (b *BufferPool) Put(data *bytes.Buffer) bool {
  if b.maxBufSize <= data.Cap() {
    // discard, dont keep too big size buffer in heap and release it
    return false
  }

  if data.Cap() < b.bufSize {
    // increase bufSize to reduce call to internal bytes.grow
    data.Grow(b.bufSize)
  }

  data.Reset()

  select {
  case b.pool <- data:
    // free capacity
    return true
  default:
    // full capacity, discard it
    return false
  }
}

func (b *BufferPool) Len() int {
  return len(b.pool)
}

func (b *BufferPool) Cap() int {
  return cap(b.pool)
}

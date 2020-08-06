package bp

type BytePool struct {
  pool    chan []byte
  bufSize int
  ch      CalibrateHandler
}

func NewBytePool(poolSize int, bufSize int, funcs ...optionFunc) *BytePool {
  opt := new(option)
  for _, fn := range funcs {
    fn(opt)
  }

  b := &BytePool{
    pool:    make(chan []byte, poolSize),
    bufSize: bufSize,
    ch:      opt.calibrator,
  }

  if opt.preload {
    b.calibrate()
  }

  return b
}

func (b *BytePool) calibrate() {
  if b.ch != nil {
    b.ch.CalibrateBytePool(b)
  }
}

func (b *BytePool) GetRef() *ByteRef {
  data := b.Get()

  ref := newByteRef(data, b)
  ref.setFinalizer()
  return ref
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

  b.calibrate()

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

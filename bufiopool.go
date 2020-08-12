package bp

import(
  "io"
  "bufio"
)

const(
  defaultBufioSize int = 4 * 1024
)

type BufioReaderPool struct {
  pool    chan *bufio.Reader
  bufSize int
  strict  bool
  ch      CalibrateHandler
}

func NewBufioReaderPool(poolSize int, funcs ...optionFunc) *BufioReaderPool {
  return newBufioReaderPool(poolSize, defaultBufioSize, false, funcs...)
}

func NewBufioReaderSizePool(poolSize int, bufSize int, funcs ...optionFunc) *BufioReaderPool {
  return newBufioReaderPool(poolSize, bufSize, true, funcs...)
}

func newBufioReaderPool(poolSize int, bufSize int, sizeStrict bool, funcs ...optionFunc) *BufioReaderPool {
  opt := newOption()
  for _, fn := range funcs {
    fn(opt)
  }

  b := &BufioReaderPool{
    pool:    make(chan *bufio.Reader, poolSize),
    bufSize: bufSize,
    strict:  sizeStrict,
    ch:      opt.calibrator,
  }

  if opt.preload {
    b.calibrate()
  }
  return b
}

func (b *BufioReaderPool) calibrate() {
  if b.ch != nil {
    b.ch.CalibrateBufioReaderPool(b)
  }
}

func (b *BufioReaderPool) GetRef(r io.Reader) *BufioReaderRef {
  data := b.Get(r)

  ref := newBufioReaderRef(data, b)
  ref.setFinalizer()
  return ref
}

func (b *BufioReaderPool) Get(r io.Reader) *bufio.Reader {
  var br *bufio.Reader
  select {
  case br = <-b.pool:
    // reuse exists pool
  default:
    // create *bufio.Reader
    br = bufio.NewReaderSize(nil, b.bufSize)
  }
  br.Reset(r)
  return br
}

func (b *BufioReaderPool) Put(br *bufio.Reader) bool {
  br.Reset(nil)

  if br.Size() < b.bufSize {
    // discard
    return false
  }

  if b.bufSize < br.Size() {
    if b.strict {
      // discard, same buffer size only
      return false
    }
  }

  b.calibrate()

  select {
  case b.pool <- br:
    // free capacity
    return true
  default:
    // full capacity, discard it
    return false
  }
}

func (b *BufioReaderPool) Len() int {
  return len(b.pool)
}

func (b *BufioReaderPool) Cap() int {
  return cap(b.pool)
}

type BufioWriterPool struct {
  pool    chan *bufio.Writer
  bufSize int
  strict  bool
  ch      CalibrateHandler
}

func NewBufioWriterPool(poolSize int, funcs ...optionFunc) *BufioWriterPool {
  return newBufioWriterPool(poolSize, defaultBufioSize, false, funcs...)
}

func NewBufioWriterSizePool(poolSize int, bufSize int, funcs ...optionFunc) *BufioWriterPool {
  return newBufioWriterPool(poolSize, bufSize, true, funcs...)
}

func newBufioWriterPool(poolSize int, bufSize int, sizeStrict bool, funcs ...optionFunc) *BufioWriterPool {
  opt := newOption()
  for _, fn := range funcs {
    fn(opt)
  }

  b := &BufioWriterPool{
    pool:    make(chan *bufio.Writer, poolSize),
    bufSize: bufSize,
    strict:  sizeStrict,
    ch:      opt.calibrator,
  }

  if opt.preload {
    b.calibrate()
  }
  return b
}

func (b *BufioWriterPool) calibrate() {
  if b.ch != nil {
    b.ch.CalibrateBufioWriterPool(b)
  }
}

func (b *BufioWriterPool) GetRef(w io.Writer) *BufioWriterRef {
  data := b.Get(w)

  ref := newBufioWriterRef(data, b)
  ref.setFinalizer()
  return ref
}

func (b *BufioWriterPool) Get(w io.Writer) *bufio.Writer {
  var bw *bufio.Writer
  select {
  case bw = <-b.pool:
    // reuse exists pool
  default:
    // create *bufio.Writer
    bw = bufio.NewWriterSize(nil, b.bufSize)
  }
  bw.Reset(w)
  return bw
}

func (b *BufioWriterPool) Put(bw *bufio.Writer) bool {
  bw.Reset(nil)

  if bw.Size() < b.bufSize {
    // discard
    return false
  }

  if b.bufSize < bw.Size() {
    if b.strict {
      // discard, same buffer size only
      return false
    }
  }

  b.calibrate()

  select {
  case b.pool <- bw:
    // free capacity
    return true
  default:
    // full capacity, discard it
    return false
  }
}

func (b *BufioWriterPool) Len() int {
  return len(b.pool)
}

func (b *BufioWriterPool) Cap() int {
  return cap(b.pool)
}

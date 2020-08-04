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
  ch      CalibrateHandler
}

func NewBufioReaderPool(poolSize int, funcs ...optionFunc) *BufioReaderPool {
  return NewBufioReaderSizePool(poolSize, defaultBufioSize, funcs...)
}

func NewBufioReaderSizePool(poolSize int, bufSize int, funcs ...optionFunc) *BufioReaderPool {
  opt := new(option)
  for _, fn := range funcs {
    fn(opt)
  }

  b := &BufioReaderPool{
    pool:    make(chan *bufio.Reader, poolSize),
    bufSize: bufSize,
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
    b.calibrate()
  }
  br.Reset(r)
  return br
}

func (b *BufioReaderPool) Put(br *bufio.Reader) bool {
  br.Reset(nil)

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
  ch      CalibrateHandler
}

func NewBufioWriterPool(poolSize int, funcs ...optionFunc) *BufioWriterPool {
  return NewBufioWriterSizePool(poolSize, defaultBufioSize, funcs...)
}

func NewBufioWriterSizePool(poolSize int, bufSize int, funcs ...optionFunc) *BufioWriterPool {
  opt := new(option)
  for _, fn := range funcs {
    fn(opt)
  }

  b := &BufioWriterPool{
    pool:    make(chan *bufio.Writer, poolSize),
    bufSize: bufSize,
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
    b.calibrate()
  }
  bw.Reset(w)
  return bw
}

func (b *BufioWriterPool) Put(bw *bufio.Writer) bool {
  bw.Reset(nil)

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

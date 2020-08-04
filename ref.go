package bp

import(
  "runtime"
  "bytes"
  "bufio"
  "image"
)

type Ref interface {
  setFinalizer()
  Release()
}

func finalizeRef(ref Ref) {
  ref.Release()
}

// compile check
var(
  _ Ref = (*ByteRef)(nil)
  _ Ref = (*BufferRef)(nil)
  _ Ref = (*BufioReaderRef)(nil)
  _ Ref = (*BufioWriterRef)(nil)
  _ Ref = (*ImageRGBARef)(nil)
  _ Ref = (*ImageYCbCrRef)(nil)
)

type ByteRef struct {
  B      []byte
  pool   *BytePool
  closed bool
}

func newByteRef(data []byte, pool *BytePool) *ByteRef {
  return &ByteRef{
    B:      data,
    pool:   pool,
    closed: false,
  }
}

func (b *ByteRef) setFinalizer() {
  runtime.SetFinalizer(b, finalizeRef)
}

func (b *ByteRef) Release() {
  if b.closed != true {
    runtime.SetFinalizer(b, nil) // clear finalizer
    b.pool.Put(b.B)
    b.closed = true
  }
}

type BufferRef struct {
  Buf    *bytes.Buffer
  pool   *BufferPool
  closed bool
}

func newBufferRef(data *bytes.Buffer, pool *BufferPool) *BufferRef {
  return &BufferRef{
    Buf:    data,
    pool:   pool,
    closed: false,
  }
}

func (b *BufferRef) setFinalizer() {
  runtime.SetFinalizer(b, finalizeRef)
}

func (b *BufferRef) Release() {
  if b.closed != true {
    runtime.SetFinalizer(b, nil) // clear
    b.pool.Put(b.Buf)
    b.closed = true
  }
}

type BufioReaderRef struct {
  Buf   *bufio.Reader
  pool  *BufioReaderPool
  closed bool
}

func newBufioReaderRef(data *bufio.Reader, pool *BufioReaderPool) *BufioReaderRef {
  return &BufioReaderRef{
    Buf:   data,
    pool:   pool,
    closed: false,
  }
}

func (b *BufioReaderRef) setFinalizer() {
  runtime.SetFinalizer(b, finalizeRef)
}

func (b *BufioReaderRef) Release() {
  if b.closed != true {
    runtime.SetFinalizer(b, nil) // clear
    b.pool.Put(b.Buf)
    b.closed = true
  }
}

type BufioWriterRef struct {
  Buf   *bufio.Writer
  pool  *BufioWriterPool
  closed bool
}

func newBufioWriterRef(data *bufio.Writer, pool *BufioWriterPool) *BufioWriterRef {
  return &BufioWriterRef{
    Buf:   data,
    pool:   pool,
    closed: false,
  }
}

func (b *BufioWriterRef) setFinalizer() {
  runtime.SetFinalizer(b, finalizeRef)
}

func (b *BufioWriterRef) Release() {
  if b.closed != true {
    runtime.SetFinalizer(b, nil) // clear
    b.pool.Put(b.Buf)
    b.closed = true
  }
}

type ImageRGBARef struct {
  Img    *image.RGBA
  pix    []uint8
  pool   *ImageRGBAPool
  closed bool
}

func newImageRGBARef(pix []uint8, img *image.RGBA, pool *ImageRGBAPool) *ImageRGBARef {
  return &ImageRGBARef{
    Img:    img,
    pix:    pix,
    pool:   pool,
    closed: false,
  }
}

func (b *ImageRGBARef) setFinalizer() {
  runtime.SetFinalizer(b, finalizeRef)
}

func (b *ImageRGBARef) Release() {
  if b.closed != true {
    runtime.SetFinalizer(b, nil) // clear
    b.pool.Put(b.pix)
    b.Img = nil
    b.closed = true
  }
}

type ImageYCbCrRef struct {
  Img    *image.YCbCr
  pix    []uint8
  pool   *ImageYCbCrPool
  closed bool
}

func newImageYCbCrRef(pix []uint8, img *image.YCbCr, pool *ImageYCbCrPool) *ImageYCbCrRef {
  return &ImageYCbCrRef{
    Img:    img,
    pix:    pix,
    pool:   pool,
    closed: false,
  }
}

func (b *ImageYCbCrRef) setFinalizer() {
  runtime.SetFinalizer(b, finalizeRef)
}

func (b *ImageYCbCrRef) Release() {
  if b.closed != true {
    runtime.SetFinalizer(b, nil) // clear
    b.pool.Put(b.pix)
    b.Img = nil
    b.closed = true
  }
}

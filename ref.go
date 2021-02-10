package bp

import (
	"bufio"
	"bytes"
	"image"
	"runtime"
	"sync/atomic"
)

const (
	refInit   int32 = 0
	refClosed int32 = 1
)

type Ref interface {
	isClosed() bool
	setFinalizer()
	Release()
}

func finalizeRef(ref Ref) {
	ref.Release()
}

// compile check
var (
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
	closed int32
}

func newByteRef(data []byte, pool *BytePool) *ByteRef {
	return &ByteRef{
		B:      data,
		pool:   pool,
		closed: refInit,
	}
}

func (b *ByteRef) Bytes() []byte {
	return b.B
}

func (b *ByteRef) isClosed() bool {
	return atomic.LoadInt32(&b.closed) == refClosed
}

func (b *ByteRef) setFinalizer() {
	runtime.SetFinalizer(b, finalizeRef)
}

func (b *ByteRef) Release() {
	if atomic.CompareAndSwapInt32(&b.closed, refInit, refClosed) {
		runtime.SetFinalizer(b, nil) // clear finalizer
		b.pool.Put(b.B)
		b.B = nil
		b.pool = nil
	}
}

type BufferRef struct {
	Buf    *bytes.Buffer
	pool   *BufferPool
	closed int32
}

func newBufferRef(data *bytes.Buffer, pool *BufferPool) *BufferRef {
	return &BufferRef{
		Buf:    data,
		pool:   pool,
		closed: refInit,
	}
}

func (b *BufferRef) Buffer() *bytes.Buffer {
	return b.Buf
}

func (b *BufferRef) isClosed() bool {
	return atomic.LoadInt32(&b.closed) == refClosed
}

func (b *BufferRef) setFinalizer() {
	runtime.SetFinalizer(b, finalizeRef)
}

func (b *BufferRef) Release() {
	if atomic.CompareAndSwapInt32(&b.closed, refInit, refClosed) {
		runtime.SetFinalizer(b, nil) // clear
		b.pool.Put(b.Buf)
		b.Buf = nil
		b.pool = nil
	}
}

type BufioReaderRef struct {
	Buf    *bufio.Reader
	pool   *BufioReaderPool
	closed int32
}

func newBufioReaderRef(data *bufio.Reader, pool *BufioReaderPool) *BufioReaderRef {
	return &BufioReaderRef{
		Buf:    data,
		pool:   pool,
		closed: refInit,
	}
}

func (b *BufioReaderRef) Reader() *bufio.Reader {
	return b.Buf
}

func (b *BufioReaderRef) isClosed() bool {
	return atomic.LoadInt32(&b.closed) == refClosed
}

func (b *BufioReaderRef) setFinalizer() {
	runtime.SetFinalizer(b, finalizeRef)
}

func (b *BufioReaderRef) Release() {
	if atomic.CompareAndSwapInt32(&b.closed, refInit, refClosed) {
		runtime.SetFinalizer(b, nil) // clear
		b.pool.Put(b.Buf)
		b.Buf = nil
		b.pool = nil
	}
}

type BufioWriterRef struct {
	Buf    *bufio.Writer
	pool   *BufioWriterPool
	closed int32
}

func newBufioWriterRef(data *bufio.Writer, pool *BufioWriterPool) *BufioWriterRef {
	return &BufioWriterRef{
		Buf:    data,
		pool:   pool,
		closed: refInit,
	}
}

func (b *BufioWriterRef) Writer() *bufio.Writer {
	return b.Buf
}

func (b *BufioWriterRef) isClosed() bool {
	return atomic.LoadInt32(&b.closed) == refClosed
}

func (b *BufioWriterRef) setFinalizer() {
	runtime.SetFinalizer(b, finalizeRef)
}

func (b *BufioWriterRef) Release() {
	if atomic.CompareAndSwapInt32(&b.closed, refInit, refClosed) {
		runtime.SetFinalizer(b, nil) // clear
		b.pool.Put(b.Buf)
		b.Buf = nil
		b.pool = nil
	}
}

type ImageRGBARef struct {
	Img    *image.RGBA
	pix    []uint8
	pool   *ImageRGBAPool
	closed int32
}

func newImageRGBARef(pix []uint8, img *image.RGBA, pool *ImageRGBAPool) *ImageRGBARef {
	return &ImageRGBARef{
		Img:    img,
		pix:    pix,
		pool:   pool,
		closed: refInit,
	}
}

func (b *ImageRGBARef) Image() *image.RGBA {
	return b.Img
}

func (b *ImageRGBARef) isClosed() bool {
	return atomic.LoadInt32(&b.closed) == refClosed
}

func (b *ImageRGBARef) setFinalizer() {
	runtime.SetFinalizer(b, finalizeRef)
}

func (b *ImageRGBARef) Release() {
	if atomic.CompareAndSwapInt32(&b.closed, refInit, refClosed) {
		runtime.SetFinalizer(b, nil) // clear
		b.pool.Put(b.pix)
		b.Img = nil
		b.pix = nil
		b.pool = nil
	}
}

type ImageNRGBARef struct {
	Img    *image.NRGBA
	pix    []uint8
	pool   *ImageNRGBAPool
	closed int32
}

func newImageNRGBARef(pix []uint8, img *image.NRGBA, pool *ImageNRGBAPool) *ImageNRGBARef {
	return &ImageNRGBARef{
		Img:    img,
		pix:    pix,
		pool:   pool,
		closed: refInit,
	}
}

func (b *ImageNRGBARef) Image() *image.NRGBA {
	return b.Img
}

func (b *ImageNRGBARef) isClosed() bool {
	return atomic.LoadInt32(&b.closed) == refClosed
}

func (b *ImageNRGBARef) setFinalizer() {
	runtime.SetFinalizer(b, finalizeRef)
}

func (b *ImageNRGBARef) Release() {
	if atomic.CompareAndSwapInt32(&b.closed, refInit, refClosed) {
		runtime.SetFinalizer(b, nil) // clear
		b.pool.Put(b.pix)
		b.Img = nil
		b.pix = nil
		b.pool = nil
	}
}

type ImageYCbCrRef struct {
	Img    *image.YCbCr
	pix    []uint8
	pool   *ImageYCbCrPool
	closed int32
}

func newImageYCbCrRef(pix []uint8, img *image.YCbCr, pool *ImageYCbCrPool) *ImageYCbCrRef {
	return &ImageYCbCrRef{
		Img:    img,
		pix:    pix,
		pool:   pool,
		closed: refInit,
	}
}

func (b *ImageYCbCrRef) Image() *image.YCbCr {
	return b.Img
}

func (b *ImageYCbCrRef) isClosed() bool {
	return atomic.LoadInt32(&b.closed) == refClosed
}

func (b *ImageYCbCrRef) setFinalizer() {
	runtime.SetFinalizer(b, finalizeRef)
}

func (b *ImageYCbCrRef) Release() {
	if atomic.CompareAndSwapInt32(&b.closed, refInit, refClosed) {
		runtime.SetFinalizer(b, nil) // clear
		b.pool.Put(b.pix)
		b.Img = nil
		b.pix = nil
		b.pool = nil
	}
}

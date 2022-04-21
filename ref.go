package bp

import (
	"bufio"
	"bytes"
	"image"
	"runtime"
	"sync/atomic"
	"time"
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
	runtime.SetFinalizer(ref, nil) // clear finalizer
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
	_ Ref = (*TickerRef)(nil)
	_ Ref = (*TimerRef)(nil)
)

type ByteRef struct {
	B      []byte
	pool   ByteGetPut
	closed int32
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
		b.pool.Put(b.B)
	}
}

func newByteRef(data []byte, pool ByteGetPut) *ByteRef {
	return &ByteRef{
		B:      data,
		pool:   pool,
		closed: refInit,
	}
}

type BufferRef struct {
	Buf    *bytes.Buffer
	pool   BytesBufferGetPut
	closed int32
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
		b.pool.Put(b.Buf)
	}
}

func newBufferRef(data *bytes.Buffer, pool BytesBufferGetPut) *BufferRef {
	return &BufferRef{
		Buf:    data,
		pool:   pool,
		closed: refInit,
	}
}

type BufioReaderRef struct {
	Buf    *bufio.Reader
	pool   BufioReaderGetPut
	closed int32
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
		b.pool.Put(b.Buf)
	}
}

func newBufioReaderRef(data *bufio.Reader, pool BufioReaderGetPut) *BufioReaderRef {
	return &BufioReaderRef{
		Buf:    data,
		pool:   pool,
		closed: refInit,
	}
}

type BufioWriterRef struct {
	Buf    *bufio.Writer
	pool   BufioWriterGetPut
	closed int32
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
		b.pool.Put(b.Buf)
	}
}

func newBufioWriterRef(data *bufio.Writer, pool BufioWriterGetPut) *BufioWriterRef {
	return &BufioWriterRef{
		Buf:    data,
		pool:   pool,
		closed: refInit,
	}
}

type ImageRGBARef struct {
	Img    *image.RGBA
	pix    []byte
	pool   ImageRGBAGetPut
	closed int32
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
		b.pool.Put(b.pix)
	}
}

func newImageRGBARef(pix []byte, img *image.RGBA, pool ImageRGBAGetPut) *ImageRGBARef {
	return &ImageRGBARef{
		Img:    img,
		pix:    pix,
		pool:   pool,
		closed: refInit,
	}
}

type ImageNRGBARef struct {
	Img    *image.NRGBA
	pix    []byte
	pool   ImageNRGBAGetPut
	closed int32
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
		b.pool.Put(b.pix)
	}
}

func newImageNRGBARef(pix []byte, img *image.NRGBA, pool ImageNRGBAGetPut) *ImageNRGBARef {
	return &ImageNRGBARef{
		Img:    img,
		pix:    pix,
		pool:   pool,
		closed: refInit,
	}
}

type ImageYCbCrRef struct {
	Img    *image.YCbCr
	pix    []byte
	pool   ImageYCbCrGetPut
	closed int32
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
		b.pool.Put(b.pix)
	}
}

func newImageYCbCrRef(pix []byte, img *image.YCbCr, pool ImageYCbCrGetPut) *ImageYCbCrRef {
	return &ImageYCbCrRef{
		Img:    img,
		pix:    pix,
		pool:   pool,
		closed: refInit,
	}
}

type TickerRef struct {
	T      *time.Ticker
	pool   TickerGetPut
	closed int32
}

func (b *TickerRef) Ticker() *time.Ticker {
	return b.T
}

func (b *TickerRef) isClosed() bool {
	return atomic.LoadInt32(&b.closed) == refClosed
}

func (b *TickerRef) setFinalizer() {
	runtime.SetFinalizer(b, finalizeRef)
}

func (b *TickerRef) Release() {
	if atomic.CompareAndSwapInt32(&b.closed, refInit, refClosed) {
		b.pool.Put(b.T)
	}
}

func newTickerRef(ticker *time.Ticker, pool TickerGetPut) *TickerRef {
	return &TickerRef{
		T:      ticker,
		pool:   pool,
		closed: refInit,
	}
}

type TimerRef struct {
	T      *time.Timer
	pool   TimerGetPut
	closed int32
}

func (b *TimerRef) Ticker() *time.Timer {
	return b.T
}

func (b *TimerRef) isClosed() bool {
	return atomic.LoadInt32(&b.closed) == refClosed
}

func (b *TimerRef) setFinalizer() {
	runtime.SetFinalizer(b, finalizeRef)
}

func (b *TimerRef) Release() {
	if atomic.CompareAndSwapInt32(&b.closed, refInit, refClosed) {
		b.pool.Put(b.T)
	}
}

func newTimerRef(timer *time.Timer, pool TimerGetPut) *TimerRef {
	return &TimerRef{
		T:      timer,
		pool:   pool,
		closed: refInit,
	}
}

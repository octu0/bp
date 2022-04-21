package bp

import (
	"bufio"
	"bytes"
	"io"
	"time"
)

type PoolSize interface {
	Len() int
	Cap() int
}

type ByteGetPut interface {
	GetRef() *ByteRef
	Get() []byte
	Put([]byte) bool
}

type BytesBufferGetPut interface {
	GetRef() *BufferRef
	Get() *bytes.Buffer
	Put(*bytes.Buffer) bool
}

type BufioReaderGetPut interface {
	GetRef(io.Reader) *BufioReaderRef
	Get(io.Reader) *bufio.Reader
	Put(*bufio.Reader) bool
}

type BufioWriterGetPut interface {
	GetRef(io.Writer) *BufioWriterRef
	Get(io.Writer) *bufio.Writer
	Put(*bufio.Writer) bool
}

type ImageRGBAGetPut interface {
	GetRef() *ImageRGBARef
	Put([]byte) bool
}

type ImageNRGBAGetPut interface {
	GetRef() *ImageNRGBARef
	Put([]byte) bool
}

type ImageYCbCrGetPut interface {
	GetRef() *ImageYCbCrRef
	Put([]byte) bool
}

type TickerGetPut interface {
	GetRef(time.Duration) *TickerRef
	Get(time.Duration) *time.Ticker
	Put(*time.Ticker) bool
}

type TimerGetPut interface {
	GetRef(time.Duration) *TimerRef
	Get(time.Duration) *time.Timer
	Put(*time.Timer) bool
}

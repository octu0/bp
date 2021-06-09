package bp

import (
	"bufio"
	"bytes"
	"io"
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

# `bp`

[![Apache License](https://img.shields.io/github/license/octu0/bp)](https://github.com/octu0/bp/blob/master/LICENSE)
[![GoDoc](https://godoc.org/github.com/octu0/bp?status.svg)](https://godoc.org/github.com/octu0/bp)
[![Go Report Card](https://goreportcard.com/badge/github.com/octu0/bp)](https://goreportcard.com/report/github.com/octu0/bp)
[![Releases](https://img.shields.io/github/v/release/octu0/bp)](https://github.com/octu0/chanque/bp)

`bp` implements buffer pool of various objects such as byte array (`[]byte`) or [*bytes.Buffer](http://golang.org/pkg/bytes/#Buffer) / [*image.RGBA](https://golang.org/pkg/image/#RGBA) and [*bufio.Reader](https://golang.org/pkg/bufio/#Reader).  
It is inspired by [bpool](https://github.com/oxtoacart/bpool) and its many features are similar.

`bp` provides the following pool types
- `bp.BufferPool` which provides fixed-size pool of [*bytes.Buffers](http://golang.org/pkg/bytes/#Buffer)
- `bp.BytePool` which provides fixed-size pool of `[]byte` slice 
- `bp.BufioReaderPool` which provides fixed-size pool of [*bufio.Reader](https://golang.org/pkg/bufio/#Reader)
- `bp.BufioWriterPool` which provides fixed-size pool of [*bufio.Writer](https://golang.org/pkg/bufio/#Writer)
- `bp.ImageRGBAPool` which provides fixed-size pool of [*image.RGBA](https://golang.org/pkg/image/#RGBA) 
- `bp.ImageYCbCrPool` which provides fixed-size pool of [*image.YCbCr](https://golang.org/pkg/image/#YCbCr) 

It also provides a MultiPool to bundle multiple pools:

- MultiBytePool
- MultiBufferPool
- MultiImageRGBAPool
- MultiImageYCbCrPool

In addition, `bp` provides an easy to manipulate object interface to prevent forgetting to put it back into the pool

- `bp.ByteRef`
- `bp.BufferRef`
- `bp.BufioReaderRef`
- `bp.BufioWriterRef`
- `bp.ImageRGBARef`
- `bp.ImageYCbCrRef`

## Installation

```bash
go get github.com/octu0/bp
```

## Example

Here's a quick example for using `bp.BufferPool`. We create a pool of the desired size, call the `Get()` method to obtain a buffer for use, and call `Put(buf)` to return the buffer to the pool.

```go
var (
  bufpool := bp.NewBufferPool(1000, 128) // capacity 1000 items, each buffer initial 128 Byte pre-sized
)

func main() {
  // get buffer from pool
  buf := bufpool.Get()
  ...
  ...
  // return buffer to pool
  bufpool.Put(buf)
}
```

## License

Apache 2.0, see LICENSE file for details.

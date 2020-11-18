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

## io.Copy vs CopyIOPool.Copy Benchmark

CopyIOPool.Copy to reduce allocation of [io.Copy](https://golang.org/pkg/io/#Copy)

```bash
$ go test -bench=BenchmarkIoCopy -benchmem ./
goos: darwin
goarch: amd64
pkg: github.com/octu0/bp
BenchmarkIoCopy-8                	  528212	      2270 ns/op	   32816 B/op	       3 allocs/op
BenchmarkIoCopyPoolDefault-8     	  662061	      1888 ns/op	   16624 B/op	       7 allocs/op
BenchmarkIoCopyPoolFixedSize-8   	 2311072	       521 ns/op	      48 B/op	       2 allocs/op
PASS
ok  	github.com/octu0/bp	4.244s
```

## ioutil.ReadAll vs CopyIOPool.ReadAll Benchmark

similarly, CopyIOPool.ReadAll reduces allocation of [ioutil.ReadAll](https://golang.org/pkg/io/ioutil/#ReadAll)

```bash
$ go test -bench=BenchmarkIoReadAll -benchmem ./
goos: darwin
goarch: amd64
pkg: github.com/octu0/bp
BenchmarkIoReadAllIoUtil-8          	    3774	    315602 ns/op	 4192308 B/op	      13 allocs/op
BenchmarkIoReadAllPoolDefault-8     	    3688	    356025 ns/op	 4063475 B/op	      13 allocs/op
BenchmarkIoReadAllPoolFixedSize-8   	    3208	    348119 ns/op	 4046906 B/op	       8 allocs/op
PASS
ok  	github.com/octu0/bp	3.761s
```

## License

Apache 2.0, see LICENSE file for details.

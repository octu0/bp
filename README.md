# `bp`

[![Apache License](https://img.shields.io/github/license/octu0/bp)](https://github.com/octu0/bp/blob/master/LICENSE)
[![GoDoc](https://godoc.org/github.com/octu0/bp?status.svg)](https://godoc.org/github.com/octu0/bp)
[![Go Report Card](https://goreportcard.com/badge/github.com/octu0/bp)](https://goreportcard.com/report/github.com/octu0/bp)
[![Releases](https://img.shields.io/github/v/release/octu0/bp)](https://github.com/octu0/bp/releases)

`bp` implements buffer pool of various objects such as byte array (`[]byte`) or [*bytes.Buffer](http://golang.org/pkg/bytes/#Buffer) / [*image.RGBA](https://golang.org/pkg/image/#RGBA) and [*bufio.Reader](https://golang.org/pkg/bufio/#Reader).  
It is inspired by [bpool](https://github.com/oxtoacart/bpool) and its many features are similar.

`bp` provides the following pool types
- `bp.BufferPool` which provides fixed-size pool of [*bytes.Buffers](http://golang.org/pkg/bytes/#Buffer)
- `bp.BytePool` which provides fixed-size pool of `[]byte` slice 
- `bp.MmapBytePool` Same as BytePool, but uses mmap to allocate the slices
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

# Benchmark

## `bytes.Buffer`: sync.Pool vs BufferPool

```bash
$ go test -v -bench=BenchmarkBufferPool -benchmem ./
goos: darwin
goarch: amd64
pkg: github.com/octu0/bp
BenchmarkBufferPool/default/8-8         	 1000000	      1119 ns/op	      32 B/op	       1 allocs/op
BenchmarkBufferPool/default/4096-8      	  965864	      1289 ns/op	      32 B/op	       1 allocs/op
BenchmarkBufferPool/syncpool/8-8        	  991191	      1206 ns/op	      48 B/op	       1 allocs/op
BenchmarkBufferPool/syncpool/4096-8     	  967411	      1265 ns/op	      48 B/op	       1 allocs/op
BenchmarkBufferPool/bufferpool/8-8      	  898627	      1361 ns/op	      48 B/op	       1 allocs/op
BenchmarkBufferPool/bufferpool/4096-8   	  811107	      1430 ns/op	      48 B/op	       1 allocs/op
PASS
ok  	github.com/octu0/bp	7.274s
```

## `[]byte`: sync.Pool vs BytePool

```bash
$ go test -v -bench=BenchmarkBytePool -benchmem ./
goos: darwin
goarch: amd64
pkg: github.com/octu0/bp
BenchmarkBytePool/default/8-8         	 1000000	      1064 ns/op	      16 B/op	       1 allocs/op
BenchmarkBytePool/default/4096-8      	  902665	      1132 ns/op	      16 B/op	       1 allocs/op
BenchmarkBytePool/syncpool/8-8        	 1000000	      1165 ns/op	      64 B/op	       2 allocs/op
BenchmarkBytePool/syncpool/4096-8     	 1000000	      1059 ns/op	      64 B/op	       2 allocs/op
BenchmarkBytePool/bytepool/8-8        	 1027453	      1246 ns/op	      32 B/op	       1 allocs/op
BenchmarkBytePool/bytepool/4096-8     	  935844	      1299 ns/op	      32 B/op	       1 allocs/op
PASS
ok  	github.com/octu0/bp	7.878s
```

## `bufio.Reader`: sync.Pool vs BufioReaderPool

```bash
$ go test -v -bench=BenchmarkBufioReaderPool -benchmem ./
goos: darwin
goarch: amd64
pkg: github.com/octu0/bp
BenchmarkBufioReaderPool/default/8-8         	 1000000	      1190 ns/op	    1056 B/op	       3 allocs/op
BenchmarkBufioReaderPool/default/4096-8      	  766180	      1500 ns/op	    5136 B/op	       3 allocs/op
BenchmarkBufioReaderPool/syncpool/8-8        	 1000000	      1339 ns/op	    1056 B/op	       2 allocs/op
BenchmarkBufioReaderPool/syncpool/4096-8     	  896155	      1337 ns/op	    1058 B/op	       2 allocs/op
BenchmarkBufioReaderPool/bufiopool/8-8       	  831057	      1485 ns/op	    1168 B/op	       4 allocs/op
BenchmarkBufioReaderPool/bufiopool/4096-8    	  836048	      1356 ns/op	    1056 B/op	       2 allocs/op
PASS
ok  	github.com/octu0/bp	7.354s
```

## `image.RGBA`: sync.Pool vs `ImageRGBAPool`

```bash
$ go test -v -bench=BenchmarkImageRGBAPool -benchmem ./
goos: darwin
goarch: amd64
pkg: github.com/octu0/bp
BenchmarkImageRGBAPool/default/360-8         	   41775	     28822 ns/op	  925764 B/op	       2 allocs/op
BenchmarkImageRGBAPool/default/1080-8        	    2203	    506796 ns/op	 8294806 B/op	       2 allocs/op
BenchmarkImageRGBAPool/syncpool/360-8        	 1000000	      1070 ns/op	      28 B/op	       1 allocs/op
BenchmarkImageRGBAPool/syncpool/1080-8       	 1019760	      1137 ns/op	      89 B/op	       1 allocs/op
BenchmarkImageRGBAPool/imagepool/360-8       	  818752	      1382 ns/op	     138 B/op	       3 allocs/op
BenchmarkImageRGBAPool/imagepool/1080-8      	  976969	      1448 ns/op	     137 B/op	       3 allocs/op
PASS
ok  	github.com/octu0/bp	9.034s
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

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
$ go test -bench=BenchmarkBufferPool -benchmem ./
goos: darwin
goarch: amd64
pkg: github.com/octu0/bp
cpu: Intel(R) Core(TM) i7-8569U CPU @ 2.80GHz
BenchmarkBufferPool/default/8-8         	 1224789	       990.2 ns/op	      32 B/op	       1 allocs/op
BenchmarkBufferPool/default/4096-8      	 1000000	      1103 ns/op	      32 B/op	       1 allocs/op
BenchmarkBufferPool/syncpool/8-8        	 1000000	      1055 ns/op	      48 B/op	       1 allocs/op
BenchmarkBufferPool/syncpool/4096-8     	 1000000	      1108 ns/op	      48 B/op	       1 allocs/op
BenchmarkBufferPool/bufferpool/8-8      	 1000000	      1155 ns/op	      48 B/op	       1 allocs/op
BenchmarkBufferPool/bufferpool/4096-8   	  981954	      1189 ns/op	      48 B/op	       1 allocs/op
PASS
ok  	github.com/octu0/bp	7.899s
```

## `[]byte`: sync.Pool vs BytePool

```bash
$ go test -bench=BenchmarkBytePool -benchmem ./
goos: darwin
goarch: amd64
pkg: github.com/octu0/bp
cpu: Intel(R) Core(TM) i7-8569U CPU @ 2.80GHz
BenchmarkBytePool/default/8-8         	 1294598	       941.2 ns/op	      16 B/op	       1 allocs/op
BenchmarkBytePool/default/4096-8      	 1000000	      1031 ns/op	      16 B/op	       1 allocs/op
BenchmarkBytePool/syncpool/8-8        	 1000000	      1069 ns/op	      48 B/op	       2 allocs/op
BenchmarkBytePool/syncpool/4096-8     	 1000000	      1075 ns/op	      48 B/op	       2 allocs/op
BenchmarkBytePool/bytepool/8-8        	 1213732	      1023 ns/op	      24 B/op	       1 allocs/op
BenchmarkBytePool/bytepool/4096-8     	 1000000	      1127 ns/op	      24 B/op	       1 allocs/op
PASS
ok  	github.com/octu0/bp	8.780s
```

## `bufio.Reader`: sync.Pool vs BufioReaderPool

```bash
$ go test -bench=BenchmarkBufioReaderPool -benchmem ./
goos: darwin
goarch: amd64
pkg: github.com/octu0/bp
cpu: Intel(R) Core(TM) i7-8569U CPU @ 2.80GHz
BenchmarkBufioReaderPool/default/8-8         	 1000000	      1120 ns/op	    1056 B/op	       3 allocs/op
BenchmarkBufioReaderPool/default/4096-8      	  803418	      1542 ns/op	    5136 B/op	       3 allocs/op
BenchmarkBufioReaderPool/syncpool/8-8        	 1000000	      1102 ns/op	    1048 B/op	       2 allocs/op
BenchmarkBufioReaderPool/syncpool/4096-8     	 1000000	      1110 ns/op	    1051 B/op	       2 allocs/op
BenchmarkBufioReaderPool/bufiopool/8-8       	 1000000	      1290 ns/op	    1160 B/op	       4 allocs/op
BenchmarkBufioReaderPool/bufiopool/4096-8    	  918162	      1279 ns/op	    1048 B/op	       2 allocs/op
PASS
ok  	github.com/octu0/bp	7.147s
```

## `image.RGBA`: sync.Pool vs `ImageRGBAPool`

```bash
$ go test -bench=BenchmarkImageRGBAPool -benchmem ./
goos: darwin
goarch: amd64
pkg: github.com/octu0/bp
cpu: Intel(R) Core(TM) i7-8569U CPU @ 2.80GHz
BenchmarkImageRGBAPool/default/360-8         	   38256	     29911 ns/op	  925739 B/op	       3 allocs/op
BenchmarkImageRGBAPool/default/1080-8        	    1995	    591596 ns/op	 8286139 B/op	       3 allocs/op
BenchmarkImageRGBAPool/syncpool/360-8        	 1316624	      1042 ns/op	      26 B/op	       1 allocs/op
BenchmarkImageRGBAPool/syncpool/1080-8       	 1000000	      1098 ns/op	      82 B/op	       1 allocs/op
BenchmarkImageRGBAPool/imagepool/360-8       	 1006983	      1326 ns/op	     153 B/op	       3 allocs/op
BenchmarkImageRGBAPool/imagepool/1080-8      	  871188	      1343 ns/op	     154 B/op	       3 allocs/op
PASS
ok  	github.com/octu0/bp	10.637s
```

## io.Copy vs CopyIOPool.Copy Benchmark

CopyIOPool.Copy to reduce allocation of [io.Copy](https://golang.org/pkg/io/#Copy)

```bash
$ go test -bench=BenchmarkIoCopy -benchmem ./
goos: darwin
goarch: amd64
pkg: github.com/octu0/bp
cpu: Intel(R) Core(TM) i7-8569U CPU @ 2.80GHz
BenchmarkIoCopy-8                	  452194	      2650 ns/op	   32816 B/op	       3 allocs/op
BenchmarkIoCopyPoolDefault-8     	  575596	      2064 ns/op	   16608 B/op	       7 allocs/op
BenchmarkIoCopyPoolFixedSize-8   	 2274831	       522.1 ns/op	      48 B/op	       2 allocs/op
PASS
ok  	github.com/octu0/bp	4.189s
```

## ioutil.ReadAll vs CopyIOPool.ReadAll Benchmark

similarly, CopyIOPool.ReadAll reduces allocation of [ioutil.ReadAll](https://golang.org/pkg/io/ioutil/#ReadAll)

```bash
$ go test -bench=BenchmarkIoReadAll -benchmem ./
goos: darwin
goarch: amd64
pkg: github.com/octu0/bp
cpu: Intel(R) Core(TM) i7-8569U CPU @ 2.80GHz
BenchmarkIoReadAllIoUtil-8          	    2628	    455815 ns/op	 5862972 B/op	      30 allocs/op
BenchmarkIoReadAllPoolDefault-8     	    3057	    378599 ns/op	 4063444 B/op	      13 allocs/op
BenchmarkIoReadAllPoolFixedSize-8   	    3180	    378923 ns/op	 4046892 B/op	       8 allocs/op
PASS
ok  	github.com/octu0/bp	3.718s
```

## License

Apache 2.0, see LICENSE file for details.

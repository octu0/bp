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
- `bp.CopyIOPool` which provides fixed-size pool of [io.CopyBuffer](https://golang.org/pkg/io#CopyBuffer) and [io.ReadAll](https://golang.org/pkg/io#ReadAll)
- `bp.TickerPool` which provides fixed-size pool of [*time.Ticker](https://golang.org/pkg/time#Ticker)
- `bp.TimerPool` which provides fixed-size pool of [*time.Timer](https://golang.org/pkg/time#Timer)

It also provides a MultiPool to bundle multiple pools:

- MultiBytePool
- MultiMmapBytePool
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
BenchmarkBufferPool/default/8-8         	 1608904	       768.9 ns/op	      32 B/op	       1 allocs/op
BenchmarkBufferPool/default/4096-8      	 1421576	       791.1 ns/op	      32 B/op	       1 allocs/op
BenchmarkBufferPool/syncpool/8-8        	 1584180	       743.3 ns/op	      48 B/op	       1 allocs/op
BenchmarkBufferPool/syncpool/4096-8     	 1505594	       798.4 ns/op	      48 B/op	       1 allocs/op
BenchmarkBufferPool/bufferpool/8-8      	 1439310	       917.2 ns/op	      48 B/op	       1 allocs/op
BenchmarkBufferPool/bufferpool/4096-8   	 1225413	       967.4 ns/op	      48 B/op	       1 allocs/op
PASS
ok  	github.com/octu0/bp	12.309s
```

## `[]byte`: sync.Pool vs BytePool

```bash
$ go test -bench=BenchmarkBytePool -benchmem ./
goos: darwin
goarch: amd64
pkg: github.com/octu0/bp
cpu: Intel(R) Core(TM) i7-8569U CPU @ 2.80GHz
BenchmarkBytePool/default/8-8         	 1827867	       653.8 ns/op	      16 B/op	       1 allocs/op
BenchmarkBytePool/default/4096-8      	 1562638	       788.3 ns/op	      16 B/op	       1 allocs/op
BenchmarkBytePool/syncpool/8-8        	 1643428	       763.2 ns/op	      48 B/op	       2 allocs/op
BenchmarkBytePool/syncpool/4096-8     	 1586283	       803.8 ns/op	      48 B/op	       2 allocs/op
BenchmarkBytePool/bytepool/8-8        	 1357020	       904.1 ns/op	      24 B/op	       1 allocs/op
BenchmarkBytePool/bytepool/4096-8     	 1359921	       846.6 ns/op	      24 B/op	       1 allocs/op
PASS
ok  	github.com/octu0/bp	12.105s
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
BenchmarkImageRGBAPool/default/360-8         	   31308	     36206 ns/op	  925814 B/op	       3 allocs/op
BenchmarkImageRGBAPool/default/1080-8        	    4639	    821211 ns/op	 8282518 B/op	       3 allocs/op
BenchmarkImageRGBAPool/syncpool/360-8        	 1657608	       684.7 ns/op	      20 B/op	       1 allocs/op
BenchmarkImageRGBAPool/syncpool/1080-8       	 1635321	       696.2 ns/op	      56 B/op	       1 allocs/op
BenchmarkImageRGBAPool/imagepool/360-8       	 1000000	      1190 ns/op	     151 B/op	       3 allocs/op
BenchmarkImageRGBAPool/imagepool/1080-8      	 1000000	      1105 ns/op	     151 B/op	       3 allocs/op
PASS
ok  	github.com/octu0/bp	11.502s
```

## io.Copy vs CopyIOPool.Copy Benchmark

CopyIOPool.Copy to reduce allocation of [io.Copy](https://golang.org/pkg/io/#Copy)

```bash
$ go test -bench=BenchmarkIoCopy -benchmem ./
goos: darwin
goarch: amd64
pkg: github.com/octu0/bp
cpu: Intel(R) Core(TM) i7-8569U CPU @ 2.80GHz
BenchmarkIoCopy-8                	  461892	      2637 ns/op	   32816 B/op	       3 allocs/op
BenchmarkIoCopyPoolDefault-8     	  619404	      1982 ns/op	   16608 B/op	       7 allocs/op
BenchmarkIoCopyPoolFixedSize-8   	 2858724	       410.9 ns/op	      48 B/op	       2 allocs/op
PASS
ok  	github.com/octu0/bp	4.126s
```

## ioutil.ReadAll vs CopyIOPool.ReadAll Benchmark

similarly, CopyIOPool.ReadAll reduces allocation of [io.ReadAll](https://golang.org/pkg/io#ReadAll)

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

## `*time.Ticker`: sync.Pool vs `TickerPool`

```bash
$ go test -bench=BenchmarkTickerPool -benchmem ./
goos: darwin
goarch: amd64
pkg: github.com/octu0/bp
cpu: Intel(R) Core(TM) i7-8569U CPU @ 2.80GHz
BenchmarkTickerPool/default-8         	  114298	     10731 ns/op	     224 B/op	       4 allocs/op
BenchmarkTickerPool/syncpool-8        	  110748	     10863 ns/op	      32 B/op	       1 allocs/op
BenchmarkTickerPool/pool-8            	  110414	     10867 ns/op	      32 B/op	       1 allocs/op
PASS
ok  	github.com/octu0/bp	3.994s
```

## `*time.Timer`: sync.Pool vs `TimerPool`

```bash
$ go test -bench=BenchmarkTimerPool -benchmem ./
goos: darwin
goarch: amd64
pkg: github.com/octu0/bp
cpu: Intel(R) Core(TM) i7-8569U CPU @ 2.80GHz
BenchmarkTimerPool/default-8         	  115437	     10889 ns/op	     224 B/op	       4 allocs/op
BenchmarkTimerPool/syncpool-8        	  107302	     11099 ns/op	      32 B/op	       1 allocs/op
BenchmarkTimerPool/pool-8            	  106790	     11061 ns/op	      32 B/op	       1 allocs/op
PASS
ok  	github.com/octu0/bp	4.002s
```

## License

Apache 2.0, see LICENSE file for details.

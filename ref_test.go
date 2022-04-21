package bp

import (
	"bufio"
	"bytes"
	"image"
	"strings"
	"testing"
	"time"

	"github.com/octu0/chanque"
)

func BenchmarkRef(b *testing.B) {
	r4096 := strings.NewReader(strings.Repeat("@", 4096))
	w4096 := []byte(strings.Repeat("@", 4096))

	b.Run("byte", func(tb *testing.B) {
		tb.Run("Get", func(bb *testing.B) {
			e := chanque.NewExecutor(10, 10)
			defer e.Release()

			p := NewBytePool(e.MaxWorker(), 4096)
			for i := 0; i < tb.N; i += 1 {
				e.Submit(func() {
					s := p.Get()
					defer p.Put(s)
					r4096.Read(s)
				})
			}
		})
		tb.Run("GetRef", func(bb *testing.B) {
			e := chanque.NewExecutor(10, 10)
			defer e.Release()

			p := NewBytePool(e.MaxWorker(), 8)
			for i := 0; i < tb.N; i += 1 {
				e.Submit(func() {
					s := p.GetRef()
					defer s.Release()
					r4096.Read(s.B)
				})
			}
		})
	})
	b.Run("buffer", func(tb *testing.B) {
		tb.Run("Get", func(bb *testing.B) {
			e := chanque.NewExecutor(10, 10)
			defer e.Release()

			p := NewBufferPool(e.MaxWorker(), 4096)
			for i := 0; i < tb.N; i += 1 {
				e.Submit(func() {
					s := p.Get()
					defer p.Put(s)
					s.Write(w4096)
				})
			}
		})
		tb.Run("GetRef", func(bb *testing.B) {
			e := chanque.NewExecutor(10, 10)
			defer e.Release()

			p := NewBufferPool(e.MaxWorker(), 4096)
			for i := 0; i < tb.N; i += 1 {
				e.Submit(func() {
					s := p.GetRef()
					defer s.Release()
					s.Buf.Write(w4096)
				})
			}
		})
	})
}

func TestRefValue(t *testing.T) {
	t.Run("byte", func(tt *testing.T) {
		b := newByteRef([]byte{10, 123, 4, 56}, nil)
		if bytes.Equal(b.Bytes(), []byte{10, 123, 4, 56}) != true {
			tt.Errorf("same value")
		}
	})
	t.Run("buffer", func(tt *testing.T) {
		bf := bytes.NewBuffer(nil)
		bf.Write([]byte("hello"))

		b := newBufferRef(bf, nil)
		d := b.Buffer()
		if bytes.Equal(d.Bytes(), []byte("hello")) != true {
			tt.Errorf("same value")
		}
		if d.Len() != bf.Len() {
			tt.Errorf("same len")
		}
		if d.Cap() != bf.Cap() {
			tt.Errorf("same cap")
		}
	})
	t.Run("bufio.Reader", func(tt *testing.T) {
		b := newBufioReaderRef(bufio.NewReaderSize(nil, 123), nil)
		d := b.Reader()
		if d.Size() != 123 {
			tt.Errorf("same size")
		}
	})
	t.Run("bufio.Writer", func(tt *testing.T) {
		b := newBufioWriterRef(bufio.NewWriterSize(nil, 123), nil)
		d := b.Writer()
		if d.Size() != 123 {
			tt.Errorf("same size")
		}
	})
	t.Run("image.RGBA", func(tt *testing.T) {
		img := image.NewRGBA(image.Rect(0, 0, 10, 10))
		b := newImageRGBARef(img.Pix, img, nil)

		d := b.Image()
		if d.Bounds().Eq(image.Rect(0, 0, 10, 10)) != true {
			tt.Errorf("same size")
		}
		if bytes.Equal(d.Pix, img.Pix) != true {
			tt.Errorf("same value")
		}
	})
	t.Run("image.NRGBA", func(tt *testing.T) {
		img := image.NewNRGBA(image.Rect(0, 0, 10, 10))
		b := newImageNRGBARef(img.Pix, img, nil)

		d := b.Image()
		if d.Bounds().Eq(image.Rect(0, 0, 10, 10)) != true {
			tt.Errorf("same size")
		}
		if bytes.Equal(d.Pix, img.Pix) != true {
			tt.Errorf("same value")
		}
	})
	t.Run("image.YCbCr", func(tt *testing.T) {
		img := image.NewYCbCr(image.Rect(0, 0, 10, 10), image.YCbCrSubsampleRatio420)
		c := cap(img.Y) + cap(img.Cb) + cap(img.Cr)
		v := make([]byte, c)
		copy(v[0:cap(img.Y)], img.Y)
		copy(v[cap(img.Y):cap(img.Y)+cap(img.Cb)], img.Cb)
		copy(v[cap(img.Y)+cap(img.Cb):cap(img.Y)+cap(img.Cb)+cap(img.Cr)], img.Cr)
		b := newImageYCbCrRef(v, img, nil)

		d := b.Image()
		if d.Bounds().Eq(image.Rect(0, 0, 10, 10)) != true {
			tt.Errorf("same size")
		}
		if bytes.Equal(v[0:cap(img.Y)], img.Y) != true {
			tt.Errorf("same value")
		}
		if bytes.Equal(v[cap(img.Y):cap(img.Y)+cap(img.Cb)], img.Cb) != true {
			tt.Errorf("same value")
		}
		if bytes.Equal(v[cap(img.Y)+cap(img.Cb):cap(img.Y)+cap(img.Cb)+cap(img.Cr)], img.Cr) != true {
			tt.Errorf("same value")
		}
	})
}

func TestRefRelease(t *testing.T) {
	testRelease := func(tt *testing.T, r Ref) {
		r.Release()
		if r.isClosed() != true {
			tt.Errorf("closed")
		}
		r.Release()
		if r.isClosed() != true {
			tt.Errorf("double close ok")
		}
		r.setFinalizer()
		r.Release()
		if r.isClosed() != true {
			tt.Errorf("setFinalizer after call ok")
		}
	}
	t.Run("byte", func(tt *testing.T) {
		p := NewBytePool(1, 1)
		b := newByteRef([]byte{}, p)
		testRelease(tt, b)
	})
	t.Run("byte_mmap", func(tt *testing.T) {
		p := NewMmapBytePool(1, 1)
		b := newByteRef([]byte{}, p)
		testRelease(tt, b)
	})
	t.Run("buffer", func(tt *testing.T) {
		p := NewBufferPool(1, 1)
		bf := bytes.NewBuffer(nil)
		bf.Write([]byte{})
		b := newBufferRef(bf, p)
		testRelease(tt, b)
	})
	t.Run("bufio.Reader", func(tt *testing.T) {
		p := NewBufioReaderPool(1)
		b := newBufioReaderRef(bufio.NewReaderSize(nil, 123), p)
		testRelease(tt, b)
	})
	t.Run("bufio.Writer", func(tt *testing.T) {
		p := NewBufioWriterPool(1)
		b := newBufioWriterRef(bufio.NewWriterSize(nil, 123), p)
		testRelease(tt, b)
	})
	t.Run("image.RGBA", func(tt *testing.T) {
		img := image.NewRGBA(image.Rect(0, 0, 10, 10))
		p := NewImageRGBAPool(1, image.Rect(0, 0, 10, 10))
		b := newImageRGBARef(img.Pix, img, p)
		testRelease(tt, b)
	})
	t.Run("image.NRGBA", func(tt *testing.T) {
		img := image.NewNRGBA(image.Rect(0, 0, 10, 10))
		p := NewImageNRGBAPool(1, image.Rect(0, 0, 10, 10))
		b := newImageNRGBARef(img.Pix, img, p)
		testRelease(tt, b)
	})
	t.Run("image.YCbCrPool", func(tt *testing.T) {
		img := image.NewYCbCr(image.Rect(0, 0, 10, 10), image.YCbCrSubsampleRatio420)
		p := NewImageYCbCrPool(1, image.Rect(0, 0, 10, 10), image.YCbCrSubsampleRatio420)
		b := newImageYCbCrRef([]byte{}, img, p)
		testRelease(tt, b)
	})
	t.Run("time.Ticker", func(tt *testing.T) {
		p := NewTickerPool(1)
		b := newTickerRef(time.NewTicker(1*time.Millisecond), p)
		testRelease(tt, b)
	})
	t.Run("time.Timer", func(tt *testing.T) {
		p := NewTimerPool(1)
		b := newTimerRef(time.NewTimer(1*time.Millisecond), p)
		testRelease(tt, b)
	})
}

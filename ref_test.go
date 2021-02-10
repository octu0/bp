package bp

import (
	"bufio"
	"bytes"
	"image"
	"testing"
)

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
}

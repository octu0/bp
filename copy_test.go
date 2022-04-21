package bp

import (
	"bytes"
	"io"
	"testing"
)

type testWriter struct {
	io.Writer
}

func newTestWriter() *testWriter {
	return new(testWriter)
}

func (w *testWriter) Write(data []byte) (int, error) {
	return len(data), nil
}

type testReader struct {
	offset, size int
	io.Reader
}

func newTestReader() *testReader {
	return &testReader{
		offset: 0,
		size:   1024 * 1024,
	}
}

func (r *testReader) Read(b []byte) (int, error) {
	if r.offset == r.size {
		return 0, io.EOF
	}
	remain := r.size - r.offset
	if len(b) < remain {
		remain = len(b)
	}
	r.offset += remain
	return remain, nil
}

func BenchmarkIoCopy(b *testing.B) {
	for i := 0; i < b.N; i += 1 {
		w := newTestWriter()
		r := newTestReader()
		io.Copy(w, r)
	}
}

func BenchmarkIoCopyPoolDefault(b *testing.B) {
	for i := 0; i < b.N; i += 1 {
		w := newTestWriter()
		r := newTestReader()
		Copy(w, r)
	}
}

func BenchmarkIoCopyPoolFixedSize(b *testing.B) {
	cp := NewCopyIOPool(100, 16*1024)
	for i := 0; i < b.N; i += 1 {
		w := newTestWriter()
		r := newTestReader()
		cp.Copy(w, r)
	}
}

func BenchmarkIoReadAllIoUtil(b *testing.B) {
	for i := 0; i < b.N; i += 1 {
		r := newTestReader()
		_, _ = io.ReadAll(r)
	}
}

func BenchmarkIoReadAllPoolDefault(b *testing.B) {
	for i := 0; i < b.N; i += 1 {
		r := newTestReader()
		_, _ = ReadAll(r)
	}
}

func BenchmarkIoReadAllPoolFixedSize(b *testing.B) {
	cp := NewCopyIOPool(100, 16*1024)
	for i := 0; i < b.N; i += 1 {
		r := newTestReader()
		_, _ = cp.ReadAll(r)
	}
}

func TestCopyIOPoolCopy(t *testing.T) {
	t.Run("existsbuf", func(tt *testing.T) {
		dst := bytes.NewBuffer([]byte("hello"))
		src := bytes.NewReader([]byte("world"))

		tt.Logf("prev dst: '%s'", dst.Bytes())
		cp := NewCopyIOPool(10, 4)
		cp.Copy(dst, src)
		tt.Logf("curr dst: '%s'", dst.Bytes())

		if string(dst.Bytes()) != "helloworld" {
			tt.Errorf("no append? '%s'", dst.Bytes())
		}
	})
	t.Run("small", func(tt *testing.T) {
		src := bytes.Repeat([]byte("a"), 1000)
		dst := bytes.NewBuffer(nil)

		cp := NewCopyIOPool(10, 1)
		cp.Copy(dst, bytes.NewReader(src))
		if bytes.Equal(dst.Bytes(), src) != true {
			tt.Errorf("not same bytes")
		}
	})
	t.Run("large", func(tt *testing.T) {
		src := bytes.Repeat([]byte("a"), 1000)
		dst := bytes.NewBuffer(nil)

		cp := NewCopyIOPool(10, 4*1024*1024)
		cp.Copy(dst, bytes.NewReader(src))
		if bytes.Equal(dst.Bytes(), src) != true {
			tt.Errorf("not same bytes")
		}
	})
	t.Run("same", func(tt *testing.T) {
		src := bytes.Repeat([]byte("a"), 1000)
		dst := bytes.NewBuffer(nil)

		cp := NewCopyIOPool(10, 1000)
		cp.Copy(dst, bytes.NewReader(src))
		if bytes.Equal(dst.Bytes(), src) != true {
			tt.Errorf("not same bytes")
		}
	})
}

func TestCopyIOPoolReadAll(t *testing.T) {
	t.Run("str", func(tt *testing.T) {
		src := []byte("helloworld")

		cp := NewCopyIOPool(10, 4)
		data, err := cp.ReadAll(bytes.NewReader(src))
		if err != nil {
			tt.Fatalf(err.Error())
		}
		tt.Logf("readAll = '%s'", data)
		if string(data) != "helloworld" {
			tt.Errorf("ReadAll bytes?: '%s'", data)
		}
	})
	t.Run("small", func(tt *testing.T) {
		src := bytes.Repeat([]byte("b"), 1000)

		cp := NewCopyIOPool(10, 1)
		data, err := cp.ReadAll(bytes.NewReader(src))
		if err != nil {
			tt.Fatalf(err.Error())
		}
		if bytes.Equal(data, src) != true {
			tt.Errorf("not same bytes")
		}
	})
	t.Run("large", func(tt *testing.T) {
		src := bytes.Repeat([]byte("b"), 1000)

		cp := NewCopyIOPool(10, 4*1024*1024)
		data, err := cp.ReadAll(bytes.NewReader(src))
		if err != nil {
			tt.Fatalf(err.Error())
		}
		if bytes.Equal(data, src) != true {
			tt.Errorf("not same bytes")
		}
	})
	t.Run("same", func(tt *testing.T) {
		src := bytes.Repeat([]byte("b"), 1000)

		cp := NewCopyIOPool(10, 1000)
		data, err := cp.ReadAll(bytes.NewReader(src))
		if err != nil {
			tt.Fatalf(err.Error())
		}
		if bytes.Equal(data, src) != true {
			tt.Errorf("not same bytes")
		}
	})
}

package bp

import (
	"bytes"
	"errors"
	"io"
)

const (
	defaultCopyIOSize int = 16 * 1024
)

var (
	ErrIOReadNagativeRead = errors.New("negative count from io.Read")
)

type CopyIOPool struct {
	pool *BytePool
}

func (c *CopyIOPool) Copy(dst io.Writer, src io.Reader) (int64, error) {
	buf := c.pool.Get()
	defer c.pool.Put(buf)

	return io.CopyBuffer(dst, src, buf)
}

func (c *CopyIOPool) ReadAll(src io.Reader) ([]byte, error) {
	buf := c.pool.Get()
	defer c.pool.Put(buf)

	size := int64(0)
	out := bytes.NewBuffer(make([]byte, 0, c.pool.bufSize))
	for {
		n, err := src.Read(buf)
		if n < 0 {
			return []byte{}, ErrIOReadNagativeRead
		}
		size += int64(n)
		if err == io.EOF {
			return out.Bytes(), nil
		}
		if err != nil {
			return []byte{}, err
		}
		out.Write(buf[:n])
	}
}

func (c *CopyIOPool) Len() int {
	return c.pool.Len()
}

func (c *CopyIOPool) Cap() int {
	return c.pool.Cap()
}

func NewCopyIOPool(poolSize int, bufSize int, funcs ...optionFunc) *CopyIOPool {
	return &CopyIOPool{
		pool: NewBytePool(poolSize, bufSize, funcs...),
	}
}

func Copy(dst io.Writer, src io.Reader) (int64, error) {
	c := NewCopyIOPool(1, defaultCopyIOSize)
	return c.Copy(dst, src)
}

func ReadAll(src io.Reader) ([]byte, error) {
	c := NewCopyIOPool(1, defaultCopyIOSize)
	return c.ReadAll(src)
}

// +build aix darwin dragonfly freebsd linux netbsd openbsd solaris

package bp

import (
	"golang.org/x/sys/unix"
)

const (
	DefaultMmapAlignment int = 8
	mmapPerm                 = unix.PROT_READ | unix.PROT_WRITE
	mmapFlag                 = unix.MAP_ANON | unix.MAP_PRIVATE
)

func mmapAlign(size int, align int) int {
	return ((size + align) >> 3) << 3
}

func defaultMmapAlign(size int) int {
	// default aligment 32-byte
	return mmapAlign(size, DefaultMmapAlignment)
}

type MmapBytePool struct {
	pool      chan []byte
	bufSize   int
	alignSize int
}

func NewMmapBytePool(poolSize, bufSize int, funcs ...optionFunc) *MmapBytePool {
	opt := newOption()
	for _, fn := range funcs {
		fn(opt)
	}

	b := &MmapBytePool{
		pool:      make(chan []byte, poolSize),
		bufSize:   bufSize,
		alignSize: defaultMmapAlign(bufSize),
	}

	if opt.preload {
		b.preload(opt.preloadRate)
	}

	return b
}

func (b *MmapBytePool) preload(rate float64) {
	if 0 < cap(b.pool) {
		preloadSize := int(float64(cap(b.pool)) * rate)
		buffers, err := unix.Mmap(-1, 0, b.alignSize*preloadSize, mmapPerm, mmapFlag)
		if err != nil {
			buffers = make([]byte, b.alignSize*preloadSize) // fallback
		}
		for 0 < len(buffers) {
			b.Put(buffers[:b.bufSize:b.alignSize])
			buffers = buffers[b.alignSize:]
		}
	}
}

func (b *MmapBytePool) GetRef() *ByteRef {
	data := b.Get()

	ref := newByteRef(data, b)
	ref.setFinalizer()
	return ref
}

func (b *MmapBytePool) Get() []byte {
	select {
	case data := <-b.pool:
		// reuse exists pool
		return data[:b.bufSize]
	default:
		// create from mmap
		buf, err := unix.Mmap(-1, 0, b.alignSize, mmapPerm, mmapFlag)
		if err != nil {
			buf = make([]byte, b.alignSize) // fallback
		}
		return buf[:b.bufSize]
	}
}

func (b *MmapBytePool) Put(data []byte) bool {
	if cap(data) != b.alignSize {
		// discard
		return false
	}

	select {
	case b.pool <- data[:b.bufSize]:
		// free capacity
		return true
	default:
		// full capacity, discard it
		unix.Munmap(data)
		return false
	}
}

func (b *MmapBytePool) Len() int {
	return len(b.pool)
}

func (b *MmapBytePool) Cap() int {
	return cap(b.pool)
}

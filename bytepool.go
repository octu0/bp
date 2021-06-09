package bp

type BytePool struct {
	pool       chan []byte
	bufSize    int
	maxBufSize int
}

func NewBytePool(poolSize int, bufSize int, funcs ...optionFunc) *BytePool {
	opt := newOption()
	for _, fn := range funcs {
		fn(opt)
	}

	b := &BytePool{
		pool:       make(chan []byte, poolSize),
		bufSize:    bufSize,
		maxBufSize: int(opt.maxBufSizeFactor * float64(bufSize)),
	}

	if b.maxBufSize < 1 {
		b.maxBufSize = bufSize
	}

	if opt.preload {
		b.preload(opt.preloadRate)
	}

	return b
}

func (b *BytePool) GetRef() *ByteRef {
	data := b.Get()

	ref := newByteRef(data, b)
	ref.setFinalizer()
	return ref
}

func (b *BytePool) preload(rate float64) {
	if 0 < cap(b.pool) {
		preloadSize := int(float64(cap(b.pool)) * rate)
		for i := 0; i < preloadSize; i += 1 {
			b.Put(make([]byte, b.bufSize))
		}
	}
}

func (b *BytePool) Get() []byte {
	select {
	case data := <-b.pool:
		// reuse exists pool
		return data[:b.bufSize]
	default:
		// create []byte
		return make([]byte, b.bufSize)
	}
}

func (b *BytePool) Put(data []byte) bool {
	if b.maxBufSize < cap(data) {
		// discard, dont keep too big size byte in heap and release it
		return false
	}

	if cap(data) < b.bufSize {
		// discard small buffer
		return false
	}

	select {
	case b.pool <- data[:b.bufSize:b.bufSize]:
		// free capacity
		return true
	default:
		// full capacity, discard it
		return false
	}
}

func (b *BytePool) Len() int {
	return len(b.pool)
}

func (b *BytePool) Cap() int {
	return cap(b.pool)
}

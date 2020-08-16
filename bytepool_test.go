package bp

import (
	"testing"
)

func TestBytePoolBufSize(t *testing.T) {
	bufSize := 8
	t.Run("getsamecap", func(tt *testing.T) {
		p := NewBytePool(10, bufSize)
		d1 := p.Get()
		d2 := p.Get()
		if cap(d1) != bufSize {
			tt.Errorf("buf size = %d", bufSize)
		}
		if cap(d2) != bufSize {
			tt.Errorf("buf size = %d", bufSize)
		}
	})
	t.Run("getput/samecap", func(tt *testing.T) {
		p := NewBytePool(10, bufSize)
		d1 := p.Get()
		p.Put(d1)
		d2 := p.Get()
		p.Put(d2)

		d3 := p.Get()
		d4 := p.Get()
		if cap(d3) != bufSize {
			tt.Errorf("buf size = %d", bufSize)
		}
		if cap(d4) != bufSize {
			tt.Errorf("buf size = %d", bufSize)
		}
	})
	t.Run("getput/largecap", func(tt *testing.T) {
		p := NewBytePool(10, bufSize)
		p.Put(make([]byte, 123))

		d1 := p.Get()
		if cap(d1) != bufSize {
			tt.Errorf("discard over max buf 123 byte")
		}
		if len(d1) != bufSize {
			tt.Errorf("discard over max buf 123 byte")
		}
	})
}

func TestBytePoolDiscard(t *testing.T) {
	bufSize := 8
	t.Run("freecap/samesize", func(tt *testing.T) {
		p := NewBytePool(10, bufSize)

		d1 := make([]byte, bufSize)
		if p.Put(d1) != true {
			tt.Errorf("freecap %d", p.Cap())
		}
	})
	t.Run("fullcap/samesize", func(tt *testing.T) {
		p := NewBytePool(2, bufSize)
		p.Put(make([]byte, bufSize))
		p.Put(make([]byte, bufSize))

		d1 := make([]byte, bufSize)
		if p.Put(d1) {
			tt.Errorf("fulled capacity %d", p.Cap())
		}
	})
	t.Run("freecap/largesize", func(tt *testing.T) {
		p := NewBytePool(10, bufSize)
		if p.Put(make([]byte, bufSize)) != true {
			tt.Errorf("put ok")
		}

		if p.Put(make([]byte, bufSize+1)) != true {
			tt.Errorf("put ok nearby size")
		}

		if p.Put(make([]byte, 100)) {
			tt.Errorf("discard over max buf size")
		}
	})
	t.Run("freecap/smallsize", func(tt *testing.T) {
		p := NewBytePool(10, bufSize)
		if p.Put(make([]byte, bufSize)) != true {
			tt.Errorf("put ok")
		}

		if p.Put(make([]byte, 1)) {
			tt.Errorf("discard small buf")
		}
	})
}

func TestBytePoolLenCap(t *testing.T) {
	t.Run("getput", func(tt *testing.T) {
		bufSize := 1
		p := NewBytePool(10, bufSize)
		if 0 != p.Len() {
			tt.Errorf("initial len 0")
		}
		if 10 != p.Cap() {
			tt.Errorf("initial cap 10")
		}

		data := p.Get()
		if 0 != p.Len() {
			tt.Errorf("initial len 0")
		}
		if 10 != p.Cap() {
			tt.Errorf("initial cap 10")
		}
		p.Put(data)

		if 1 != p.Len() {
			tt.Errorf("put one")
		}
		if 10 != p.Cap() {
			tt.Errorf("initial cap 10")
		}

		d1 := p.Get()
		if 0 != p.Len() {
			tt.Errorf("aquire pool")
		}
		p.Put(d1)
		if 1 != p.Len() {
			tt.Errorf("release pool")
		}
	})
	t.Run("maxcap", func(tt *testing.T) {
		bufSize := 1
		p := NewBytePool(10, bufSize)
		s := make([][]byte, 0)
		for i := 0; i < 10; i += 1 {
			d := p.Get()
			s = append(s, d)
		}
		for _, d := range s {
			p.Put(d)
		}

		if 10 != p.Len() {
			tt.Errorf("fill-ed pool: %d", p.Len())
		}
		if 10 != p.Cap() {
			tt.Errorf("max capacity = 10")
		}

		d1 := make([]byte, bufSize)
		d2 := make([]byte, bufSize)
		p.Put(d1)
		p.Put(d2)

		if 10 != p.Len() {
			tt.Errorf("fixed size pool")
		}
		if 10 != p.Cap() {
			tt.Errorf("max capacity = 10")
		}
	})
}

func TestBytePoolPreload(t *testing.T) {
	p := NewBytePool(12, 8, Preload(true))
	l := int(float64(12) * defaultPreloadRate)
	if p.Len() != l {
		t.Errorf("preloaded buffer = %d", p.Len())
	}
}

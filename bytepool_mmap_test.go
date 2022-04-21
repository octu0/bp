//go:build aix || darwin || dragonfly || freebsd || linux || netbsd || openbsd || solaris
// +build aix darwin dragonfly freebsd linux netbsd openbsd solaris

package bp

import (
	"runtime"
	"strings"
	"testing"

	"github.com/octu0/chanque"
)

func BenchmarkMmapBytePool(b *testing.B) {
	k8 := strings.NewReader(strings.Repeat("@", 8))
	k4096 := strings.NewReader(strings.Repeat("@", 4096))
	run := func(name string, fn func(*testing.B)) {
		m1 := new(runtime.MemStats)
		runtime.ReadMemStats(m1)

		b.Run(name, fn)

		m2 := new(runtime.MemStats)
		runtime.ReadMemStats(m2)
		/*
			b.Logf(
				"%-20s\tTotalAlloc=%5d\tStackInUse=%5d",
				name,
				int64(m2.TotalAlloc)-int64(m1.TotalAlloc),
				int64(m2.StackInuse)-int64(m1.StackInuse),
				//int64(m2.HeapSys)  - int64(m1.HeapSys),
				//int64(m2.HeapIdle)   - int64(m1.HeapIdle),
			)
		*/
	}
	run("heap/8", func(tb *testing.B) {
		e := chanque.NewExecutor(10, 10)
		defer e.Release()

		p := NewBytePool(e.MaxWorker(), 8)

		for i := 0; i < tb.N; i += 1 {
			e.Submit(func() {
				s := p.Get()
				k8.Read(s)
				p.Put(s)
			})
		}
	})
	run("heap/4096", func(tb *testing.B) {
		e := chanque.NewExecutor(10, 10)
		defer e.Release()

		p := NewBytePool(e.MaxWorker(), 4096)

		for i := 0; i < tb.N; i += 1 {
			e.Submit(func() {
				s := p.Get()
				k4096.Read(s)
				p.Put(s)
			})
		}
	})
	run("mmap/8", func(tb *testing.B) {
		e := chanque.NewExecutor(10, 10)
		defer e.Release()

		p := NewMmapBytePool(e.MaxWorker(), 8)

		for i := 0; i < tb.N; i += 1 {
			e.Submit(func() {
				s := p.Get()
				k8.Read(s)
				p.Put(s)
			})
		}
	})
	run("mmap/4096", func(tb *testing.B) {
		e := chanque.NewExecutor(10, 10)
		defer e.Release()

		p := NewMmapBytePool(e.MaxWorker(), 4096)

		for i := 0; i < tb.N; i += 1 {
			e.Submit(func() {
				s := p.Get()
				k4096.Read(s)
				p.Put(s)
			})
		}
	})
}

func TestMmapBytePoolBufSize(t *testing.T) {
	bufSize := 12
	t.Run("getsamecap", func(tt *testing.T) {
		p := NewMmapBytePool(10, bufSize)
		d1 := p.Get()
		d2 := p.Get()
		if len(d1) != bufSize {
			tt.Errorf("buf size = %d", bufSize)
		}
		if len(d2) != bufSize {
			tt.Errorf("buf size = %d", bufSize)
		}
	})
	t.Run("getput/samecap", func(tt *testing.T) {
		p := NewMmapBytePool(10, bufSize)
		d1 := p.Get()
		p.Put(d1)
		d2 := p.Get()
		p.Put(d2)

		d3 := p.Get()
		d4 := p.Get()
		if len(d3) != bufSize {
			tt.Errorf("buf size = %d", bufSize)
		}
		if len(d4) != bufSize {
			tt.Errorf("buf size = %d", bufSize)
		}
	})
	t.Run("getput/largecap", func(tt *testing.T) {
		p := NewMmapBytePool(10, bufSize)
		p.Put(make([]byte, 123))

		if p.Len() != 0 {
			tt.Errorf("discard over max buf 123 byte")
		}

		d1 := p.Get()
		tt.Logf("len=%d cap=%d", len(d1), cap(d1))
	})
	t.Run("getput/samealigncap", func(tt *testing.T) {
		align := defaultMmapAlign(bufSize)
		p := NewMmapBytePool(10, bufSize)
		p.Put(make([]byte, align))

		if p.Len() != 1 {
			tt.Errorf("no discard. same alignment size")
		}

		d1 := p.Get()
		tt.Logf("len=%d cap=%d", len(d1), cap(d1))
	})
}

func TestMmapBytePoolDiscard(t *testing.T) {
	bufSize := 7
	t.Run("freecap/samesize", func(tt *testing.T) {
		p := NewMmapBytePool(10, bufSize)
		a := defaultMmapAlign(bufSize)

		d1 := make([]byte, bufSize)
		if p.Put(d1) {
			tt.Errorf("discard, not alignment size: %d", a)
		}
		if p.Len() != 0 {
			tt.Errorf("discard, not alignment size: %d", a)
		}
	})
	t.Run("fullcap/samealignsize", func(tt *testing.T) {
		p := NewMmapBytePool(2, bufSize)
		a := defaultMmapAlign(bufSize)
		p.Put(make([]byte, a))
		p.Put(make([]byte, a))

		if p.Len() != 2 {
			tt.Errorf("same alignsize putable")
		}

		d1 := make([]byte, a)
		if p.Put(d1) {
			tt.Errorf("fulled capacity %d", p.Cap())
		}
	})
	t.Run("freecap/largesize", func(tt *testing.T) {
		p := NewMmapBytePool(10, bufSize)
		a := defaultMmapAlign(bufSize)

		if p.Put(make([]byte, bufSize+5)) {
			tt.Errorf("%d <> %d", bufSize+5, a)
		}

		if p.Put(make([]byte, 100)) {
			tt.Errorf("same align size only: %d", a)
		}
	})
}

func TestMmapBytePoolLenCap(t *testing.T) {
	t.Run("getput", func(tt *testing.T) {
		bufSize := 1
		p := NewMmapBytePool(10, bufSize)
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
			tt.Errorf("acquire pool")
		}
		p.Put(d1)
		if 1 != p.Len() {
			tt.Errorf("release pool")
		}
	})
	t.Run("maxcap", func(tt *testing.T) {
		bufSize := 1
		p := NewMmapBytePool(10, bufSize)
		s := make([][]byte, 0)
		for i := 0; i < 10; i += 1 {
			d := p.Get()
			s = append(s, d)
		}
		// fill it
		for _, d := range s {
			p.Put(d)
		}

		if 10 != p.Len() {
			tt.Errorf("fill-ed pool: %d", p.Len())
		}
		if 10 != p.Cap() {
			tt.Errorf("max capacity = 10")
		}

		// discard it
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

func TestMmapBytePoolPreload(t *testing.T) {
	p := NewMmapBytePool(12, 8, Preload(true))
	l := int(float64(12) * defaultPreloadRate)
	if p.Len() != l {
		t.Errorf("preloaded buffer = %d", p.Len())
	}
}

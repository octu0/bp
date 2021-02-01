package bp

import (
	"bytes"
	"runtime"
	"strings"
	"sync"
	"testing"

	"github.com/octu0/chanque"
)

func BenchmarkBufferPool(b *testing.B) {
	k8 := []byte(strings.Repeat("@", 8))
	k4096 := []byte(strings.Repeat("@", 4096))
	run := func(name string, fn func(*testing.B)) {
		m1 := new(runtime.MemStats)
		runtime.ReadMemStats(m1)

		b.Run(name, fn)

		m2 := new(runtime.MemStats)
		runtime.ReadMemStats(m2)
		b.Logf(
			"%-20s\tTotalAlloc=%5d\tStackInUse=%5d",
			name,
			int64(m2.TotalAlloc)-int64(m1.TotalAlloc),
			int64(m2.StackInuse)-int64(m1.StackInuse),
			//int64(m2.HeapSys)  - int64(m1.HeapSys),
			//int64(m2.HeapIdle)   - int64(m1.HeapIdle),
		)
	}
	run("default/8", func(tb *testing.B) {
		e := chanque.NewExecutor(10, 10)
		defer e.Release()

		for i := 0; i < tb.N; i += 1 {
			e.Submit(func() {
				s := bytes.NewBuffer(make([]byte, 0, 8))
				s.Write(k8)
			})
		}
	})
	run("default/4096", func(tb *testing.B) {
		e := chanque.NewExecutor(10, 10)
		defer e.Release()

		for i := 0; i < tb.N; i += 1 {
			e.Submit(func() {
				s := bytes.NewBuffer(make([]byte, 0, 4096))
				s.Write(k4096)
			})
		}
	})
	run("syncpool/8", func(tb *testing.B) {
		e := chanque.NewExecutor(10, 10)
		defer e.Release()

		p := &sync.Pool{
			New: func() interface{} {
				return bytes.NewBuffer(make([]byte, 0, 8))
			},
		}
		for i := 0; i < tb.N; i += 1 {
			e.Submit(func() {
				s := p.Get().(*bytes.Buffer)
				s.Write(k8)
				s.Reset()
				p.Put(s)
			})
		}
	})
	run("syncpool/4096", func(tb *testing.B) {
		e := chanque.NewExecutor(10, 10)
		defer e.Release()

		p := &sync.Pool{
			New: func() interface{} {
				return bytes.NewBuffer(make([]byte, 0, 4096))
			},
		}
		for i := 0; i < tb.N; i += 1 {
			e.Submit(func() {
				s := p.Get().(*bytes.Buffer)
				s.Write(k4096)
				s.Reset()
				p.Put(s)
			})
		}
	})
	run("bufferpool/8", func(tb *testing.B) {
		e := chanque.NewExecutor(10, 10)
		defer e.Release()

		p := NewBufferPool(e.MaxWorker(), 8)

		for i := 0; i < tb.N; i += 1 {
			e.Submit(func() {
				s := p.Get()
				s.Write(k8)
				p.Put(s)
			})
		}
	})
	run("bufferpool/4096", func(tb *testing.B) {
		e := chanque.NewExecutor(10, 10)
		defer e.Release()

		p := NewBufferPool(e.MaxWorker(), 4096)

		for i := 0; i < tb.N; i += 1 {
			e.Submit(func() {
				s := p.Get()
				s.Write(k4096)
				p.Put(s)
			})
		}
	})
}

func TestBufferPoolBufSize(t *testing.T) {
	bufSize := 8
	t.Run("getsamecap", func(tt *testing.T) {
		p := NewBufferPool(10, bufSize)

		d1 := p.Get()
		d2 := p.Get()
		if d1.Cap() != bufSize {
			tt.Errorf("buf size = %d", bufSize)
		}
		if d2.Cap() != bufSize {
			tt.Errorf("buf size = %d", bufSize)
		}
	})
	t.Run("getput/samecap", func(tt *testing.T) {
		p := NewBufferPool(10, bufSize)
		d1 := p.Get()
		p.Put(d1)
		d2 := p.Get()
		p.Put(d2)

		d3 := p.Get()
		d4 := p.Get()
		if d3.Cap() != bufSize {
			tt.Errorf("buf size = %d", bufSize)
		}
		if d4.Cap() != bufSize {
			tt.Errorf("buf size = %d", bufSize)
		}
	})
	t.Run("getput/largecap", func(tt *testing.T) {
		p := NewBufferPool(10, bufSize)
		p.Put(bytes.NewBuffer(make([]byte, 0, 20)))

		d1 := p.Get()
		if d1.Cap() == bufSize {
			tt.Errorf("manually set buffer capacity be different: %d", d1.Cap())
		}
		if d1.Len() != 0 {
			tt.Errorf("manually set buffer be Reset")
		}
	})
}

func TestBufferPoolDiscard(t *testing.T) {
	bufSize := 8
	t.Run("freecap/samesize", func(tt *testing.T) {
		p := NewBufferPool(10, bufSize)

		d1 := bytes.NewBuffer(make([]byte, 0, bufSize))
		if p.Put(d1) != true {
			tt.Errorf("freecap %d", p.Cap())
		}
	})
	t.Run("fullcap/samesize", func(tt *testing.T) {
		p := NewBufferPool(2, bufSize)
		p.Put(bytes.NewBuffer(make([]byte, 0, bufSize)))
		p.Put(bytes.NewBuffer(make([]byte, 0, bufSize)))

		d1 := bytes.NewBuffer(make([]byte, 0, bufSize))
		if p.Put(d1) {
			tt.Errorf("fulled capacity %d", p.Cap())
		}
	})
	t.Run("freecap/largesize", func(tt *testing.T) {
		p := NewBufferPool(10, bufSize)
		p.Put(bytes.NewBuffer(make([]byte, 0, bufSize)))

		if p.Put(bytes.NewBuffer(make([]byte, 0, 100))) {
			tt.Errorf("maxBufSize put ng")
		}
		if p.Put(bytes.NewBuffer(make([]byte, 100))) {
			tt.Errorf("maxBufSize put ng")
		}

		if p.Put(bytes.NewBuffer(make([]byte, 0, 30))) != true {
			tt.Errorf("less than maxBufSize put no")
		}
		if p.Put(bytes.NewBuffer(make([]byte, 30))) != true {
			tt.Errorf("less than maxBufSize put no")
		}
	})
	t.Run("freecap/smallsize", func(tt *testing.T) {
		p := NewBufferPool(10, bufSize)
		p.Put(bytes.NewBuffer(make([]byte, 0, bufSize)))

		if p.Put(bytes.NewBuffer(nil)) != true {
			tt.Errorf("put ok // call bytes.Grow")
		}
		if p.Put(bytes.NewBuffer(make([]byte, 0, 1))) != true {
			tt.Errorf("put ok // call bytes.Grow")
		}
		if p.Put(bytes.NewBuffer(make([]byte, 1))) != true {
			tt.Errorf("put ok // call bytes.Grow")
		}
	})
}

func TestBufferPoolLenCap(t *testing.T) {
	t.Run("getput", func(tt *testing.T) {
		bufSize := 1
		p := NewBufferPool(10, bufSize)
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
		p := NewBufferPool(10, bufSize)
		s := make([]*bytes.Buffer, 0)

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

		d1 := bytes.NewBuffer(make([]byte, 0, bufSize))
		d2 := bytes.NewBuffer(make([]byte, 0, bufSize))
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

func TestBufferPoolPreload(t *testing.T) {
	p := NewBufferPool(12, 8, Preload(true))
	l := int(float64(12) * defaultPreloadRate)
	if p.Len() != l {
		t.Errorf("preloaded buffer = %d", p.Len())
	}
}

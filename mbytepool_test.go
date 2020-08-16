package bp

import (
	"testing"
)

func TestMultiBytePoolNew(t *testing.T) {
	t.Run("sorted", func(tt *testing.T) {
		mp := NewMultiBytePool(
			MultiBytePoolSize(10, 4),
			MultiBytePoolSize(10, 8),
			MultiBytePoolSize(10, 16),
		)
		if mp.pools[0].bufSize != 4 {
			tt.Errorf("sorted")
		}
		if mp.pools[2].bufSize != 16 {
			tt.Errorf("sorted")
		}
		for _, p := range mp.pools {
			if p.Len() != 0 {
				tt.Errorf("initial pool size 0")
			}
		}
	})
	t.Run("preload", func(tt *testing.T) {
		mp := NewMultiBytePool(
			MultiBytePoolSize(10, 4),
			MultiBytePoolSize(10, 8),
			MultiBytePoolSize(10, 16),
			MultiBytePoolOption(
				Preload(true),
				PreloadRate(0.5),
			),
		)
		if mp.pools[0].bufSize != 4 {
			tt.Errorf("sorted head")
		}
		if mp.pools[2].bufSize != 16 {
			tt.Errorf("sorted tail")
		}
		for _, p := range mp.pools {
			l := int(float64(p.Cap()) * 0.5)
			if p.Len() != l {
				tt.Errorf("preloaded %d", l)
			}
		}
	})
}

func TestMultiBytePoolPutGet(t *testing.T) {
	t.Run("getput", func(tt *testing.T) {
		mp := NewMultiBytePool(
			MultiBytePoolSize(10, 4),
			MultiBytePoolSize(10, 8),
			MultiBytePoolSize(10, 16),
		)
		d1 := mp.Get(1)
		d2 := mp.Get(4)
		if mp.Put(d1) != true {
			tt.Errorf("release ok / freecap")
		}
		if mp.Put(d2) != true {
			tt.Errorf("release ok / freecap")
		}
		if mp.pools[0].Len() != 2 {
			tt.Errorf("release pools[0] size = 1 and size = 4")
		}

		d3 := mp.Get(5)  // pools[0] < 5 < pools[1]
		d4 := mp.Get(10) // pools[1] < 10 < pools[2]
		if mp.Put(d3) != true {
			tt.Errorf("release ok")
		}
		if mp.Put(d4) != true {
			tt.Errorf("release ok")
		}
		if mp.pools[1].Len() != 1 {
			tt.Errorf("release pools[1] size = 5")
		}
		if mp.pools[2].Len() != 1 {
			tt.Errorf("release pools[1] size = 10")
		}

		d5 := mp.Get(1024)
		if mp.Put(d5) {
			tt.Errorf("discard large pool")
		}
	})
	t.Run("getref", func(tt *testing.T) {
		mp := NewMultiBytePool(
			MultiBytePoolSize(10, 4),
			MultiBytePoolSize(10, 8),
			MultiBytePoolSize(10, 16),
		)
		d1 := mp.GetRef(3)
		if cap(d1.Bytes()) != 4 {
			tt.Errorf("use pools[0]")
		}
		d1.Release()
		if mp.pools[0].Len() != 1 {
			tt.Errorf("released pools[0]")
		}

		d2 := mp.GetRef(10)
		if cap(d2.Bytes()) != 16 {
			tt.Errorf("use pools[2]")
		}
		d2.Release()
		if mp.pools[2].Len() != 1 {
			tt.Errorf("release pools[2]")
		}

		d3 := mp.GetRef(1024)
		d3.Release()
		if mp.pools[2].Len() != 1 {
			tt.Errorf("too large discard pools[2]")
		}

		d4 := mp.GetRef(20)
		d4.Release()
		if mp.pools[2].Len() != 2 {
			tt.Errorf("nearby size release buffer")
		}
	})
}

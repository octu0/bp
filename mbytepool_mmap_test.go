// +build aix darwin dragonfly freebsd linux netbsd openbsd solaris

package bp

import (
	"testing"
)

func TestMultiMmapBytePoolNew(t *testing.T) {
	t.Run("sorted", func(tt *testing.T) {
		mp := NewMultiMmapBytePool(
			MultiMmapBytePoolSize(10, 4),
			MultiMmapBytePoolSize(10, 8),
			MultiMmapBytePoolSize(10, 16),
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
		mp := NewMultiMmapBytePool(
			MultiMmapBytePoolSize(10, 4),
			MultiMmapBytePoolSize(10, 8),
			MultiMmapBytePoolSize(10, 16),
			MultiMmapBytePoolOption(
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

func TestMultiMmapBytePoolPutGet(t *testing.T) {
	t.Run("getput", func(tt *testing.T) {
		mp := NewMultiMmapBytePool(
			MultiMmapBytePoolSize(10, 7),
			MultiMmapBytePoolSize(10, 16),
			MultiMmapBytePoolSize(10, 24),
		)
		d1 := mp.Get(1) // 7 < pools[0]
		d2 := mp.Get(8) // 8 <= pools[0]
		if mp.Put(d1) != true {
			tt.Errorf("release ok / freecap")
		}
		if mp.Put(d2) != true {
			tt.Errorf("release ok / freecap")
		}
		if mp.pools[0].Len() != 2 {
			tt.Errorf("pool[0,1,2] = %d,%d,%d", mp.pools[0].Len(), mp.pools[1].Len(), mp.pools[2].Len())
		}

		d3 := mp.Get(9)  // pools[0] < 9 <= pools[1]
		d4 := mp.Get(16) //  pools[0] < 16 <= pools[2]
		if mp.Put(d3) != true {
			tt.Errorf("release ok")
		}
		if mp.Put(d4) != true {
			tt.Errorf("release ok")
		}
		if mp.pools[1].Len() != 2 {
			tt.Errorf("pool[0,1,2] = %d,%d,%d", mp.pools[0].Len(), mp.pools[1].Len(), mp.pools[2].Len())
		}
		if mp.pools[2].Len() != 0 {
			tt.Errorf("pool[0,1,2] = %d,%d,%d", mp.pools[0].Len(), mp.pools[1].Len(), mp.pools[2].Len())
		}

		d5 := mp.Get(1024)
		if mp.Put(d5) {
			tt.Errorf("discard large pool")
		}
	})
	t.Run("getref", func(tt *testing.T) {
		mp := NewMultiMmapBytePool(
			MultiMmapBytePoolSize(10, 7),
			MultiMmapBytePoolSize(10, 8),
			MultiMmapBytePoolSize(10, 16),
		)
		d1 := mp.GetRef(3)
		tt.Logf("pools[0] cap=%d", cap(d1.Bytes()))
		d1.Release()
		if mp.pools[0].Len() != 1 {
			tt.Errorf("released pools[0]")
		}

		d2 := mp.GetRef(10)
		tt.Logf("pools[1] cap=%d", cap(d2.Bytes()))
		d2.Release()
		if mp.pools[1].Len() != 1 {
			tt.Errorf("release pools[1]")
		}

		d3 := mp.GetRef(24)
		tt.Logf("pools[2] cap=%d", cap(d3.Bytes()))
		d3.Release()
		if mp.pools[2].Len() != 1 {
			tt.Errorf("release pools[2]")
		}

		d4 := mp.GetRef(1024)
		d4.Release()
		if mp.pools[2].Len() != 1 {
			tt.Errorf("discard, large alignment")
		}

		d5 := mp.GetRef(25)
		d5.Release()
		if mp.pools[2].Len() != 1 {
			tt.Errorf("discard, large alignment")
		}
	})
}

package bp

import (
	"fmt"
	"image"
	"testing"
)

func TestMultiImageRGBAPoolNew(t *testing.T) {
	t.Run("sorted", func(tt *testing.T) {
		mp := NewMultiImageRGBAPool(
			MultiImagePoolSize(10, image.Rect(0, 0, 640, 640)),
			MultiImagePoolSize(10, image.Rect(0, 0, 640, 360)),
			MultiImagePoolSize(10, image.Rect(0, 0, 1280, 720)),
			MultiImagePoolSize(10, image.Rect(0, 0, 360, 640)),
			MultiImagePoolSize(10, image.Rect(0, 0, 720, 1280)),
		)
		rects := make([]string, len(mp.pools))
		for i, p := range mp.pools {
			rects[i] = fmt.Sprintf("%dx%d", p.rect.Dx(), p.rect.Dy())
		}
		order := []string{
			"360x640",
			"640x360",
			"640x640",
			"720x1280",
			"1280x720",
		}
		for i, s := range order {
			if rects[i] != s {
				tt.Errorf("sorted expect:pools[%d]=%s", i, s)
			}
		}
	})
	t.Run("preload", func(tt *testing.T) {
		mp := NewMultiImageRGBAPool(
			MultiImagePoolSize(10, image.Rect(0, 0, 640, 360)),
			MultiImagePoolSize(10, image.Rect(0, 0, 1280, 720)),
			MultiImagePoolSize(10, image.Rect(0, 0, 360, 640)),
			MultiImagePoolOption(
				Preload(true),
				PreloadRate(0.5),
			),
		)
		if mp.pools[0].rect.Eq(image.Rect(0, 0, 360, 640)) != true {
			tt.Errorf("sorted head")
		}
		if mp.pools[2].rect.Eq(image.Rect(0, 0, 1280, 720)) != true {
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

func TestMultiImageRGBAPoolPutGet(t *testing.T) {
	t.Run("getput", func(tt *testing.T) {
		mp := NewMultiImageRGBAPool(
			MultiImagePoolSize(10, image.Rect(0, 0, 640, 360)),
			MultiImagePoolSize(10, image.Rect(0, 0, 1280, 720)),
			MultiImagePoolSize(10, image.Rect(0, 0, 360, 640)),
		)
		d1 := mp.GetRef(image.Rect(0, 0, 100, 100)) // 100x100 < pools[0]
		d2 := mp.GetRef(image.Rect(0, 0, 360, 120)) // 360x120 < pools[0]
		d3 := mp.GetRef(image.Rect(0, 0, 360, 700)) // pools[1] < 360x700 < pools[2]
		d4 := mp.GetRef(image.Rect(0, 0, 640, 320)) // 640x320 < pools[1]
		d5 := mp.GetRef(image.Rect(0, 0, 768, 432)) // pools[1] < 768x432 < pools[2]
		if mp.Put(d1.pix, d1.Img.Bounds()) != true {
			tt.Errorf("release ok / free cap")
		}
		if mp.pools[0].Len() != 1 {
			tt.Errorf("release pool[0] 100x100")
		}
		if mp.Put(d2.pix, d2.Img.Bounds()) != true {
			tt.Errorf("release ok / free cap")
		}
		if mp.pools[0].Len() != 2 {
			tt.Errorf("release pool[0] 360x120")
		}
		if mp.Put(d3.pix, d3.Img.Bounds()) != true {
			tt.Errorf("release ok / free cap")
		}
		if mp.pools[2].Len() != 1 {
			tt.Errorf("release pool[2] 360x700")
		}
		if mp.Put(d4.pix, d4.Img.Bounds()) != true {
			tt.Errorf("release ok / free cap")
		}
		if mp.pools[1].Len() != 1 {
			tt.Errorf("release pool[1] 640x320")
		}
		if mp.Put(d5.pix, d5.Img.Bounds()) != true {
			tt.Errorf("release ok / free cap")
		}
		if mp.pools[2].Len() != 2 {
			tt.Errorf("release pool[2] 768x432")
		}
	})
	t.Run("getref", func(tt *testing.T) {
		mp := NewMultiImageRGBAPool(
			MultiImagePoolSize(10, image.Rect(0, 0, 640, 360)),
			MultiImagePoolSize(10, image.Rect(0, 0, 1280, 720)),
			MultiImagePoolSize(10, image.Rect(0, 0, 360, 640)),
		)
		d1 := mp.GetRef(image.Rect(0, 0, 100, 100)) // 100x100 < pools[0]
		d2 := mp.GetRef(image.Rect(0, 0, 360, 120)) // 360x120 < pools[0]
		d3 := mp.GetRef(image.Rect(0, 0, 360, 700)) // pools[1] < 360x700 < pools[2]
		d4 := mp.GetRef(image.Rect(0, 0, 640, 320)) // 640x320 < pools[1]
		d5 := mp.GetRef(image.Rect(0, 0, 768, 432)) // pools[1] < 768x432 < pools[2]
		d1.Release()
		if mp.pools[0].Len() != 1 {
			tt.Errorf("release pool[0] 100x100")
		}
		d2.Release()
		if mp.pools[0].Len() != 2 {
			tt.Errorf("release pool[0] 360x120")
		}
		d3.Release()
		if mp.pools[2].Len() != 1 {
			tt.Errorf("release pool[2] 360x700")
		}
		d4.Release()
		if mp.pools[1].Len() != 1 {
			tt.Errorf("release pool[1] 640x320")
		}
		d5.Release()
		if mp.pools[2].Len() != 2 {
			tt.Errorf("release pool[2] 768x432")
		}

		d6 := mp.GetRef(image.Rect(0, 0, 1920, 1080))
		d6.Release()
		if mp.pools[2].Len() == 2 {
			tt.Errorf("put ok large pix")
		}
	})
}

func TestMultiImageNRGBAPoolNew(t *testing.T) {
	t.Run("sorted", func(tt *testing.T) {
		mp := NewMultiImageNRGBAPool(
			MultiImagePoolSize(10, image.Rect(0, 0, 640, 640)),
			MultiImagePoolSize(10, image.Rect(0, 0, 640, 360)),
			MultiImagePoolSize(10, image.Rect(0, 0, 1280, 720)),
			MultiImagePoolSize(10, image.Rect(0, 0, 360, 640)),
			MultiImagePoolSize(10, image.Rect(0, 0, 720, 1280)),
		)
		rects := make([]string, len(mp.pools))
		for i, p := range mp.pools {
			rects[i] = fmt.Sprintf("%dx%d", p.rect.Dx(), p.rect.Dy())
		}
		order := []string{
			"360x640",
			"640x360",
			"640x640",
			"720x1280",
			"1280x720",
		}
		for i, s := range order {
			if rects[i] != s {
				tt.Errorf("sorted expect:pools[%d]=%s", i, s)
			}
		}
	})
	t.Run("preload", func(tt *testing.T) {
		mp := NewMultiImageNRGBAPool(
			MultiImagePoolSize(10, image.Rect(0, 0, 640, 360)),
			MultiImagePoolSize(10, image.Rect(0, 0, 1280, 720)),
			MultiImagePoolSize(10, image.Rect(0, 0, 360, 640)),
			MultiImagePoolOption(
				Preload(true),
				PreloadRate(0.5),
			),
		)
		if mp.pools[0].rect.Eq(image.Rect(0, 0, 360, 640)) != true {
			tt.Errorf("sorted head")
		}
		if mp.pools[2].rect.Eq(image.Rect(0, 0, 1280, 720)) != true {
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

func TestMultiImageNRGBAPoolPutGet(t *testing.T) {
	t.Run("getput", func(tt *testing.T) {
		mp := NewMultiImageNRGBAPool(
			MultiImagePoolSize(10, image.Rect(0, 0, 640, 360)),
			MultiImagePoolSize(10, image.Rect(0, 0, 1280, 720)),
			MultiImagePoolSize(10, image.Rect(0, 0, 360, 640)),
		)
		d1 := mp.GetRef(image.Rect(0, 0, 100, 100)) // 100x100 < pools[0]
		d2 := mp.GetRef(image.Rect(0, 0, 360, 120)) // 360x120 < pools[0]
		d3 := mp.GetRef(image.Rect(0, 0, 360, 700)) // pools[1] < 360x700 < pools[2]
		d4 := mp.GetRef(image.Rect(0, 0, 640, 320)) // 640x320 < pools[1]
		d5 := mp.GetRef(image.Rect(0, 0, 768, 432)) // pools[1] < 768x432 < pools[2]
		if mp.Put(d1.pix, d1.Img.Bounds()) != true {
			tt.Errorf("release ok / free cap")
		}
		if mp.pools[0].Len() != 1 {
			tt.Errorf("release pool[0] 100x100")
		}
		if mp.Put(d2.pix, d2.Img.Bounds()) != true {
			tt.Errorf("release ok / free cap")
		}
		if mp.pools[0].Len() != 2 {
			tt.Errorf("release pool[0] 360x120")
		}
		if mp.Put(d3.pix, d3.Img.Bounds()) != true {
			tt.Errorf("release ok / free cap")
		}
		if mp.pools[2].Len() != 1 {
			tt.Errorf("release pool[2] 360x700")
		}
		if mp.Put(d4.pix, d4.Img.Bounds()) != true {
			tt.Errorf("release ok / free cap")
		}
		if mp.pools[1].Len() != 1 {
			tt.Errorf("release pool[1] 640x320")
		}
		if mp.Put(d5.pix, d5.Img.Bounds()) != true {
			tt.Errorf("release ok / free cap")
		}
		if mp.pools[2].Len() != 2 {
			tt.Errorf("release pool[2] 768x432")
		}
	})
	t.Run("getref", func(tt *testing.T) {
		mp := NewMultiImageNRGBAPool(
			MultiImagePoolSize(10, image.Rect(0, 0, 640, 360)),
			MultiImagePoolSize(10, image.Rect(0, 0, 1280, 720)),
			MultiImagePoolSize(10, image.Rect(0, 0, 360, 640)),
		)
		d1 := mp.GetRef(image.Rect(0, 0, 100, 100)) // 100x100 < pools[0]
		d2 := mp.GetRef(image.Rect(0, 0, 360, 120)) // 360x120 < pools[0]
		d3 := mp.GetRef(image.Rect(0, 0, 360, 700)) // pools[1] < 360x700 < pools[2]
		d4 := mp.GetRef(image.Rect(0, 0, 640, 320)) // 640x320 < pools[1]
		d5 := mp.GetRef(image.Rect(0, 0, 768, 432)) // pools[1] < 768x432 < pools[2]
		d1.Release()
		if mp.pools[0].Len() != 1 {
			tt.Errorf("release pool[0] 100x100")
		}
		d2.Release()
		if mp.pools[0].Len() != 2 {
			tt.Errorf("release pool[0] 360x120")
		}
		d3.Release()
		if mp.pools[2].Len() != 1 {
			tt.Errorf("release pool[2] 360x700")
		}
		d4.Release()
		if mp.pools[1].Len() != 1 {
			tt.Errorf("release pool[1] 640x320")
		}
		d5.Release()
		if mp.pools[2].Len() != 2 {
			tt.Errorf("release pool[2] 768x432")
		}

		d6 := mp.GetRef(image.Rect(0, 0, 1920, 1080))
		d6.Release()
		if mp.pools[2].Len() == 2 {
			tt.Errorf("put ok large pix")
		}
	})
}

func TestMultiImageYCbCrPoolNew(t *testing.T) {
	t.Run("sorted", func(tt *testing.T) {
		mp := NewMultiImageYCbCrPool(
			image.YCbCrSubsampleRatio420,
			MultiImagePoolSize(10, image.Rect(0, 0, 640, 640)),
			MultiImagePoolSize(10, image.Rect(0, 0, 640, 360)),
			MultiImagePoolSize(10, image.Rect(0, 0, 1280, 720)),
			MultiImagePoolSize(10, image.Rect(0, 0, 360, 640)),
			MultiImagePoolSize(10, image.Rect(0, 0, 720, 1280)),
		)
		rects := make([]string, len(mp.pools))
		for i, p := range mp.pools {
			rects[i] = fmt.Sprintf("%dx%d", p.rect.Dx(), p.rect.Dy())
		}
		order := []string{
			"360x640",
			"640x360",
			"640x640",
			"720x1280",
			"1280x720",
		}
		for i, s := range order {
			if rects[i] != s {
				tt.Errorf("sorted expect:pools[%d]=%s", i, s)
			}
		}
	})
	t.Run("preload", func(tt *testing.T) {
		mp := NewMultiImageYCbCrPool(
			image.YCbCrSubsampleRatio420,
			MultiImagePoolSize(10, image.Rect(0, 0, 640, 360)),
			MultiImagePoolSize(10, image.Rect(0, 0, 1280, 720)),
			MultiImagePoolSize(10, image.Rect(0, 0, 360, 640)),
			MultiImagePoolOption(
				Preload(true),
				PreloadRate(0.5),
			),
		)
		if mp.pools[0].rect.Eq(image.Rect(0, 0, 360, 640)) != true {
			tt.Errorf("sorted head")
		}
		if mp.pools[2].rect.Eq(image.Rect(0, 0, 1280, 720)) != true {
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

func TestMultiImageYCbCrPoolPutGet(t *testing.T) {
	t.Run("getput", func(tt *testing.T) {
		mp := NewMultiImageYCbCrPool(
			image.YCbCrSubsampleRatio420,
			MultiImagePoolSize(10, image.Rect(0, 0, 640, 360)),
			MultiImagePoolSize(10, image.Rect(0, 0, 1280, 720)),
			MultiImagePoolSize(10, image.Rect(0, 0, 360, 640)),
		)
		d1 := mp.GetRef(image.Rect(0, 0, 100, 100)) // 100x100 < pools[0]
		d2 := mp.GetRef(image.Rect(0, 0, 360, 120)) // 360x120 < pools[0]
		d3 := mp.GetRef(image.Rect(0, 0, 360, 700)) // pools[1] < 360x700 < pools[2]
		d4 := mp.GetRef(image.Rect(0, 0, 640, 320)) // 640x320 < pools[1]
		d5 := mp.GetRef(image.Rect(0, 0, 768, 432)) // pools[1] < 768x432 < pools[2]
		if mp.Put(d1.pix, d1.Img.Bounds()) != true {
			tt.Errorf("release ok / free cap")
		}
		if mp.pools[0].Len() != 1 {
			tt.Errorf("release pool[0] 100x100")
		}
		if mp.Put(d2.pix, d2.Img.Bounds()) != true {
			tt.Errorf("release ok / free cap")
		}
		if mp.pools[0].Len() != 2 {
			tt.Errorf("release pool[0] 360x120")
		}
		if mp.Put(d3.pix, d3.Img.Bounds()) != true {
			tt.Errorf("release ok / free cap")
		}
		if mp.pools[2].Len() != 1 {
			tt.Errorf("release pool[2] 360x700")
		}
		if mp.Put(d4.pix, d4.Img.Bounds()) != true {
			tt.Errorf("release ok / free cap")
		}
		if mp.pools[1].Len() != 1 {
			tt.Errorf("release pool[1] 640x320")
		}
		if mp.Put(d5.pix, d5.Img.Bounds()) != true {
			tt.Errorf("release ok / free cap")
		}
		if mp.pools[2].Len() != 2 {
			tt.Errorf("release pool[2] 768x432")
		}
	})
	t.Run("getref", func(tt *testing.T) {
		mp := NewMultiImageYCbCrPool(
			image.YCbCrSubsampleRatio420,
			MultiImagePoolSize(10, image.Rect(0, 0, 640, 360)),
			MultiImagePoolSize(10, image.Rect(0, 0, 1280, 720)),
			MultiImagePoolSize(10, image.Rect(0, 0, 360, 640)),
		)
		d1 := mp.GetRef(image.Rect(0, 0, 100, 100)) // 100x100 < pools[0]
		d2 := mp.GetRef(image.Rect(0, 0, 360, 120)) // 360x120 < pools[0]
		d3 := mp.GetRef(image.Rect(0, 0, 360, 700)) // pools[1] < 360x700 < pools[2]
		d4 := mp.GetRef(image.Rect(0, 0, 640, 320)) // 640x320 < pools[1]
		d5 := mp.GetRef(image.Rect(0, 0, 768, 432)) // pools[1] < 768x432 < pools[2]
		d1.Release()
		if mp.pools[0].Len() != 1 {
			tt.Errorf("release pool[0] 100x100")
		}
		d2.Release()
		if mp.pools[0].Len() != 2 {
			tt.Errorf("release pool[0] 360x120")
		}
		d3.Release()
		if mp.pools[2].Len() != 1 {
			tt.Errorf("release pool[2] 360x700")
		}
		d4.Release()
		if mp.pools[1].Len() != 1 {
			tt.Errorf("release pool[1] 640x320")
		}
		d5.Release()
		if mp.pools[2].Len() != 2 {
			tt.Errorf("release pool[2] 768x432")
		}

		d6 := mp.GetRef(image.Rect(0, 0, 1920, 1080))
		d6.Release()
		if mp.pools[2].Len() == 2 {
			tt.Errorf("put ok large pix")
		}
	})
}

package bp

import(
  "testing"
  "image"
)

func TestMultiImageRGBAPool(t *testing.T) {
  mp := NewMultiImageRGBAPool(
    MultiImagePoolSize(10, image.Rect(0, 0, 640, 360)),
    MultiImagePoolSize(10, image.Rect(0, 0, 1280, 720)),
    MultiImagePoolSize(10, image.Rect(0, 0, 360, 640)),
    MultiImagePoolSize(10, image.Rect(0, 0, 720, 1280)),
  )
  b := mp.GetRef(image.Rect(0, 0, 120, 30))
  println("------", cap(b.pix))
}

func TestMultiImageYCbCrPool(t *testing.T) {
  mp := NewMultiImageYCbCrPool(
    image.YCbCrSubsampleRatio420,
    MultiImagePoolSize(10, image.Rect(0, 0, 640, 360)),
    MultiImagePoolSize(10, image.Rect(0, 0, 1280, 720)),
    MultiImagePoolSize(10, image.Rect(0, 0, 360, 640)),
    MultiImagePoolSize(10, image.Rect(0, 0, 720, 1280)),
  )
  b := mp.GetRef(image.Rect(0, 0, 120, 30))
  println("------", cap(b.pix))
}

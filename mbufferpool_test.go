package bp

import(
  "testing"
)

func TestMultiBufferPoolNew(t *testing.T) {
  mp := NewMultiBufferPool(
    MultiBufferPoolSize(10, 4),
    MultiBufferPoolSize(10, 8),
    MultiBufferPoolSize(10, 16),
  )
  b := mp.Get(1)
  println("-----", b.Cap())
}

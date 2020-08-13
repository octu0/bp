package bp

import(
  "testing"
)

func TestNewMultiBytePoolNew(t *testing.T) {
  mp := NewMultiBytePool(
    MultiBytePoolSize(10, 4),
    MultiBytePoolSize(10, 8),
    MultiBytePoolSize(10, 16),
  )
  b := mp.Get(1)
  println("----", cap(b))
}

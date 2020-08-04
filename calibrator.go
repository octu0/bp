package bp

import(
  "sync/atomic"
)

const(
  defaultCallMod uint64 = 1000
)

type CalibrateHandler interface {
  CalibrateBytePool(*BytePool)
  CalibrateBufferPool(*BufferPool)
  CalibrateBufioReaderPool(*BufioReaderPool)
  CalibrateBufioWriterPool(*BufioWriterPool)
  CalibrateImageRGBAPool(*ImageRGBAPool)
  CalibrateImageYCbCrPool(*ImageYCbCrPool)
}

var(
  _ CalibrateHandler = (*capacityUtilCalibrator)(nil)
)

type capacityUtilCalibrator struct {
  utilRate float64
  counter  uint64
}

func (c *capacityUtilCalibrator) increAndRun() bool {
  newValue := atomic.AddUint64(&c.counter, 1)
  if 0 == ((newValue - 1) % defaultCallMod) {
    return true
  }
  return false
}

func (c *capacityUtilCalibrator) CalibrateBytePool(p *BytePool) {
}

func (c *capacityUtilCalibrator) CalibrateBufferPool(p *BufferPool) {
}

func (c *capacityUtilCalibrator) CalibrateBufioReaderPool(p *BufioReaderPool) {
}

func (c *capacityUtilCalibrator) CalibrateBufioWriterPool(p *BufioWriterPool) {
}

func (c *capacityUtilCalibrator) CalibrateImageRGBAPool(p *ImageRGBAPool) {
}

func (c *capacityUtilCalibrator) CalibrateImageYCbCrPool(p *ImageYCbCrPool) {
}

func CapacityFillRate(rate float64) CalibrateHandler {
  return &capacityUtilCalibrator{
    utilRate: rate,
    counter:  uint64(0),
  }
}

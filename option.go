package bp

type optionFunc func(*option)

type option struct {
  calibrator       CalibrateHandler
  preload          bool
  maxBufSizeFactor float64
}

func newOption() *option {
  return &option{
    calibrator:       nil,
    preload:          false,
    maxBufSizeFactor: 4.0,
  }
}

func Calibrator(c CalibrateHandler) optionFunc {
  return func(opt *option) {
    opt.calibrator = c
  }
}

func Preload(enable bool) optionFunc {
  return func(opt *option) {
    opt.preload = enable
  }
}

func MaxBufSizeFactor(factor float64) optionFunc {
  return func(opt *option) {
    opt.maxBufSizeFactor = factor
  }
}

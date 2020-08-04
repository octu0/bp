package bp

type optionFunc func(*option)

type option struct {
  calibrator CalibrateHandler
  preload    bool
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

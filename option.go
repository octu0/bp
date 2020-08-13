package bp

type optionFunc func(*option)

const(
  defaultPreloadEnable    bool    = false
  defaultMaxBufSizeFactor float64 = 4.0
)

type option struct {
  preload          bool
  maxBufSizeFactor float64
}

func newOption() *option {
  return &option{
    preload:          defaultPreloadEnable,
    maxBufSizeFactor: defaultMaxBufSizeFactor,
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

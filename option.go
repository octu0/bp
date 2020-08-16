package bp

type optionFunc func(*option)

const (
	defaultPreloadEnable    bool    = false
	defaultPreloadRate      float64 = 0.25
	defaultMaxBufSizeFactor float64 = 4.0
)

type option struct {
	preload          bool
	preloadRate      float64
	maxBufSizeFactor float64
}

func newOption() *option {
	return &option{
		preload:          defaultPreloadEnable,
		preloadRate:      defaultPreloadRate,
		maxBufSizeFactor: defaultMaxBufSizeFactor,
	}
}

func Preload(enable bool) optionFunc {
	return func(opt *option) {
		opt.preload = enable
	}
}

func PreloadRate(rate float64) optionFunc {
	return func(opt *option) {
		opt.preloadRate = rate
	}
}

func MaxBufSizeFactor(factor float64) optionFunc {
	return func(opt *option) {
		opt.maxBufSizeFactor = factor
	}
}

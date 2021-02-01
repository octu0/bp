package bp

type optionFunc func(*option)

const (
	defaultPreloadEnable    bool    = false
	defaultPreloadRate      float64 = 0.25
	defaultMaxBufSizeFactor float64 = 1.25
	defaultAutoGrowEnable   bool    = false
)

type option struct {
	preload          bool
	preloadRate      float64
	maxBufSizeFactor float64
	autoGrow         bool
}

func newOption() *option {
	return &option{
		preload:          defaultPreloadEnable,
		preloadRate:      defaultPreloadRate,
		maxBufSizeFactor: defaultMaxBufSizeFactor,
		autoGrow:         defaultAutoGrowEnable,
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

func AutoGrow(enable bool) optionFunc {
	return func(opt *option) {
		opt.autoGrow = enable
	}
}

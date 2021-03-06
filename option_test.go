package bp

import (
	"testing"
)

func TestOption(t *testing.T) {
	opt := newOption()
	if opt.preload != defaultPreloadEnable {
		t.Errorf("default preload option = %v", defaultPreloadEnable)
	}
	if opt.maxBufSizeFactor != defaultMaxBufSizeFactor {
		t.Errorf("default max bufsize factor = %v", defaultMaxBufSizeFactor)
	}

	options := []optionFunc{
		Preload(true),
		PreloadRate(0.85),
		MaxBufSizeFactor(12.3),
	}
	for _, fn := range options {
		fn(opt)
	}

	if opt.preload != true {
		t.Errorf("option set preload = true")
	}
	if opt.preloadRate != 0.85 {
		t.Errorf("option set preload rate = 0.85")
	}
	if opt.maxBufSizeFactor != 12.3 {
		t.Errorf("option set max bufsize factor = 12.3")
	}
}

package bp

import (
	"time"
)

type TickerPool struct {
	pool chan *time.Ticker
}

func (b *TickerPool) GetRef(dur time.Duration) *TickerRef {
	t := b.Get(dur)

	ref := newTickerRef(t, b)
	ref.setFinalizer()
	return ref
}

func (b *TickerPool) Get(dur time.Duration) *time.Ticker {
	select {
	case t := <-b.pool:
		// reuse exists pool
		t.Reset(dur)
		return t
	default:
		// create *time.Ticker
		return time.NewTicker(dur)
	}
}

func (b *TickerPool) Put(t *time.Ticker) bool {
	t.Stop()

	select {
	case b.pool <- t:
		return true
	default:
		return false
	}
}

func (b *TickerPool) Len() int {
	return len(b.pool)
}

func (b *TickerPool) Cap() int {
	return cap(b.pool)
}

func NewTickerPool(poolSize int, funcs ...optionFunc) *TickerPool {
	opt := newOption()
	for _, fn := range funcs {
		fn(opt)
	}

	return &TickerPool{
		pool: make(chan *time.Ticker, poolSize),
	}
}

type TimerPool struct {
	pool chan *time.Timer
}

func (b *TimerPool) GetRef(dur time.Duration) *TimerRef {
	t := b.Get(dur)

	ref := newTimerRef(t, b)
	ref.setFinalizer()
	return ref
}

func (b *TimerPool) Get(dur time.Duration) *time.Timer {
	select {
	case t := <-b.pool:
		// reuse exists pool
		t.Reset(dur)
		return t
	default:
		// create *time.Ticker
		return time.NewTimer(dur)
	}
}

func (b *TimerPool) Put(t *time.Timer) bool {
	if t.Stop() != true {
		select {
		case <-t.C:
			// drain
		default:
			// free
		}
	}

	select {
	case b.pool <- t:
		return true
	default:
		return false
	}
}

func (b *TimerPool) Len() int {
	return len(b.pool)
}

func (b *TimerPool) Cap() int {
	return cap(b.pool)
}

func NewTimerPool(poolSize int, funcs ...optionFunc) *TimerPool {
	opt := newOption()
	for _, fn := range funcs {
		fn(opt)
	}

	return &TimerPool{
		pool: make(chan *time.Timer, poolSize),
	}
}

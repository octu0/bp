package bp

import (
	"sync"
	"testing"
	"time"

	"github.com/octu0/chanque"
)

func BenchmarkTickerPool(b *testing.B) {
	dur := 1 * time.Millisecond
	b.Run("default", func(tb *testing.B) {
		e := chanque.NewExecutor(100, 100)
		defer e.Release()

		wg := new(sync.WaitGroup)
		wg.Add(tb.N)
		for i := 0; i < tb.N; i += 1 {
			e.Submit(func(_wg *sync.WaitGroup) chanque.Job {
				return func() {
					defer _wg.Done()

					t := time.NewTicker(dur)
					defer t.Stop()
					select {
					case <-t.C:
					}
				}
			}(wg))
		}
		wg.Wait()
	})
	b.Run("syncpool", func(tb *testing.B) {
		e := chanque.NewExecutor(100, 100)
		defer e.Release()

		pool := new(sync.Pool)

		wg := new(sync.WaitGroup)
		wg.Add(tb.N)
		for i := 0; i < tb.N; i += 1 {
			e.Submit(func(_wg *sync.WaitGroup) chanque.Job {
				return func() {
					defer _wg.Done()

					v := pool.Get()
					var t *time.Ticker
					if v != nil {
						t = v.(*time.Ticker)
						t.Reset(dur)
					} else {
						t = time.NewTicker(dur)
					}
					defer func() {
						t.Stop()
						pool.Put(t)
					}()

					select {
					case <-t.C:
					}
				}
			}(wg))
		}
		wg.Wait()
	})
	b.Run("pool", func(tb *testing.B) {
		e := chanque.NewExecutor(100, 100)
		defer e.Release()

		pool := NewTickerPool(e.MaxWorker())

		wg := new(sync.WaitGroup)
		wg.Add(tb.N)
		for i := 0; i < tb.N; i += 1 {
			e.Submit(func(_wg *sync.WaitGroup) chanque.Job {
				return func() {
					defer _wg.Done()

					t := pool.Get(dur)
					defer pool.Put(t)

					select {
					case <-t.C:
					}
				}
			}(wg))
		}
		wg.Wait()
	})
}

func BenchmarkTimerPool(b *testing.B) {
	dur := 1 * time.Millisecond
	b.Run("default/timer", func(tb *testing.B) {
		e := chanque.NewExecutor(100, 100)
		defer e.Release()

		wg := new(sync.WaitGroup)
		wg.Add(tb.N)
		for i := 0; i < tb.N; i += 1 {
			e.Submit(func(_wg *sync.WaitGroup) chanque.Job {
				return func() {
					defer _wg.Done()

					t := time.NewTimer(dur)
					select {
					case <-t.C:
					}
				}
			}(wg))
		}
		wg.Wait()
	})
	b.Run("default/timeafter", func(tb *testing.B) {
		e := chanque.NewExecutor(100, 100)
		defer e.Release()

		wg := new(sync.WaitGroup)
		wg.Add(tb.N)
		for i := 0; i < tb.N; i += 1 {
			e.Submit(func(_wg *sync.WaitGroup) chanque.Job {
				return func() {
					defer _wg.Done()

					select {
					case <-time.After(dur):
					}
				}
			}(wg))
		}
		wg.Wait()
	})
	b.Run("syncpool", func(tb *testing.B) {
		e := chanque.NewExecutor(100, 100)
		defer e.Release()

		pool := new(sync.Pool)

		wg := new(sync.WaitGroup)
		wg.Add(tb.N)
		for i := 0; i < tb.N; i += 1 {
			e.Submit(func(_wg *sync.WaitGroup) chanque.Job {
				return func() {
					defer _wg.Done()

					v := pool.Get()
					var t *time.Timer
					if v != nil {
						t = v.(*time.Timer)
						t.Reset(dur)
					} else {
						t = time.NewTimer(dur)
					}
					defer func() {
						if t.Stop() != true {
							select {
							case <-t.C:
							default:
							}
						}
						pool.Put(t)
					}()

					select {
					case <-t.C:
					}
				}
			}(wg))
		}
		wg.Wait()
	})
	b.Run("pool", func(tb *testing.B) {
		e := chanque.NewExecutor(100, 100)
		defer e.Release()

		pool := NewTimerPool(e.MaxWorker())

		wg := new(sync.WaitGroup)
		wg.Add(tb.N)
		for i := 0; i < tb.N; i += 1 {
			e.Submit(func(_wg *sync.WaitGroup) chanque.Job {
				return func() {
					defer _wg.Done()

					t := pool.Get(dur)
					defer pool.Put(t)

					select {
					case <-t.C:
					}
				}
			}(wg))
		}
		wg.Wait()
	})
}

func TestTickerPool(t *testing.T) {
	t.Run("dupStop", func(tt *testing.T) {
		pool := NewTickerPool(10)

		t1 := pool.Get(1 * time.Second)
		t1.Stop()
		pool.Put(t1)
	})

	t.Run("reuse", func(tt *testing.T) {
		pool := NewTickerPool(10)

		t1 := pool.Get(1 * time.Second)
		t2 := pool.Get(2 * time.Second)
		t3 := pool.Get(3 * time.Second)
		pool.Put(t1)
		pool.Put(t2)
		pool.Put(t3)

		for i := 0; i < 10; i += 1 {
			dur := time.Duration(i+10) * time.Millisecond

			tick := pool.Get(dur)
			start := time.Now()
			select {
			case <-tick.C:
				if time.Since(start) < dur {
					tt.Errorf("expired ticker")
				}
			default:
				// maybe new ticker
			}

			start = time.Now()
			counter := 0
			for j := 0; j < 10; j += 1 {
				select {
				case <-tick.C:
					counter += 1
				case <-time.After(200 * time.Millisecond):
					tt.Errorf("failed to ticker reset")
				}
			}
			if counter != 10 {
				tt.Errorf("tick is not emitted correct num of times. %d", counter)
			}
			if time.Since(start) < (dur * 10) {
				tt.Errorf("incorrect elaspse time of tick. %s", time.Since(start))
			}
			pool.Put(tick)
		}
	})
}

func TestTimerPool(t *testing.T) {
	t.Run("dupStop", func(tt *testing.T) {
		pool := NewTimerPool(10)

		t1 := pool.Get(1 * time.Second)
		t1.Stop()
		pool.Put(t1)
	})

	t.Run("reuse", func(tt *testing.T) {
		pool := NewTimerPool(10)

		t1 := pool.Get(1 * time.Second)
		t2 := pool.Get(2 * time.Second)
		t3 := pool.Get(3 * time.Second)
		pool.Put(t1)
		pool.Put(t2)
		pool.Put(t3)

		for i := 0; i < 10; i += 1 {
			dur := time.Duration(i+10) * time.Millisecond

			tick := pool.Get(dur)
			start := time.Now()
			select {
			case <-tick.C:
				if time.Since(start) < dur {
					tt.Errorf("expired ticker")
				}
			default:
				// maybe new ticker
			}

			start = time.Now()
			select {
			case <-tick.C:
			case <-time.After(200 * time.Millisecond):
				tt.Errorf("failed to ticker reset")
			}
			if time.Since(start) < dur {
				tt.Errorf("incorrect elaspse time of timer. %s", time.Since(start))
			}
			pool.Put(tick)
		}
	})
}

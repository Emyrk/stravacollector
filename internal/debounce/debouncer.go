package debounce

import (
	"sync"
	"time"
)

type Debouncer struct {
	mu    sync.Mutex
	delay time.Duration
	last  time.Time
}

func New(delay time.Duration) *Debouncer {
	return &Debouncer{
		delay: delay,
	}
}

func (d *Debouncer) Do(f func()) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if time.Since(d.last) < d.delay {
		return
	}

	f()
	d.last = time.Now()
}

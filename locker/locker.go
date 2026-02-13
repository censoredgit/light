package locker

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"time"
)

const acquireAttemptsTimeout = time.Second * 3
const retryLockTimeout = time.Millisecond * 250

var ErrLockTimeOut = errors.New("lock timeout")

type entry struct {
	m sync.RWMutex
	c atomic.Int32
}

type Locker struct {
	lock    sync.RWMutex
	storage map[string]*entry
	cfg     *Config
}

func (e *entry) RUnlock() {
	e.m.RUnlock()
	e.c.Add(1)
}

func (e *entry) Unlock() {
	e.m.Unlock()
	e.c.Add(-1)
}

func (l *Locker) len() int {
	l.lock.Lock()
	defer l.lock.Unlock()
	return len(l.storage)
}

type WriteUnlocker interface {
	Unlock()
}

type ReadUnlocker interface {
	RUnlock()
}

func New(cfg *Config) *Locker {
	l := &Locker{
		lock:    sync.RWMutex{},
		storage: make(map[string]*entry),
		cfg:     cfg,
	}

	go l.runGC()

	return l
}

func (l *Locker) WriteLock(id string) WriteUnlocker {
	e := l.getOrCreate(id)
	e.m.Lock()
	e.c.Add(1)

	return e
}

func (l *Locker) ReadLock(id string) ReadUnlocker {
	e := l.getOrCreate(id)
	e.m.RLock()
	e.c.Add(1)

	return e
}

func (l *Locker) SimpleLock(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), acquireAttemptsTimeout)
	defer cancel()

	return l.SimpleLockWithContext(ctx, id)
}

func (l *Locker) SimpleLockWithContext(ctx context.Context, id string) error {
	e := l.getOrCreate(id)

	for {
		select {
		case <-ctx.Done():
			return ErrLockTimeOut
		default:
			if e.m.TryLock() {
				return nil
			}
			time.Sleep(retryLockTimeout)
		}
	}
}

func (l *Locker) ReleaseSimpleLock(id string) {
	e := l.getOrCreate(id)
	if e.m.TryLock() {
		e.c.Add(1)
	}
	e.Unlock()
}

func (l *Locker) getOrCreate(id string) *entry {
	var e *entry
	var ok bool

	l.lock.RLock()
	if e, ok = l.storage[id]; ok {
		l.lock.RUnlock()
		return e
	}
	l.lock.RUnlock()
	l.lock.Lock()
	defer l.lock.Unlock()

	if e, ok = l.storage[id]; !ok {
		e = &entry{
			m: sync.RWMutex{},
			c: atomic.Int32{},
		}
		l.storage[id] = e
	}

	return e
}

func (l *Locker) runGC() {
	if l.cfg.GCTimeout <= 0 {
		return
	}

	for {
		time.Sleep(l.cfg.GCTimeout)

		l.lock.RLock()
		garbage := make([]string, 0, max(1, len(l.storage)/2))
		for id, e := range l.storage {
			if e.c.Load() == 0 {
				garbage = append(garbage, id)
			}
		}
		l.lock.RUnlock()

		if len(garbage) > 0 {
			l.lock.Lock()
			for id, e := range l.storage {
				if e.c.Load() == 0 {
					delete(l.storage, id)
				}
			}
			l.lock.Unlock()
		}
	}
}

package locker

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"
)

func TestConcurrentReadLock(t *testing.T) {
	locker := New(&Config{GCTimeout: 0})

	wg := sync.WaitGroup{}
	wg.Add(1_000_000)

	for i := 0; i < 1_000_000; i++ {
		go func() {
			c := locker.ReadLock("test")
			c.RUnlock()
			defer wg.Done()
		}()
	}

	wg.Wait()
}

func TestLockerGarbage(t *testing.T) {
	duration := 5 * time.Second

	locker := New(&Config{GCTimeout: duration})
	for range 2 {
		un := locker.WriteLock("test1")
		un.Unlock()

		un = locker.WriteLock("test2")
		un.Unlock()

		if locker.len() == 0 {
			t.Error("locker should not be empty")
		}

		time.Sleep(duration + 1*time.Second)
		if ll := locker.len(); ll > 0 {
			t.Errorf("locker should be empty, current length is %d", ll)
		}
	}
}

func TestManyLockerInstances(t *testing.T) {
	duration := 5 * time.Second

	wg := sync.WaitGroup{}
	for range 100_000 {
		wg.Add(1)

		go func() {
			defer wg.Done()

			locker := New(&Config{GCTimeout: duration})
			un := locker.WriteLock("test1")
			un.Unlock()

			un = locker.WriteLock("test2")
			un.Unlock()

			if locker.len() == 0 {
				t.Error("locker should not be empty")
			}
			time.Sleep(duration + 1*time.Second)
			if ll := locker.len(); ll > 0 {
				t.Errorf("locker should be empty, current length is %d", ll)
			}
		}()
	}

	wg.Wait()
}

func TestSimpleLock(t *testing.T) {
	locker := New(&Config{})

	err := locker.SimpleLock("test")
	if err != nil {
		t.Error(err)
	}

	locker.ReleaseSimpleLock("test")
}

func TestSimpleLockWaiting(t *testing.T) {
	locker := New(&Config{})
	err := locker.SimpleLock("test")
	if err != nil {
		t.Error(err)
	}

	go func() {
		time.Sleep(time.Second * 10)
		locker.ReleaseSimpleLock("test")
	}()

	for i := 0; i < 5; i++ {
		err = locker.SimpleLock("test")
		if !errors.Is(err, ErrLockTimeOut) {
			t.Error(err)
		} else {
			break
		}
	}

	locker.ReleaseSimpleLock("test")
}

func TestSimpleLockCtxWaiting(t *testing.T) {
	locker := New(&Config{})

	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(time.Second*1))
	defer cancel()

	err := locker.SimpleLock("test")
	if err != nil {
		t.Error(err)
	}

	go func() {
		time.Sleep(time.Second * 10)
		locker.ReleaseSimpleLock("test")
	}()

	for i := 0; i < 10; i++ {
		err = locker.SimpleLockWithContext(ctx, "test")
		if !errors.Is(err, ErrLockTimeOut) {
			t.Error(err)
		} else {
			break
		}
	}

	locker.ReleaseSimpleLock("test")
}

func TestSimpleLockCtx(t *testing.T) {
	locker := New(&Config{})

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	err := locker.SimpleLockWithContext(ctx, "test")
	if err != nil {
		t.Error(err)
	}

	locker.ReleaseSimpleLock("test")
}

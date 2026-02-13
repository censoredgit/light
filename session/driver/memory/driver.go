package memory

import (
	"fmt"
	"github.com/censoredgit/light/locker"
	"github.com/censoredgit/light/session"
	"time"
)

const DefaultGarbageSchedulerTime = 60 * time.Minute
const DefaultGarbageListInitCap = 1000

type driver struct {
	locker               *locker.Locker
	lifeTime             time.Duration
	garbageSchedulerTime time.Duration
	garbageListInitCap   uint
	garbage              []string
	storage              map[string]*session.Data
}

func Setup(
	locker *locker.Locker,
	lifeTime time.Duration,
	garbageSchedulerTime time.Duration,
	garbageListInitCap uint,
) session.Driver {
	return &driver{
		locker:               locker,
		lifeTime:             lifeTime,
		garbageSchedulerTime: garbageSchedulerTime,
		garbageListInitCap:   garbageListInitCap,
		garbage:              make([]string, garbageListInitCap),
		storage:              make(map[string]*session.Data),
	}
}

func (d *driver) Init() error {
	go d.runGarbageScheduler()
	return nil
}

func (d *driver) Open(id string) (*session.Data, error) {
	err := d.locker.SimpleLock(id)
	if err != nil {
		return nil, fmt.Errorf("session open error: %w", err)
	}

	if data, exists := d.storage[id]; exists {
		return data, nil
	}

	data := session.NewData(id, d.lifeTime)
	d.storage[id] = data
	return data, nil
}

func (d *driver) Close(data *session.Data) error {
	if data.IsNew() || data.IsModified() {
		d.storage[data.Id()] = data
	}

	d.locker.ReleaseSimpleLock(data.Id())

	return nil
}

func (d *driver) runGarbageScheduler() {
	for {
		time.Sleep(d.garbageSchedulerTime)

		d.collectGarbage()
		d.clearGarbage()
	}
}

func (d *driver) collectGarbage() {
	currentTime := time.Now()
	var unlocker locker.ReadUnlocker

	for id, data := range d.storage {
		unlocker = d.locker.ReadLock(id)
		if data.Expire().Before(currentTime) {
			d.garbage = append(d.garbage, id)
		}
		unlocker.RUnlock()
	}
}

func (d *driver) clearGarbage() {
	if len(d.garbage) == 0 {
		return
	}

	var unlocker locker.WriteUnlocker
	for _, id := range d.garbage {
		unlocker = d.locker.WriteLock(id)
		delete(d.storage, id)
		unlocker.Unlock()
	}

	d.garbage = make([]string, 0, d.garbageListInitCap)
}

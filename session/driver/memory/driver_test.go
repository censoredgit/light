package memory

import (
	"github.com/censoredgit/light/locker"
	"testing"
	"time"
)

func TestDriver(t *testing.T) {
	const sessId = "111"

	d := Setup(
		locker.New(&locker.Config{
			GCTimeout: time.Second * 5,
		},
		),
		time.Minute,
		DefaultGarbageSchedulerTime,
		DefaultGarbageListInitCap,
	)

	err := d.Init()
	if err != nil {
		t.Error(err)
	}

	data, err := d.Open(sessId)
	if err != nil {
		t.Error(err)
	}

	defer func() {
		err := d.Close(data)
		if err != nil {
			t.Error(err)
		}
	}()

	data.Set("test", "test data")

	if data.Get("test") != "test data" {
		t.Error("sess file corrupted")
	}

}

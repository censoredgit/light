package file

import (
	"github.com/censoredgit/light/locker"
	"log/slog"
	"os"
	"testing"
	"time"
)

func TestDriver(t *testing.T) {
	const sessId = "111"

	d := Setup(os.TempDir(),
		locker.New(&locker.Config{
			GCTimeout: time.Second * 5,
		}),
		slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelWarn,
		})),
		time.Minute*10,
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

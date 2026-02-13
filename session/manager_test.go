package session

import (
	"github.com/censoredgit/light/session/hasher"
	"log/slog"
	"net/http"
	"os"
	"testing"
	"time"
)

func makeManager() *Manager {
	return MustSetup(&Config{
		Hasher:     hasher.Md5Hasher{},
		Salt:       "test",
		TTL:        time.Minute,
		CookieName: "test",
		Driver:     &dummyDriver{},
		Logger: slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelWarn,
		})),
	})
}

func TestManager(t *testing.T) {
	m := makeManager()

	existsSessId := make(map[string]struct{})

	var d, nData *Data
	var cookie *http.Cookie
	var err error
	for range 100_000 {
		d, err = m.Init("")
		if err != nil {
			t.Fatal(err)
		}

		if _, has := existsSessId[d.Id()]; has {
			t.Fatal(d.Id(), " already exists")
			return
		}

		existsSessId[d.Id()] = struct{}{}
	}

	for v := range existsSessId {
		nData, err = m.Init(v)
		if err != nil {
			t.Fatal(err)
		}

		nData.Set("test", "test1")

		if !nData.IsModified() {
			t.Fatal("data should be modified")
		}

		if nData.IsModified() {
			t.Fatal("data should be not modified")
		}

		cookie = m.ToCookie(nData)
		if cookie.Value != nData.Id() {
			t.Fatal("cookie value not eq data.id")
		}

		if cookie.Expires.Compare(nData.expireTime) != 0 {
			t.Fatal("cookie expires not eq data.expireTime")
		}

		err = m.Close(nData)
		if err != nil {
			t.Fatal(err)
		}
	}
}

type dummyDriver struct {
	storage map[string]*Data
}

func (d *dummyDriver) Init() error {
	d.storage = make(map[string]*Data)
	return nil
}

func (d *dummyDriver) Open(id string) (*Data, error) {
	d.storage[id] = NewData(id, time.Minute)
	return d.storage[id], nil
}

func (d *dummyDriver) Close(data *Data) error {
	if data.IsNew() || data.IsModified() {
		d.storage[data.Id()] = data
	}

	delete(d.storage, data.Id())

	return nil
}

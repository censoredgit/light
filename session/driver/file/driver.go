package file

import (
	"encoding/json"
	"fmt"
	"github.com/censoredgit/light/locker"
	"github.com/censoredgit/light/session"
	"io"
	"log/slog"
	"os"
	"path"
	"time"
)

const DefaultGarbageSchedulerTime = 60 * time.Minute
const DefaultGarbageListInitCap = 1000

type driver struct {
	path                 string
	locker               *locker.Locker
	log                  *slog.Logger
	garbage              []string
	lifeTime             time.Duration
	garbageSchedulerTime time.Duration
	garbageListInitCap   uint
}

func Setup(
	path string,
	locker *locker.Locker,
	log *slog.Logger,
	lifeTime time.Duration,
	garbageSchedulerTime time.Duration,
	garbageListInitCap uint,
) session.Driver {
	return &driver{
		path:                 path,
		locker:               locker,
		log:                  log.With(slog.String("session.driver", "file")),
		lifeTime:             lifeTime,
		garbage:              make([]string, 0, garbageListInitCap),
		garbageSchedulerTime: garbageSchedulerTime,
		garbageListInitCap:   garbageListInitCap,
	}
}

func (d *driver) Init() error {
	f, err := os.Stat(d.path)
	if err != nil {
		return err
	}
	if !f.IsDir() {
		return fmt.Errorf("%w: %s", session.ErrNotExists, d.path)
	}

	go d.runGarbageScheduler()

	return nil
}

func (d *driver) Open(id string) (*session.Data, error) {
	err := d.locker.SimpleLock(id)
	if err != nil {
		return nil, fmt.Errorf("session open error: %w", err)
	}

	f, err := os.OpenFile(d.composeSessionFilePath(id), os.O_CREATE|os.O_RDWR, os.ModePerm)
	if err != nil {
		d.locker.ReleaseSimpleLock(id)
		return nil, fmt.Errorf("session open error: %w", err)
	}
	defer f.Close()

	fileData, err := io.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("session open error: %w", err)
	}

	data, err := decode(id, string(fileData))
	if err != nil {
		return session.NewData(id, d.lifeTime), nil
	}

	return data, nil
}

func (d *driver) Close(data *session.Data) error {
	var err error
	if data.IsNew() || data.IsModified() {
		err = writeSessionFile(d.composeSessionFilePath(data.Id()), data)
	}

	d.locker.ReleaseSimpleLock(data.Id())

	return err
}

func (d *driver) composeSessionFilePath(id string) string {
	return path.Clean(path.Join(d.path, id) + ".json")
}

func (d *driver) runGarbageScheduler() {
	var err error
	for {
		time.Sleep(d.garbageSchedulerTime)

		if err = d.collectGarbage(); err != nil {
			d.log.Error(err.Error())
			continue
		}
		d.clearGarbage()
	}
}

func (d *driver) collectGarbage() error {
	dirList, err := os.ReadDir(d.path)
	if err != nil {
		return err
	}

	currentTime := time.Now()

	for _, f := range dirList {
		if f.IsDir() {
			continue
		}

		fi, err := f.Info()
		if err != nil {
			continue
		}

		if fi.ModTime().Add(d.lifeTime).Before(currentTime) {
			d.garbage = append(d.garbage, path.Join(d.path, f.Name()))
		}
	}

	return nil
}

func (d *driver) clearGarbage() {
	if len(d.garbage) == 0 {
		return
	}

	var err error
	for _, f := range d.garbage {
		err = os.Remove(f)
		if err != nil {
			d.log.Warn(err.Error())
		}
	}

	d.garbage = make([]string, 0, d.garbageListInitCap)
}

func writeSessionFile(filePath string, data *session.Data) error {
	strJson, err := encode(data)
	if err != nil {
		return err
	}

	return os.WriteFile(filePath, []byte(strJson), os.ModePerm)
}

func encode(sessionData *session.Data) (string, error) {
	stc := struct {
		Data   map[string]string `json:"data"`
		Expire time.Time         `json:"expire"`
	}{}
	stc.Data = sessionData.All()
	stc.Expire = sessionData.Expire()

	strJson, err := json.Marshal(stc)
	if err != nil {
		return "", err
	}

	return string(strJson), nil
}

func decode(id string, strJson string) (*session.Data, error) {
	stc := struct {
		Data   map[string]string `json:"data"`
		Expire time.Time         `json:"expire"`
	}{
		Data:   make(map[string]string),
		Expire: time.Now(),
	}
	err := json.Unmarshal([]byte(strJson), &stc)
	if err != nil {
		return nil, err
	}

	sessionData := session.NewExistsData(id, stc.Expire.Sub(time.Now()))

	for k, v := range stc.Data {
		sessionData.Set(k, v)
	}
	return sessionData, nil
}

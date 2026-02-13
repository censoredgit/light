package session

import (
	"crypto/rand"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"sync"
	"time"
)

var ErrNotExists = errors.New("not exists")

type Driver interface {
	Init() error
	Open(id string) (*Data, error)
	Close(data *Data) error
}

type Hasher interface {
	Sum(b []byte) string
	BlockSize() int
}

type Manager struct {
	hasher              Hasher
	driver              Driver
	logger              *slog.Logger
	salt                string
	saltLength          int
	ttl                 time.Duration
	cookieName          string
	sessionIdBufferPool sync.Pool
}

func MustSetup(cfg *Config) *Manager {
	if cfg.Hasher == nil {
		panic("Hasher must be provided")
	}
	if cfg.Driver == nil {
		panic("Driver must be provided")
	}
	if cfg.CookieName == "" {
		panic("Cookie name must not be empty")
	}
	if cfg.Salt == "" {
		panic("Salt must not be empty")
	}
	if cfg.TTL == 0 {
		panic("TTL must not be zero")
	}

	if err := cfg.Driver.Init(); err != nil {
		panic(err.Error())
	}

	saltLength := len(cfg.Salt)

	return &Manager{
		hasher:     cfg.Hasher,
		driver:     cfg.Driver,
		logger:     cfg.Logger,
		salt:       cfg.Salt,
		saltLength: saltLength,
		ttl:        cfg.TTL,
		cookieName: cfg.CookieName,
		sessionIdBufferPool: sync.Pool{New: func() any {
			return newIdBuffer(saltLength + cfg.Hasher.BlockSize())
		}},
	}
}

func (m *Manager) InitByRequest(req *http.Request) (*Data, error) {
	c, err := req.Cookie(m.cookieName)
	if err != nil {
		return m.Init("")
	}
	err = c.Valid()
	if err != nil {
		return m.Init("")
	}

	return m.Init(c.Value)
}

func (m *Manager) Init(id string) (*Data, error) {
	var err error

	if id == "" {
		id, err = m.generateId()
		if err != nil {
			return nil, fmt.Errorf("init error: %w", err)
		}
	}

	return m.driver.Open(id)
}

func (m *Manager) Close(data *Data) error {
	return m.driver.Close(data)
}

func (m *Manager) ToCookie(data *Data) *http.Cookie {
	data.isNew = false
	return &http.Cookie{
		Name:    m.cookieName,
		Value:   data.Id(),
		Path:    "/",
		Expires: data.Expire(),
	}
}

func (m *Manager) generateId() (string, error) {
	buf := m.sessionIdBufferPool.Get().(*idBuffer)
	defer m.sessionIdBufferPool.Put(buf)
	buf.Rewind()

	_ = buf.ReadString(m.salt)

	_, err := buf.Read(rand.Reader)
	if err != nil {
		return "", err
	}

	return m.hasher.Sum(buf.Bytes()), nil
}

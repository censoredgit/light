package session

import (
	"log/slog"
	"time"
)

type Config struct {
	Salt       string
	TTL        time.Duration
	CookieName string
	Driver     Driver
	Logger     *slog.Logger
	Hasher     Hasher
}

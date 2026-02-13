package controller

import (
	"github.com/censoredgit/light/session"
	"log/slog"
)

type Config struct {
	Protocol       string
	Host           string
	Port           string
	InternalHost   string
	ExternalHost   string
	Logger         *slog.Logger
	SessionManager *session.Manager
	MaxUploadSize  int64
	Debug          bool
	Templates      struct {
		RootPath string
		Page500  string
		Page404  string
	}
	UserProvider   UserProvider
	StaticPath     string
	CsrfFieldName  string
	LoginRouteName string
}

package controller

import (
	"fmt"
	"io/fs"
	"net/http"
	"strings"
)

type staticFileSystem struct {
	fs http.FileSystem
}

func newStaticFileSystem(fs http.FileSystem) http.FileSystem {
	return &staticFileSystem{
		fs: fs,
	}
}

func (s *staticFileSystem) Open(name string) (http.File, error) {
	if strings.HasSuffix(name, "/") {
		return nil, fs.ErrPermission
	}

	f, err := s.fs.Open(name)
	if err != nil {
		return nil, err
	}

	fStat, err := f.Stat()
	if err != nil {
		iLog.Warn(fmt.Errorf("static file system: stat error %w", err).Error())
		return nil, fs.ErrPermission
	}

	if fStat.IsDir() {
		return nil, fs.ErrPermission
	}

	return f, nil
}

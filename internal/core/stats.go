package core

import (
	"sync/atomic"

	"go-vless-client/internal/config"
)

// Stats атомарно считает трафик текущей сессии.
type Stats struct {
	upload   atomic.Int64
	download atomic.Int64
}

func NewStats() *Stats {
	return &Stats{}
}

func (s *Stats) AddUpload(n int64) {
	s.upload.Add(n)
}

func (s *Stats) AddDownload(n int64) {
	s.download.Add(n)
}

func (s *Stats) Reset() {
	s.upload.Store(0)
	s.download.Store(0)
}

func (s *Stats) Get() config.Stats {
	return config.Stats{
		Upload:   s.upload.Load(),
		Download: s.download.Load(),
	}
}
package internal

import (
	"database/sql"
	"sync/atomic"

	"github.com/blevesearch/bleve"
)

const (
	StatusCreated = iota
	StatusRunned
	StatusStopping
	StatusStopped
)

type Service interface {
	Init(conf *AppConfig, db *sql.DB, index bleve.Index) error
	Run() error
	Name() string
	Stop()
}

type BaseService struct {
	Service

	status uint32
}

func (s *BaseService) Status() uint32 {
	return atomic.LoadUint32(&s.status)
}

func (s *BaseService) SetStatus(v uint32) {
	atomic.StoreUint32(&s.status, v)
}

func (s *BaseService) Stop() {
	s.SetStatus(StatusStopping)
}

func (s *BaseService) IsNeedStop() bool {
	return s.Status() == StatusStopping
}

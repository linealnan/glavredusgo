package service

import (
	"sync/atomic"
)

const (
	StatusCreated = iota
	StatusRunned
	StatusStopping
	StatusStopped
)

type Service interface {
	Init() error
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

package cache

import (
	"errors"
	"sync"
	"time"
)

type TokenCacheService interface {
	Validate(token string) error
	Revoke(token string) error

	Verify(token string) bool
}

type MemoryCacheService struct {
	Locker         sync.RWMutex
	TokenStatusMap map[string]*TokenStatus
	ExpireInterval int64
}

type TokenStatus struct {
	Validity   bool
	InsertTime int64
}

func NewMemoryCacheService(expireInterval int64) *MemoryCacheService {
	return &MemoryCacheService{
		Locker:         sync.RWMutex{},
		TokenStatusMap: make(map[string]*TokenStatus),
		ExpireInterval: expireInterval,
	}
}

func (m *MemoryCacheService) Validate(token string) error {
	m.Locker.Lock()
	defer m.Locker.Unlock()
	if _, ok := m.TokenStatusMap[token]; ok {
		return errors.New("token has already been validated")
	}
	status := &TokenStatus{
		Validity:   true,
		InsertTime: time.Now().Unix(),
	}
	m.TokenStatusMap[token] = status
	return nil
}

func (m *MemoryCacheService) Revoke(token string) error {
	m.Locker.Lock()
	defer m.Locker.Unlock()
	if status, ok := m.TokenStatusMap[token]; ok {
		status.Validity = false
	}
	return nil
}

func (m *MemoryCacheService) Verify(token string) bool {
	m.Locker.RLock()
	defer m.Locker.RUnlock()
	if status, ok := m.TokenStatusMap[token]; ok {
		return status.Validity
	}
	return false
}

func (m *MemoryCacheService) CronJobExpire() {
	nowUnix := time.Now().Unix()

	m.Locker.Lock()
	defer m.Locker.Unlock()
	for _, v := range m.TokenStatusMap {
		if v.InsertTime+m.ExpireInterval < nowUnix {
			v.Validity = false
		}
	}
}

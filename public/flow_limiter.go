package public

import (
	"golang.org/x/time/rate"
	"sync"
)

var FlowLimiterHandler *FlowLimiter

type FlowLimiter struct {
	FlowLimiterMap   map[string]*FlowLimiterItem
	FlowLimiterSlice []*FlowLimiterItem
	Mutex            sync.RWMutex
}

type FlowLimiterItem struct {
	ServiceName string
	Limiter     *rate.Limiter
}

func init() {
	FlowLimiterHandler = NewFlowLimiter()
}

func NewFlowLimiter() *FlowLimiter {
	return &FlowLimiter{
		FlowLimiterMap:   make(map[string]*FlowLimiterItem),
		FlowLimiterSlice: []*FlowLimiterItem{},
		Mutex:            sync.RWMutex{},
	}
}
func (l *FlowLimiter) GetLimiter(serviceName string, qps int64) (*rate.Limiter, error) {
	for _, item := range l.FlowLimiterSlice {
		if item.ServiceName == serviceName {
			return item.Limiter, nil
		}
	}

	newLimiter := rate.NewLimiter(rate.Limit(qps), int(qps*3))
	newItem := &FlowLimiterItem{
		ServiceName: serviceName,
		Limiter:     newLimiter,
	}
	l.FlowLimiterSlice = append(l.FlowLimiterSlice, newItem)
	l.Mutex.Lock()
	defer l.Mutex.Unlock()
	l.FlowLimiterMap[serviceName] = newItem
	return newLimiter, nil
}

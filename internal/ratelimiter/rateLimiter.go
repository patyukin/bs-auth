package ratelimiter

import (
	"fmt"
	"sync"
	"time"
)

type RequestCounter struct {
	requestTimes map[string][]time.Time
	mu           sync.Mutex
}

func NewRequestCounter() *RequestCounter {
	return &RequestCounter{}
}

type RequestCounterInterface interface {
	Increment(key string) error
}

func (c *RequestCounter) Increment(key string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	rt, ok := c.requestTimes[key]
	if !ok {
		c.requestTimes[key] = []time.Time{time.Now()}
		return nil
	}

	// Удаляем все записи, которые старше 1 минуты
	for i, t := range rt {
		if time.Since(t) > 1*time.Minute {
			rt = append(rt[:i], rt[i+1:]...)
		}
	}

	if len(rt) >= 100 {
		return fmt.Errorf("too Many Requests")
	}

	rt = append(rt, time.Now())
	c.requestTimes[key] = rt

	return nil
}

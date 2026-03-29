package service

import (
	"sync"
	"sync/atomic"
)

type ClickCounter struct {
	counts map[int64]*atomic.Int64
	mu     sync.RWMutex
}

func NewClickCounter() *ClickCounter {
	return &ClickCounter{
		counts: make(map[int64]*atomic.Int64),
	}
}

func (c *ClickCounter) Increment(linkID int64) {
	c.mu.RLock()
	counter, exists := c.counts[linkID]
	c.mu.RUnlock()

	if !exists {
		c.mu.Lock()
		if _, exists = c.counts[linkID]; !exists {
			counter = &atomic.Int64{}
			c.counts[linkID] = counter
		} else {
			counter = c.counts[linkID]
		}
		c.mu.Unlock()
	}

	counter.Add(1)
}

func (c *ClickCounter) GetCount(linkID int64) int64 {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if counter, exists := c.counts[linkID]; exists {
		return counter.Load()
	}

	return 0
}

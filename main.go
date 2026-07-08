package vpn

import (
	"slices"
	"sync"
	"time"
)

type VPN struct {
	Host      string
	Country   string
	IsBlocked bool
}

type ControlVPN struct {
	mu sync.Mutex

	ticker   *time.Ticker
	Vpn      []VPN
	capacity int
}

func NewControlVPN(cap int, duration time.Duration) *ControlVPN {
	c := &ControlVPN{
		capacity: cap,
		Vpn:      make([]VPN, 0, cap),
		ticker:   time.NewTicker(time.Second * duration),
	}

	go func() {
		c.cron()
	}()

	return c
}

func (c *ControlVPN) Add(vpn VPN) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	if !slices.Contains(c.Vpn, vpn) {
		c.Vpn = append(c.Vpn, vpn)
		return true
	}
	return false
}

func (c *ControlVPN) Del(vpn VPN) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	idx := slices.Index(c.Vpn, vpn)
	if idx >= 0 {
		slices.Delete(c.Vpn, idx, 1)
		return true
	}
	return false
}

func (c *ControlVPN) cron() {
	for {
		select {
		case <-c.ticker.C:
			c.optimizeCap()
		}
	}
}

func (c *ControlVPN) optimizeCap() {
	c.mu.Lock()
	defer c.mu.Unlock()
	if len(c.Vpn)*3 < cap(c.Vpn) {
		capacity := len(c.Vpn) * 2
		if capacity < c.capacity {
			capacity = c.capacity
		}
		newSlice := make([]VPN, len(c.Vpn), capacity)
		copy(newSlice, c.Vpn)
		c.Vpn = newSlice
	}
}

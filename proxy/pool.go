package proxy

import (
	"fmt"
	"sync"
)

type Pool struct {
	proxies []*Proxy
	mu      sync.RWMutex
}

func NewPool(targets []string) (*Pool, error) {
	proxies := make([]*Proxy, 0, len(targets))

	for _, target := range targets {
		proxy, err := NewProxy(target)
		if err != nil {
			return nil, fmt.Errorf("failed to create proxy %q: %w", target, err)
		}

		proxies = append(proxies, proxy)
	}

	return &Pool{proxies: proxies}, nil
}

func (p *Pool) All() []*Proxy {
	p.mu.RLock()
	defer p.mu.RUnlock()

	copied := make([]*Proxy, len(p.proxies))
	copy(copied, p.proxies)
	return copied
}

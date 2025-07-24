package strategy

import (
	"net/http/httputil"
	"sync"
	"tightrope/proxy"
)

type RoundRobin struct {
	current int
	pool    *proxy.Pool
	mu      sync.Mutex
}

func NewRoundRobin(pool *proxy.Pool) *RoundRobin {
	return &RoundRobin{
		current: 0,
		pool:    pool,
	}
}

func (r *RoundRobin) NextProxy() *httputil.ReverseProxy {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.pool.Len() == 0 {
		return nil
	}

	reverseProxy := r.pool.At(r.current)
	r.current = (r.current + 1) % r.pool.Len()

	return reverseProxy.Backend
}

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
	proxies := r.pool.All()
	if len(proxies) == 0 {
		return nil
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	for i := 0; i < len(proxies); i++ {
		p := proxies[r.current%len(proxies)]
		r.current++
		if p.Healthy.Load() {
			return p.Backend
		}
	}

	return nil
}

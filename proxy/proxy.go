package proxy

import (
	"fmt"
	"net/http/httputil"
	"net/url"
	"sync/atomic"
)

type Proxy struct {
	URL     *url.URL
	Backend *httputil.ReverseProxy
	Healthy atomic.Bool
}

func NewProxy(target string) (*Proxy, error) {
	u, err := url.Parse(target)
	if err != nil {
		return nil, fmt.Errorf("failed to parse target url: %v", err)
	}

	proxy := httputil.NewSingleHostReverseProxy(u)

	return &Proxy{
		URL:     u,
		Backend: proxy,
		Healthy: atomic.Bool{},
	}, nil
}

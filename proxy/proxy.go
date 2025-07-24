package proxy

import (
	"fmt"
	"net/http/httputil"
	"net/url"
)

type Proxy struct {
	Backend *httputil.ReverseProxy
	Healthy bool //Unused, but can be used to track health status
}

func NewProxy(target string) (*Proxy, error) {
	u, err := url.Parse(target)
	if err != nil {
		return nil, fmt.Errorf("failed to parse target url: %v", err)
	}

	proxy := httputil.NewSingleHostReverseProxy(u)

	return &Proxy{
		Backend: proxy,
		Healthy: true,
	}, nil
}

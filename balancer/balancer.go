package balancer

import (
	"log/slog"
	"net/http"
	"net/http/httputil"
)

type Strategy interface {
	NextProxy() *httputil.ReverseProxy
}
type Balancer struct {
	logger   *slog.Logger
	strategy Strategy
}

func New(logger *slog.Logger, strategy Strategy) *Balancer {
	return &Balancer{
		logger:   logger,
		strategy: strategy,
	}
}

func (b *Balancer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	b.logger.Info("Received request", slog.String("method", r.Method), slog.String("url", r.URL.String()))
	proxy := b.strategy.NextProxy()
	if proxy == nil {
		http.Error(w, "No available servers", http.StatusServiceUnavailable)
		return
	}
	b.logger.Info("Forwarding request to proxy")

	proxy.ServeHTTP(w, r)
}

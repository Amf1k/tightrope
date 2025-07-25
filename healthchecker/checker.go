package healthchecker

import (
	"context"
	"log/slog"
	"net"
	"tightrope/proxy"
	"time"
)

type HealthChecker struct {
	logger *slog.Logger
	pool   *proxy.Pool
}

func New(logger *slog.Logger, pool *proxy.Pool) *HealthChecker {
	return &HealthChecker{
		logger: logger,
		pool:   pool,
	}
}

func (hc *HealthChecker) Run(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Second) // Check every 5 seconds
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			proxies := hc.pool.All()
			hc.logger.Info("Running health check for all proxies", slog.Int("count", len(proxies)))
			for _, p := range proxies {
				go hc.check(p)
			}
		case <-ctx.Done():
			return
		}
	}
}

func (hc *HealthChecker) check(p *proxy.Proxy) {
	hc.logger.Info("Checking health of proxy", slog.String("url", p.URL.String()))
	conn, err := net.DialTimeout("tcp", p.URL.Host, 2*time.Second)
	if err != nil {
		hc.logger.Error("Health check failed", slog.String("url", p.URL.String()), slog.Any("error", err))
		p.Healthy.Store(false)
		return
	}
	defer conn.Close()

	hc.logger.Info("Health check succeeded", slog.String("url", p.URL.String()))
	p.Healthy.Store(true)
}

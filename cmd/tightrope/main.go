package main

import (
	"context"
	"errors"
	"flag"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"tightrope/balancer"
	"tightrope/healthchecker"
	"tightrope/proxy"
	"tightrope/strategy"
	"time"
)

func main() {
	var addr string
	flag.StringVar(&addr, "addr", ":8081", "HTTP server address")
	flag.Parse()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	pool, err := proxy.NewPool([]string{"http://localhost:5000", "http://localhost:5001", "http://localhost:5002"})
	if err != nil {
		logger.Error("Failed to create proxy pool", slog.Any("error", err))
		return
	}

	hc := healthchecker.New(logger, pool)

	go func() {
		logger.Info("Starting health checker")
		hc.Run(ctx)
		logger.Info("Health checker stopped")
	}()

	s := strategy.NewRoundRobin(pool)

	b := balancer.New(logger, s)

	server := &http.Server{
		Addr:    addr,
		Handler: b,
	}

	go func() {
		logger.Info("Starting server", slog.String("address", addr))
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			logger.Error("Server error", slog.Any("error", err))
		}
		logger.Info("Server stopped")
	}()

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error("Server shutdown error", slog.Any("error", err))
	}

	logger.Info("Server shutdown")
}

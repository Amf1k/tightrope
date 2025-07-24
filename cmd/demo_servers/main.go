package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	var count int
	flag.IntVar(&count, "count", 3, "Number of servers to start")

	flag.Parse()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	if count <= 0 {
		panic("count must be greater than 0")
	}

	logger.Info("Starting servers", slog.Int("count", count))
	for i := 0; i < count; i++ {
		go startServer(logger, 5000+i) // Start servers on ports 5000, 5001, ..., 5000 + count - 1
	}

	<-ctx.Done()
	logger.Info("Received shutdown signal, stopping servers")
}

func startServer(logger *slog.Logger, port int) {
	logger.Info("Starting servers", slog.Int("port", port))

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		logger.Info("Received request", slog.String("method", r.Method), slog.String("url", r.URL.String()), slog.Int("port", port))
		w.WriteHeader(http.StatusOK)
		res := fmt.Sprintf("Hello From Server: %d", port)
		w.Write([]byte(res))
	})

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server on port %d failed: %v", port, err)
	}

}

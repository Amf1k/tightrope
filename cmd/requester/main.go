package main

import (
	"flag"
	"io"
	"log/slog"
	"net/http"
	"os"
)

func main() {
	var url string
	flag.StringVar(&url, "url", "http://localhost:8081", "URL to send requests to")
	var count int
	flag.IntVar(&count, "count", 10, "Number of requests to send")

	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	client := &http.Client{}

	requestFunc := func(idx int) {
		logger.Info("Sending request", slog.Int("request_index", idx), slog.String("url", url))
		resp, err := client.Get(url)
		if err != nil {
			logger.Error("Failed to send request", slog.Int("request_index", idx), slog.Any("error", err))
			return
		}
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)

		logger.Info("Received response", slog.Int("request_index", idx), slog.Int("status_code", resp.StatusCode), slog.String("body", string(body)))
	}

	for i := 0; i < count; i++ {
		requestFunc(i)
	}
}

package main

import (
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	r := chi.NewRouter()

	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok\n"))
	})

	r.Handle("/metrics", promhttp.Handler())

	// Stub JSON endpoints backed by downstream gRPC in future
	r.Get("/api/fairrent/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"SERVING"}`))
	})

	r.Get("/api/zklease/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"SERVING"}`))
	})

	addr := ":8080"
	if v := os.Getenv("GATEWAY_ADDR"); v != "" {
		addr = v
	}
	logger.Info("gateway listening", zap.String("addr", addr))
	if err := http.ListenAndServe(addr, r); err != nil {
		logger.Fatal("listen", zap.Error(err))
	}
}

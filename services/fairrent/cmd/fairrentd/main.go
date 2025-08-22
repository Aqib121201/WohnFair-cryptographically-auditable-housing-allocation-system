package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/wohnfair/wohnfair/services/fairrent/api"
	"github.com/wohnfair/wohnfair/services/fairrent/internal/scheduler"
	"github.com/wohnfair/wohnfair/services/fairrent/internal/telemetry"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	port     = flag.Int("port", 50051, "gRPC server port")
	logLevel = flag.String("log-level", "info", "Log level (debug, info, warn, error)")
	configFile = flag.String("config", "", "Configuration file path")
)

func main() {
	flag.Parse()

	// Initialize logger
	logger := initLogger()
	defer logger.Sync()

	logger.Info("Starting FairRent service",
		zap.Int("port", *port),
		zap.String("log_level", *logLevel),
	)

	// Initialize telemetry
	if err := telemetry.InitTracer("fairrent", "0.1.0"); err != nil {
		logger.Fatal("Failed to initialize tracer", zap.Error(err))
	}

	// Load configuration
	config := scheduler.DefaultConfig()
	if *configFile != "" {
		if err := loadConfig(*configFile, config); err != nil {
			logger.Fatal("Failed to load configuration", zap.Error(err))
		}
	}

	// Create scheduler
	scheduler := scheduler.NewFairRent(config, logger)

	// Create and start server
	server := api.NewServer(scheduler, logger, *port)

	// Start server in goroutine
	go func() {
		if err := server.Start(); err != nil {
			logger.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	// Start metrics server
	go startMetricsServer(logger)

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down FairRent service")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	server.Stop()

	logger.Info("FairRent service stopped")
}

// initLogger initializes the logger
func initLogger() *zap.Logger {
	var level zapcore.Level
	switch *logLevel {
	case "debug":
		level = zapcore.DebugLevel
	case "info":
		level = zapcore.InfoLevel
	case "warn":
		level = zapcore.WarnLevel
	case "error":
		level = zapcore.ErrorLevel
	default:
		level = zapcore.InfoLevel
	}

	config := zap.NewProductionConfig()
	config.Level = zap.NewAtomicLevelAt(level)
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	logger, err := config.Build()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}

	return logger
}

// loadConfig loads configuration from file
func loadConfig(configFile string, config *scheduler.Config) error {
	// TODO: Implement configuration loading from YAML/JSON
	// For now, just log that we're using defaults
	return nil
}

// startMetricsServer starts the Prometheus metrics server
func startMetricsServer(logger *zap.Logger) {
	// This would typically run on a different port
	// For now, we'll just log that metrics are available via gRPC
	logger.Info("Prometheus metrics available via gRPC server")
	
	// In a real deployment, you might want a separate HTTP server for metrics
	// mux := http.NewServeMux()
	// mux.Handle("/metrics", promhttp.Handler())
	// 
	// server := &http.Server{
	//     Addr:    ":9090",
	//     Handler: mux,
	// }
	// 
	// if err := server.ListenAndServe(); err != nil {
	//     logger.Error("Metrics server failed", zap.Error(err))
	// }
}

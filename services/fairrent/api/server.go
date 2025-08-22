package api

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/wohnfair/wohnfair/services/gen/wohnfair/fairrent/v1"
	"github.com/wohnfair/wohnfair/services/fairrent/internal/scheduler"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Server implements the FairRent gRPC service
type Server struct {
	fairrentv1.UnimplementedFairRentServiceServer
	
	// Dependencies
	scheduler *scheduler.FairRent
	logger    *zap.Logger
	
	// gRPC server
	grpcServer *grpc.Server
	healthServer *health.Server
	
	// Configuration
	port int
}

// NewServer creates a new FairRent server
func NewServer(scheduler *scheduler.FairRent, logger *zap.Logger, port int) *Server {
	// Create gRPC server with middleware
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_prometheus.UnaryServerInterceptor,
			otelgrpc.UnaryServerInterceptor(),
		)),
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			grpc_prometheus.StreamServerInterceptor,
			otelgrpc.StreamServerInterceptor(),
		)),
	)
	
	// Create health server
	healthServer := health.NewServer()
	
	server := &Server{
		scheduler:    scheduler,
		logger:       logger,
		grpcServer:   grpcServer,
		healthServer: healthServer,
		port:         port,
	}
	
	// Register services
	fairrentv1.RegisterFairRentServiceServer(grpcServer, server)
	grpc_health_v1.RegisterHealthServer(grpcServer, healthServer)
	
	// Enable reflection for development
	reflection.Register(grpcServer)
	
	// Register Prometheus metrics
	grpc_prometheus.Register(grpcServer)
	
	// Set health status
	healthServer.SetServingStatus("", grpc_health_v1.HealthCheckResponse_SERVING)
	
	return server
}

// Start starts the gRPC server
func (s *Server) Start() error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}
	
	s.logger.Info("Starting FairRent gRPC server",
		zap.Int("port", s.port),
	)
	
	// Start metrics server
	go s.startMetricsServer()
	
	// Start gRPC server
	if err := s.grpcServer.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %w", err)
	}
	
	return nil
}

// Stop gracefully stops the server
func (s *Server) Stop() {
	s.logger.Info("Stopping FairRent gRPC server")
	
	// Set health status to not serving
	s.healthServer.SetServingStatus("", grpc_health_v1.HealthCheckResponse_NOT_SERVING)
	
	// Graceful shutdown
	if s.grpcServer != nil {
		s.grpcServer.GracefulStop()
	}
}

// startMetricsServer starts the Prometheus metrics server
func (s *Server) startMetricsServer() {
	// This would typically run on a different port
	// For now, we'll just log that metrics are available via gRPC
	s.logger.Info("Prometheus metrics available via gRPC server")
}

// Enqueue implements the Enqueue RPC method
func (s *Server) Enqueue(ctx context.Context, req *fairrentv1.EnqueueRequest) (*fairrentv1.EnqueueResponse, error) {
	start := time.Now()
	
	s.logger.Info("Enqueue request received",
		zap.String("user_id", req.UserId.Value),
		zap.String("user_group", req.UserGroup.String()),
		zap.Int("urgency", int(req.Urgency)),
	)
	
	// Validate request
	if err := s.validateEnqueueRequest(req); err != nil {
		s.logger.Error("Enqueue request validation failed",
			zap.Error(err),
			zap.String("user_id", req.UserId.Value),
		)
		return nil, err
	}
	
	// Process request
	resp, err := s.scheduler.Enqueue(ctx, req)
	if err != nil {
		s.logger.Error("Failed to enqueue request",
			zap.Error(err),
			zap.String("user_id", req.UserId.Value),
		)
		return nil, err
	}
	
	// Log success
	s.logger.Info("Request enqueued successfully",
		zap.String("ticket_id", resp.TicketId.Value),
		zap.String("user_id", req.UserId.Value),
		zap.Duration("processing_time", time.Since(start)),
	)
	
	return resp, nil
}

// ScheduleNext implements the ScheduleNext RPC method
func (s *Server) ScheduleNext(ctx context.Context, req *fairrentv1.ScheduleNextRequest) (*fairrentv1.ScheduleNextResponse, error) {
	start := time.Now()
	
	s.logger.Info("ScheduleNext request received",
		zap.Int("available_properties", len(req.AvailableProperties)),
	)
	
	// Process request
	resp, err := s.scheduler.ScheduleNext(ctx, req)
	if err != nil {
		s.logger.Error("Failed to schedule next request",
			zap.Error(err),
		)
		return nil, err
	}
	
	// Log success
	s.logger.Info("Next request scheduled successfully",
		zap.String("ticket_id", resp.TicketId.Value),
		zap.String("user_id", resp.UserId.Value),
		zap.Float64("fairness_score", resp.FairnessScore),
		zap.Duration("processing_time", time.Since(start)),
	)
	
	return resp, nil
}

// PeekPosition implements the PeekPosition RPC method
func (s *Server) PeekPosition(ctx context.Context, req *fairrentv1.PeekPositionRequest) (*fairrentv1.PeekPositionResponse, error) {
	s.logger.Debug("PeekPosition request received",
		zap.String("ticket_id", req.TicketId.Value),
	)
	
	// Process request
	resp, err := s.scheduler.PeekPosition(ctx, req)
	if err != nil {
		s.logger.Error("Failed to peek position",
			zap.Error(err),
			zap.String("ticket_id", req.TicketId.Value),
		)
		return nil, err
	}
	
	return resp, nil
}

// GetMetrics implements the GetMetrics RPC method
func (s *Server) GetMetrics(ctx context.Context, req *fairrentv1.GetMetricsRequest) (*fairrentv1.FairnessMetrics, error) {
	s.logger.Debug("GetMetrics request received")
	
	// Get metrics from scheduler
	metrics, err := s.scheduler.GetMetrics(ctx)
	if err != nil {
		s.logger.Error("Failed to get metrics",
			zap.Error(err),
		)
		return nil, err
	}
	
	return metrics, nil
}

// UpdateRequest implements the UpdateRequest RPC method
func (s *Server) UpdateRequest(ctx context.Context, req *fairrentv1.UpdateRequestRequest) (*fairrentv1.UpdateRequestResponse, error) {
	s.logger.Info("UpdateRequest received",
		zap.String("ticket_id", req.TicketId.Value),
	)
	
	// TODO: Implement update logic
	return nil, fmt.Errorf("UpdateRequest not yet implemented")
}

// CancelRequest implements the CancelRequest RPC method
func (s *Server) CancelRequest(ctx context.Context, req *fairrentv1.CancelRequestRequest) (*fairrentv1.CancelRequestResponse, error) {
	s.logger.Info("CancelRequest received",
		zap.String("ticket_id", req.TicketId.Value),
		zap.String("reason", req.Reason),
	)
	
	// TODO: Implement cancellation logic
	return nil, fmt.Errorf("CancelRequest not yet implemented")
}

// GetQueueStatus implements the GetQueueStatus RPC method
func (s *Server) GetQueueStatus(ctx context.Context, req *fairrentv1.GetQueueStatusRequest) (*fairrentv1.QueueStatus, error) {
	s.logger.Debug("GetQueueStatus request received")
	
	// TODO: Implement queue status logic
	return nil, fmt.Errorf("GetQueueStatus not yet implemented")
}

// Health implements the Health RPC method
func (s *Server) Health(ctx context.Context, req *fairrentv1.HealthRequest) (*fairrentv1.HealthResponse, error) {
	return &fairrentv1.HealthResponse{
		Status:  fairrentv1.HealthResponse_SERVING,
		Message: "FairRent service is healthy",
		Timestamp: &timestamppb.Timestamp{
			Seconds: time.Now().Unix(),
		},
	}, nil
}

// validateEnqueueRequest validates the enqueue request
func (s *Server) validateEnqueueRequest(req *fairrentv1.EnqueueRequest) error {
	if req.UserId == nil || req.UserId.Value == "" {
		return fmt.Errorf("user_id is required")
	}
	
	if req.UserGroup == fairrentv1.UserGroup_USER_GROUP_UNSPECIFIED {
		return fmt.Errorf("user_group is required")
	}
	
	if req.Urgency == fairrentv1.UrgencyLevel_URGENCY_LEVEL_UNSPECIFIED {
		return fmt.Errorf("urgency level is required")
	}
	
	if req.FinancialConstraints != nil {
		if req.FinancialConstraints.MaxMonthlyRent <= 0 {
			return fmt.Errorf("max_monthly_rent must be positive")
		}
	}
	
	return nil
}

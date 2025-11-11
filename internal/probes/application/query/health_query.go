package query

import (
	"context"

	"github.com/MaisamV/wallet/internal/probes/entity"
	"github.com/MaisamV/wallet/internal/probes/ports"
	"github.com/MaisamV/wallet/platform/logger"
)

// GetHealthQuery represents a query to get system health status
type GetHealthQuery struct{}

// GetHealthQueryHandler handles health check queries
type GetHealthQueryHandler struct {
	logger          logger.Logger
	databaseChecker ports.DatabaseChecker
}

// NewGetHealthQueryHandler creates a new health query handler
func NewGetHealthQueryHandler(logger logger.Logger, databaseChecker ports.DatabaseChecker) *GetHealthQueryHandler {
	return &GetHealthQueryHandler{
		logger:          logger,
		databaseChecker: databaseChecker,
	}
}

// Handle executes the health check query
func (h *GetHealthQueryHandler) Handle(ctx context.Context, query GetHealthQuery) (*entity.HealthResponse, error) {
	h.logger.Info().Msg("Starting health check")
	response := entity.NewHealthResponse()

	// Check database connectivity
	if h.databaseChecker != nil {
		h.logger.Debug().Msg("Checking database connectivity")
		dbHealthy, dbResponseTime, err := h.databaseChecker.CheckDatabase(ctx)
		if err != nil {
			h.logger.Error().Err(err).Msg("Database health check failed")
			response.AddCheck("database", entity.CheckStatusDown, 0)
		} else {
			status := entity.CheckStatusUp
			if !dbHealthy {
				status = entity.CheckStatusDown
				h.logger.Warn().Msg("Database is not healthy")
			} else {
				h.logger.Info().Int64("response_time_ms", dbResponseTime.Milliseconds()).Msg("Database health check passed")
			}
			response.AddCheck("database", status, dbResponseTime.Milliseconds())
		}
	}

	// Determine overall status
	response.DetermineOverallStatus()
	h.logger.Info().Bool("is_healthy", response.IsHealthy()).Msg("Health check completed")

	return response, nil
}

// HealthService implements the HealthService port
type HealthService struct {
	logger       logger.Logger
	queryHandler *GetHealthQueryHandler
}

// NewHealthService creates a new health service
func NewHealthService(logger logger.Logger, queryHandler *GetHealthQueryHandler) *HealthService {
	return &HealthService{
		logger:       logger,
		queryHandler: queryHandler,
	}
}

// GetHealthStatus returns the current health status
func (s *HealthService) GetHealthStatus(ctx context.Context) (*entity.HealthResponse, error) {
	s.logger.Debug().Msg("Health status requested")
	return s.queryHandler.Handle(ctx, GetHealthQuery{})
}

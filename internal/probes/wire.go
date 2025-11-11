package probes

import (
	healthQuery "github.com/MaisamV/wallet/internal/probes/application/query"
	pingQuery "github.com/MaisamV/wallet/internal/probes/application/query"
	healthInfra "github.com/MaisamV/wallet/internal/probes/infrastructure"
	healthHttp "github.com/MaisamV/wallet/internal/probes/presentation/http"
	pingHttp "github.com/MaisamV/wallet/internal/probes/presentation/http"
	"github.com/MaisamV/wallet/platform/config"
	"github.com/MaisamV/wallet/platform/logger"
	"github.com/google/wire"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ProvidePingQueryHandler provides a ping query handler
func ProvidePingQueryHandler(logger logger.Logger) *pingQuery.PingQueryHandler {
	return pingQuery.NewPingQueryHandler(logger)
}

// ProvidePingHandler provides a ping HTTP handler
func ProvidePingHandler(logger logger.Logger, pingQueryHandler *pingQuery.PingQueryHandler) *pingHttp.PingHandler {
	return pingHttp.NewPingHandler(logger, pingQueryHandler)
}

// ProvideDatabaseChecker provides a database checker
func ProvideDatabaseChecker(cfg *config.Config, logger logger.Logger, db *pgxpool.Pool) *healthInfra.DatabaseChecker {
	return healthInfra.NewDatabaseChecker(logger, db, cfg.Health.DatabaseTimeout)
}

// ProvideHealthQueryHandler provides a health query handler
func ProvideHealthQueryHandler(logger logger.Logger, databaseChecker *healthInfra.DatabaseChecker) *healthQuery.GetHealthQueryHandler {
	return healthQuery.NewGetHealthQueryHandler(logger, databaseChecker)
}

// ProvideHealthService provides a health service
func ProvideHealthService(logger logger.Logger, healthQueryHandler *healthQuery.GetHealthQueryHandler) *healthQuery.HealthService {
	return healthQuery.NewHealthService(logger, healthQueryHandler)
}

// ProvideLivenessQueryHandler provides a liveness query handler
func ProvideLivenessQueryHandler(logger logger.Logger) *healthQuery.GetLivenessQueryHandler {
	return healthQuery.NewGetLivenessQueryHandler(logger)
}

// ProvideLivenessService provides a liveness service
func ProvideLivenessService(logger logger.Logger, livenessQueryHandler *healthQuery.GetLivenessQueryHandler) *healthQuery.LivenessService {
	return healthQuery.NewLivenessService(logger, livenessQueryHandler)
}

// ProvideHealthHandler provides a health HTTP handler
func ProvideHealthHandler(logger logger.Logger, healthService *healthQuery.HealthService, livenessService *healthQuery.LivenessService) *healthHttp.HealthHandler {
	return healthHttp.NewHealthHandler(logger, healthService, livenessService)
}

// ProbesSet is a wire provider set for all probes dependencies
var ProbesSet = wire.NewSet(
	ProvidePingQueryHandler,
	ProvidePingHandler,
	ProvideDatabaseChecker,
	ProvideHealthQueryHandler,
	ProvideHealthService,
	ProvideLivenessQueryHandler,
	ProvideLivenessService,
	ProvideHealthHandler,
)

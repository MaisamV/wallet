package platform

import (
	"github.com/MaisamV/wallet/platform/config"
	"github.com/MaisamV/wallet/platform/database"
	"github.com/MaisamV/wallet/platform/http"
	"github.com/MaisamV/wallet/platform/logger"
	"github.com/google/wire"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ProvideLogger provides a logger instance
func ProvideLogger(cfg *config.Config) logger.Logger {
	return logger.NewWithLevel(cfg.Logging.Level)
}

// ProvideConfig provides a configuration instance
func ProvideConfig() (*config.Config, error) {
	return config.Load()
}

// ProvideDatabase provides a database connection pool
func ProvideDatabase(cfg *config.Config, log logger.Logger) (*pgxpool.Pool, error) {
	return database.NewConnection(cfg.Database, log)
}

// ProvideHTTPServer provides an HTTP server instance
func ProvideHTTPServer(cfg *config.Config, log logger.Logger) *http.Server {
	return http.NewServer(cfg.Server, log)
}

// PlatformSet is a wire provider set for all platform dependencies
var PlatformSet = wire.NewSet(
	ProvideLogger,
	ProvideConfig,
	ProvideDatabase,
	ProvideHTTPServer,
)

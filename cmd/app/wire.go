//go:build wireinject
// +build wireinject

package main

import (
	"github.com/MaisamV/wallet/internal/probes"
	probesHttp "github.com/MaisamV/wallet/internal/probes/presentation/http"
	"github.com/MaisamV/wallet/internal/swagger"
	swaggerHttp "github.com/MaisamV/wallet/internal/swagger/presentation/http"
	wallet "github.com/MaisamV/wallet/internal/wallet"
	infrastructure "github.com/MaisamV/wallet/internal/wallet/infrastructure/repo"
	walletHttp "github.com/MaisamV/wallet/internal/wallet/presentation/http"
	"github.com/MaisamV/wallet/platform"
	"github.com/MaisamV/wallet/platform/config"
	"github.com/MaisamV/wallet/platform/http"
	"github.com/MaisamV/wallet/platform/logger"
	"github.com/google/wire"
)

// Application holds all the application dependencies
type Application struct {
	Config     *config.Config
	Logger     logger.Logger
	HTTPServer *http.Server
	Probes     *ProbesModule
	Swagger    *SwaggerModule
	Wallet     *WalletModule
}

// ProbesModule holds all probes-related dependencies
type ProbesModule struct {
	PingHandler   *probesHttp.PingHandler
	HealthHandler *probesHttp.HealthHandler
}

// SwaggerModule holds all swagger-related dependencies
type SwaggerModule struct {
	DocsHandler *swaggerHttp.DocsHandler
}

type WalletModule struct {
	WalletHandler *walletHttp.WalletHandler
	Repo          *infrastructure.PgxWalletRepo
}

// InitializeApplication creates and initializes the application with all dependencies
func InitializeApplication() (*Application, error) {
	wire.Build(
		// Platform providers
		platform.PlatformSet,

		// Internal module providers
		probes.ProbesSet,
		swagger.SwaggerSet,
		wallet.WalletSet,

		// Application structure providers
		ProvideProbesModule,
		ProvideSwaggerModule,
		ProvideWalletModule,
		ProvideApplication,
	)
	return &Application{}, nil
}

// ProvideProbesModule provides the probes module
func ProvideProbesModule(
	pingHandler *probesHttp.PingHandler,
	healthHandler *probesHttp.HealthHandler,
) *ProbesModule {
	return &ProbesModule{
		PingHandler:   pingHandler,
		HealthHandler: healthHandler,
	}
}

// ProvideSwaggerModule provides the swagger module
func ProvideSwaggerModule(
	docsHandler *swaggerHttp.DocsHandler,
) *SwaggerModule {
	return &SwaggerModule{
		DocsHandler: docsHandler,
	}
}

// ProvideSwaggerModule provides the swagger module
func ProvideWalletModule(
	handler *walletHttp.WalletHandler,
	repo *infrastructure.PgxWalletRepo,
) *WalletModule {
	return &WalletModule{
		WalletHandler: handler,
		Repo:          repo,
	}
}

// ProvideApplication provides the main application structure
func ProvideApplication(
	config *config.Config,
	logger logger.Logger,
	httpServer *http.Server,
	probesModule *ProbesModule,
	swaggerModule *SwaggerModule,
	walletModule *WalletModule,
) *Application {
	return &Application{
		Config:     config,
		Logger:     logger,
		HTTPServer: httpServer,
		Probes:     probesModule,
		Swagger:    swaggerModule,
		Wallet:     walletModule,
	}
}

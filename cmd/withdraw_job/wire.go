//go:build wireinject
// +build wireinject

package main

import (
	wallet "github.com/MaisamV/wallet/internal/wallet"
	"github.com/MaisamV/wallet/internal/wallet/application/command"
	infrastructure "github.com/MaisamV/wallet/internal/wallet/infrastructure/repo"
	"github.com/MaisamV/wallet/platform"
	"github.com/MaisamV/wallet/platform/config"
	"github.com/MaisamV/wallet/platform/logger"
	"github.com/google/wire"
)

// Application holds all the application dependencies
type Application struct {
	Config *config.Config
	Logger logger.Logger
	Wallet *WalletModule
}

type WalletModule struct {
	WithdrawHandler *command.WithdrawCommandHandler
	Repo            *infrastructure.PgxWalletRepo
}

// InitializeApplication creates and initializes the application with all dependencies
func InitializeApplication() (*Application, error) {
	wire.Build(
		// Platform providers
		platform.PlatformSet,

		wallet.WalletSet,

		// Application structure providers
		ProvideWalletModule,
		ProvideApplication,
	)
	return &Application{}, nil
}

// ProvideSwaggerModule provides the swagger module
func ProvideWalletModule(
	handler *command.WithdrawCommandHandler,
	repo *infrastructure.PgxWalletRepo,
) *WalletModule {
	return &WalletModule{
		WithdrawHandler: handler,
		Repo:            repo,
	}
}

// ProvideApplication provides the main application structure
func ProvideApplication(
	config *config.Config,
	logger logger.Logger,
	walletModule *WalletModule,
) *Application {
	return &Application{
		Config: config,
		Logger: logger,
		Wallet: walletModule,
	}
}

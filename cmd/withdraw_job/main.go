package main

import (
	"context"
	"github.com/MaisamV/wallet/internal/wallet/application/command"
	"github.com/MaisamV/wallet/platform/logger"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	app, err := InitializeApplication()
	if err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}
	defer app.Wallet.Repo.Close()
	app.Logger.Info().Msg("Starting release job")
	releaseConfig := app.Config.Release

	worker := NewWithdrawWorker(
		app.Wallet.WithdrawHandler,
		app.Logger,
		releaseConfig.Interval,
		releaseConfig.BatchSize,
		releaseConfig.WorkerCount,
	)
	worker.Start()

	// Gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	app.Logger.Info().Msg("Finished gracefully")
}

type WithdrawWorker struct {
	handler     *command.WithdrawCommandHandler
	logger      logger.Logger
	interval    time.Duration
	batchSize   int
	workerCount int
	stop        chan any
}

func NewWithdrawWorker(handler *command.WithdrawCommandHandler, logger logger.Logger, interval time.Duration, batchSize int, workerCount int) *WithdrawWorker {
	return &WithdrawWorker{
		handler:     handler,
		logger:      logger,
		interval:    interval,
		batchSize:   batchSize,
		workerCount: workerCount,
		stop:        make(chan any),
	}
}

func (w *WithdrawWorker) Start() {
	for i := 0; i < 1; i++ {
		go w.workerLoop(i)
	}
}

func (w *WithdrawWorker) Stop() {
	close(w.stop)
}

func (w *WithdrawWorker) workerLoop(id int) {
	ticker := time.NewTicker(w.interval)
	w.logger.Info().Int("worker", id).Msg("started")

	for {
		select {
		case <-ticker.C:
			w.withdraw(id)

		case <-w.stop:
			w.logger.Info().Int("worker", id).Msg("stopped")
			return
		}
	}
}

func (w *WithdrawWorker) withdraw(id int) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	w.handler.Start()

	cmd := command.WithdrawCommand{Limit: w.batchSize}
	err := w.handler.Handle(ctx, cmd)
	if err != nil {
		w.logger.Error().Err(err).Int("worker", id).Msg("withdraw failed")
		return
	}
}

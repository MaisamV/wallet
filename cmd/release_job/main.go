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

	worker := NewReleaseWorker(
		app.Wallet.ReleaseHandler,
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

type ReleaseWorker struct {
	handler     *command.ReleaseCommandHandler
	logger      logger.Logger
	interval    time.Duration
	batchSize   int
	workerCount int
	stop        chan any
}

func NewReleaseWorker(
	handler *command.ReleaseCommandHandler,
	logger logger.Logger,
	interval time.Duration,
	batchSize int,
	workerCount int,
) *ReleaseWorker {
	return &ReleaseWorker{
		handler:     handler,
		logger:      logger,
		interval:    interval,
		batchSize:   batchSize,
		workerCount: workerCount,
		stop:        make(chan any),
	}
}

func (w *ReleaseWorker) Start() {
	for i := 0; i < w.workerCount; i++ {
		go w.workerLoop(i)
	}
}

func (w *ReleaseWorker) Stop() {
	close(w.stop)
}

func (w *ReleaseWorker) workerLoop(id int) {
	ticker := time.NewTicker(w.interval)
	w.logger.Info().Int("worker", id).Msg("started")

	for {
		select {
		case <-ticker.C:
			w.release(id)

		case <-w.stop:
			w.logger.Info().Int("worker", id).Msg("stopped")
			return
		}
	}
}

func (w *ReleaseWorker) release(id int) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	cmd := command.ReleaseCommand{BatchSize: w.batchSize}
	transactions, err := w.handler.Handle(ctx, cmd)
	if err != nil {
		w.logger.Error().Err(err).Int("worker", id).Msg("release failed")
		return
	}

	for _, txn := range transactions {
		w.logger.Info().Str("ID", txn.ID.String()).Msg("transaction released")
	}
}

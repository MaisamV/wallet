package infrastructure

import (
	"context"
	"fmt"
	"github.com/MaisamV/wallet/platform/config"
	"github.com/MaisamV/wallet/platform/database"
	"github.com/MaisamV/wallet/platform/logger"
	"github.com/gofrs/uuid/v5"
	"sync"
	"testing"
	"time"
)

type Job struct {
	ID   int64
	UUID uuid.UUID
}

// BenchmarkCharge_multipleUsersConcurrent benchmarks the WalletRepo.Charge function for different users
func BenchmarkCharge_multipleUsersConcurrent(b *testing.B) {
	noopLogger := logger.NewNoopLogger()
	cfg, err := config.Load()
	if err != nil {
		fmt.Println(err)
	}
	pool, err := database.NewConnection(cfg.Database, noopLogger)
	repo := NewPgxWalletRepo(noopLogger, pool)
	const jobsNum = 1000
	const maxWorkers = 35

	uuids := make([]uuid.UUID, jobsNum)
	for i := 0; i < jobsNum; i++ {
		u, err := uuid.NewV7()
		if err != nil {
			b.Fatalf("failed to generate UUID: %v", err)
		}
		uuids[i] = u
	}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		var allFinishedWg sync.WaitGroup
		allFinishedWg.Add(maxWorkers)
		jobs := make(chan Job, jobsNum)
		// Prepare jobs
		for i := int64(0); i < jobsNum; i++ {
			jobs <- Job{ID: i, UUID: uuids[i]}
		}
		close(jobs)

		var allStartedWg sync.WaitGroup
		allStartedWg.Add(maxWorkers)
		barrier := make(chan any)

		for w := 0; w < maxWorkers; w++ {
			go func(workerID int) {
				//wait until all goroutines start
				allStartedWg.Done()
				<-barrier

				for job := range jobs {
					_, _ = repo.Charge(context.Background(), job.ID, &job.UUID, 1000, nil)
				}
				allFinishedWg.Done()
			}(w)
		}
		//wait until all goroutines start
		allStartedWg.Wait()
		//close the barrier to let the goroutines start processing together
		start := time.Now()
		close(barrier)

		// Wait for all jobs to complete
		allFinishedWg.Wait()
		totalTime := time.Since(start)

		b.Logf("Run %d: TPS: %.2f: total time %v", n+1, float64(jobsNum)/totalTime.Seconds(), totalTime)
	}
}

// BenchmarkCharge_singleUsersConcurrent benchmarks the WalletRepo.Charge function for a single users
func BenchmarkCharge_singleUsersConcurrent(b *testing.B) {
	noopLogger := logger.NewNoopLogger()
	cfg, err := config.Load()
	if err != nil {
		fmt.Println(err)
	}
	pool, err := database.NewConnection(cfg.Database, noopLogger)
	repo := NewPgxWalletRepo(noopLogger, pool)
	const jobsNum = 1000
	const maxWorkers = 35

	uuids := make([]uuid.UUID, jobsNum)
	for i := 0; i < jobsNum; i++ {
		u, err := uuid.NewV7()
		if err != nil {
			b.Fatalf("failed to generate UUID: %v", err)
		}
		uuids[i] = u
	}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		var allFinishedWg sync.WaitGroup
		allFinishedWg.Add(maxWorkers)
		jobs := make(chan Job, jobsNum)
		// Prepare jobs
		for i := int64(0); i < jobsNum; i++ {
			jobs <- Job{ID: i, UUID: uuids[i]}
		}
		close(jobs)

		var allStartedWg sync.WaitGroup
		allStartedWg.Add(maxWorkers)
		barrier := make(chan any)

		for w := 0; w < maxWorkers; w++ {
			go func(workerID int) {
				//wait until all goroutines start
				allStartedWg.Done()
				<-barrier

				for job := range jobs {
					// Charging only user 0 wallet
					_, _ = repo.Charge(context.Background(), 0, &job.UUID, 1000, nil)
				}
				allFinishedWg.Done()
			}(w)
		}
		//wait until all goroutines start
		allStartedWg.Wait()
		//close the barrier to let the goroutines start processing together
		start := time.Now()
		close(barrier)

		// Wait for all jobs to complete
		allFinishedWg.Wait()
		totalTime := time.Since(start)

		b.Logf("Run %d: TPS: %.2f: total time %v", n+1, float64(jobsNum)/totalTime.Seconds(), totalTime)
	}
}

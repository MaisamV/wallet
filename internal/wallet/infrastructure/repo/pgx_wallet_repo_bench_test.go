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
	var err error
	repo := Init()
	defer repo.Close()
	const jobsNum = 1000
	const maxWorkers = 35

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		uuids := make([]uuid.UUID, jobsNum)
		for i := 0; i < jobsNum; i++ {
			u, err := uuid.NewV7()
			if err != nil {
				b.Fatalf("failed to generate UUID: %v", err)
			}
			uuids[i] = u
		}

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

		failedCount := 0
		for w := 0; w < maxWorkers; w++ {
			go func(workerID int) {
				//wait until all goroutines start
				allStartedWg.Done()
				<-barrier

				for job := range jobs {
					_, err = repo.Charge(context.Background(), job.ID, &job.UUID, 1000, nil)
					if err != nil {
						failedCount++
						fmt.Printf("Failed: %v", err)
					}
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

		b.Logf("\nRun %d\nTPS: %.2f\nFailed Count: %d\ntotal time %v", n+1, float64(jobsNum)/totalTime.Seconds(), failedCount, totalTime)
	}
}

// BenchmarkCharge_singleUsersConcurrent benchmarks the WalletRepo.Charge function for a single users
func BenchmarkCharge_singleUsersConcurrent(b *testing.B) {
	var err error
	repo := Init()
	defer repo.Close()
	const jobsNum = 1000
	const maxWorkers = 35

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		uuids := make([]uuid.UUID, jobsNum)
		for i := 0; i < jobsNum; i++ {
			u, err := uuid.NewV7()
			if err != nil {
				b.Fatalf("failed to generate UUID: %v", err)
			}
			uuids[i] = u
		}

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

		failedCount := 0
		for w := 0; w < maxWorkers; w++ {
			go func(workerID int) {
				//wait until all goroutines start
				allStartedWg.Done()
				<-barrier

				for job := range jobs {
					// Charging only user 0 wallet
					_, err = repo.Charge(context.Background(), 0, &job.UUID, 1000, nil)
					if err != nil {
						failedCount++
						fmt.Printf("Failed: %v", err)
					}
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

		b.Logf("\nRun %d\nTPS: %.2f\nFailed Count: %d\ntotal time %v", n+1, float64(jobsNum)/totalTime.Seconds(), failedCount, totalTime)
	}
}

// BenchmarkDebit_singleUsersConcurrent benchmarks the WalletRepo.Debit function for a single users
func BenchmarkDebit_singleUsersConcurrent(b *testing.B) {
	var err error
	repo := Init()
	defer repo.Close()

	u, err := uuid.NewV7()
	_, err = repo.Charge(context.Background(), 0, &u, 1000000, nil)
	if err != nil {
		panic(err)
	}

	const jobsNum = 1000
	const maxWorkers = 35
	releaseTime := time.Now().Add(10 * time.Minute)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		uuids := make([]uuid.UUID, jobsNum)
		for i := 0; i < jobsNum; i++ {
			u, err := uuid.NewV7()
			if err != nil {
				b.Fatalf("failed to generate UUID: %v", err)
			}
			uuids[i] = u
		}

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

		failedCount := 0
		for w := 0; w < maxWorkers; w++ {
			go func(workerID int) {
				//wait until all goroutines start
				allStartedWg.Done()
				<-barrier

				for job := range jobs {
					// Charging only user 0 wallet
					_, err = repo.Debit(context.Background(), 0, &job.UUID, 1000, &releaseTime)
					if err != nil {
						failedCount++
						fmt.Printf("Failed: %v", err)
					}
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

		b.Logf("\nRun %d\nTPS: %.2f\nFailed Count: %d\ntotal time %v", n+1, float64(jobsNum)/totalTime.Seconds(), failedCount, totalTime)
	}
}

// BenchmarkDebit_multipleUsersConcurrent benchmarks the WalletRepo.Debit function for multiple users
func BenchmarkDebit_multipleUsersConcurrent(b *testing.B) {
	var err error
	repo := Init()
	defer repo.Close()

	const jobsNum = 1000
	const maxWorkers = 35
	releaseTime := time.Now().Add(10 * time.Minute)

	for i := int64(0); i < jobsNum; i++ {
		u, err := uuid.NewV7()
		if err != nil {
			panic(err)
		}
		_, err = repo.Charge(context.Background(), i, &u, 1000000, nil)
		if err != nil {
			panic(err)
		}
	}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		uuids := make([]uuid.UUID, jobsNum)
		for i := 0; i < jobsNum; i++ {
			u, err := uuid.NewV7()
			if err != nil {
				b.Fatalf("failed to generate UUID: %v", err)
			}
			uuids[i] = u
		}

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

		failedCount := 0
		for w := 0; w < maxWorkers; w++ {
			go func(workerID int) {
				//wait until all goroutines start
				allStartedWg.Done()
				<-barrier

				for job := range jobs {
					// Charging only user 0 wallet
					_, err = repo.Debit(context.Background(), job.ID, &job.UUID, 1000, &releaseTime)
					if err != nil {
						failedCount++
						fmt.Printf("Failed: %v", err)
					}
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

		b.Logf("\nRun %d\nTPS: %.2f\nFailed Count: %d\ntotal time %v", n+1, float64(jobsNum)/totalTime.Seconds(), failedCount, totalTime)
	}
}

func Init() *PgxWalletRepo {
	noopLogger := logger.NewNoopLogger()
	cfg, err := config.Load()
	if err != nil {
		fmt.Println(err)
	}
	pool, err := database.NewConnection(cfg.TestDatabase, noopLogger)
	ctx := context.Background()
	_, err = pool.Exec(ctx, "TRUNCATE wallets, transactions RESTART IDENTITY CASCADE")
	if err != nil {
		panic(err)
	}
	repo := NewPgxWalletRepo(noopLogger, pool)
	return repo
}

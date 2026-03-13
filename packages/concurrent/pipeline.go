package concurrent

import (
	"context"
	"fmt"
	"sync"
)

// =============================================================================
// Pipeline — orchestrates all stages
//
// Architecture:
//
//   [generate] ──jobCh──► [workerPool × N] ──resultCh×N──► [fanIn] ──mergedCh──► [aggregate]
//       ↑                        ↑                              ↑
//   Stage 1                  Stage 2                        Stage 3 + 4
//   (source)              (fan-out: 1→N)               (fan-in: N→1 + sink)
//
// =============================================================================

// Pipeline holds configuration for the concurrent execution pipeline
type Pipeline struct {
	WorkerCount   int // number of parallel workers (fan-out degree)
	JobBufSize    int // buffer size of the job channel
	MaxConcurrent int // max simultaneous "HTTP calls" via semaphore (rate limit)
}

// Run executes all pipeline stages and returns aggregated stats.
// The pipeline respects ctx cancellation at every stage.
func (p *Pipeline) Run(ctx context.Context, jobList []Job) Stats {
	// Stage 1: Generator — source goroutine, closes jobCh when all jobs are sent
	jobCh := generate(ctx, jobList, p.JobBufSize)

	// Shared semaphore across all workers — limits actual concurrency of I/O
	sem := newSemaphore(p.MaxConcurrent)

	// Stage 2: Worker Pool (fan-out) — N workers each get their own result channel
	resultChans := workerPool(ctx, jobCh, p.WorkerCount, sem)

	// Stage 3: Fan-in — merge N result channels into one ordered stream
	mergedCh := fanIn(ctx, resultChans...)

	// Stage 4: Aggregate — single consumer, no synchronization needed
	return aggregate(mergedCh)
}

// =============================================================================
// Stage 1: Generator
// Pattern: goroutine as data source, close channel to signal completion
// =============================================================================

func generate(ctx context.Context, jobs []Job, bufSize int) <-chan Job {
	// Buffered channel: generator can run ahead of workers up to bufSize jobs
	jobCh := make(chan Job, bufSize)

	go func() {
		defer close(jobCh) // closing signals workers that no more jobs are coming

		for _, job := range jobs {
			select {
			case jobCh <- job:
			case <-ctx.Done():
				fmt.Println("[Generator] pipeline canceled, stopping job dispatch")
				return
			}
		}
		fmt.Printf("[Generator] all %d jobs dispatched\n", len(jobs))
	}()

	return jobCh
}

// =============================================================================
// Stage 3: Fan-in
// Pattern: merge N channels into 1 using a goroutine per input channel + WaitGroup
// Each forwarder goroutine drains its input channel and forwards to merged.
// WaitGroup tracks when all forwarders finish so we know when to close merged.
// =============================================================================

func fanIn(ctx context.Context, channels ...<-chan Result) <-chan Result {
	// Buffer = number of workers so each worker can always deliver without blocking merged
	merged := make(chan Result, len(channels))
	var wg sync.WaitGroup

	// forward reads from one worker's result channel and forwards to merged
	forward := func(ch <-chan Result) {
		defer wg.Done()
		// range over channel: loops until ch is closed (worker done)
		for result := range ch {
			select {
			case merged <- result:
			case <-ctx.Done():
				// drain ch to unblock the worker so it can exit cleanly
				for range ch {
				}
				return
			}
		}
	}

	wg.Add(len(channels))
	for _, ch := range channels {
		go forward(ch) // one goroutine per input channel
	}

	// Closer goroutine: waits for all forwarders, then closes merged
	// This is the only goroutine allowed to close merged — prevents double-close
	go func() {
		wg.Wait()
		close(merged)
	}()

	return merged
}

// =============================================================================
// Stage 4: Aggregator
// Pattern: single goroutine as sole consumer — no mutex needed for Stats
// =============================================================================

func aggregate(results <-chan Result) Stats {
	var stats Stats
	for result := range results { // range exits when merged channel is closed
		fmt.Println(result)
		stats.Record(result)
	}
	return stats
}

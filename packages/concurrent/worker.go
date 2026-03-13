package concurrent

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"
)

// =============================================================================
// Semaphore — buffered channel as a counting semaphore
// Only used to rate-limit the HTTP request step.
// 10 workers run freely through all other steps; only 3 can do HTTP at once.
// =============================================================================

type semaphore chan struct{}

func newSemaphore(n int) semaphore { return make(semaphore, n) }

func (s semaphore) AcquireCtx(ctx context.Context) error {
	select {
	case s <- struct{}{}:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (s semaphore) Release() { <-s }

// =============================================================================
// Worker Pool — Fan-out
// =============================================================================

func workerPool(ctx context.Context, jobs <-chan Job, n int, sem semaphore) []<-chan Result {
	resultChans := make([]<-chan Result, n)
	for i := range n {
		resultChans[i] = startWorker(ctx, i+1, jobs, sem)
	}
	return resultChans
}

func startWorker(ctx context.Context, id int, jobs <-chan Job, sem semaphore) <-chan Result {
	out := make(chan Result)

	go func() {
		defer close(out)
		for {
			select {
			case <-ctx.Done():
				return
			case job, ok := <-jobs:
				if !ok {
					return
				}
				result := processWithRetry(ctx, id, job, sem)
				select {
				case out <- result:
				case <-ctx.Done():
					return
				}
			}
		}
	}()

	return out
}

// =============================================================================
// Retry
// =============================================================================

func processWithRetry(ctx context.Context, workerID int, job Job, sem semaphore) Result {
	var lastResult Result

	for attempt := 1; attempt <= job.MaxRetry+1; attempt++ {
		if attempt > 1 {
			backoff := time.Duration(attempt*attempt) * 80 * time.Millisecond
			fmt.Printf("[Worker-%02d] Job #%03d retry attempt=%d backoff=%v\n",
				workerID, job.ID, attempt, backoff)
			select {
			case <-time.After(backoff):
			case <-ctx.Done():
				return Result{Job: job, WorkerID: workerID, Err: ctx.Err(), Attempt: attempt}
			}
		}

		lastResult = process(ctx, workerID, job, sem)
		lastResult.Attempt = attempt

		if lastResult.Err == nil {
			return lastResult
		}
		if errors.Is(lastResult.Err, context.Canceled) || errors.Is(lastResult.Err, context.DeadlineExceeded) {
			return lastResult
		}
	}

	return lastResult
}

// =============================================================================
// process — multi-step pipeline inside each worker
//
// Step 1: Parse params      — CPU work, no I/O, no rate limit needed
// Step 2: Query local cache — fast local I/O, no rate limit needed
// Step 3: HTTP request      — ★ external I/O, rate-limited by semaphore ★
//                              (skipped on cache hit)
// Step 4: Parse response    — CPU work, no rate limit needed
// Step 5: Write to DB       — local I/O, fast enough, no rate limit needed
//
// The semaphore only wraps Step 3, so 10 workers can run Steps 1,2,4,5
// in parallel freely — only their HTTP calls are capped at 3 simultaneous.
// =============================================================================

func process(ctx context.Context, workerID int, job Job, sem semaphore) Result {
	result := Result{Job: job, WorkerID: workerID}

	// Step 1: Parse params
	t := time.Now()
	simulate(ctx, 10, 30) // fast CPU work
	result.Steps.ParseParams = time.Since(t)
	if ctx.Err() != nil {
		return Result{Job: job, WorkerID: workerID, Err: ctx.Err()}
	}

	// Step 2: Query local cache (30% hit rate)
	t = time.Now()
	cacheHit := rand.Intn(100) < 30
	simulate(ctx, 20, 50) // local cache lookup
	result.Steps.QueryCache = time.Since(t)
	if ctx.Err() != nil {
		return Result{Job: job, WorkerID: workerID, Err: ctx.Err()}
	}

	if cacheHit {
		// Cache hit: skip HTTP entirely, go straight to write DB
		result.CacheHit = true
		result.Status = 200

		// Step 4 (parse response skipped on cache hit)
		// Step 5: Write to DB
		t = time.Now()
		simulate(ctx, 30, 60)
		result.Steps.WriteDB = time.Since(t)
		return result
	}

	// Step 3: HTTP request — THE ONLY STEP BEHIND THE SEMAPHORE
	// Workers that reach here must wait if 3 others are already doing HTTP.
	// Meanwhile those workers are FREE to run Steps 1 & 2 for the next job.
	if err := sem.AcquireCtx(ctx); err != nil {
		return Result{Job: job, WorkerID: workerID, Err: err}
	}

	jobCtx, cancel := context.WithTimeout(ctx, job.Timeout)
	t = time.Now()

	httpDuration := time.Duration(rand.Intn(300)+100) * time.Millisecond
	willFail := rand.Intn(100) < 20

	var httpErr error
	select {
	case <-time.After(httpDuration):
		if willFail {
			httpErr = fmt.Errorf("connection refused")
		} else {
			result.Status = 200
		}
	case <-jobCtx.Done():
		httpErr = fmt.Errorf("http timeout: %w", jobCtx.Err())
	}

	result.Steps.HTTPRequest = time.Since(t)
	cancel()
	sem.Release() // release as soon as HTTP is done — before Steps 4 & 5

	if httpErr != nil {
		return Result{Job: job, WorkerID: workerID, Err: httpErr, Steps: result.Steps}
	}

	// Step 4: Parse response
	t = time.Now()
	simulate(ctx, 10, 30)
	result.Steps.ParseResp = time.Since(t)
	if ctx.Err() != nil {
		return Result{Job: job, WorkerID: workerID, Err: ctx.Err(), Steps: result.Steps}
	}

	// Step 5: Write to DB
	t = time.Now()
	simulate(ctx, 30, 80)
	result.Steps.WriteDB = time.Since(t)

	return result
}

// simulate blocks for a random duration in [minMs, maxMs) — represents CPU or local I/O work
func simulate(ctx context.Context, minMs, maxMs int) {
	d := time.Duration(rand.Intn(maxMs-minMs)+minMs) * time.Millisecond
	select {
	case <-time.After(d):
	case <-ctx.Done():
	}
}

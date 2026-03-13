package concurrent

import (
	"context"
	"fmt"
	"time"
)

func Concurrent() {
	fmt.Println("=== Concurrent Pipeline: Health Check Simulator ===")
	fmt.Println()

	// Simulate a microservice health check job list
	urls := []string{
		"https://api.service-a.internal/health",
		"https://api.service-b.internal/health",
		"https://api.service-c.internal/health",
		"https://db.replica-1.internal/ping",
		"https://db.replica-2.internal/ping",
		"https://cache.redis-1.internal/info",
		"https://cache.redis-2.internal/info",
		"https://queue.kafka-1.internal/status",
		"https://queue.kafka-2.internal/status",
		"https://cdn.assets.internal/health",
		"https://auth.service.internal/health",
		"https://payment.gateway.internal/ping",
		"https://notification.service.internal/health",
		"https://search.elastic.internal/health",
		"https://metrics.prometheus.internal/health",
	}

	jobs := make([]Job, len(urls))
	for i, url := range urls {
		jobs[i] = Job{
			ID:       i + 1,
			URL:      url,
			Timeout:  350 * time.Millisecond, // per-job timeout
			MaxRetry: 2,                      // retry up to 2 extra times on failure
		}
	}

	// Overall pipeline timeout — if the whole thing takes longer than this, cancel everything
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	pipeline := Pipeline{
		WorkerCount:   10, // 10 workers — freely parallel on parse/cache/db steps
		JobBufSize:    8,  // job channel buffer: generator runs up to 8 ahead of workers
		MaxConcurrent: 3,  // semaphore: max 3 simultaneous HTTP calls regardless of worker count
	}

	start := time.Now()
	stats := pipeline.Run(ctx, jobs)
	elapsed := time.Since(start)

	stats.Print()
	fmt.Printf("Total elapsed: %v\n", elapsed.Round(time.Millisecond))
}

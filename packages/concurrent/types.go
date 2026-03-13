package concurrent

import (
	"context"
	"errors"
	"fmt"
	"time"
)

type Job struct {
	ID       int
	URL      string
	Timeout  time.Duration
	MaxRetry int
}

// StepDurations records how long each pipeline step took — makes the multi-step
// nature of each worker visible in the output
type StepDurations struct {
	ParseParams  time.Duration
	QueryCache   time.Duration
	HTTPRequest  time.Duration // only this step is rate-limited by semaphore
	ParseResp    time.Duration
	WriteDB      time.Duration
}

type Result struct {
	Job       Job
	WorkerID  int
	Steps     StepDurations
	CacheHit  bool
	Status    int
	Err       error
	Attempt   int
}

func (r Result) totalDuration() time.Duration {
	s := r.Steps
	return s.ParseParams + s.QueryCache + s.HTTPRequest + s.ParseResp + s.WriteDB
}

func (r Result) String() string {
	if r.Err != nil {
		return fmt.Sprintf("[Worker-%02d] Job #%03d  %-42s  FAIL    err=%v (attempt %d)",
			r.WorkerID, r.Job.ID, r.Job.URL, r.Err, r.Attempt)
	}
	source := "http"
	if r.CacheHit {
		source = "cache"
	}
	return fmt.Sprintf(
		"[Worker-%02d] Job #%03d  %-42s  OK  status=%d  source=%-5s  "+
			"parse=%v cache=%v http=%v resp=%v db=%v  total=%v  (attempt %d)",
		r.WorkerID, r.Job.ID, r.Job.URL, r.Status, source,
		r.Steps.ParseParams.Round(time.Millisecond),
		r.Steps.QueryCache.Round(time.Millisecond),
		r.Steps.HTTPRequest.Round(time.Millisecond),
		r.Steps.ParseResp.Round(time.Millisecond),
		r.Steps.WriteDB.Round(time.Millisecond),
		r.totalDuration().Round(time.Millisecond),
		r.Attempt,
	)
}

type Stats struct {
	Total     int
	Success   int
	CacheHits int
	Failed    int
	Canceled  int
	TotalHTTPDuration time.Duration
}

func (s *Stats) Record(r Result) {
	s.Total++
	switch {
	case errors.Is(r.Err, context.Canceled), errors.Is(r.Err, context.DeadlineExceeded):
		s.Canceled++
	case r.Err != nil:
		s.Failed++
	default:
		s.Success++
		if r.CacheHit {
			s.CacheHits++
		}
		s.TotalHTTPDuration += r.Steps.HTTPRequest
	}
}

func (s Stats) Print() {
	avgHTTP := time.Duration(0)
	httpJobs := s.Success - s.CacheHits
	if httpJobs > 0 {
		avgHTTP = s.TotalHTTPDuration / time.Duration(httpJobs)
	}
	fmt.Printf("\n=== Pipeline Stats ===\n")
	fmt.Printf("Total: %d | Success: %d (cache hits: %d) | Failed: %d | Canceled: %d | Avg HTTP: %v\n",
		s.Total, s.Success, s.CacheHits, s.Failed, s.Canceled, avgHTTP.Round(time.Millisecond))
}

package mockdata

import (
	"context"
	"fmt"
	"time"

	"go-playground/internal/microdemo/config"
	microdemo "go-playground/kitex_gen/microdemo"
)

func SleepContext(ctx context.Context, d time.Duration) error {
	if d <= 0 {
		return nil
	}
	timer := time.NewTimer(d)
	defer timer.Stop()

	select {
	case <-timer.C:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func ProfileDelay(scenario string) time.Duration {
	switch scenario {
	case config.ScenarioProfileSlow, config.ScenarioBothSlow, config.ScenarioProfileSlowRecFail:
		return 450 * time.Millisecond
	default:
		return 80 * time.Millisecond
	}
}

func RecommendDelay(scenario string) time.Duration {
	switch scenario {
	case config.ScenarioRecommendSlow, config.ScenarioBothSlow:
		return 420 * time.Millisecond
	default:
		return 120 * time.Millisecond
	}
}

func RecommendShouldFail(scenario string) bool {
	return scenario == config.ScenarioRecommendFail || scenario == config.ScenarioProfileSlowRecFail
}

func Profile(userID int64) *microdemo.ProfileResponse {
	cities := []string{"Beijing", "Shanghai", "Hangzhou", "Shenzhen"}
	return &microdemo.ProfileResponse{
		UserId: userID,
		Name:   fmt.Sprintf("user-%d", userID),
		City:   cities[int(userID)%len(cities)],
		Bio:    "Distributed systems learner who likes tracing RPC latency.",
	}
}

func Recommendations(userID int64, limit int32) *microdemo.RecommendResponse {
	if limit <= 0 {
		limit = 3
	}
	reasons := []string{
		"Because your profile says you like Go microservices.",
		"Popular among similar users in your city.",
		"Fresh content with low latency from the feed cache.",
		"Fallback ranking used after one downstream degraded.",
	}
	items := make([]*microdemo.RecommendItem, 0, limit)
	for i := int32(0); i < limit; i++ {
		items = append(items, &microdemo.RecommendItem{
			PostId: userID*100 + int64(i) + 1,
			Title:  fmt.Sprintf("Recommended post %d for user %d", i+1, userID),
			Reason: reasons[int(i)%len(reasons)],
		})
	}
	return &microdemo.RecommendResponse{
		UserId: userID,
		Items:  items,
	}
}

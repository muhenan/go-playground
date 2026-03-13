package recommend

import (
	"context"
	"fmt"
	"log"

	"go-playground/internal/microdemo/mockdata"
	microdemo "go-playground/kitex_gen/microdemo"
)

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s *Service) GetRecommendations(ctx context.Context, req *microdemo.RecommendRequest) (*microdemo.RecommendResponse, error) {
	scenario := req.GetScenario()
	delay := mockdata.RecommendDelay(scenario)
	log.Printf("[recommend] user=%d scenario=%q delay=%v limit=%d", req.GetUserId(), scenario, delay, req.GetLimit())

	if err := mockdata.SleepContext(ctx, delay); err != nil {
		return nil, err
	}
	if mockdata.RecommendShouldFail(scenario) {
		return nil, fmt.Errorf("recommendation engine unavailable for scenario %q", scenario)
	}

	return mockdata.Recommendations(req.GetUserId(), req.GetLimit()), nil
}

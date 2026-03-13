package profile

import (
	"context"
	"log"

	"go-playground/internal/microdemo/mockdata"
	microdemo "go-playground/kitex_gen/microdemo"
)

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s *Service) GetProfile(ctx context.Context, req *microdemo.ProfileRequest) (*microdemo.ProfileResponse, error) {
	scenario := req.GetScenario()
	delay := mockdata.ProfileDelay(scenario)
	log.Printf("[profile] user=%d scenario=%q delay=%v", req.GetUserId(), scenario, delay)

	if err := mockdata.SleepContext(ctx, delay); err != nil {
		return nil, err
	}

	return mockdata.Profile(req.GetUserId()), nil
}

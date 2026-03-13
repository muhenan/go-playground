package config

import "time"

const (
	ProfileServiceName   = "ProfileService"
	RecommendServiceName = "RecommendService"

	ProfileServiceAddr   = "127.0.0.1:9001"
	RecommendServiceAddr = "127.0.0.1:9002"
	GatewayHTTPAddr      = "127.0.0.1:8081"

	DefaultGatewayTimeout = 350 * time.Millisecond
)

const (
	ScenarioNormal             = ""
	ScenarioProfileSlow        = "profile_slow"
	ScenarioRecommendSlow      = "recommend_slow"
	ScenarioRecommendFail      = "recommend_fail"
	ScenarioBothSlow           = "both_slow"
	ScenarioProfileSlowRecFail = "profile_slow_recommend_fail"
)

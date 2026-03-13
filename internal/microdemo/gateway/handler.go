package gateway

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go-playground/internal/microdemo/config"
	microdemo "go-playground/kitex_gen/microdemo"
	"go-playground/kitex_gen/microdemo/profileservice"
	"go-playground/kitex_gen/microdemo/recommendservice"
)

type Handler struct {
	profileClient   profileservice.Client
	recommendClient recommendservice.Client
}

type FeedResponse struct {
	UserID          int64                      `json:"userId"`
	Scenario        string                     `json:"scenario,omitempty"`
	TimeoutMS       int64                      `json:"timeoutMs"`
	ElapsedMS       int64                      `json:"elapsedMs"`
	Degraded        bool                       `json:"degraded"`
	TimedOut        bool                       `json:"timedOut"`
	Profile         *microdemo.ProfileResponse `json:"profile,omitempty"`
	Recommendations []*microdemo.RecommendItem `json:"recommendations,omitempty"`
	Errors          []string                   `json:"errors,omitempty"`
}

type profileResult struct {
	resp *microdemo.ProfileResponse
	err  error
}

type recommendResult struct {
	resp *microdemo.RecommendResponse
	err  error
}

func NewHandler(profileClient profileservice.Client, recommendClient recommendservice.Client) *Handler {
	return &Handler{
		profileClient:   profileClient,
		recommendClient: recommendClient,
	}
}

func (h *Handler) RegisterRoutes(r *gin.Engine) {
	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})
	r.GET("/feed/:userID", h.GetFeed)
}

func (h *Handler) GetFeed(c *gin.Context) {
	userID, err := strconv.ParseInt(c.Param("userID"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	limit := int32(3)
	if raw := c.Query("limit"); raw != "" {
		v, parseErr := strconv.ParseInt(raw, 10, 32)
		if parseErr != nil || v <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "limit must be a positive integer"})
			return
		}
		limit = int32(v)
	}

	timeout := config.DefaultGatewayTimeout
	if raw := c.Query("timeout_ms"); raw != "" {
		v, parseErr := strconv.Atoi(raw)
		if parseErr != nil || v <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "timeout_ms must be a positive integer"})
			return
		}
		timeout = time.Duration(v) * time.Millisecond
	}

	scenario := c.Query("scenario")
	ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
	defer cancel()

	start := time.Now()
	resp := FeedResponse{
		UserID:    userID,
		Scenario:  scenario,
		TimeoutMS: timeout.Milliseconds(),
	}

	profileCh := make(chan profileResult, 1)
	recommendCh := make(chan recommendResult, 1)

	go h.fetchProfile(ctx, userID, scenario, profileCh)
	go h.fetchRecommendations(ctx, userID, limit, scenario, recommendCh)

	var (
		pendingProfile   <-chan profileResult   = profileCh
		pendingRecommend <-chan recommendResult = recommendCh
	)

	for pendingProfile != nil || pendingRecommend != nil {
		select {
		case result := <-pendingProfile:
			pendingProfile = nil
			if result.err != nil {
				resp.Degraded = true
				resp.Errors = append(resp.Errors, fmt.Sprintf("profile rpc failed: %v", result.err))
				continue
			}
			resp.Profile = result.resp
		case result := <-pendingRecommend:
			pendingRecommend = nil
			if result.err != nil {
				resp.Degraded = true
				resp.Errors = append(resp.Errors, fmt.Sprintf("recommend rpc failed: %v", result.err))
				continue
			}
			resp.Recommendations = result.resp.Items
		case <-ctx.Done():
			resp.TimedOut = true
			resp.Degraded = true
			if pendingProfile != nil {
				resp.Errors = append(resp.Errors, "profile rpc timed out or was canceled")
				pendingProfile = nil
			}
			if pendingRecommend != nil {
				resp.Errors = append(resp.Errors, "recommend rpc timed out or was canceled")
				pendingRecommend = nil
			}
		}
	}

	resp.ElapsedMS = time.Since(start).Milliseconds()

	status := http.StatusOK
	if resp.TimedOut && resp.Profile == nil && len(resp.Recommendations) == 0 {
		status = http.StatusGatewayTimeout
	}
	c.JSON(status, resp)
}

func (h *Handler) fetchProfile(ctx context.Context, userID int64, scenario string, out chan<- profileResult) {
	result := profileResult{
		resp: nil,
		err:  nil,
	}
	resp, err := h.profileClient.GetProfile(ctx, &microdemo.ProfileRequest{
		UserId:   userID,
		Scenario: optionalString(scenario),
	})
	result.resp = resp
	result.err = err

	select {
	case out <- result:
	case <-ctx.Done():
	}
}

func (h *Handler) fetchRecommendations(ctx context.Context, userID int64, limit int32, scenario string, out chan<- recommendResult) {
	result := recommendResult{
		resp: nil,
		err:  nil,
	}

	var lastErr error
	for attempt := 1; attempt <= 2; attempt++ {
		resp, err := h.recommendClient.GetRecommendations(ctx, &microdemo.RecommendRequest{
			UserId:   userID,
			Limit:    limit,
			Scenario: optionalString(scenario),
		})
		if err == nil {
			result.resp = resp
			lastErr = nil
			break
		}
		lastErr = err
		if ctx.Err() != nil {
			break
		}
		log.Printf("[gateway] recommend attempt=%d failed: %v", attempt, err)
	}
	result.err = lastErr

	select {
	case out <- result:
	case <-ctx.Done():
	}
}

func optionalString(v string) *string {
	if v == "" {
		return nil
	}
	return &v
}

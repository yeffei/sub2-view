package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/pkg/pagination"
	"github.com/Wei-Shaw/sub2api/internal/pkg/usagestats"
	middleware2 "github.com/Wei-Shaw/sub2api/internal/server/middleware"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

type userUsageRepoCapture struct {
	service.UsageLogRepository
	listParams  pagination.PaginationParams
	listFilters usagestats.UsageLogFilters
	listRows    []service.UsageLog
	stats       *usagestats.UsageStats
	statsUserID int64
	modelStats  []usagestats.ModelStat
}

func (s *userUsageRepoCapture) ListWithFilters(ctx context.Context, params pagination.PaginationParams, filters usagestats.UsageLogFilters) ([]service.UsageLog, *pagination.PaginationResult, error) {
	s.listParams = params
	s.listFilters = filters
	return s.listRows, &pagination.PaginationResult{
		Total:    int64(len(s.listRows)),
		Page:     params.Page,
		PageSize: params.PageSize,
		Pages:    0,
	}, nil
}

func (s *userUsageRepoCapture) GetUserStatsAggregated(_ context.Context, userID int64, _, _ time.Time) (*usagestats.UsageStats, error) {
	s.statsUserID = userID
	if s.stats != nil {
		return s.stats, nil
	}
	return &usagestats.UsageStats{}, nil
}

func (s *userUsageRepoCapture) GetUserModelStats(_ context.Context, _ int64, _, _ time.Time) ([]usagestats.ModelStat, error) {
	return s.modelStats, nil
}

func newUserUsageRequestTypeTestRouter(repo *userUsageRepoCapture) *gin.Engine {
	gin.SetMode(gin.TestMode)
	usageSvc := service.NewUsageService(repo, nil, nil, nil)
	handler := NewUsageHandler(usageSvc, nil, nil, nil)
	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set(string(middleware2.ContextKeyUser), middleware2.AuthSubject{UserID: 42})
		c.Next()
	})
	router.GET("/usage", handler.List)
	router.GET("/usage/stats", handler.Stats)
	router.GET("/usage/dashboard/models", handler.DashboardModels)
	return router
}

func TestUserUsageListRequestTypePriority(t *testing.T) {
	repo := &userUsageRepoCapture{}
	router := newUserUsageRequestTypeTestRouter(repo)

	req := httptest.NewRequest(http.MethodGet, "/usage?request_type=ws_v2&stream=bad", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	require.Equal(t, int64(42), repo.listFilters.UserID)
	require.NotNil(t, repo.listFilters.RequestType)
	require.Equal(t, int16(service.RequestTypeWSV2), *repo.listFilters.RequestType)
	require.Nil(t, repo.listFilters.Stream)
}

func TestUserUsageListInvalidRequestType(t *testing.T) {
	repo := &userUsageRepoCapture{}
	router := newUserUsageRequestTypeTestRouter(repo)

	req := httptest.NewRequest(http.MethodGet, "/usage?request_type=invalid", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestUserUsageListInvalidStream(t *testing.T) {
	repo := &userUsageRepoCapture{}
	router := newUserUsageRequestTypeTestRouter(repo)

	req := httptest.NewRequest(http.MethodGet, "/usage?stream=invalid", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestUserUsageListSupportedFilters(t *testing.T) {
	repo := &userUsageRepoCapture{}
	router := newUserUsageRequestTypeTestRouter(repo)

	req := httptest.NewRequest(http.MethodGet, "/usage?model=gpt-5&billing_type=1&request_type=sync&start_date=2026-03-01&end_date=2026-03-02", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	require.Equal(t, int64(42), repo.listFilters.UserID)
	require.Equal(t, "gpt-5", repo.listFilters.Model)
	require.Empty(t, repo.listFilters.ModelFilterSource)
	require.NotNil(t, repo.listFilters.BillingType)
	require.Equal(t, int8(1), *repo.listFilters.BillingType)
	require.NotNil(t, repo.listFilters.RequestType)
	require.Equal(t, int16(service.RequestTypeSync), *repo.listFilters.RequestType)
	require.NotNil(t, repo.listFilters.StartTime)
	require.NotNil(t, repo.listFilters.EndTime)
}

func TestUserUsageListKeepsBillingWithoutAdminFields(t *testing.T) {
	ipAddress := "203.0.113.10"
	upstreamModel := "upstream-private-model"
	billingTier := "internal-tier"
	channelID := int64(99)
	accountRateMultiplier := 1.7
	accountStatsCost := 0.12
	repo := &userUsageRepoCapture{
		listRows: []service.UsageLog{{
			ID:                    1,
			UserID:                42,
			APIKeyID:              7,
			AccountID:             5,
			RequestID:             "req_user_billing",
			Model:                 "gpt-5",
			InputCost:             0.01,
			OutputCost:            0.02,
			CacheCreationCost:     0.03,
			CacheReadCost:         0.04,
			TotalCost:             0.10,
			ActualCost:            0.08,
			RateMultiplier:        0.8,
			IPAddress:             &ipAddress,
			UpstreamModel:         &upstreamModel,
			BillingTier:           &billingTier,
			ChannelID:             &channelID,
			AccountRateMultiplier: &accountRateMultiplier,
			AccountStatsCost:      &accountStatsCost,
		}},
	}
	router := newUserUsageRequestTypeTestRouter(repo)

	req := httptest.NewRequest(http.MethodGet, "/usage", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	body := rec.Body.String()
	require.Contains(t, body, `"input_cost":0.01`)
	require.Contains(t, body, `"output_cost":0.02`)
	require.Contains(t, body, `"cache_creation_cost":0.03`)
	require.Contains(t, body, `"cache_read_cost":0.04`)
	require.Contains(t, body, `"total_cost":0.1`)
	require.Contains(t, body, `"actual_cost":0.08`)
	require.Contains(t, body, `"rate_multiplier":0.8`)
	require.NotContains(t, body, "ip_address")
	require.NotContains(t, body, "upstream_endpoint")
	require.NotContains(t, body, "account_rate_multiplier")
	require.NotContains(t, body, "account_stats_cost")
	require.NotContains(t, body, "upstream_model")
	require.NotContains(t, body, "billing_tier")
	require.NotContains(t, body, "channel_id")
	require.NotContains(t, body, `"account":`)
}

func TestUserUsageStatsUsesScopedFilters(t *testing.T) {
	accountCost := 0.12
	repo := &userUsageRepoCapture{
		stats: &usagestats.UsageStats{
			TotalCost:        0.10,
			TotalActualCost:  0.08,
			TotalAccountCost: &accountCost,
			UpstreamEndpoints: []usagestats.EndpointStat{{
				Endpoint: "/v1/responses",
			}},
			EndpointPaths: []usagestats.EndpointStat{{
				Endpoint: "/v1/chat/completions -> /v1/responses",
			}},
		},
	}
	router := newUserUsageRequestTypeTestRouter(repo)

	req := httptest.NewRequest(http.MethodGet, "/usage/stats?group_id=9&request_type=sync&billing_mode=token&start_date=2026-03-01&end_date=2026-03-02", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	require.Equal(t, int64(42), repo.statsUserID)
	require.Contains(t, rec.Body.String(), `"total_cost":0.1`)
	require.Contains(t, rec.Body.String(), `"total_actual_cost":0.08`)
	require.NotContains(t, rec.Body.String(), "total_account_cost")
	require.NotContains(t, rec.Body.String(), "upstream_endpoints")
	require.NotContains(t, rec.Body.String(), "endpoint_paths")
}

func TestUserUsageDashboardModelsOmitsAccountCost(t *testing.T) {
	repo := &userUsageRepoCapture{
		modelStats: []usagestats.ModelStat{{
			Model:       "gpt-5",
			Requests:    2,
			TotalTokens: 30,
			Cost:        0.10,
			ActualCost:  0.08,
			AccountCost: 0.07,
		}},
	}
	router := newUserUsageRequestTypeTestRouter(repo)

	req := httptest.NewRequest(http.MethodGet, "/usage/dashboard/models?start_date=2026-03-01&end_date=2026-03-02", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	body := rec.Body.String()
	require.Contains(t, body, `"cost":0.1`)
	require.Contains(t, body, `"actual_cost":0.08`)
	require.NotContains(t, body, "account_cost")
}

func TestUserUsageDashboardModelsDoesNotExposeAdminModelSources(t *testing.T) {
	repo := &userUsageRepoCapture{}
	router := newUserUsageRequestTypeTestRouter(repo)

	req := httptest.NewRequest(http.MethodGet, "/usage/dashboard/models?model_source=upstream&start_date=2026-03-01&end_date=2026-03-02", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	require.NotContains(t, rec.Body.String(), "account_cost")
}

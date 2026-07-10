//go:build integration

package repository

import (
	"context"
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/pkg/usagestats"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/stretchr/testify/require"
)

func TestUsageLog_GetStatsWithFilters_AggregatesAndEndpoints(t *testing.T) {
	ctx := context.Background()
	tx := testEntTx(t)
	client := tx.Client()
	repo := newUsageLogRepositoryWithSQL(client, tx)

	user := mustCreateUser(t, client, &service.User{Email: "stats@test.com"})
	apiKey := mustCreateApiKey(t, client, &service.APIKey{UserID: user.ID, Key: "sk-stats-1", Name: "k"})
	account := mustCreateAccount(t, client, &service.Account{Name: "acc-stats"})

	now := time.Now().UTC()
	inboundEndpoint := "/v1/messages"
	upstreamEndpoint := "/v1/responses"
	logs := []*service.UsageLog{
		{
			UserID: user.ID, APIKeyID: apiKey.ID, AccountID: account.ID,
			Model: "claude-3", InputTokens: 2, OutputTokens: 3,
			CacheReadTokens: 4, TotalCost: 0.5, ActualCost: 0.4, CreatedAt: now,
			InboundEndpoint: &inboundEndpoint, UpstreamEndpoint: &upstreamEndpoint,
		},
		{
			UserID: user.ID, APIKeyID: apiKey.ID, AccountID: account.ID,
			Model: "claude-3", InputTokens: 4, OutputTokens: 5,
			CacheCreationTokens: 7, TotalCost: 0.6, ActualCost: 0.5, CreatedAt: now,
			InboundEndpoint: &inboundEndpoint, UpstreamEndpoint: &upstreamEndpoint,
		},
		{
			UserID: user.ID, APIKeyID: apiKey.ID, AccountID: account.ID,
			Model: "claude-3", InputTokens: 6, OutputTokens: 7,
			TotalCost: 0.7, ActualCost: 0.6, CreatedAt: now,
			InboundEndpoint: &inboundEndpoint, UpstreamEndpoint: &upstreamEndpoint,
		},
	}
	for _, log := range logs {
		_, err := repo.Create(ctx, log)
		require.NoError(t, err)
	}

	start := now.Add(-1 * time.Hour)
	end := now.Add(1 * time.Hour)
	// 按本测试创建的 user 维度过滤:集成库为共享实例,其它用 testEntClient 的兄弟测试会留下
	// 已提交的 usage_log 行(含零 token 的失败请求),不限定 user 会把它们计入 TotalRequests。
	stats, err := repo.GetStatsWithFilters(ctx, usagestats.UsageLogFilters{UserID: user.ID, StartTime: &start, EndTime: &end})
	require.NoError(t, err)
	require.Equal(t, int64(3), stats.TotalRequests)
	require.Equal(t, int64(12), stats.TotalInputTokens)
	require.Equal(t, int64(15), stats.TotalOutputTokens)
	require.Equal(t, int64(7), stats.TotalCacheCreationTokens)
	require.Equal(t, int64(4), stats.TotalCacheReadTokens)
	require.Equal(t, int64(1), stats.CacheReadHitRequests)
	require.Equal(t, int64(1), stats.CacheCreationRequests)
	require.InDelta(t, 1.5, stats.TotalActualCost, 1e-9)
	require.InDelta(t, 1.0/3.0, stats.CacheReadHitRatio, 1e-9)
	require.InDelta(t, 4.0, stats.AverageCacheReadTokensPerHit, 1e-9)
	require.InDelta(t, 4.0, stats.AverageActualInputTokens, 1e-9)
	require.NotEmpty(t, stats.Endpoints)
	require.NotEmpty(t, stats.UpstreamEndpoints)
	require.NotEmpty(t, stats.EndpointPaths)
	require.Equal(t, int64(12), stats.Endpoints[0].InputTokens)
	require.Equal(t, int64(15), stats.Endpoints[0].OutputTokens)
	require.Equal(t, int64(7), stats.Endpoints[0].CacheCreationTokens)
	require.Equal(t, int64(4), stats.Endpoints[0].CacheReadTokens)
	require.InDelta(t, 4.0, stats.Endpoints[0].AverageActualInputTokens, 1e-9)
}

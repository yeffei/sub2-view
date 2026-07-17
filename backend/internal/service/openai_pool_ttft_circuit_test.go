package service

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type highTTFTAccountRepoStub struct {
	stubOpenAIAccountRepo
	accountIDs []int64
	untils     []time.Time
	reasons    []string
}

func (r *highTTFTAccountRepoStub) SetTempUnschedulable(_ context.Context, accountID int64, until time.Time, reason string) error {
	r.accountIDs = append(r.accountIDs, accountID)
	r.untils = append(r.untils, until)
	r.reasons = append(r.reasons, reason)
	return nil
}

func TestRateLimitService_OpenAIPoolModeHighTTFT_TracksConsecutiveSevereSamples(t *testing.T) {
	svc := &RateLimitService{}
	for i := 1; i <= openAIPoolModeHighTTFTConsecutive; i++ {
		require.Equal(t, i, svc.recordOpenAIPoolModeHighTTFT(7001, OpenAIPoolModeHighTTFTThresholdMs))
	}
}

func TestRateLimitService_OpenAIPoolModeHighTTFT_LowSampleResetsStreak(t *testing.T) {
	svc := &RateLimitService{}
	require.Equal(t, 1, svc.recordOpenAIPoolModeHighTTFT(7002, OpenAIPoolModeHighTTFTThresholdMs))
	require.Equal(t, 2, svc.recordOpenAIPoolModeHighTTFT(7002, OpenAIPoolModeHighTTFTThresholdMs))
	require.Zero(t, svc.recordOpenAIPoolModeHighTTFT(7002, OpenAIPoolModeHighTTFTThresholdMs-1))
	require.Equal(t, 1, svc.recordOpenAIPoolModeHighTTFT(7002, OpenAIPoolModeHighTTFTThresholdMs))
}

func TestRateLimitService_OpenAIPoolModeHighTTFT_OpensCircuitAfterConsecutiveSamples(t *testing.T) {
	repo := &highTTFTAccountRepoStub{}
	svc := &RateLimitService{accountRepo: repo}
	account := &Account{
		ID:       7003,
		Platform: PlatformOpenAI,
		Type:     AccountTypeAPIKey,
		Credentials: map[string]any{
			"pool_mode": true,
		},
	}
	observed := OpenAIPoolModeHighTTFTThresholdMs

	require.False(t, svc.HandleOpenAIPoolModeHighTTFT(context.Background(), account, &observed))
	require.False(t, svc.HandleOpenAIPoolModeHighTTFT(context.Background(), account, &observed))
	require.True(t, svc.HandleOpenAIPoolModeHighTTFT(context.Background(), account, &observed))

	require.Equal(t, []int64{account.ID}, repo.accountIDs)
	require.Len(t, repo.untils, 1)
	require.True(t, repo.untils[0].After(time.Now()))
	require.Contains(t, repo.reasons[0], "pool_mode_high_ttft")
	require.Zero(t, svc.recordOpenAIPoolModeHighTTFT(account.ID, observed-1))
}

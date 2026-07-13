package service

import (
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/stretchr/testify/require"
)

func TestGatewayServiceStickyHealthUsesRealTTFT(t *testing.T) {
	ttftThreshold := 6000
	svc := &GatewayService{
		cfg: &config.Config{Gateway: config.GatewayConfig{
			OpenAIScheduler: config.GatewayOpenAISchedulerConfig{
				StickyEscapeEnabled:   true,
				StickyEscapeTTFTMs:    ttftThreshold,
				StickyEscapeErrorRate: 0.5,
			},
		}},
		accountRuntimeStats: newOpenAIAccountRuntimeStats(),
	}
	account := &Account{ID: 10, Status: "active", Schedulable: true}

	require.False(t, svc.shouldClearStickySession(account, "gpt-4o"))
	probeLatency := 100
	svc.ReportAccountScheduleResult(account.ID, true, &probeLatency)
	require.False(t, svc.shouldClearStickySession(account, "gpt-4o"))

	realTTFT := 12000
	for i := 0; i < 5; i++ {
		svc.ReportAccountScheduleResult(account.ID, true, &realTTFT)
	}
	require.True(t, svc.shouldClearStickySession(account, "gpt-4o"))
}

func TestGatewayServiceStickyHealthReportsErrors(t *testing.T) {
	svc := &GatewayService{
		cfg: &config.Config{Gateway: config.GatewayConfig{
			OpenAIScheduler: config.GatewayOpenAISchedulerConfig{
				StickyEscapeEnabled:   true,
				StickyEscapeTTFTMs:    15000,
				StickyEscapeErrorRate: 0.5,
			},
		}},
		accountRuntimeStats: newOpenAIAccountRuntimeStats(),
	}
	account := &Account{ID: 11, Status: "active", Schedulable: true}

	for i := 0; i < 5; i++ {
		svc.ReportAccountScheduleResult(account.ID, false, nil)
	}
	require.True(t, svc.shouldClearStickySession(account, "gpt-4o"))
}

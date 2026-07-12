//go:build unit

package service

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/pkg/openai_compat"
	"github.com/stretchr/testify/require"
)

type poolSnapshotPruneRepoStub struct {
	AccountMonitorRepository
	cutoffs []time.Time
	deleted int64
}

func (s *poolSnapshotPruneRepoStub) DeletePoolAvailabilityBefore(_ context.Context, before time.Time) (int64, error) {
	s.cutoffs = append(s.cutoffs, before)
	return s.deleted, nil
}

func TestPoolAvailabilitySnapshotPruneRunsDailyWithThirtyDayRetention(t *testing.T) {
	now := time.Date(2026, 7, 12, 9, 0, 0, 0, time.UTC)
	repo := &poolSnapshotPruneRepoStub{deleted: 12}
	svc := NewAccountMonitorService(repo, nil, nil)
	svc.now = func() time.Time { return now }

	deleted, err := svc.prunePoolAvailabilitySnapshots(t.Context())
	require.NoError(t, err)
	require.Equal(t, int64(12), deleted)
	require.Equal(t, []time.Time{now.Add(-30 * 24 * time.Hour)}, repo.cutoffs)

	deleted, err = svc.prunePoolAvailabilitySnapshots(t.Context())
	require.NoError(t, err)
	require.Zero(t, deleted)
	require.Len(t, repo.cutoffs, 1)

	now = now.Add(24 * time.Hour)
	_, err = svc.prunePoolAvailabilitySnapshots(t.Context())
	require.NoError(t, err)
	require.Len(t, repo.cutoffs, 2)
}

func TestPoolRuntimeWeightRequiresThreeObservationsAndRecovers(t *testing.T) {
	now := time.Date(2026, 7, 12, 10, 0, 0, 0, time.UTC)
	state := (*PoolRuntimeWeightState)(nil)
	for i := 0; i < 2; i++ {
		state = nextPoolRuntimeWeightState(1, 2, state, 0.5, "probe_failed", now.Add(time.Duration(i)*time.Minute))
		require.Equal(t, 1.0, state.Factor)
	}
	state = nextPoolRuntimeWeightState(1, 2, state, 0.5, "probe_failed", now.Add(2*time.Minute))
	require.Equal(t, 0.75, state.Factor)

	for i := 0; i < 2; i++ {
		state = nextPoolRuntimeWeightState(1, 2, state, 1.0, "healthy", now.Add(time.Duration(3+i)*time.Minute))
		require.Equal(t, 0.75, state.Factor)
	}
	state = nextPoolRuntimeWeightState(1, 2, state, 1.0, "healthy", now.Add(5*time.Minute))
	require.Equal(t, 1.0, state.Factor)
}

func TestUpstreamHealthAlertRequiresThreeObservationsRemindsAndResolves(t *testing.T) {
	now := time.Date(2026, 7, 12, 10, 0, 0, 0, time.UTC)
	var events []upstreamHealthAlertEvent
	svc := NewAccountMonitorService(nil, nil, nil)
	svc.now = func() time.Time { return now }
	svc.healthAlertEmitter = func(event upstreamHealthAlertEvent) {
		events = append(events, event)
	}
	observation := upstreamHealthAlertObservation{
		key: "account_probe_failed:pool:1:account:2", alertType: "account_probe_failed",
		poolID: 1, accountID: 2, message: "上游账号连续探测失败",
	}

	svc.reconcileUpstreamHealthAlerts([]upstreamHealthAlertObservation{observation})
	svc.reconcileUpstreamHealthAlerts([]upstreamHealthAlertObservation{observation})
	require.Empty(t, events)

	svc.reconcileUpstreamHealthAlerts([]upstreamHealthAlertObservation{observation})
	require.Len(t, events, 1)
	require.Equal(t, "firing", events[0].status)

	now = now.Add(upstreamHealthAlertReminderCooldown - time.Minute)
	svc.reconcileUpstreamHealthAlerts([]upstreamHealthAlertObservation{observation})
	require.Len(t, events, 1)

	now = now.Add(time.Minute)
	svc.reconcileUpstreamHealthAlerts([]upstreamHealthAlertObservation{observation})
	require.Len(t, events, 2)
	require.Equal(t, "reminder", events[1].status)

	svc.reconcileUpstreamHealthAlerts(nil)
	require.Len(t, events, 3)
	require.Equal(t, "resolved", events[2].status)
	require.Empty(t, svc.healthAlertStates)
}

func TestPoolHealthAlertCapacityDeduplicatesExpandedMembers(t *testing.T) {
	accounts := map[int64]*Account{
		1: {ID: 1, Status: StatusActive, Schedulable: true},
		2: {ID: 2, Status: StatusActive, Schedulable: false},
	}
	members := []UpstreamPoolMember{
		{AccountID: 1, Enabled: true},
		{AccountID: 1, Enabled: true, SourceType: "set"},
		{AccountID: 2, Enabled: true},
	}

	available, configured := poolHealthAlertCapacity(members, accounts)
	require.Equal(t, 1, available)
	require.Equal(t, 2, configured)
}

func TestAutoWeightTargetUsesRuntimeStateAndPoolRelativeLatency(t *testing.T) {
	now := time.Date(2026, 7, 12, 10, 0, 0, 0, time.UTC)
	resetAt := now.Add(5 * time.Minute)
	target, reason, observed := autoWeightTarget(&Account{RateLimitResetAt: &resetAt}, nil, 0, now)
	require.True(t, observed)
	require.Equal(t, 0.25, target)
	require.Equal(t, "rate_limited", reason)

	fastLatency := 70
	target, reason, observed = autoWeightTarget(&Account{}, &AccountMonitorHistoryRow{
		Status: MonitorStatusOperational, LatencyMs: &fastLatency,
	}, 100, now)
	require.True(t, observed)
	require.Equal(t, 1.25, target)
	require.Equal(t, "faster_than_pool", reason)

	slowLatency := 180
	target, reason, observed = autoWeightTarget(&Account{}, &AccountMonitorHistoryRow{
		Status: MonitorStatusOperational, LatencyMs: &slowLatency,
	}, 100, now)
	require.True(t, observed)
	require.Equal(t, 0.5, target)
	require.Equal(t, "much_slower_than_pool", reason)
}

func TestUpstreamPoolAutoWeightPolicyPreservesRoutingPolicy(t *testing.T) {
	policy := SetUpstreamPoolAccountTypeStrategyPolicyJSON(nil, UpstreamPoolAccountTypeStrategyOAuthPreferred)
	policy = SetUpstreamPoolAutoWeightPolicyJSON(policy, true)
	require.True(t, UpstreamPoolAutoWeightEnabledFromPolicyJSON(policy))
	require.Equal(t, UpstreamPoolAccountTypeStrategyOAuthPreferred, UpstreamPoolAccountTypeStrategyFromPolicyJSON(policy))
}

func TestNormalizeUpstreamPoolRejectsAutoWeightOutsideOpenAI(t *testing.T) {
	err := normalizeUpstreamPoolForCreate(&UpstreamPool{Name: "anthropic-pool", Code: "anthropic-pool", Platform: PlatformAnthropic, AutoWeightEnabled: true})
	require.ErrorContains(t, err, "auto weight currently supports OpenAI pools only")
}

func TestAccountProbeAPIModeUsesResponsesWhenSupported(t *testing.T) {
	account := &Account{
		Platform: PlatformOpenAI,
		Type:     AccountTypeAPIKey,
		Extra: map[string]any{
			openai_compat.ExtraKeyResponsesSupported: true,
		},
	}

	require.Equal(t, MonitorAPIModeResponses, accountProbeAPIMode(account))
}

func TestRunAccountProbeTarget_ResponsesUsesOnlyHiAndOneToken(t *testing.T) {
	h := &openAICaptureHandler{}
	endpoint := setupFakeOpenAI(t, h)

	row := runAccountProbeTarget(t.Context(), accountProbeTarget{
		AccountID: 1,
		PoolID:    2,
		Provider:  MonitorProviderOpenAI,
		Model:     "gpt-5.4",
		Endpoint:  endpoint,
		APIKey:    "sk-test",
		APIMode:   MonitorAPIModeResponses,
	})

	require.Equal(t, MonitorStatusOperational, row.Status)
	require.Equal(t, providerOpenAIResponsesPath, h.lastPath)
	require.Equal(t, monitorLightweightPrompt, strings.TrimSpace(stringFromAny(h.lastBody["input"])))
	require.Equal(t, float64(monitorLightweightMaxTokens), h.lastBody["max_output_tokens"])
	require.NotContains(t, h.lastBody, "instructions")
	require.Equal(t, "Bearer sk-test", h.lastHeaders.Get("Authorization"))
}

func TestRunAccountProbeTarget_ChatUsesOnlyHiAndOneToken(t *testing.T) {
	h := &openAICaptureHandler{}
	endpoint := setupFakeOpenAI(t, h)

	row := runAccountProbeTarget(t.Context(), accountProbeTarget{
		AccountID: 1,
		PoolID:    2,
		Provider:  MonitorProviderOpenAI,
		Model:     "gpt-5.4",
		Endpoint:  endpoint,
		APIKey:    "sk-test",
		APIMode:   MonitorAPIModeChatCompletions,
	})

	require.Equal(t, MonitorStatusOperational, row.Status)
	require.Equal(t, providerOpenAIPath, h.lastPath)
	messages, _ := h.lastBody["messages"].([]any)
	require.Len(t, messages, 1)
	msg, _ := messages[0].(map[string]any)
	require.Equal(t, monitorLightweightPrompt, strings.TrimSpace(stringFromAny(msg["content"])))
	require.Equal(t, float64(monitorLightweightMaxTokens), h.lastBody["max_tokens"])
	require.Equal(t, "Bearer sk-test", h.lastHeaders.Get("Authorization"))
}

func TestRunAccountProbeTarget_AnthropicBearerAuthScheme(t *testing.T) {
	h := &captureHandler{respondText: monitorLightweightPrompt}
	endpoint := setupFakeAnthropic(t, h)

	row := runAccountProbeTarget(t.Context(), accountProbeTarget{
		AccountID:  1,
		PoolID:     2,
		Provider:   MonitorProviderAnthropic,
		Model:      "claude-opus-4-8",
		Endpoint:   endpoint,
		APIKey:     "sk-test",
		AuthScheme: AnthropicAPIKeyAuthSchemeAuthorizationBearer,
	})

	require.Equal(t, MonitorStatusOperational, row.Status)
	messages, _ := h.lastBody["messages"].([]any)
	require.Len(t, messages, 1)
	msg, _ := messages[0].(map[string]any)
	require.Equal(t, monitorLightweightPrompt, strings.TrimSpace(stringFromAny(msg["content"])))
	require.Equal(t, float64(monitorLightweightMaxTokens), h.lastBody["max_tokens"])
	require.Equal(t, "Bearer sk-test", h.lastHeaders.Get("Authorization"))
	require.Empty(t, h.lastHeaders.Get("x-api-key"))
	require.Equal(t, monitorAnthropicAPIVersion, h.lastHeaders.Get("anthropic-version"))
}

func TestRunTestBackgroundWithMode_LightweightResponsesUsesMinimalViableProbe(t *testing.T) {
	h := &openAICaptureHandler{}
	endpoint := setupFakeOpenAI(t, h)
	account := &Account{
		ID:       10,
		Platform: PlatformOpenAI,
		Type:     AccountTypeAPIKey,
		Credentials: map[string]any{
			"api_key":  "sk-test",
			"base_url": endpoint,
		},
		Extra: map[string]any{
			openai_compat.ExtraKeyResponsesSupported: true,
		},
	}
	svc := &AccountTestService{accountRepo: &stubOpenAIAccountRepo{accounts: []Account{*account}}}

	result, err := svc.RunTestBackgroundWithMode(t.Context(), account.ID, "gpt-5.4", AccountTestModeLightweight)

	require.NoError(t, err)
	require.Equal(t, "success", result.Status)
	require.Equal(t, providerOpenAIResponsesPath, h.lastPath)
	require.Equal(t, monitorLightweightPrompt, strings.TrimSpace(stringFromAny(h.lastBody["input"])))
	require.Equal(t, float64(monitorLightweightMaxTokens), h.lastBody["max_output_tokens"])
	require.NotContains(t, h.lastBody, "instructions")
}

func TestRunTestBackgroundWithMode_LightweightChatUsesMinimalViableProbe(t *testing.T) {
	h := &openAICaptureHandler{}
	endpoint := setupFakeOpenAI(t, h)
	account := &Account{
		ID:       11,
		Platform: PlatformOpenAI,
		Type:     AccountTypeAPIKey,
		Credentials: map[string]any{
			"api_key":  "sk-test",
			"base_url": endpoint,
		},
		Extra: map[string]any{
			openai_compat.ExtraKeyResponsesSupported: false,
		},
	}
	svc := &AccountTestService{accountRepo: &stubOpenAIAccountRepo{accounts: []Account{*account}}}

	result, err := svc.RunTestBackgroundWithMode(t.Context(), account.ID, "gpt-5.4", AccountTestModeLightweight)

	require.NoError(t, err)
	require.Equal(t, "success", result.Status)
	require.Equal(t, providerOpenAIPath, h.lastPath)
	messages, _ := h.lastBody["messages"].([]any)
	require.Len(t, messages, 1)
	msg, _ := messages[0].(map[string]any)
	require.Equal(t, monitorLightweightPrompt, strings.TrimSpace(stringFromAny(msg["content"])))
	require.Equal(t, float64(monitorLightweightMaxTokens), h.lastBody["max_tokens"])
	require.NotContains(t, h.lastBody, "instructions")
}

func TestAccountProbeAPIModeUsesChatWhenResponsesUnsupported(t *testing.T) {
	account := &Account{
		Platform: PlatformOpenAI,
		Type:     AccountTypeAPIKey,
		Extra: map[string]any{
			openai_compat.ExtraKeyResponsesSupported: false,
		},
	}

	require.Equal(t, MonitorAPIModeChatCompletions, accountProbeAPIMode(account))
}

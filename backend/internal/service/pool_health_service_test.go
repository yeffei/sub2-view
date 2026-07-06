package service

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestPoolHealthAggregatesAccountProbeHistory(t *testing.T) {
	groupID := int64(7)
	latency := 123
	pingLatency := 45
	now := time.Now().UTC()

	pool := UpstreamPool{
		ID:       1,
		Name:     "OpenAI Pool",
		Platform: PlatformOpenAI,
		Enabled:  true,
	}
	bindings := []UpstreamPoolBinding{
		{ID: 1, PoolID: pool.ID, GroupID: groupID, GroupName: "openai-main", Platform: PlatformOpenAI, Models: []string{"gpt-4o"}, Enabled: true},
	}
	members := []UpstreamPoolMember{
		{PoolID: pool.ID, AccountID: 1001, AccountName: "acc-main", AccountPlatform: PlatformOpenAI, Enabled: true, Weight: 1},
		{PoolID: pool.ID, AccountID: 1002, AccountName: "acc-backup", AccountPlatform: PlatformOpenAI, Enabled: true, Weight: 1},
	}
	accountsByID := map[int64]*Account{
		1001: {ID: 1001, Name: "acc-main", Platform: PlatformOpenAI, Status: StatusActive, Schedulable: true},
		1002: {ID: 1002, Name: "acc-backup", Platform: PlatformOpenAI, Status: StatusActive, Schedulable: true},
	}
	facts := &poolAccountFacts{
		latestByAccount: map[int64][]*AccountMonitorLatest{
			1001: {{AccountID: 1001, Model: "gpt-4o", Status: MonitorStatusOperational, LatencyMs: &latency, PingLatencyMs: &pingLatency, CheckedAt: now}},
			1002: {{AccountID: 1002, Model: "gpt-4o", Status: MonitorStatusFailed, CheckedAt: now}},
		},
		avail7ByAccount: map[int64][]*AccountMonitorAvailability{
			1001: {{AccountID: 1001, Model: "gpt-4o", WindowDays: monitorAvailability7Days, AvailabilityPct: 100}},
			1002: {{AccountID: 1002, Model: "gpt-4o", WindowDays: monitorAvailability7Days, AvailabilityPct: 0}},
		},
		avail15ByAccount: map[int64][]*AccountMonitorAvailability{},
		avail30ByAccount: map[int64][]*AccountMonitorAvailability{},
		historyByAccount: map[int64][]*AccountMonitorHistoryEntry{
			1001: {{AccountID: 1001, Model: "gpt-4o", Status: MonitorStatusOperational, LatencyMs: &latency, PingLatencyMs: &pingLatency, CheckedAt: now}},
			1002: {{AccountID: 1002, Model: "gpt-4o", Status: MonitorStatusFailed, CheckedAt: now}},
		},
	}

	detail := buildPoolHealthDetail(pool, bindings, members, accountsByID, facts)

	require.Equal(t, MonitorStatusOperational, detail.Status)
	require.Equal(t, 2, detail.HealthyMemberCount)
	require.Equal(t, 2, detail.TotalMemberCount)
	require.Equal(t, latency, *detail.BestLatencyMs)
	require.Equal(t, pingLatency, *detail.BestPingLatencyMs)
	require.Len(t, detail.Lines, 2)
	require.Equal(t, int64(1001), detail.Lines[0].AccountID)
	require.Equal(t, "gpt-4o", detail.Lines[0].ProbeModel)
	require.Equal(t, 50.0, detail.Availability7d)
}

func TestPoolHealthUsesAccountStateWithoutMonitorLines(t *testing.T) {
	pool := UpstreamPool{ID: 1, Name: "OpenAI Pool", Platform: PlatformOpenAI, Enabled: true}
	bindings := []UpstreamPoolBinding{
		{ID: 1, PoolID: pool.ID, GroupID: 7, GroupName: "openai-main", Platform: PlatformOpenAI, Enabled: true},
	}
	members := []UpstreamPoolMember{{PoolID: pool.ID, AccountID: 1001, Enabled: true}}
	accountsByID := map[int64]*Account{
		1001: {ID: 1001, Platform: PlatformOpenAI, Status: StatusActive, Schedulable: true},
	}
	facts := emptyPoolAccountFacts()

	detail := buildPoolHealthDetail(pool, bindings, members, accountsByID, facts)

	require.Equal(t, MonitorStatusOperational, detail.Status)
	require.Equal(t, 100.0, detail.Availability7d)
	require.Equal(t, 1, detail.HealthyMemberCount)
	require.Equal(t, 1, detail.TotalMemberCount)
	require.Len(t, detail.Lines, 1)
}

func TestPoolHealthServingStatusIgnoresPartialUnschedulableMembers(t *testing.T) {
	pool := UpstreamPool{ID: 1, Name: "OpenAI Pool", Platform: PlatformOpenAI, Enabled: true}
	bindings := []UpstreamPoolBinding{
		{ID: 1, PoolID: pool.ID, GroupID: 7, GroupName: "openai-main", Platform: PlatformOpenAI, Enabled: true},
	}
	members := []UpstreamPoolMember{
		{PoolID: pool.ID, AccountID: 1001, Enabled: true},
		{PoolID: pool.ID, AccountID: 1002, Enabled: true},
	}
	accountsByID := map[int64]*Account{
		1001: {ID: 1001, Platform: PlatformOpenAI, Status: StatusActive, Schedulable: true},
		1002: {ID: 1002, Platform: PlatformOpenAI, Status: StatusActive, Schedulable: false},
	}

	detail := buildPoolHealthDetail(pool, bindings, members, accountsByID, emptyPoolAccountFacts())

	require.Equal(t, MonitorStatusOperational, detail.Status)
	require.Equal(t, 1, detail.HealthyMemberCount)
	require.Equal(t, 1, detail.DegradedMemberCount)
	require.Equal(t, 50.0, detail.Availability7d)
}

func TestPoolHealthUsesFailedProbeOverHealthyMemberState(t *testing.T) {
	now := time.Now().UTC()
	pool := UpstreamPool{ID: 1, Name: "OpenAI Pool", Platform: PlatformOpenAI, Enabled: true}
	bindings := []UpstreamPoolBinding{
		{ID: 1, PoolID: pool.ID, GroupID: 7, GroupName: "openai-main", Platform: PlatformOpenAI, Models: []string{"gpt-4o"}, Enabled: true},
	}
	members := []UpstreamPoolMember{
		{PoolID: pool.ID, AccountID: 1001, Enabled: true},
		{PoolID: pool.ID, AccountID: 1002, Enabled: true},
	}
	accountsByID := map[int64]*Account{
		1001: {ID: 1001, Platform: PlatformOpenAI, Status: StatusActive, Schedulable: true},
		1002: {ID: 1002, Platform: PlatformOpenAI, Status: StatusActive, Schedulable: true},
	}
	facts := emptyPoolAccountFacts()
	facts.latestByAccount[1001] = []*AccountMonitorLatest{{AccountID: 1001, Model: "gpt-4o", Status: MonitorStatusFailed, CheckedAt: now}}
	facts.latestByAccount[1002] = []*AccountMonitorLatest{{AccountID: 1002, Model: "gpt-4o", Status: MonitorStatusFailed, CheckedAt: now}}

	detail := buildPoolHealthDetail(pool, bindings, members, accountsByID, facts)

	require.Equal(t, MonitorStatusFailed, detail.Status)
	require.Equal(t, 2, detail.HealthyMemberCount)
	require.Equal(t, 2, detail.TotalMemberCount)
}

func TestPoolHealthRequiresPoolBinding(t *testing.T) {
	pool := UpstreamPool{ID: 1, Name: "OpenAI Pool", Platform: PlatformOpenAI, Enabled: true}
	members := []UpstreamPoolMember{{PoolID: pool.ID, AccountID: 1001, Enabled: true}}
	accountsByID := map[int64]*Account{
		1001: {ID: 1001, Platform: PlatformOpenAI, Status: StatusActive, Schedulable: true},
	}
	facts := emptyPoolAccountFacts()

	detail := buildPoolHealthDetail(pool, nil, members, accountsByID, facts)

	require.Equal(t, MonitorStatusError, detail.Status)
	require.Empty(t, detail.GroupName)
}

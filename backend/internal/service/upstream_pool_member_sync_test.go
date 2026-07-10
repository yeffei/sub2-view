package service

import "testing"

func TestBuildUpstreamPoolMemberSyncResult_MembershipOnlyPreservesManualFields(t *testing.T) {
	priority := 7
	concurrency := 3
	override := false
	pool := &UpstreamPool{ID: 11, Platform: PlatformOpenAI}
	current := []UpstreamPoolMember{{
		ID:                     21,
		PoolID:                 11,
		AccountID:              101,
		AccountName:            "manual",
		AccountPlatform:        PlatformOpenAI,
		Enabled:                false,
		ManualDrained:          true,
		Weight:                 55,
		SchedulableOverride:    &override,
		PriorityOverride:       &priority,
		MaxConcurrencyOverride: &concurrency,
		Notes:                  "keep me",
		SourceType:             "direct",
		Editable:               true,
	}}
	targets := []UpstreamPoolMember{{
		PoolID:          11,
		AccountID:       101,
		AccountName:     "manual",
		AccountPlatform: PlatformOpenAI,
		Enabled:         true,
		ManualDrained:   false,
		Weight:          100,
		Notes:           "synced",
	}}

	result := BuildUpstreamPoolMemberSyncResult(pool, current, targets, UpstreamPoolMemberSyncModeMembershipOnly)
	if result.UpdateCount != 0 {
		t.Fatalf("UpdateCount = %d, want 0", result.UpdateCount)
	}
	if result.SkipCount != 1 {
		t.Fatalf("SkipCount = %d, want 1", result.SkipCount)
	}
	if result.OverwriteRiskCount != 0 {
		t.Fatalf("OverwriteRiskCount = %d, want 0", result.OverwriteRiskCount)
	}
}

func TestBuildUpstreamPoolMemberSyncResult_OverwriteReportsManualFields(t *testing.T) {
	priority := 7
	pool := &UpstreamPool{ID: 11, Platform: PlatformOpenAI}
	current := []UpstreamPoolMember{{
		ID:               21,
		PoolID:           11,
		AccountID:        101,
		AccountName:      "manual",
		AccountPlatform:  PlatformOpenAI,
		Enabled:          false,
		ManualDrained:    true,
		Weight:           55,
		PriorityOverride: &priority,
		Notes:            "keep me",
		SourceType:       "direct",
		Editable:         true,
	}}
	targets := []UpstreamPoolMember{{
		PoolID:          11,
		AccountID:       101,
		AccountName:     "manual",
		AccountPlatform: PlatformOpenAI,
		Enabled:         true,
		ManualDrained:   false,
		Weight:          100,
		Notes:           "synced",
	}}

	result := BuildUpstreamPoolMemberSyncResult(pool, current, targets, UpstreamPoolMemberSyncModeOverwriteSchedulerFields)
	if result.UpdateCount != 1 {
		t.Fatalf("UpdateCount = %d, want 1", result.UpdateCount)
	}
	if result.OverwriteRiskCount != 1 {
		t.Fatalf("OverwriteRiskCount = %d, want 1", result.OverwriteRiskCount)
	}
	if len(result.Updates) != 1 || len(result.Updates[0].Overwrites) == 0 {
		t.Fatalf("Updates = %+v, want overwrite details", result.Updates)
	}
}

func TestEnabledUpstreamPoolBindingGroupIDsFromList(t *testing.T) {
	got := enabledUpstreamPoolBindingGroupIDsFromList(10, []UpstreamPoolBinding{
		{ID: 1, PoolID: 10, GroupID: 102, Enabled: true},
		{ID: 2, PoolID: 10, GroupID: 101, Enabled: true},
		{ID: 3, PoolID: 10, GroupID: 102, Enabled: true},
		{ID: 4, PoolID: 10, GroupID: 103, Enabled: false},
		{ID: 5, PoolID: 11, GroupID: 104, Enabled: true},
	})
	want := []int64{101, 102}
	if len(got) != len(want) {
		t.Fatalf("len = %d, want %d: %+v", len(got), len(want), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("got = %+v, want %+v", got, want)
		}
	}
}

func TestAccountMatchesUpstreamPoolMemberSyncStrategy(t *testing.T) {
	if !accountMatchesUpstreamPoolMemberSyncStrategy(&Account{Type: AccountTypeOAuth}, UpstreamPoolAccountTypeStrategyOAuthOnly) {
		t.Fatal("oauth account should match oauth_only")
	}
	if !accountMatchesUpstreamPoolMemberSyncStrategy(&Account{Type: AccountTypeSetupToken}, UpstreamPoolAccountTypeStrategyOAuthOnly) {
		t.Fatal("setup-token account should match oauth_only")
	}
	if accountMatchesUpstreamPoolMemberSyncStrategy(&Account{Type: AccountTypeAPIKey}, UpstreamPoolAccountTypeStrategyOAuthOnly) {
		t.Fatal("apikey account should not match oauth_only")
	}
	if !accountMatchesUpstreamPoolMemberSyncStrategy(&Account{Type: AccountTypeAPIKey}, UpstreamPoolAccountTypeStrategyOAuthPreferred) {
		t.Fatal("preferred strategies should keep fallback account types eligible for sync")
	}
}

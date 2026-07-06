package service

import (
	"context"
	"errors"
	"testing"
)

func TestLogRoutingExplanationWritesDirectSinkEvent(t *testing.T) {
	logSink, restore := captureStructuredLog(t)
	defer restore()

	groupID := int64(2)
	stickyEnabled := true
	stickyTTFT := 6000
	stickyErrorRate := 0.3
	observedTTFT := 7200

	logRoutingExplanation(context.Background(), &groupID, "gpt-5.5", "session-hash-123456", RoutingExplanation{
		Engine:                     "openai",
		Layer:                      "load_balance",
		Reason:                     "load_balance",
		SelectedAccountID:          9,
		SelectedAccountType:        "apikey",
		CandidateCount:             2,
		TopK:                       2,
		LoadSkew:                   1,
		LatencyMs:                  12,
		CacheAffinityKeyHash:       "abc123",
		CacheAffinityTopKIDs:       []int64{9, 10},
		RequiredCapability:         "chat_completions",
		Skipped:                    map[string]int{"runtime_blocked": 1},
		PoolID:                     1,
		PoolCode:                   "1",
		PoolName:                   "Codex Plus",
		StickyEscapeSource:         "pool",
		StickyEscapeEnabled:        &stickyEnabled,
		StickyEscapeTTFTMs:         &stickyTTFT,
		StickyEscapeErrorRate:      &stickyErrorRate,
		StickyEscapeTriggered:      true,
		StickyEscapeReason:         "ttft",
		StickyEscapeObservedTTFTMs: &observedTTFT,
	}, nil)

	logSink.mu.Lock()
	defer logSink.mu.Unlock()
	if len(logSink.events) != 1 {
		t.Fatalf("events len = %d, want 1", len(logSink.events))
	}

	event := logSink.events[0]
	if event.Level != "info" {
		t.Fatalf("level = %q, want info", event.Level)
	}
	if event.Component != "routing.explanation" {
		t.Fatalf("component = %q, want routing.explanation", event.Component)
	}
	if event.Message != "routing_explanation" {
		t.Fatalf("message = %q, want routing_explanation", event.Message)
	}
	if got := event.Fields["platform"]; got != "openai" {
		t.Fatalf("platform field = %v, want openai", got)
	}
	if got := event.Fields["pool_id"]; got != int64(1) {
		t.Fatalf("pool_id field = %v, want 1", got)
	}
	if got := event.Fields["account_id"]; got != int64(9) {
		t.Fatalf("account_id field = %v, want 9", got)
	}
	if got := event.Fields["sticky_escape_triggered"]; got != true {
		t.Fatalf("sticky_escape_triggered field = %v, want true", got)
	}
	if got := event.Fields["cache_affinity_key_hash"]; got != "abc123" {
		t.Fatalf("cache_affinity_key_hash field = %v, want abc123", got)
	}
	if got, ok := event.Fields["cache_affinity_top_k_account_ids"].([]int64); !ok || len(got) != 2 || got[0] != 9 || got[1] != 10 {
		t.Fatalf("cache_affinity_top_k_account_ids field = %v, want [9 10]", event.Fields["cache_affinity_top_k_account_ids"])
	}
	if got, ok := event.Fields["skipped"].(map[string]int); !ok || got["runtime_blocked"] != 1 {
		t.Fatalf("skipped field = %v, want runtime_blocked=1", event.Fields["skipped"])
	}
}

func TestLogRoutingExplanationErrorUsesWarnLevel(t *testing.T) {
	logSink, restore := captureStructuredLog(t)
	defer restore()

	logRoutingExplanation(context.Background(), nil, "gpt-5.5", "", RoutingExplanation{
		Engine: "openai",
		Reason: "fallback_load_balance",
	}, errors.New("no available accounts"))

	logSink.mu.Lock()
	defer logSink.mu.Unlock()
	if len(logSink.events) != 1 {
		t.Fatalf("events len = %d, want 1", len(logSink.events))
	}
	event := logSink.events[0]
	if event.Level != "warn" {
		t.Fatalf("level = %q, want warn", event.Level)
	}
	if got := event.Fields["error"]; got != "no available accounts" {
		t.Fatalf("error field = %v, want no available accounts", got)
	}
}

func TestNewPoolRoutingExplanationIncludesPoolFields(t *testing.T) {
	explanation := newPoolRoutingExplanation("anthropic", "selection", "load_balance", &Account{
		ID:   12,
		Type: AccountTypeAPIKey,
	}, 3, &UpstreamPoolResolvedBinding{
		Pool: &UpstreamPool{
			ID:   7,
			Code: "claude-code",
			Name: "Claude Code",
		},
	})

	if explanation.Engine != "anthropic" {
		t.Fatalf("engine = %q, want anthropic", explanation.Engine)
	}
	if explanation.SelectedAccountID != 12 {
		t.Fatalf("selected account = %d, want 12", explanation.SelectedAccountID)
	}
	if explanation.CandidateCount != 3 {
		t.Fatalf("candidate count = %d, want 3", explanation.CandidateCount)
	}
	if explanation.PoolID != 7 || explanation.PoolCode != "claude-code" || explanation.PoolName != "Claude Code" {
		t.Fatalf("pool fields = (%d, %q, %q), want (7, claude-code, Claude Code)", explanation.PoolID, explanation.PoolCode, explanation.PoolName)
	}
}

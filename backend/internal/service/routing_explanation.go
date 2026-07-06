package service

import (
	"context"
	"strings"

	"github.com/Wei-Shaw/sub2api/internal/pkg/logger"
)

type RoutingExplanation struct {
	Engine                        string         `json:"engine,omitempty"`
	Layer                         string         `json:"layer,omitempty"`
	Reason                        string         `json:"reason,omitempty"`
	SelectedAccountID             int64          `json:"selected_account_id,omitempty"`
	SelectedAccountType           string         `json:"selected_account_type,omitempty"`
	CandidateCount                int            `json:"candidate_count,omitempty"`
	TopK                          int            `json:"top_k,omitempty"`
	LoadSkew                      float64        `json:"load_skew,omitempty"`
	LatencyMs                     int64          `json:"latency_ms,omitempty"`
	CacheAffinityKeyHash          string         `json:"cache_affinity_key_hash,omitempty"`
	CacheAffinityTopKIDs          []int64        `json:"cache_affinity_top_k_account_ids,omitempty"`
	WaitPlanned                   bool           `json:"wait_planned,omitempty"`
	RequiredTransport             string         `json:"required_transport,omitempty"`
	RequiredCapability            string         `json:"required_capability,omitempty"`
	RequiredImageCapability       string         `json:"required_image_capability,omitempty"`
	RequireCompact                bool           `json:"require_compact,omitempty"`
	ExcludedCount                 int            `json:"excluded_count,omitempty"`
	Fallback                      bool           `json:"fallback,omitempty"`
	PoolID                        int64          `json:"pool_id,omitempty"`
	PoolCode                      string         `json:"pool_code,omitempty"`
	PoolName                      string         `json:"pool_name,omitempty"`
	StickyEscapeSource            string         `json:"sticky_escape_source,omitempty"`
	StickyEscapeEnabled           *bool          `json:"sticky_escape_enabled,omitempty"`
	StickyEscapeTTFTMs            *int           `json:"sticky_escape_ttft_ms,omitempty"`
	StickyEscapeErrorRate         *float64       `json:"sticky_escape_error_rate,omitempty"`
	StickyEscapeTriggered         bool           `json:"sticky_escape_triggered,omitempty"`
	StickyEscapeReason            string         `json:"sticky_escape_reason,omitempty"`
	StickyEscapeObservedTTFTMs    *int           `json:"sticky_escape_observed_ttft_ms,omitempty"`
	StickyEscapeObservedErrorRate *float64       `json:"sticky_escape_observed_error_rate,omitempty"`
	Skipped                       map[string]int `json:"skipped,omitempty"`
}

func WithRoutingExplanation(ctx context.Context, value RoutingExplanation, bridgeOldKeys bool) context.Context {
	return updateRequestMetadata(ctx, bridgeOldKeys, func(md *RequestMetadata) {
		v := value
		md.RoutingExplanation = &v
	}, nil)
}

func RoutingExplanationFromContext(ctx context.Context) (RoutingExplanation, bool) {
	if md := metadataFromContext(ctx); md != nil && md.RoutingExplanation != nil {
		return *md.RoutingExplanation, true
	}
	return RoutingExplanation{}, false
}

func routingReasonFromOpenAIDecision(decision OpenAIAccountScheduleDecision, waitPlanned bool, fallback bool) string {
	if fallback {
		if waitPlanned {
			return "fallback_load_balance_wait"
		}
		return "fallback_load_balance"
	}
	switch decision.Layer {
	case openAIAccountScheduleLayerPreviousResponse:
		return "previous_response_sticky"
	case openAIAccountScheduleLayerSessionSticky:
		if waitPlanned {
			return "session_sticky_wait"
		}
		return "session_sticky"
	case openAIAccountScheduleLayerLoadBalance:
		if waitPlanned {
			return "load_balance_wait"
		}
		return "load_balance"
	default:
		if waitPlanned {
			return "selection_wait"
		}
		return "selection"
	}
}

func newOpenAIRoutingExplanation(
	decision OpenAIAccountScheduleDecision,
	selection *AccountSelectionResult,
	excludedIDs map[int64]struct{},
	requiredTransport OpenAIUpstreamTransport,
	requiredCapability OpenAIEndpointCapability,
	requiredImageCapability OpenAIImagesCapability,
	requireCompact bool,
	policy OpenAIRoutingPolicy,
	stickyEscapeConfig openAIStickyEscapeConfig,
	stickyEscapeSource string,
	fallback bool,
) RoutingExplanation {
	waitPlanned := selection != nil && selection.WaitPlan != nil
	stickyEscapeEnabled := stickyEscapeConfig.enabled
	stickyEscapeTTFT := int(stickyEscapeConfig.ttftMs)
	stickyEscapeErrorRate := stickyEscapeConfig.errorRate
	explanation := RoutingExplanation{
		Engine:                        "openai",
		Layer:                         strings.TrimSpace(decision.Layer),
		Reason:                        routingReasonFromOpenAIDecision(decision, waitPlanned, fallback),
		SelectedAccountID:             decision.SelectedAccountID,
		SelectedAccountType:           strings.TrimSpace(decision.SelectedAccountType),
		CandidateCount:                decision.CandidateCount,
		TopK:                          decision.TopK,
		LoadSkew:                      decision.LoadSkew,
		LatencyMs:                     decision.LatencyMs,
		CacheAffinityKeyHash:          strings.TrimSpace(decision.CacheAffinityKeyHash),
		CacheAffinityTopKIDs:          append([]int64(nil), decision.CacheAffinityTopKIDs...),
		WaitPlanned:                   waitPlanned,
		RequiredTransport:             string(requiredTransport),
		RequiredCapability:            string(requiredCapability),
		RequiredImageCapability:       string(requiredImageCapability),
		RequireCompact:                requireCompact,
		ExcludedCount:                 len(excludedIDs),
		Fallback:                      fallback,
		PoolID:                        policy.PoolID,
		PoolCode:                      strings.TrimSpace(policy.PoolCode),
		PoolName:                      strings.TrimSpace(policy.PoolName),
		StickyEscapeSource:            strings.TrimSpace(stickyEscapeSource),
		StickyEscapeEnabled:           &stickyEscapeEnabled,
		StickyEscapeTTFTMs:            &stickyEscapeTTFT,
		StickyEscapeErrorRate:         &stickyEscapeErrorRate,
		StickyEscapeTriggered:         decision.StickyEscapeTriggered,
		StickyEscapeReason:            strings.TrimSpace(decision.StickyEscapeReason),
		StickyEscapeObservedTTFTMs:    decision.StickyEscapeObservedTTFTMs,
		StickyEscapeObservedErrorRate: decision.StickyEscapeObservedErrorRate,
		Skipped:                       cloneStringIntMap(decision.Skipped),
	}
	if selection != nil && selection.Account != nil {
		explanation.SelectedAccountID = selection.Account.ID
		if explanation.SelectedAccountType == "" {
			explanation.SelectedAccountType = selection.Account.Type
		}
	}
	return explanation
}

func newPoolRoutingExplanation(engine, layer, reason string, account *Account, candidateCount int, resolved *UpstreamPoolResolvedBinding) RoutingExplanation {
	explanation := RoutingExplanation{
		Engine:         strings.TrimSpace(engine),
		Layer:          strings.TrimSpace(layer),
		Reason:         strings.TrimSpace(reason),
		CandidateCount: candidateCount,
	}
	if account != nil {
		explanation.SelectedAccountID = account.ID
		explanation.SelectedAccountType = strings.TrimSpace(account.Type)
	}
	if resolved != nil && resolved.Pool != nil {
		explanation.PoolID = resolved.Pool.ID
		explanation.PoolCode = strings.TrimSpace(resolved.Pool.Code)
		explanation.PoolName = strings.TrimSpace(resolved.Pool.Name)
	}
	return explanation
}

func cloneStringIntMap(input map[string]int) map[string]int {
	if len(input) == 0 {
		return nil
	}
	out := make(map[string]int, len(input))
	for key, value := range input {
		if strings.TrimSpace(key) == "" || value <= 0 {
			continue
		}
		out[key] = value
	}
	return out
}

func logRoutingExplanation(ctx context.Context, groupID *int64, requestedModel string, sessionHash string, explanation RoutingExplanation, err error) {
	_ = ctx
	fields := map[string]any{
		"component":                 "routing.explanation",
		"engine":                    explanation.Engine,
		"platform":                  explanation.Engine,
		"layer":                     explanation.Layer,
		"reason":                    explanation.Reason,
		"group_id":                  derefGroupID(groupID),
		"model":                     requestedModel,
		"session":                   shortSessionHash(sessionHash),
		"account_id":                explanation.SelectedAccountID,
		"selected_account_type":     explanation.SelectedAccountType,
		"candidate_count":           explanation.CandidateCount,
		"top_k":                     explanation.TopK,
		"load_skew":                 explanation.LoadSkew,
		"routing_latency_ms":        explanation.LatencyMs,
		"wait_planned":              explanation.WaitPlanned,
		"required_transport":        explanation.RequiredTransport,
		"required_capability":       explanation.RequiredCapability,
		"required_image_capability": explanation.RequiredImageCapability,
		"require_compact":           explanation.RequireCompact,
		"excluded_count":            explanation.ExcludedCount,
		"fallback":                  explanation.Fallback,
	}
	if explanation.PoolID > 0 {
		fields["pool_id"] = explanation.PoolID
	}
	if explanation.PoolCode != "" {
		fields["pool_code"] = explanation.PoolCode
	}
	if explanation.PoolName != "" {
		fields["pool_name"] = explanation.PoolName
	}
	if explanation.StickyEscapeSource != "" {
		fields["sticky_escape_source"] = explanation.StickyEscapeSource
	}
	if explanation.CacheAffinityKeyHash != "" {
		fields["cache_affinity_key_hash"] = explanation.CacheAffinityKeyHash
	}
	if len(explanation.CacheAffinityTopKIDs) > 0 {
		fields["cache_affinity_top_k_account_ids"] = explanation.CacheAffinityTopKIDs
	}
	if explanation.StickyEscapeEnabled != nil {
		fields["sticky_escape_enabled"] = *explanation.StickyEscapeEnabled
	}
	if explanation.StickyEscapeTTFTMs != nil {
		fields["sticky_escape_ttft_ms"] = *explanation.StickyEscapeTTFTMs
	}
	if explanation.StickyEscapeErrorRate != nil {
		fields["sticky_escape_error_rate"] = *explanation.StickyEscapeErrorRate
	}
	if explanation.StickyEscapeTriggered {
		fields["sticky_escape_triggered"] = true
	}
	if explanation.StickyEscapeReason != "" {
		fields["sticky_escape_reason"] = explanation.StickyEscapeReason
	}
	if explanation.StickyEscapeObservedTTFTMs != nil {
		fields["sticky_escape_observed_ttft_ms"] = *explanation.StickyEscapeObservedTTFTMs
	}
	if explanation.StickyEscapeObservedErrorRate != nil {
		fields["sticky_escape_observed_error_rate"] = *explanation.StickyEscapeObservedErrorRate
	}
	if len(explanation.Skipped) > 0 {
		fields["skipped"] = explanation.Skipped
	}
	level := "info"
	if err != nil {
		fields["error"] = err.Error()
		level = "warn"
	}
	logger.WriteSinkEvent(level, "routing.explanation", "routing_explanation", fields)
}

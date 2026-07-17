package service

import (
	"math/rand/v2"
	"strings"

	"github.com/Wei-Shaw/sub2api/internal/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
)

const (
	cacheInstrumentationComponent = "cache.instrumentation"

	cacheFamilyOpenAIResponses             = "openai_responses"
	cacheFamilyOpenAIChatToResponses       = "openai_chat_to_responses"
	cacheFamilyAnthropicToOpenAIResponses  = "anthropic_to_openai_responses"
	cacheInstrumentationStateGinContextKey = "__cache_instrumentation_state"
)

type CacheInstrumentationSnapshot struct {
	CacheFamily                 string
	SessionHashShort            string
	SessionSignalSource         string
	UpstreamSessionAnchorSource string
	PromptCacheKeyPresent       bool
	PromptCacheKeyAutoInjected  bool
	PreviousResponseIDPresent   bool
	StickyAccountHit            bool
	SelectedAccountID           int64
	StickyAccountID             int64
	AccountSwitchHappened       bool
}

type cacheInstrumentationState struct {
	snapshot CacheInstrumentationSnapshot
}

type cacheInstrumentationEmitInput struct {
	Snapshot            *CacheInstrumentationSnapshot
	RequestID           string
	UserID              int64
	AccountID           int64
	Platform            string
	RequestedModel      string
	UpstreamModel       string
	InboundEndpoint     string
	UpstreamEndpoint    string
	RequestType         string
	CacheReadTokens     int
	CacheCreationTokens int
}

var cacheInstrumentationSampleFloat64 = func() float64 {
	return rand.Float64()
}

func CaptureCacheInstrumentationSnapshot(c *gin.Context) *CacheInstrumentationSnapshot {
	state := getCacheInstrumentationState(c)
	if state == nil {
		return nil
	}
	snapshot := state.snapshot
	if strings.TrimSpace(snapshot.CacheFamily) == "" {
		return nil
	}
	return &snapshot
}

func RecordCacheInstrumentationRouting(c *gin.Context, sessionHash string, decision OpenAIAccountScheduleDecision) {
	state := getOrCreateCacheInstrumentationState(c)
	if state == nil {
		return
	}
	if session := shortSessionHash(strings.TrimSpace(sessionHash)); session != "" {
		state.snapshot.SessionHashShort = session
	}
	if decision.StickyPreviousHit || decision.StickySessionHit {
		state.snapshot.StickyAccountHit = true
	}
	if decision.SelectedAccountID > 0 {
		state.snapshot.SelectedAccountID = decision.SelectedAccountID
	}
	if decision.StickyAccountID > 0 {
		state.snapshot.StickyAccountID = decision.StickyAccountID
	}
	if decision.StickyAccountID > 0 && decision.SelectedAccountID > 0 && decision.StickyAccountID != decision.SelectedAccountID {
		state.snapshot.AccountSwitchHappened = true
	}
}

func MarkCacheInstrumentationAccountSwitch(c *gin.Context) {
	state := getOrCreateCacheInstrumentationState(c)
	if state == nil {
		return
	}
	state.snapshot.AccountSwitchHappened = true
}

func RecordCacheInstrumentationUpstreamSessionAnchor(c *gin.Context, source string) {
	state := getOrCreateCacheInstrumentationState(c)
	if state == nil {
		return
	}
	state.snapshot.UpstreamSessionAnchorSource = strings.TrimSpace(source)
}

func setCacheInstrumentationOpenAIResponsesState(c *gin.Context, body []byte, promptCacheKey string, promptCacheKeyAutoInjected bool, previousResponseIDPresent bool) {
	setOpenAIResponsesAutoUpstreamSessionAnchor(c, body, previousResponseIDPresent)
	state := getOrCreateCacheInstrumentationState(c)
	if state == nil {
		return
	}
	state.snapshot.CacheFamily = cacheFamilyOpenAIResponses
	state.snapshot.SessionSignalSource = resolveOpenAISessionSignalSource(c, body)
	state.snapshot.PromptCacheKeyPresent = strings.TrimSpace(promptCacheKey) != ""
	state.snapshot.PromptCacheKeyAutoInjected = promptCacheKeyAutoInjected
	state.snapshot.PreviousResponseIDPresent = previousResponseIDPresent
}

func setCacheInstrumentationOpenAIChatState(c *gin.Context, body []byte, promptCacheKey string, promptCacheKeyAutoInjected bool) {
	state := getOrCreateCacheInstrumentationState(c)
	if state == nil {
		return
	}
	state.snapshot.CacheFamily = cacheFamilyOpenAIChatToResponses
	state.snapshot.SessionSignalSource = resolveOpenAISessionSignalSource(c, body)
	if state.snapshot.SessionSignalSource == "" && promptCacheKeyAutoInjected {
		state.snapshot.SessionSignalSource = "content_fallback"
	}
	state.snapshot.PromptCacheKeyPresent = strings.TrimSpace(promptCacheKey) != ""
	state.snapshot.PromptCacheKeyAutoInjected = promptCacheKeyAutoInjected
	state.snapshot.PreviousResponseIDPresent = false
}

func setCacheInstrumentationOpenAIMessagesState(
	c *gin.Context,
	body []byte,
	promptCacheKey string,
	promptCacheKeyAutoInjected bool,
	previousResponseIDPresent bool,
	sessionSignalSource string,
) {
	state := getOrCreateCacheInstrumentationState(c)
	if state == nil {
		return
	}
	state.snapshot.CacheFamily = cacheFamilyAnthropicToOpenAIResponses
	if normalized := strings.TrimSpace(sessionSignalSource); normalized != "" {
		state.snapshot.SessionSignalSource = normalized
	} else {
		state.snapshot.SessionSignalSource = resolveOpenAISessionSignalSource(c, body)
	}
	state.snapshot.PromptCacheKeyPresent = strings.TrimSpace(promptCacheKey) != ""
	state.snapshot.PromptCacheKeyAutoInjected = promptCacheKeyAutoInjected
	state.snapshot.PreviousResponseIDPresent = previousResponseIDPresent
}

func emitCacheInstrumentationForOpenAI(input *OpenAIRecordUsageInput, usageLog *UsageLog) {
	if input == nil || usageLog == nil {
		return
	}
	emitCacheInstrumentationSample(cacheInstrumentationEmitInput{
		Snapshot:            input.CacheInstrumentationSnapshot,
		RequestID:           usageLog.RequestID,
		UserID:              usageLog.UserID,
		AccountID:           usageLog.AccountID,
		Platform:            accountPlatformOrDefault(input.Account),
		RequestedModel:      usageLog.RequestedModel,
		UpstreamModel:       derefStringPtr(usageLog.UpstreamModel),
		InboundEndpoint:     derefStringPtr(usageLog.InboundEndpoint),
		UpstreamEndpoint:    derefStringPtr(usageLog.UpstreamEndpoint),
		RequestType:         usageLog.EffectiveRequestType().String(),
		CacheReadTokens:     usageLog.CacheReadTokens,
		CacheCreationTokens: usageLog.CacheCreationTokens,
	})
}

func emitCacheInstrumentationForGateway(input *recordUsageCoreInput, usageLog *UsageLog) {
	if input == nil || usageLog == nil {
		return
	}
	emitCacheInstrumentationSample(cacheInstrumentationEmitInput{
		Snapshot:            input.CacheInstrumentationSnapshot,
		RequestID:           usageLog.RequestID,
		UserID:              usageLog.UserID,
		AccountID:           usageLog.AccountID,
		Platform:            accountPlatformOrDefault(input.Account),
		RequestedModel:      usageLog.RequestedModel,
		UpstreamModel:       derefStringPtr(usageLog.UpstreamModel),
		InboundEndpoint:     derefStringPtr(usageLog.InboundEndpoint),
		UpstreamEndpoint:    derefStringPtr(usageLog.UpstreamEndpoint),
		RequestType:         usageLog.EffectiveRequestType().String(),
		CacheReadTokens:     usageLog.CacheReadTokens,
		CacheCreationTokens: usageLog.CacheCreationTokens,
	})
}

func emitCacheInstrumentationSample(input cacheInstrumentationEmitInput) {
	snapshot := input.Snapshot
	if snapshot == nil || strings.TrimSpace(snapshot.CacheFamily) == "" {
		return
	}
	sampleRate, sampleReason := cacheInstrumentationSamplePolicy(snapshot)
	if sampleRate <= 0 {
		return
	}
	if sampleRate < 1 && cacheInstrumentationSampleFloat64() >= sampleRate {
		return
	}

	fields := map[string]any{
		"request_id":                     strings.TrimSpace(input.RequestID),
		"user_id":                        input.UserID,
		"account_id":                     input.AccountID,
		"platform":                       strings.TrimSpace(input.Platform),
		"model":                          firstNonEmptyCacheInstrumentationString(input.RequestedModel, input.UpstreamModel),
		"cache_family":                   snapshot.CacheFamily,
		"session_signal_source":          snapshot.SessionSignalSource,
		"upstream_session_anchor_source": firstNonEmptyCacheInstrumentationString(snapshot.UpstreamSessionAnchorSource, "none"),
		"prompt_cache_key_present":       snapshot.PromptCacheKeyPresent,
		"prompt_cache_key_auto_injected": snapshot.PromptCacheKeyAutoInjected,
		"previous_response_id_present":   snapshot.PreviousResponseIDPresent,
		"sticky_account_hit":             snapshot.StickyAccountHit,
		"account_switch_happened":        snapshot.AccountSwitchHappened,
		"inbound_endpoint":               strings.TrimSpace(input.InboundEndpoint),
		"upstream_endpoint":              strings.TrimSpace(input.UpstreamEndpoint),
		"request_type":                   strings.TrimSpace(input.RequestType),
		"requested_model":                strings.TrimSpace(input.RequestedModel),
		"upstream_model":                 strings.TrimSpace(input.UpstreamModel),
		"cache_read_tokens":              input.CacheReadTokens,
		"cache_creation_tokens":          input.CacheCreationTokens,
		"sample_rate":                    sampleRate,
		"sample_reason":                  sampleReason,
	}
	if snapshot.SessionHashShort != "" {
		fields["session_hash_short"] = snapshot.SessionHashShort
	}
	if snapshot.SelectedAccountID > 0 {
		fields["selected_account_id"] = snapshot.SelectedAccountID
	}
	if snapshot.StickyAccountID > 0 {
		fields["sticky_account_id"] = snapshot.StickyAccountID
	}

	logger.WriteSinkEvent("info", cacheInstrumentationComponent, "cache_instrumentation", fields)
}

func cacheInstrumentationSamplePolicy(snapshot *CacheInstrumentationSnapshot) (float64, string) {
	if snapshot == nil {
		return 0, ""
	}
	if snapshot.AccountSwitchHappened || (snapshot.StickyAccountID > 0 && !snapshot.StickyAccountHit) {
		return 1, "sticky_miss_or_switch"
	}
	if snapshot.PreviousResponseIDPresent {
		return 0.10, "previous_response_id"
	}
	return 0.01, "default"
}

func resolveOpenAISessionSignalSource(c *gin.Context, body []byte) string {
	if c != nil {
		if strings.TrimSpace(c.GetHeader("session_id")) != "" {
			return "header_session_id"
		}
		if strings.TrimSpace(c.GetHeader("conversation_id")) != "" {
			return "header_conversation_id"
		}
	}
	if len(body) > 0 && strings.TrimSpace(gjson.GetBytes(body, "prompt_cache_key").String()) != "" {
		return "body_prompt_cache_key"
	}
	if len(body) > 0 && strings.TrimSpace(deriveOpenAIContentSessionSeed(body)) != "" {
		return "content_fallback"
	}
	return ""
}

func getOrCreateCacheInstrumentationState(c *gin.Context) *cacheInstrumentationState {
	if c == nil {
		return nil
	}
	if raw, ok := c.Get(cacheInstrumentationStateGinContextKey); ok {
		if state, ok := raw.(*cacheInstrumentationState); ok && state != nil {
			return state
		}
	}
	state := &cacheInstrumentationState{}
	c.Set(cacheInstrumentationStateGinContextKey, state)
	return state
}

func getCacheInstrumentationState(c *gin.Context) *cacheInstrumentationState {
	if c == nil {
		return nil
	}
	raw, ok := c.Get(cacheInstrumentationStateGinContextKey)
	if !ok {
		return nil
	}
	state, _ := raw.(*cacheInstrumentationState)
	return state
}

func accountPlatformOrDefault(account *Account) string {
	if account == nil {
		return ""
	}
	return strings.TrimSpace(account.Platform)
}

func derefStringPtr(v *string) string {
	if v == nil {
		return ""
	}
	return strings.TrimSpace(*v)
}

func firstNonEmptyCacheInstrumentationString(values ...string) string {
	for _, value := range values {
		if trimmed := strings.TrimSpace(value); trimmed != "" {
			return trimmed
		}
	}
	return ""
}

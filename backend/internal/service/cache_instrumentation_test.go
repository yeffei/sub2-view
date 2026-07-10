package service

import (
	"net/http/httptest"
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

type cacheInstrumentationCaptureSink struct {
	events []*logger.LogEvent
}

func (s *cacheInstrumentationCaptureSink) WriteLogEvent(event *logger.LogEvent) {
	if event == nil {
		return
	}
	s.events = append(s.events, event)
}

func TestCaptureCacheInstrumentationSnapshot(t *testing.T) {
	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	req := httptest.NewRequest("POST", "/v1/responses", nil)
	req.Header.Set("session_id", "sess_123")
	c.Request = req

	setCacheInstrumentationOpenAIResponsesState(c, []byte(`{"model":"gpt-5","previous_response_id":"resp_1"}`), "compat-key", true, true)
	RecordCacheInstrumentationUpstreamSessionAnchor(c, "header_session_id")
	RecordCacheInstrumentationRouting(c, "session_hash_123456", OpenAIAccountScheduleDecision{
		StickyAccountID:   11,
		SelectedAccountID: 22,
	})

	snapshot := CaptureCacheInstrumentationSnapshot(c)
	require.NotNil(t, snapshot)
	require.Equal(t, cacheFamilyOpenAIResponses, snapshot.CacheFamily)
	require.Equal(t, "session_", snapshot.SessionHashShort)
	require.Equal(t, "header_session_id", snapshot.SessionSignalSource)
	require.Equal(t, "header_session_id", snapshot.UpstreamSessionAnchorSource)
	require.True(t, snapshot.PromptCacheKeyPresent)
	require.True(t, snapshot.PromptCacheKeyAutoInjected)
	require.True(t, snapshot.PreviousResponseIDPresent)
	require.Equal(t, int64(11), snapshot.StickyAccountID)
	require.Equal(t, int64(22), snapshot.SelectedAccountID)
	require.True(t, snapshot.AccountSwitchHappened)
	require.False(t, snapshot.StickyAccountHit)
}

func TestEmitCacheInstrumentationSample_FullSampleOnAccountSwitch(t *testing.T) {
	sink := &cacheInstrumentationCaptureSink{}
	logger.SetSink(sink)
	originalSampler := cacheInstrumentationSampleFloat64
	cacheInstrumentationSampleFloat64 = func() float64 { return 0.99 }
	t.Cleanup(func() {
		cacheInstrumentationSampleFloat64 = originalSampler
		logger.SetSink(nil)
	})

	emitCacheInstrumentationSample(cacheInstrumentationEmitInput{
		Snapshot: &CacheInstrumentationSnapshot{
			CacheFamily:                 cacheFamilyAnthropicToOpenAIResponses,
			SessionHashShort:            "abcd1234",
			SessionSignalSource:         "anthropic_digest",
			UpstreamSessionAnchorSource: "prompt_cache_key",
			PromptCacheKeyPresent:       true,
			SelectedAccountID:           22,
			StickyAccountID:             11,
			AccountSwitchHappened:       true,
		},
		RequestID:           "local:req-1",
		UserID:              1,
		AccountID:           22,
		Platform:            "openai",
		RequestedModel:      "claude-sonnet-4-5",
		UpstreamModel:       "gpt-5",
		InboundEndpoint:     "/v1/messages",
		UpstreamEndpoint:    "/backend-api/codex/responses",
		RequestType:         "stream",
		CacheReadTokens:     0,
		CacheCreationTokens: 1280,
	})

	require.Len(t, sink.events, 1)
	event := sink.events[0]
	require.Equal(t, cacheInstrumentationComponent, event.Component)
	require.Equal(t, "cache_instrumentation", event.Message)
	require.Equal(t, "sticky_miss_or_switch", event.Fields["sample_reason"])
	require.Equal(t, 1.0, event.Fields["sample_rate"])
	require.Equal(t, int64(22), event.Fields["selected_account_id"])
	require.Equal(t, int64(11), event.Fields["sticky_account_id"])
	require.Equal(t, "anthropic_digest", event.Fields["session_signal_source"])
	require.Equal(t, "prompt_cache_key", event.Fields["upstream_session_anchor_source"])
	require.Equal(t, "/v1/messages", event.Fields["inbound_endpoint"])
}

func TestEmitCacheInstrumentationSample_SkipsWhenDefaultSampleMisses(t *testing.T) {
	sink := &cacheInstrumentationCaptureSink{}
	logger.SetSink(sink)
	originalSampler := cacheInstrumentationSampleFloat64
	cacheInstrumentationSampleFloat64 = func() float64 { return 0.5 }
	t.Cleanup(func() {
		cacheInstrumentationSampleFloat64 = originalSampler
		logger.SetSink(nil)
	})

	emitCacheInstrumentationSample(cacheInstrumentationEmitInput{
		Snapshot: &CacheInstrumentationSnapshot{
			CacheFamily:           cacheFamilyOpenAIChatToResponses,
			SessionSignalSource:   "content_fallback",
			PromptCacheKeyPresent: true,
		},
		RequestID:      "local:req-2",
		UserID:         1,
		AccountID:      2,
		Platform:       "openai",
		RequestedModel: "gpt-5",
		RequestType:    "sync",
	})

	require.Empty(t, sink.events)
}

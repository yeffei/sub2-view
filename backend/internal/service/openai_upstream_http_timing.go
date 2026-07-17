package service

import (
	"net/http"
	"net/http/httptrace"
	"sync/atomic"
	"time"

	"github.com/gin-gonic/gin"
)

// openAIUpstreamHTTPTrace captures transport events without changing the
// HTTPUpstream interface or the configured transport implementation.
type openAIUpstreamHTTPTrace struct {
	startedAtUnixNano      int64
	requestWrittenUnixNano atomic.Int64
	firstResponseUnixNano  atomic.Int64
}

func traceOpenAIUpstreamHTTPRequest(req *http.Request) (*http.Request, *openAIUpstreamHTTPTrace) {
	if req == nil {
		return nil, nil
	}

	traceState := &openAIUpstreamHTTPTrace{startedAtUnixNano: time.Now().UnixNano()}
	trace := &httptrace.ClientTrace{
		WroteRequest: func(httptrace.WroteRequestInfo) {
			traceState.requestWrittenUnixNano.CompareAndSwap(0, time.Now().UnixNano())
		},
		GotFirstResponseByte: func() {
			traceState.firstResponseUnixNano.CompareAndSwap(0, time.Now().UnixNano())
		},
	}
	return req.WithContext(httptrace.WithClientTrace(req.Context(), trace)), traceState
}

func recordOpenAIUpstreamHTTPTrace(c *gin.Context, traceState *openAIUpstreamHTTPTrace) {
	if c == nil || traceState == nil || traceState.startedAtUnixNano <= 0 {
		return
	}

	startedAt := traceState.startedAtUnixNano
	requestWrittenAt := traceState.requestWrittenUnixNano.Load()
	firstResponseAt := traceState.firstResponseUnixNano.Load()

	if requestWrittenAt > 0 {
		SetOpsLatencyMs(c, OpsUpstreamRequestWrittenMsKey, nonNegativeDurationMs(requestWrittenAt-startedAt))
	}
	if firstResponseAt > 0 {
		SetOpsLatencyMs(c, OpsUpstreamFirstByteMsKey, nonNegativeDurationMs(firstResponseAt-startedAt))
	}
	if firstResponseAt > 0 && requestWrittenAt > 0 {
		SetOpsLatencyMs(c, OpsUpstreamPostWriteWaitMsKey, nonNegativeDurationMs(firstResponseAt-requestWrittenAt))
	}
}

func nonNegativeDurationMs(deltaNanos int64) int64 {
	if deltaNanos <= 0 {
		return 0
	}
	return deltaNanos / int64(time.Millisecond)
}

package service

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestOpenAIUpstreamHTTPTrace_RecordsRequestWriteAndFirstByte(t *testing.T) {
	gin.SetMode(gin.TestMode)
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("ok"))
	}))
	defer upstream.Close()

	req, err := http.NewRequest(http.MethodPost, upstream.URL, bytes.NewReader([]byte("request-body")))
	require.NoError(t, err)
	tracedReq, traceState := traceOpenAIUpstreamHTTPRequest(req)
	resp, err := http.DefaultClient.Do(tracedReq)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Eventually(t, func() bool {
		return traceState.requestWrittenUnixNano.Load() > 0 && traceState.firstResponseUnixNano.Load() > 0
	}, time.Second, time.Millisecond)

	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	recordOpenAIUpstreamHTTPTrace(c, traceState)

	_, requestWritten := c.Get(OpsUpstreamRequestWrittenMsKey)
	_, firstByte := c.Get(OpsUpstreamFirstByteMsKey)
	require.True(t, requestWritten)
	require.True(t, firstByte)
	postWriteWait, ok := c.Get(OpsUpstreamPostWriteWaitMsKey)
	if ok {
		require.GreaterOrEqual(t, postWriteWait.(int64), int64(0))
	}
}

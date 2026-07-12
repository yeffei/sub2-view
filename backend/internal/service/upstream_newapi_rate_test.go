package service

import (
	"context"
	"io"
	"net/http"
	"strings"
	"sync"
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/pkg/tlsfingerprint"
	"github.com/stretchr/testify/require"
)

type newAPIRateProbeHTTPUpstream struct {
	mu       sync.Mutex
	logCalls int
}

func (u *newAPIRateProbeHTTPUpstream) Do(req *http.Request, _ string, _ int64, _ int) (*http.Response, error) {
	return u.respond(req)
}

func (u *newAPIRateProbeHTTPUpstream) DoWithTLS(req *http.Request, _ string, _ int64, _ int, _ *tlsfingerprint.Profile) (*http.Response, error) {
	return u.respond(req)
}

func (u *newAPIRateProbeHTTPUpstream) respond(req *http.Request) (*http.Response, error) {
	u.mu.Lock()
	defer u.mu.Unlock()

	status := http.StatusOK
	body := `{}`
	header := http.Header{"Content-Type": []string{"application/json"}}
	switch {
	case req.Method == http.MethodGet && strings.HasSuffix(req.URL.Path, "/api/log/token"):
		u.logCalls++
		if u.logCalls == 1 {
			body = `{"success":true,"message":"","data":[]}`
		} else {
			body = `{"success":true,"message":"","data":[{"type":2,"group":"gpt-pro","model_name":"gpt-5.4-mini","request_id":"req-probe","other":"{\"group_ratio\":0.13,\"user_group_ratio\":-1}"}]}`
		}
	case req.Method == http.MethodGet && strings.HasSuffix(req.URL.Path, "/v1/models"):
		body = `{"data":[{"id":"gpt-5.4-mini"}]}`
	case req.Method == http.MethodPost && strings.HasSuffix(req.URL.Path, "/v1/responses"):
		header.Set("X-Oneapi-Request-Id", "req-probe")
		body = `{"id":"resp_probe"}`
	default:
		status = http.StatusNotFound
	}

	return &http.Response{
		StatusCode: status,
		Header:     header,
		Body:       io.NopCloser(strings.NewReader(body)),
	}, nil
}

func TestProbeNewAPIUpstreamRateMultiplier(t *testing.T) {
	upstream := &newAPIRateProbeHTTPUpstream{}
	svc := &AccountTestService{httpUpstream: upstream, cfg: upstreamModelSyncTestConfig()}
	account := &Account{
		ID:          4,
		Platform:    PlatformOpenAI,
		Type:        AccountTypeAPIKey,
		Credentials: map[string]any{"api_key": "sk-test", "base_url": "https://upstream.example.com/v1"},
	}

	meta, err := svc.probeNewAPIUpstreamRateMultiplier(context.Background(), account, "https://upstream.example.com/v1", "sk-test")
	require.NoError(t, err)
	require.Equal(t, 0.13, meta.RateMultiplier)
	require.Equal(t, "gpt-pro", meta.GroupName)
	require.Equal(t, "new_api_token_log_probe", meta.RateSource)
}

func TestExtractNewAPILogRateMultiplier(t *testing.T) {
	rate, err := extractNewAPILogRateMultiplier([]byte(`"{\"group_ratio\":0.08,\"user_group_ratio\":0.06}"`))
	require.NoError(t, err)
	require.Equal(t, 0.06, rate)

	rate, err = extractNewAPILogRateMultiplier([]byte(`{"group_ratio":0.08,"user_group_ratio":-1}`))
	require.NoError(t, err)
	require.Equal(t, 0.08, rate)
}

func TestBuildNewAPITokenLogURL(t *testing.T) {
	require.Equal(t, "https://gateway.example.com/api/log/token", buildNewAPITokenLogURL("https://gateway.example.com"))
	require.Equal(t, "https://gateway.example.com/api/log/token", buildNewAPITokenLogURL("https://gateway.example.com/v1"))
}

func TestSelectNewAPIRateProbeModelFromLogs(t *testing.T) {
	logs := []newAPITokenLog{
		{Type: 5, ModelName: "gpt-5.4-mini"},
		{Type: 2, ModelName: "gpt-5.5"},
		{Type: 2, ModelName: "gpt-5.4-mini"},
	}
	require.Equal(t, "gpt-5.4-mini", selectNewAPIRateProbeModelFromLogs(logs))
}

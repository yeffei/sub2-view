package service

import (
	"context"
	"io"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/pkg/tlsfingerprint"
	"github.com/stretchr/testify/require"
)

type legacyRateProbeHTTPUpstream struct {
	mu                sync.Mutex
	usageCalls        int
	afterRequestCount int64
}

func (u *legacyRateProbeHTTPUpstream) Do(req *http.Request, _ string, _ int64, _ int) (*http.Response, error) {
	return u.respond(req)
}

func (u *legacyRateProbeHTTPUpstream) DoWithTLS(req *http.Request, _ string, _ int64, _ int, _ *tlsfingerprint.Profile) (*http.Response, error) {
	return u.respond(req)
}

func (u *legacyRateProbeHTTPUpstream) respond(req *http.Request) (*http.Response, error) {
	u.mu.Lock()
	defer u.mu.Unlock()

	status := http.StatusOK
	body := `{}`
	switch {
	case req.Method == http.MethodGet && strings.HasSuffix(req.URL.Path, "/v1/usage"):
		u.usageCalls++
		if u.usageCalls == 1 {
			body = `{"usage":{"today":{"requests":265,"cost":20.03,"actual_cost":1.71}}}`
		} else {
			body = `{"usage":{"today":{"requests":` + strconv.FormatInt(265+u.afterRequestCount, 10) + `,"cost":20.03058545,"actual_cost":1.7100761085}}}`
		}
	case req.Method == http.MethodGet && strings.HasSuffix(req.URL.Path, "/v1/models"):
		body = `{"data":[{"id":"gpt-5.4-mini"}]}`
	case req.Method == http.MethodPost && strings.HasSuffix(req.URL.Path, "/v1/responses"):
		body = `{"id":"resp_probe","usage":{"input_tokens":10,"output_tokens":1}}`
	default:
		status = http.StatusNotFound
	}

	return &http.Response{
		StatusCode: status,
		Header:     http.Header{"Content-Type": []string{"application/json"}},
		Body:       io.NopCloser(strings.NewReader(body)),
	}, nil
}

func TestProbeLegacyUpstreamRateMultiplier(t *testing.T) {
	upstream := &legacyRateProbeHTTPUpstream{afterRequestCount: 1}
	svc := &AccountTestService{httpUpstream: upstream, cfg: upstreamModelSyncTestConfig()}
	account := &Account{
		ID:       2092,
		Platform: PlatformOpenAI,
		Type:     AccountTypeAPIKey,
		Credentials: map[string]any{
			"api_key":  "sk-test",
			"base_url": "https://upstream.example.com",
		},
	}

	meta, err := svc.probeLegacyUpstreamRateMultiplier(context.Background(), account, "https://upstream.example.com", "sk-test")
	require.NoError(t, err)
	require.Equal(t, 0.13, meta.RateMultiplier)
	require.Equal(t, "legacy_usage_delta_probe", meta.RateSource)
}

func TestProbeLegacyUpstreamRateMultiplierRejectsConcurrentUsage(t *testing.T) {
	upstream := &legacyRateProbeHTTPUpstream{afterRequestCount: 2}
	svc := &AccountTestService{httpUpstream: upstream, cfg: upstreamModelSyncTestConfig()}
	account := &Account{
		ID:          2092,
		Platform:    PlatformOpenAI,
		Type:        AccountTypeAPIKey,
		Credentials: map[string]any{"api_key": "sk-test", "base_url": "https://upstream.example.com"},
	}

	_, err := svc.probeLegacyUpstreamRateMultiplier(context.Background(), account, "https://upstream.example.com", "sk-test")
	require.Error(t, err)
	require.Contains(t, err.Error(), "并发请求")
}

func TestSelectLegacyRateProbeModel(t *testing.T) {
	require.Equal(t, "gpt-5.4-mini", selectLegacyRateProbeModel([]string{"gpt-5.4", "gpt-5.4-mini"}))
	require.Equal(t, "gpt-5.6-terra", selectLegacyRateProbeModel([]string{"gpt-image-2", "gpt-5.6-terra"}))
	require.Empty(t, selectLegacyRateProbeModel([]string{"gpt-image-2", "gpt-4o-realtime-preview"}))
}

func TestBuildCompatibleResponsesURL(t *testing.T) {
	require.Equal(t, "https://gateway.example.com/v1/responses", buildCompatibleResponsesURL("https://gateway.example.com"))
	require.Equal(t, "https://gateway.example.com/v1/responses", buildCompatibleResponsesURL("https://gateway.example.com/v1"))
}

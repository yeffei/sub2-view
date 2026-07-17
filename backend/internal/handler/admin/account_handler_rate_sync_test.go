package admin

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/Wei-Shaw/sub2api/internal/pkg/tlsfingerprint"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/stretchr/testify/require"
)

type rateSyncHTTPUpstream struct{}

func (rateSyncHTTPUpstream) Do(req *http.Request, _ string, _ int64, _ int) (*http.Response, error) {
	return rateSyncResponse(req), nil
}

func (rateSyncHTTPUpstream) DoWithTLS(req *http.Request, _ string, _ int64, _ int, _ *tlsfingerprint.Profile) (*http.Response, error) {
	return rateSyncResponse(req), nil
}

func rateSyncResponse(req *http.Request) *http.Response {
	status := http.StatusNotFound
	body := `{}`
	if req.URL.Host == "ok.example.com" && strings.HasSuffix(req.URL.Path, "/v1/account/meta") {
		status = http.StatusOK
		body = `{"data":{"compatible":true,"platform":"openai","rate_multiplier":0.08,"rate_source":"user_group_override"}}`
	}
	return &http.Response{StatusCode: status, Header: make(http.Header), Body: io.NopCloser(strings.NewReader(body))}
}

func TestBatchSyncUpstreamRateMultipliersNormalizesIDsAndReportsPartialFailure(t *testing.T) {
	adminSvc := newStubAdminService()
	adminSvc.getAccountsByIDsFunc = func(_ context.Context, ids []int64) ([]*service.Account, error) {
		require.Equal(t, []int64{1, 2, 404}, ids)
		return []*service.Account{
			{
				ID: 1, Platform: service.PlatformOpenAI, Type: service.AccountTypeAPIKey,
				Credentials: map[string]any{"api_key": "sk-ok", "base_url": "https://ok.example.com/v1"},
			},
			{
				ID: 2, Platform: service.PlatformOpenAI, Type: service.AccountTypeAPIKey,
				Credentials: map[string]any{"api_key": "sk-fail", "base_url": "https://fail.example.com/v1"},
			},
		}, nil
	}
	cfg := &config.Config{}
	cfg.Security.URLAllowlist.AllowInsecureHTTP = true
	testSvc := service.NewAccountTestService(nil, nil, nil, nil, nil, rateSyncHTTPUpstream{}, cfg, nil)
	handler := NewAccountHandler(adminSvc, nil, nil, nil, nil, nil, nil, nil, testSvc, nil, nil, nil, nil, nil)

	payload, err := handler.batchSyncUpstreamRateMultipliers(context.Background(), []int64{0, 1, 2, 1, -5, 404})
	require.NoError(t, err)
	require.Equal(t, 3, payload["total"])
	require.Equal(t, 1, payload["success"])
	require.Equal(t, 2, payload["failed"])

	results := payload["results"].([]BatchSyncUpstreamRateMultiplierResult)
	require.Len(t, results, 3)
	require.True(t, results[0].Success)
	require.Equal(t, 1.0, results[0].PreviousRateMultiplier)
	require.Equal(t, 0.08, results[0].RateMultiplier)
	require.True(t, results[0].Changed)
	require.True(t, results[0].SignificantChange)
	require.False(t, results[1].Success)
	require.NotEmpty(t, results[1].Error)
	require.Equal(t, "account not found", results[2].Error)
}

func TestIsSignificantRateMultiplierChange(t *testing.T) {
	require.False(t, isSignificantRateMultiplierChange(0.1, 0.1))
	require.False(t, isSignificantRateMultiplierChange(0.1, 0.12))
	require.True(t, isSignificantRateMultiplierChange(0.1, 0.15))
	require.True(t, isSignificantRateMultiplierChange(0.1, 0.05))
	require.True(t, isSignificantRateMultiplierChange(0, 0.08))
}

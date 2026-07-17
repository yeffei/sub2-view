package service

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/Wei-Shaw/sub2api/internal/pkg/tlsfingerprint"
	"github.com/stretchr/testify/require"
)

type upstreamRateSyncAdminStub struct {
	accounts []*Account
	updated  map[int64]float64
}

func (s *upstreamRateSyncAdminStub) ListAccounts(_ context.Context, page, _ int, _, accountType, _, _, _ string, _ int64, _, _, _ string) ([]Account, int64, error) {
	if page > 1 || accountType != AccountTypeAPIKey {
		return []Account{}, int64(len(s.accounts)), nil
	}
	result := make([]Account, 0, len(s.accounts))
	for _, account := range s.accounts {
		if account != nil && account.Type == AccountTypeAPIKey {
			result = append(result, *account)
		}
	}
	return result, int64(len(result)), nil
}

func (s *upstreamRateSyncAdminStub) GetAccountsByIDs(_ context.Context, ids []int64) ([]*Account, error) {
	result := make([]*Account, 0, len(ids))
	for _, id := range ids {
		for _, account := range s.accounts {
			if account != nil && account.ID == id {
				result = append(result, account)
			}
		}
	}
	return result, nil
}

func (s *upstreamRateSyncAdminStub) UpdateAccount(_ context.Context, id int64, input *UpdateAccountInput) (*Account, error) {
	if s.updated == nil {
		s.updated = map[int64]float64{}
	}
	if input != nil && input.RateMultiplier != nil {
		s.updated[id] = *input.RateMultiplier
	}
	for _, account := range s.accounts {
		if account != nil && account.ID == id {
			copy := *account
			copy.RateMultiplier = input.RateMultiplier
			return &copy, nil
		}
	}
	return nil, nil
}

type upstreamRateSyncHTTPStub struct{}

func (upstreamRateSyncHTTPStub) Do(req *http.Request, _ string, _ int64, _ int) (*http.Response, error) {
	return upstreamRateSyncHTTPResponse(req), nil
}

func (upstreamRateSyncHTTPStub) DoWithTLS(req *http.Request, _ string, _ int64, _ int, _ *tlsfingerprint.Profile) (*http.Response, error) {
	return upstreamRateSyncHTTPResponse(req), nil
}

func upstreamRateSyncHTTPResponse(req *http.Request) *http.Response {
	status := http.StatusNotFound
	body := `{}`
	if req.URL.Host == "ok.example.com" && strings.HasSuffix(req.URL.Path, "/v1/account/meta") {
		status = http.StatusOK
		body = `{"data":{"compatible":true,"platform":"openai","rate_multiplier":0.08,"rate_source":"user_group_override"}}`
	}
	return &http.Response{StatusCode: status, Header: make(http.Header), Body: io.NopCloser(strings.NewReader(body))}
}

func TestUpstreamRateSyncServiceSyncAllReportsPartialFailure(t *testing.T) {
	admin := &upstreamRateSyncAdminStub{accounts: []*Account{
		{ID: 1, Name: "ok", Platform: PlatformOpenAI, Type: AccountTypeAPIKey, Credentials: map[string]any{"api_key": "sk-ok", "base_url": "https://ok.example.com/v1"}},
		{ID: 2, Name: "fail", Platform: PlatformOpenAI, Type: AccountTypeAPIKey, Credentials: map[string]any{"api_key": "sk-fail", "base_url": "https://fail.example.com/v1"}},
	}}
	cfg := &config.Config{}
	cfg.Security.URLAllowlist.AllowInsecureHTTP = true
	testSvc := NewAccountTestService(nil, nil, nil, nil, nil, upstreamRateSyncHTTPStub{}, cfg, nil)
	svc := NewUpstreamRateSyncService(admin, testSvc)

	result, err := svc.SyncAllAPIKeyAccounts(context.Background())
	require.NoError(t, err)
	require.Equal(t, 2, result.Total)
	require.Equal(t, 1, result.Success)
	require.Equal(t, 1, result.Failed)
	require.Equal(t, 0.08, admin.updated[1])
	require.True(t, result.Items[0].Result.Changed)
	require.True(t, result.Items[0].Result.SignificantChange)
	require.Error(t, result.Items[1].Error)
}

func TestUpstreamRateSyncScheduleRunsThreeTimesInShanghai(t *testing.T) {
	location, err := time.LoadLocation("Asia/Shanghai")
	require.NoError(t, err)
	schedule, err := upstreamRateSyncCronParser.Parse(upstreamRateSyncSchedule)
	require.NoError(t, err)

	next := schedule.Next(time.Date(2026, 7, 12, 1, 0, 0, 0, location))
	require.Equal(t, 2, next.Hour())
	next = schedule.Next(next)
	require.Equal(t, 10, next.Hour())
	next = schedule.Next(next)
	require.Equal(t, 18, next.Hour())
	next = schedule.Next(next)
	require.Equal(t, 2, next.Hour())
	require.Equal(t, 13, next.Day())
}

func TestIsUpstreamRateSyncEnabledIsOptIn(t *testing.T) {
	disabled := NewSettingService(&gatewayTTLSettingRepo{data: map[string]string{}}, &config.Config{})
	require.False(t, disabled.IsUpstreamRateSyncEnabled(context.Background()))

	enabled := NewSettingService(&gatewayTTLSettingRepo{data: map[string]string{SettingKeyUpstreamRateSyncEnabled: "true"}}, &config.Config{})
	require.True(t, enabled.IsUpstreamRateSyncEnabled(context.Background()))
}

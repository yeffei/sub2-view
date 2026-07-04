//go:build unit

package service

import (
	"context"
	"errors"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/Wei-Shaw/sub2api/internal/pkg/pagination"
	"github.com/stretchr/testify/require"
)

var _ AccountRepository = (*upstreamPoolRecoveryProbeAccountRepoStub)(nil)
var _ UpstreamPoolRepository = (*upstreamPoolRecoveryProbeRepoStub)(nil)

type upstreamPoolRecoveryProbeAccountRepoStub struct {
	accounts       map[int64]*Account
	tempUnschedIDs []int64
	tempReasons    []string
}

func (r *upstreamPoolRecoveryProbeAccountRepoStub) Create(ctx context.Context, account *Account) error {
	return nil
}

func (r *upstreamPoolRecoveryProbeAccountRepoStub) GetByID(ctx context.Context, id int64) (*Account, error) {
	if r == nil || r.accounts == nil {
		return nil, errors.New("not found")
	}
	account, ok := r.accounts[id]
	if !ok {
		return nil, errors.New("not found")
	}
	return account, nil
}

func (r *upstreamPoolRecoveryProbeAccountRepoStub) GetByIDs(ctx context.Context, ids []int64) ([]*Account, error) {
	return nil, nil
}

func (r *upstreamPoolRecoveryProbeAccountRepoStub) ExistsByID(ctx context.Context, id int64) (bool, error) {
	_, ok := r.accounts[id]
	return ok, nil
}

func (r *upstreamPoolRecoveryProbeAccountRepoStub) GetByCRSAccountID(ctx context.Context, crsAccountID string) (*Account, error) {
	return nil, nil
}

func (r *upstreamPoolRecoveryProbeAccountRepoStub) FindByExtraField(ctx context.Context, key string, value any) ([]Account, error) {
	return nil, nil
}

func (r *upstreamPoolRecoveryProbeAccountRepoStub) ListCRSAccountIDs(ctx context.Context) (map[string]int64, error) {
	return nil, nil
}

func (r *upstreamPoolRecoveryProbeAccountRepoStub) Update(ctx context.Context, account *Account) error {
	return nil
}
func (r *upstreamPoolRecoveryProbeAccountRepoStub) Delete(ctx context.Context, id int64) error {
	return nil
}
func (r *upstreamPoolRecoveryProbeAccountRepoStub) List(ctx context.Context, params pagination.PaginationParams) ([]Account, *pagination.PaginationResult, error) {
	return nil, nil, nil
}
func (r *upstreamPoolRecoveryProbeAccountRepoStub) ListWithFilters(ctx context.Context, params pagination.PaginationParams, platform, accountType, status, anomalyReason, search string, groupID int64, privacyMode string) ([]Account, *pagination.PaginationResult, error) {
	return nil, nil, nil
}
func (r *upstreamPoolRecoveryProbeAccountRepoStub) ListByGroup(ctx context.Context, groupID int64) ([]Account, error) {
	return nil, nil
}
func (r *upstreamPoolRecoveryProbeAccountRepoStub) ListActive(ctx context.Context) ([]Account, error) {
	return nil, nil
}
func (r *upstreamPoolRecoveryProbeAccountRepoStub) ListOAuthRefreshCandidates(ctx context.Context) ([]Account, error) {
	return nil, nil
}
func (r *upstreamPoolRecoveryProbeAccountRepoStub) ListByPlatform(ctx context.Context, platform string) ([]Account, error) {
	return nil, nil
}
func (r *upstreamPoolRecoveryProbeAccountRepoStub) UpdateLastUsed(ctx context.Context, id int64) error {
	return nil
}
func (r *upstreamPoolRecoveryProbeAccountRepoStub) BatchUpdateLastUsed(ctx context.Context, updates map[int64]time.Time) error {
	return nil
}
func (r *upstreamPoolRecoveryProbeAccountRepoStub) SetError(ctx context.Context, id int64, errorMsg string) error {
	return nil
}
func (r *upstreamPoolRecoveryProbeAccountRepoStub) ClearError(ctx context.Context, id int64) error {
	return nil
}
func (r *upstreamPoolRecoveryProbeAccountRepoStub) SetSchedulable(ctx context.Context, id int64, schedulable bool) error {
	return nil
}
func (r *upstreamPoolRecoveryProbeAccountRepoStub) AutoPauseExpiredAccounts(ctx context.Context, now time.Time) (int64, error) {
	return 0, nil
}
func (r *upstreamPoolRecoveryProbeAccountRepoStub) GetGroups(ctx context.Context, accountID int64) ([]Group, error) {
	return nil, nil
}
func (r *upstreamPoolRecoveryProbeAccountRepoStub) BindGroups(ctx context.Context, accountID int64, groupIDs []int64) error {
	return nil
}
func (r *upstreamPoolRecoveryProbeAccountRepoStub) ListSchedulable(ctx context.Context) ([]Account, error) {
	return nil, nil
}
func (r *upstreamPoolRecoveryProbeAccountRepoStub) ListSchedulableByGroupID(ctx context.Context, groupID int64) ([]Account, error) {
	return nil, nil
}
func (r *upstreamPoolRecoveryProbeAccountRepoStub) ListSchedulableByPlatform(ctx context.Context, platform string) ([]Account, error) {
	return nil, nil
}
func (r *upstreamPoolRecoveryProbeAccountRepoStub) ListSchedulableByGroupIDAndPlatform(ctx context.Context, groupID int64, platform string) ([]Account, error) {
	return nil, nil
}
func (r *upstreamPoolRecoveryProbeAccountRepoStub) ListSchedulableByPlatforms(ctx context.Context, platforms []string) ([]Account, error) {
	return nil, nil
}
func (r *upstreamPoolRecoveryProbeAccountRepoStub) ListSchedulableByGroupIDAndPlatforms(ctx context.Context, groupID int64, platforms []string) ([]Account, error) {
	return nil, nil
}
func (r *upstreamPoolRecoveryProbeAccountRepoStub) ListSchedulableUngroupedByPlatform(ctx context.Context, platform string) ([]Account, error) {
	return nil, nil
}
func (r *upstreamPoolRecoveryProbeAccountRepoStub) ListSchedulableUngroupedByPlatforms(ctx context.Context, platforms []string) ([]Account, error) {
	return nil, nil
}
func (r *upstreamPoolRecoveryProbeAccountRepoStub) SetRateLimited(ctx context.Context, id int64, resetAt time.Time) error {
	return nil
}
func (r *upstreamPoolRecoveryProbeAccountRepoStub) SetModelRateLimit(ctx context.Context, id int64, scope string, resetAt time.Time, reason ...string) error {
	return nil
}
func (r *upstreamPoolRecoveryProbeAccountRepoStub) SetOverloaded(ctx context.Context, id int64, until time.Time) error {
	return nil
}
func (r *upstreamPoolRecoveryProbeAccountRepoStub) SetTempUnschedulable(ctx context.Context, id int64, until time.Time, reason string) error {
	r.tempUnschedIDs = append(r.tempUnschedIDs, id)
	r.tempReasons = append(r.tempReasons, reason)
	if account := r.accounts[id]; account != nil {
		account.TempUnschedulableUntil = &until
		account.TempUnschedulableReason = reason
	}
	return nil
}
func (r *upstreamPoolRecoveryProbeAccountRepoStub) ClearTempUnschedulable(ctx context.Context, id int64) error {
	return nil
}
func (r *upstreamPoolRecoveryProbeAccountRepoStub) ClearRateLimit(ctx context.Context, id int64) error {
	return nil
}
func (r *upstreamPoolRecoveryProbeAccountRepoStub) ClearAntigravityQuotaScopes(ctx context.Context, id int64) error {
	return nil
}
func (r *upstreamPoolRecoveryProbeAccountRepoStub) ClearModelRateLimits(ctx context.Context, id int64) error {
	return nil
}
func (r *upstreamPoolRecoveryProbeAccountRepoStub) UpdateSessionWindow(ctx context.Context, id int64, start, end *time.Time, status string) error {
	return nil
}
func (r *upstreamPoolRecoveryProbeAccountRepoStub) UpdateSessionWindowEnd(ctx context.Context, id int64, end time.Time) error {
	return nil
}
func (r *upstreamPoolRecoveryProbeAccountRepoStub) UpdateExtra(ctx context.Context, id int64, updates map[string]any) error {
	return nil
}
func (r *upstreamPoolRecoveryProbeAccountRepoStub) BulkUpdate(ctx context.Context, ids []int64, updates AccountBulkUpdate) (int64, error) {
	return 0, nil
}
func (r *upstreamPoolRecoveryProbeAccountRepoStub) IncrementQuotaUsed(ctx context.Context, id int64, amount float64) error {
	return nil
}
func (r *upstreamPoolRecoveryProbeAccountRepoStub) ResetQuotaUsed(ctx context.Context, id int64) error {
	return nil
}
func (r *upstreamPoolRecoveryProbeAccountRepoStub) RevertProxyFallback(ctx context.Context, accountID int64) error {
	return nil
}

type upstreamPoolRecoveryProbeRepoStub struct {
	pools   []UpstreamPool
	members map[int64][]UpstreamPoolMember
}

func (r *upstreamPoolRecoveryProbeRepoStub) ListUpstreamPools(ctx context.Context) ([]UpstreamPool, error) {
	return append([]UpstreamPool(nil), r.pools...), nil
}

func (r *upstreamPoolRecoveryProbeRepoStub) GetUpstreamPoolByID(ctx context.Context, id int64) (*UpstreamPool, error) {
	return nil, ErrUpstreamPoolNotFound
}

func (r *upstreamPoolRecoveryProbeRepoStub) CreateUpstreamPool(ctx context.Context, input *UpstreamPool) (*UpstreamPool, error) {
	return input, nil
}

func (r *upstreamPoolRecoveryProbeRepoStub) UpdateUpstreamPool(ctx context.Context, input *UpstreamPool) (*UpstreamPool, error) {
	return input, nil
}

func (r *upstreamPoolRecoveryProbeRepoStub) DeleteUpstreamPool(ctx context.Context, id int64) error {
	return nil
}

func (r *upstreamPoolRecoveryProbeRepoStub) ListUpstreamPoolMembers(ctx context.Context, poolID int64) ([]UpstreamPoolMember, error) {
	return append([]UpstreamPoolMember(nil), r.members[poolID]...), nil
}

func (r *upstreamPoolRecoveryProbeRepoStub) GetUpstreamPoolMemberByID(ctx context.Context, id int64) (*UpstreamPoolMember, error) {
	return nil, ErrUpstreamPoolNotFound
}

func (r *upstreamPoolRecoveryProbeRepoStub) CreateUpstreamPoolMember(ctx context.Context, input *UpstreamPoolMember) (*UpstreamPoolMember, error) {
	return input, nil
}

func (r *upstreamPoolRecoveryProbeRepoStub) UpdateUpstreamPoolMember(ctx context.Context, input *UpstreamPoolMember) (*UpstreamPoolMember, error) {
	return input, nil
}

func (r *upstreamPoolRecoveryProbeRepoStub) DeleteUpstreamPoolMember(ctx context.Context, id int64) error {
	return nil
}

func (r *upstreamPoolRecoveryProbeRepoStub) ListUpstreamPoolBindings(ctx context.Context) ([]UpstreamPoolBinding, error) {
	return nil, nil
}

func (r *upstreamPoolRecoveryProbeRepoStub) GetUpstreamPoolBindingByID(ctx context.Context, id int64) (*UpstreamPoolBinding, error) {
	return nil, ErrUpstreamPoolNotFound
}

func (r *upstreamPoolRecoveryProbeRepoStub) CreateUpstreamPoolBinding(ctx context.Context, input *UpstreamPoolBinding) (*UpstreamPoolBinding, error) {
	return input, nil
}

func (r *upstreamPoolRecoveryProbeRepoStub) UpdateUpstreamPoolBinding(ctx context.Context, input *UpstreamPoolBinding) (*UpstreamPoolBinding, error) {
	return input, nil
}

func (r *upstreamPoolRecoveryProbeRepoStub) DeleteUpstreamPoolBinding(ctx context.Context, id int64) error {
	return nil
}

func (r *upstreamPoolRecoveryProbeRepoStub) ListEnabledMemberAccountIDsByGroupAndPlatform(ctx context.Context, groupID int64, platform string) (map[int64]struct{}, error) {
	return nil, nil
}

func (r *upstreamPoolRecoveryProbeRepoStub) GetResolvedBindingByGroupAndPlatform(ctx context.Context, groupID int64, platform string) (*UpstreamPoolResolvedBinding, error) {
	return nil, nil
}

func (r *upstreamPoolRecoveryProbeRepoStub) GetOpenAIRoutingPolicy(ctx context.Context, groupID int64) (*OpenAIRoutingPolicy, error) {
	return nil, nil
}

type upstreamPoolRecoveryProbeTestSvc struct {
	mu      sync.Mutex
	calls   []int64
	modes   []string
	results map[int64]string
}

func (s *upstreamPoolRecoveryProbeTestSvc) RunTestBackground(ctx context.Context, accountID int64, modelID string) (*ScheduledTestResult, error) {
	return s.RunTestBackgroundWithMode(ctx, accountID, modelID, AccountTestModeDefault)
}

func (s *upstreamPoolRecoveryProbeTestSvc) RunTestBackgroundWithMode(ctx context.Context, accountID int64, modelID string, mode string) (*ScheduledTestResult, error) {
	s.mu.Lock()
	s.calls = append(s.calls, accountID)
	s.modes = append(s.modes, mode)
	s.mu.Unlock()
	if s.results != nil {
		if status := strings.TrimSpace(s.results[accountID]); status != "" {
			return &ScheduledTestResult{Status: status}, nil
		}
	}
	if accountID == 101 {
		return &ScheduledTestResult{Status: "success"}, nil
	}
	return &ScheduledTestResult{Status: "failed"}, nil
}

type upstreamPoolRecoveryProbeRateLimitStub struct {
	mu    sync.Mutex
	calls []int64
}

func (s *upstreamPoolRecoveryProbeRateLimitStub) RecoverAccountAfterSuccessfulTest(ctx context.Context, accountID int64) (*SuccessfulTestRecoveryResult, error) {
	s.mu.Lock()
	s.calls = append(s.calls, accountID)
	s.mu.Unlock()
	return &SuccessfulTestRecoveryResult{ClearedError: true, ClearedRateLimit: true}, nil
}

func TestUpstreamPoolRecoveryProbeRunnerService_ProbesRecoverablePoolMembers(t *testing.T) {
	accountRepo := &upstreamPoolRecoveryProbeAccountRepoStub{
		accounts: map[int64]*Account{
			101: {ID: 101, Status: StatusError, Schedulable: true, Platform: PlatformOpenAI},
			102: {ID: 102, Status: StatusActive, Schedulable: true, Platform: PlatformOpenAI},
			103: {ID: 103, Status: StatusActive, Schedulable: true, Platform: PlatformOpenAI, RateLimitResetAt: ptrTimeUpstreamPoolRecoveryProbe(time.Now().Add(10 * time.Minute))},
		},
	}
	upstreamRepo := &upstreamPoolRecoveryProbeRepoStub{
		pools: []UpstreamPool{
			{ID: 1, Enabled: true, Platform: PlatformOpenAI},
			{ID: 2, Enabled: true, Platform: PlatformAnthropic},
		},
		members: map[int64][]UpstreamPoolMember{
			1: {
				{AccountID: 101, Enabled: true},
				{AccountID: 102, Enabled: true, ManualDrained: true},
				{AccountID: 103, Enabled: true},
			},
			2: {
				{AccountID: 999, Enabled: true},
			},
		},
	}
	testSvc := &upstreamPoolRecoveryProbeTestSvc{}
	rateLimitSvc := &upstreamPoolRecoveryProbeRateLimitStub{}
	svc := NewUpstreamPoolRecoveryProbeRunnerService(upstreamRepo, accountRepo, testSvc, rateLimitSvc, &config.Config{})

	svc.runScheduled()

	require.ElementsMatch(t, []int64{101, 103}, testSvc.calls)
	require.ElementsMatch(t, []string{AccountTestModeCompact, AccountTestModeCompact}, testSvc.modes)
	require.Equal(t, []int64{101}, rateLimitSvc.calls)
}

func TestUpstreamPoolRecoveryProbeRunnerService_SkipsHealthyAccount(t *testing.T) {
	accountRepo := &upstreamPoolRecoveryProbeAccountRepoStub{
		accounts: map[int64]*Account{
			201: {ID: 201, Status: StatusActive, Schedulable: true, Platform: PlatformOpenAI},
		},
	}
	upstreamRepo := &upstreamPoolRecoveryProbeRepoStub{
		pools: []UpstreamPool{{ID: 1, Enabled: true, Platform: PlatformOpenAI}},
		members: map[int64][]UpstreamPoolMember{
			1: {{AccountID: 201, Enabled: true}},
		},
	}
	testSvc := &upstreamPoolRecoveryProbeTestSvc{}
	rateLimitSvc := &upstreamPoolRecoveryProbeRateLimitStub{}
	svc := NewUpstreamPoolRecoveryProbeRunnerService(upstreamRepo, accountRepo, testSvc, rateLimitSvc, &config.Config{})

	svc.runScheduled()

	require.Empty(t, testSvc.calls)
	require.Empty(t, rateLimitSvc.calls)
}

func TestUpstreamPoolRecoveryProbeRunnerService_ExtendsTempUnschedOnFailedProbe(t *testing.T) {
	now := time.Now()
	accountRepo := &upstreamPoolRecoveryProbeAccountRepoStub{
		accounts: map[int64]*Account{
			301: {
				ID:                      301,
				Status:                  StatusActive,
				Schedulable:             true,
				Platform:                PlatformOpenAI,
				TempUnschedulableUntil:  ptrTimeUpstreamPoolRecoveryProbe(now.Add(30 * time.Second)),
				TempUnschedulableReason: `{"matched_keyword":"pool_mode_5xx"}`,
				LastUsedAt:              ptrTimeUpstreamPoolRecoveryProbe(now.Add(-10 * time.Minute)),
			},
		},
	}
	upstreamRepo := &upstreamPoolRecoveryProbeRepoStub{
		pools: []UpstreamPool{{ID: 1, Enabled: true, Platform: PlatformOpenAI}},
		members: map[int64][]UpstreamPoolMember{
			1: {{AccountID: 301, Enabled: true}},
		},
	}
	testSvc := &upstreamPoolRecoveryProbeTestSvc{results: map[int64]string{301: "failed"}}
	rateLimitSvc := &upstreamPoolRecoveryProbeRateLimitStub{}
	svc := NewUpstreamPoolRecoveryProbeRunnerService(upstreamRepo, accountRepo, testSvc, rateLimitSvc, &config.Config{})

	svc.runScheduled()

	require.Equal(t, []int64{301}, testSvc.calls)
	require.Empty(t, rateLimitSvc.calls)
	require.Equal(t, []int64{301}, accountRepo.tempUnschedIDs)
	require.Len(t, accountRepo.tempReasons, 1)
	require.Contains(t, accountRepo.tempReasons[0], "probe_failed")
	require.NotNil(t, accountRepo.accounts[301].TempUnschedulableUntil)
	require.True(t, accountRepo.accounts[301].TempUnschedulableUntil.After(now.Add(time.Minute)))
}

func TestUpstreamPoolRecoveryProbeRunnerService_UsesPolicyFailedProbeExtension(t *testing.T) {
	now := time.Now()
	accountRepo := &upstreamPoolRecoveryProbeAccountRepoStub{
		accounts: map[int64]*Account{
			302: {
				ID:                      302,
				Status:                  StatusActive,
				Schedulable:             true,
				Platform:                PlatformOpenAI,
				TempUnschedulableUntil:  ptrTimeUpstreamPoolRecoveryProbe(now.Add(30 * time.Second)),
				TempUnschedulableReason: `{"matched_keyword":"pool_mode_5xx"}`,
				LastUsedAt:              ptrTimeUpstreamPoolRecoveryProbe(now.Add(-10 * time.Minute)),
			},
		},
	}
	upstreamRepo := &upstreamPoolRecoveryProbeRepoStub{
		pools: []UpstreamPool{{
			ID:       1,
			Enabled:  true,
			Platform: PlatformOpenAI,
			PolicyJSON: map[string]any{
				"circuit_breaker": map[string]any{
					"half_open_probe_failed_extension_seconds": 45,
				},
			},
		}},
		members: map[int64][]UpstreamPoolMember{
			1: {{AccountID: 302, Enabled: true}},
		},
	}
	testSvc := &upstreamPoolRecoveryProbeTestSvc{results: map[int64]string{302: "failed"}}
	svc := NewUpstreamPoolRecoveryProbeRunnerService(upstreamRepo, accountRepo, testSvc, &upstreamPoolRecoveryProbeRateLimitStub{}, &config.Config{})

	svc.runScheduled()

	require.Equal(t, []int64{302}, accountRepo.tempUnschedIDs)
	require.NotNil(t, accountRepo.accounts[302].TempUnschedulableUntil)
	require.True(t, accountRepo.accounts[302].TempUnschedulableUntil.After(now.Add(40*time.Second)))
}

func ptrTimeUpstreamPoolRecoveryProbe(t time.Time) *time.Time { return &t }

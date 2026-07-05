//go:build unit

package service

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/Wei-Shaw/sub2api/internal/pkg/openai_compat"
	"github.com/Wei-Shaw/sub2api/internal/pkg/pagination"
	"github.com/stretchr/testify/require"
)

var _ AccountRepository = (*upstreamPoolRecoveryProbeAccountRepoStub)(nil)
var _ UpstreamPoolRepository = (*upstreamPoolRecoveryProbeRepoStub)(nil)

type upstreamPoolRecoveryProbeAccountRepoStub struct {
	accounts       map[int64]*Account
	groupsByID     map[int64][]Group
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
	if r == nil || r.accounts == nil {
		return nil, nil
	}
	out := make([]*Account, 0, len(ids))
	for _, id := range ids {
		if account, ok := r.accounts[id]; ok {
			out = append(out, account)
		}
	}
	return out, nil
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
	if r != nil && r.groupsByID != nil {
		return append([]Group(nil), r.groupsByID[accountID]...), nil
	}
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
	for _, id := range ids {
		account := r.accounts[id]
		if account == nil {
			continue
		}
		if updates.Name != nil {
			account.Name = *updates.Name
		}
		if updates.Status != nil {
			account.Status = *updates.Status
		}
		if updates.Schedulable != nil {
			account.Schedulable = *updates.Schedulable
		}
	}
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
	pools          []UpstreamPool
	members        map[int64][]UpstreamPoolMember
	bindings       []UpstreamPoolBinding
	createdMembers []UpstreamPoolMember
	updatedMembers []UpstreamPoolMember
	deletedMembers []int64
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
	if input != nil {
		r.createdMembers = append(r.createdMembers, *input)
	}
	return input, nil
}

func (r *upstreamPoolRecoveryProbeRepoStub) UpdateUpstreamPoolMember(ctx context.Context, input *UpstreamPoolMember) (*UpstreamPoolMember, error) {
	if input != nil {
		r.updatedMembers = append(r.updatedMembers, *input)
	}
	return input, nil
}

func (r *upstreamPoolRecoveryProbeRepoStub) DeleteUpstreamPoolMember(ctx context.Context, id int64) error {
	r.deletedMembers = append(r.deletedMembers, id)
	for poolID, members := range r.members {
		filtered := members[:0]
		for _, member := range members {
			if member.ID != id {
				filtered = append(filtered, member)
			}
		}
		r.members[poolID] = filtered
	}
	return nil
}

func (r *upstreamPoolRecoveryProbeRepoStub) ListUpstreamAccountSets(ctx context.Context) ([]UpstreamAccountSet, error) {
	return nil, nil
}

func (r *upstreamPoolRecoveryProbeRepoStub) GetUpstreamAccountSetByID(ctx context.Context, id int64) (*UpstreamAccountSet, error) {
	return nil, ErrUpstreamPoolNotFound
}

func (r *upstreamPoolRecoveryProbeRepoStub) CreateUpstreamAccountSet(ctx context.Context, input *UpstreamAccountSet) (*UpstreamAccountSet, error) {
	return input, nil
}

func (r *upstreamPoolRecoveryProbeRepoStub) UpdateUpstreamAccountSet(ctx context.Context, input *UpstreamAccountSet) (*UpstreamAccountSet, error) {
	return input, nil
}

func (r *upstreamPoolRecoveryProbeRepoStub) DeleteUpstreamAccountSet(ctx context.Context, id int64) error {
	return nil
}

func (r *upstreamPoolRecoveryProbeRepoStub) ListUpstreamAccountSetMembers(ctx context.Context, setID int64) ([]UpstreamAccountSetMember, error) {
	return nil, nil
}

func (r *upstreamPoolRecoveryProbeRepoStub) AddUpstreamAccountSetMembers(ctx context.Context, setID int64, accountIDs []int64) error {
	return nil
}

func (r *upstreamPoolRecoveryProbeRepoStub) DeleteUpstreamAccountSetMember(ctx context.Context, setID, accountID int64) error {
	return nil
}

func (r *upstreamPoolRecoveryProbeRepoStub) ListUpstreamPoolMemberSets(ctx context.Context, poolID int64) ([]UpstreamPoolMemberSet, error) {
	return nil, nil
}

func (r *upstreamPoolRecoveryProbeRepoStub) GetUpstreamPoolMemberSetByID(ctx context.Context, id int64) (*UpstreamPoolMemberSet, error) {
	return nil, ErrUpstreamPoolNotFound
}

func (r *upstreamPoolRecoveryProbeRepoStub) CreateUpstreamPoolMemberSet(ctx context.Context, input *UpstreamPoolMemberSet) (*UpstreamPoolMemberSet, error) {
	return input, nil
}

func (r *upstreamPoolRecoveryProbeRepoStub) UpdateUpstreamPoolMemberSet(ctx context.Context, input *UpstreamPoolMemberSet) (*UpstreamPoolMemberSet, error) {
	return input, nil
}

func (r *upstreamPoolRecoveryProbeRepoStub) DeleteUpstreamPoolMemberSet(ctx context.Context, id int64) error {
	return nil
}

func (r *upstreamPoolRecoveryProbeRepoStub) ListUpstreamPoolBindings(ctx context.Context) ([]UpstreamPoolBinding, error) {
	return append([]UpstreamPoolBinding(nil), r.bindings...), nil
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

type upstreamPoolRecoveryProbeHTTPHandler struct {
	mu       sync.Mutex
	requests []map[string]any
}

func (h *upstreamPoolRecoveryProbeHTTPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() { _ = r.Body.Close() }()
	var body map[string]any
	_ = json.NewDecoder(r.Body).Decode(&body)

	h.mu.Lock()
	h.requests = append(h.requests, body)
	h.mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if r.URL.Path == providerOpenAIResponsesPath {
		_ = json.NewEncoder(w).Encode(map[string]any{
			"output": []map[string]any{
				{"type": "message", "content": []map[string]any{{"type": "output_text", "text": "Hi"}}},
			},
		})
		return
	}
	_ = json.NewEncoder(w).Encode(map[string]any{
		"choices": []map[string]any{
			{"message": map[string]any{"content": "Hi"}},
		},
	})
}

func (h *upstreamPoolRecoveryProbeHTTPHandler) bodies() []map[string]any {
	h.mu.Lock()
	defer h.mu.Unlock()
	out := make([]map[string]any, len(h.requests))
	copy(out, h.requests)
	return out
}

func TestUpstreamPoolRecoveryProbeRunnerService_ProbesRecoverablePoolMembers(t *testing.T) {
	swapMonitorHTTPClient(t)
	httpHandler := &upstreamPoolRecoveryProbeHTTPHandler{}
	httpServer := httptest.NewServer(httpHandler)
	t.Cleanup(httpServer.Close)

	accountRepo := &upstreamPoolRecoveryProbeAccountRepoStub{
		accounts: map[int64]*Account{
			101: {ID: 101, Status: StatusError, Schedulable: true, Platform: PlatformOpenAI, Type: AccountTypeAPIKey, Credentials: map[string]any{"api_key": "sk-test", "base_url": httpServer.URL}},
			102: {ID: 102, Status: StatusActive, Schedulable: true, Platform: PlatformOpenAI, Type: AccountTypeAPIKey, Credentials: map[string]any{"api_key": "sk-test", "base_url": httpServer.URL}},
			103: {ID: 103, Status: StatusActive, Schedulable: true, Platform: PlatformOpenAI, Type: AccountTypeAPIKey, Credentials: map[string]any{"api_key": "sk-test", "base_url": httpServer.URL}, RateLimitResetAt: ptrTimeUpstreamPoolRecoveryProbe(time.Now().Add(10 * time.Minute))},
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

	require.Empty(t, testSvc.calls)
	require.Empty(t, testSvc.modes)
	require.ElementsMatch(t, []int64{101, 103}, rateLimitSvc.calls)

	bodies := httpHandler.bodies()
	require.Len(t, bodies, 2)
	for _, body := range bodies {
		messages, _ := body["messages"].([]any)
		require.Len(t, messages, 1)
		msg, _ := messages[0].(map[string]any)
		require.Equal(t, monitorLightweightPrompt, strings.TrimSpace(stringFromAny(msg["content"])))
		require.Equal(t, float64(monitorLightweightMaxTokens), body["max_tokens"])
	}
}

func TestUpstreamPoolRecoveryProbeRunnerService_ProbesAnthropicPoolMembersViaAccountTest(t *testing.T) {
	accountRepo := &upstreamPoolRecoveryProbeAccountRepoStub{
		accounts: map[int64]*Account{
			401: {
				ID:          401,
				Status:      StatusError,
				Schedulable: true,
				Platform:    PlatformAnthropic,
				Type:        AccountTypeOAuth,
				LastUsedAt:  ptrTimeUpstreamPoolRecoveryProbe(time.Now().Add(-10 * time.Minute)),
			},
		},
	}
	upstreamRepo := &upstreamPoolRecoveryProbeRepoStub{
		pools: []UpstreamPool{{ID: 4, Enabled: true, Platform: PlatformAnthropic}},
		members: map[int64][]UpstreamPoolMember{
			4: {{AccountID: 401, Enabled: true}},
		},
	}
	testSvc := &upstreamPoolRecoveryProbeTestSvc{results: map[int64]string{401: "success"}}
	rateLimitSvc := &upstreamPoolRecoveryProbeRateLimitStub{}
	svc := NewUpstreamPoolRecoveryProbeRunnerService(upstreamRepo, accountRepo, testSvc, rateLimitSvc, &config.Config{})

	svc.runScheduled()

	require.Equal(t, []int64{401}, testSvc.calls)
	require.Equal(t, []int64{401}, rateLimitSvc.calls)
}

func TestAdminServiceSyncUpstreamPoolMembersForAccount_AnthropicOAuth(t *testing.T) {
	ctx := context.Background()
	account := &Account{
		ID:          501,
		Name:        "claude-oauth",
		Platform:    PlatformAnthropic,
		Type:        AccountTypeOAuth,
		Status:      StatusActive,
		Schedulable: true,
		GroupIDs:    []int64{42},
	}
	accountRepo := &upstreamPoolRecoveryProbeAccountRepoStub{
		accounts: map[int64]*Account{501: account},
	}
	upstreamRepo := &upstreamPoolRecoveryProbeRepoStub{
		bindings: []UpstreamPoolBinding{
			{ID: 1, GroupID: 42, PoolID: 7, Platform: PlatformAnthropic, Enabled: true},
			{ID: 2, GroupID: 42, PoolID: 8, Platform: PlatformOpenAI, Enabled: true},
		},
		members: map[int64][]UpstreamPoolMember{
			7: {},
			8: {},
		},
	}
	svc := &adminServiceImpl{
		accountRepo:      accountRepo,
		upstreamPoolRepo: upstreamRepo,
	}

	require.NoError(t, svc.syncUpstreamPoolMembersForAccount(ctx, nil, account))

	require.Len(t, upstreamRepo.createdMembers, 1)
	require.Equal(t, int64(7), upstreamRepo.createdMembers[0].PoolID)
	require.Equal(t, int64(501), upstreamRepo.createdMembers[0].AccountID)
	require.Equal(t, PlatformAnthropic, upstreamRepo.createdMembers[0].AccountPlatform)
	require.True(t, upstreamRepo.createdMembers[0].Enabled)
	require.False(t, upstreamRepo.createdMembers[0].ManualDrained)
	require.Empty(t, upstreamRepo.deletedMembers)
}

func TestAdminServiceSyncUpstreamPoolMembersForAccount_RemovesPreviousPlatformPoolMember(t *testing.T) {
	ctx := context.Background()
	previous := &Account{
		ID:          502,
		Name:        "openai-key",
		Platform:    PlatformOpenAI,
		Type:        AccountTypeAPIKey,
		Status:      StatusActive,
		Schedulable: true,
		GroupIDs:    []int64{42},
	}
	current := &Account{
		ID:          502,
		Name:        "claude-oauth",
		Platform:    PlatformAnthropic,
		Type:        AccountTypeOAuth,
		Status:      StatusActive,
		Schedulable: true,
		GroupIDs:    []int64{43},
	}
	upstreamRepo := &upstreamPoolRecoveryProbeRepoStub{
		bindings: []UpstreamPoolBinding{
			{ID: 1, GroupID: 42, PoolID: 7, Platform: PlatformOpenAI, Enabled: true},
			{ID: 2, GroupID: 43, PoolID: 8, Platform: PlatformAnthropic, Enabled: true},
		},
		members: map[int64][]UpstreamPoolMember{
			7: {{ID: 70, PoolID: 7, AccountID: 502, AccountPlatform: PlatformOpenAI}},
			8: {},
		},
	}
	svc := &adminServiceImpl{
		accountRepo:      &upstreamPoolRecoveryProbeAccountRepoStub{},
		upstreamPoolRepo: upstreamRepo,
	}

	require.NoError(t, svc.syncUpstreamPoolMembersForAccount(ctx, previous, current))

	require.Equal(t, []int64{70}, upstreamRepo.deletedMembers)
	require.Len(t, upstreamRepo.createdMembers, 1)
	require.Equal(t, int64(8), upstreamRepo.createdMembers[0].PoolID)
	require.Equal(t, PlatformAnthropic, upstreamRepo.createdMembers[0].AccountPlatform)
}

func TestAdminServiceBulkUpdateAccounts_SyncsAnthropicUpstreamPoolMembers(t *testing.T) {
	ctx := context.Background()
	groupIDs := []int64{42}
	schedulable := false
	accountRepo := &upstreamPoolRecoveryProbeAccountRepoStub{
		accounts: map[int64]*Account{
			503: {
				ID:          503,
				Name:        "claude-oauth",
				Platform:    PlatformAnthropic,
				Type:        AccountTypeOAuth,
				Status:      StatusActive,
				Schedulable: true,
			},
		},
	}
	upstreamRepo := &upstreamPoolRecoveryProbeRepoStub{
		bindings: []UpstreamPoolBinding{
			{ID: 1, GroupID: 42, PoolID: 9, Platform: PlatformAnthropic, Enabled: true},
		},
		members: map[int64][]UpstreamPoolMember{
			9: {{ID: 90, PoolID: 9, AccountID: 503, AccountPlatform: PlatformAnthropic, Enabled: true}},
		},
	}
	svc := &adminServiceImpl{
		accountRepo:      accountRepo,
		groupRepo:        &groupRepoStubForAdmin{getByID: &Group{ID: 42, Name: "claude"}},
		upstreamPoolRepo: upstreamRepo,
	}

	result, err := svc.BulkUpdateAccounts(ctx, &BulkUpdateAccountsInput{
		AccountIDs:            []int64{503},
		GroupIDs:              &groupIDs,
		Schedulable:           &schedulable,
		SkipMixedChannelCheck: true,
	})

	require.NoError(t, err)
	require.Equal(t, 1, result.Success)
	require.Empty(t, result.FailedIDs)
	require.Len(t, upstreamRepo.updatedMembers, 1)
	require.Equal(t, int64(9), upstreamRepo.updatedMembers[0].PoolID)
	require.False(t, upstreamRepo.updatedMembers[0].Enabled)
	require.True(t, upstreamRepo.updatedMembers[0].ManualDrained)
}

func TestUpstreamPoolRecoveryProbeRunnerService_ResponsesProbeUsesOnlyHiAndOneToken(t *testing.T) {
	swapMonitorHTTPClient(t)
	httpHandler := &upstreamPoolRecoveryProbeHTTPHandler{}
	httpServer := httptest.NewServer(httpHandler)
	t.Cleanup(httpServer.Close)

	accountRepo := &upstreamPoolRecoveryProbeAccountRepoStub{
		accounts: map[int64]*Account{
			111: {
				ID:           111,
				Status:       StatusError,
				Schedulable:  true,
				Platform:     PlatformOpenAI,
				Type:         AccountTypeAPIKey,
				Credentials:  map[string]any{"api_key": "sk-test", "base_url": httpServer.URL},
				Extra:        map[string]any{openai_compat.ExtraKeyResponsesSupported: true},
				ErrorMessage: "runtime error",
			},
		},
	}
	upstreamRepo := &upstreamPoolRecoveryProbeRepoStub{
		pools:   []UpstreamPool{{ID: 1, Enabled: true, Platform: PlatformOpenAI}},
		members: map[int64][]UpstreamPoolMember{1: {{AccountID: 111, Enabled: true}}},
	}
	testSvc := &upstreamPoolRecoveryProbeTestSvc{}
	rateLimitSvc := &upstreamPoolRecoveryProbeRateLimitStub{}
	svc := NewUpstreamPoolRecoveryProbeRunnerService(upstreamRepo, accountRepo, testSvc, rateLimitSvc, &config.Config{})

	svc.runScheduled()

	require.Empty(t, testSvc.calls)
	require.Equal(t, []int64{111}, rateLimitSvc.calls)

	bodies := httpHandler.bodies()
	require.Len(t, bodies, 1)
	require.Equal(t, monitorLightweightPrompt, strings.TrimSpace(stringFromAny(bodies[0]["input"])))
	require.Equal(t, float64(monitorLightweightMaxTokens), bodies[0]["max_output_tokens"])
	require.NotContains(t, bodies[0], "instructions")
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
				Type:                    AccountTypeAPIKey,
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

	require.Empty(t, testSvc.calls)
	require.Empty(t, testSvc.modes)
	require.Empty(t, rateLimitSvc.calls)
	require.Equal(t, []int64{301}, accountRepo.tempUnschedIDs)
	require.Len(t, accountRepo.tempReasons, 1)
	require.Contains(t, accountRepo.tempReasons[0], "probe_error")
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
				Type:                    AccountTypeAPIKey,
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

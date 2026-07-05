package service

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"sync"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/pkg/claude"
	"github.com/Wei-Shaw/sub2api/internal/pkg/geminicli"
	openaipkg "github.com/Wei-Shaw/sub2api/internal/pkg/openai"
	"github.com/Wei-Shaw/sub2api/internal/pkg/openai_compat"
	"github.com/alitto/pond/v2"
	"golang.org/x/sync/errgroup"
)

type AccountMonitorRepository interface {
	InsertHistoryBatch(ctx context.Context, rows []*AccountMonitorHistoryRow) error
	ListLatestForAccountIDs(ctx context.Context, ids []int64) (map[int64][]*AccountMonitorLatest, error)
	ComputeAvailabilityForAccounts(ctx context.Context, ids []int64, windowDays int) (map[int64][]*AccountMonitorAvailability, error)
	// ListHistorySinceForAccounts loads recent probe facts for account IDs. When primaryModels is empty,
	// all models are returned; pool health then filters by each pool's own probe model.
	ListHistorySinceForAccounts(ctx context.Context, ids []int64, primaryModels map[int64]string, since time.Time) (map[int64][]*AccountMonitorHistoryEntry, error)
	DeleteHistoryBefore(ctx context.Context, before time.Time) (int64, error)
}

type AccountMonitorService struct {
	repo             AccountMonitorRepository
	upstreamPoolRepo UpstreamPoolRepository
	accountRepo      AccountRepository
}

type accountProbeTarget struct {
	AccountID  int64
	PoolID     int64
	GroupID    int64
	Provider   string
	Model      string
	Endpoint   string
	APIKey     string
	APIMode    string
	AuthScheme string
}

func NewAccountMonitorService(
	repo AccountMonitorRepository,
	upstreamPoolRepo UpstreamPoolRepository,
	accountRepo AccountRepository,
) *AccountMonitorService {
	return &AccountMonitorService{repo: repo, upstreamPoolRepo: upstreamPoolRepo, accountRepo: accountRepo}
}

func (s *AccountMonitorService) RunOnce(ctx context.Context) error {
	if s == nil || s.repo == nil || s.upstreamPoolRepo == nil || s.accountRepo == nil {
		return nil
	}
	targets, err := s.buildProbeTargets(ctx)
	if err != nil {
		return err
	}
	if len(targets) == 0 {
		return nil
	}
	rows := s.runProbeTargets(ctx, targets)
	if err := s.repo.InsertHistoryBatch(ctx, rows); err != nil {
		return fmt.Errorf("insert account monitor history: %w", err)
	}
	return nil
}

func (s *AccountMonitorService) buildProbeTargets(ctx context.Context) ([]accountProbeTarget, error) {
	pools, err := s.upstreamPoolRepo.ListUpstreamPools(ctx)
	if err != nil {
		return nil, fmt.Errorf("list upstream pools: %w", err)
	}
	bindings, err := s.upstreamPoolRepo.ListUpstreamPoolBindings(ctx)
	if err != nil {
		return nil, fmt.Errorf("list upstream pool bindings: %w", err)
	}
	enabledPools := filterEnabledPools(pools)
	poolBindings := groupPoolBindings(bindings, enabledPools)
	membersByPool, accountsByID := s.loadPoolMembersAndAccounts(ctx, enabledPools)

	targets := make([]accountProbeTarget, 0)
	seen := make(map[string]struct{})
	for _, pool := range enabledPools {
		model := probeModelForPool(pool.Platform, poolBindings[pool.ID])
		if model == "" {
			continue
		}
		groupID, _ := primaryPoolHealthGroup(poolBindings[pool.ID])
		for _, member := range membersByPool[pool.ID] {
			account := accountsByID[member.AccountID]
			if !accountProbeCandidate(member, account, pool.Platform) {
				continue
			}
			endpoint, apiKey := accountProbeEndpointAndKey(account)
			if endpoint == "" || apiKey == "" {
				continue
			}
			key := fmt.Sprintf("%d:%d:%s", pool.ID, account.ID, model)
			if _, ok := seen[key]; ok {
				continue
			}
			seen[key] = struct{}{}
			target := accountProbeTarget{
				AccountID:  account.ID,
				PoolID:     pool.ID,
				Provider:   pool.Platform,
				Model:      account.GetMappedModel(model),
				Endpoint:   endpoint,
				APIKey:     apiKey,
				APIMode:    accountProbeAPIMode(account),
				AuthScheme: accountProbeAuthScheme(account),
			}
			if groupID != nil {
				target.GroupID = *groupID
			}
			targets = append(targets, target)
		}
	}
	return targets, nil
}

func (s *AccountMonitorService) loadPoolMembersAndAccounts(ctx context.Context, pools []UpstreamPool) (map[int64][]UpstreamPoolMember, map[int64]*Account) {
	helper := &PoolHealthService{upstreamPoolRepo: s.upstreamPoolRepo, accountRepo: s.accountRepo}
	return helper.loadPoolMembersAndAccounts(ctx, pools)
}

func (s *AccountMonitorService) runProbeTargets(ctx context.Context, targets []accountProbeTarget) []*AccountMonitorHistoryRow {
	rows := make([]*AccountMonitorHistoryRow, 0, len(targets))
	var mu sync.Mutex
	var eg errgroup.Group
	sem := make(chan struct{}, monitorWorkerConcurrency)
	for _, target := range targets {
		target := target
		eg.Go(func() error {
			select {
			case sem <- struct{}{}:
				defer func() { <-sem }()
			case <-ctx.Done():
				return nil
			}
			row := runAccountProbeTarget(ctx, target)
			mu.Lock()
			rows = append(rows, row)
			mu.Unlock()
			return nil
		})
	}
	_ = eg.Wait()
	return rows
}

func runAccountProbeTarget(ctx context.Context, target accountProbeTarget) *AccountMonitorHistoryRow {
	pingMs := pingEndpointOrigin(ctx, target.Endpoint)
	result := runCheckForModel(ctx, target.Provider, target.Endpoint, target.APIKey, target.Model, &CheckOptions{
		Lightweight:      true,
		APIMode:          target.APIMode,
		APIKeyAuthScheme: target.AuthScheme,
	})
	result.PingLatencyMs = pingMs
	poolID := target.PoolID
	var groupID *int64
	if target.GroupID > 0 {
		id := target.GroupID
		groupID = &id
	}
	return &AccountMonitorHistoryRow{
		AccountID:     target.AccountID,
		PoolID:        &poolID,
		GroupID:       groupID,
		Provider:      target.Provider,
		Model:         result.Model,
		Status:        result.Status,
		LatencyMs:     result.LatencyMs,
		PingLatencyMs: result.PingLatencyMs,
		Message:       result.Message,
		CheckedAt:     result.CheckedAt,
	}
}

func accountProbeAPIMode(account *Account) string {
	if account == nil || account.Platform != PlatformOpenAI || account.Type != AccountTypeAPIKey {
		return ""
	}
	if openai_compat.ResolveResponsesSupport(account.Extra) == openai_compat.ResponsesSupportYes {
		return MonitorAPIModeResponses
	}
	return MonitorAPIModeChatCompletions
}

func accountProbeAuthScheme(account *Account) string {
	if account == nil || account.Platform != PlatformAnthropic || account.Type != AccountTypeAPIKey {
		return ""
	}
	return account.GetAnthropicAPIKeyAuthScheme()
}

func accountProbeCandidate(member UpstreamPoolMember, account *Account, poolPlatform string) bool {
	if !poolMemberSchedulable(member, account) {
		return false
	}
	if account == nil || account.Type != AccountTypeAPIKey {
		return false
	}
	if !strings.EqualFold(account.Platform, poolPlatform) {
		return false
	}
	switch account.Platform {
	case PlatformOpenAI, PlatformAnthropic, PlatformGemini:
		return true
	default:
		return false
	}
}

func accountProbeEndpointAndKey(account *Account) (string, string) {
	if account == nil {
		return "", ""
	}
	switch account.Platform {
	case PlatformOpenAI:
		return account.GetOpenAIBaseURL(), account.GetOpenAIApiKey()
	case PlatformAnthropic:
		return account.GetBaseURL(), account.GetCredential("api_key")
	case PlatformGemini:
		return account.GetGeminiBaseURL(geminicli.AIStudioBaseURL), account.GetCredential("api_key")
	default:
		return "", ""
	}
}

func probeModelForPool(provider string, bindings []UpstreamPoolBinding) string {
	for _, binding := range bindings {
		for _, model := range binding.Models {
			if isConcreteProbeModel(model) {
				return strings.TrimSpace(model)
			}
		}
	}
	switch provider {
	case PlatformOpenAI:
		return openaipkg.DefaultTestModel
	case PlatformAnthropic:
		return claude.DefaultTestModel
	case PlatformGemini:
		return geminicli.DefaultTestModel
	default:
		return ""
	}
}

func isConcreteProbeModel(model string) bool {
	model = strings.TrimSpace(model)
	if model == "" || model == "*" {
		return false
	}
	return !strings.Contains(model, "*")
}

type AccountMonitorRunner struct {
	svc            *AccountMonitorService
	settingService *SettingService
	pool           pond.Pool
	parentCtx      context.Context
	parentCancel   context.CancelFunc
	wg             sync.WaitGroup
}

func NewAccountMonitorRunner(svc *AccountMonitorService, settingService *SettingService) *AccountMonitorRunner {
	ctx, cancel := context.WithCancel(context.Background())
	return &AccountMonitorRunner{
		svc:            svc,
		settingService: settingService,
		pool:           pond.NewPool(1),
		parentCtx:      ctx,
		parentCancel:   cancel,
	}
}

func (r *AccountMonitorRunner) Start() {
	if r == nil || r.svc == nil {
		return
	}
	r.wg.Add(1)
	go r.loop()
}

func (r *AccountMonitorRunner) Stop() {
	if r == nil {
		return
	}
	r.parentCancel()
	r.wg.Wait()
	r.pool.StopAndWait()
}

func (r *AccountMonitorRunner) loop() {
	defer r.wg.Done()
	r.fire()
	timer := time.NewTimer(r.interval())
	defer timer.Stop()
	for {
		select {
		case <-r.parentCtx.Done():
			return
		case <-timer.C:
			r.fire()
			timer.Reset(r.interval())
		}
	}
}

func (r *AccountMonitorRunner) interval() time.Duration {
	if r.settingService == nil {
		return 5 * time.Minute
	}
	runtime := r.settingService.GetChannelMonitorRuntime(context.Background())
	if runtime.DefaultIntervalSeconds <= 0 {
		return 5 * time.Minute
	}
	return time.Duration(runtime.DefaultIntervalSeconds) * time.Second
}

func (r *AccountMonitorRunner) fire() {
	if r.settingService != nil && !r.settingService.GetChannelMonitorRuntime(r.parentCtx).Enabled {
		return
	}
	if _, ok := r.pool.TrySubmit(func() {
		ctx, cancel := context.WithTimeout(r.parentCtx, accountMonitorRunTimeout())
		defer cancel()
		if err := r.svc.RunOnce(ctx); err != nil {
			slog.Warn("account_monitor: run failed", "error", err)
		}
	}); !ok {
		slog.Debug("account_monitor: previous run still active, skip")
	}
}

func accountMonitorRunTimeout() time.Duration {
	return monitorRequestTimeout + monitorPingTimeout + monitorRunOneBuffer
}

func ProvideAccountMonitorRunner(svc *AccountMonitorService, settingService *SettingService) *AccountMonitorRunner {
	r := NewAccountMonitorRunner(svc, settingService)
	r.Start()
	return r
}

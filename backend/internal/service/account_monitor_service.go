package service

import (
	"context"
	"fmt"
	"log/slog"
	"math"
	"sort"
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

const (
	poolAvailabilitySnapshotRetention     = 30 * 24 * time.Hour
	poolAvailabilitySnapshotPruneInterval = 24 * time.Hour
	poolRuntimeWeightTTL                  = 30 * time.Minute
)

type AccountMonitorRepository interface {
	InsertHistoryBatch(ctx context.Context, rows []*AccountMonitorHistoryRow) error
	InsertPoolAvailabilitySnapshots(ctx context.Context, rows []*PoolAvailabilitySnapshotRow) error
	ListPoolAvailabilitySince(ctx context.Context, poolIDs []int64, since time.Time) (map[int64][]*PoolAvailabilitySnapshotEntry, error)
	ListLatestForAccountIDs(ctx context.Context, ids []int64) (map[int64][]*AccountMonitorLatest, error)
	ComputeAvailabilityForAccounts(ctx context.Context, ids []int64, windowDays int) (map[int64][]*AccountMonitorAvailability, error)
	// ListHistorySinceForAccounts loads recent probe facts for account IDs. When primaryModels is empty,
	// all models are returned; pool health then filters by each pool's own probe model.
	ListHistorySinceForAccounts(ctx context.Context, ids []int64, primaryModels map[int64]string, since time.Time) (map[int64][]*AccountMonitorHistoryEntry, error)
	DeleteHistoryBefore(ctx context.Context, before time.Time) (int64, error)
	DeletePoolAvailabilityBefore(ctx context.Context, before time.Time) (int64, error)
	ListPoolRuntimeWeightStates(ctx context.Context, poolIDs []int64) (map[int64]map[int64]*PoolRuntimeWeightState, error)
	UpsertPoolRuntimeWeightStates(ctx context.Context, states []*PoolRuntimeWeightState) error
	ListRecentAccountRuntimeHealth(ctx context.Context, accountIDs []int64, since time.Time) (map[int64]AccountRuntimeHealthSnapshot, error)
	InsertUpstreamCapacitySnapshots(ctx context.Context, rows []*UpstreamCapacitySnapshotRow) error
	ListUpstreamCapacitySnapshotStats(ctx context.Context, setIDs []int64, since time.Time) (map[int64]UpstreamCapacitySnapshotStats, error)
}

type AccountMonitorService struct {
	repo                 AccountMonitorRepository
	upstreamPoolRepo     UpstreamPoolRepository
	accountRepo          AccountRepository
	pruneMu              sync.Mutex
	nextSnapshotPruneAt  time.Time
	healthAlertMu        sync.Mutex
	healthAlertStates    map[string]*upstreamHealthAlertState
	healthAlertEmitter   func(upstreamHealthAlertEvent)
	now                  func() time.Time
	capacityReader       UpstreamCapacityPressureReader
	gatewayService       *GatewayService
	openAIGatewayService *OpenAIGatewayService
}

func (s *AccountMonitorService) SetCapacityPressureReader(reader UpstreamCapacityPressureReader) {
	if s != nil {
		s.capacityReader = reader
	}
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
	gatewayService *GatewayService,
	openAIGatewayService *OpenAIGatewayService,
) *AccountMonitorService {
	return &AccountMonitorService{
		repo: repo, upstreamPoolRepo: upstreamPoolRepo, accountRepo: accountRepo,
		gatewayService: gatewayService, openAIGatewayService: openAIGatewayService,
		healthAlertStates:  make(map[string]*upstreamHealthAlertState),
		healthAlertEmitter: emitUpstreamHealthAlert,
		now:                time.Now,
	}
}

func (s *AccountMonitorService) RunOnce(ctx context.Context) error {
	if s == nil || s.repo == nil || s.upstreamPoolRepo == nil || s.accountRepo == nil {
		return nil
	}
	snapshots, err := s.buildPoolAvailabilitySnapshots(ctx)
	if err != nil {
		return err
	}
	if err := s.repo.InsertPoolAvailabilitySnapshots(ctx, snapshots); err != nil {
		return fmt.Errorf("insert pool availability snapshots: %w", err)
	}
	if s.capacityReader != nil {
		pressures, pressureErr := s.capacityReader.ListUpstreamCapacityPressures(ctx)
		if pressureErr != nil {
			return fmt.Errorf("list upstream capacity pressures: %w", pressureErr)
		}
		rows := make([]*UpstreamCapacitySnapshotRow, 0, len(pressures))
		for _, pressure := range pressures {
			loadRate := 0
			if pressure.CapacityLimit > 0 {
				loadRate = (pressure.CurrentConcurrency + pressure.WaitingCount) * 100 / pressure.CapacityLimit
			}
			rows = append(rows, &UpstreamCapacitySnapshotRow{SetID: pressure.SetID, CapacityLimit: pressure.CapacityLimit, CurrentConcurrency: pressure.CurrentConcurrency, WaitingCount: pressure.WaitingCount, LoadRate: loadRate, CheckedAt: s.currentTime()})
		}
		if err := s.repo.InsertUpstreamCapacitySnapshots(ctx, rows); err != nil {
			return fmt.Errorf("insert upstream capacity snapshots: %w", err)
		}
	}
	if _, err := s.prunePoolAvailabilitySnapshots(ctx); err != nil {
		slog.Warn("account_monitor: prune pool availability snapshots failed", "error", err)
	}
	targets, err := s.buildProbeTargets(ctx)
	if err != nil {
		return err
	}
	var rows []*AccountMonitorHistoryRow
	if len(targets) > 0 {
		rows = s.runProbeTargets(ctx, targets)
		if err := s.repo.InsertHistoryBatch(ctx, rows); err != nil {
			return fmt.Errorf("insert account monitor history: %w", err)
		}
	}
	if err := s.adjustPoolRuntimeWeights(ctx, rows); err != nil {
		slog.Warn("account_monitor: adjust pool runtime weights failed", "error", err)
	}
	if err := s.evaluateUpstreamHealthAlerts(ctx, rows); err != nil {
		slog.Warn("account_monitor: evaluate upstream health alerts failed", "error", err)
	}
	return nil
}

func (s *AccountMonitorService) adjustPoolRuntimeWeights(ctx context.Context, probeRows []*AccountMonitorHistoryRow) error {
	pools, err := s.upstreamPoolRepo.ListUpstreamPools(ctx)
	if err != nil {
		return fmt.Errorf("list upstream pools for auto weight: %w", err)
	}
	autoPools := make([]UpstreamPool, 0, len(pools))
	poolIDs := make([]int64, 0, len(pools))
	for _, pool := range pools {
		if pool.Enabled && UpstreamPoolAutoWeightModeFromPolicyJSON(pool.PolicyJSON) != "off" && strings.EqualFold(pool.Platform, PlatformOpenAI) {
			autoPools = append(autoPools, pool)
			poolIDs = append(poolIDs, pool.ID)
		}
	}
	if len(autoPools) == 0 {
		return nil
	}
	statesByPool, err := s.repo.ListPoolRuntimeWeightStates(ctx, poolIDs)
	if err != nil {
		return err
	}
	membersByPool, accountsByID := s.loadPoolMembersAndAccounts(ctx, autoPools)
	accountIDs := make([]int64, 0)
	seenAccountIDs := make(map[int64]struct{})
	for _, members := range membersByPool {
		for _, member := range members {
			if member.AccountID > 0 {
				if _, ok := seenAccountIDs[member.AccountID]; !ok {
					seenAccountIDs[member.AccountID] = struct{}{}
					accountIDs = append(accountIDs, member.AccountID)
				}
			}
		}
	}
	realHealth, err := s.repo.ListRecentAccountRuntimeHealth(ctx, accountIDs, s.currentTime().UTC().Add(-30*time.Minute))
	if err != nil {
		slog.Warn("account_monitor: list recent account runtime health failed; falling back to probes", "error", err)
		realHealth = nil
	}
	rowsByPool := make(map[int64]map[int64]*AccountMonitorHistoryRow, len(autoPools))
	for _, row := range probeRows {
		if row == nil || row.PoolID == nil || *row.PoolID <= 0 || row.AccountID <= 0 {
			continue
		}
		if rowsByPool[*row.PoolID] == nil {
			rowsByPool[*row.PoolID] = map[int64]*AccountMonitorHistoryRow{}
		}
		rowsByPool[*row.PoolID][row.AccountID] = row
	}
	now := s.currentTime().UTC()
	updates := make([]*PoolRuntimeWeightState, 0)
	for _, pool := range autoPools {
		members := activeAutoWeightMembers(membersByPool[pool.ID], accountsByID)
		if len(members) < 2 {
			continue
		}
		medianLatency := medianOperationalProbeLatency(rowsByPool[pool.ID])
		medianRealTTFT := medianRuntimeTTFT(realHealth, members)
		for _, member := range members {
			account := accountsByID[member.AccountID]
			target, reason, observed := autoWeightTargetWithRuntime(account, rowsByPool[pool.ID][member.AccountID], medianLatency, realHealth[member.AccountID], medianRealTTFT, now)
			if !observed {
				continue
			}
			current := (*PoolRuntimeWeightState)(nil)
			if statesByPool[pool.ID] != nil {
				current = statesByPool[pool.ID][member.AccountID]
			}
			if current != nil && now.Sub(current.LastObservedAt) > poolRuntimeWeightTTL {
				current = nil
			}
			updates = append(updates, nextPoolRuntimeWeightState(pool.ID, member.AccountID, current, target, reason, now))
		}
	}
	return s.repo.UpsertPoolRuntimeWeightStates(ctx, updates)
}

func (s *AccountMonitorService) seedGatewayRuntimeHealth(health map[int64]AccountRuntimeHealthSnapshot) {
	for accountID, snapshot := range health {
		if snapshot.SampleCount < 3 || snapshot.P95TTFTMs <= 0 {
			continue
		}
		if s.gatewayService != nil {
			s.gatewayService.SeedAccountRuntimeTTFT(accountID, snapshot.P95TTFTMs)
		}
		if s.openAIGatewayService != nil {
			s.openAIGatewayService.SeedOpenAIAccountRuntimeTTFT(accountID, snapshot.P95TTFTMs)
		}
	}
}

func (s *AccountMonitorService) currentTime() time.Time {
	if s != nil && s.now != nil {
		return s.now()
	}
	return time.Now()
}

func (s *AccountMonitorService) WarmRuntimeHealth(ctx context.Context) error {
	if s == nil || s.repo == nil || s.upstreamPoolRepo == nil || s.accountRepo == nil {
		return nil
	}
	pools, err := s.upstreamPoolRepo.ListUpstreamPools(ctx)
	if err != nil {
		return err
	}
	_, accounts := s.loadPoolMembersAndAccounts(ctx, pools)
	accountIDs := make([]int64, 0, len(accounts))
	for accountID := range accounts {
		accountIDs = append(accountIDs, accountID)
	}
	health, err := s.repo.ListRecentAccountRuntimeHealth(ctx, accountIDs, s.currentTime().UTC().Add(-30*time.Minute))
	if err != nil {
		return err
	}
	s.seedGatewayRuntimeHealth(health)
	return nil
}

func activeAutoWeightMembers(members []UpstreamPoolMember, accounts map[int64]*Account) []UpstreamPoolMember {
	out := make([]UpstreamPoolMember, 0, len(members))
	for _, member := range members {
		account := accounts[member.AccountID]
		if !member.Enabled || member.ManualDrained || (member.SchedulableOverride != nil && !*member.SchedulableOverride) || account == nil || !account.IsActive() {
			continue
		}
		out = append(out, member)
	}
	return out
}

func medianOperationalProbeLatency(rows map[int64]*AccountMonitorHistoryRow) int {
	latencies := make([]int, 0, len(rows))
	for _, row := range rows {
		if row != nil && row.Status == MonitorStatusOperational && row.LatencyMs != nil && *row.LatencyMs > 0 {
			latencies = append(latencies, *row.LatencyMs)
		}
	}
	if len(latencies) == 0 {
		return 0
	}
	sort.Ints(latencies)
	middle := len(latencies) / 2
	if len(latencies)%2 == 1 {
		return latencies[middle]
	}
	return (latencies[middle-1] + latencies[middle]) / 2
}

func autoWeightTarget(account *Account, row *AccountMonitorHistoryRow, medianLatency int, now time.Time) (float64, string, bool) {
	return autoWeightTargetWithRuntime(account, row, medianLatency, AccountRuntimeHealthSnapshot{}, 0, now)
}

func medianRuntimeTTFT(health map[int64]AccountRuntimeHealthSnapshot, members []UpstreamPoolMember) int {
	values := make([]int, 0, len(members))
	for _, member := range members {
		if snapshot, ok := health[member.AccountID]; ok && snapshot.SampleCount >= 3 && snapshot.P95TTFTMs > 0 {
			values = append(values, snapshot.P95TTFTMs)
		}
	}
	if len(values) == 0 {
		return 0
	}
	sort.Ints(values)
	middle := len(values) / 2
	if len(values)%2 == 1 {
		return values[middle]
	}
	return (values[middle-1] + values[middle]) / 2
}

func autoWeightTargetWithRuntime(account *Account, row *AccountMonitorHistoryRow, medianLatency int, runtime AccountRuntimeHealthSnapshot, medianRealTTFT int, now time.Time) (float64, string, bool) {
	if account == nil {
		return 1, "", false
	}
	if account.RateLimitResetAt != nil && now.Before(*account.RateLimitResetAt) {
		return 0.25, "rate_limited", true
	}
	if account.OverloadUntil != nil && now.Before(*account.OverloadUntil) {
		return 0.5, "overloaded", true
	}
	if account.TempUnschedulableUntil != nil && now.Before(*account.TempUnschedulableUntil) {
		return 0.5, "temporarily_unschedulable", true
	}
	if runtime.SampleCount >= 3 && runtime.P95TTFTMs > 0 {
		ratio := 0.0
		if medianRealTTFT > 0 {
			ratio = float64(runtime.P95TTFTMs) / float64(medianRealTTFT)
		}
		switch {
		case runtime.P95TTFTMs > 15000 && (medianRealTTFT <= 0 || ratio >= 2):
			return 0.5, "real_ttft_much_slower", true
		case runtime.P95TTFTMs > 15000 || ratio >= 1.5:
			return 0.75, "real_ttft_slower", true
		case ratio <= 0.75:
			return 1.25, "real_ttft_faster", true
		default:
			return 1, "real_ttft_healthy", true
		}
	}
	if row == nil {
		return 1, "", false
	}
	switch row.Status {
	case MonitorStatusFailed, MonitorStatusError:
		return 0.5, "probe_failed", true
	case MonitorStatusDegraded:
		return 0.75, "probe_degraded", true
	case MonitorStatusOperational:
		if row.LatencyMs == nil || *row.LatencyMs <= 0 || medianLatency <= 0 {
			return 1, "healthy", true
		}
		ratio := float64(*row.LatencyMs) / float64(medianLatency)
		switch {
		case ratio <= 0.75:
			return 1.25, "faster_than_pool", true
		case ratio >= 1.75:
			return 0.5, "much_slower_than_pool", true
		case ratio >= 1.25:
			return 0.75, "slower_than_pool", true
		default:
			return 1, "healthy", true
		}
	default:
		return 1, "", false
	}
}

func nextPoolRuntimeWeightState(poolID, accountID int64, current *PoolRuntimeWeightState, target float64, reason string, now time.Time) *PoolRuntimeWeightState {
	factor := 1.0
	previousTarget := 1.0
	healthyStreak, unhealthyStreak := 0, 0
	if current != nil {
		factor = current.Factor
		previousTarget = current.TargetFactor
		healthyStreak = current.HealthyStreak
		unhealthyStreak = current.UnhealthyStreak
	}
	target = math.Max(0.25, math.Min(1.25, target))
	direction := 0
	if target > factor {
		direction = 1
	} else if target < factor {
		direction = -1
	}
	if direction == 0 {
		healthyStreak, unhealthyStreak = 0, 0
	} else if direction > 0 {
		if previousTarget != target {
			healthyStreak = 0
		}
		healthyStreak++
		unhealthyStreak = 0
		if healthyStreak >= 3 {
			factor = math.Min(target, factor+0.25)
		}
	} else {
		if previousTarget != target {
			unhealthyStreak = 0
		}
		unhealthyStreak++
		healthyStreak = 0
		if unhealthyStreak >= 3 {
			factor = math.Max(target, factor-0.25)
		}
	}
	return &PoolRuntimeWeightState{
		PoolID: poolID, AccountID: accountID, Factor: factor, TargetFactor: target,
		HealthyStreak: healthyStreak, UnhealthyStreak: unhealthyStreak,
		Reason: reason, LastObservedAt: now, UpdatedAt: now,
	}
}

func (s *AccountMonitorService) prunePoolAvailabilitySnapshots(ctx context.Context) (int64, error) {
	if s == nil || s.repo == nil {
		return 0, nil
	}
	now := time.Now().UTC()
	if s.now != nil {
		now = s.now().UTC()
	}

	s.pruneMu.Lock()
	if !s.nextSnapshotPruneAt.IsZero() && now.Before(s.nextSnapshotPruneAt) {
		s.pruneMu.Unlock()
		return 0, nil
	}
	s.nextSnapshotPruneAt = now.Add(poolAvailabilitySnapshotPruneInterval)
	s.pruneMu.Unlock()

	deleted, err := s.repo.DeletePoolAvailabilityBefore(ctx, now.Add(-poolAvailabilitySnapshotRetention))
	if err != nil {
		return deleted, err
	}
	if deleted > 0 {
		slog.Info("account_monitor: pruned pool availability snapshots", "deleted", deleted)
	}
	return deleted, nil
}

func (s *AccountMonitorService) buildPoolAvailabilitySnapshots(ctx context.Context) ([]*PoolAvailabilitySnapshotRow, error) {
	pools, err := s.upstreamPoolRepo.ListUpstreamPools(ctx)
	if err != nil {
		return nil, fmt.Errorf("list upstream pools for availability: %w", err)
	}
	enabledPools := filterEnabledPools(pools)
	membersByPool, accountsByID := s.loadPoolMembersAndAccounts(ctx, enabledPools)
	checkedAt := time.Now().UTC()
	rows := make([]*PoolAvailabilitySnapshotRow, 0, len(enabledPools))
	for _, pool := range enabledPools {
		members := membersByPool[pool.ID]
		available := 0
		for _, member := range members {
			if poolMemberSchedulable(member, accountsByID[member.AccountID]) {
				available++
			}
		}
		status := MonitorStatusFailed
		if available > 0 {
			status = MonitorStatusOperational
		}
		rows = append(rows, &PoolAvailabilitySnapshotRow{
			PoolID: pool.ID, Status: status, TotalMembers: len(members),
			AvailableMembers: available, CheckedAt: checkedAt,
		})
	}
	return rows, nil
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
	warmCtx, cancel := context.WithTimeout(r.parentCtx, 5*time.Second)
	if err := r.svc.WarmRuntimeHealth(warmCtx); err != nil {
		slog.Warn("account_monitor: warm runtime health failed", "error", err)
	}
	cancel()
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

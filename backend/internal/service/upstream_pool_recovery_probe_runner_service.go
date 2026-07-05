package service

import (
	"context"
	"encoding/json"
	"strings"
	"sync"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/Wei-Shaw/sub2api/internal/pkg/logger"
	"github.com/Wei-Shaw/sub2api/internal/pkg/openai_compat"
	"github.com/robfig/cron/v3"
)

const upstreamPoolRecoveryProbeDefaultCron = "* * * * *"

var upstreamPoolRecoveryProbeCronParser = cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)

const upstreamPoolRecoveryProbeMinAge = 2 * time.Minute
const upstreamPoolRecoveryProbeFailedTempUnschedExtension = 2 * time.Minute

// UpstreamPoolRecoveryProbeRunnerService periodically probes upstream pool members
// that are still carrying recoverable runtime failure state, then restores them
// through the existing account recovery path when the probe succeeds.
type UpstreamPoolRecoveryProbeRunnerService struct {
	upstreamPoolRepo UpstreamPoolRepository
	accountRepo      AccountRepository
	accountTestSvc   interface {
		RunTestBackground(ctx context.Context, accountID int64, modelID string) (*ScheduledTestResult, error)
		RunTestBackgroundWithMode(ctx context.Context, accountID int64, modelID string, mode string) (*ScheduledTestResult, error)
	}
	rateLimitSvc interface {
		RecoverAccountAfterSuccessfulTest(ctx context.Context, accountID int64) (*SuccessfulTestRecoveryResult, error)
	}
	cfg *config.Config

	cron      *cron.Cron
	startOnce sync.Once
	stopOnce  sync.Once
}

func NewUpstreamPoolRecoveryProbeRunnerService(
	upstreamPoolRepo UpstreamPoolRepository,
	accountRepo AccountRepository,
	accountTestSvc interface {
		RunTestBackground(ctx context.Context, accountID int64, modelID string) (*ScheduledTestResult, error)
		RunTestBackgroundWithMode(ctx context.Context, accountID int64, modelID string, mode string) (*ScheduledTestResult, error)
	},
	rateLimitSvc interface {
		RecoverAccountAfterSuccessfulTest(ctx context.Context, accountID int64) (*SuccessfulTestRecoveryResult, error)
	},
	cfg *config.Config,
) *UpstreamPoolRecoveryProbeRunnerService {
	return &UpstreamPoolRecoveryProbeRunnerService{
		upstreamPoolRepo: upstreamPoolRepo,
		accountRepo:      accountRepo,
		accountTestSvc:   accountTestSvc,
		rateLimitSvc:     rateLimitSvc,
		cfg:              cfg,
	}
}

func (s *UpstreamPoolRecoveryProbeRunnerService) Start() {
	if s == nil {
		return
	}
	s.startOnce.Do(func() {
		loc := time.Local
		if s.cfg != nil && strings.TrimSpace(s.cfg.Timezone) != "" {
			if parsed, err := time.LoadLocation(strings.TrimSpace(s.cfg.Timezone)); err == nil && parsed != nil {
				loc = parsed
			}
		}

		c := cron.New(cron.WithParser(upstreamPoolRecoveryProbeCronParser), cron.WithLocation(loc))
		if _, err := c.AddFunc(upstreamPoolRecoveryProbeDefaultCron, func() { s.runScheduled() }); err != nil {
			logger.LegacyPrintf("service.upstream_pool_recovery_probe", "[UpstreamPoolRecoveryProbe] not started (invalid schedule): %v", err)
			return
		}
		s.cron = c
		s.cron.Start()
		logger.LegacyPrintf("service.upstream_pool_recovery_probe", "[UpstreamPoolRecoveryProbe] started (tick=every 1 minute, recovery-only)")
	})
}

func (s *UpstreamPoolRecoveryProbeRunnerService) Stop() {
	if s == nil {
		return
	}
	s.stopOnce.Do(func() {
		if s.cron != nil {
			ctx := s.cron.Stop()
			select {
			case <-ctx.Done():
			case <-time.After(3 * time.Second):
				logger.LegacyPrintf("service.upstream_pool_recovery_probe", "[UpstreamPoolRecoveryProbe] cron stop timed out")
			}
		}
	})
}

func (s *UpstreamPoolRecoveryProbeRunnerService) runScheduled() {
	if s == nil || s.upstreamPoolRepo == nil || s.accountRepo == nil || s.accountTestSvc == nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	pools, err := s.upstreamPoolRepo.ListUpstreamPools(ctx)
	if err != nil {
		logger.LegacyPrintf("service.upstream_pool_recovery_probe", "[UpstreamPoolRecoveryProbe] ListUpstreamPools error: %v", err)
		return
	}

	var candidates []int64
	for i := range pools {
		pool := pools[i]
		if !pool.Enabled || !isUpstreamPoolRecoverySupportedPlatform(pool.Platform) {
			continue
		}
		members, err := s.upstreamPoolRepo.ListUpstreamPoolMembers(ctx, pool.ID)
		if err != nil {
			logger.LegacyPrintf("service.upstream_pool_recovery_probe", "[UpstreamPoolRecoveryProbe] pool=%d ListUpstreamPoolMembers error: %v", pool.ID, err)
			continue
		}
		for j := range members {
			member := members[j]
			if !member.Enabled || member.ManualDrained {
				continue
			}
			if member.SchedulableOverride != nil && !*member.SchedulableOverride {
				continue
			}
			candidates = append(candidates, member.AccountID)
			if pool.PolicyJSON != nil {
				s.cachePoolPolicyForProbe(member.AccountID, pool.PolicyJSON)
			}
		}
	}

	if len(candidates) == 0 {
		return
	}

	sem := make(chan struct{}, 3)
	var wg sync.WaitGroup
	seen := make(map[int64]struct{}, len(candidates))
	for _, accountID := range candidates {
		if accountID <= 0 {
			continue
		}
		if _, ok := seen[accountID]; ok {
			continue
		}
		seen[accountID] = struct{}{}
		sem <- struct{}{}
		wg.Add(1)
		go func(id int64) {
			defer wg.Done()
			defer func() { <-sem }()
			s.runOneCandidate(ctx, id)
		}(accountID)
	}
	wg.Wait()
}

func (s *UpstreamPoolRecoveryProbeRunnerService) runOneCandidate(ctx context.Context, accountID int64) {
	account, err := s.accountRepo.GetByID(ctx, accountID)
	if err != nil || account == nil {
		return
	}
	if account.Status != StatusActive && account.Status != StatusError {
		return
	}
	if account.Status != StatusError && !hasRecoverableRuntimeState(account) {
		return
	}
	if account.LastUsedAt != nil && time.Since(*account.LastUsedAt) < upstreamPoolRecoveryProbeMinAge {
		return
	}
	if !isUpstreamPoolRecoverySupportedAccount(account) {
		return
	}

	result := runUpstreamPoolRecoveryLightweightProbe(ctx, account)
	if result == nil && s.accountTestSvc != nil {
		result, _ = s.accountTestSvc.RunTestBackground(ctx, account.ID, probeModelForPool(account.Platform, nil))
	}
	if result == nil {
		logger.LegacyPrintf("service.upstream_pool_recovery_probe", "[UpstreamPoolRecoveryProbe] account=%d lightweight probe skipped", accountID)
		s.extendTempUnschedOnFailedProbe(ctx, account, "probe_error")
		return
	}
	if result == nil || result.Status != "success" {
		s.extendTempUnschedOnFailedProbe(ctx, account, "probe_failed")
		return
	}

	if s.rateLimitSvc == nil {
		logger.LegacyPrintf("service.upstream_pool_recovery_probe", "[UpstreamPoolRecoveryProbe] account=%d probe success but recovery service missing", accountID)
		return
	}

	recovery, err := s.rateLimitSvc.RecoverAccountAfterSuccessfulTest(ctx, accountID)
	if err != nil {
		logger.LegacyPrintf("service.upstream_pool_recovery_probe", "[UpstreamPoolRecoveryProbe] account=%d recovery failed: %v", accountID, err)
		return
	}
	if recovery == nil {
		return
	}
	if recovery.ClearedError || recovery.ClearedRateLimit {
		logger.LegacyPrintf("service.upstream_pool_recovery_probe", "[UpstreamPoolRecoveryProbe] account=%d recovered error=%v rate_limit=%v", accountID, recovery.ClearedError, recovery.ClearedRateLimit)
	}
}

func isUpstreamPoolRecoverySupportedPlatform(platform string) bool {
	switch platform {
	case PlatformOpenAI, PlatformAnthropic:
		return true
	default:
		return false
	}
}

func isUpstreamPoolRecoverySupportedAccount(account *Account) bool {
	if account == nil {
		return false
	}
	switch account.Platform {
	case PlatformOpenAI:
		return account.Type == AccountTypeAPIKey
	case PlatformAnthropic:
		return account.Type == AccountTypeOAuth || account.Type == AccountTypeSetupToken || account.Type == AccountTypeAPIKey || account.Type == AccountTypeServiceAccount || account.Type == AccountTypeBedrock
	default:
		return false
	}
}

func runUpstreamPoolRecoveryLightweightProbe(ctx context.Context, account *Account) *ScheduledTestResult {
	if account == nil || account.Platform != PlatformOpenAI || account.Type != AccountTypeAPIKey {
		return nil
	}
	endpoint := account.GetOpenAIBaseURL()
	apiKey := account.GetOpenAIApiKey()
	if strings.TrimSpace(endpoint) == "" || strings.TrimSpace(apiKey) == "" {
		return nil
	}

	model := account.GetMappedModel(probeModelForPool(PlatformOpenAI, nil))
	opts := &CheckOptions{Lightweight: true}
	if openai_compat.ResolveResponsesSupport(account.Extra) == openai_compat.ResponsesSupportYes {
		opts.APIMode = MonitorAPIModeResponses
	}

	started := time.Now()
	check := runCheckForModel(ctx, MonitorProviderOpenAI, endpoint, apiKey, model, opts)
	finished := time.Now()
	status := "failed"
	if check != nil && (check.Status == MonitorStatusOperational || check.Status == MonitorStatusDegraded) {
		status = "success"
	}

	result := &ScheduledTestResult{
		Status:     status,
		StartedAt:  started,
		FinishedAt: finished,
		CreatedAt:  finished,
	}
	if check != nil {
		if check.Message != "" {
			result.ErrorMessage = check.Message
		}
		if check.LatencyMs != nil {
			result.LatencyMs = int64(*check.LatencyMs)
		}
	}
	return result
}

var upstreamPoolRecoveryProbePolicyCache sync.Map // key: int64(accountID), value: OpenAIRoutingPolicy

func (s *UpstreamPoolRecoveryProbeRunnerService) cachePoolPolicyForProbe(accountID int64, policyJSON map[string]any) {
	if accountID <= 0 || len(policyJSON) == 0 {
		return
	}
	policy := OpenAIRoutingPolicy{HasBinding: true}
	ApplyOpenAIRoutingPolicyJSON(&policy, policyJSON)
	upstreamPoolRecoveryProbePolicyCache.Store(accountID, policy)
}

func (s *UpstreamPoolRecoveryProbeRunnerService) extendTempUnschedOnFailedProbe(ctx context.Context, account *Account, reason string) {
	if s == nil || s.accountRepo == nil || account == nil || account.TempUnschedulableUntil == nil {
		return
	}
	now := time.Now()
	if !now.Before(*account.TempUnschedulableUntil) {
		return
	}
	extension := upstreamPoolRecoveryProbeFailedTempUnschedExtension
	if cached, ok := upstreamPoolRecoveryProbePolicyCache.Load(account.ID); ok {
		if policy, _ := cached.(OpenAIRoutingPolicy); policy.HasBinding {
			extension = policy.EffectiveHalfOpenProbeFailedExtension(extension)
		}
	}
	until := now.Add(extension)
	if account.TempUnschedulableUntil.After(until) {
		until = *account.TempUnschedulableUntil
	}
	state := &TempUnschedState{
		UntilUnix:       until.Unix(),
		TriggeredAtUnix: now.Unix(),
		StatusCode:      0,
		MatchedKeyword:  strings.TrimSpace(reason),
		RuleIndex:       -2,
		ErrorMessage:    "upstream pool recovery probe did not pass",
	}
	encodedReason := "upstream pool recovery probe did not pass"
	if raw, err := json.Marshal(state); err == nil {
		encodedReason = string(raw)
	}
	if err := s.accountRepo.SetTempUnschedulable(ctx, account.ID, until, encodedReason); err != nil {
		logger.LegacyPrintf("service.upstream_pool_recovery_probe", "[UpstreamPoolRecoveryProbe] account=%d extend temp-unsched failed: %v", account.ID, err)
		return
	}
	logger.LegacyPrintf("service.upstream_pool_recovery_probe", "[UpstreamPoolRecoveryProbe] account=%d half-open probe failed, temp-unsched extended until=%s", account.ID, until.Format(time.RFC3339))
}

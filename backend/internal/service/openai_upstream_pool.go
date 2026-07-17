package service

import (
	"context"
	"errors"
	"strings"
	"sync"
	"time"
)

const openAIRoutingMemberCacheTTL = 5 * time.Second

type cachedOpenAIRoutingPolicy struct {
	policy    OpenAIRoutingPolicy
	expiresAt time.Time
}

type cachedOpenAIRoutingMembers struct {
	hasBinding    bool
	memberConfigs map[int64]UpstreamPoolResolvedMemberConfig
	expiresAt     time.Time
}

type routingMemberCacheKey struct {
	groupID  int64
	platform string
}

type openAIRoutingGroupIDContextKey struct{}

func WithOpenAIRoutingGroupID(ctx context.Context, groupID *int64) context.Context {
	if ctx == nil || groupID == nil || *groupID <= 0 {
		return ctx
	}
	return context.WithValue(ctx, openAIRoutingGroupIDContextKey{}, *groupID)
}

func OpenAIRoutingGroupIDFromContext(ctx context.Context) *int64 {
	if ctx == nil {
		return nil
	}
	groupID, ok := ctx.Value(openAIRoutingGroupIDContextKey{}).(int64)
	if !ok || groupID <= 0 {
		return nil
	}
	return &groupID
}

func (s *OpenAIGatewayService) SetUpstreamPoolRepository(repo UpstreamPoolRepository) {
	if s == nil {
		return
	}
	s.upstreamPoolRepo = repo
	s.openAIRoutingPolicyCache = sync.Map{}
	s.openAIRoutingMemberCache = sync.Map{}
}

func (s *OpenAIGatewayService) tryAcquireAccountSlotForAccount(ctx context.Context, account *Account) (*AcquireResult, error) {
	if s.concurrencyService == nil {
		return &AcquireResult{Acquired: true, ReleaseFunc: func() {}}, nil
	}
	if account == nil {
		return nil, errors.New("account is required")
	}
	return s.concurrencyService.AcquireAccountSlotWithScope(ctx, account.ConcurrencyScope())
}

func (s *OpenAIGatewayService) isAccountAllowedByResolvedPool(ctx context.Context, groupID *int64, platform string, accountID int64) bool {
	if s == nil || s.upstreamPoolRepo == nil || groupID == nil || *groupID <= 0 || accountID <= 0 || strings.TrimSpace(platform) == "" {
		return true
	}
	hasBinding, memberConfigs, err := s.getResolvedPoolMemberConfigs(ctx, *groupID, platform)
	if err != nil || !hasBinding {
		return true
	}
	_, ok := memberConfigs[accountID]
	return ok
}

func (s *OpenAIGatewayService) getOpenAIRoutingPolicy(ctx context.Context, groupID *int64) OpenAIRoutingPolicy {
	if s == nil || s.upstreamPoolRepo == nil || groupID == nil || *groupID <= 0 {
		return OpenAIRoutingPolicy{}
	}
	if cached, ok := s.openAIRoutingPolicyCache.Load(*groupID); ok {
		if entry, _ := cached.(*cachedOpenAIRoutingPolicy); entry != nil && time.Now().Before(entry.expiresAt) {
			return entry.policy
		}
	}
	policy, err := s.upstreamPoolRepo.GetOpenAIRoutingPolicy(ctx, *groupID)
	if err != nil || policy == nil {
		return OpenAIRoutingPolicy{}
	}
	s.openAIRoutingPolicyCache.Store(*groupID, &cachedOpenAIRoutingPolicy{policy: *policy, expiresAt: time.Now().Add(5 * time.Second)})
	return *policy
}

func (s *OpenAIGatewayService) OpenAIRoutingPoolMode5xxCooldown(ctx context.Context, groupID *int64, defaultCooldown time.Duration) time.Duration {
	return s.getOpenAIRoutingPolicy(ctx, groupID).EffectivePoolMode5xxCooldown(defaultCooldown)
}

func (s *OpenAIGatewayService) OpenAIRoutingMaxFailoverHops(ctx context.Context, groupID *int64, defaultMax int) int {
	return s.getOpenAIRoutingPolicy(ctx, groupID).EffectiveMaxFailoverHops(defaultMax)
}

func (s *OpenAIGatewayService) getResolvedPoolMemberConfigs(ctx context.Context, groupID int64, platform string) (bool, map[int64]UpstreamPoolResolvedMemberConfig, error) {
	if s == nil || s.upstreamPoolRepo == nil || groupID <= 0 || strings.TrimSpace(platform) == "" {
		return false, nil, nil
	}
	key := routingMemberCacheKey{groupID: groupID, platform: strings.TrimSpace(platform)}
	if cached, ok := s.openAIRoutingMemberCache.Load(key); ok {
		if entry, _ := cached.(*cachedOpenAIRoutingMembers); entry != nil && time.Now().Before(entry.expiresAt) {
			return entry.hasBinding, entry.memberConfigs, nil
		}
	}
	resolved, err := s.upstreamPoolRepo.GetResolvedBindingByGroupAndPlatform(ctx, groupID, platform)
	if err != nil {
		return false, nil, err
	}
	hasBinding := resolved != nil && resolved.Binding != nil && resolved.Pool != nil
	memberConfigs := map[int64]UpstreamPoolResolvedMemberConfig{}
	if hasBinding && resolved.Pool.Enabled {
		for accountID, cfg := range resolved.MemberConfigs {
			memberConfigs[accountID] = cfg
		}
		if len(memberConfigs) == 0 {
			for accountID := range resolved.MemberIDs {
				memberConfigs[accountID] = UpstreamPoolResolvedMemberConfig{AccountID: accountID, Weight: 100}
			}
		}
	}
	s.openAIRoutingMemberCache.Store(key, &cachedOpenAIRoutingMembers{hasBinding: hasBinding, memberConfigs: memberConfigs, expiresAt: time.Now().Add(openAIRoutingMemberCacheTTL)})
	return hasBinding, memberConfigs, nil
}

func (s *OpenAIGatewayService) applyResolvedPoolMemberOverrides(ctx context.Context, groupID *int64, platform string, account *Account) *Account {
	if account == nil {
		return nil
	}
	if s == nil || s.upstreamPoolRepo == nil || groupID == nil || *groupID <= 0 || strings.TrimSpace(platform) == "" {
		return account
	}
	hasBinding, memberConfigs, err := s.getResolvedPoolMemberConfigs(ctx, *groupID, platform)
	if err != nil || !hasBinding {
		return account
	}
	memberCfg, ok := memberConfigs[account.ID]
	if !ok {
		return nil
	}
	return applyResolvedPoolMemberConfig(account, memberCfg)
}

func applyResolvedPoolMemberConfig(account *Account, memberCfg UpstreamPoolResolvedMemberConfig) *Account {
	if account == nil {
		return nil
	}
	cloned := *account
	if memberCfg.SchedulableOverride != nil {
		cloned.Schedulable = *memberCfg.SchedulableOverride
	}
	if memberCfg.PriorityOverride != nil {
		cloned.Priority = *memberCfg.PriorityOverride
	}
	if memberCfg.MaxConcurrencyOverride != nil && *memberCfg.MaxConcurrencyOverride > 0 {
		override := *memberCfg.MaxConcurrencyOverride
		cloned.Concurrency = override
		cloned.LoadFactor = &override
	}
	if memberCfg.Weight > 0 {
		cloned.PoolMemberWeight = memberCfg.Weight
	} else {
		cloned.PoolMemberWeight = 100
	}
	return &cloned
}

func (s *OpenAIGatewayService) applyResolvedPoolMembershipToAccounts(ctx context.Context, groupID *int64, platform string, accounts []Account) ([]Account, error) {
	if s == nil || s.upstreamPoolRepo == nil || groupID == nil || *groupID <= 0 || strings.TrimSpace(platform) == "" {
		return accounts, nil
	}
	hasBinding, memberConfigs, err := s.getResolvedPoolMemberConfigs(ctx, *groupID, platform)
	if err != nil {
		return nil, err
	}
	if !hasBinding {
		return accounts, nil
	}
	filtered := make([]Account, 0, len(accounts))
	for i := range accounts {
		memberCfg, ok := memberConfigs[accounts[i].ID]
		if !ok {
			continue
		}
		if resolved := applyResolvedPoolMemberConfig(&accounts[i], memberCfg); resolved != nil {
			filtered = append(filtered, *resolved)
		}
	}
	return filtered, nil
}

package service

import (
	"context"
	"strconv"
	"strings"
)

type resolvedUpstreamPoolBindingContextKey struct{}

type resolvedUpstreamPoolBindingContextValue struct {
	key     string
	binding *UpstreamPoolResolvedBinding
}

func (s *GatewayService) filterAccountsByResolvedUpstreamPool(ctx context.Context, groupID *int64, platform string, accounts []Account) []Account {
	if len(accounts) == 0 {
		return accounts
	}
	resolved := s.resolvedUpstreamPoolBindingForRouting(ctx, groupID, platform)
	if resolved == nil || resolved.Binding == nil || resolved.Pool == nil {
		return accounts
	}
	if len(resolved.MemberConfigs) == 0 {
		return nil
	}
	filtered := make([]Account, 0, len(accounts))
	for i := range accounts {
		if overridden := applyResolvedUpstreamPoolMemberOverridesToAccount(&accounts[i], resolved.MemberConfigs); overridden != nil {
			filtered = append(filtered, *overridden)
		}
	}
	return filtered
}

func (s *GatewayService) applyResolvedUpstreamPoolMemberOverrides(ctx context.Context, groupID *int64, platform string, account *Account) *Account {
	if account == nil {
		return nil
	}
	resolved := s.resolvedUpstreamPoolBindingForRouting(ctx, groupID, platform)
	if resolved == nil || resolved.Binding == nil || resolved.Pool == nil {
		return account
	}
	return applyResolvedUpstreamPoolMemberOverridesToAccount(account, resolved.MemberConfigs)
}

func (s *GatewayService) isAccountAllowedByResolvedPool(ctx context.Context, groupID *int64, platform string, accountID int64) bool {
	resolved := s.resolvedUpstreamPoolBindingForRouting(ctx, groupID, platform)
	if resolved == nil || resolved.Binding == nil || resolved.Pool == nil {
		return true
	}
	_, allowed := resolved.MemberConfigs[accountID]
	return allowed
}

func (s *GatewayService) resolvedUpstreamPoolBindingForRouting(ctx context.Context, groupID *int64, platform string) *UpstreamPoolResolvedBinding {
	platform = strings.TrimSpace(platform)
	if s == nil || s.upstreamPoolRepo == nil || groupID == nil || *groupID <= 0 || platform == "" {
		return nil
	}
	key := strconv.FormatInt(*groupID, 10) + ":" + platform
	if cached, ok := ctx.Value(resolvedUpstreamPoolBindingContextKey{}).(resolvedUpstreamPoolBindingContextValue); ok && cached.key == key {
		return cached.binding
	}
	resolved, err := s.upstreamPoolRepo.GetResolvedBindingByGroupAndPlatform(ctx, *groupID, platform)
	if err != nil || resolved == nil || resolved.Binding == nil || resolved.Pool == nil {
		return nil
	}
	if !resolved.Pool.Enabled {
		resolved.MemberIDs = map[int64]struct{}{}
		resolved.MemberConfigs = map[int64]UpstreamPoolResolvedMemberConfig{}
		return resolved
	}
	return resolved
}

func withResolvedUpstreamPoolBindingContext(ctx context.Context, groupID *int64, platform string, resolved *UpstreamPoolResolvedBinding) context.Context {
	if ctx == nil || groupID == nil || *groupID <= 0 || strings.TrimSpace(platform) == "" {
		return ctx
	}
	key := strconv.FormatInt(*groupID, 10) + ":" + strings.TrimSpace(platform)
	return context.WithValue(ctx, resolvedUpstreamPoolBindingContextKey{}, resolvedUpstreamPoolBindingContextValue{key: key, binding: resolved})
}

func (s *GatewayService) WithResolvedUpstreamPoolBindingForRoutingContext(ctx context.Context, groupID *int64, platform string) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	resolved := s.resolvedUpstreamPoolBindingForRouting(ctx, groupID, platform)
	return withResolvedUpstreamPoolBindingContext(ctx, groupID, platform, resolved)
}

func applyResolvedUpstreamPoolMemberOverridesToAccount(account *Account, memberConfigs map[int64]UpstreamPoolResolvedMemberConfig) *Account {
	if account == nil || len(memberConfigs) == 0 {
		return nil
	}
	memberCfg, ok := memberConfigs[account.ID]
	if !ok {
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

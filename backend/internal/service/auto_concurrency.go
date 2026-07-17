package service

import (
	"context"
	"fmt"
	"sync"
	"time"
)

const (
	autoConcurrencyBootstrapAccount = 20
	autoConcurrencyBootstrapGroup   = 50
	autoConcurrencyUnboundedCeiling = 100
	autoConcurrencyRecoverSamples   = 12
	autoConcurrencySlowUpstreamMs   = 15_000
	autoConcurrencyStateTTL         = 24 * time.Hour
)

// AutoConcurrencyState is shared by all gateway processes when the backing
// concurrency cache supports it. Limit is an effective runtime limit, never a
// replacement for an administrator's configured hard ceiling.
type AutoConcurrencyState struct {
	Limit          int    `json:"limit"`
	HealthySamples int    `json:"healthy_samples"`
	Reason         string `json:"reason"`
	UpdatedAtUnix  int64  `json:"updated_at_unix"`
}

// AutoConcurrencyStateCache is deliberately optional so existing cache
// implementations remain compatible. Production Redis implements it; tests
// and legacy cache adapters use the process-local fallback.
type AutoConcurrencyStateCache interface {
	GetAutoConcurrencyState(ctx context.Context, key string) (*AutoConcurrencyState, error)
	SetAutoConcurrencyState(ctx context.Context, key string, state *AutoConcurrencyState, ttl time.Duration) error
}

type autoConcurrencyController struct {
	cache AutoConcurrencyStateCache

	mu    sync.Mutex
	local map[string]AutoConcurrencyState
}

func newAutoConcurrencyController(cache ConcurrencyCache) *autoConcurrencyController {
	stateCache, _ := cache.(AutoConcurrencyStateCache)
	return &autoConcurrencyController{
		cache: stateCache,
		local: make(map[string]AutoConcurrencyState),
	}
}

func (c *autoConcurrencyController) resolveScope(ctx context.Context, input AccountConcurrencyScope) (AccountConcurrencyScope, error) {
	if c == nil || input.AccountID <= 0 {
		return input, nil
	}

	accountCeiling := autoConcurrencyCeiling(input.AccountLimit)
	accountLimit, err := c.limit(ctx, autoConcurrencyAccountKey(input.AccountID), accountCeiling, autoConcurrencyBootstrapAccount)
	if err != nil {
		return input, err
	}
	input.AccountLimit = accountLimit

	if input.Capacity == nil || input.Capacity.GroupID <= 0 || input.Capacity.GroupLimit <= 0 {
		return input, nil
	}

	capacity := *input.Capacity
	groupCeiling := capacity.GroupLimit
	groupLimit, err := c.limit(ctx, autoConcurrencyGroupKey(capacity.GroupID), groupCeiling, autoConcurrencyBootstrapGroup)
	if err != nil {
		return input, err
	}
	capacity.GroupLimit = groupLimit

	// Capacity-member hard limits remain absolute ceilings when configured.
	// The dynamic account limit is always applied, including when the legacy
	// member limit was left empty.
	if capacity.MemberHardLimit > 0 && capacity.MemberHardLimit < input.AccountLimit {
		input.AccountLimit = capacity.MemberHardLimit
	}
	capacity.MemberHardLimit = input.AccountLimit
	if capacity.MemberSoftShare <= 0 || capacity.MemberSoftShare > input.AccountLimit {
		capacity.MemberSoftShare = input.AccountLimit
	}
	input.Capacity = &capacity
	return input, nil
}

func (c *autoConcurrencyController) observe(ctx context.Context, scope AccountConcurrencyScope, upstreamWaitMs int64) error {
	if c == nil || scope.AccountID <= 0 || upstreamWaitMs <= 0 {
		return nil
	}
	slow := upstreamWaitMs >= autoConcurrencySlowUpstreamMs
	if err := c.adjust(ctx, autoConcurrencyAccountKey(scope.AccountID), autoConcurrencyCeiling(scope.AccountLimit), autoConcurrencyBootstrapAccount, slow, "upstream_post_write_wait"); err != nil {
		return err
	}
	if scope.Capacity == nil || scope.Capacity.GroupID <= 0 || scope.Capacity.GroupLimit <= 0 {
		return nil
	}
	return c.adjust(ctx, autoConcurrencyGroupKey(scope.Capacity.GroupID), scope.Capacity.GroupLimit, autoConcurrencyBootstrapGroup, slow, "upstream_post_write_wait")
}

func (c *autoConcurrencyController) contract(ctx context.Context, scope AccountConcurrencyScope, reason string) error {
	if c == nil || scope.AccountID <= 0 {
		return nil
	}
	if err := c.adjust(ctx, autoConcurrencyAccountKey(scope.AccountID), autoConcurrencyCeiling(scope.AccountLimit), autoConcurrencyBootstrapAccount, true, reason); err != nil {
		return err
	}
	if scope.Capacity == nil || scope.Capacity.GroupID <= 0 || scope.Capacity.GroupLimit <= 0 {
		return nil
	}
	return c.adjust(ctx, autoConcurrencyGroupKey(scope.Capacity.GroupID), scope.Capacity.GroupLimit, autoConcurrencyBootstrapGroup, true, reason)
}

func (c *autoConcurrencyController) limit(ctx context.Context, key string, ceiling, bootstrap int) (int, error) {
	state, err := c.get(ctx, key)
	if err != nil {
		return 0, err
	}
	if state == nil || state.Limit <= 0 {
		return minPositive(ceiling, bootstrap), nil
	}
	return minPositive(ceiling, state.Limit), nil
}

func (c *autoConcurrencyController) adjust(ctx context.Context, key string, ceiling, bootstrap int, slow bool, reason string) error {
	state, err := c.get(ctx, key)
	if err != nil {
		return err
	}
	if state == nil {
		state = &AutoConcurrencyState{Limit: minPositive(ceiling, bootstrap)}
	}
	state.Limit = minPositive(ceiling, state.Limit)
	if slow {
		state.Limit = maxAutoConcurrencyInt(1, (state.Limit+1)/2)
		state.HealthySamples = 0
		state.Reason = reason
	} else {
		state.HealthySamples++
		if state.HealthySamples >= autoConcurrencyRecoverSamples && state.Limit < ceiling {
			step := maxAutoConcurrencyInt(1, state.Limit/10)
			state.Limit = minPositive(ceiling, state.Limit+step)
			state.HealthySamples = 0
			state.Reason = "healthy_recovery"
		}
	}
	state.UpdatedAtUnix = time.Now().Unix()
	return c.set(ctx, key, state)
}

func (c *autoConcurrencyController) get(ctx context.Context, key string) (*AutoConcurrencyState, error) {
	if c.cache != nil {
		return c.cache.GetAutoConcurrencyState(ctx, key)
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	state, ok := c.local[key]
	if !ok {
		return nil, nil
	}
	copy := state
	return &copy, nil
}

func (c *autoConcurrencyController) set(ctx context.Context, key string, state *AutoConcurrencyState) error {
	if state == nil {
		return nil
	}
	if c.cache != nil {
		return c.cache.SetAutoConcurrencyState(ctx, key, state, autoConcurrencyStateTTL)
	}
	c.mu.Lock()
	c.local[key] = *state
	c.mu.Unlock()
	return nil
}

func autoConcurrencyCeiling(value int) int {
	if value > 0 {
		return value
	}
	return autoConcurrencyUnboundedCeiling
}

func minPositive(a, b int) int {
	if a <= 0 {
		return b
	}
	if b <= 0 || a < b {
		return a
	}
	return b
}

func maxAutoConcurrencyInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func autoConcurrencyAccountKey(accountID int64) string {
	return fmt.Sprintf("account:%d", accountID)
}

func autoConcurrencyGroupKey(groupID int64) string {
	return fmt.Sprintf("capacity:%d", groupID)
}

package service

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"time"
)

var (
	ErrUpstreamPoolNotFound = errors.New("upstream pool not found")
)

const (
	UpstreamPoolPlatformOpenAI = PlatformOpenAI

	UpstreamPoolSchedulerModeBasic    = "basic"
	UpstreamPoolSchedulerModeAdvanced = "advanced"
)

type UpstreamPool struct {
	ID                             int64
	Name                           string
	Code                           string
	Platform                       string
	Description                    string
	Enabled                        bool
	SchedulerMode                  string
	DefaultRequiredCapability      string
	DefaultRequiredTransport       string
	StickyEnabled                  bool
	StickyTTLSeconds               int
	StickyEscapeEnabled            bool
	StickyEscapeErrorRateThreshold float64
	StickyEscapeTTFTMSThreshold    int
	LoadBalanceEnabled             bool
	FailoverEnabled                bool
	TopK                           int
	MaxFailoverHops                int
	WaitTimeoutMS                  int
	MaxWaiting                     int
	PolicyJSON                     map[string]any
	CreatedAt                      time.Time
	UpdatedAt                      time.Time
}

type UpstreamPoolMember struct {
	ID                            int64
	PoolID                        int64
	AccountID                     int64
	AccountName                   string
	AccountPlatform               string
	AccountStatus                 string
	AccountSchedulable            bool
	RuntimeStatus                 string
	RuntimeReason                 string
	RuntimeErrorRate              *float64
	RuntimeTTFTMs                 *int
	RuntimeLastUsedAt             *time.Time
	RuntimeRateLimitResetAt       *time.Time
	RuntimeOverloadUntil          *time.Time
	RuntimeTempUnschedulableUntil *time.Time
	Enabled                       bool
	SchedulableOverride           *bool
	ManualDrained                 bool
	Weight                        int
	PriorityOverride              *int
	MaxConcurrencyOverride        *int
	Notes                         string
	JoinedAt                      time.Time
	UpdatedAt                     time.Time
}

type UpstreamPoolBinding struct {
	ID               int64
	GroupID          int64
	GroupName        string
	GroupPlatform    string
	PoolID           int64
	Platform         string
	Models           []string
	RequestPathScope []string
	Priority         int
	Enabled          bool
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

type CreateUpstreamPoolInput struct {
	Name                           string
	Code                           string
	Platform                       string
	Description                    string
	Enabled                        bool
	SchedulerMode                  string
	DefaultRequiredCapability      string
	DefaultRequiredTransport       string
	StickyEnabled                  bool
	StickyTTLSeconds               int
	StickyEscapeEnabled            bool
	StickyEscapeErrorRateThreshold float64
	StickyEscapeTTFTMSThreshold    int
	LoadBalanceEnabled             bool
	FailoverEnabled                bool
	TopK                           int
	MaxFailoverHops                int
	WaitTimeoutMS                  int
	MaxWaiting                     int
	PolicyJSON                     map[string]any
}

type UpdateUpstreamPoolInput struct {
	Name                           *string
	Code                           *string
	Platform                       *string
	Description                    *string
	Enabled                        *bool
	SchedulerMode                  *string
	DefaultRequiredCapability      *string
	DefaultRequiredTransport       *string
	StickyEnabled                  *bool
	StickyTTLSeconds               *int
	StickyEscapeEnabled            *bool
	StickyEscapeErrorRateThreshold *float64
	StickyEscapeTTFTMSThreshold    *int
	LoadBalanceEnabled             *bool
	FailoverEnabled                *bool
	TopK                           *int
	MaxFailoverHops                *int
	WaitTimeoutMS                  *int
	MaxWaiting                     *int
	PolicyJSON                     map[string]any
}

type CreateUpstreamPoolMemberInput struct {
	AccountID              int64
	Enabled                bool
	SchedulableOverride    *bool
	ManualDrained          bool
	Weight                 int
	PriorityOverride       *int
	MaxConcurrencyOverride *int
	Notes                  string
}

type UpdateUpstreamPoolMemberInput struct {
	Enabled                *bool
	SchedulableOverride    *bool
	ManualDrained          *bool
	Weight                 *int
	PriorityOverride       *int
	MaxConcurrencyOverride *int
	Notes                  *string
}

type CreateUpstreamPoolBindingInput struct {
	GroupID          int64
	PoolID           int64
	Platform         string
	Models           []string
	RequestPathScope []string
	Priority         int
	Enabled          bool
}

type UpdateUpstreamPoolBindingInput struct {
	GroupID          *int64
	PoolID           *int64
	Platform         *string
	Models           []string
	RequestPathScope []string
	Priority         *int
	Enabled          *bool
}

type UpstreamPoolResolvedBinding struct {
	Binding       *UpstreamPoolBinding
	Pool          *UpstreamPool
	MemberIDs     map[int64]struct{}
	MemberConfigs map[int64]UpstreamPoolResolvedMemberConfig
}

type UpstreamPoolResolvedMemberConfig struct {
	AccountID              int64
	Weight                 int
	PriorityOverride       *int
	MaxConcurrencyOverride *int
	SchedulableOverride    *bool
}

type OpenAIAccountRuntimeSnapshot struct {
	AccountID int64
	ErrorRate float64
	TTFTMs    *int
}

type OpenAIAccountRuntimeObserver interface {
	SnapshotOpenAIAccountRuntime(accountIDs []int64) map[int64]OpenAIAccountRuntimeSnapshot
}

// OpenAIRoutingPolicy 是 OpenAI 路由层实际使用的池策略快照。
// HasBinding=false 表示该 group 当前没有命中可用的 upstream pool 绑定，
// 路由层应回退到默认调度行为。
type OpenAIRoutingPolicy struct {
	HasBinding                     bool
	PoolID                         int64
	PoolCode                       string
	PoolName                       string
	SchedulerMode                  string
	StickyEnabled                  bool
	StickyEscapeEnabled            *bool
	StickyEscapeErrorRateThreshold float64
	StickyEscapeTTFTMSThreshold    int
	LoadBalanceEnabled             bool
	FailoverEnabled                bool
	TopK                           int
	MaxFailoverHops                int
	WaitTimeout                    time.Duration
	MaxWaiting                     int
	CacheAffinityEnabled           *bool
	PoolMode5xxCooldown            time.Duration
	HalfOpenProbeFailedExtension   time.Duration
}

func (p OpenAIRoutingPolicy) EffectiveTopK(defaultTopK int) int {
	if defaultTopK <= 0 {
		defaultTopK = 1
	}
	if !p.HasBinding {
		return defaultTopK
	}
	if !p.LoadBalanceEnabled {
		return 1
	}
	if p.TopK > 0 {
		return p.TopK
	}
	return defaultTopK
}

func (p OpenAIRoutingPolicy) EffectiveMaxFailoverHops(defaultMax int) int {
	if !p.HasBinding {
		return defaultMax
	}
	if !p.FailoverEnabled {
		return 0
	}
	if p.MaxFailoverHops > 0 {
		return p.MaxFailoverHops
	}
	return defaultMax
}

func (p OpenAIRoutingPolicy) EffectiveCacheAffinityEnabled(defaultEnabled bool) bool {
	if !p.HasBinding || p.CacheAffinityEnabled == nil {
		return defaultEnabled
	}
	return *p.CacheAffinityEnabled
}

func (p OpenAIRoutingPolicy) EffectivePoolMode5xxCooldown(defaultCooldown time.Duration) time.Duration {
	if !p.HasBinding || p.PoolMode5xxCooldown <= 0 {
		return defaultCooldown
	}
	return p.PoolMode5xxCooldown
}

func (p OpenAIRoutingPolicy) EffectiveHalfOpenProbeFailedExtension(defaultExtension time.Duration) time.Duration {
	if !p.HasBinding || p.HalfOpenProbeFailedExtension <= 0 {
		return defaultExtension
	}
	return p.HalfOpenProbeFailedExtension
}

func ApplyOpenAIRoutingPolicyJSON(policy *OpenAIRoutingPolicy, policyJSON map[string]any) {
	if policy == nil || len(policyJSON) == 0 {
		return
	}
	if cacheAffinity := policyJSONMap(policyJSON, "cache_affinity"); len(cacheAffinity) > 0 {
		if enabled, ok := policyJSONBool(cacheAffinity, "enabled"); ok {
			policy.CacheAffinityEnabled = &enabled
		}
	}
	if circuitBreaker := policyJSONMap(policyJSON, "circuit_breaker"); len(circuitBreaker) > 0 {
		if minutes, ok := policyJSONInt(circuitBreaker, "openai_pool_mode_5xx_cooldown_minutes"); ok && minutes > 0 {
			policy.PoolMode5xxCooldown = time.Duration(minutes) * time.Minute
		}
		if seconds, ok := policyJSONInt(circuitBreaker, "half_open_probe_failed_extension_seconds"); ok && seconds > 0 {
			policy.HalfOpenProbeFailedExtension = time.Duration(seconds) * time.Second
		}
	}
}

func policyJSONMap(input map[string]any, key string) map[string]any {
	raw, ok := input[key]
	if !ok || raw == nil {
		return nil
	}
	if typed, ok := raw.(map[string]any); ok {
		return typed
	}
	if rawBytes, ok := raw.(json.RawMessage); ok {
		var out map[string]any
		if err := json.Unmarshal(rawBytes, &out); err == nil {
			return out
		}
	}
	return nil
}

func policyJSONBool(input map[string]any, key string) (bool, bool) {
	raw, ok := input[key]
	if !ok {
		return false, false
	}
	switch typed := raw.(type) {
	case bool:
		return typed, true
	case string:
		switch strings.ToLower(strings.TrimSpace(typed)) {
		case "true", "1", "yes", "on":
			return true, true
		case "false", "0", "no", "off":
			return false, true
		}
	}
	return false, false
}

func policyJSONInt(input map[string]any, key string) (int, bool) {
	raw, ok := input[key]
	if !ok || raw == nil {
		return 0, false
	}
	switch typed := raw.(type) {
	case int:
		return typed, true
	case int64:
		return int(typed), true
	case float64:
		return int(typed), true
	case json.Number:
		if value, err := typed.Int64(); err == nil {
			return int(value), true
		}
	case string:
		if value, err := strconv.Atoi(strings.TrimSpace(typed)); err == nil {
			return value, true
		}
	}
	return 0, false
}

type UpstreamPoolRepository interface {
	ListUpstreamPools(ctx context.Context) ([]UpstreamPool, error)
	GetUpstreamPoolByID(ctx context.Context, id int64) (*UpstreamPool, error)
	CreateUpstreamPool(ctx context.Context, input *UpstreamPool) (*UpstreamPool, error)
	UpdateUpstreamPool(ctx context.Context, input *UpstreamPool) (*UpstreamPool, error)
	DeleteUpstreamPool(ctx context.Context, id int64) error
	ListUpstreamPoolMembers(ctx context.Context, poolID int64) ([]UpstreamPoolMember, error)
	GetUpstreamPoolMemberByID(ctx context.Context, id int64) (*UpstreamPoolMember, error)
	CreateUpstreamPoolMember(ctx context.Context, input *UpstreamPoolMember) (*UpstreamPoolMember, error)
	UpdateUpstreamPoolMember(ctx context.Context, input *UpstreamPoolMember) (*UpstreamPoolMember, error)
	DeleteUpstreamPoolMember(ctx context.Context, id int64) error
	ListUpstreamPoolBindings(ctx context.Context) ([]UpstreamPoolBinding, error)
	GetUpstreamPoolBindingByID(ctx context.Context, id int64) (*UpstreamPoolBinding, error)
	CreateUpstreamPoolBinding(ctx context.Context, input *UpstreamPoolBinding) (*UpstreamPoolBinding, error)
	UpdateUpstreamPoolBinding(ctx context.Context, input *UpstreamPoolBinding) (*UpstreamPoolBinding, error)
	DeleteUpstreamPoolBinding(ctx context.Context, id int64) error
	ListEnabledMemberAccountIDsByGroupAndPlatform(ctx context.Context, groupID int64, platform string) (map[int64]struct{}, error)
	GetResolvedBindingByGroupAndPlatform(ctx context.Context, groupID int64, platform string) (*UpstreamPoolResolvedBinding, error)
	GetOpenAIRoutingPolicy(ctx context.Context, groupID int64) (*OpenAIRoutingPolicy, error)
}

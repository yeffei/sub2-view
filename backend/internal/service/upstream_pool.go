package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	infraerrors "github.com/Wei-Shaw/sub2api/internal/pkg/errors"
)

var (
	ErrUpstreamPoolNotFound = errors.New("upstream pool not found")
)

const (
	UpstreamPoolPlatformOpenAI = PlatformOpenAI

	UpstreamPoolSchedulerModeBasic    = "basic"
	UpstreamPoolSchedulerModeAdvanced = "advanced"

	UpstreamPoolAccountTypeStrategyAll             = "all"
	UpstreamPoolAccountTypeStrategyOAuthPreferred  = "oauth_preferred"
	UpstreamPoolAccountTypeStrategyOAuthOnly       = "oauth_only"
	UpstreamPoolAccountTypeStrategyAPIKeyPreferred = "apikey_preferred"

	UpstreamPoolMemberSyncModeMembershipOnly           = "membership_only"
	UpstreamPoolMemberSyncModeOverwriteSchedulerFields = "overwrite_scheduler_fields"
)

type UpstreamPool struct {
	ID                             int64
	Name                           string
	AccountTypeStrategy            string
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
	AccountType                   string
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
	SourceType                    string
	SourceSetID                   *int64
	SourceSetName                 string
	Editable                      bool
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

type UpstreamAccountSet struct {
	ID           int64
	Name         string
	Code         string
	Platform     string
	Description  string
	Enabled      bool
	AccountCount int
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type UpstreamAccountSetMember struct {
	SetID                         int64
	AccountID                     int64
	AccountName                   string
	AccountPlatform               string
	AccountType                   string
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
	AddedAt                       time.Time
}

type UpstreamPoolMemberSet struct {
	ID          int64
	PoolID      int64
	SetID       int64
	SetName     string
	SetCode     string
	SetPlatform string
	Enabled     bool
	Notes       string
	JoinedAt    time.Time
	UpdatedAt   time.Time
}

type UpstreamPoolPlatformUsage struct {
	DirectMemberCount int
	MemberSetCount    int
	BindingCount      int
}

func (u UpstreamPoolPlatformUsage) IsLocked() bool {
	return u.DirectMemberCount > 0 || u.MemberSetCount > 0 || u.BindingCount > 0
}

type UpstreamAccountSetPlatformUsage struct {
	MemberCount      int
	PoolBindingCount int
}

func (u UpstreamAccountSetPlatformUsage) IsLocked() bool {
	return u.MemberCount > 0 || u.PoolBindingCount > 0
}

type UpstreamPoolMemberSyncMode string

type UpstreamPoolMemberSyncTarget struct {
	AccountID int64
	Enabled   bool
}

type UpstreamPoolMemberSyncPreviewInput struct {
	Mode UpstreamPoolMemberSyncMode
}

type UpstreamPoolMemberSyncApplyInput struct {
	Mode UpstreamPoolMemberSyncMode
}

type UpstreamPoolMemberSyncChange struct {
	AccountID       int64
	AccountName     string
	AccountPlatform string
	AccountType     string
	Action          string
	Overwrites      []string
}

type UpstreamPoolMemberSyncResult struct {
	PoolID             int64
	Platform           string
	Mode               UpstreamPoolMemberSyncMode
	CreateCount        int
	UpdateCount        int
	DeleteCount        int
	SkipCount          int
	OverwriteRiskCount int
	Creates            []UpstreamPoolMemberSyncChange
	Updates            []UpstreamPoolMemberSyncChange
	Deletes            []UpstreamPoolMemberSyncChange
	Skips              []UpstreamPoolMemberSyncChange
}

type CreateUpstreamPoolInput struct {
	Name                           string
	AccountTypeStrategy            string
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
	AccountTypeStrategy            *string
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
	Enabled                   *bool
	SchedulableOverrideSet    bool
	SchedulableOverride       *bool
	ManualDrained             *bool
	Weight                    *int
	PriorityOverrideSet       bool
	PriorityOverride          *int
	MaxConcurrencyOverrideSet bool
	MaxConcurrencyOverride    *int
	NotesSet                  bool
	Notes                     *string
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

type CreateUpstreamAccountSetInput struct {
	Name        string
	Code        string
	Platform    string
	Description string
	Enabled     bool
}

type UpdateUpstreamAccountSetInput struct {
	Name        *string
	Code        *string
	Platform    *string
	Description *string
	Enabled     *bool
}

type AddUpstreamAccountSetMembersInput struct {
	AccountIDs []int64
}

type CreateUpstreamPoolMemberSetInput struct {
	SetID   int64
	Enabled bool
	Notes   string
}

type UpdateUpstreamPoolMemberSetInput struct {
	Enabled *bool
	Notes   *string
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
	AccountTypeStrategy            string
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
	if p.MaxFailoverHops >= 0 {
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

func NormalizeUpstreamPoolAccountTypeStrategy(value string) string {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case UpstreamPoolAccountTypeStrategyOAuthPreferred:
		return UpstreamPoolAccountTypeStrategyOAuthPreferred
	case UpstreamPoolAccountTypeStrategyOAuthOnly:
		return UpstreamPoolAccountTypeStrategyOAuthOnly
	case UpstreamPoolAccountTypeStrategyAPIKeyPreferred:
		return UpstreamPoolAccountTypeStrategyAPIKeyPreferred
	default:
		return UpstreamPoolAccountTypeStrategyAll
	}
}

func UpstreamPoolAccountTypeStrategyFromPolicyJSON(policyJSON map[string]any) string {
	if routing := policyJSONMap(policyJSON, "routing"); len(routing) > 0 {
		if strategy, ok := policyJSONString(routing, "account_type_strategy"); ok {
			return NormalizeUpstreamPoolAccountTypeStrategy(strategy)
		}
	}
	if strategy, ok := policyJSONString(policyJSON, "account_type_strategy"); ok {
		return NormalizeUpstreamPoolAccountTypeStrategy(strategy)
	}
	return UpstreamPoolAccountTypeStrategyAll
}

func SetUpstreamPoolAccountTypeStrategyPolicyJSON(policyJSON map[string]any, strategy string) map[string]any {
	if policyJSON == nil {
		policyJSON = map[string]any{}
	}
	normalized := NormalizeUpstreamPoolAccountTypeStrategy(strategy)
	routing := policyJSONMap(policyJSON, "routing")
	if routing == nil {
		routing = map[string]any{}
	}
	routing["account_type_strategy"] = normalized
	policyJSON["routing"] = routing
	return policyJSON
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
	if routing := policyJSONMap(policyJSON, "routing"); len(routing) > 0 {
		if strategy, ok := policyJSONString(routing, "account_type_strategy"); ok {
			policy.AccountTypeStrategy = NormalizeUpstreamPoolAccountTypeStrategy(strategy)
		}
	}
	if strategy, ok := policyJSONString(policyJSON, "account_type_strategy"); ok && policy.AccountTypeStrategy == "" {
		policy.AccountTypeStrategy = NormalizeUpstreamPoolAccountTypeStrategy(strategy)
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

func policyJSONString(input map[string]any, key string) (string, bool) {
	raw, ok := input[key]
	if !ok || raw == nil {
		return "", false
	}
	switch typed := raw.(type) {
	case string:
		value := strings.TrimSpace(typed)
		return value, value != ""
	case json.RawMessage:
		var value string
		if err := json.Unmarshal(typed, &value); err == nil {
			value = strings.TrimSpace(value)
			return value, value != ""
		}
	}
	return "", false
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
	ListUpstreamAccountSets(ctx context.Context) ([]UpstreamAccountSet, error)
	GetUpstreamAccountSetByID(ctx context.Context, id int64) (*UpstreamAccountSet, error)
	CreateUpstreamAccountSet(ctx context.Context, input *UpstreamAccountSet) (*UpstreamAccountSet, error)
	UpdateUpstreamAccountSet(ctx context.Context, input *UpstreamAccountSet) (*UpstreamAccountSet, error)
	DeleteUpstreamAccountSet(ctx context.Context, id int64) error
	ListUpstreamAccountSetMembers(ctx context.Context, setID int64) ([]UpstreamAccountSetMember, error)
	AddUpstreamAccountSetMembers(ctx context.Context, setID int64, accountIDs []int64) error
	DeleteUpstreamAccountSetMember(ctx context.Context, setID, accountID int64) error
	ListUpstreamPoolMemberSets(ctx context.Context, poolID int64) ([]UpstreamPoolMemberSet, error)
	GetUpstreamPoolMemberSetByID(ctx context.Context, id int64) (*UpstreamPoolMemberSet, error)
	CreateUpstreamPoolMemberSet(ctx context.Context, input *UpstreamPoolMemberSet) (*UpstreamPoolMemberSet, error)
	UpdateUpstreamPoolMemberSet(ctx context.Context, input *UpstreamPoolMemberSet) (*UpstreamPoolMemberSet, error)
	DeleteUpstreamPoolMemberSet(ctx context.Context, id int64) error
	ListUpstreamPoolBindings(ctx context.Context) ([]UpstreamPoolBinding, error)
	GetUpstreamPoolBindingByID(ctx context.Context, id int64) (*UpstreamPoolBinding, error)
	CreateUpstreamPoolBinding(ctx context.Context, input *UpstreamPoolBinding) (*UpstreamPoolBinding, error)
	UpdateUpstreamPoolBinding(ctx context.Context, input *UpstreamPoolBinding) (*UpstreamPoolBinding, error)
	DeleteUpstreamPoolBinding(ctx context.Context, id int64) error
	ListEnabledMemberAccountIDsByGroupAndPlatform(ctx context.Context, groupID int64, platform string) (map[int64]struct{}, error)
	GetResolvedBindingByGroupAndPlatform(ctx context.Context, groupID int64, platform string) (*UpstreamPoolResolvedBinding, error)
	GetOpenAIRoutingPolicy(ctx context.Context, groupID int64) (*OpenAIRoutingPolicy, error)
}

type UpstreamPoolPlatformUsageReader interface {
	GetUpstreamPoolPlatformUsage(ctx context.Context, poolID int64) (UpstreamPoolPlatformUsage, error)
	GetUpstreamAccountSetPlatformUsage(ctx context.Context, setID int64) (UpstreamAccountSetPlatformUsage, error)
}

type UpstreamPoolMemberSyncer interface {
	SyncUpstreamPoolDirectMembers(ctx context.Context, poolID int64, targets []UpstreamPoolMember, mode UpstreamPoolMemberSyncMode) (*UpstreamPoolMemberSyncResult, error)
}

func NormalizeUpstreamPoolMemberSyncMode(mode UpstreamPoolMemberSyncMode) UpstreamPoolMemberSyncMode {
	switch strings.TrimSpace(string(mode)) {
	case UpstreamPoolMemberSyncModeOverwriteSchedulerFields:
		return UpstreamPoolMemberSyncModeOverwriteSchedulerFields
	default:
		return UpstreamPoolMemberSyncModeMembershipOnly
	}
}

func NewUpstreamPoolBadRequest(reason, message string) error {
	return infraerrors.BadRequest(reason, message)
}

func NewUpstreamPoolPlatformLockedError(resource string, usage any) error {
	return infraerrors.Conflict(
		"UPSTREAM_POOL_PLATFORM_LOCKED",
		fmt.Sprintf("%s platform cannot be changed while it has members or bindings", resource),
	).WithMetadata(map[string]string{
		"resource": fmt.Sprintf("%s", resource),
		"usage":    fmt.Sprintf("%+v", usage),
	})
}

func BuildUpstreamPoolMemberSyncResult(pool *UpstreamPool, current []UpstreamPoolMember, targets []UpstreamPoolMember, mode UpstreamPoolMemberSyncMode) *UpstreamPoolMemberSyncResult {
	if pool == nil {
		return &UpstreamPoolMemberSyncResult{}
	}
	mode = NormalizeUpstreamPoolMemberSyncMode(mode)
	result := &UpstreamPoolMemberSyncResult{
		PoolID:   pool.ID,
		Platform: pool.Platform,
		Mode:     mode,
	}
	targetByAccountID := make(map[int64]UpstreamPoolMember, len(targets))
	for _, target := range targets {
		if target.AccountID <= 0 {
			continue
		}
		targetByAccountID[target.AccountID] = target
	}
	directByAccountID := make(map[int64]UpstreamPoolMember, len(current))
	for _, member := range current {
		if member.AccountID <= 0 {
			continue
		}
		if strings.EqualFold(strings.TrimSpace(member.SourceType), "account_set") || member.ID <= 0 {
			result.Skips = append(result.Skips, UpstreamPoolMemberSyncChange{
				AccountID:       member.AccountID,
				AccountName:     member.AccountName,
				AccountPlatform: member.AccountPlatform,
				AccountType:     member.AccountType,
				Action:          "skip_account_set_member",
			})
			continue
		}
		directByAccountID[member.AccountID] = member
	}
	for _, target := range targets {
		member, exists := directByAccountID[target.AccountID]
		if !exists {
			result.Creates = append(result.Creates, upstreamPoolMemberSyncChangeFromTarget(target, "create", nil))
			continue
		}
		if mode == UpstreamPoolMemberSyncModeOverwriteSchedulerFields {
			overwrites := upstreamPoolMemberSyncOverwrites(member, target)
			change := upstreamPoolMemberSyncChangeFromTarget(target, "update", overwrites)
			change.AccountID = member.AccountID
			if change.AccountName == "" {
				change.AccountName = member.AccountName
			}
			result.Updates = append(result.Updates, change)
			if len(overwrites) > 0 {
				result.OverwriteRiskCount++
			}
			continue
		}
		result.Skips = append(result.Skips, upstreamPoolMemberSyncChangeFromTarget(target, "skip_existing_direct_member", nil))
	}
	for _, member := range directByAccountID {
		if _, ok := targetByAccountID[member.AccountID]; ok {
			continue
		}
		result.Deletes = append(result.Deletes, upstreamPoolMemberSyncChangeFromTarget(member, "delete", upstreamPoolMemberSyncDeleteOverwrites(member)))
		if len(upstreamPoolMemberSyncDeleteOverwrites(member)) > 0 {
			result.OverwriteRiskCount++
		}
	}
	result.CreateCount = len(result.Creates)
	result.UpdateCount = len(result.Updates)
	result.DeleteCount = len(result.Deletes)
	result.SkipCount = len(result.Skips)
	return result
}

func upstreamPoolMemberSyncChangeFromTarget(member UpstreamPoolMember, action string, overwrites []string) UpstreamPoolMemberSyncChange {
	return UpstreamPoolMemberSyncChange{
		AccountID:       member.AccountID,
		AccountName:     member.AccountName,
		AccountPlatform: member.AccountPlatform,
		AccountType:     member.AccountType,
		Action:          action,
		Overwrites:      append([]string{}, overwrites...),
	}
}

func upstreamPoolMemberSyncOverwrites(current, target UpstreamPoolMember) []string {
	overwrites := make([]string, 0, 6)
	if current.Enabled != target.Enabled {
		overwrites = append(overwrites, "enabled")
	}
	if current.ManualDrained != target.ManualDrained {
		overwrites = append(overwrites, "manual_drained")
	}
	if current.Weight != target.Weight {
		overwrites = append(overwrites, "weight")
	}
	if current.SchedulableOverride != nil {
		overwrites = append(overwrites, "schedulable_override")
	}
	if current.PriorityOverride != nil {
		overwrites = append(overwrites, "priority_override")
	}
	if current.MaxConcurrencyOverride != nil {
		overwrites = append(overwrites, "max_concurrency_override")
	}
	if strings.TrimSpace(current.Notes) != strings.TrimSpace(target.Notes) {
		overwrites = append(overwrites, "notes")
	}
	return overwrites
}

func upstreamPoolMemberSyncDeleteOverwrites(member UpstreamPoolMember) []string {
	overwrites := []string{"membership"}
	if member.Weight != 100 {
		overwrites = append(overwrites, "weight")
	}
	if member.SchedulableOverride != nil {
		overwrites = append(overwrites, "schedulable_override")
	}
	if member.PriorityOverride != nil {
		overwrites = append(overwrites, "priority_override")
	}
	if member.MaxConcurrencyOverride != nil {
		overwrites = append(overwrites, "max_concurrency_override")
	}
	if member.ManualDrained {
		overwrites = append(overwrites, "manual_drained")
	}
	if strings.TrimSpace(member.Notes) != "" {
		overwrites = append(overwrites, "notes")
	}
	return overwrites
}

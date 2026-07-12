package dto

import (
	"time"

	"github.com/Wei-Shaw/sub2api/internal/service"
)

const upstreamPoolTimeLayout = "2006-01-02T15:04:05Z07:00"

type UpstreamPool struct {
	ID                             int64          `json:"id"`
	Name                           string         `json:"name"`
	Code                           string         `json:"code"`
	Platform                       string         `json:"platform"`
	Description                    string         `json:"description"`
	Enabled                        bool           `json:"enabled"`
	SchedulerMode                  string         `json:"scheduler_mode"`
	AccountTypeStrategy            string         `json:"account_type_strategy"`
	DefaultRequiredCapability      string         `json:"default_required_capability"`
	DefaultRequiredTransport       string         `json:"default_required_transport"`
	StickyEnabled                  bool           `json:"sticky_enabled"`
	StickyTTLSeconds               int            `json:"sticky_ttl_seconds"`
	StickyEscapeEnabled            bool           `json:"sticky_escape_enabled"`
	StickyEscapeErrorRateThreshold float64        `json:"sticky_escape_error_rate_threshold"`
	StickyEscapeTTFTMSThreshold    int            `json:"sticky_escape_ttft_ms_threshold"`
	LoadBalanceEnabled             bool           `json:"load_balance_enabled"`
	AutoWeightEnabled              bool           `json:"auto_weight_enabled"`
	FailoverEnabled                bool           `json:"failover_enabled"`
	TopK                           int            `json:"top_k"`
	MaxFailoverHops                int            `json:"max_failover_hops"`
	WaitTimeoutMS                  int            `json:"wait_timeout_ms"`
	MaxWaiting                     int            `json:"max_waiting"`
	PolicyJSON                     map[string]any `json:"policy_json"`
	CreatedAt                      string         `json:"created_at"`
	UpdatedAt                      string         `json:"updated_at"`
}

type UpstreamPoolMember struct {
	ID                            int64    `json:"id"`
	PoolID                        int64    `json:"pool_id"`
	AccountID                     int64    `json:"account_id"`
	AccountName                   string   `json:"account_name"`
	AccountPlatform               string   `json:"account_platform"`
	AccountType                   string   `json:"account_type"`
	AccountStatus                 string   `json:"account_status"`
	AccountSchedulable            bool     `json:"account_schedulable"`
	RuntimeStatus                 string   `json:"runtime_status"`
	RuntimeReason                 string   `json:"runtime_reason"`
	RuntimeErrorRate              *float64 `json:"runtime_error_rate"`
	RuntimeTTFTMs                 *int     `json:"runtime_ttft_ms"`
	RuntimeLastUsedAt             *string  `json:"runtime_last_used_at"`
	RuntimeRateLimitResetAt       *string  `json:"runtime_rate_limit_reset_at"`
	RuntimeOverloadUntil          *string  `json:"runtime_overload_until"`
	RuntimeTempUnschedulableUntil *string  `json:"runtime_temp_unschedulable_until"`
	RuntimeWeightFactor           float64  `json:"runtime_weight_factor"`
	EffectiveWeight               int      `json:"effective_weight"`
	RuntimeWeightReason           string   `json:"runtime_weight_reason"`
	RuntimeWeightUpdatedAt        *string  `json:"runtime_weight_updated_at"`
	Enabled                       bool     `json:"enabled"`
	SchedulableOverride           *bool    `json:"schedulable_override"`
	ManualDrained                 bool     `json:"manual_drained"`
	Weight                        int      `json:"weight"`
	PriorityOverride              *int     `json:"priority_override"`
	MaxConcurrencyOverride        *int     `json:"max_concurrency_override"`
	Notes                         string   `json:"notes"`
	JoinedAt                      string   `json:"joined_at"`
	UpdatedAt                     string   `json:"updated_at"`
	SourceType                    string   `json:"source_type"`
	SourceSetID                   *int64   `json:"source_set_id"`
	SourceSetName                 string   `json:"source_set_name"`
	Editable                      bool     `json:"editable"`
}

type UpstreamPoolBinding struct {
	ID               int64    `json:"id"`
	GroupID          int64    `json:"group_id"`
	GroupName        string   `json:"group_name"`
	GroupPlatform    string   `json:"group_platform"`
	PoolID           int64    `json:"pool_id"`
	Platform         string   `json:"platform"`
	Models           []string `json:"models"`
	RequestPathScope []string `json:"request_path_scope"`
	Priority         int      `json:"priority"`
	Enabled          bool     `json:"enabled"`
	CreatedAt        string   `json:"created_at"`
	UpdatedAt        string   `json:"updated_at"`
}

type UpstreamAccountSet struct {
	ID           int64  `json:"id"`
	Name         string `json:"name"`
	Code         string `json:"code"`
	Platform     string `json:"platform"`
	Description  string `json:"description"`
	Enabled      bool   `json:"enabled"`
	AccountCount int    `json:"account_count"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

type UpstreamAccountSetMember struct {
	SetID                         int64    `json:"set_id"`
	AccountID                     int64    `json:"account_id"`
	AccountName                   string   `json:"account_name"`
	AccountPlatform               string   `json:"account_platform"`
	AccountType                   string   `json:"account_type"`
	AccountStatus                 string   `json:"account_status"`
	AccountSchedulable            bool     `json:"account_schedulable"`
	RuntimeStatus                 string   `json:"runtime_status"`
	RuntimeReason                 string   `json:"runtime_reason"`
	RuntimeErrorRate              *float64 `json:"runtime_error_rate"`
	RuntimeTTFTMs                 *int     `json:"runtime_ttft_ms"`
	RuntimeLastUsedAt             *string  `json:"runtime_last_used_at"`
	RuntimeRateLimitResetAt       *string  `json:"runtime_rate_limit_reset_at"`
	RuntimeOverloadUntil          *string  `json:"runtime_overload_until"`
	RuntimeTempUnschedulableUntil *string  `json:"runtime_temp_unschedulable_until"`
	AddedAt                       string   `json:"added_at"`
}

type UpstreamPoolMemberSet struct {
	ID          int64  `json:"id"`
	PoolID      int64  `json:"pool_id"`
	SetID       int64  `json:"set_id"`
	SetName     string `json:"set_name"`
	SetCode     string `json:"set_code"`
	SetPlatform string `json:"set_platform"`
	Enabled     bool   `json:"enabled"`
	Notes       string `json:"notes"`
	JoinedAt    string `json:"joined_at"`
	UpdatedAt   string `json:"updated_at"`
}

type UpstreamPoolMemberSyncChange struct {
	AccountID       int64    `json:"account_id"`
	AccountName     string   `json:"account_name"`
	AccountPlatform string   `json:"account_platform"`
	AccountType     string   `json:"account_type"`
	Action          string   `json:"action"`
	Overwrites      []string `json:"overwrites"`
}

type UpstreamPoolMemberSyncResult struct {
	PoolID             int64                          `json:"pool_id"`
	Platform           string                         `json:"platform"`
	Mode               string                         `json:"mode"`
	CreateCount        int                            `json:"create_count"`
	UpdateCount        int                            `json:"update_count"`
	DeleteCount        int                            `json:"delete_count"`
	SkipCount          int                            `json:"skip_count"`
	OverwriteRiskCount int                            `json:"overwrite_risk_count"`
	Creates            []UpstreamPoolMemberSyncChange `json:"creates"`
	Updates            []UpstreamPoolMemberSyncChange `json:"updates"`
	Deletes            []UpstreamPoolMemberSyncChange `json:"deletes"`
	Skips              []UpstreamPoolMemberSyncChange `json:"skips"`
}

func UpstreamPoolFromService(pool *service.UpstreamPool) *UpstreamPool {
	if pool == nil {
		return nil
	}
	return &UpstreamPool{
		ID:                             pool.ID,
		Name:                           pool.Name,
		Code:                           pool.Code,
		Platform:                       pool.Platform,
		Description:                    pool.Description,
		Enabled:                        pool.Enabled,
		SchedulerMode:                  pool.SchedulerMode,
		AccountTypeStrategy:            pool.AccountTypeStrategy,
		DefaultRequiredCapability:      pool.DefaultRequiredCapability,
		DefaultRequiredTransport:       pool.DefaultRequiredTransport,
		StickyEnabled:                  pool.StickyEnabled,
		StickyTTLSeconds:               pool.StickyTTLSeconds,
		StickyEscapeEnabled:            pool.StickyEscapeEnabled,
		StickyEscapeErrorRateThreshold: pool.StickyEscapeErrorRateThreshold,
		StickyEscapeTTFTMSThreshold:    pool.StickyEscapeTTFTMSThreshold,
		LoadBalanceEnabled:             pool.LoadBalanceEnabled,
		AutoWeightEnabled:              pool.AutoWeightEnabled,
		FailoverEnabled:                pool.FailoverEnabled,
		TopK:                           pool.TopK,
		MaxFailoverHops:                pool.MaxFailoverHops,
		WaitTimeoutMS:                  pool.WaitTimeoutMS,
		MaxWaiting:                     pool.MaxWaiting,
		PolicyJSON:                     pool.PolicyJSON,
		CreatedAt:                      pool.CreatedAt.Format(upstreamPoolTimeLayout),
		UpdatedAt:                      pool.UpdatedAt.Format(upstreamPoolTimeLayout),
	}
}

func UpstreamPoolMemberFromService(member *service.UpstreamPoolMember) *UpstreamPoolMember {
	if member == nil {
		return nil
	}
	return &UpstreamPoolMember{
		ID:                            member.ID,
		PoolID:                        member.PoolID,
		AccountID:                     member.AccountID,
		AccountName:                   member.AccountName,
		AccountPlatform:               member.AccountPlatform,
		AccountType:                   member.AccountType,
		AccountStatus:                 member.AccountStatus,
		AccountSchedulable:            member.AccountSchedulable,
		RuntimeStatus:                 member.RuntimeStatus,
		RuntimeReason:                 member.RuntimeReason,
		RuntimeErrorRate:              member.RuntimeErrorRate,
		RuntimeTTFTMs:                 member.RuntimeTTFTMs,
		RuntimeLastUsedAt:             formatOptionalTime(member.RuntimeLastUsedAt),
		RuntimeRateLimitResetAt:       formatOptionalTime(member.RuntimeRateLimitResetAt),
		RuntimeOverloadUntil:          formatOptionalTime(member.RuntimeOverloadUntil),
		RuntimeTempUnschedulableUntil: formatOptionalTime(member.RuntimeTempUnschedulableUntil),
		RuntimeWeightFactor:           member.RuntimeWeightFactor,
		EffectiveWeight:               member.EffectiveWeight,
		RuntimeWeightReason:           member.RuntimeWeightReason,
		RuntimeWeightUpdatedAt:        formatOptionalTime(member.RuntimeWeightUpdatedAt),
		Enabled:                       member.Enabled,
		SchedulableOverride:           member.SchedulableOverride,
		ManualDrained:                 member.ManualDrained,
		Weight:                        member.Weight,
		PriorityOverride:              member.PriorityOverride,
		MaxConcurrencyOverride:        member.MaxConcurrencyOverride,
		Notes:                         member.Notes,
		JoinedAt:                      member.JoinedAt.Format(upstreamPoolTimeLayout),
		UpdatedAt:                     member.UpdatedAt.Format(upstreamPoolTimeLayout),
		SourceType:                    member.SourceType,
		SourceSetID:                   member.SourceSetID,
		SourceSetName:                 member.SourceSetName,
		Editable:                      member.Editable,
	}
}

func formatOptionalTime(value *time.Time) *string {
	if value == nil {
		return nil
	}
	formatted := value.Format(upstreamPoolTimeLayout)
	return &formatted
}

func UpstreamPoolBindingFromService(binding *service.UpstreamPoolBinding) *UpstreamPoolBinding {
	if binding == nil {
		return nil
	}
	return &UpstreamPoolBinding{
		ID:               binding.ID,
		GroupID:          binding.GroupID,
		GroupName:        binding.GroupName,
		GroupPlatform:    binding.GroupPlatform,
		PoolID:           binding.PoolID,
		Platform:         binding.Platform,
		Models:           binding.Models,
		RequestPathScope: binding.RequestPathScope,
		Priority:         binding.Priority,
		Enabled:          binding.Enabled,
		CreatedAt:        binding.CreatedAt.Format(upstreamPoolTimeLayout),
		UpdatedAt:        binding.UpdatedAt.Format(upstreamPoolTimeLayout),
	}
}

func UpstreamAccountSetFromService(item *service.UpstreamAccountSet) *UpstreamAccountSet {
	if item == nil {
		return nil
	}
	return &UpstreamAccountSet{
		ID:           item.ID,
		Name:         item.Name,
		Code:         item.Code,
		Platform:     item.Platform,
		Description:  item.Description,
		Enabled:      item.Enabled,
		AccountCount: item.AccountCount,
		CreatedAt:    item.CreatedAt.Format(upstreamPoolTimeLayout),
		UpdatedAt:    item.UpdatedAt.Format(upstreamPoolTimeLayout),
	}
}

func UpstreamAccountSetMemberFromService(item *service.UpstreamAccountSetMember) *UpstreamAccountSetMember {
	if item == nil {
		return nil
	}
	return &UpstreamAccountSetMember{
		SetID:                         item.SetID,
		AccountID:                     item.AccountID,
		AccountName:                   item.AccountName,
		AccountPlatform:               item.AccountPlatform,
		AccountType:                   item.AccountType,
		AccountStatus:                 item.AccountStatus,
		AccountSchedulable:            item.AccountSchedulable,
		RuntimeStatus:                 item.RuntimeStatus,
		RuntimeReason:                 item.RuntimeReason,
		RuntimeErrorRate:              item.RuntimeErrorRate,
		RuntimeTTFTMs:                 item.RuntimeTTFTMs,
		RuntimeLastUsedAt:             formatOptionalTime(item.RuntimeLastUsedAt),
		RuntimeRateLimitResetAt:       formatOptionalTime(item.RuntimeRateLimitResetAt),
		RuntimeOverloadUntil:          formatOptionalTime(item.RuntimeOverloadUntil),
		RuntimeTempUnschedulableUntil: formatOptionalTime(item.RuntimeTempUnschedulableUntil),
		AddedAt:                       item.AddedAt.Format(upstreamPoolTimeLayout),
	}
}

func UpstreamPoolMemberSetFromService(item *service.UpstreamPoolMemberSet) *UpstreamPoolMemberSet {
	if item == nil {
		return nil
	}
	return &UpstreamPoolMemberSet{
		ID:          item.ID,
		PoolID:      item.PoolID,
		SetID:       item.SetID,
		SetName:     item.SetName,
		SetCode:     item.SetCode,
		SetPlatform: item.SetPlatform,
		Enabled:     item.Enabled,
		Notes:       item.Notes,
		JoinedAt:    item.JoinedAt.Format(upstreamPoolTimeLayout),
		UpdatedAt:   item.UpdatedAt.Format(upstreamPoolTimeLayout),
	}
}

func UpstreamPoolMemberSyncResultFromService(item *service.UpstreamPoolMemberSyncResult) *UpstreamPoolMemberSyncResult {
	if item == nil {
		return nil
	}
	return &UpstreamPoolMemberSyncResult{
		PoolID:             item.PoolID,
		Platform:           item.Platform,
		Mode:               string(item.Mode),
		CreateCount:        item.CreateCount,
		UpdateCount:        item.UpdateCount,
		DeleteCount:        item.DeleteCount,
		SkipCount:          item.SkipCount,
		OverwriteRiskCount: item.OverwriteRiskCount,
		Creates:            upstreamPoolMemberSyncChangesFromService(item.Creates),
		Updates:            upstreamPoolMemberSyncChangesFromService(item.Updates),
		Deletes:            upstreamPoolMemberSyncChangesFromService(item.Deletes),
		Skips:              upstreamPoolMemberSyncChangesFromService(item.Skips),
	}
}

func upstreamPoolMemberSyncChangesFromService(items []service.UpstreamPoolMemberSyncChange) []UpstreamPoolMemberSyncChange {
	out := make([]UpstreamPoolMemberSyncChange, 0, len(items))
	for _, item := range items {
		out = append(out, UpstreamPoolMemberSyncChange{
			AccountID:       item.AccountID,
			AccountName:     item.AccountName,
			AccountPlatform: item.AccountPlatform,
			AccountType:     item.AccountType,
			Action:          item.Action,
			Overwrites:      append([]string{}, item.Overwrites...),
		})
	}
	return out
}

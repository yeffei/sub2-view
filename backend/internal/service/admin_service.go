package service

import (
	"context"
	"errors"
	"fmt"
	"hash/fnv"
	"net/http"
	"sort"
	"strings"
	"time"

	dbent "github.com/Wei-Shaw/sub2api/ent"
	infraerrors "github.com/Wei-Shaw/sub2api/internal/pkg/errors"
	"github.com/Wei-Shaw/sub2api/internal/pkg/pagination"
)

// AdminService interface defines admin management operations
type AdminService interface {
	// User management
	ListUsers(ctx context.Context, page, pageSize int, filters UserListFilters, sortBy, sortOrder string) ([]User, int64, error)
	GetUser(ctx context.Context, id int64) (*User, error)
	GetUserIncludeDeleted(ctx context.Context, id int64) (*User, error)
	CreateUser(ctx context.Context, input *CreateUserInput) (*User, error)
	UpdateUser(ctx context.Context, id int64, input *UpdateUserInput) (*User, error)
	DeleteUser(ctx context.Context, id int64) error
	UpdateUserBalance(ctx context.Context, userID int64, balance float64, operation string, notes string) (*User, error)
	BatchUpdateConcurrency(ctx context.Context, userIDs []int64, value int, mode string) (int, error)
	GetUserAPIKeys(ctx context.Context, userID int64, page, pageSize int, sortBy, sortOrder string) ([]APIKey, int64, error)
	GetUserUsageStats(ctx context.Context, userID int64, period string) (any, error)
	GetUserRPMStatus(ctx context.Context, userID int64) (*UserRPMStatus, error)
	// GetUserBalanceHistory returns paginated balance/concurrency change records for a user.
	// codeType is optional - pass empty string to return all types.
	// Also returns totalRecharged (sum of all positive balance top-ups).
	GetUserBalanceHistory(ctx context.Context, userID int64, page, pageSize int, codeType string) ([]RedeemCode, int64, float64, error)
	BindUserAuthIdentity(ctx context.Context, userID int64, input AdminBindAuthIdentityInput) (*AdminBoundAuthIdentity, error)

	// Group management
	ListGroups(ctx context.Context, page, pageSize int, platform, status, search string, isExclusive *bool, sortBy, sortOrder string) ([]Group, int64, error)
	GetAllGroups(ctx context.Context) ([]Group, error)
	GetAllGroupsByPlatform(ctx context.Context, platform string) ([]Group, error)
	// GetAllGroupsIncludingInactive returns all groups regardless of status (active + disabled),
	// ordered by sort_order then id. Used by the API Key group filter dropdown.
	GetAllGroupsIncludingInactive(ctx context.Context) ([]Group, error)
	GetGroup(ctx context.Context, id int64) (*Group, error)
	GetGroupModelsListCandidates(ctx context.Context, id int64, platform string) ([]string, error)
	CreateGroup(ctx context.Context, input *CreateGroupInput) (*Group, error)
	UpdateGroup(ctx context.Context, id int64, input *UpdateGroupInput) (*Group, error)
	DeleteGroup(ctx context.Context, id int64) error
	GetGroupAPIKeys(ctx context.Context, groupID int64, page, pageSize int) ([]APIKey, int64, error)
	GetGroupRateMultipliers(ctx context.Context, groupID int64) ([]UserGroupRateEntry, error)
	ClearGroupRateMultipliers(ctx context.Context, groupID int64) error
	BatchSetGroupRateMultipliers(ctx context.Context, groupID int64, entries []GroupRateMultiplierInput) error
	ClearGroupRPMOverrides(ctx context.Context, groupID int64) error
	BatchSetGroupRPMOverrides(ctx context.Context, groupID int64, entries []GroupRPMOverrideInput) error
	UpdateGroupSortOrders(ctx context.Context, updates []GroupSortOrderUpdate) error

	// Upstream pool management
	ListUpstreamPools(ctx context.Context) ([]UpstreamPool, error)
	GetUpstreamPoolByID(ctx context.Context, id int64) (*UpstreamPool, error)
	CreateUpstreamPool(ctx context.Context, input *CreateUpstreamPoolInput) (*UpstreamPool, error)
	UpdateUpstreamPool(ctx context.Context, id int64, input *UpdateUpstreamPoolInput) (*UpstreamPool, error)
	DeleteUpstreamPool(ctx context.Context, id int64) error
	ListUpstreamPoolMembers(ctx context.Context, poolID int64) ([]UpstreamPoolMember, error)
	PreviewUpstreamPoolMemberSync(ctx context.Context, poolID int64, input *UpstreamPoolMemberSyncPreviewInput) (*UpstreamPoolMemberSyncResult, error)
	ApplyUpstreamPoolMemberSync(ctx context.Context, poolID int64, input *UpstreamPoolMemberSyncApplyInput) (*UpstreamPoolMemberSyncResult, error)
	CreateUpstreamPoolMember(ctx context.Context, poolID int64, input *CreateUpstreamPoolMemberInput) (*UpstreamPoolMember, error)
	UpdateUpstreamPoolMember(ctx context.Context, id int64, input *UpdateUpstreamPoolMemberInput) (*UpstreamPoolMember, error)
	DeleteUpstreamPoolMember(ctx context.Context, id int64) error
	ListUpstreamAccountSets(ctx context.Context) ([]UpstreamAccountSet, error)
	CreateUpstreamAccountSet(ctx context.Context, input *CreateUpstreamAccountSetInput) (*UpstreamAccountSet, error)
	UpdateUpstreamAccountSet(ctx context.Context, id int64, input *UpdateUpstreamAccountSetInput) (*UpstreamAccountSet, error)
	DeleteUpstreamAccountSet(ctx context.Context, id int64) error
	ListUpstreamAccountSetMembers(ctx context.Context, setID int64) ([]UpstreamAccountSetMember, error)
	AddUpstreamAccountSetMembers(ctx context.Context, setID int64, input *AddUpstreamAccountSetMembersInput) error
	DeleteUpstreamAccountSetMember(ctx context.Context, setID, accountID int64) error
	ListUpstreamPoolMemberSets(ctx context.Context, poolID int64) ([]UpstreamPoolMemberSet, error)
	CreateUpstreamPoolMemberSet(ctx context.Context, poolID int64, input *CreateUpstreamPoolMemberSetInput) (*UpstreamPoolMemberSet, error)
	UpdateUpstreamPoolMemberSet(ctx context.Context, id int64, input *UpdateUpstreamPoolMemberSetInput) (*UpstreamPoolMemberSet, error)
	DeleteUpstreamPoolMemberSet(ctx context.Context, id int64) error
	ListUpstreamPoolBindings(ctx context.Context) ([]UpstreamPoolBinding, error)
	CreateUpstreamPoolBinding(ctx context.Context, input *CreateUpstreamPoolBindingInput) (*UpstreamPoolBinding, error)
	UpdateUpstreamPoolBinding(ctx context.Context, id int64, input *UpdateUpstreamPoolBindingInput) (*UpstreamPoolBinding, error)
	DeleteUpstreamPoolBinding(ctx context.Context, id int64) error

	// API Key management (admin)
	AdminUpdateAPIKeyGroupID(ctx context.Context, keyID int64, groupID *int64) (*AdminUpdateAPIKeyGroupIDResult, error)
	AdminResetAPIKeyRateLimitUsage(ctx context.Context, keyID int64) (*APIKey, error)

	// ReplaceUserGroup 替换用户的专属分组：授予新分组权限、迁移 Key、移除旧分组权限
	ReplaceUserGroup(ctx context.Context, userID, oldGroupID, newGroupID int64) (*ReplaceUserGroupResult, error)

	// Account management
	ListAccounts(ctx context.Context, page, pageSize int, platform, accountType, status, anomalyReason, search string, groupID int64, privacyMode string, sortBy, sortOrder string) ([]Account, int64, error)
	// ListAccountsForSchedulerScoreFilter 返回符合过滤条件的全部账号（不分页），
	// 作为账号列表页计算 OpenAI 调度分数的过滤范围池。
	ListAccountsForSchedulerScoreFilter(ctx context.Context, platform, accountType, status, search string, groupID int64, privacyMode string) ([]Account, error)
	// ListOpenAISchedulableAccountsForSchedulerScore 返回指定分组（nil 为未分组）内
	// 可调度的 OpenAI 账号，用于按组计算调度分数。
	ListOpenAISchedulableAccountsForSchedulerScore(ctx context.Context, groupID *int64) ([]Account, error)
	GetAccount(ctx context.Context, id int64) (*Account, error)
	GetAccountsByIDs(ctx context.Context, ids []int64) ([]*Account, error)
	CreateAccount(ctx context.Context, input *CreateAccountInput) (*Account, error)
	// DuplicateAccount creates an independent account from an existing account's configuration.
	// First-class runtime columns are intentionally reset by the normal account creation path.
	DuplicateAccount(ctx context.Context, id int64, actorScope, operationKey string) (*Account, error)
	// RecoverDuplicateAccount returns a previously committed duplicate for an ambiguous retry.
	// It never creates an account.
	RecoverDuplicateAccount(ctx context.Context, id int64, actorScope, operationKey string) (*Account, error)
	UpdateAccount(ctx context.Context, id int64, input *UpdateAccountInput) (*Account, error)
	// UpdateAccountExtra 仅对 Extra 做 JSONB 增量合并（key 级覆盖），不会影响其它字段或运行态键。
	// 用于刷新流程持久化 account_uuid / org_uuid 等少量键，避免被全量快照覆盖。
	UpdateAccountExtra(ctx context.Context, id int64, updates map[string]any) error
	DeleteAccount(ctx context.Context, id int64) error
	RefreshAccountCredentials(ctx context.Context, id int64) (*Account, error)
	ClearAccountError(ctx context.Context, id int64) (*Account, error)
	SetAccountError(ctx context.Context, id int64, errorMsg string) error
	// EnsureOpenAIPrivacy 检查 OpenAI OAuth 账号 privacy_mode，未设置则尝试关闭训练数据共享并持久化。
	EnsureOpenAIPrivacy(ctx context.Context, account *Account) string
	// EnsureAntigravityPrivacy 检查 Antigravity OAuth 账号 privacy_mode，未设置则调用 setUserSettings 并持久化。
	EnsureAntigravityPrivacy(ctx context.Context, account *Account) string
	// ForceOpenAIPrivacy 强制重新设置 OpenAI OAuth 账号隐私，无论当前状态。
	ForceOpenAIPrivacy(ctx context.Context, account *Account) string
	// ForceAntigravityPrivacy 强制重新设置 Antigravity OAuth 账号隐私，无论当前状态。
	ForceAntigravityPrivacy(ctx context.Context, account *Account) string
	SetAccountSchedulable(ctx context.Context, id int64, schedulable bool) (*Account, error)
	BulkUpdateAccounts(ctx context.Context, input *BulkUpdateAccountsInput) (*BulkUpdateAccountsResult, error)
	CheckMixedChannelRisk(ctx context.Context, currentAccountID int64, currentAccountPlatform string, groupIDs []int64) error
	// RevertAccountProxyFallback 将账号的 proxy_id 切回 proxy_fallback_origin_id，并清空 origin 字段。
	// 若账号不存在返回 ErrAccountNotFound；若账号存在但不在 fallback 状态，返回 ErrAccountNotInFallback。
	RevertAccountProxyFallback(ctx context.Context, id int64) error

	// Proxy management
	ListProxies(ctx context.Context, page, pageSize int, protocol, status, search string, sortBy, sortOrder string) ([]Proxy, int64, error)
	ListProxiesWithAccountCount(ctx context.Context, page, pageSize int, protocol, status, search string, sortBy, sortOrder string) ([]ProxyWithAccountCount, int64, error)
	GetAllProxies(ctx context.Context) ([]Proxy, error)
	GetAllProxiesWithAccountCount(ctx context.Context) ([]ProxyWithAccountCount, error)
	GetProxy(ctx context.Context, id int64) (*Proxy, error)
	GetProxiesByIDs(ctx context.Context, ids []int64) ([]Proxy, error)
	CreateProxy(ctx context.Context, input *CreateProxyInput) (*Proxy, error)
	UpdateProxy(ctx context.Context, id int64, input *UpdateProxyInput) (*Proxy, error)
	DeleteProxy(ctx context.Context, id int64) error
	BatchDeleteProxies(ctx context.Context, ids []int64) (*ProxyBatchDeleteResult, error)
	GetProxyAccounts(ctx context.Context, proxyID int64) ([]ProxyAccountSummary, error)
	CheckProxyExists(ctx context.Context, host string, port int, username, password string) (bool, error)
	TestProxy(ctx context.Context, id int64) (*ProxyTestResult, error)
	CheckProxyQuality(ctx context.Context, id int64) (*ProxyQualityCheckResult, error)

	// Redeem code management
	ListRedeemCodes(ctx context.Context, page, pageSize int, codeType, status, search string, sortBy, sortOrder string) ([]RedeemCode, int64, error)
	GetRedeemCode(ctx context.Context, id int64) (*RedeemCode, error)
	GenerateRedeemCodes(ctx context.Context, input *GenerateRedeemCodesInput) ([]RedeemCode, error)
	DeleteRedeemCode(ctx context.Context, id int64) error
	BatchDeleteRedeemCodes(ctx context.Context, ids []int64) (int64, error)
	ExpireRedeemCode(ctx context.Context, id int64) (*RedeemCode, error)
	ResetAccountQuota(ctx context.Context, id int64) error
}

// CreateUserInput represents input for creating a new user via admin operations.
type CreateUserInput struct {
	Email         string
	Password      string
	Username      string
	Notes         string
	Role          string // 空字符串表示使用默认角色(user);合法值 admin/user
	Balance       *float64
	Concurrency   int
	RPMLimit      int
	AllowedGroups []int64
	// ActorAdminID 执行本次操作的管理员ID(来自JWT)，仅用于权限敏感操作的审计日志。
	ActorAdminID int64
}

type UpdateUserInput struct {
	Email         string
	Password      string
	Username      *string
	Notes         *string
	Role          string   // 空字符串表示"未提供"(不修改);合法值 admin/user
	Balance       *float64 // 使用指针区分"未提供"和"设置为0"
	Concurrency   *int     // 使用指针区分"未提供"和"设置为0"
	RPMLimit      *int     // 使用指针区分"未提供"和"设置为0"
	Status        string
	AllowedGroups *[]int64 // 使用指针区分"未提供"和"设置为空数组"
	// GroupRates 用户专属分组倍率配置
	// map[groupID]*rate，nil 表示删除该分组的专属倍率
	GroupRates map[int64]*float64
	// ActorAdminID 执行本次操作的管理员ID(来自JWT)，仅用于权限敏感操作的审计日志。
	ActorAdminID int64
}

type AdminBindAuthIdentityInput struct {
	ProviderType    string
	ProviderKey     string
	ProviderSubject string
	Issuer          *string
	Metadata        map[string]any
	Channel         *AdminBindAuthIdentityChannelInput
}

type AdminBindAuthIdentityChannelInput struct {
	Channel        string
	ChannelAppID   string
	ChannelSubject string
	Metadata       map[string]any
}

type AdminBoundAuthIdentity struct {
	UserID          int64                          `json:"user_id"`
	ProviderType    string                         `json:"provider_type"`
	ProviderKey     string                         `json:"provider_key"`
	ProviderSubject string                         `json:"provider_subject"`
	VerifiedAt      *time.Time                     `json:"verified_at,omitempty"`
	Issuer          *string                        `json:"issuer,omitempty"`
	Metadata        map[string]any                 `json:"metadata"`
	CreatedAt       time.Time                      `json:"created_at"`
	UpdatedAt       time.Time                      `json:"updated_at"`
	Channel         *AdminBoundAuthIdentityChannel `json:"channel,omitempty"`
}

type AdminBoundAuthIdentityChannel struct {
	Channel        string         `json:"channel"`
	ChannelAppID   string         `json:"channel_app_id"`
	ChannelSubject string         `json:"channel_subject"`
	Metadata       map[string]any `json:"metadata"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
}

type CreateGroupInput struct {
	Name             string
	Description      string
	Platform         string
	RateMultiplier   float64
	IsExclusive      bool
	SubscriptionType string   // standard/subscription
	DailyLimitUSD    *float64 // 日限额 (USD)
	WeeklyLimitUSD   *float64 // 周限额 (USD)
	MonthlyLimitUSD  *float64 // 月限额 (USD)
	// 图片生成计费配置（仅 antigravity 平台使用）
	AllowImageGeneration         bool
	AllowBatchImageGeneration    bool
	ImageRateIndependent         bool
	ImageRateMultiplier          *float64
	BatchImageDiscountMultiplier *float64
	BatchImageHoldMultiplier     *float64
	VideoRateIndependent         bool
	VideoRateMultiplier          *float64
	// 高峰时段倍率配置（PeakRateMultiplier 为 nil 时按 1.0 处理）
	PeakRateEnabled    bool
	PeakStart          string
	PeakEnd            string
	PeakRateMultiplier *float64
	ImagePrice1K       *float64
	ImagePrice2K       *float64
	ImagePrice4K       *float64
	VideoPrice480P     *float64
	VideoPrice720P     *float64
	VideoPrice1080P    *float64
	// Codex alpha/search 网页搜索单次价格（USD/次，仅 openai 平台使用）；nil/负数按默认价 0.01 处理
	WebSearchPricePerCall *float64
	ClaudeCodeOnly        bool   // 仅允许 Claude Code 客户端
	FallbackGroupID       *int64 // 降级分组 ID
	// 无效请求兜底分组 ID（仅 anthropic 平台使用）
	FallbackGroupIDOnInvalidRequest *int64
	// 模型路由配置（仅 anthropic 平台使用）
	ModelRouting        map[string][]int64
	ModelRoutingEnabled bool // 是否启用模型路由
	MCPXMLInject        *bool
	// 支持的模型系列（仅 antigravity 平台使用）
	SupportedModelScopes []string
	// OpenAI Messages 调度配置（仅 openai 平台使用）
	AllowMessagesDispatch       bool
	DefaultMappedModel          string
	RequireOAuthOnly            bool
	RequirePrivacySet           bool
	MessagesDispatchModelConfig OpenAIMessagesDispatchModelConfig
	ModelsListConfig            GroupModelsListConfig
	// RPMLimit 分组 RPM 上限（0 = 不限制）
	RPMLimit int
	// 从指定分组复制账号（创建分组后在同一事务内绑定）
	CopyAccountsFromGroupIDs []int64
}

type UpdateGroupInput struct {
	Name             string
	Description      *string
	Platform         string
	RateMultiplier   *float64 // 使用指针以支持设置为0
	IsExclusive      *bool
	Status           string
	SubscriptionType string   // standard/subscription
	DailyLimitUSD    *float64 // 日限额 (USD)
	WeeklyLimitUSD   *float64 // 周限额 (USD)
	MonthlyLimitUSD  *float64 // 月限额 (USD)
	// 图片生成计费配置（仅 antigravity 平台使用）
	AllowImageGeneration         *bool
	AllowBatchImageGeneration    *bool
	ImageRateIndependent         *bool
	ImageRateMultiplier          *float64
	BatchImageDiscountMultiplier *float64
	BatchImageHoldMultiplier     *float64
	VideoRateIndependent         *bool
	VideoRateMultiplier          *float64
	// 高峰时段倍率配置（nil 表示不修改）
	PeakRateEnabled    *bool
	PeakStart          *string
	PeakEnd            *string
	PeakRateMultiplier *float64
	ImagePrice1K       *float64
	ImagePrice2K       *float64
	ImagePrice4K       *float64
	VideoPrice480P     *float64
	VideoPrice720P     *float64
	VideoPrice1080P    *float64
	// Codex alpha/search 网页搜索单次价格（USD/次）；nil 表示不修改，负数表示清除回默认价 0.01
	WebSearchPricePerCall *float64
	ClaudeCodeOnly        *bool  // 仅允许 Claude Code 客户端
	FallbackGroupID       *int64 // 降级分组 ID
	// 无效请求兜底分组 ID（仅 anthropic 平台使用）
	FallbackGroupIDOnInvalidRequest *int64
	// 模型路由配置（仅 anthropic 平台使用）
	ModelRouting        map[string][]int64
	ModelRoutingEnabled *bool // 是否启用模型路由
	MCPXMLInject        *bool
	// 支持的模型系列（仅 antigravity 平台使用）
	SupportedModelScopes *[]string
	// OpenAI Messages 调度配置（仅 openai 平台使用）
	AllowMessagesDispatch       *bool
	DefaultMappedModel          *string
	RequireOAuthOnly            *bool
	RequirePrivacySet           *bool
	MessagesDispatchModelConfig *OpenAIMessagesDispatchModelConfig
	ModelsListConfig            *GroupModelsListConfig
	// RPMLimit 分组 RPM 上限（0 = 不限制），nil 表示未提供不改动。
	RPMLimit *int
	// 从指定分组复制账号（同步操作：先清空当前分组的账号绑定，再绑定源分组的账号）
	CopyAccountsFromGroupIDs []int64
}

type CreateAccountInput struct {
	Name               string
	Notes              *string
	Platform           string
	Type               string
	Credentials        map[string]any
	Extra              map[string]any
	ProxyID            *int64
	Concurrency        int
	Priority           int
	RateMultiplier     *float64 // 账号计费倍率（>=0，允许 0）
	LoadFactor         *int
	GroupIDs           []int64
	ExpiresAt          *int64
	AutoPauseOnExpired *bool
	// SkipDefaultGroupBind prevents auto-binding to platform default group when GroupIDs is empty.
	SkipDefaultGroupBind bool
	// SkipMixedChannelCheck skips the mixed channel risk check when binding groups.
	// This should only be set when the caller has explicitly confirmed the risk.
	SkipMixedChannelCheck bool
}

// ShadowOptions is the input for CreateShadow.
// The shadow holds no credentials; the scheduler delegates to the parent account.
type ShadowOptions struct {
	Name        string
	Priority    int
	Concurrency int
	GroupIDs    []int64
}

type UpdateAccountInput struct {
	Name                  string
	Notes                 *string
	Type                  string // Account type: oauth, setup-token, apikey
	Credentials           map[string]any
	Extra                 map[string]any
	ProxyID               *int64
	Concurrency           *int     // 使用指针区分"未提供"和"设置为0"
	Priority              *int     // 使用指针区分"未提供"和"设置为0"
	RateMultiplier        *float64 // 账号计费倍率（>=0，允许 0）
	LoadFactor            *int
	Status                string
	GroupIDs              *[]int64
	ExpiresAt             *int64
	AutoPauseOnExpired    *bool
	SkipMixedChannelCheck bool // 跳过混合渠道检查（用户已确认风险）
}

// BulkUpdateAccountsInput describes the payload for bulk updating accounts.
type BulkUpdateAccountsInput struct {
	AccountIDs     []int64
	Filters        *BulkUpdateAccountFilters
	Name           string
	ProxyID        *int64
	Concurrency    *int
	Priority       *int
	RateMultiplier *float64 // 账号计费倍率（>=0，允许 0）
	LoadFactor     *int
	Status         string
	Schedulable    *bool
	GroupIDs       *[]int64
	Credentials    map[string]any
	Extra          map[string]any
	// SkipMixedChannelCheck skips the mixed channel risk check when binding groups.
	// This should only be set when the caller has explicitly confirmed the risk.
	SkipMixedChannelCheck bool
}

type BulkUpdateAccountFilters struct {
	Platform      string
	Type          string
	Status        string
	AnomalyReason string
	Group         string
	Search        string
	PrivacyMode   string
}

// BulkUpdateAccountResult captures the result for a single account update.
type BulkUpdateAccountResult struct {
	AccountID int64  `json:"account_id"`
	Success   bool   `json:"success"`
	Error     string `json:"error,omitempty"`
}

// AdminUpdateAPIKeyGroupIDResult is the result of AdminUpdateAPIKeyGroupID.
type AdminUpdateAPIKeyGroupIDResult struct {
	APIKey                 *APIKey
	AutoGrantedGroupAccess bool   // true if a new exclusive group permission was auto-added
	GrantedGroupID         *int64 // the group ID that was auto-granted
	GrantedGroupName       string // the group name that was auto-granted
}

// ReplaceUserGroupResult 分组替换操作的结果
type ReplaceUserGroupResult struct {
	MigratedKeys int64 // 迁移的 Key 数量
}

// UserRPMStatus describes a user's current per-minute RPM usage.
type UserRPMStatus struct {
	UserRPMUsed  int                  `json:"user_rpm_used"`
	UserRPMLimit int                  `json:"user_rpm_limit"`
	PerGroup     []UserGroupRPMStatus `json:"per_group"`
}

// UserGroupRPMStatus describes current per-minute RPM usage for one user/group pair.
type UserGroupRPMStatus struct {
	GroupID   int64  `json:"group_id"`
	GroupName string `json:"group_name"`
	Used      int    `json:"used"`
	Limit     int    `json:"limit"`
	Source    string `json:"source"` // "group" | "override"
}

// BulkUpdateAccountsResult is the aggregated response for bulk updates.
type BulkUpdateAccountsResult struct {
	Success    int                       `json:"success"`
	Failed     int                       `json:"failed"`
	SuccessIDs []int64                   `json:"success_ids"`
	FailedIDs  []int64                   `json:"failed_ids"`
	Results    []BulkUpdateAccountResult `json:"results"`
}

type CreateProxyInput struct {
	Name           string
	Protocol       string
	Host           string
	Port           int
	Username       string
	Password       string
	ExpiresAt      *time.Time
	FallbackMode   string
	BackupProxyID  *int64
	ExpiryWarnDays int
}

type UpdateProxyInput struct {
	Name           string
	Protocol       string
	Host           string
	Port           int
	Username       string
	Password       string
	Status         string
	ExpiresAt      *time.Time
	FallbackMode   string
	BackupProxyID  *int64
	ExpiryWarnDays int
}

type GenerateRedeemCodesInput struct {
	Count        int
	Type         string
	Value        float64
	GroupID      *int64 // 订阅类型专用：关联的分组ID
	ValidityDays int    // 订阅类型专用：有效天数
	ExpiresAt    *time.Time
}

type ProxyBatchDeleteResult struct {
	DeletedIDs []int64                   `json:"deleted_ids"`
	Skipped    []ProxyBatchDeleteSkipped `json:"skipped"`
}

type ProxyBatchDeleteSkipped struct {
	ID     int64  `json:"id"`
	Reason string `json:"reason"`
}

// ProxyTestResult represents the result of testing a proxy
type ProxyTestResult struct {
	Success     bool   `json:"success"`
	Message     string `json:"message"`
	LatencyMs   int64  `json:"latency_ms,omitempty"`
	IPAddress   string `json:"ip_address,omitempty"`
	City        string `json:"city,omitempty"`
	Region      string `json:"region,omitempty"`
	Country     string `json:"country,omitempty"`
	CountryCode string `json:"country_code,omitempty"`
}

type ProxyQualityCheckResult struct {
	ProxyID        int64                   `json:"proxy_id"`
	Score          int                     `json:"score"`
	Grade          string                  `json:"grade"`
	Summary        string                  `json:"summary"`
	ExitIP         string                  `json:"exit_ip,omitempty"`
	Country        string                  `json:"country,omitempty"`
	CountryCode    string                  `json:"country_code,omitempty"`
	BaseLatencyMs  int64                   `json:"base_latency_ms,omitempty"`
	PassedCount    int                     `json:"passed_count"`
	WarnCount      int                     `json:"warn_count"`
	FailedCount    int                     `json:"failed_count"`
	ChallengeCount int                     `json:"challenge_count"`
	CheckedAt      int64                   `json:"checked_at"`
	Items          []ProxyQualityCheckItem `json:"items"`
}

type ProxyQualityCheckItem struct {
	Target     string `json:"target"`
	Status     string `json:"status"` // pass/warn/fail/challenge
	HTTPStatus int    `json:"http_status,omitempty"`
	LatencyMs  int64  `json:"latency_ms,omitempty"`
	Message    string `json:"message,omitempty"`
	CFRay      string `json:"cf_ray,omitempty"`
}

// ProxyExitInfo represents proxy exit information from ip-api.com
type ProxyExitInfo struct {
	IP          string
	City        string
	Region      string
	Country     string
	CountryCode string
}

// ProxyExitInfoProber tests proxy connectivity and retrieves exit information
type ProxyExitInfoProber interface {
	ProbeProxy(ctx context.Context, proxyURL string) (*ProxyExitInfo, int64, error)
}

type groupExistenceBatchReader interface {
	ExistsByIDs(ctx context.Context, ids []int64) (map[int64]bool, error)
}

type proxyQualityTarget struct {
	Target          string
	URL             string
	Method          string
	AllowedStatuses map[int]struct{}
}

var proxyQualityTargets = []proxyQualityTarget{
	{
		Target: "openai",
		URL:    "https://api.openai.com/v1/models",
		Method: http.MethodGet,
		AllowedStatuses: map[int]struct{}{
			http.StatusUnauthorized: {},
		},
	},
	{
		Target: "anthropic",
		URL:    "https://api.anthropic.com/v1/messages",
		Method: http.MethodGet,
		AllowedStatuses: map[int]struct{}{
			http.StatusUnauthorized:     {},
			http.StatusMethodNotAllowed: {},
			http.StatusNotFound:         {},
			http.StatusBadRequest:       {},
		},
	},
	{
		Target: "gemini",
		URL:    "https://generativelanguage.googleapis.com/$discovery/rest?version=v1beta",
		Method: http.MethodGet,
		AllowedStatuses: map[int]struct{}{
			http.StatusOK: {},
		},
	},
}

const (
	proxyQualityRequestTimeout        = 15 * time.Second
	proxyQualityResponseHeaderTimeout = 10 * time.Second
	proxyQualityMaxBodyBytes          = int64(8 * 1024)
	proxyQualityClientUserAgent       = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/136.0.0.0 Safari/537.36"
)

var ErrRPMStatusUnavailable = infraerrors.New(http.StatusNotImplemented, "RPM_STATUS_UNAVAILABLE", "RPM cache not available")

// adminServiceImpl implements AdminService
type adminServiceImpl struct {
	userRepo             UserRepository
	groupRepo            GroupRepository
	accountRepo          AdminAccountRepository
	upstreamPoolRepo     UpstreamPoolRepository
	accountDuplicateRepo AccountDuplicateRepository
	proxyRepo            ProxyRepository
	apiKeyRepo           APIKeyRepository
	redeemCodeRepo       RedeemCodeRepository
	userGroupRateRepo    UserGroupRateRepository
	userRPMCache         UserRPMCache
	billingCacheService  *BillingCacheService
	proxyProber          ProxyExitInfoProber
	proxyLatencyCache    ProxyLatencyCache
	authCacheInvalidator APIKeyAuthCacheInvalidator
	entClient            *dbent.Client // 用于开启数据库事务
	settingService       *SettingService
	defaultSubAssigner   DefaultSubscriptionAssigner
	userSubRepo          UserSubscriptionRepository
	privacyClientFactory PrivacyClientFactory
	openAIAccountRuntime OpenAIAccountRuntimeObserver
	runtimeBlocker       AccountRuntimeBlocker
	accountMonitorRepo   AccountMonitorRepository
	concurrencyService   *ConcurrencyService
}

type userGroupRateBatchReader interface {
	GetByUserIDs(ctx context.Context, userIDs []int64) (map[int64]map[int64]float64, error)
}

// NewAdminService creates a new AdminService
func NewAdminService(
	userRepo UserRepository,
	groupRepo GroupRepository,
	accountRepo AdminAccountRepository,
	upstreamPoolRepo UpstreamPoolRepository,
	proxyRepo ProxyRepository,
	apiKeyRepo APIKeyRepository,
	redeemCodeRepo RedeemCodeRepository,
	userGroupRateRepo UserGroupRateRepository,
	userRPMCache UserRPMCache,
	billingCacheService *BillingCacheService,
	proxyProber ProxyExitInfoProber,
	proxyLatencyCache ProxyLatencyCache,
	authCacheInvalidator APIKeyAuthCacheInvalidator,
	entClient *dbent.Client,
	settingService *SettingService,
	defaultSubAssigner DefaultSubscriptionAssigner,
	userSubRepo UserSubscriptionRepository,
	privacyClientFactory PrivacyClientFactory,
	openAIAccountRuntime OpenAIAccountRuntimeObserver,
	runtimeBlocker AccountRuntimeBlocker,
	concurrencyService *ConcurrencyService,
) AdminService {
	return &adminServiceImpl{
		userRepo:             userRepo,
		groupRepo:            groupRepo,
		accountRepo:          accountRepo,
		upstreamPoolRepo:     upstreamPoolRepo,
		accountDuplicateRepo: accountRepo,
		proxyRepo:            proxyRepo,
		apiKeyRepo:           apiKeyRepo,
		redeemCodeRepo:       redeemCodeRepo,
		userGroupRateRepo:    userGroupRateRepo,
		userRPMCache:         userRPMCache,
		billingCacheService:  billingCacheService,
		proxyProber:          proxyProber,
		proxyLatencyCache:    proxyLatencyCache,
		authCacheInvalidator: authCacheInvalidator,
		entClient:            entClient,
		settingService:       settingService,
		defaultSubAssigner:   defaultSubAssigner,
		userSubRepo:          userSubRepo,
		privacyClientFactory: privacyClientFactory,
		openAIAccountRuntime: openAIAccountRuntime,
		runtimeBlocker:       runtimeBlocker,
		concurrencyService:   concurrencyService,
	}
}
func (s *adminServiceImpl) SetAccountMonitorRepository(repo AccountMonitorRepository) {
	if s == nil {
		return
	}
	s.accountMonitorRepo = repo
}

func (s *adminServiceImpl) ListUpstreamPools(ctx context.Context) ([]UpstreamPool, error) {
	if s == nil || s.upstreamPoolRepo == nil {
		return []UpstreamPool{}, nil
	}
	pools, err := s.upstreamPoolRepo.ListUpstreamPools(ctx)
	if err != nil {
		return nil, err
	}
	return pools, nil
}

func (s *adminServiceImpl) GetUpstreamPoolByID(ctx context.Context, id int64) (*UpstreamPool, error) {
	if s == nil || s.upstreamPoolRepo == nil {
		return nil, ErrUpstreamPoolNotFound
	}
	return s.upstreamPoolRepo.GetUpstreamPoolByID(ctx, id)
}

func (s *adminServiceImpl) CreateUpstreamPool(ctx context.Context, input *CreateUpstreamPoolInput) (*UpstreamPool, error) {
	if s == nil || s.upstreamPoolRepo == nil || input == nil {
		return nil, ErrUpstreamPoolNotFound
	}

	pool := &UpstreamPool{
		Name:                           strings.TrimSpace(input.Name),
		Code:                           strings.TrimSpace(input.Code),
		Platform:                       strings.TrimSpace(input.Platform),
		Description:                    strings.TrimSpace(input.Description),
		Enabled:                        input.Enabled,
		SchedulerMode:                  strings.TrimSpace(input.SchedulerMode),
		AccountTypeStrategy:            strings.TrimSpace(input.AccountTypeStrategy),
		DefaultRequiredCapability:      strings.TrimSpace(input.DefaultRequiredCapability),
		DefaultRequiredTransport:       strings.TrimSpace(input.DefaultRequiredTransport),
		StickyEnabled:                  input.StickyEnabled,
		StickyTTLSeconds:               input.StickyTTLSeconds,
		StickyEscapeEnabled:            input.StickyEscapeEnabled,
		StickyEscapeErrorRateThreshold: input.StickyEscapeErrorRateThreshold,
		StickyEscapeTTFTMSThreshold:    input.StickyEscapeTTFTMSThreshold,
		LoadBalanceEnabled:             input.LoadBalanceEnabled,
		AutoWeightEnabled:              input.AutoWeightEnabled,
		AutoWeightMode:                 strings.TrimSpace(input.AutoWeightMode),
		FailoverEnabled:                input.FailoverEnabled,
		TopK:                           input.TopK,
		MaxFailoverHops:                input.MaxFailoverHops,
		WaitTimeoutMS:                  input.WaitTimeoutMS,
		MaxWaiting:                     input.MaxWaiting,
		PolicyJSON:                     input.PolicyJSON,
	}
	if err := normalizeUpstreamPoolForCreate(pool); err != nil {
		return nil, err
	}
	created, err := s.upstreamPoolRepo.CreateUpstreamPool(ctx, pool)
	if err != nil {
		return nil, err
	}
	return created, nil
}

func (s *adminServiceImpl) UpdateUpstreamPool(ctx context.Context, id int64, input *UpdateUpstreamPoolInput) (*UpstreamPool, error) {
	if s == nil || s.upstreamPoolRepo == nil || input == nil {
		return nil, ErrUpstreamPoolNotFound
	}
	pool, err := s.upstreamPoolRepo.GetUpstreamPoolByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if input.Name != nil {
		pool.Name = strings.TrimSpace(*input.Name)
	}
	if input.Code != nil {
		pool.Code = strings.TrimSpace(*input.Code)
	}
	if input.Platform != nil {
		nextPlatform := strings.TrimSpace(*input.Platform)
		if !strings.EqualFold(nextPlatform, pool.Platform) {
			if err := s.ensureUpstreamPoolPlatformMutable(ctx, pool.ID); err != nil {
				return nil, err
			}
		}
		pool.Platform = nextPlatform
	}
	if input.Description != nil {
		pool.Description = strings.TrimSpace(*input.Description)
	}
	if input.Enabled != nil {
		pool.Enabled = *input.Enabled
	}
	if input.SchedulerMode != nil {
		pool.SchedulerMode = strings.TrimSpace(*input.SchedulerMode)
	}
	if input.AccountTypeStrategy != nil {
		pool.AccountTypeStrategy = strings.TrimSpace(*input.AccountTypeStrategy)
	}
	if input.DefaultRequiredCapability != nil {
		pool.DefaultRequiredCapability = strings.TrimSpace(*input.DefaultRequiredCapability)
	}
	if input.DefaultRequiredTransport != nil {
		pool.DefaultRequiredTransport = strings.TrimSpace(*input.DefaultRequiredTransport)
	}
	if input.StickyEnabled != nil {
		pool.StickyEnabled = *input.StickyEnabled
	}
	if input.StickyTTLSeconds != nil {
		pool.StickyTTLSeconds = *input.StickyTTLSeconds
	}
	if input.StickyEscapeEnabled != nil {
		pool.StickyEscapeEnabled = *input.StickyEscapeEnabled
	}
	if input.StickyEscapeErrorRateThreshold != nil {
		pool.StickyEscapeErrorRateThreshold = *input.StickyEscapeErrorRateThreshold
	}
	if input.StickyEscapeTTFTMSThreshold != nil {
		pool.StickyEscapeTTFTMSThreshold = *input.StickyEscapeTTFTMSThreshold
	}
	if input.LoadBalanceEnabled != nil {
		pool.LoadBalanceEnabled = *input.LoadBalanceEnabled
	}
	if input.AutoWeightEnabled != nil {
		pool.AutoWeightEnabled = *input.AutoWeightEnabled
	}
	if input.AutoWeightMode != nil {
		pool.AutoWeightMode = strings.TrimSpace(*input.AutoWeightMode)
	}
	if input.FailoverEnabled != nil {
		pool.FailoverEnabled = *input.FailoverEnabled
	}
	if input.TopK != nil {
		pool.TopK = *input.TopK
	}
	if input.MaxFailoverHops != nil {
		pool.MaxFailoverHops = *input.MaxFailoverHops
	}
	if input.WaitTimeoutMS != nil {
		pool.WaitTimeoutMS = *input.WaitTimeoutMS
	}
	if input.MaxWaiting != nil {
		pool.MaxWaiting = *input.MaxWaiting
	}
	if input.PolicyJSON != nil {
		pool.PolicyJSON = input.PolicyJSON
	}
	if err := normalizeUpstreamPoolForCreate(pool); err != nil {
		return nil, err
	}
	updated, err := s.upstreamPoolRepo.UpdateUpstreamPool(ctx, pool)
	if err != nil {
		return nil, err
	}
	return updated, nil
}

func (s *adminServiceImpl) DeleteUpstreamPool(ctx context.Context, id int64) error {
	if s == nil || s.upstreamPoolRepo == nil {
		return ErrUpstreamPoolNotFound
	}
	return s.upstreamPoolRepo.DeleteUpstreamPool(ctx, id)
}

func (s *adminServiceImpl) ListUpstreamPoolMembers(ctx context.Context, poolID int64) ([]UpstreamPoolMember, error) {
	if s == nil || s.upstreamPoolRepo == nil {
		return []UpstreamPoolMember{}, nil
	}
	members, err := s.upstreamPoolRepo.ListUpstreamPoolMembers(ctx, poolID)
	if err != nil {
		return nil, err
	}
	if s.accountRepo == nil || len(members) == 0 {
		return members, nil
	}
	accountIDs := make([]int64, 0, len(members))
	for i := range members {
		if members[i].AccountID > 0 {
			accountIDs = append(accountIDs, members[i].AccountID)
		}
	}
	accountMap, runtimeSnapshots, monitorLatest, err := s.loadUpstreamAccountRuntimeState(ctx, accountIDs)
	if err != nil {
		return nil, err
	}
	for i := range members {
		account := accountMap[members[i].AccountID]
		if account == nil {
			continue
		}
		s.applyUpstreamAccountRuntimeToPoolMember(&members[i], account, runtimeSnapshots[members[i].AccountID], monitorLatest[members[i].AccountID])
	}
	return members, nil
}

func (s *adminServiceImpl) PreviewUpstreamPoolMemberSync(ctx context.Context, poolID int64, input *UpstreamPoolMemberSyncPreviewInput) (*UpstreamPoolMemberSyncResult, error) {
	if s == nil || s.upstreamPoolRepo == nil || s.accountRepo == nil || input == nil {
		return nil, ErrUpstreamPoolNotFound
	}
	mode := NormalizeUpstreamPoolMemberSyncMode(input.Mode)
	pool, err := s.upstreamPoolRepo.GetUpstreamPoolByID(ctx, poolID)
	if err != nil {
		return nil, err
	}
	targets, err := s.buildUpstreamPoolMemberSyncTargets(ctx, pool)
	if err != nil {
		return nil, err
	}
	members, err := s.upstreamPoolRepo.ListUpstreamPoolMembers(ctx, pool.ID)
	if err != nil {
		return nil, err
	}
	return BuildUpstreamPoolMemberSyncResult(pool, members, targets, mode), nil
}

func (s *adminServiceImpl) ApplyUpstreamPoolMemberSync(ctx context.Context, poolID int64, input *UpstreamPoolMemberSyncApplyInput) (*UpstreamPoolMemberSyncResult, error) {
	if s == nil || s.upstreamPoolRepo == nil || s.accountRepo == nil || input == nil {
		return nil, ErrUpstreamPoolNotFound
	}
	syncer, ok := s.upstreamPoolRepo.(UpstreamPoolMemberSyncer)
	if !ok {
		return nil, infraerrors.InternalServer("UPSTREAM_POOL_SYNC_UNAVAILABLE", "upstream pool member sync is not available")
	}
	mode := NormalizeUpstreamPoolMemberSyncMode(input.Mode)
	pool, err := s.upstreamPoolRepo.GetUpstreamPoolByID(ctx, poolID)
	if err != nil {
		return nil, err
	}
	targets, err := s.buildUpstreamPoolMemberSyncTargets(ctx, pool)
	if err != nil {
		return nil, err
	}
	return syncer.SyncUpstreamPoolDirectMembers(ctx, pool.ID, targets, mode)
}

func (s *adminServiceImpl) buildUpstreamPoolMemberSyncTargets(ctx context.Context, pool *UpstreamPool) ([]UpstreamPoolMember, error) {
	if pool == nil {
		return nil, ErrUpstreamPoolNotFound
	}
	accounts, err := s.listUpstreamPoolMemberSyncAccounts(ctx, pool)
	if err != nil {
		return nil, err
	}
	targets := make([]UpstreamPoolMember, 0, len(accounts))
	for i := range accounts {
		account := accounts[i]
		if !strings.EqualFold(strings.TrimSpace(account.Platform), strings.TrimSpace(pool.Platform)) {
			continue
		}
		if !accountMatchesUpstreamPoolMemberSyncStrategy(&account, pool.AccountTypeStrategy) {
			continue
		}
		enabled := account.IsActive() && account.IsSchedulable()
		targets = append(targets, UpstreamPoolMember{
			PoolID:          pool.ID,
			AccountID:       account.ID,
			AccountName:     account.Name,
			AccountPlatform: account.Platform,
			AccountType:     account.Type,
			AccountStatus:   account.Status,
			Enabled:         enabled,
			ManualDrained:   !enabled,
			Weight:          100,
			Notes:           "从账号管理同步",
		})
	}
	return targets, nil
}

func (s *adminServiceImpl) listUpstreamPoolMemberSyncAccounts(ctx context.Context, pool *UpstreamPool) ([]Account, error) {
	groupIDs, err := s.enabledUpstreamPoolBindingGroupIDs(ctx, pool.ID)
	if err != nil {
		return nil, err
	}
	if len(groupIDs) == 0 {
		accounts, _, err := s.accountRepo.ListWithFilters(
			ctx,
			pagination.PaginationParams{Page: 1, PageSize: 10000, SortBy: "id", SortOrder: "ASC"},
			pool.Platform,
			"",
			"",
			"",
			"",
			0,
			"",
		)
		return accounts, err
	}

	seen := make(map[int64]struct{})
	accounts := make([]Account, 0)
	for _, groupID := range groupIDs {
		groupAccounts, err := s.accountRepo.ListByGroup(ctx, groupID)
		if err != nil {
			return nil, err
		}
		for _, account := range groupAccounts {
			if _, ok := seen[account.ID]; ok {
				continue
			}
			seen[account.ID] = struct{}{}
			accounts = append(accounts, account)
		}
	}
	sort.SliceStable(accounts, func(i, j int) bool {
		return accounts[i].ID < accounts[j].ID
	})
	return accounts, nil
}

func (s *adminServiceImpl) enabledUpstreamPoolBindingGroupIDs(ctx context.Context, poolID int64) ([]int64, error) {
	if s == nil || s.upstreamPoolRepo == nil {
		return nil, nil
	}
	bindings, err := s.upstreamPoolRepo.ListUpstreamPoolBindings(ctx)
	if err != nil {
		return nil, err
	}
	return enabledUpstreamPoolBindingGroupIDsFromList(poolID, bindings), nil
}

func accountMatchesUpstreamPoolMemberSyncStrategy(account *Account, strategy string) bool {
	if account == nil {
		return false
	}
	if NormalizeUpstreamPoolAccountTypeStrategy(strategy) == UpstreamPoolAccountTypeStrategyOAuthOnly {
		return account.Type == AccountTypeOAuth || account.Type == AccountTypeSetupToken
	}
	return true
}

func enabledUpstreamPoolBindingGroupIDsFromList(poolID int64, bindings []UpstreamPoolBinding) []int64 {
	seen := make(map[int64]struct{})
	groupIDs := make([]int64, 0)
	for _, binding := range bindings {
		if binding.PoolID != poolID || !binding.Enabled || binding.GroupID <= 0 {
			continue
		}
		if _, ok := seen[binding.GroupID]; ok {
			continue
		}
		seen[binding.GroupID] = struct{}{}
		groupIDs = append(groupIDs, binding.GroupID)
	}
	sort.Slice(groupIDs, func(i, j int) bool { return groupIDs[i] < groupIDs[j] })
	return groupIDs
}

func (s *adminServiceImpl) ListUpstreamAccountSetMembers(ctx context.Context, setID int64) ([]UpstreamAccountSetMember, error) {
	if s == nil || s.upstreamPoolRepo == nil {
		return []UpstreamAccountSetMember{}, nil
	}
	members, err := s.upstreamPoolRepo.ListUpstreamAccountSetMembers(ctx, setID)
	if err != nil {
		return nil, err
	}
	if s.accountRepo == nil || len(members) == 0 {
		return members, nil
	}
	accountIDs := make([]int64, 0, len(members))
	for i := range members {
		if members[i].AccountID > 0 {
			accountIDs = append(accountIDs, members[i].AccountID)
		}
	}
	accountMap, runtimeSnapshots, monitorLatest, err := s.loadUpstreamAccountRuntimeState(ctx, accountIDs)
	if err != nil {
		return nil, err
	}
	for i := range members {
		account := accountMap[members[i].AccountID]
		if account == nil {
			continue
		}
		s.applyUpstreamAccountRuntimeToAccountSetMember(&members[i], account, runtimeSnapshots[members[i].AccountID], monitorLatest[members[i].AccountID])
	}
	return members, nil
}

func (s *adminServiceImpl) loadUpstreamAccountRuntimeState(ctx context.Context, accountIDs []int64) (map[int64]*Account, map[int64]OpenAIAccountRuntimeSnapshot, map[int64][]*AccountMonitorLatest, error) {
	accounts, err := s.accountRepo.GetByIDs(ctx, accountIDs)
	if err != nil {
		return nil, nil, nil, err
	}
	accountMap := make(map[int64]*Account, len(accounts))
	for _, account := range accounts {
		if account == nil {
			continue
		}
		accountMap[account.ID] = account
	}
	runtimeSnapshots := map[int64]OpenAIAccountRuntimeSnapshot{}
	if s.openAIAccountRuntime != nil && len(accountIDs) > 0 {
		runtimeSnapshots = s.openAIAccountRuntime.SnapshotOpenAIAccountRuntime(accountIDs)
	}
	monitorLatest := map[int64][]*AccountMonitorLatest{}
	if s.accountMonitorRepo != nil && len(accountIDs) > 0 {
		if latest, err := s.accountMonitorRepo.ListLatestForAccountIDs(ctx, accountIDs); err == nil {
			monitorLatest = latest
		}
	}
	return accountMap, runtimeSnapshots, monitorLatest, nil
}

func (s *adminServiceImpl) applyUpstreamAccountRuntimeToPoolMember(member *UpstreamPoolMember, account *Account, snapshot OpenAIAccountRuntimeSnapshot, latest []*AccountMonitorLatest) {
	if member == nil || account == nil {
		return
	}
	member.AccountName = account.Name
	member.AccountPlatform = account.Platform
	member.AccountType = account.Type
	member.AccountStatus = account.Status
	member.AccountSchedulable = account.IsSchedulable()
	member.RuntimeLastUsedAt = account.LastUsedAt
	member.RuntimeRateLimitResetAt = account.RateLimitResetAt
	member.RuntimeOverloadUntil = account.OverloadUntil
	member.RuntimeTempUnschedulableUntil = account.TempUnschedulableUntil
	member.RuntimeStatus, member.RuntimeReason = summarizeUpstreamPoolMemberRuntime(account)
	member.RuntimeStatus, member.RuntimeReason = summarizeUpstreamPoolMemberRuntimeWithMonitor(
		member.RuntimeStatus,
		member.RuntimeReason,
		latest,
	)
	if snapshot.TTFTMs != nil || snapshot.ErrorRate > 0 {
		errorRate := snapshot.ErrorRate
		member.RuntimeErrorRate = &errorRate
		member.RuntimeTTFTMs = snapshot.TTFTMs
	}
}

func (s *adminServiceImpl) applyUpstreamAccountRuntimeToAccountSetMember(member *UpstreamAccountSetMember, account *Account, snapshot OpenAIAccountRuntimeSnapshot, latest []*AccountMonitorLatest) {
	if member == nil || account == nil {
		return
	}
	member.AccountName = account.Name
	member.AccountPlatform = account.Platform
	member.AccountType = account.Type
	member.AccountStatus = account.Status
	member.AccountSchedulable = account.IsSchedulable()
	member.RuntimeLastUsedAt = account.LastUsedAt
	member.RuntimeRateLimitResetAt = account.RateLimitResetAt
	member.RuntimeOverloadUntil = account.OverloadUntil
	member.RuntimeTempUnschedulableUntil = account.TempUnschedulableUntil
	member.RuntimeStatus, member.RuntimeReason = summarizeUpstreamPoolMemberRuntime(account)
	member.RuntimeStatus, member.RuntimeReason = summarizeUpstreamPoolMemberRuntimeWithMonitor(
		member.RuntimeStatus,
		member.RuntimeReason,
		latest,
	)
	if snapshot.TTFTMs != nil || snapshot.ErrorRate > 0 {
		errorRate := snapshot.ErrorRate
		member.RuntimeErrorRate = &errorRate
		member.RuntimeTTFTMs = snapshot.TTFTMs
	}
}

func summarizeUpstreamPoolMemberRuntime(account *Account) (string, string) {
	if account == nil {
		return "unknown", ""
	}
	now := time.Now()
	if account.Status == StatusError {
		reason := strings.TrimSpace(account.ErrorMessage)
		if reason == "" {
			reason = "账号错误，等待恢复探针"
		}
		return "error_recovering", reason
	}
	if account.RateLimitResetAt != nil && now.Before(*account.RateLimitResetAt) {
		return "rate_limited", "限流中，等待窗口恢复"
	}
	if account.OverloadUntil != nil && now.Before(*account.OverloadUntil) {
		return "overloaded", "高负载保护中"
	}
	if account.TempUnschedulableUntil != nil && now.Before(*account.TempUnschedulableUntil) {
		reason := strings.TrimSpace(account.TempUnschedulableReason)
		if reason == "" {
			reason = "临时不可调度，等待恢复探针"
		}
		return "temp_unschedulable", reason
	}
	if !account.Schedulable {
		return "disabled", "账号自身已关闭调度"
	}
	if account.IsSchedulable() {
		return "healthy", "正常参与调度"
	}
	return "degraded", "当前未进入可调度状态"
}

func summarizeUpstreamPoolMemberRuntimeWithMonitor(baseStatus, baseReason string, latest []*AccountMonitorLatest) (string, string) {
	if len(latest) == 0 {
		return baseStatus, baseReason
	}
	status := ""
	kinds := make([]string, 0, len(latest))
	for _, row := range latest {
		if row == nil {
			continue
		}
		status = worseUpstreamPoolMemberMonitorStatus(status, row.Status)
		kinds = append(kinds, upstreamPoolMemberMonitorKind(row.Provider, row.Model, row.Status))
	}
	if status == "" || status == MonitorStatusOperational {
		if baseStatus == "healthy" && len(kinds) > 0 {
			return baseStatus, "最近探针正常：" + strings.Join(kinds, " / ")
		}
		return baseStatus, baseReason
	}
	reason := strings.TrimSpace("最近探针未通过：" + strings.Join(kinds, " / "))
	if reason == "最近探针未通过：" {
		reason = "最近探针未通过"
	}
	switch status {
	case MonitorStatusFailed, MonitorStatusError:
		return "error_recovering", reason
	case MonitorStatusDegraded:
		if baseStatus == "healthy" {
			return "degraded", reason
		}
		return baseStatus, firstNonEmptyTrimmed(baseReason, reason)
	default:
		return baseStatus, baseReason
	}
}

func worseUpstreamPoolMemberMonitorStatus(current, next string) string {
	if monitorStatusRank(next) > monitorStatusRank(current) {
		return next
	}
	return current
}

func upstreamPoolMemberMonitorKind(provider, model, status string) string {
	label := strings.TrimSpace(model)
	switch strings.ToLower(strings.TrimSpace(provider)) {
	case PlatformAnthropic:
		label = "Claude Code " + label
	case PlatformOpenAI:
		label = "OpenAI " + label
	default:
		label = strings.TrimSpace(provider + " " + label)
	}
	if strings.TrimSpace(status) != "" {
		label += "=" + strings.TrimSpace(status)
	}
	return strings.TrimSpace(label)
}

func (s *adminServiceImpl) CreateUpstreamPoolMember(ctx context.Context, poolID int64, input *CreateUpstreamPoolMemberInput) (*UpstreamPoolMember, error) {
	if s == nil || s.upstreamPoolRepo == nil || s.accountRepo == nil || input == nil {
		return nil, ErrUpstreamPoolNotFound
	}
	pool, err := s.upstreamPoolRepo.GetUpstreamPoolByID(ctx, poolID)
	if err != nil {
		return nil, err
	}
	account, err := s.accountRepo.GetByID(ctx, input.AccountID)
	if err != nil {
		return nil, err
	}
	if !strings.EqualFold(strings.TrimSpace(account.Platform), strings.TrimSpace(pool.Platform)) {
		return nil, errors.New("account platform does not match pool platform")
	}
	member := &UpstreamPoolMember{
		PoolID:                 pool.ID,
		AccountID:              account.ID,
		AccountName:            account.Name,
		AccountPlatform:        account.Platform,
		Enabled:                input.Enabled,
		SchedulableOverride:    input.SchedulableOverride,
		ManualDrained:          input.ManualDrained,
		Weight:                 input.Weight,
		PriorityOverride:       input.PriorityOverride,
		MaxConcurrencyOverride: input.MaxConcurrencyOverride,
		Notes:                  strings.TrimSpace(input.Notes),
	}
	if err := normalizeUpstreamPoolMemberForCreate(member); err != nil {
		return nil, err
	}
	created, err := s.upstreamPoolRepo.CreateUpstreamPoolMember(ctx, member)
	if err != nil {
		return nil, err
	}
	created.AccountName = account.Name
	created.AccountPlatform = account.Platform
	return created, nil
}

func (s *adminServiceImpl) UpdateUpstreamPoolMember(ctx context.Context, id int64, input *UpdateUpstreamPoolMemberInput) (*UpstreamPoolMember, error) {
	if s == nil || s.upstreamPoolRepo == nil || input == nil {
		return nil, ErrUpstreamPoolNotFound
	}
	member, err := s.upstreamPoolRepo.GetUpstreamPoolMemberByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if input.Enabled != nil {
		member.Enabled = *input.Enabled
	}
	if input.SchedulableOverrideSet {
		member.SchedulableOverride = input.SchedulableOverride
	}
	if input.ManualDrained != nil {
		member.ManualDrained = *input.ManualDrained
	}
	if input.Weight != nil {
		member.Weight = *input.Weight
	}
	if input.PriorityOverrideSet {
		member.PriorityOverride = input.PriorityOverride
	}
	if input.MaxConcurrencyOverrideSet {
		member.MaxConcurrencyOverride = input.MaxConcurrencyOverride
	}
	if input.NotesSet {
		if input.Notes == nil {
			member.Notes = ""
		} else {
			member.Notes = strings.TrimSpace(*input.Notes)
		}
	}
	if err := normalizeUpstreamPoolMemberForCreate(member); err != nil {
		return nil, err
	}
	updated, err := s.upstreamPoolRepo.UpdateUpstreamPoolMember(ctx, member)
	if err != nil {
		return nil, err
	}
	return updated, nil
}

func (s *adminServiceImpl) DeleteUpstreamPoolMember(ctx context.Context, id int64) error {
	if s == nil || s.upstreamPoolRepo == nil {
		return ErrUpstreamPoolNotFound
	}
	return s.upstreamPoolRepo.DeleteUpstreamPoolMember(ctx, id)
}

func (s *adminServiceImpl) ListUpstreamAccountSets(ctx context.Context) ([]UpstreamAccountSet, error) {
	if s == nil || s.upstreamPoolRepo == nil {
		return []UpstreamAccountSet{}, nil
	}
	return s.upstreamPoolRepo.ListUpstreamAccountSets(ctx)
}

func (s *adminServiceImpl) ListUpstreamCapacityPressures(ctx context.Context) ([]UpstreamCapacityPressure, error) {
	if s == nil || s.upstreamPoolRepo == nil {
		return []UpstreamCapacityPressure{}, nil
	}
	sets, err := s.upstreamPoolRepo.ListUpstreamAccountSets(ctx)
	if err != nil {
		return nil, err
	}
	since := time.Now().Add(-5 * time.Minute)
	setIDs := make([]int64, 0, len(sets))
	for _, set := range sets {
		if set.SharedConcurrencyLimit != nil && *set.SharedConcurrencyLimit > 0 {
			setIDs = append(setIDs, set.ID)
		}
	}
	snapshotStats := map[int64]UpstreamCapacitySnapshotStats{}
	if s.accountMonitorRepo != nil {
		if values, statsErr := s.accountMonitorRepo.ListUpstreamCapacitySnapshotStats(ctx, setIDs, since); statsErr == nil {
			snapshotStats = values
		}
	}
	out := make([]UpstreamCapacityPressure, 0)
	for _, set := range sets {
		if set.SharedConcurrencyLimit == nil || *set.SharedConcurrencyLimit <= 0 {
			continue
		}
		members, err := s.upstreamPoolRepo.ListUpstreamAccountSetMembers(ctx, set.ID)
		if err != nil {
			return nil, err
		}
		pressure := UpstreamCapacityPressure{
			SetID:         set.ID,
			SetName:       set.Name,
			SetCode:       set.Code,
			Platform:      set.Platform,
			Enabled:       set.Enabled,
			CapacityLimit: *set.SharedConcurrencyLimit,
			Members:       make([]UpstreamCapacityMemberPressure, 0, len(members)),
		}
		for _, member := range members {
			if member.CapacityHardLimit == nil && member.CapacitySoftShare == nil {
				continue
			}
			memberPressure := UpstreamCapacityMemberPressure{
				AccountID:            member.AccountID,
				AccountName:          member.AccountName,
				HardConcurrencyLimit: clonePositiveIntPointer(member.CapacityHardLimit),
				SoftConcurrencyShare: clonePositiveIntPointer(member.CapacitySoftShare),
			}
			if s.concurrencyService != nil {
				counts, countErr := s.concurrencyService.GetCapacitySlotCounts(ctx, set.ID, member.AccountID)
				if countErr != nil {
					return nil, countErr
				}
				if counts.GroupConcurrency > pressure.CurrentConcurrency {
					pressure.CurrentConcurrency = counts.GroupConcurrency
				}
				memberPressure.CurrentConcurrency = counts.MemberConcurrency
				waiting, waitErr := s.concurrencyService.GetAccountWaitingCount(ctx, member.AccountID)
				if waitErr != nil {
					return nil, waitErr
				}
				memberPressure.WaitingCount = waiting
				pressure.WaitingCount += waiting
			}
			denominator := pressure.CapacityLimit
			if memberPressure.HardConcurrencyLimit != nil && *memberPressure.HardConcurrencyLimit > 0 {
				denominator = *memberPressure.HardConcurrencyLimit
			} else if memberPressure.SoftConcurrencyShare != nil && *memberPressure.SoftConcurrencyShare > 0 {
				denominator = *memberPressure.SoftConcurrencyShare
			}
			if denominator > 0 {
				memberPressure.LoadRate = (memberPressure.CurrentConcurrency + memberPressure.WaitingCount) * 100 / denominator
			}
			pressure.Members = append(pressure.Members, memberPressure)
		}
		if s.concurrencyService != nil {
			metrics, metricsErr := s.concurrencyService.GetCapacityMetrics(ctx, set.ID, since)
			if metricsErr != nil {
				return nil, metricsErr
			}
			pressure.GroupFullCount = metrics.GroupFullCount
			pressure.MemberFullCount = metrics.MemberFullCount
			pressure.BorrowedSlotCount = metrics.BorrowedSlotCount
		}
		pressure.AvailableCapacity = pressure.CapacityLimit - pressure.CurrentConcurrency
		if pressure.AvailableCapacity < 0 {
			pressure.AvailableCapacity = 0
		}
		pressure.PeakConcurrency5m = snapshotStats[set.ID].PeakConcurrency5m
		pressure.P95LoadRate5m = snapshotStats[set.ID].P95LoadRate5m
		maxMemberConcurrency := 0
		for _, member := range pressure.Members {
			if member.CurrentConcurrency > maxMemberConcurrency {
				maxMemberConcurrency = member.CurrentConcurrency
			}
		}
		if pressure.CurrentConcurrency > 0 {
			pressure.SchedulingConcentration = maxMemberConcurrency * 100 / pressure.CurrentConcurrency
		}
		out = append(out, pressure)
	}
	return out, nil
}

func (s *adminServiceImpl) CreateUpstreamAccountSet(ctx context.Context, input *CreateUpstreamAccountSetInput) (*UpstreamAccountSet, error) {
	if s == nil || s.upstreamPoolRepo == nil || input == nil {
		return nil, ErrUpstreamPoolNotFound
	}
	item := &UpstreamAccountSet{
		Name:                   strings.TrimSpace(input.Name),
		Code:                   strings.TrimSpace(input.Code),
		Platform:               strings.TrimSpace(input.Platform),
		Description:            strings.TrimSpace(input.Description),
		Enabled:                input.Enabled,
		SharedConcurrencyLimit: clonePositiveIntPointer(input.SharedConcurrencyLimit),
	}
	if item.Code == "" {
		item.Code = buildAutoUpstreamAccountSetCode(item.Platform, item.Name)
	}
	if err := normalizeUpstreamAccountSetForCreate(item); err != nil {
		return nil, err
	}
	return s.upstreamPoolRepo.CreateUpstreamAccountSet(ctx, item)
}

func (s *adminServiceImpl) UpdateUpstreamAccountSet(ctx context.Context, id int64, input *UpdateUpstreamAccountSetInput) (*UpstreamAccountSet, error) {
	if s == nil || s.upstreamPoolRepo == nil || input == nil {
		return nil, ErrUpstreamPoolNotFound
	}
	item, err := s.upstreamPoolRepo.GetUpstreamAccountSetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if input.Name != nil {
		item.Name = strings.TrimSpace(*input.Name)
	}
	if input.Code != nil {
		item.Code = strings.TrimSpace(*input.Code)
	}
	if input.Platform != nil {
		nextPlatform := strings.TrimSpace(*input.Platform)
		if !strings.EqualFold(nextPlatform, item.Platform) {
			if err := s.ensureUpstreamAccountSetPlatformMutable(ctx, item.ID); err != nil {
				return nil, err
			}
		}
		item.Platform = nextPlatform
	}
	if input.Description != nil {
		item.Description = strings.TrimSpace(*input.Description)
	}
	if input.Enabled != nil {
		item.Enabled = *input.Enabled
	}
	if input.SharedConcurrencyLimitSet {
		item.SharedConcurrencyLimit = clonePositiveIntPointer(input.SharedConcurrencyLimit)
	}
	if err := normalizeUpstreamAccountSetForCreate(item); err != nil {
		return nil, err
	}
	return s.upstreamPoolRepo.UpdateUpstreamAccountSet(ctx, item)
}

func (s *adminServiceImpl) DeleteUpstreamAccountSet(ctx context.Context, id int64) error {
	if s == nil || s.upstreamPoolRepo == nil {
		return ErrUpstreamPoolNotFound
	}
	return s.upstreamPoolRepo.DeleteUpstreamAccountSet(ctx, id)
}

func (s *adminServiceImpl) AddUpstreamAccountSetMembers(ctx context.Context, setID int64, input *AddUpstreamAccountSetMembersInput) error {
	if s == nil || s.upstreamPoolRepo == nil || s.accountRepo == nil || input == nil {
		return ErrUpstreamPoolNotFound
	}
	item, err := s.upstreamPoolRepo.GetUpstreamAccountSetByID(ctx, setID)
	if err != nil {
		return err
	}
	accountIDs := uniquePositiveInt64s(input.AccountIDs)
	if len(accountIDs) == 0 {
		return errors.New("account_ids is required")
	}
	accounts, err := s.accountRepo.GetByIDs(ctx, accountIDs)
	if err != nil {
		return err
	}
	if len(accounts) != len(accountIDs) {
		return errors.New("some accounts were not found")
	}
	for _, account := range accounts {
		if account == nil {
			return errors.New("some accounts were not found")
		}
		if !strings.EqualFold(strings.TrimSpace(account.Platform), strings.TrimSpace(item.Platform)) {
			return errors.New("account platform does not match account set platform")
		}
	}
	return s.upstreamPoolRepo.AddUpstreamAccountSetMembers(ctx, setID, accountIDs)
}

func (s *adminServiceImpl) DeleteUpstreamAccountSetMember(ctx context.Context, setID, accountID int64) error {
	if s == nil || s.upstreamPoolRepo == nil {
		return ErrUpstreamPoolNotFound
	}
	return s.upstreamPoolRepo.DeleteUpstreamAccountSetMember(ctx, setID, accountID)
}

func (s *adminServiceImpl) UpdateUpstreamAccountSetMemberCapacity(ctx context.Context, setID, accountID int64, hardLimit, softShare *int) error {
	if s == nil || s.upstreamPoolRepo == nil {
		return ErrUpstreamPoolNotFound
	}
	writer, ok := s.upstreamPoolRepo.(interface {
		UpdateUpstreamAccountSetMemberCapacity(context.Context, int64, int64, *int, *int) error
	})
	if !ok {
		return errors.New("capacity member configuration is unavailable")
	}
	return writer.UpdateUpstreamAccountSetMemberCapacity(ctx, setID, accountID, hardLimit, softShare)
}

func (s *adminServiceImpl) ListUpstreamPoolMemberSets(ctx context.Context, poolID int64) ([]UpstreamPoolMemberSet, error) {
	if s == nil || s.upstreamPoolRepo == nil {
		return []UpstreamPoolMemberSet{}, nil
	}
	return s.upstreamPoolRepo.ListUpstreamPoolMemberSets(ctx, poolID)
}

func (s *adminServiceImpl) CreateUpstreamPoolMemberSet(ctx context.Context, poolID int64, input *CreateUpstreamPoolMemberSetInput) (*UpstreamPoolMemberSet, error) {
	if s == nil || s.upstreamPoolRepo == nil || input == nil {
		return nil, ErrUpstreamPoolNotFound
	}
	pool, err := s.upstreamPoolRepo.GetUpstreamPoolByID(ctx, poolID)
	if err != nil {
		return nil, err
	}
	setItem, err := s.upstreamPoolRepo.GetUpstreamAccountSetByID(ctx, input.SetID)
	if err != nil {
		return nil, err
	}
	if !strings.EqualFold(strings.TrimSpace(pool.Platform), strings.TrimSpace(setItem.Platform)) {
		return nil, errors.New("account set platform does not match pool platform")
	}
	item := &UpstreamPoolMemberSet{
		PoolID:      pool.ID,
		SetID:       setItem.ID,
		SetName:     setItem.Name,
		SetCode:     setItem.Code,
		SetPlatform: setItem.Platform,
		Enabled:     input.Enabled,
		Notes:       strings.TrimSpace(input.Notes),
	}
	if err := normalizeUpstreamPoolMemberSetForCreate(item); err != nil {
		return nil, err
	}
	return s.upstreamPoolRepo.CreateUpstreamPoolMemberSet(ctx, item)
}

func (s *adminServiceImpl) UpdateUpstreamPoolMemberSet(ctx context.Context, id int64, input *UpdateUpstreamPoolMemberSetInput) (*UpstreamPoolMemberSet, error) {
	if s == nil || s.upstreamPoolRepo == nil || input == nil {
		return nil, ErrUpstreamPoolNotFound
	}
	item, err := s.upstreamPoolRepo.GetUpstreamPoolMemberSetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if input.Enabled != nil {
		item.Enabled = *input.Enabled
	}
	if input.Notes != nil {
		item.Notes = strings.TrimSpace(*input.Notes)
	}
	if err := normalizeUpstreamPoolMemberSetForCreate(item); err != nil {
		return nil, err
	}
	return s.upstreamPoolRepo.UpdateUpstreamPoolMemberSet(ctx, item)
}

func (s *adminServiceImpl) DeleteUpstreamPoolMemberSet(ctx context.Context, id int64) error {
	if s == nil || s.upstreamPoolRepo == nil {
		return ErrUpstreamPoolNotFound
	}
	return s.upstreamPoolRepo.DeleteUpstreamPoolMemberSet(ctx, id)
}

func (s *adminServiceImpl) ListUpstreamPoolBindings(ctx context.Context) ([]UpstreamPoolBinding, error) {
	if s == nil || s.upstreamPoolRepo == nil {
		return []UpstreamPoolBinding{}, nil
	}
	bindings, err := s.upstreamPoolRepo.ListUpstreamPoolBindings(ctx)
	if err != nil {
		return nil, err
	}
	return bindings, nil
}

func (s *adminServiceImpl) CreateUpstreamPoolBinding(ctx context.Context, input *CreateUpstreamPoolBindingInput) (*UpstreamPoolBinding, error) {
	if s == nil || s.upstreamPoolRepo == nil || s.groupRepo == nil || input == nil {
		return nil, ErrUpstreamPoolNotFound
	}
	group, err := s.groupRepo.GetByIDLite(ctx, input.GroupID)
	if err != nil {
		return nil, err
	}
	pool, err := s.upstreamPoolRepo.GetUpstreamPoolByID(ctx, input.PoolID)
	if err != nil {
		return nil, err
	}
	binding := &UpstreamPoolBinding{
		GroupID:          group.ID,
		GroupName:        group.Name,
		GroupPlatform:    group.Platform,
		PoolID:           pool.ID,
		Platform:         pool.Platform,
		Models:           append([]string{}, input.Models...),
		RequestPathScope: append([]string{}, input.RequestPathScope...),
		Priority:         input.Priority,
		Enabled:          input.Enabled,
	}
	if err := normalizeUpstreamPoolBindingForCreate(binding); err != nil {
		return nil, err
	}
	if err := validateUpstreamPoolBindingPlatform(group, pool, binding.Platform); err != nil {
		return nil, err
	}
	created, err := s.upstreamPoolRepo.CreateUpstreamPoolBinding(ctx, binding)
	if err != nil {
		return nil, err
	}
	created.GroupName = group.Name
	created.GroupPlatform = group.Platform
	return created, nil
}

func (s *adminServiceImpl) UpdateUpstreamPoolBinding(ctx context.Context, id int64, input *UpdateUpstreamPoolBindingInput) (*UpstreamPoolBinding, error) {
	if s == nil || s.upstreamPoolRepo == nil || s.groupRepo == nil || input == nil {
		return nil, ErrUpstreamPoolNotFound
	}
	binding, err := s.upstreamPoolRepo.GetUpstreamPoolBindingByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if input.GroupID != nil {
		binding.GroupID = *input.GroupID
	}
	if input.PoolID != nil {
		binding.PoolID = *input.PoolID
	}
	if input.Platform != nil {
		binding.Platform = strings.TrimSpace(*input.Platform)
	}
	if input.Models != nil {
		binding.Models = append([]string{}, input.Models...)
	}
	if input.RequestPathScope != nil {
		binding.RequestPathScope = append([]string{}, input.RequestPathScope...)
	}
	if input.Priority != nil {
		binding.Priority = *input.Priority
	}
	if input.Enabled != nil {
		binding.Enabled = *input.Enabled
	}
	group, err := s.groupRepo.GetByIDLite(ctx, binding.GroupID)
	if err != nil {
		return nil, err
	}
	pool, err := s.upstreamPoolRepo.GetUpstreamPoolByID(ctx, binding.PoolID)
	if err != nil {
		return nil, err
	}
	binding.Platform = pool.Platform
	if err := normalizeUpstreamPoolBindingForCreate(binding); err != nil {
		return nil, err
	}
	if err := validateUpstreamPoolBindingPlatform(group, pool, binding.Platform); err != nil {
		return nil, err
	}
	updated, err := s.upstreamPoolRepo.UpdateUpstreamPoolBinding(ctx, binding)
	if err != nil {
		return nil, err
	}
	updated.GroupName = group.Name
	updated.GroupPlatform = group.Platform
	return updated, nil
}

func (s *adminServiceImpl) DeleteUpstreamPoolBinding(ctx context.Context, id int64) error {
	if s == nil || s.upstreamPoolRepo == nil {
		return ErrUpstreamPoolNotFound
	}
	return s.upstreamPoolRepo.DeleteUpstreamPoolBinding(ctx, id)
}

func (s *adminServiceImpl) ensureUpstreamPoolPlatformMutable(ctx context.Context, poolID int64) error {
	reader, ok := s.upstreamPoolRepo.(UpstreamPoolPlatformUsageReader)
	if !ok {
		return nil
	}
	usage, err := reader.GetUpstreamPoolPlatformUsage(ctx, poolID)
	if err != nil {
		return err
	}
	if usage.IsLocked() {
		return NewUpstreamPoolPlatformLockedError("upstream_pool", usage)
	}
	return nil
}

func (s *adminServiceImpl) ensureUpstreamAccountSetPlatformMutable(ctx context.Context, setID int64) error {
	reader, ok := s.upstreamPoolRepo.(UpstreamPoolPlatformUsageReader)
	if !ok {
		return nil
	}
	usage, err := reader.GetUpstreamAccountSetPlatformUsage(ctx, setID)
	if err != nil {
		return err
	}
	if usage.IsLocked() {
		return NewUpstreamPoolPlatformLockedError("upstream_account_set", usage)
	}
	return nil
}

func validateUpstreamPoolBindingPlatform(group *Group, pool *UpstreamPool, bindingPlatform string) error {
	if group == nil || pool == nil {
		return ErrUpstreamPoolNotFound
	}
	if !strings.EqualFold(strings.TrimSpace(pool.Platform), strings.TrimSpace(bindingPlatform)) {
		return NewUpstreamPoolBadRequest("UPSTREAM_POOL_BINDING_PLATFORM_MISMATCH", "binding platform must match pool platform")
	}
	if !strings.EqualFold(strings.TrimSpace(group.Platform), strings.TrimSpace(pool.Platform)) {
		return NewUpstreamPoolBadRequest("UPSTREAM_POOL_BINDING_GROUP_PLATFORM_MISMATCH", "group platform must match pool platform")
	}
	return nil
}

func normalizeUpstreamPoolForCreate(pool *UpstreamPool) error {
	if pool == nil {
		return ErrUpstreamPoolNotFound
	}
	pool.Name = strings.TrimSpace(pool.Name)
	pool.Code = strings.TrimSpace(pool.Code)
	pool.Platform = strings.TrimSpace(pool.Platform)
	pool.Description = strings.TrimSpace(pool.Description)
	pool.SchedulerMode = strings.TrimSpace(pool.SchedulerMode)
	pool.DefaultRequiredCapability = strings.TrimSpace(pool.DefaultRequiredCapability)
	pool.DefaultRequiredTransport = strings.TrimSpace(pool.DefaultRequiredTransport)
	if pool.Name == "" {
		return errors.New("name is required")
	}
	if pool.Code == "" {
		return errors.New("code is required")
	}
	if pool.Platform == "" {
		return errors.New("platform is required")
	}
	if pool.AutoWeightEnabled && !strings.EqualFold(pool.Platform, PlatformOpenAI) {
		return NewUpstreamPoolBadRequest("INVALID_UPSTREAM_POOL_AUTO_WEIGHT_PLATFORM", "auto weight currently supports OpenAI pools only")
	}
	if pool.SchedulerMode == "" {
		pool.SchedulerMode = UpstreamPoolSchedulerModeAdvanced
	}
	if pool.StickyTTLSeconds <= 0 {
		return NewUpstreamPoolBadRequest("INVALID_UPSTREAM_POOL_STICKY_TTL", "sticky_ttl_seconds must be > 0")
	}
	if pool.StickyEscapeErrorRateThreshold <= 0 || pool.StickyEscapeErrorRateThreshold > 1 {
		return NewUpstreamPoolBadRequest("INVALID_UPSTREAM_POOL_STICKY_ESCAPE_ERROR_RATE", "sticky_escape_error_rate_threshold must be > 0 and <= 1")
	}
	if pool.StickyEscapeTTFTMSThreshold <= 0 {
		return NewUpstreamPoolBadRequest("INVALID_UPSTREAM_POOL_STICKY_ESCAPE_TTFT", "sticky_escape_ttft_ms_threshold must be > 0")
	}
	if pool.TopK <= 0 {
		return NewUpstreamPoolBadRequest("INVALID_UPSTREAM_POOL_TOP_K", "top_k must be > 0")
	}
	if pool.MaxFailoverHops < 0 {
		return NewUpstreamPoolBadRequest("INVALID_UPSTREAM_POOL_MAX_FAILOVER_HOPS", "max_failover_hops must be >= 0")
	}
	if pool.WaitTimeoutMS < 0 {
		return NewUpstreamPoolBadRequest("INVALID_UPSTREAM_POOL_WAIT_TIMEOUT", "wait_timeout_ms must be >= 0")
	}
	if pool.MaxWaiting < 0 {
		return NewUpstreamPoolBadRequest("INVALID_UPSTREAM_POOL_MAX_WAITING", "max_waiting must be >= 0")
	}
	if pool.PolicyJSON == nil {
		pool.PolicyJSON = map[string]any{}
	}
	pool.AccountTypeStrategy = NormalizeUpstreamPoolAccountTypeStrategy(pool.AccountTypeStrategy)
	if pool.AccountTypeStrategy == UpstreamPoolAccountTypeStrategyAll {
		pool.AccountTypeStrategy = UpstreamPoolAccountTypeStrategyFromPolicyJSON(pool.PolicyJSON)
	}
	pool.PolicyJSON = SetUpstreamPoolAccountTypeStrategyPolicyJSON(pool.PolicyJSON, pool.AccountTypeStrategy)
	if pool.AutoWeightMode == "" {
		if pool.AutoWeightEnabled {
			pool.AutoWeightMode = "active"
		} else {
			pool.AutoWeightMode = UpstreamPoolAutoWeightModeFromPolicyJSON(pool.PolicyJSON)
		}
	}
	pool.PolicyJSON = SetUpstreamPoolAutoWeightModePolicyJSON(pool.PolicyJSON, pool.AutoWeightMode)
	pool.AutoWeightEnabled = pool.AutoWeightMode != "off"
	return nil
}

func normalizeUpstreamPoolMemberForCreate(member *UpstreamPoolMember) error {
	if member == nil {
		return ErrUpstreamPoolNotFound
	}
	if member.PoolID <= 0 {
		return errors.New("pool_id is required")
	}
	if member.AccountID <= 0 {
		return errors.New("account_id is required")
	}
	if member.Weight <= 0 {
		member.Weight = 100
	}
	if member.Notes == "" {
		member.Notes = ""
	}
	return nil
}

func normalizeUpstreamPoolBindingForCreate(binding *UpstreamPoolBinding) error {
	if binding == nil {
		return ErrUpstreamPoolNotFound
	}
	binding.Platform = strings.TrimSpace(binding.Platform)
	if binding.GroupID <= 0 {
		return errors.New("group_id is required")
	}
	if binding.PoolID <= 0 {
		return errors.New("pool_id is required")
	}
	if binding.Platform == "" {
		return errors.New("platform is required")
	}
	if binding.Priority < 0 {
		return NewUpstreamPoolBadRequest("INVALID_UPSTREAM_POOL_BINDING_PRIORITY", "priority must be >= 0")
	}
	if binding.Models == nil {
		binding.Models = []string{}
	}
	if binding.RequestPathScope == nil {
		binding.RequestPathScope = []string{}
	}
	return nil
}

func normalizeUpstreamAccountSetForCreate(item *UpstreamAccountSet) error {
	if item == nil {
		return ErrUpstreamPoolNotFound
	}
	item.Name = strings.TrimSpace(item.Name)
	item.Code = strings.TrimSpace(item.Code)
	item.Platform = strings.TrimSpace(item.Platform)
	item.Description = strings.TrimSpace(item.Description)
	if item.Name == "" {
		return errors.New("name is required")
	}
	if item.Code == "" {
		return errors.New("code is required")
	}
	if item.Platform == "" {
		return errors.New("platform is required")
	}
	if item.SharedConcurrencyLimit != nil && *item.SharedConcurrencyLimit <= 0 {
		return errors.New("shared_concurrency_limit must be greater than 0")
	}
	return nil
}

func clonePositiveIntPointer(value *int) *int {
	if value == nil {
		return nil
	}
	cloned := *value
	return &cloned
}

func buildAutoUpstreamAccountSetCode(platform, name string) string {
	normalizedPlatform := strings.ToLower(strings.TrimSpace(platform))
	if normalizedPlatform == "" {
		normalizedPlatform = "set"
	}

	normalizedName := strings.ToLower(strings.TrimSpace(name))
	var builder strings.Builder
	lastDash := false
	for _, r := range normalizedName {
		switch {
		case r >= 'a' && r <= 'z':
			builder.WriteRune(r)
			lastDash = false
		case r >= '0' && r <= '9':
			builder.WriteRune(r)
			lastDash = false
		case r == '-' || r == '_' || r == '.' || r == ' ':
			if builder.Len() > 0 && !lastDash {
				builder.WriteByte('-')
				lastDash = true
			}
		}
		if builder.Len() >= 24 {
			break
		}
	}

	baseName := strings.Trim(builder.String(), "-")
	if baseName == "" {
		baseName = "set"
	}

	h := fnv.New32a()
	_, _ = h.Write([]byte(normalizedPlatform + ":" + strings.TrimSpace(name)))
	return fmt.Sprintf("%s-%s-%08x", normalizedPlatform, baseName, h.Sum32())
}

func normalizeUpstreamPoolMemberSetForCreate(item *UpstreamPoolMemberSet) error {
	if item == nil {
		return ErrUpstreamPoolNotFound
	}
	if item.PoolID <= 0 {
		return errors.New("pool_id is required")
	}
	if item.SetID <= 0 {
		return errors.New("set_id is required")
	}
	item.Notes = strings.TrimSpace(item.Notes)
	return nil
}

func uniquePositiveInt64s(values []int64) []int64 {
	if len(values) == 0 {
		return nil
	}
	seen := make(map[int64]struct{}, len(values))
	out := make([]int64, 0, len(values))
	for _, value := range values {
		if value <= 0 {
			continue
		}
		if _, ok := seen[value]; ok {
			continue
		}
		seen[value] = struct{}{}
		out = append(out, value)
	}
	return out
}

// User management implementations
func (s *adminServiceImpl) syncUpstreamPoolMembersForAccount(ctx context.Context, previousAccount, currentAccount *Account) error {
	if s == nil || s.upstreamPoolRepo == nil || s.accountRepo == nil {
		return nil
	}

	previousSupported := isAccountSupportedByUpstreamPoolAutoSync(previousAccount)
	currentSupported := isAccountSupportedByUpstreamPoolAutoSync(currentAccount)

	if !previousSupported && !currentSupported {
		return nil
	}

	var accountID int64
	var currentGroupIDs []int64
	platforms := make(map[string]struct{})
	if currentSupported {
		accountID = currentAccount.ID
		platforms[currentAccount.Platform] = struct{}{}
		if len(currentAccount.GroupIDs) > 0 {
			currentGroupIDs = append(currentGroupIDs, currentAccount.GroupIDs...)
		} else {
			accountGroups, err := s.accountRepo.GetGroups(ctx, currentAccount.ID)
			if err != nil {
				return fmt.Errorf("get account groups: %w", err)
			}
			currentGroupIDs = make([]int64, 0, len(accountGroups))
			for _, group := range accountGroups {
				if group.ID > 0 {
					currentGroupIDs = append(currentGroupIDs, group.ID)
				}
			}
		}
	}
	if previousSupported {
		platforms[previousAccount.Platform] = struct{}{}
	}
	if accountID == 0 && previousSupported {
		accountID = previousAccount.ID
	}

	bindings, err := s.upstreamPoolRepo.ListUpstreamPoolBindings(ctx)
	if err != nil {
		return fmt.Errorf("list upstream pool bindings: %w", err)
	}

	targetPoolIDs := make(map[int64]struct{})
	platformPoolIDs := make(map[int64]struct{})
	for _, binding := range bindings {
		if !binding.Enabled {
			continue
		}
		if _, ok := platforms[binding.Platform]; !ok {
			continue
		}
		platformPoolIDs[binding.PoolID] = struct{}{}
		if currentSupported && binding.Platform == currentAccount.Platform {
			for _, groupID := range currentGroupIDs {
				if binding.GroupID == groupID {
					targetPoolIDs[binding.PoolID] = struct{}{}
				}
			}
		}
	}

	for poolID := range platformPoolIDs {
		members, err := s.upstreamPoolRepo.ListUpstreamPoolMembers(ctx, poolID)
		if err != nil {
			return fmt.Errorf("list upstream pool members: %w", err)
		}

		var existing *UpstreamPoolMember
		for i := range members {
			if members[i].AccountID == accountID {
				existing = &members[i]
				break
			}
		}

		if !currentSupported {
			if existing != nil {
				if err := s.upstreamPoolRepo.DeleteUpstreamPoolMember(ctx, existing.ID); err != nil {
					return fmt.Errorf("delete upstream pool member: %w", err)
				}
			}
			continue
		}

		if _, ok := targetPoolIDs[poolID]; !ok {
			if existing != nil {
				if err := s.upstreamPoolRepo.DeleteUpstreamPoolMember(ctx, existing.ID); err != nil {
					return fmt.Errorf("delete upstream pool member: %w", err)
				}
			}
			continue
		}

		wantEnabled := currentAccount.IsSchedulable() && currentAccount.IsActive()
		if existing == nil {
			_, err = s.upstreamPoolRepo.CreateUpstreamPoolMember(ctx, &UpstreamPoolMember{
				PoolID:                 poolID,
				AccountID:              accountID,
				AccountName:            currentAccount.Name,
				AccountPlatform:        currentAccount.Platform,
				Enabled:                wantEnabled,
				SchedulableOverride:    nil,
				ManualDrained:          !wantEnabled,
				Weight:                 100,
				PriorityOverride:       nil,
				MaxConcurrencyOverride: nil,
				Notes:                  "synced from account",
			})
			if err != nil {
				return fmt.Errorf("create upstream pool member: %w", err)
			}
			continue
		}

		existing.AccountName = currentAccount.Name
		existing.AccountPlatform = currentAccount.Platform
		existing.Enabled = wantEnabled
		existing.ManualDrained = !wantEnabled
		existing.UpdatedAt = time.Now()
		if _, err = s.upstreamPoolRepo.UpdateUpstreamPoolMember(ctx, existing); err != nil {
			return fmt.Errorf("update upstream pool member: %w", err)
		}
	}
	return nil
}

func isAccountSupportedByUpstreamPoolAutoSync(account *Account) bool {
	if account == nil {
		return false
	}
	switch account.Platform {
	case PlatformOpenAI:
		return account.Type == AccountTypeAPIKey
	case PlatformAnthropic:
		switch account.Type {
		case AccountTypeOAuth, AccountTypeSetupToken, AccountTypeAPIKey, AccountTypeServiceAccount, AccountTypeBedrock:
			return true
		default:
			return false
		}
	default:
		return false
	}
}

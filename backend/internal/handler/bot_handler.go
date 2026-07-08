package handler

import (
	"strconv"
	"strings"
	"time"

	dbent "github.com/Wei-Shaw/sub2api/ent"
	"github.com/Wei-Shaw/sub2api/internal/handler/quotaview"
	infraerrors "github.com/Wei-Shaw/sub2api/internal/pkg/errors"
	"github.com/Wei-Shaw/sub2api/internal/pkg/response"
	"github.com/Wei-Shaw/sub2api/internal/service"

	"github.com/gin-gonic/gin"
)

type BotHandler struct {
	adminService          service.AdminService
	paymentService        *service.PaymentService
	poolHealthService     *service.PoolHealthService
	userPlatformQuotaRepo service.UserPlatformQuotaRepository
}

func NewBotHandler(
	adminService service.AdminService,
	paymentService *service.PaymentService,
	poolHealthService *service.PoolHealthService,
	userPlatformQuotaRepo service.UserPlatformQuotaRepository,
) *BotHandler {
	return &BotHandler{
		adminService:          adminService,
		paymentService:        paymentService,
		poolHealthService:     poolHealthService,
		userPlatformQuotaRepo: userPlatformQuotaRepo,
	}
}

type botUserResponse struct {
	ID              int64      `json:"id"`
	Email           string     `json:"email"`
	Username        string     `json:"username"`
	Role            string     `json:"role"`
	Status          string     `json:"status"`
	Balance         float64    `json:"balance"`
	Concurrency     int        `json:"concurrency"`
	RPMLimit        int        `json:"rpm_limit"`
	AllowedGroups   []int64    `json:"allowed_groups"`
	LastLoginAt     *time.Time `json:"last_login_at,omitempty"`
	LastActiveAt    *time.Time `json:"last_active_at,omitempty"`
	LastUsedAt      *time.Time `json:"last_used_at,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	TotalRecharged  float64    `json:"total_recharged"`
	BalanceNotifyOn bool       `json:"balance_notify_enabled"`
}

type botAPIKeyResponse struct {
	ID              int64      `json:"id"`
	UserID          int64      `json:"user_id"`
	Name            string     `json:"name"`
	GroupID         *int64     `json:"group_id,omitempty"`
	Status          string     `json:"status"`
	MaskedKey       string     `json:"masked_key"`
	KeyTail         string     `json:"key_tail"`
	LastUsedAt      *time.Time `json:"last_used_at,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	Quota           float64    `json:"quota"`
	QuotaUsed       float64    `json:"quota_used"`
	QuotaRemaining  float64    `json:"quota_remaining"`
	ExpiresAt       *time.Time `json:"expires_at,omitempty"`
	RateLimit5h     float64    `json:"rate_limit_5h"`
	RateLimit1d     float64    `json:"rate_limit_1d"`
	RateLimit7d     float64    `json:"rate_limit_7d"`
	Usage5h         float64    `json:"usage_5h"`
	Usage1d         float64    `json:"usage_1d"`
	Usage7d         float64    `json:"usage_7d"`
	Window5hStart   *time.Time `json:"window_5h_start,omitempty"`
	Window1dStart   *time.Time `json:"window_1d_start,omitempty"`
	Window7dStart   *time.Time `json:"window_7d_start,omitempty"`
	DaysUntilExpiry int        `json:"days_until_expiry"`
}

type botLineHealthItem struct {
	ID                  int64                     `json:"id"`
	Name                string                    `json:"name"`
	Provider            string                    `json:"provider"`
	GroupID             *int64                    `json:"group_id,omitempty"`
	GroupName           string                    `json:"group_name"`
	Status              string                    `json:"status"`
	Availability7d      float64                   `json:"availability_7d"`
	BestLatencyMs       *int                      `json:"best_latency_ms"`
	BestPingLatencyMs   *int                      `json:"best_ping_latency_ms"`
	HealthyMemberCount  int                       `json:"healthy_member_count"`
	DegradedMemberCount int                       `json:"degraded_member_count"`
	FailedMemberCount   int                       `json:"failed_member_count"`
	TotalMemberCount    int                       `json:"total_member_count"`
	Timeline            []poolHealthTimelinePoint `json:"timeline"`
}

type botOrderResponse struct {
	ID                  int64      `json:"id"`
	UserID              int64      `json:"user_id"`
	UserEmail           string     `json:"user_email"`
	UserName            string     `json:"user_name"`
	Amount              float64    `json:"amount"`
	PayAmount           float64    `json:"pay_amount"`
	PaymentType         string     `json:"payment_type"`
	OrderType           string     `json:"order_type"`
	PlanID              *int64     `json:"plan_id,omitempty"`
	SubscriptionGroupID *int64     `json:"subscription_group_id,omitempty"`
	SubscriptionDays    *int       `json:"subscription_days,omitempty"`
	Status              string     `json:"status"`
	RefundAmount        float64    `json:"refund_amount"`
	ExpiresAt           time.Time  `json:"expires_at"`
	PaidAt              *time.Time `json:"paid_at,omitempty"`
	CompletedAt         *time.Time `json:"completed_at,omitempty"`
	FailedAt            *time.Time `json:"failed_at,omitempty"`
	FailedReason        *string    `json:"failed_reason,omitempty"`
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at"`
}

type botResolveUserRequest struct {
	Login string `json:"login"`
}

type botResolvedUserResponse struct {
	UserID   int64  `json:"user_id"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Status   string `json:"status"`
}

func (h *BotHandler) ListLineHealth(c *gin.Context) {
	views, err := h.poolHealthService.ListUserPoolHealth(c.Request.Context())
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	items := make([]botLineHealthItem, 0, len(views))
	for _, view := range views {
		items = append(items, botLineHealthView(view))
	}
	response.Success(c, gin.H{"items": items})
}

func (h *BotHandler) GetUser(c *gin.Context) {
	userID, ok := parseBotUserID(c)
	if !ok {
		return
	}
	user, err := h.adminService.GetUser(c.Request.Context(), userID)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, botUserFromService(user))
}

func (h *BotHandler) GetUserUsage(c *gin.Context) {
	userID, ok := parseBotUserID(c)
	if !ok {
		return
	}
	period := c.DefaultQuery("period", "today")
	stats, err := h.adminService.GetUserUsageStats(c.Request.Context(), userID, period)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, stats)
}

func (h *BotHandler) GetUserPlatformQuotas(c *gin.Context) {
	userID, ok := parseBotUserID(c)
	if !ok {
		return
	}
	if h.userPlatformQuotaRepo == nil {
		response.Success(c, gin.H{"platform_quotas": []any{}})
		return
	}
	if _, err := h.adminService.GetUser(c.Request.Context(), userID); err != nil {
		response.ErrorFrom(c, err)
		return
	}
	records, err := h.userPlatformQuotaRepo.ListByUser(c.Request.Context(), userID)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	now := time.Now().UTC()
	out := make([]map[string]any, 0, len(records))
	for _, record := range records {
		out = append(out, quotaview.LazyZeroQuotaForResponse(record, now, true))
	}
	response.Success(c, gin.H{"platform_quotas": out})
}

func (h *BotHandler) GetUserAPIKeys(c *gin.Context) {
	userID, ok := parseBotUserID(c)
	if !ok {
		return
	}
	page, pageSize := response.ParsePagination(c)
	sortBy := c.DefaultQuery("sort_by", "created_at")
	sortOrder := c.DefaultQuery("sort_order", "desc")

	keys, total, err := h.adminService.GetUserAPIKeys(c.Request.Context(), userID, page, pageSize, sortBy, sortOrder)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	out := make([]botAPIKeyResponse, 0, len(keys))
	for i := range keys {
		out = append(out, botAPIKeyFromService(&keys[i]))
	}
	response.Paginated(c, out, total, page, pageSize)
}

func (h *BotHandler) GetUserOrders(c *gin.Context) {
	userID, ok := parseBotUserID(c)
	if !ok {
		return
	}
	page, pageSize := response.ParsePagination(c)
	orders, total, err := h.paymentService.AdminListOrders(c.Request.Context(), userID, service.OrderListParams{
		Page:        page,
		PageSize:    pageSize,
		Status:      c.Query("status"),
		OrderType:   c.Query("order_type"),
		PaymentType: c.Query("payment_type"),
	})
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	out := make([]botOrderResponse, 0, len(orders))
	for _, order := range orders {
		out = append(out, botOrderFromEnt(order))
	}
	response.Paginated(c, out, int64(total), page, pageSize)
}

func (h *BotHandler) ResolveUser(c *gin.Context) {
	var req botResolveUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
	}
	login := strings.TrimSpace(req.Login)
	if login == "" {
		response.BadRequest(c, "login is required")
		return
	}
	users, _, err := h.adminService.ListUsers(c.Request.Context(), 1, 10, service.UserListFilters{Search: login}, "id", "asc")
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	var matched *service.User
	for i := range users {
		if strings.EqualFold(strings.TrimSpace(users[i].Email), login) || strings.EqualFold(strings.TrimSpace(users[i].Username), login) {
			matched = &users[i]
			break
		}
	}
	if matched == nil {
		response.ErrorFrom(c, service.ErrUserNotFound)
		return
	}
	response.Success(c, botResolvedUserResponse{
		UserID:   matched.ID,
		Email:    matched.Email,
		Username: matched.Username,
		Status:   matched.Status,
	})
}

func parseBotUserID(c *gin.Context) (int64, bool) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		response.ErrorFrom(c, infraerrors.BadRequest("INVALID_USER_ID", "invalid user id"))
		return 0, false
	}
	return id, true
}

func botUserFromService(user *service.User) botUserResponse {
	return botUserResponse{
		ID:              user.ID,
		Email:           user.Email,
		Username:        user.Username,
		Role:            user.Role,
		Status:          user.Status,
		Balance:         user.Balance,
		Concurrency:     user.Concurrency,
		RPMLimit:        user.RPMLimit,
		AllowedGroups:   user.AllowedGroups,
		LastLoginAt:     user.LastLoginAt,
		LastActiveAt:    user.LastActiveAt,
		LastUsedAt:      user.LastUsedAt,
		CreatedAt:       user.CreatedAt,
		UpdatedAt:       user.UpdatedAt,
		TotalRecharged:  user.TotalRecharged,
		BalanceNotifyOn: user.BalanceNotifyEnabled,
	}
}

func botAPIKeyFromService(key *service.APIKey) botAPIKeyResponse {
	return botAPIKeyResponse{
		ID:              key.ID,
		UserID:          key.UserID,
		Name:            key.Name,
		GroupID:         key.GroupID,
		Status:          key.Status,
		MaskedKey:       maskAPIKey(key.Key),
		KeyTail:         keyTail(key.Key),
		LastUsedAt:      key.LastUsedAt,
		CreatedAt:       key.CreatedAt,
		UpdatedAt:       key.UpdatedAt,
		Quota:           key.Quota,
		QuotaUsed:       key.QuotaUsed,
		QuotaRemaining:  key.GetQuotaRemaining(),
		ExpiresAt:       key.ExpiresAt,
		RateLimit5h:     key.RateLimit5h,
		RateLimit1d:     key.RateLimit1d,
		RateLimit7d:     key.RateLimit7d,
		Usage5h:         key.EffectiveUsage5h(),
		Usage1d:         key.EffectiveUsage1d(),
		Usage7d:         key.EffectiveUsage7d(),
		Window5hStart:   key.Window5hStart,
		Window1dStart:   key.Window1dStart,
		Window7dStart:   key.Window7dStart,
		DaysUntilExpiry: key.GetDaysUntilExpiry(),
	}
}

func botLineHealthView(view *service.PoolHealthView) botLineHealthItem {
	return botLineHealthItem{
		ID:                  view.ID,
		Name:                view.Name,
		Provider:            view.Provider,
		GroupID:             view.GroupID,
		GroupName:           view.GroupName,
		Status:              view.Status,
		Availability7d:      view.Availability7d,
		BestLatencyMs:       view.BestLatencyMs,
		BestPingLatencyMs:   view.BestPingLatencyMs,
		HealthyMemberCount:  view.HealthyMemberCount,
		DegradedMemberCount: view.DegradedMemberCount,
		FailedMemberCount:   view.FailedMemberCount,
		TotalMemberCount:    view.TotalMemberCount,
		Timeline:            poolHealthTimelineToResponse(view.Timeline),
	}
}

func botOrderFromEnt(order *dbent.PaymentOrder) botOrderResponse {
	if order == nil {
		return botOrderResponse{}
	}
	return botOrderResponse{
		ID:                  order.ID,
		UserID:              order.UserID,
		UserEmail:           order.UserEmail,
		UserName:            order.UserName,
		Amount:              order.Amount,
		PayAmount:           order.PayAmount,
		PaymentType:         order.PaymentType,
		OrderType:           order.OrderType,
		PlanID:              order.PlanID,
		SubscriptionGroupID: order.SubscriptionGroupID,
		SubscriptionDays:    order.SubscriptionDays,
		Status:              order.Status,
		RefundAmount:        order.RefundAmount,
		ExpiresAt:           order.ExpiresAt,
		PaidAt:              order.PaidAt,
		CompletedAt:         order.CompletedAt,
		FailedAt:            order.FailedAt,
		FailedReason:        order.FailedReason,
		CreatedAt:           order.CreatedAt,
		UpdatedAt:           order.UpdatedAt,
	}
}

func maskAPIKey(key string) string {
	key = strings.TrimSpace(key)
	if key == "" {
		return ""
	}
	tail := keyTail(key)
	prefix := key
	if len(prefix) > 6 {
		prefix = prefix[:6]
	}
	if tail == "" || tail == prefix {
		return "****"
	}
	return prefix + "..." + tail
}

func keyTail(key string) string {
	key = strings.TrimSpace(key)
	if len(key) <= 4 {
		return ""
	}
	return key[len(key)-4:]
}

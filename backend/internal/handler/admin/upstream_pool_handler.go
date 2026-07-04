package admin

import (
	"context"
	"strconv"
	"strings"

	"github.com/Wei-Shaw/sub2api/internal/pkg/response"
	"github.com/Wei-Shaw/sub2api/internal/handler/dto"
	"github.com/Wei-Shaw/sub2api/internal/service"

	"github.com/gin-gonic/gin"
)

type UpstreamPoolHandler struct {
	adminService service.AdminService
}

func NewUpstreamPoolHandler(adminService service.AdminService) *UpstreamPoolHandler {
	return &UpstreamPoolHandler{adminService: adminService}
}

func (h *UpstreamPoolHandler) List(c *gin.Context) {
	pools, err := h.adminService.ListUpstreamPools(c.Request.Context())
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	out := make([]dto.UpstreamPool, 0, len(pools))
	for i := range pools {
		out = append(out, *dto.UpstreamPoolFromService(&pools[i]))
	}
	response.Success(c, out)
}

func (h *UpstreamPoolHandler) GetByID(c *gin.Context) {
	poolID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || poolID <= 0 {
		response.BadRequest(c, "Invalid pool ID")
		return
	}

	pool, err := h.adminService.GetUpstreamPoolByID(c.Request.Context(), poolID)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, dto.UpstreamPoolFromService(pool))
}

func (h *UpstreamPoolHandler) GetMembers(c *gin.Context) {
	poolID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || poolID <= 0 {
		response.BadRequest(c, "Invalid pool ID")
		return
	}

	members, err := h.adminService.ListUpstreamPoolMembers(c.Request.Context(), poolID)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	out := make([]dto.UpstreamPoolMember, 0, len(members))
	for i := range members {
		out = append(out, *dto.UpstreamPoolMemberFromService(&members[i]))
	}
	response.Success(c, out)
}

func (h *UpstreamPoolHandler) GetBindings(c *gin.Context) {
	bindings, err := h.adminService.ListUpstreamPoolBindings(c.Request.Context())
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	out := make([]dto.UpstreamPoolBinding, 0, len(bindings))
	for i := range bindings {
		out = append(out, *dto.UpstreamPoolBindingFromService(&bindings[i]))
	}
	response.Success(c, out)
}

func (h *UpstreamPoolHandler) Create(c *gin.Context) {
	var req upstreamPoolWriteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
	}
	executeAdminIdempotentJSON(c, "admin.upstream_pools.create", req, service.DefaultWriteIdempotencyTTL(), func(ctx context.Context) (any, error) {
		pool, err := h.adminService.CreateUpstreamPool(ctx, &service.CreateUpstreamPoolInput{
			Name:                             stringValue(req.Name),
			Code:                             stringValue(req.Code),
			Platform:                         stringValue(req.Platform),
			Description:                      stringValue(req.Description),
			Enabled:                          boolValue(req.Enabled, true),
			SchedulerMode:                    stringValue(req.SchedulerMode),
			DefaultRequiredCapability:        stringValue(req.DefaultRequiredCapability),
			DefaultRequiredTransport:         stringValue(req.DefaultRequiredTransport),
			StickyEnabled:                    boolValue(req.StickyEnabled, true),
			StickyTTLSeconds:                 intValue(req.StickyTTLSeconds, 1800),
			StickyEscapeEnabled:              boolValue(req.StickyEscapeEnabled, true),
			StickyEscapeErrorRateThreshold:    float64Value(req.StickyEscapeErrorRateThreshold, 0.3000),
			StickyEscapeTTFTMSThreshold:      intValue(req.StickyEscapeTTFTMSThreshold, 6000),
			LoadBalanceEnabled:               boolValue(req.LoadBalanceEnabled, true),
			FailoverEnabled:                  boolValue(req.FailoverEnabled, true),
			TopK:                             intValue(req.TopK, 2),
			MaxFailoverHops:                  intValue(req.MaxFailoverHops, 3),
			WaitTimeoutMS:                    intValue(req.WaitTimeoutMS, 30000),
			MaxWaiting:                       intValue(req.MaxWaiting, 100),
			PolicyJSON:                       req.PolicyJSON,
		})
		if err != nil {
			return nil, err
		}
		return dto.UpstreamPoolFromService(pool), nil
	})
}

func (h *UpstreamPoolHandler) Update(c *gin.Context) {
	poolID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || poolID <= 0 {
		response.BadRequest(c, "Invalid pool ID")
		return
	}

	var req upstreamPoolWriteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
	}
	executeAdminIdempotentJSON(c, "admin.upstream_pools.update", map[string]any{"id": poolID, "payload": req}, service.DefaultWriteIdempotencyTTL(), func(ctx context.Context) (any, error) {
		pool, err := h.adminService.UpdateUpstreamPool(ctx, poolID, &service.UpdateUpstreamPoolInput{
			Name:                              req.Name,
			Code:                              req.Code,
			Platform:                          req.Platform,
			Description:                       req.Description,
			Enabled:                           req.Enabled,
			SchedulerMode:                     req.SchedulerMode,
			DefaultRequiredCapability:         req.DefaultRequiredCapability,
			DefaultRequiredTransport:          req.DefaultRequiredTransport,
			StickyEnabled:                     req.StickyEnabled,
			StickyTTLSeconds:                  req.StickyTTLSeconds,
			StickyEscapeEnabled:               req.StickyEscapeEnabled,
			StickyEscapeErrorRateThreshold:     req.StickyEscapeErrorRateThreshold,
			StickyEscapeTTFTMSThreshold:       req.StickyEscapeTTFTMSThreshold,
			LoadBalanceEnabled:                req.LoadBalanceEnabled,
			FailoverEnabled:                   req.FailoverEnabled,
			TopK:                              req.TopK,
			MaxFailoverHops:                   req.MaxFailoverHops,
			WaitTimeoutMS:                     req.WaitTimeoutMS,
			MaxWaiting:                        req.MaxWaiting,
			PolicyJSON:                        req.PolicyJSON,
		})
		if err != nil {
			return nil, err
		}
		return dto.UpstreamPoolFromService(pool), nil
	})
}

func (h *UpstreamPoolHandler) Delete(c *gin.Context) {
	poolID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || poolID <= 0 {
		response.BadRequest(c, "Invalid pool ID")
		return
	}
	if err := h.adminService.DeleteUpstreamPool(c.Request.Context(), poolID); err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, gin.H{"message": "Upstream pool deleted successfully"})
}

func (h *UpstreamPoolHandler) CreateMember(c *gin.Context) {
	poolID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || poolID <= 0 {
		response.BadRequest(c, "Invalid pool ID")
		return
	}
	var req upstreamPoolMemberWriteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
	}
	executeAdminIdempotentJSON(c, "admin.upstream_pools.members.create", map[string]any{"pool_id": poolID, "payload": req}, service.DefaultWriteIdempotencyTTL(), func(ctx context.Context) (any, error) {
		member, err := h.adminService.CreateUpstreamPoolMember(ctx, poolID, &service.CreateUpstreamPoolMemberInput{
			AccountID:              req.AccountID,
			Enabled:                boolValue(req.Enabled, true),
			SchedulableOverride:    req.SchedulableOverride,
			ManualDrained:          boolValue(req.ManualDrained, false),
			Weight:                 intValue(req.Weight, 100),
			PriorityOverride:       req.PriorityOverride,
			MaxConcurrencyOverride: req.MaxConcurrencyOverride,
			Notes:                  stringValue(req.Notes),
		})
		if err != nil {
			return nil, err
		}
		return dto.UpstreamPoolMemberFromService(member), nil
	})
}

func (h *UpstreamPoolHandler) UpdateMember(c *gin.Context) {
	memberID, err := strconv.ParseInt(c.Param("member_id"), 10, 64)
	if err != nil || memberID <= 0 {
		response.BadRequest(c, "Invalid member ID")
		return
	}
	var req upstreamPoolMemberWriteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
	}
	executeAdminIdempotentJSON(c, "admin.upstream_pools.members.update", map[string]any{"member_id": memberID, "payload": req}, service.DefaultWriteIdempotencyTTL(), func(ctx context.Context) (any, error) {
		member, err := h.adminService.UpdateUpstreamPoolMember(ctx, memberID, &service.UpdateUpstreamPoolMemberInput{
			Enabled:                req.Enabled,
			SchedulableOverride:    req.SchedulableOverride,
			ManualDrained:          req.ManualDrained,
			Weight:                 req.Weight,
			PriorityOverride:       req.PriorityOverride,
			MaxConcurrencyOverride: req.MaxConcurrencyOverride,
			Notes:                  req.Notes,
		})
		if err != nil {
			return nil, err
		}
		return dto.UpstreamPoolMemberFromService(member), nil
	})
}

func (h *UpstreamPoolHandler) DeleteMember(c *gin.Context) {
	memberID, err := strconv.ParseInt(c.Param("member_id"), 10, 64)
	if err != nil || memberID <= 0 {
		response.BadRequest(c, "Invalid member ID")
		return
	}
	if err := h.adminService.DeleteUpstreamPoolMember(c.Request.Context(), memberID); err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, gin.H{"message": "Upstream pool member deleted successfully"})
}

func (h *UpstreamPoolHandler) CreateBinding(c *gin.Context) {
	var req upstreamPoolBindingWriteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
	}
	executeAdminIdempotentJSON(c, "admin.upstream_pools.bindings.create", req, service.DefaultWriteIdempotencyTTL(), func(ctx context.Context) (any, error) {
		binding, err := h.adminService.CreateUpstreamPoolBinding(ctx, &service.CreateUpstreamPoolBindingInput{
			GroupID:          int64Value(req.GroupID, 0),
			PoolID:           int64Value(req.PoolID, 0),
			Platform:         stringValue(req.Platform),
			Models:           append([]string{}, req.Models...),
			RequestPathScope: append([]string{}, req.RequestPathScope...),
			Priority:         intValue(req.Priority, 100),
			Enabled:          boolValue(req.Enabled, true),
		})
		if err != nil {
			return nil, err
		}
		return dto.UpstreamPoolBindingFromService(binding), nil
	})
}

func (h *UpstreamPoolHandler) UpdateBinding(c *gin.Context) {
	bindingID, err := strconv.ParseInt(c.Param("binding_id"), 10, 64)
	if err != nil || bindingID <= 0 {
		response.BadRequest(c, "Invalid binding ID")
		return
	}
	var req upstreamPoolBindingWriteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
	}
	executeAdminIdempotentJSON(c, "admin.upstream_pools.bindings.update", map[string]any{"binding_id": bindingID, "payload": req}, service.DefaultWriteIdempotencyTTL(), func(ctx context.Context) (any, error) {
		binding, err := h.adminService.UpdateUpstreamPoolBinding(ctx, bindingID, &service.UpdateUpstreamPoolBindingInput{
			GroupID:          req.GroupID,
			PoolID:           req.PoolID,
			Platform:         req.Platform,
			Models:           append([]string{}, req.Models...),
			RequestPathScope: append([]string{}, req.RequestPathScope...),
			Priority:         req.Priority,
			Enabled:          req.Enabled,
		})
		if err != nil {
			return nil, err
		}
		return dto.UpstreamPoolBindingFromService(binding), nil
	})
}

func (h *UpstreamPoolHandler) DeleteBinding(c *gin.Context) {
	bindingID, err := strconv.ParseInt(c.Param("binding_id"), 10, 64)
	if err != nil || bindingID <= 0 {
		response.BadRequest(c, "Invalid binding ID")
		return
	}
	if err := h.adminService.DeleteUpstreamPoolBinding(c.Request.Context(), bindingID); err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, gin.H{"message": "Upstream pool binding deleted successfully"})
}

type upstreamPoolWriteRequest struct {
	Name                             *string        `json:"name"`
	Code                             *string        `json:"code"`
	Platform                         *string        `json:"platform"`
	Description                      *string        `json:"description"`
	Enabled                          *bool          `json:"enabled"`
	SchedulerMode                    *string        `json:"scheduler_mode"`
	DefaultRequiredCapability        *string        `json:"default_required_capability"`
	DefaultRequiredTransport         *string        `json:"default_required_transport"`
	StickyEnabled                    *bool          `json:"sticky_enabled"`
	StickyTTLSeconds                 *int           `json:"sticky_ttl_seconds"`
	StickyEscapeEnabled              *bool          `json:"sticky_escape_enabled"`
	StickyEscapeErrorRateThreshold    *float64       `json:"sticky_escape_error_rate_threshold"`
	StickyEscapeTTFTMSThreshold      *int           `json:"sticky_escape_ttft_ms_threshold"`
	LoadBalanceEnabled               *bool          `json:"load_balance_enabled"`
	FailoverEnabled                  *bool          `json:"failover_enabled"`
	TopK                             *int           `json:"top_k"`
	MaxFailoverHops                  *int           `json:"max_failover_hops"`
	WaitTimeoutMS                    *int           `json:"wait_timeout_ms"`
	MaxWaiting                       *int           `json:"max_waiting"`
	PolicyJSON                       map[string]any `json:"policy_json"`
}

type upstreamPoolMemberWriteRequest struct {
	AccountID              int64   `json:"account_id"`
	Enabled                *bool   `json:"enabled"`
	SchedulableOverride    *bool   `json:"schedulable_override"`
	ManualDrained          *bool   `json:"manual_drained"`
	Weight                 *int    `json:"weight"`
	PriorityOverride       *int    `json:"priority_override"`
	MaxConcurrencyOverride *int    `json:"max_concurrency_override"`
	Notes                  *string `json:"notes"`
}

type upstreamPoolBindingWriteRequest struct {
	GroupID          *int64   `json:"group_id"`
	PoolID           *int64   `json:"pool_id"`
	Platform         *string  `json:"platform"`
	Models           []string `json:"models"`
	RequestPathScope []string `json:"request_path_scope"`
	Priority         *int     `json:"priority"`
	Enabled          *bool    `json:"enabled"`
}

func stringValue(v *string) string {
	if v == nil {
		return ""
	}
	return strings.TrimSpace(*v)
}

func boolValue(v *bool, fallback bool) bool {
	if v == nil {
		return fallback
	}
	return *v
}

func intValue(v *int, fallback int) int {
	if v == nil {
		return fallback
	}
	return *v
}

func int64Value(v *int64, fallback int64) int64 {
	if v == nil {
		return fallback
	}
	return *v
}

func float64Value(v *float64, fallback float64) float64 {
	if v == nil {
		return fallback
	}
	return *v
}

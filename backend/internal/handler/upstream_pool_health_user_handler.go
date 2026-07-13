package handler

import (
	"strconv"
	"time"

	infraerrors "github.com/Wei-Shaw/sub2api/internal/pkg/errors"
	"github.com/Wei-Shaw/sub2api/internal/pkg/response"
	"github.com/Wei-Shaw/sub2api/internal/service"

	"github.com/gin-gonic/gin"
)

type UpstreamPoolHealthUserHandler struct {
	poolHealthService *service.PoolHealthService
	settingService    *service.SettingService
}

func NewUpstreamPoolHealthUserHandler(
	poolHealthService *service.PoolHealthService,
	settingService *service.SettingService,
) *UpstreamPoolHealthUserHandler {
	return &UpstreamPoolHealthUserHandler{
		poolHealthService: poolHealthService,
		settingService:    settingService,
	}
}

func (h *UpstreamPoolHealthUserHandler) featureEnabled(c *gin.Context) bool {
	if h.settingService == nil {
		return true
	}
	return h.settingService.GetChannelMonitorRuntime(c.Request.Context()).Enabled
}

type poolHealthListItem struct {
	ID                int64                     `json:"id"`
	Name              string                    `json:"name"`
	Status            string                    `json:"status"`
	CapacityStatus    string                    `json:"capacity_status"`
	Availability7d    float64                   `json:"availability_7d"`
	BestLatencyMs     *int                      `json:"best_latency_ms"`
	BestPingLatencyMs *int                      `json:"best_ping_latency_ms"`
	Timeline          []poolHealthTimelinePoint `json:"timeline"`
}

type poolHealthTimelinePoint struct {
	Status        string `json:"status"`
	LatencyMs     *int   `json:"latency_ms"`
	PingLatencyMs *int   `json:"ping_latency_ms"`
	CheckedAt     string `json:"checked_at"`
}

type poolHealthDetailResponse struct {
	ID                int64                     `json:"id"`
	Name              string                    `json:"name"`
	Status            string                    `json:"status"`
	CapacityStatus    string                    `json:"capacity_status"`
	Availability7d    float64                   `json:"availability_7d"`
	Availability15d   float64                   `json:"availability_15d"`
	Availability30d   float64                   `json:"availability_30d"`
	BestLatencyMs     *int                      `json:"best_latency_ms"`
	BestPingLatencyMs *int                      `json:"best_ping_latency_ms"`
	Timeline          []poolHealthTimelinePoint `json:"timeline"`
}

func poolHealthViewToItem(v *service.PoolHealthView) poolHealthListItem {
	return poolHealthListItem{
		ID:                v.ID,
		Name:              v.Name,
		Status:            v.Status,
		CapacityStatus:    valueOrCapacityStatus(v.CapacityStatus, v.Status),
		Availability7d:    v.Availability7d,
		BestLatencyMs:     v.BestLatencyMs,
		BestPingLatencyMs: v.BestPingLatencyMs,
		Timeline:          poolHealthTimelineToResponse(v.Timeline),
	}
}

func poolHealthDetailToResponse(d *service.PoolHealthDetail) *poolHealthDetailResponse {
	return &poolHealthDetailResponse{
		ID:                d.ID,
		Name:              d.Name,
		Status:            d.Status,
		CapacityStatus:    valueOrCapacityStatus(d.CapacityStatus, d.Status),
		Availability7d:    d.Availability7d,
		Availability15d:   d.Availability15d,
		Availability30d:   d.Availability30d,
		BestLatencyMs:     d.BestLatencyMs,
		BestPingLatencyMs: d.BestPingLatencyMs,
		Timeline:          poolHealthTimelineToResponse(d.Timeline),
	}
}

func valueOrCapacityStatus(value, status string) string {
	if value != "" {
		return value
	}
	switch status {
	case service.MonitorStatusOperational:
		return "ample"
	case service.MonitorStatusDegraded:
		return "observe"
	case service.MonitorStatusFailed:
		return "tight"
	default:
		return "queueing"
	}
}

func poolHealthTimelineToResponse(points []service.PoolHealthTimelinePoint) []poolHealthTimelinePoint {
	out := make([]poolHealthTimelinePoint, 0, len(points))
	for _, point := range points {
		out = append(out, poolHealthTimelinePoint{
			Status:        point.Status,
			LatencyMs:     point.LatencyMs,
			PingLatencyMs: point.PingLatencyMs,
			CheckedAt:     point.CheckedAt.UTC().Format(time.RFC3339),
		})
	}
	return out
}

func (h *UpstreamPoolHealthUserHandler) List(c *gin.Context) {
	if !h.featureEnabled(c) {
		response.Success(c, gin.H{"items": []poolHealthListItem{}})
		return
	}
	views, err := h.poolHealthService.ListUserPoolHealth(c.Request.Context())
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	items := make([]poolHealthListItem, 0, len(views))
	for _, view := range views {
		items = append(items, poolHealthViewToItem(view))
	}
	response.Success(c, gin.H{"items": items})
}

func (h *UpstreamPoolHealthUserHandler) Get(c *gin.Context) {
	if !h.featureEnabled(c) {
		response.ErrorFrom(c, service.ErrPoolHealthNotFound)
		return
	}
	id, ok := parseUserPoolHealthID(c)
	if !ok {
		return
	}
	detail, err := h.poolHealthService.GetUserPoolHealthDetail(c.Request.Context(), id)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, poolHealthDetailToResponse(detail))
}

func parseUserPoolHealthID(c *gin.Context) (int64, bool) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		response.ErrorFrom(c, infraerrors.BadRequest("INVALID_POOL_ID", "invalid upstream pool id"))
		return 0, false
	}
	return id, true
}

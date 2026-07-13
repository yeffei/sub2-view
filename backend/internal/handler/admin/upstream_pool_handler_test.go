package admin

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/Wei-Shaw/sub2api/internal/service"
)

func TestUpstreamPoolMemberUpdateRequestNullableFields(t *testing.T) {
	var req upstreamPoolMemberUpdateRequest
	if err := json.Unmarshal([]byte(`{
		"schedulable_override": null,
		"priority_override": null,
		"max_concurrency_override": null,
		"notes": null
	}`), &req); err != nil {
		t.Fatalf("unmarshal nullable member update request: %v", err)
	}

	if !req.SchedulableOverride.Set || req.SchedulableOverride.Value != nil {
		t.Fatalf("schedulable_override = %+v, want set nil", req.SchedulableOverride)
	}
	if !req.PriorityOverride.Set || req.PriorityOverride.Value != nil {
		t.Fatalf("priority_override = %+v, want set nil", req.PriorityOverride)
	}
	if !req.MaxConcurrencyOverride.Set || req.MaxConcurrencyOverride.Value != nil {
		t.Fatalf("max_concurrency_override = %+v, want set nil", req.MaxConcurrencyOverride)
	}
	if !req.Notes.Set || req.Notes.Value != nil {
		t.Fatalf("notes = %+v, want set nil", req.Notes)
	}

	var empty upstreamPoolMemberUpdateRequest
	if err := json.Unmarshal([]byte(`{}`), &empty); err != nil {
		t.Fatalf("unmarshal empty member update request: %v", err)
	}
	if empty.SchedulableOverride.Set || empty.PriorityOverride.Set || empty.MaxConcurrencyOverride.Set || empty.Notes.Set {
		t.Fatalf("empty request should not mark nullable fields as set: %+v", empty)
	}
}

func TestUpstreamAccountSetWriteRequestSharedLimitNullable(t *testing.T) {
	var req upstreamAccountSetWriteRequest
	if err := json.Unmarshal([]byte(`{"shared_concurrency_limit":null}`), &req); err != nil {
		t.Fatalf("unmarshal shared limit: %v", err)
	}
	if !req.SharedConcurrencyLimit.Set || req.SharedConcurrencyLimit.Value != nil {
		t.Fatalf("shared limit = %+v, want set nil", req.SharedConcurrencyLimit)
	}
	var empty upstreamAccountSetWriteRequest
	if err := json.Unmarshal([]byte(`{}`), &empty); err != nil {
		t.Fatalf("unmarshal empty account set request: %v", err)
	}
	if empty.SharedConcurrencyLimit.Set {
		t.Fatal("empty account set request should not mark shared limit as set")
	}
}

func TestUpstreamPoolMemberSyncPreviewRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := newStubAdminService()
	h := NewUpstreamPoolHandler(svc)
	router := gin.New()
	router.POST("/admin/upstream-pools/:id/member-sync/preview", h.PreviewMemberSync)

	req := httptest.NewRequest(http.MethodPost, "/admin/upstream-pools/11/member-sync/preview", bytes.NewBufferString(`{"mode":"overwrite_scheduler_fields"}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, body=%s", rec.Code, rec.Body.String())
	}
	var body struct {
		Data struct {
			PoolID int64  `json:"pool_id"`
			Mode   string `json:"mode"`
		} `json:"data"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}
	if body.Data.PoolID != 11 || body.Data.Mode != "overwrite_scheduler_fields" {
		t.Fatalf("response data = %+v", body.Data)
	}
}

func TestUpstreamPoolCapacityPressures(t *testing.T) {
	gin.SetMode(gin.TestMode)
	hardLimit := 1000
	softShare := 800
	svc := newStubAdminService()
	svc.capacityPressures = []service.UpstreamCapacityPressure{{
		SetID:              7,
		SetName:            "聪明共享容量",
		SetCode:            "smart-capacity",
		Platform:           "openai",
		Enabled:            true,
		CapacityLimit:      3000,
		CurrentConcurrency: 2400,
		AvailableCapacity:  600,
		WaitingCount:       2,
		GroupFullCount:     3,
		MemberFullCount:    4,
		BorrowedSlotCount:  5,
		Members: []service.UpstreamCapacityMemberPressure{{
			AccountID:            11,
			AccountName:          "smart-1",
			HardConcurrencyLimit: &hardLimit,
			SoftConcurrencyShare: &softShare,
			CurrentConcurrency:   900,
			WaitingCount:         1,
			LoadRate:             90,
		}},
	}}
	h := NewUpstreamPoolHandler(svc)
	router := gin.New()
	router.GET("/admin/upstream-pools/capacity-pressures", h.GetCapacityPressures)

	req := httptest.NewRequest(http.MethodGet, "/admin/upstream-pools/capacity-pressures", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, body=%s", rec.Code, rec.Body.String())
	}
	var body struct {
		Data []struct {
			SetID             int64 `json:"set_id"`
			CapacityLimit     int   `json:"capacity_limit"`
			BorrowedSlotCount int   `json:"borrowed_slot_count"`
			Members           []struct {
				AccountID int64 `json:"account_id"`
				LoadRate  int   `json:"load_rate"`
			} `json:"members"`
		} `json:"data"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}
	if len(body.Data) != 1 || body.Data[0].SetID != 7 || body.Data[0].CapacityLimit != 3000 || body.Data[0].BorrowedSlotCount != 5 {
		t.Fatalf("response data = %+v", body.Data)
	}
	if len(body.Data[0].Members) != 1 || body.Data[0].Members[0].AccountID != 11 || body.Data[0].Members[0].LoadRate != 90 {
		t.Fatalf("response members = %+v", body.Data[0].Members)
	}
}

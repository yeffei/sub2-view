package admin

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
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

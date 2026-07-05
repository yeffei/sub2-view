package admin

import (
	"encoding/json"
	"testing"
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

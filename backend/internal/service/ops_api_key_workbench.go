package service

import (
	"context"
	"strings"
	"time"
)

type APIKeyLatestErrorSummary struct {
	ErrorID        int64     `json:"error_id"`
	CreatedAt      time.Time `json:"created_at"`
	Category       string    `json:"category"`
	Message        string    `json:"message"`
	ReasonCode     string    `json:"reason_code,omitempty"`
	ActionCode     string    `json:"action_code,omitempty"`
	RequestedModel string    `json:"requested_model,omitempty"`
	UpstreamModel  string    `json:"upstream_model,omitempty"`
	Retryable      bool      `json:"retryable,omitempty"`
	Temporary      bool      `json:"temporary,omitempty"`
}

type APIKeyErrorSummary struct {
	APIKeyID      int64                    `json:"api_key_id"`
	ErrorCount24h int64                    `json:"error_count_24h"`
	LatestError   *APIKeyLatestErrorSummary `json:"latest_error,omitempty"`
}

func (s *OpsService) GetUserAPIKeyErrorSummaries(ctx context.Context, userID int64, apiKeyIDs []int64, startTime time.Time) (map[int64]*APIKeyErrorSummary, error) {
	if s == nil || s.opsRepo == nil {
		return map[int64]*APIKeyErrorSummary{}, nil
	}
	return s.opsRepo.GetUserAPIKeyErrorSummaries(ctx, userID, apiKeyIDs, startTime)
}

func newAPIKeyLatestErrorSummary(detail *OpsErrorLogDetail) *APIKeyLatestErrorSummary {
	if detail == nil {
		return nil
	}
	diagnosis := buildUserErrorDiagnosis(detail)
	model := strings.TrimSpace(detail.RequestedModel)
	if model == "" {
		model = strings.TrimSpace(detail.Model)
	}
	out := &APIKeyLatestErrorSummary{
		ErrorID:        detail.ID,
		CreatedAt:      detail.CreatedAt,
		Category:       MapUserErrorCategory(detail.Phase, detail.Type),
		Message:        strings.TrimSpace(detail.Message),
		RequestedModel: model,
		UpstreamModel:  strings.TrimSpace(detail.UpstreamModel),
	}
	if diagnosis != nil {
		out.ReasonCode = diagnosis.ReasonCode
		out.ActionCode = diagnosis.ActionCode
		out.Retryable = diagnosis.Retryable
		out.Temporary = diagnosis.Temporary
		if out.RequestedModel == "" {
			out.RequestedModel = diagnosis.RequestedModel
		}
		if out.UpstreamModel == "" {
			out.UpstreamModel = diagnosis.UpstreamModel
		}
	}
	return out
}

// NewAPIKeyLatestErrorSummaryForRepo lets repository-level aggregate queries
// reuse the same user-safe diagnosis mapping without duplicating it there.
func NewAPIKeyLatestErrorSummaryForRepo(detail *OpsErrorLogDetail) *APIKeyLatestErrorSummary {
	return newAPIKeyLatestErrorSummary(detail)
}

package service

import (
	"context"
	"strings"
)

// ModelAvailabilityDiagnosis distinguishes a model configuration miss from a
// transient no-account condition.
type ModelAvailabilityDiagnosis struct {
	HasAccountsInPool bool
	HasModelSupport   bool
}

// ModelAvailabilityDiagnoser is shared by gateway handlers when classifying
// no-available-account failures.
type ModelAvailabilityDiagnoser interface {
	DiagnoseModelAvailabilityForPlatform(ctx context.Context, groupID *int64, requestedModel, platform string) ModelAvailabilityDiagnosis
}

// DiagnoseModelAvailabilityForPlatform ignores transient account state. On an
// internal lookup failure it keeps the conservative service-unavailable path.
func (s *GatewayService) DiagnoseModelAvailabilityForPlatform(ctx context.Context, groupID *int64, requestedModel, platform string) ModelAvailabilityDiagnosis {
	if s == nil {
		return ModelAvailabilityDiagnosis{HasAccountsInPool: true, HasModelSupport: true}
	}
	requestedModel = strings.TrimSpace(requestedModel)
	platform = strings.TrimSpace(platform)
	if requestedModel == "" || platform == "" {
		return ModelAvailabilityDiagnosis{HasAccountsInPool: true, HasModelSupport: true}
	}

	accounts, _, err := s.listSchedulableAccounts(ctx, groupID, platform, false)
	if err != nil {
		return ModelAvailabilityDiagnosis{HasAccountsInPool: true, HasModelSupport: true}
	}

	diag := ModelAvailabilityDiagnosis{}
	for i := range accounts {
		diag.HasAccountsInPool = true
		if s.isModelSupportedByAccountWithContext(ctx, &accounts[i], requestedModel) {
			diag.HasModelSupport = true
			return diag
		}
	}
	return diag
}

package service

import "strings"

const (
	AccountAnomalyReasonExpiredPause      = "expired_pause"
	AccountAnomalyReasonTempUnschedulable = "temp_unschedulable"
	AccountAnomalyReasonRateLimited       = "rate_limited"
	AccountAnomalyReasonOverloaded        = "overloaded"
	AccountAnomalyReasonAuthFailed        = "auth_failed"
	AccountAnomalyReasonErrorState        = "error_state"
	AccountAnomalyReasonQuotaExhausted    = "quota_exhausted"
	AccountAnomalyReasonManualPause       = "manual_pause"
	AccountAnomalyReasonInactive          = "inactive"
)

var accountAuthFailureKeywords = []string{
	"401",
	"403",
	"unauthorized",
	"forbidden",
	"invalid api key",
	"api key not found",
	"authentication",
	"token expired",
	"refresh token",
	"invalid_grant",
	"session expired",
	"login",
}

func IsAccountAuthFailureMessage(message string) bool {
	normalized := strings.ToLower(strings.TrimSpace(message))
	if normalized == "" {
		return false
	}
	for _, keyword := range accountAuthFailureKeywords {
		if strings.Contains(normalized, keyword) {
			return true
		}
	}
	return false
}

package service

import (
	"encoding/json"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const userErrorBodyMaxRunes = 800

var htmlTagPattern = regexp.MustCompile(`(?s)<[^>]*>`)

// UserErrorRequest 是面向终端用户的"错误请求"精简脱敏视图（白名单）。
// 严禁包含 client_ip / user_agent / account / api_key_prefix / upstream_endpoint /
// user_email 等敏感或内部字段。注：message（网关标准化错误描述）与 key_name
// （用户自有 API Key 名称，KeysView 中本就可见）经产品决策对该用户开放；
// error_body 仅在详情接口（GetUserErrorRequestDetail）按归属校验后返回。
type UserErrorRequest struct {
	ID              int64     `json:"id"`
	CreatedAt       time.Time `json:"created_at"`
	Model           string    `json:"model"`
	InboundEndpoint string    `json:"inbound_endpoint"`
	StatusCode      int       `json:"status_code"`
	Category        string    `json:"category"`
	Platform        string    `json:"platform"`
	Message         string    `json:"message"`
	KeyName         string    `json:"key_name"`
	KeyDeleted      bool      `json:"key_deleted"`
}

// UserErrorRequestList 是用户错误请求分页结果。
type UserErrorRequestList struct {
	Items    []*UserErrorRequest `json:"items"`
	Total    int                 `json:"total"`
	Page     int                 `json:"page"`
	PageSize int                 `json:"page_size"`
}

// MapUserErrorCategory 把后端 error_phase + error_type 映射为用户侧粗分类码。
// 返回的是稳定的分类 code（前端做 i18n），不是展示文案。
func MapUserErrorCategory(phase, errType string) string {
	switch phase {
	case "auth":
		return "auth"
	case "routing":
		return "service_unavailable"
	case "upstream", "network":
		return "upstream"
	case "internal":
		return "internal"
	case "request":
		switch errType {
		case "rate_limit_error":
			return "rate_limit"
		case "billing_error", "subscription_error":
			return "quota"
		case "invalid_request_error":
			return "invalid_request"
		case "cyber_policy":
			return "cyber"
		}
	}
	return "other"
}

// CategoryToFilter 把用户侧分类码反向映射为后端过滤条件（plain ANY）。
// 未知分类返回两个空切片（即不施加分类过滤）。
// 注意："other" 与未知分类都走 default 返回空切片——"other" 无对应的 phase/type 组合，无法精确反查，因此等价于不过滤。
func CategoryToFilter(category string) (phases []string, errorTypes []string) {
	switch category {
	case "auth":
		return []string{"auth"}, nil
	case "service_unavailable":
		return []string{"routing"}, nil
	case "upstream":
		return []string{"upstream", "network"}, nil
	case "internal":
		return []string{"internal"}, nil
	case "rate_limit":
		return nil, []string{"rate_limit_error"}
	case "quota":
		return nil, []string{"billing_error", "subscription_error"}
	case "invalid_request":
		return nil, []string{"invalid_request_error"}
	case "cyber":
		return []string{"request"}, []string{"cyber_policy"}
	default:
		return nil, nil
	}
}

// ToUserErrorRequest 把内部 OpsErrorLog 裁剪为用户安全视图。
func ToUserErrorRequest(e *OpsErrorLog) *UserErrorRequest {
	if e == nil {
		return nil
	}
	model := e.RequestedModel
	if model == "" {
		model = e.Model
	}
	return &UserErrorRequest{
		ID:              e.ID,
		CreatedAt:       e.CreatedAt,
		Model:           model,
		InboundEndpoint: e.InboundEndpoint,
		StatusCode:      e.StatusCode,
		Category:        MapUserErrorCategory(e.Phase, e.Type),
		Platform:        e.Platform,
		Message:         summarizeUserErrorMessage(e.Message, e.StatusCode),
		KeyName:         e.APIKeyName,
		KeyDeleted:      e.APIKeyDeleted,
	}
}

// UserErrorRequestDetail 是错误请求详情的脱敏视图(点击单行查看)。
// 在 UserErrorRequest 基础上额外暴露 error_body(上游错误响应正文)与 upstream_status_code;
// 仍严禁任何内部/敏感字段。
type UserErrorRequestDetail struct {
	UserErrorRequest
	ErrorBody          string `json:"error_body"`
	ErrorBodyKind      string `json:"error_body_kind,omitempty"`
	UpstreamStatusCode *int   `json:"upstream_status_code,omitempty"`
	Diagnosis          *UserErrorDiagnosis `json:"diagnosis,omitempty"`
}

// UserErrorDiagnosis is a user-safe, product-facing diagnosis derived from the
// internal ops error detail. It exposes only stable machine codes and a small
// amount of non-sensitive evidence for frontend explanation rendering.
type UserErrorDiagnosis struct {
	ReasonCode    string `json:"reason_code"`
	ActionCode    string `json:"action_code,omitempty"`
	Retryable     bool   `json:"retryable,omitempty"`
	Temporary     bool   `json:"temporary,omitempty"`
	RequestedModel string `json:"requested_model,omitempty"`
	UpstreamModel string `json:"upstream_model,omitempty"`
}

// ToUserErrorRequestDetail 把内部 OpsErrorLogDetail 裁剪为用户安全详情视图。
func ToUserErrorRequestDetail(e *OpsErrorLogDetail) *UserErrorRequestDetail {
	if e == nil {
		return nil
	}
	base := ToUserErrorRequest(&e.OpsErrorLog)
	body, kind := summarizeUserErrorBody(e.ErrorBody, e.UpstreamStatusCode)
	return &UserErrorRequestDetail{
		UserErrorRequest:   *base,
		ErrorBody:          body,
		ErrorBodyKind:      kind,
		UpstreamStatusCode: e.UpstreamStatusCode,
		Diagnosis:          buildUserErrorDiagnosis(e),
	}
}

func summarizeUserErrorBody(body string, upstreamStatusCode *int) (string, string) {
	trimmed := strings.TrimSpace(body)
	if trimmed == "" {
		return "", ""
	}

	if looksLikeHTML(trimmed) {
		return truncateUserErrorBody(summarizeHTMLBody(trimmed, upstreamStatusCode)), "html"
	}

	var parsed any
	if json.Unmarshal([]byte(trimmed), &parsed) == nil {
		if message := extractJSONErrorMessage(parsed); message != "" {
			return truncateUserErrorBody(message), "json"
		}
	}

	return truncateUserErrorBody(normalizeWhitespace(trimmed)), "text"
}

func summarizeUserErrorMessage(message string, statusCode int) string {
	trimmed := strings.TrimSpace(message)
	if trimmed == "" {
		return ""
	}
	if looksLikeHTML(trimmed) {
		status := statusCode
		if status <= 0 {
			return truncateUserErrorBody(summarizeHTMLBody(trimmed, nil))
		}
		return truncateUserErrorBody(summarizeHTMLBody(trimmed, &status))
	}
	return truncateUserErrorBody(normalizeWhitespace(trimmed))
}

func looksLikeHTML(body string) bool {
	lower := strings.ToLower(strings.TrimSpace(body))
	return strings.HasPrefix(lower, "<!doctype html") ||
		strings.HasPrefix(lower, "<html") ||
		strings.Contains(lower, "<head") ||
		strings.Contains(lower, "<body") ||
		strings.Contains(lower, "</html>")
}

func summarizeHTMLBody(body string, upstreamStatusCode *int) string {
	parts := make([]string, 0, 3)
	if title := extractHTMLTitle(body); title != "" {
		parts = append(parts, title)
	}
	if upstreamStatusCode != nil && *upstreamStatusCode > 0 {
		parts = append(parts, "HTTP "+strconv.Itoa(*upstreamStatusCode))
	}
	text := htmlTagPattern.ReplaceAllString(body, " ")
	text = strings.NewReplacer("&nbsp;", " ", "&amp;", "&", "&lt;", "<", "&gt;", ">", "&quot;", "\"", "&#39;", "'").Replace(text)
	text = normalizeWhitespace(text)
	if text != "" {
		parts = append(parts, text)
	}
	if len(parts) == 0 {
		return "上游返回了 HTML 错误页"
	}
	return strings.Join(parts, " · ")
}

func extractHTMLTitle(body string) string {
	lower := strings.ToLower(body)
	start := strings.Index(lower, "<title")
	if start < 0 {
		return ""
	}
	gt := strings.Index(body[start:], ">")
	if gt < 0 {
		return ""
	}
	contentStart := start + gt + 1
	end := strings.Index(strings.ToLower(body[contentStart:]), "</title>")
	if end < 0 {
		return ""
	}
	return normalizeWhitespace(body[contentStart : contentStart+end])
}

func extractJSONErrorMessage(value any) string {
	switch typed := value.(type) {
	case map[string]any:
		if errorValue, ok := typed["error"]; ok {
			if message := extractJSONErrorMessage(errorValue); message != "" {
				return message
			}
		}
		for _, key := range []string{"message", "error_message", "detail", "description"} {
			if raw, ok := typed[key].(string); ok && strings.TrimSpace(raw) != "" {
				return normalizeWhitespace(raw)
			}
		}
	case string:
		return normalizeWhitespace(typed)
	}
	return ""
}

func normalizeWhitespace(value string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(value)), " ")
}

func truncateUserErrorBody(value string) string {
	runes := []rune(strings.TrimSpace(value))
	if len(runes) <= userErrorBodyMaxRunes {
		return string(runes)
	}
	return string(runes[:userErrorBodyMaxRunes]) + "..."
}

func buildUserErrorDiagnosis(e *OpsErrorLogDetail) *UserErrorDiagnosis {
	if e == nil {
		return nil
	}

	category := MapUserErrorCategory(e.Phase, e.Type)
	combined := strings.ToLower(strings.Join([]string{
		e.Message,
		e.UpstreamErrorMessage,
		e.UpstreamErrorDetail,
		e.ErrorBody,
	}, " "))

	out := &UserErrorDiagnosis{
		RequestedModel: strings.TrimSpace(e.RequestedModel),
		UpstreamModel:  strings.TrimSpace(e.UpstreamModel),
	}
	if out.RequestedModel == "" {
		out.RequestedModel = strings.TrimSpace(e.Model)
	}

	switch category {
	case "auth":
		out.ActionCode = "keys_connection_test"
		if e.APIKeyDeleted {
			out.ReasonCode = "auth_key_deleted"
			return out
		}
		out.ReasonCode = "auth_invalid_credentials"
		return out
	case "quota":
		out.ActionCode = "profile_balance_notify"
		if e.Type == "subscription_error" || containsAny(combined, "subscription", "plan", "expired", "inactive") {
			out.ReasonCode = "quota_subscription_exhausted"
			return out
		}
		out.ReasonCode = "quota_balance_exhausted"
		return out
	case "rate_limit":
		out.ReasonCode = "rate_limit_window_exhausted"
		out.ActionCode = "dashboard_requests_focus"
		out.Retryable = true
		out.Temporary = true
		return out
	case "invalid_request":
		out.ActionCode = "usage_review_payload"
		if containsAny(combined,
			"model not found",
			"unsupported model",
			"does not exist",
			"unknown model",
			"not support this model",
		) {
			out.ReasonCode = "request_model_not_supported"
			return out
		}
		if containsAny(combined,
			"context length",
			"prompt too long",
			"maximum context",
			"too many tokens",
			"max tokens",
			"payload too large",
			"request too large",
		) {
			out.ReasonCode = "request_payload_too_large"
			return out
		}
		out.ReasonCode = "request_invalid"
		return out
	case "service_unavailable":
		out.ActionCode = "usage_retry_later"
		modelUnsupportedCount := extractSelectionFailureCount(combined, "model_unsupported")
		modelRateLimitedCount := extractSelectionFailureCount(combined, "model_rate_limited")
		if out.RequestedModel != "" && (modelUnsupportedCount > 0 || containsAny(combined,
			"channel pricing restriction",
		)) {
			out.ReasonCode = "service_model_not_available"
			return out
		}
		if modelRateLimitedCount > 0 {
			out.ReasonCode = "service_model_rate_limited"
			out.Retryable = true
			out.Temporary = true
			return out
		}
		if out.RequestedModel != "" && containsAny(combined,
			"supporting model:",
			"supporting model ",
		) {
			out.ReasonCode = "service_model_not_available"
			return out
		}
		out.ReasonCode = "service_no_route_available"
		out.Retryable = true
		out.Temporary = true
		return out
	case "upstream":
		out.ActionCode = "usage_retry_later"
		out.Retryable = true
		out.Temporary = true
		if e.UpstreamStatusCode != nil && *e.UpstreamStatusCode >= 500 {
			out.ReasonCode = "upstream_temporarily_unavailable"
			return out
		}
		out.ReasonCode = "upstream_transport_error"
		return out
	case "internal":
		out.ReasonCode = "internal_gateway_error"
		out.ActionCode = "usage_retry_later"
		return out
	case "cyber":
		out.ReasonCode = "cyber_policy_blocked"
		out.ActionCode = "usage_review_cyber"
		return out
	default:
		out.ReasonCode = "unknown"
		out.ActionCode = "usage_retry_later"
		return out
	}
}

func containsAny(text string, needles ...string) bool {
	for _, needle := range needles {
		if needle != "" && strings.Contains(text, needle) {
			return true
		}
	}
	return false
}

func extractSelectionFailureCount(text string, key string) int {
	if text == "" || key == "" {
		return 0
	}
	marker := key + "="
	idx := strings.Index(text, marker)
	if idx < 0 {
		return 0
	}
	start := idx + len(marker)
	end := start
	for end < len(text) && text[end] >= '0' && text[end] <= '9' {
		end++
	}
	if end == start {
		return 0
	}
	value, err := strconv.Atoi(text[start:end])
	if err != nil {
		return 0
	}
	return value
}

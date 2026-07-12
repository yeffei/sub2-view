package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const upstreamAccountMetaBodyLimit int64 = 1 << 20

type UpstreamAccountMeta struct {
	Compatible          bool
	Platform            string
	GroupID             *int64
	GroupName           string
	RateMultiplier      float64
	RateSource          string
	SubscriptionType    string
	ImageRateMultiplier float64
}

type UpstreamAccountMetaSyncErrorKind string

const (
	UpstreamAccountMetaSyncErrorConfiguration UpstreamAccountMetaSyncErrorKind = "configuration"
	UpstreamAccountMetaSyncErrorUnsupported   UpstreamAccountMetaSyncErrorKind = "unsupported"
	UpstreamAccountMetaSyncErrorUpstream      UpstreamAccountMetaSyncErrorKind = "upstream"
)

type UpstreamAccountMetaSyncError struct {
	Kind    UpstreamAccountMetaSyncErrorKind
	Message string
	Err     error
}

func (e *UpstreamAccountMetaSyncError) Error() string {
	if e == nil {
		return ""
	}
	if e.Err == nil {
		return e.Message
	}
	return e.Message + ": " + e.Err.Error()
}

func (e *UpstreamAccountMetaSyncError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Err
}

func (e *UpstreamAccountMetaSyncError) SafeMessage() string {
	if e == nil || strings.TrimSpace(e.Message) == "" {
		return "Failed to sync upstream account metadata"
	}
	return e.Message
}

func newUpstreamAccountMetaConfigError(message string, err error) error {
	return &UpstreamAccountMetaSyncError{Kind: UpstreamAccountMetaSyncErrorConfiguration, Message: message, Err: err}
}

func newUpstreamAccountMetaUnsupportedError(message string, err error) error {
	return &UpstreamAccountMetaSyncError{Kind: UpstreamAccountMetaSyncErrorUnsupported, Message: message, Err: err}
}

func newUpstreamAccountMetaUpstreamError(message string, err error) error {
	return &UpstreamAccountMetaSyncError{Kind: UpstreamAccountMetaSyncErrorUpstream, Message: message, Err: err}
}

func (s *AccountTestService) FetchCompatibleUpstreamAccountMeta(ctx context.Context, account *Account) (*UpstreamAccountMeta, error) {
	if s == nil {
		return nil, newUpstreamAccountMetaConfigError("Account test service is not configured", nil)
	}
	if account == nil {
		return nil, newUpstreamAccountMetaConfigError("Account is required", nil)
	}
	if account.Type != AccountTypeAPIKey {
		return nil, newUpstreamAccountMetaUnsupportedError("Upstream rate sync only supports API-key accounts", nil)
	}
	if s.httpUpstream == nil {
		return nil, newUpstreamAccountMetaConfigError("Upstream HTTP client is not configured", nil)
	}

	apiKey := strings.TrimSpace(account.GetCredential("api_key"))
	if apiKey == "" {
		apiKey = strings.TrimSpace(account.GetOpenAIApiKey())
	}
	if apiKey == "" {
		return nil, newUpstreamAccountMetaConfigError("No upstream API key is available", nil)
	}

	baseURL := ""
	if account.IsOpenAI() {
		baseURL = account.GetOpenAIBaseURL()
	} else {
		baseURL = account.GetBaseURL()
	}
	if strings.TrimSpace(baseURL) == "" {
		return nil, newUpstreamAccountMetaConfigError("Upstream base URL is required for rate sync", nil)
	}
	normalizedBaseURL, err := s.validateUpstreamBaseURL(baseURL)
	if err != nil {
		return nil, newUpstreamAccountMetaConfigError("Invalid upstream base URL", err)
	}

	metaURL := buildCompatibleAccountMetaURL(normalizedBaseURL)
	meta, err := s.fetchCompatibleUpstreamAccountMetaURL(ctx, account, metaURL, apiKey)
	if err == nil {
		return meta, nil
	}
	var syncErr *UpstreamAccountMetaSyncError
	if !errors.As(err, &syncErr) || syncErr.Kind != UpstreamAccountMetaSyncErrorUnsupported {
		return nil, err
	}

	usageURL := buildCompatibleUsageURL(normalizedBaseURL)
	meta, usageErr := s.fetchCompatibleUpstreamAccountMetaURL(ctx, account, usageURL, apiKey)
	if usageErr != nil {
		var usageSyncErr *UpstreamAccountMetaSyncError
		if errors.As(usageErr, &usageSyncErr) && usageSyncErr.Kind == UpstreamAccountMetaSyncErrorUnsupported {
			newAPIMeta, newAPIErr := s.probeNewAPIUpstreamRateMultiplier(ctx, account, normalizedBaseURL, apiKey)
			if newAPIErr == nil {
				return newAPIMeta, nil
			}
			var newAPISyncErr *UpstreamAccountMetaSyncError
			if !errors.As(newAPIErr, &newAPISyncErr) || newAPISyncErr.Kind != UpstreamAccountMetaSyncErrorUnsupported {
				return nil, newAPIErr
			}

			legacyMeta, legacyErr := s.probeLegacyUpstreamRateMultiplier(ctx, account, normalizedBaseURL, apiKey)
			if legacyErr == nil {
				return legacyMeta, nil
			}
			var legacySyncErr *UpstreamAccountMetaSyncError
			if !errors.As(legacyErr, &legacySyncErr) || legacySyncErr.Kind != UpstreamAccountMetaSyncErrorUnsupported {
				return nil, legacyErr
			}

			pricingURL := buildCompatibleModelPricingURL(normalizedBaseURL)
			if s.probeCompatibleModelPricingURL(ctx, account, pricingURL, apiKey) {
				return nil, newUpstreamAccountMetaUnsupportedError("上游只暴露模型定价 /api/pricing，不提供账号级倍率元数据；不能自动写入账号倍率", usageErr)
			}
			return nil, newUpstreamAccountMetaUnsupportedError("上游未暴露兼容的账号倍率接口", usageErr)
		}
		return nil, usageErr
	}
	return meta, nil
}

func (s *AccountTestService) fetchCompatibleUpstreamAccountMetaURL(ctx context.Context, account *Account, targetURL string, apiKey string) (*UpstreamAccountMeta, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, targetURL, nil)
	if err != nil {
		return nil, newUpstreamAccountMetaConfigError("Invalid upstream account metadata URL", err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	proxyURL := upstreamModelsProxyURL(account)
	resp, err := s.doUpstreamModelsRequest(req, proxyURL, account)
	if err != nil {
		return nil, newUpstreamAccountMetaUpstreamError("Failed to request upstream account metadata", err)
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(io.LimitReader(resp.Body, upstreamAccountMetaBodyLimit+1))
	if err != nil {
		return nil, newUpstreamAccountMetaUpstreamError("Failed to read upstream account metadata", err)
	}
	if int64(len(body)) > upstreamAccountMetaBodyLimit {
		return nil, newUpstreamAccountMetaUpstreamError("Upstream account metadata response is too large", fmt.Errorf("response exceeds %d bytes", upstreamAccountMetaBodyLimit))
	}
	if resp.StatusCode == http.StatusNotFound {
		return nil, newUpstreamAccountMetaUnsupportedError("Upstream does not expose compatible account metadata", nil)
	}
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return nil, newUpstreamAccountMetaUpstreamError(
			fmt.Sprintf("Upstream account metadata request failed with HTTP %d", resp.StatusCode),
			fmt.Errorf("upstream account metadata returned HTTP %d", resp.StatusCode),
		)
	}

	meta, err := extractCompatibleUpstreamAccountMeta(body)
	if err != nil {
		return nil, newUpstreamAccountMetaUpstreamError("Upstream account metadata response was not valid JSON", err)
	}
	if meta == nil || !meta.Compatible {
		return nil, newUpstreamAccountMetaUnsupportedError("Upstream account metadata is not compatible", nil)
	}
	if meta.RateMultiplier < 0 {
		return nil, newUpstreamAccountMetaUpstreamError("Upstream returned an invalid rate multiplier", nil)
	}
	return meta, nil
}

func buildCompatibleAccountMetaURL(base string) string {
	normalized := strings.TrimRight(strings.TrimSpace(base), "/")
	lower := strings.ToLower(normalized)
	switch {
	case strings.HasSuffix(lower, "/v1/account/meta"):
		return normalized
	case strings.HasSuffix(lower, "/v1"):
		return normalized + "/account/meta"
	default:
		return normalized + "/v1/account/meta"
	}
}

func buildCompatibleUsageURL(base string) string {
	normalized := strings.TrimRight(strings.TrimSpace(base), "/")
	lower := strings.ToLower(normalized)
	switch {
	case strings.HasSuffix(lower, "/v1/usage"):
		return normalized
	case strings.HasSuffix(lower, "/v1"):
		return normalized + "/usage"
	default:
		return normalized + "/v1/usage"
	}
}

func buildCompatibleModelPricingURL(base string) string {
	normalized := strings.TrimRight(strings.TrimSpace(base), "/")
	lower := strings.ToLower(normalized)
	switch {
	case strings.HasSuffix(lower, "/api/pricing"):
		return normalized
	case strings.HasSuffix(lower, "/v1"):
		return strings.TrimRight(normalized[:len(normalized)-3], "/") + "/api/pricing"
	default:
		return normalized + "/api/pricing"
	}
}

func (s *AccountTestService) probeCompatibleModelPricingURL(ctx context.Context, account *Account, targetURL string, apiKey string) bool {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, targetURL, nil)
	if err != nil {
		return false
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	resp, err := s.doUpstreamModelsRequest(req, upstreamModelsProxyURL(account), account)
	if err != nil {
		return false
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return false
	}
	body, err := io.ReadAll(io.LimitReader(resp.Body, upstreamAccountMetaBodyLimit+1))
	if err != nil || int64(len(body)) > upstreamAccountMetaBodyLimit {
		return false
	}

	var payload struct {
		Data []struct {
			ModelName  string   `json:"model_name"`
			ModelRatio *float64 `json:"model_ratio"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &payload); err != nil {
		return false
	}
	for _, item := range payload.Data {
		if strings.TrimSpace(item.ModelName) != "" && item.ModelRatio != nil {
			return true
		}
	}
	return false
}

func extractCompatibleUpstreamAccountMeta(body []byte) (*UpstreamAccountMeta, error) {
	var envelope struct {
		Data json.RawMessage `json:"data"`
	}
	raw := body
	if err := json.Unmarshal(body, &envelope); err == nil && len(envelope.Data) > 0 {
		raw = envelope.Data
	}

	var usagePayload struct {
		Billing json.RawMessage `json:"billing"`
	}
	if err := json.Unmarshal(raw, &usagePayload); err == nil && len(usagePayload.Billing) > 0 {
		raw = usagePayload.Billing
	}

	var payload struct {
		Compatible          bool    `json:"compatible"`
		Platform            string  `json:"platform"`
		GroupID             *int64  `json:"group_id"`
		GroupName           string  `json:"group_name"`
		RateMultiplier      float64 `json:"rate_multiplier"`
		RateSource          string  `json:"rate_source"`
		SubscriptionType    string  `json:"subscription_type"`
		ImageRateMultiplier float64 `json:"image_rate_multiplier"`
	}
	if err := json.Unmarshal(raw, &payload); err != nil {
		return nil, fmt.Errorf("parse upstream account metadata: %w", err)
	}
	return &UpstreamAccountMeta{
		Compatible:          payload.Compatible,
		Platform:            strings.TrimSpace(payload.Platform),
		GroupID:             payload.GroupID,
		GroupName:           strings.TrimSpace(payload.GroupName),
		RateMultiplier:      payload.RateMultiplier,
		RateSource:          strings.TrimSpace(payload.RateSource),
		SubscriptionType:    strings.TrimSpace(payload.SubscriptionType),
		ImageRateMultiplier: payload.ImageRateMultiplier,
	}, nil
}

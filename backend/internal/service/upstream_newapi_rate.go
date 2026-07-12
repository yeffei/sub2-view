package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"strings"
	"time"
)

const newAPIRateLogBodyLimit int64 = 8 << 20

type newAPITokenLog struct {
	Type      int             `json:"type"`
	Group     string          `json:"group"`
	ModelName string          `json:"model_name"`
	RequestID string          `json:"request_id"`
	Other     json.RawMessage `json:"other"`
}

func (s *AccountTestService) probeNewAPIUpstreamRateMultiplier(
	ctx context.Context,
	account *Account,
	baseURL string,
	apiKey string,
) (*UpstreamAccountMeta, error) {
	if account == nil || !account.IsOpenAI() || account.Type != AccountTypeAPIKey {
		return nil, newUpstreamAccountMetaUnsupportedError("New API 倍率探测仅支持 OpenAI API-key 账号", nil)
	}

	logURL := buildNewAPITokenLogURL(baseURL)
	initialLogs, err := s.fetchNewAPITokenLogs(ctx, account, logURL, apiKey)
	if err != nil {
		return nil, err
	}

	models, modelErr := s.FetchUpstreamSupportedModels(ctx, account)
	model := ""
	if modelErr == nil {
		model = selectLegacyRateProbeModel(models)
	}
	if model == "" {
		model = selectNewAPIRateProbeModelFromLogs(initialLogs)
	}
	if model == "" {
		return nil, newUpstreamAccountMetaUnsupportedError("New API 上游没有返回可用于低成本倍率探测的 Responses 模型", modelErr)
	}

	requestID, err := s.sendRateProbeRequest(ctx, account, baseURL, apiKey, model)
	if err != nil {
		return nil, err
	}
	if requestID == "" {
		return nil, newUpstreamAccountMetaUnsupportedError("上游未返回可用于匹配计费日志的 request id", nil)
	}

	for attempt := 0; attempt < legacyRateProbePollAttempts; attempt++ {
		select {
		case <-ctx.Done():
			return nil, newUpstreamAccountMetaUpstreamError("等待 New API 倍率探测日志时请求已取消", ctx.Err())
		case <-time.After(legacyRateProbePollInterval):
		}

		logs, fetchErr := s.fetchNewAPITokenLogs(ctx, account, logURL, apiKey)
		if fetchErr != nil {
			return nil, fetchErr
		}
		for _, logEntry := range logs {
			if strings.TrimSpace(logEntry.RequestID) != requestID {
				continue
			}
			rate, parseErr := extractNewAPILogRateMultiplier(logEntry.Other)
			if parseErr != nil {
				return nil, newUpstreamAccountMetaUpstreamError("New API 探测日志缺少有效分组倍率", parseErr)
			}
			return &UpstreamAccountMeta{
				Compatible:     true,
				Platform:       account.Platform,
				GroupName:      strings.TrimSpace(logEntry.Group),
				RateMultiplier: rate,
				RateSource:     "new_api_token_log_probe",
			}, nil
		}
	}

	return nil, newUpstreamAccountMetaUpstreamError("New API 未及时生成倍率探测日志", nil)
}

func selectNewAPIRateProbeModelFromLogs(logs []newAPITokenLog) string {
	models := make([]string, 0, len(logs))
	for _, logEntry := range logs {
		if logEntry.Type != 2 {
			continue
		}
		model := strings.TrimSpace(logEntry.ModelName)
		if model != "" {
			models = append(models, model)
		}
	}
	return selectLegacyRateProbeModel(models)
}

func (s *AccountTestService) fetchNewAPITokenLogs(
	ctx context.Context,
	account *Account,
	targetURL string,
	apiKey string,
) ([]newAPITokenLog, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, targetURL, nil)
	if err != nil {
		return nil, newUpstreamAccountMetaConfigError("Invalid New API token log URL", err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)
	account.ApplyHeaderOverrides(req.Header)

	resp, err := s.doUpstreamModelsRequest(req, upstreamModelsProxyURL(account), account)
	if err != nil {
		return nil, newUpstreamAccountMetaUpstreamError("请求 New API Key 计费日志失败", err)
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(io.LimitReader(resp.Body, newAPIRateLogBodyLimit+1))
	if err != nil {
		return nil, newUpstreamAccountMetaUpstreamError("读取 New API Key 计费日志失败", err)
	}
	if int64(len(body)) > newAPIRateLogBodyLimit {
		return nil, newUpstreamAccountMetaUpstreamError("New API Key 计费日志响应过大", nil)
	}
	if resp.StatusCode == http.StatusNotFound || resp.StatusCode == http.StatusMethodNotAllowed {
		return nil, newUpstreamAccountMetaUnsupportedError("上游未提供 New API Key 级计费日志接口", nil)
	}
	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return nil, newUpstreamAccountMetaUnsupportedError("上游不允许 API Key 读取自身计费日志", nil)
	}
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return nil, newUpstreamAccountMetaUpstreamError(
			fmt.Sprintf("New API Key 计费日志请求返回 HTTP %d", resp.StatusCode), nil,
		)
	}

	var payload struct {
		Success bool             `json:"success"`
		Message string           `json:"message"`
		Data    []newAPITokenLog `json:"data"`
	}
	if err := json.Unmarshal(body, &payload); err != nil {
		return nil, newUpstreamAccountMetaUnsupportedError("上游 Key 日志不是兼容的 New API 格式", err)
	}
	if !payload.Success {
		return nil, newUpstreamAccountMetaUnsupportedError("上游不支持读取 New API Key 计费日志", nil)
	}
	return payload.Data, nil
}

func extractNewAPILogRateMultiplier(raw json.RawMessage) (float64, error) {
	if len(raw) == 0 {
		return 0, fmt.Errorf("empty log metadata")
	}
	var encoded string
	if err := json.Unmarshal(raw, &encoded); err == nil {
		raw = json.RawMessage(encoded)
	}
	var other struct {
		GroupRatio     *float64 `json:"group_ratio"`
		UserGroupRatio *float64 `json:"user_group_ratio"`
	}
	if err := json.Unmarshal(raw, &other); err != nil {
		return 0, fmt.Errorf("parse log metadata: %w", err)
	}
	rate := other.GroupRatio
	if other.UserGroupRatio != nil && *other.UserGroupRatio >= 0 {
		rate = other.UserGroupRatio
	}
	if rate == nil || math.IsNaN(*rate) || math.IsInf(*rate, 0) || *rate < 0 || *rate > 100 {
		return 0, fmt.Errorf("invalid group ratio")
	}
	return math.Round(*rate*10000) / 10000, nil
}

func buildNewAPITokenLogURL(base string) string {
	normalized := strings.TrimRight(strings.TrimSpace(base), "/")
	lower := strings.ToLower(normalized)
	if strings.HasSuffix(lower, "/v1") {
		normalized = strings.TrimRight(normalized[:len(normalized)-3], "/")
	}
	return normalized + "/api/log/token"
}

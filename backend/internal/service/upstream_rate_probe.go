package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"strings"
	"time"
)

const (
	legacyRateProbeBodyLimit    int64 = 1 << 20
	legacyRateProbePollInterval       = 500 * time.Millisecond
	legacyRateProbePollAttempts       = 16
)

type legacyUsageSnapshot struct {
	Requests   int64
	Cost       float64
	ActualCost float64
}

func (s *AccountTestService) probeLegacyUpstreamRateMultiplier(
	ctx context.Context,
	account *Account,
	baseURL string,
	apiKey string,
) (*UpstreamAccountMeta, error) {
	if account == nil || !account.IsOpenAI() || account.Type != AccountTypeAPIKey {
		return nil, newUpstreamAccountMetaUnsupportedError("旧版用量差值倍率探测仅支持 OpenAI API-key 账号", nil)
	}

	usageURL := buildCompatibleUsageURL(baseURL)
	before, err := s.fetchLegacyUsageSnapshot(ctx, account, usageURL, apiKey)
	if err != nil {
		return nil, err
	}

	models, err := s.FetchUpstreamSupportedModels(ctx, account)
	if err != nil {
		return nil, newUpstreamAccountMetaUpstreamError("无法获取倍率探测所需的上游模型列表", err)
	}
	model := selectLegacyRateProbeModel(models)
	if model == "" {
		return nil, newUpstreamAccountMetaUnsupportedError("上游没有适合低成本倍率探测的 Responses 模型", nil)
	}

	if _, err := s.sendRateProbeRequest(ctx, account, baseURL, apiKey, model); err != nil {
		return nil, err
	}

	var after *legacyUsageSnapshot
	for attempt := 0; attempt < legacyRateProbePollAttempts; attempt++ {
		select {
		case <-ctx.Done():
			return nil, newUpstreamAccountMetaUpstreamError("等待上游倍率探测统计时请求已取消", ctx.Err())
		case <-time.After(legacyRateProbePollInterval):
		}

		after, err = s.fetchLegacyUsageSnapshot(ctx, account, usageURL, apiKey)
		if err != nil {
			return nil, err
		}
		if after.Requests > before.Requests {
			break
		}
	}

	if after == nil || after.Requests <= before.Requests {
		return nil, newUpstreamAccountMetaUpstreamError("上游用量统计尚未记录倍率探测请求", nil)
	}
	if after.Requests != before.Requests+1 {
		return nil, newUpstreamAccountMetaUpstreamError("倍率探测期间该 API Key 存在并发请求，结果已放弃以避免误写", nil)
	}

	costDelta := after.Cost - before.Cost
	actualCostDelta := after.ActualCost - before.ActualCost
	if costDelta <= 0 || actualCostDelta < 0 {
		return nil, newUpstreamAccountMetaUpstreamError("上游用量差值不足，无法计算账号倍率", nil)
	}

	rate := actualCostDelta / costDelta
	if math.IsNaN(rate) || math.IsInf(rate, 0) || rate < 0 || rate > 100 {
		return nil, newUpstreamAccountMetaUpstreamError("上游用量差值计算出了异常账号倍率", nil)
	}
	rate = math.Round(rate*10000) / 10000

	return &UpstreamAccountMeta{
		Compatible:     true,
		Platform:       account.Platform,
		RateMultiplier: rate,
		RateSource:     "legacy_usage_delta_probe",
	}, nil
}

func (s *AccountTestService) fetchLegacyUsageSnapshot(
	ctx context.Context,
	account *Account,
	targetURL string,
	apiKey string,
) (*legacyUsageSnapshot, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, targetURL, nil)
	if err != nil {
		return nil, newUpstreamAccountMetaConfigError("Invalid upstream usage URL", err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)
	account.ApplyHeaderOverrides(req.Header)

	resp, err := s.doUpstreamModelsRequest(req, upstreamModelsProxyURL(account), account)
	if err != nil {
		return nil, newUpstreamAccountMetaUpstreamError("请求上游用量统计失败", err)
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(io.LimitReader(resp.Body, legacyRateProbeBodyLimit+1))
	if err != nil {
		return nil, newUpstreamAccountMetaUpstreamError("读取上游用量统计失败", err)
	}
	if int64(len(body)) > legacyRateProbeBodyLimit {
		return nil, newUpstreamAccountMetaUpstreamError("上游用量统计响应过大", nil)
	}
	if resp.StatusCode == http.StatusNotFound {
		return nil, newUpstreamAccountMetaUnsupportedError("上游未提供 /v1/usage 倍率探测能力", nil)
	}
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return nil, newUpstreamAccountMetaUpstreamError(
			fmt.Sprintf("上游用量统计请求返回 HTTP %d", resp.StatusCode), nil,
		)
	}

	var payload struct {
		Usage struct {
			Today struct {
				Requests   int64   `json:"requests"`
				Cost       float64 `json:"cost"`
				ActualCost float64 `json:"actual_cost"`
			} `json:"today"`
		} `json:"usage"`
	}
	if err := json.Unmarshal(body, &payload); err != nil {
		return nil, newUpstreamAccountMetaUnsupportedError("上游 /v1/usage 不是兼容的用量统计格式", err)
	}
	if payload.Usage.Today.Requests < 0 || payload.Usage.Today.Cost < 0 || payload.Usage.Today.ActualCost < 0 {
		return nil, newUpstreamAccountMetaUpstreamError("上游 /v1/usage 返回了异常用量统计", nil)
	}
	return &legacyUsageSnapshot{
		Requests:   payload.Usage.Today.Requests,
		Cost:       payload.Usage.Today.Cost,
		ActualCost: payload.Usage.Today.ActualCost,
	}, nil
}

func (s *AccountTestService) sendRateProbeRequest(
	ctx context.Context,
	account *Account,
	baseURL string,
	apiKey string,
	model string,
) (string, error) {
	payload, err := json.Marshal(map[string]any{
		"model":             model,
		"input":             "Reply with OK only.",
		"max_output_tokens": 16,
		"store":             false,
	})
	if err != nil {
		return "", newUpstreamAccountMetaConfigError("无法构造倍率探测请求", err)
	}

	probeURL := buildCompatibleResponsesURL(baseURL)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, probeURL, bytes.NewReader(payload))
	if err != nil {
		return "", newUpstreamAccountMetaConfigError("Invalid upstream responses URL", err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)
	account.ApplyHeaderOverrides(req.Header)

	resp, err := s.doUpstreamModelsRequest(req, upstreamModelsProxyURL(account), account)
	if err != nil {
		return "", newUpstreamAccountMetaUpstreamError("发送上游倍率探测请求失败", err)
	}
	defer func() { _ = resp.Body.Close() }()
	_, _ = io.Copy(io.Discard, io.LimitReader(resp.Body, legacyRateProbeBodyLimit))
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return "", newUpstreamAccountMetaUpstreamError(
			fmt.Sprintf("上游倍率探测请求返回 HTTP %d", resp.StatusCode), nil,
		)
	}
	requestID := strings.TrimSpace(resp.Header.Get("X-Oneapi-Request-Id"))
	if requestID == "" {
		requestID = strings.TrimSpace(resp.Header.Get("X-Request-Id"))
	}
	return requestID, nil
}

func selectLegacyRateProbeModel(models []string) string {
	available := make(map[string]string, len(models))
	for _, model := range models {
		trimmed := strings.TrimSpace(model)
		if trimmed != "" {
			available[strings.ToLower(trimmed)] = trimmed
		}
	}
	preferred := []string{
		"gpt-5.4-mini",
		"gpt-5-mini",
		"gpt-5.2",
		"gpt-5.3-codex",
		"gpt-5.4",
		"gpt-5",
	}
	for _, candidate := range preferred {
		if model, ok := available[candidate]; ok {
			return model
		}
	}
	for _, model := range models {
		lower := strings.ToLower(strings.TrimSpace(model))
		if strings.HasPrefix(lower, "gpt-5") && !strings.Contains(lower, "image") && !strings.Contains(lower, "audio") && !strings.Contains(lower, "realtime") {
			return strings.TrimSpace(model)
		}
	}
	return ""
}

func buildCompatibleResponsesURL(base string) string {
	normalized := strings.TrimRight(strings.TrimSpace(base), "/")
	lower := strings.ToLower(normalized)
	switch {
	case strings.HasSuffix(lower, "/v1/responses"):
		return normalized
	case strings.HasSuffix(lower, "/v1"):
		return normalized + "/responses"
	default:
		return normalized + "/v1/responses"
	}
}

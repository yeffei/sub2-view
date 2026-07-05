//go:build unit

package service

import (
	"strings"
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/pkg/openai_compat"
	"github.com/stretchr/testify/require"
)

func TestAccountProbeAPIModeUsesResponsesWhenSupported(t *testing.T) {
	account := &Account{
		Platform: PlatformOpenAI,
		Type:     AccountTypeAPIKey,
		Extra: map[string]any{
			openai_compat.ExtraKeyResponsesSupported: true,
		},
	}

	require.Equal(t, MonitorAPIModeResponses, accountProbeAPIMode(account))
}

func TestRunAccountProbeTarget_ResponsesUsesOnlyHiAndOneToken(t *testing.T) {
	h := &openAICaptureHandler{}
	endpoint := setupFakeOpenAI(t, h)

	row := runAccountProbeTarget(t.Context(), accountProbeTarget{
		AccountID: 1,
		PoolID:    2,
		Provider:  MonitorProviderOpenAI,
		Model:     "gpt-5.4",
		Endpoint:  endpoint,
		APIKey:    "sk-test",
		APIMode:   MonitorAPIModeResponses,
	})

	require.Equal(t, MonitorStatusOperational, row.Status)
	require.Equal(t, providerOpenAIResponsesPath, h.lastPath)
	require.Equal(t, monitorLightweightPrompt, strings.TrimSpace(stringFromAny(h.lastBody["input"])))
	require.Equal(t, float64(monitorLightweightMaxTokens), h.lastBody["max_output_tokens"])
	require.NotContains(t, h.lastBody, "instructions")
	require.Equal(t, "Bearer sk-test", h.lastHeaders.Get("Authorization"))
}

func TestRunAccountProbeTarget_ChatUsesOnlyHiAndOneToken(t *testing.T) {
	h := &openAICaptureHandler{}
	endpoint := setupFakeOpenAI(t, h)

	row := runAccountProbeTarget(t.Context(), accountProbeTarget{
		AccountID: 1,
		PoolID:    2,
		Provider:  MonitorProviderOpenAI,
		Model:     "gpt-5.4",
		Endpoint:  endpoint,
		APIKey:    "sk-test",
		APIMode:   MonitorAPIModeChatCompletions,
	})

	require.Equal(t, MonitorStatusOperational, row.Status)
	require.Equal(t, providerOpenAIPath, h.lastPath)
	messages, _ := h.lastBody["messages"].([]any)
	require.Len(t, messages, 1)
	msg, _ := messages[0].(map[string]any)
	require.Equal(t, monitorLightweightPrompt, strings.TrimSpace(stringFromAny(msg["content"])))
	require.Equal(t, float64(monitorLightweightMaxTokens), h.lastBody["max_tokens"])
	require.Equal(t, "Bearer sk-test", h.lastHeaders.Get("Authorization"))
}

func TestRunTestBackgroundWithMode_LightweightResponsesUsesMinimalViableProbe(t *testing.T) {
	h := &openAICaptureHandler{}
	endpoint := setupFakeOpenAI(t, h)
	account := &Account{
		ID:       10,
		Platform: PlatformOpenAI,
		Type:     AccountTypeAPIKey,
		Credentials: map[string]any{
			"api_key":  "sk-test",
			"base_url": endpoint,
		},
		Extra: map[string]any{
			openai_compat.ExtraKeyResponsesSupported: true,
		},
	}
	svc := &AccountTestService{accountRepo: &stubOpenAIAccountRepo{accounts: []Account{*account}}}

	result, err := svc.RunTestBackgroundWithMode(t.Context(), account.ID, "gpt-5.4", AccountTestModeLightweight)

	require.NoError(t, err)
	require.Equal(t, "success", result.Status)
	require.Equal(t, providerOpenAIResponsesPath, h.lastPath)
	require.Equal(t, monitorLightweightPrompt, strings.TrimSpace(stringFromAny(h.lastBody["input"])))
	require.Equal(t, float64(monitorLightweightMaxTokens), h.lastBody["max_output_tokens"])
	require.NotContains(t, h.lastBody, "instructions")
}

func TestRunTestBackgroundWithMode_LightweightChatUsesMinimalViableProbe(t *testing.T) {
	h := &openAICaptureHandler{}
	endpoint := setupFakeOpenAI(t, h)
	account := &Account{
		ID:       11,
		Platform: PlatformOpenAI,
		Type:     AccountTypeAPIKey,
		Credentials: map[string]any{
			"api_key":  "sk-test",
			"base_url": endpoint,
		},
		Extra: map[string]any{
			openai_compat.ExtraKeyResponsesSupported: false,
		},
	}
	svc := &AccountTestService{accountRepo: &stubOpenAIAccountRepo{accounts: []Account{*account}}}

	result, err := svc.RunTestBackgroundWithMode(t.Context(), account.ID, "gpt-5.4", AccountTestModeLightweight)

	require.NoError(t, err)
	require.Equal(t, "success", result.Status)
	require.Equal(t, providerOpenAIPath, h.lastPath)
	messages, _ := h.lastBody["messages"].([]any)
	require.Len(t, messages, 1)
	msg, _ := messages[0].(map[string]any)
	require.Equal(t, monitorLightweightPrompt, strings.TrimSpace(stringFromAny(msg["content"])))
	require.Equal(t, float64(monitorLightweightMaxTokens), h.lastBody["max_tokens"])
	require.NotContains(t, h.lastBody, "instructions")
}

func TestAccountProbeAPIModeUsesChatWhenResponsesUnsupported(t *testing.T) {
	account := &Account{
		Platform: PlatformOpenAI,
		Type:     AccountTypeAPIKey,
		Extra: map[string]any{
			openai_compat.ExtraKeyResponsesSupported: false,
		},
	}

	require.Equal(t, MonitorAPIModeChatCompletions, accountProbeAPIMode(account))
}

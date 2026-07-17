package service

import (
	"encoding/json"
	"strconv"
	"strings"
)

func (a *Account) IsOpenAICompatible() bool {
	return a != nil && (a.Platform == PlatformOpenAI || a.Platform == PlatformGrok)
}

func (a *Account) IsChatGPTAccountFedRAMP() bool {
	if !a.IsOpenAIOAuth() || a.Credentials == nil {
		return false
	}
	v, ok := a.Credentials["chatgpt_account_is_fedramp"]
	if !ok || v == nil {
		return false
	}
	switch value := v.(type) {
	case bool:
		return value
	case string:
		parsed, err := strconv.ParseBool(strings.TrimSpace(value))
		return err == nil && parsed
	case json.Number:
		parsed, err := strconv.ParseBool(value.String())
		return err == nil && parsed
	case float64:
		return value != 0
	case int:
		return value != 0
	case int64:
		return value != 0
	default:
		return false
	}
}

func (a *Account) IsCodexCLIOnlyAppServerAllowed() bool {
	if !a.IsCodexCLIOnlyEnabled() || a.Extra == nil {
		return false
	}
	v, ok := a.Extra["codex_cli_only_allow_app_server"].(bool)
	return ok && v
}

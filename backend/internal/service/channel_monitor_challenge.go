package service

import (
	"regexp"
	"strings"
)

const monitorChallengePrompt = "Hi"
const monitorChallengeExpected = "hi"

// monitorChallengeWordRegex 提取响应中的英文单词。
var monitorChallengeWordRegex = regexp.MustCompile(`[A-Za-z]+`)

// monitorChallenge 一次 challenge 的 prompt + 期望答案。
type monitorChallenge struct {
	Prompt   string
	Expected string
}

// generateChallenge 返回最小探测 challenge。
// 监控目标只是确认“请求能通、模型能最小单位返回文本”，
// 因此固定发送单词 "Hi" 即可，不再发送 few-shot 算术题。
func generateChallenge() monitorChallenge {
	return monitorChallenge{
		Prompt:   monitorChallengePrompt,
		Expected: monitorChallengeExpected,
	}
}

// validateChallenge 在响应文本中查找 expected 单词，返回是否通过校验。
// 允许上游返回 "Hi." / "hi\n" 这类轻微格式差异，但仍要求至少包含期望单词。
func validateChallenge(responseText, expected string) bool {
	if responseText == "" || expected == "" {
		return false
	}
	matches := monitorChallengeWordRegex.FindAllString(responseText, -1)
	for _, m := range matches {
		if strings.EqualFold(m, expected) {
			return true
		}
	}
	return false
}

package service

import (
	"strings"

	"github.com/gin-gonic/gin"
)

const openAIResponsesAutoUpstreamSessionAnchorGinContextKey = "__openai_responses_auto_upstream_session_anchor"

type openAIUpstreamSessionAnchor struct {
	SessionID string
	Source    string
}

func setOpenAIResponsesAutoUpstreamSessionAnchor(c *gin.Context, body []byte, previousResponseIDPresent bool) {
	if c == nil || previousResponseIDPresent {
		return
	}
	if strings.TrimSpace(explicitOpenAISessionID(c, body)) != "" {
		return
	}
	seed := strings.TrimSpace(deriveOpenAIContentSessionSeed(body))
	if seed == "" {
		return
	}
	c.Set(openAIResponsesAutoUpstreamSessionAnchorGinContextKey, openAIUpstreamSessionAnchor{
		SessionID: seed,
		Source:    "content_fallback",
	})
}

func getOpenAIResponsesAutoUpstreamSessionAnchor(c *gin.Context) openAIUpstreamSessionAnchor {
	if c == nil {
		return openAIUpstreamSessionAnchor{}
	}
	raw, ok := c.Get(openAIResponsesAutoUpstreamSessionAnchorGinContextKey)
	if !ok {
		return openAIUpstreamSessionAnchor{}
	}
	anchor, _ := raw.(openAIUpstreamSessionAnchor)
	anchor.SessionID = strings.TrimSpace(anchor.SessionID)
	anchor.Source = strings.TrimSpace(anchor.Source)
	return anchor
}

func resolveOpenAIUpstreamSessionAnchor(c *gin.Context, promptCacheKey string) openAIUpstreamSessionAnchor {
	if c != nil {
		if sessionID := strings.TrimSpace(c.GetHeader("session_id")); sessionID != "" {
			return openAIUpstreamSessionAnchor{SessionID: sessionID, Source: "header_session_id"}
		}
		if conversationID := strings.TrimSpace(c.GetHeader("conversation_id")); conversationID != "" {
			return openAIUpstreamSessionAnchor{SessionID: conversationID, Source: "header_conversation_id"}
		}
	}
	if trimmedPromptCacheKey := strings.TrimSpace(promptCacheKey); trimmedPromptCacheKey != "" {
		return openAIUpstreamSessionAnchor{SessionID: trimmedPromptCacheKey, Source: "prompt_cache_key"}
	}
	return getOpenAIResponsesAutoUpstreamSessionAnchor(c)
}

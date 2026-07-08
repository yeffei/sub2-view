package middleware

import (
	"crypto/subtle"
	"net/http"
	"strings"

	"github.com/Wei-Shaw/sub2api/internal/config"

	"github.com/gin-gonic/gin"
)

const botServiceTokenHeader = "X-Bot-Service-Token"

func NewBotServiceAuthMiddleware(cfg *config.Config) BotServiceAuthMiddleware {
	return BotServiceAuthMiddleware(func(c *gin.Context) {
		expected := strings.TrimSpace(cfg.Bot.ServiceToken)
		if expected == "" {
			AbortWithError(c, http.StatusServiceUnavailable, "BOT_SERVICE_TOKEN_NOT_CONFIGURED", "Bot service token is not configured")
			return
		}
		token := extractBotServiceToken(c)
		if token == "" || subtle.ConstantTimeCompare([]byte(token), []byte(expected)) != 1 {
			AbortWithError(c, http.StatusUnauthorized, "INVALID_BOT_SERVICE_TOKEN", "Invalid bot service token")
			return
		}
		c.Set("auth_method", "bot_service_token")
		c.Next()
	})
}

func extractBotServiceToken(c *gin.Context) string {
	if token := strings.TrimSpace(c.GetHeader(botServiceTokenHeader)); token != "" {
		return token
	}
	authHeader := strings.TrimSpace(c.GetHeader("Authorization"))
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) == 2 && strings.EqualFold(parts[0], "Bearer") {
		return strings.TrimSpace(parts[1])
	}
	return ""
}

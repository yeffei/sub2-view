package routes

import (
	"github.com/Wei-Shaw/sub2api/internal/handler"
	"github.com/Wei-Shaw/sub2api/internal/server/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterBotRoutes(
	v1 *gin.RouterGroup,
	h *handler.Handlers,
	botAuth middleware.BotServiceAuthMiddleware,
) {
	bot := v1.Group("/bot")
	bot.Use(gin.HandlerFunc(botAuth))
	{
		bot.GET("/lines/health", h.Bot.ListLineHealth)
		bot.POST("/bindings/resolve-user", h.Bot.ResolveUser)

		users := bot.Group("/users/:id")
		{
			users.GET("", h.Bot.GetUser)
			users.GET("/usage", h.Bot.GetUserUsage)
			users.GET("/platform-quotas", h.Bot.GetUserPlatformQuotas)
			users.GET("/api-keys", h.Bot.GetUserAPIKeys)
			users.GET("/orders", h.Bot.GetUserOrders)
		}
	}
}

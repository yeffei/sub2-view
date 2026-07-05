package handler

import (
	"net/http"

	"github.com/Wei-Shaw/sub2api/internal/pkg/response"
	"github.com/Wei-Shaw/sub2api/internal/server/middleware"
	"github.com/Wei-Shaw/sub2api/internal/service"

	"github.com/gin-gonic/gin"
)

// AccountMetaHandler exposes authenticated API-key metadata for compatible
// upstream gateways. It is intentionally narrow: downstream admin sync only
// needs to know how this key is billed by the upstream.
type AccountMetaHandler struct {
	apiKeyService *service.APIKeyService
}

func NewAccountMetaHandler(apiKeyService *service.APIKeyService) *AccountMetaHandler {
	return &AccountMetaHandler{apiKeyService: apiKeyService}
}

type accountMetaResponse struct {
	Compatible          bool    `json:"compatible"`
	Platform            string  `json:"platform"`
	GroupID             *int64  `json:"group_id,omitempty"`
	GroupName           string  `json:"group_name,omitempty"`
	RateMultiplier      float64 `json:"rate_multiplier"`
	RateSource          string  `json:"rate_source"`
	SubscriptionType    string  `json:"subscription_type,omitempty"`
	ImageRateMultiplier float64 `json:"image_rate_multiplier,omitempty"`
}

// Get returns the current API key's effective upstream billing multiplier.
// GET /v1/account/meta
func (h *AccountMetaHandler) Get(c *gin.Context) {
	apiKey, ok := middleware.GetAPIKeyFromContext(c)
	if !ok || apiKey == nil {
		response.Unauthorized(c, "API key is required")
		return
	}
	if apiKey.GroupID == nil || apiKey.Group == nil {
		response.Error(c, http.StatusBadRequest, "API key is not bound to a group")
		return
	}

	rate := apiKey.Group.RateMultiplier
	rateSource := "group"
	if h != nil && h.apiKeyService != nil {
		rates, err := h.apiKeyService.GetUserGroupRates(c.Request.Context(), apiKey.UserID)
		if err != nil {
			response.Error(c, http.StatusInternalServerError, "Failed to resolve API key rate")
			return
		}
		if rates != nil {
			if override, exists := rates[*apiKey.GroupID]; exists {
				rate = override
				rateSource = "user_group_override"
			}
		}
	}
	if rate < 0 {
		rate = 1
		rateSource = "fallback"
	}

	payload := accountMetaResponse{
		Compatible:          true,
		Platform:            apiKey.Group.Platform,
		GroupID:             apiKey.GroupID,
		GroupName:           apiKey.Group.Name,
		RateMultiplier:      rate,
		RateSource:          rateSource,
		SubscriptionType:    apiKey.Group.SubscriptionType,
		ImageRateMultiplier: apiKey.Group.ImageRateMultiplier,
	}
	response.Success(c, payload)
}

package service

import (
	"backend/internal/middlewares"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type IncomingFriendsRequestsCountResponse struct {
	Count int64 `json:"count"`
}

// GetUserIncomingFriendsRequestsCount godoc
// @Summary returns user's incoming friends requests count
// @Tags Friends requests
// @Router /api/user/friends/requests/incoming/count [get]
// @Security ApiKeyAuth
// @Produce json
// @Success 200 {object} IncomingFriendsRequestsCountResponse
func (s *Service) GetUserIncomingFriendsRequestsCount(c *gin.Context) {
	authData := middlewares.GetInitDataFromContext(c)
	if authData == nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	count, err := s.db.GetIncomingFriendsRequestsCount(c, authData.User.ID)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to get incoming friend requests count: %w", err))
		return
	}
	c.JSON(http.StatusOK, IncomingFriendsRequestsCountResponse{count})
}

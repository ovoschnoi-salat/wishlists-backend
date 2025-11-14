package service

import (
	"backend/internal/middlewares"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetUserIncomingFriendsRequests godoc
// @Summary returns user's incoming friends requests
// @Tags Friends requests
// @Router /api/user/friends/requests/incoming [get]
// @Security ApiKeyAuth
// @Produce json
// @Success 200 {array} Friend
func (s *Service) GetUserIncomingFriendsRequests(c *gin.Context) {
	authData := middlewares.GetInitDataFromContext(c)
	if authData == nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	requests, err := s.db.GetIncomingFriendsRequests(c, authData.User.ID)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to get incoming friend requests: %w", err))
		return
	}
	c.JSON(http.StatusOK, mapStoreUsersToFriends(requests))
}

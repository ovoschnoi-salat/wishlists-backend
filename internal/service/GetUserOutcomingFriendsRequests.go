package service

import (
	"backend/internal/middlewares"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetUserOutcomingFriendsRequests godoc
// @Summary returns user's incoming friends requests
// @Tags Friends requests
// @Router /api/user/friends/requests/outcoming [get]
// @Security ApiKeyAuth
// @Produce json
// @Success 200 {array} Friend
func (s *Service) GetUserOutcomingFriendsRequests(c *gin.Context) {
	authData := middlewares.GetInitDataFromContext(c)
	if authData == nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	requests, err := s.db.GetOutcomingFriendsRequests(c, authData.User.ID)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("cannot get outcoming friend requests: %w", err))
		return
	}
	c.JSON(http.StatusOK, mapStoreUsersToFriends(requests))
}

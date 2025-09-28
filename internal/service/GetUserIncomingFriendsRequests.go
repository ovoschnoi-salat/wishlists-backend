package service

import (
	"backend/internal/middlewares"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetUserIncomingFriendsRequests godoc
// @Summary returns user's incoming friends requests
// @Tags Friends
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
		c.Error(err)
		c.Status(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, mapStoreUsersToFriends(requests))
}

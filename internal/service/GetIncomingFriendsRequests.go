package service

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetIncomingFriendsRequests godoc
// @Summary returns user's incoming friends requests
// @Tags friends
// @Accept json
// @Produce json
// @Success 200 {array} Friend
// @Router /user/friends/requests/incoming [get]
func (s *Service) GetIncomingFriendsRequests(c *gin.Context) {
	// TODO: Extract user ID from authentication context
	// For now, using hardcoded user ID like in GetMyWishlists
	userID := int64(1)

	requests, err := s.db.GetIncomingFriendsRequests(c, userID)
	if err != nil {
		c.Error(err)
		c.Status(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, mapStoreUsersToFriends(requests))
}

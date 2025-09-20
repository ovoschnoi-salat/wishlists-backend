package service

import (
	"backend/internal/store"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// AcceptIncomingFriendsRequest godoc
// @Summary accepts an incoming friend request
// @Tags friends
// @Accept json
// @Produce json
// @Param friend_id path int true "Friend ID"
// @Success 200 {object} map[string]string
// @Router /user/friends/requests/{friend_id}/accept [post]
func (s *Service) AcceptIncomingFriendsRequest(c *gin.Context) {
	// TODO: Extract user ID from authentication context
	// For now, using hardcoded user ID like in GetMyWishlists
	userIDTo := int64(1)

	// Get friend ID from URL parameter
	friendIDStr := c.Param("friend_id")
	friendID, err := strconv.ParseInt(friendIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid friend ID"})
		return
	}

	// Accept the friend request
	err = s.db.AcceptFriendsRequest(c, store.AcceptFriendsRequestParams{
		UserIDFrom: friendID,
		UserIDTo:   userIDTo,
	})
	if err != nil {
		c.Error(err)
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Friend request accepted successfully"})
}

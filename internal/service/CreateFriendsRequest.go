package service

import (
	"backend/internal/store"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateFriendsRequest godoc
// @Summary creates a friend request to another user
// @Tags friends
// @Accept json
// @Produce json
// @Param friend_id path int true "Friend ID"
// @Success 200 {object} map[string]string
// @Router /user/friends/requests/{friend_id} [post]
func (s *Service) CreateFriendsRequest(c *gin.Context) {
	// TODO: Extract user ID from authentication context
	// For now, using hardcoded user ID like in GetMyWishlists
	userIDFrom := int64(1)

	// Get friend ID from URL parameter
	friendIDStr := c.Param("friend_id")
	friendID, err := strconv.ParseInt(friendIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid friend ID"})
		return
	}

	// Check if trying to send request to self
	if userIDFrom == friendID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot send friend request to yourself"})
		return
	}

	// Create the friend request
	err = s.db.CreateFriendsRequest(c, store.CreateFriendsRequestParams{
		UserIDFrom: userIDFrom,
		UserIDTo:   friendID,
	})
	if err != nil {
		c.Error(err)
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Friend request sent successfully"})
}

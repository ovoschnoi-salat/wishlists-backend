package service

import (
	"backend/internal/middlewares"
	"backend/internal/store"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// AcceptUserIncomingFriendsRequest godoc
// @Summary accepts an incoming friend request
// @Tags Friends
// @Router /api/user/friend/request/accept [post]
// @Security ApiKeyAuth
// @Param friend_id query int true "Friend ID"
// @Produce json
// @Success 200 {object} map[string]string
func (s *Service) AcceptUserIncomingFriendsRequest(c *gin.Context) {
	authData := middlewares.GetInitDataFromContext(c)
	if authData == nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

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
		UserIDTo:   authData.User.ID,
	})
	if err != nil {
		c.Error(err)
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Friend request accepted successfully"})
}

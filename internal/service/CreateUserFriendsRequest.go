package service

import (
	"backend/internal/middlewares"
	"backend/internal/store"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateUserFriendsRequest godoc
// @Summary creates a friend request to another user
// @Tags Friends
// @Router /api/user/friend/request/new [post]
// @Security ApiKeyAuth
// @Param username query string true "Friend username"
// @Produce json
// @Success 200 {object} map[string]string
func (s *Service) CreateUserFriendsRequest(c *gin.Context) {
	authData := middlewares.GetInitDataFromContext(c)
	if authData == nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	friendUsernameStr := c.Query("username")
	friendUsername, err := strconv.ParseInt(friendUsernameStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid friend ID"})
		return
	}

	if authData.User.ID == friendUsername {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot send friend request to yourself"})
		return
	}

	err = s.db.CreateFriendsRequest(c, store.CreateFriendsRequestParams{
		UserIDFrom: authData.User.ID,
		UserIDTo:   friendUsername,
	})
	if err != nil {
		c.Error(err)
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Friend request sent successfully"})
}

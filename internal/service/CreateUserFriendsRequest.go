package service

import (
	"backend/internal/middlewares"
	"backend/internal/store"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
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

	if authData.User.Username == friendUsernameStr {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot send friend request to yourself"})
		return
	}

	friend, err := s.db.GetUserByUsername(c, friendUsernameStr)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.Error(err)
		c.Status(http.StatusInternalServerError)
		return
	}

	count, err := s.db.CreateFriendsRequest(c, store.CreateFriendsRequestParams{
		UserIDFrom: authData.User.ID,
		UserIDTo:   friend.ID,
	})
	if err != nil {
		c.Error(err)
		c.Status(http.StatusInternalServerError)
		return
	}
	if count < 1 {
		c.Error(errors.New("no rows updated"))
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Friend request sent successfully"})
}

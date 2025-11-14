package service

import (
	"backend/internal/middlewares"
	"backend/internal/store"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

// CreateUserFriendsRequest godoc
// @Summary creates a friend request to another user
// @Tags Friends requests
// @Router /api/user/friend/request [post]
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
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("cannot send friend request to yourself"))
		return
	}

	friend, err := s.db.GetUserByUsername(c, friendUsernameStr)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.AbortWithError(http.StatusNotFound, fmt.Errorf("user not found: %w", err))
			return
		}
		c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("cannot get user by username: %w", err))
		return
	}

	count, err := s.db.CreateFriendsRequest(c, store.CreateFriendsRequestParams{
		UserIDFrom: authData.User.ID,
		UserIDTo:   friend.ID,
	})
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("error creating friendship: %w", err))
		return
	}
	if count < 1 {
		c.AbortWithError(http.StatusInternalServerError, errors.New("error creating friendship: no rows updated"))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Friend request sent successfully"})
}

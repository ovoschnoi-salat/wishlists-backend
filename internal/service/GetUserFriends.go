package service

import (
	"backend/internal/errorResponse"
	"backend/internal/errorResponse/codes"
	"backend/internal/middlewares"
	"backend/internal/store"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Friend struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Name     string `json:"name"`
	PhotoUrl string `json:"photo_url"`
}

// GetFriends godoc
// @Summary returns user's friends list
// @Tags Friends
// @Router /api/user/friends [get]
// @Security ApiKeyAuth
// @Produce json
// @Failure 401 {object} Response
// @Failure 500 {object} Response
// @Success 200 {array} Friend
func (s *Service) GetFriends(c *gin.Context) {
	authData, authorized := middlewares.GetInitDataFromContext(c)
	if !authorized {
		errorResponse.Send(c, http.StatusUnauthorized, codes.UnauthorizedErrCode, nil)
		return
	}

	friends, err := s.db.GetFriends(c, authData.User.ID)
	if err != nil {
		errorResponse.Send(c, http.StatusInternalServerError, codes.InternalErrCode, fmt.Errorf("error getting friends: %w", err))
		return
	}

	c.JSON(http.StatusOK, mapStoreUsersToFriends(friends))
}

func mapStoreUserToFriend(user store.User) Friend {
	return Friend{
		ID:       user.ID,
		Username: user.Username,
		Name:     user.DisplayedName,
		PhotoUrl: user.PhotoUrl,
	}
}

func mapStoreUsersToFriends(users []store.User) []Friend {
	res := make([]Friend, len(users))
	for i, u := range users {
		res[i] = mapStoreUserToFriend(u)
	}
	return res
}

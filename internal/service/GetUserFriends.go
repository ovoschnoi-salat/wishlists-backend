package service

import (
	"backend/internal/middlewares"
	"backend/internal/store"
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
// @Tags friends
// @Accept json
// @Produce json
// @Success 200 {array} Friend
// @Router /user/friends [get]
// @Security ApiKeyAuth
func (s *Service) GetFriends(c *gin.Context) {
	authData := middlewares.GetInitDataFromContext(c)
	if authData == nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	friends, err := s.db.GetFriends(c, authData.User.ID)
	if err != nil {
		c.Error(err)
		c.Status(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, mapStoreUsersToFriends(friends))
}

func mapStoreUserToFriend(user store.User) Friend {
	return Friend{
		ID:       user.ID,
		Username: user.Username,
		Name:     user.Name,
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

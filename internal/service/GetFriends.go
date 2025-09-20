package service

import (
	"backend/internal/store"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Friend struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	PhotoUrl string `json:"photo_url"`
}

// GetFriends godoc
// @Summary returns user's friends list
// @Tags friends
// @Accept json
// @Produce json
// @Success 200 {array} Friend
// @Router /user/friends [get]
func (s *Service) GetFriends(c *gin.Context) {
	// TODO: Extract user ID from authentication context
	// For now, using hardcoded user ID like in GetMyWishlists
	userID := int64(1)

	friends, err := s.db.GetFriends(c, userID)
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
		Username: user.Username.String,
		PhotoUrl: user.PhotoUrl.String,
	}
}

func mapStoreUsersToFriends(users []store.User) []Friend {
	res := make([]Friend, len(users))
	for i, u := range users {
		res[i] = mapStoreUserToFriend(u)
	}
	return res
}

package service

import (
	"backend/internal/middlewares"
	"backend/internal/store"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// AcceptUserIncomingFriendsRequest godoc
// @Summary accepts an incoming friend request
// @Tags Friends requests
// @Router /api/user/friend/request/accept [post]
// @Security ApiKeyAuth
// @Param friend_id query int true "Friend ID"
// @Success 204
func (s *Service) AcceptUserIncomingFriendsRequest(c *gin.Context) {
	authData := middlewares.GetInitDataFromContext(c)
	if authData == nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	friendIDStr := c.Query("friend_id")
	friendID, err := strconv.ParseInt(friendIDStr, 10, 64)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("Invalid friend ID: %w", err))
		return
	}

	count, err := s.db.AcceptFriendsRequest(c, store.AcceptFriendsRequestParams{
		UserIDFrom: friendID,
		UserIDTo:   authData.User.ID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	if count < 1 {
		c.AbortWithError(http.StatusInternalServerError, errors.New("no rows updated"))
		return
	}

	c.Status(http.StatusNoContent)
}

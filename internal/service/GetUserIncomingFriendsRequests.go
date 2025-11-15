package service

import (
	"backend/internal/errors"
	"backend/internal/errors/codes"
	"backend/internal/middlewares"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetUserIncomingFriendsRequests godoc
// @Summary returns user's incoming friends requests
// @Tags Friends requests
// @Router /api/user/friends/requests/incoming [get]
// @Security ApiKeyAuth
// @Produce json
// @Failure 401 {object} errors.Response
// @Failure 500 {object} errors.Response
// @Success 200 {array} Friend
func (s *Service) GetUserIncomingFriendsRequests(c *gin.Context) {
	authData, authorized := middlewares.GetInitDataFromContext(c)
	if !authorized {
		errors.SendResponse(c, http.StatusInternalServerError, codes.InternalErrCode, noInitDataErr)
		return
	}

	requests, err := s.db.GetIncomingFriendsRequests(c, authData.User.ID)
	if err != nil {
		errors.SendResponse(c, http.StatusInternalServerError, codes.InternalErrCode, fmt.Errorf("failed to get incoming friend requests: %w", err))
		return
	}
	c.JSON(http.StatusOK, mapStoreUsersToFriends(requests))
}

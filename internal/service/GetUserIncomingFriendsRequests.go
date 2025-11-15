package service

import (
	"backend/internal/errorResponse"
	"backend/internal/errorResponse/codes"
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
// @Failure 401 {object} errorResponse.Response
// @Failure 500 {object} errorResponse.Response
// @Success 200 {array} Friend
func (s *Service) GetUserIncomingFriendsRequests(c *gin.Context) {
	authData, authorized := middlewares.GetInitDataFromContext(c)
	if !authorized {
		errorResponse.Send(c, http.StatusUnauthorized, codes.UnauthorizedErrCode, nil)
		return
	}

	requests, err := s.db.GetIncomingFriendsRequests(c, authData.User.ID)
	if err != nil {
		errorResponse.Send(c, http.StatusInternalServerError, codes.InternalErrCode, fmt.Errorf("failed to get incoming friend requests: %w", err))
		return
	}
	c.JSON(http.StatusOK, mapStoreUsersToFriends(requests))
}

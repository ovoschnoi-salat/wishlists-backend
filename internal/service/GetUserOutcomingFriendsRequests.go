package service

import (
	"backend/internal/errors"
	"backend/internal/errors/codes"
	"backend/internal/middlewares"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetUserOutcomingFriendsRequests godoc
// @Summary returns user's incoming friends requests
// @Tags Friends requests
// @Router /api/user/friends/requests/outcoming [get]
// @Security ApiKeyAuth
// @Produce json
// @Failure 401 {object} errors.Response
// @Failure 500 {object} errors.Response
// @Success 200 {array} Friend
func (s *Service) GetUserOutcomingFriendsRequests(c *gin.Context) {
	authData, authorized := middlewares.GetInitDataFromContext(c)
	if !authorized {
		errors.SendResponse(c, http.StatusInternalServerError, codes.InternalErrCode, noInitDataErr)
		return
	}

	requests, err := s.db.GetOutcomingFriendsRequests(c, authData.User.ID)
	if err != nil {
		errors.SendResponse(c, http.StatusInternalServerError, codes.InternalErrCode, fmt.Errorf("cannot get outcoming friend requests: %w", err))
		return
	}
	c.JSON(http.StatusOK, mapStoreUsersToFriends(requests))
}

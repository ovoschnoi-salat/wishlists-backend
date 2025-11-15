package service

import (
	"backend/internal/errorResponse"
	"backend/internal/errorResponse/codes"
	"backend/internal/middlewares"
	"backend/internal/store"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// DenyUserIncomingFriendsRequest godoc
// @Summary accepts an incoming friend request
// @Tags Friends requests
// @Router /api/user/friend/request/deny [post]
// @Security ApiKeyAuth
// @Param friend_id query int true "Friend ID"
// @Failure 400 {object} errorResponse.Response
// @Failure 401 {object} errorResponse.Response
// @Failure 500 {object} errorResponse.Response
// @Success 204
func (s *Service) DenyUserIncomingFriendsRequest(c *gin.Context) {
	authData, authorized := middlewares.GetInitDataFromContext(c)
	if !authorized {
		errorResponse.Send(c, http.StatusUnauthorized, codes.UnauthorizedErrCode, nil)
		return
	}

	friendIDStr := c.Query("friend_id")
	friendID, err := strconv.ParseInt(friendIDStr, 10, 64)
	if err != nil {
		errorResponse.Send(c, http.StatusBadRequest, codes.InvalidRequestParametersErrCode, fmt.Errorf("invalid friend_id: %w", err))
		return
	}

	count, err := s.db.DenyFriendsRequest(c, store.DenyFriendsRequestParams{
		UserIDFrom: friendID,
		UserIDTo:   authData.User.ID,
	})
	if err != nil {
		errorResponse.Send(c, http.StatusInternalServerError, codes.InternalErrCode, fmt.Errorf("cannot deny friend request: %w", err))
		return
	}
	if count == 0 {
		errorResponse.Send(c, http.StatusInternalServerError, codes.InternalErrCode, errors.New("cannot deny friend request: no rows updated"))
		return
	}

	c.Status(http.StatusNoContent)
}

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

// AcceptUserIncomingFriendsRequest godoc
// @Summary accepts an incoming friend request
// @Tags Friends requests
// @Router /api/user/friend/request/accept [post]
// @Security ApiKeyAuth
// @Param friend_id query int true "Friend ID"
// @Failure 400 {object} errorResponse.Response
// @Failure 401 {object} errorResponse.Response
// @Failure 500 {object} errorResponse.Response
// @Success 204
func (s *Service) AcceptUserIncomingFriendsRequest(c *gin.Context) {
	authData, authorized := middlewares.GetInitDataFromContext(c)
	if !authorized {
		errorResponse.Send(c, http.StatusUnauthorized, codes.UnauthorizedErrCode, nil)
		return
	}

	friendIDStr := c.Query("friend_id")
	friendID, err := strconv.ParseInt(friendIDStr, 10, 64)
	if err != nil {
		errorResponse.Send(c, http.StatusBadRequest, codes.InvalidRequestParametersErrCode, fmt.Errorf("invalid friend ID: %w", err))
		return
	}

	count, err := s.db.AcceptFriendsRequest(c, store.AcceptFriendsRequestParams{
		UserIDFrom: friendID,
		UserIDTo:   authData.User.ID,
	})
	if err != nil {
		errorResponse.Send(c, http.StatusInternalServerError, codes.InternalErrCode, fmt.Errorf("failed to accept friends request: %w", err))
		return
	}
	if count == 0 {
		errorResponse.Send(c, http.StatusBadRequest, codes.InvalidRequestErrCode, errors.New("failed to accept friends request: no rows affected"))
		return
	}

	c.Status(http.StatusNoContent)
}

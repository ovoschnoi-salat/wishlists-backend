package service

import (
	"backend/internal/middlewares"
	"backend/internal/store"
	"backend/internal/subcodeErrors"
	"backend/internal/subcodeErrors/codes"
	"errors"
	"fmt"
	"net/http"
	"strings"

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
// @Failure 400 {object} subcodeErrors.Response
// @Failure 401 {object} subcodeErrors.Response
// @Failure 500 {object} subcodeErrors.Response
// @Success 204
func (s *Service) CreateUserFriendsRequest(c *gin.Context) {
	authData, authorized := middlewares.GetInitDataFromContext(c)
	if !authorized {
		subcodeErrors.SendResponse(c, http.StatusInternalServerError, codes.InternalErrCode, noInitDataErr)
		return
	}

	friendUsernameStr := c.Query("username")

	if authData.User.Username == friendUsernameStr {
		subcodeErrors.SendResponse(c, http.StatusBadRequest, codes.CantSendRequestToYourselfErrCode, nil)
		return
	}

	friendUsernameStr = strings.TrimPrefix(friendUsernameStr, "@")
	friendUsernameStr = strings.ToLower(friendUsernameStr)

	friend, err := s.db.GetUserByUsername(c, friendUsernameStr)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			subcodeErrors.SendResponse(c, http.StatusNotFound, codes.FriendNotFoundErrCode, nil)
			return
		}
		subcodeErrors.SendResponse(c, http.StatusInternalServerError, codes.InternalErrCode, fmt.Errorf("cannot get user by username: %w", err))
		return
	}

	if !friend.OpenToRequests {
		subcodeErrors.SendResponse(c, http.StatusNotFound, codes.FriendNotFoundErrCode, nil)
		return
	}

	_, err = s.db.CreateFriendsRequest(c, store.CreateFriendsRequestParams{
		UserIDFrom: authData.User.ID,
		UserIDTo:   friend.ID,
	})
	if err != nil {
		subcodeErrors.SendResponse(c, http.StatusInternalServerError, codes.InternalErrCode, fmt.Errorf("error creating friendship: %w", err))
		return
	}

	c.Status(http.StatusNoContent)
}

package service

import (
	"backend/internal/errorResponse"
	"backend/internal/errorResponse/codes"
	"backend/internal/middlewares"
	"backend/internal/store"
	"errors"
	"fmt"
	"net/http"

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
// @Failure 400 {object} errorResponse.Response
// @Failure 401 {object} errorResponse.Response
// @Failure 500 {object} errorResponse.Response
// @Success 204
func (s *Service) CreateUserFriendsRequest(c *gin.Context) {
	authData, authorized := middlewares.GetInitDataFromContext(c)
	if !authorized {
		errorResponse.Send(c, http.StatusInternalServerError, codes.InternalErrCode, noInitDataErr)
		return
	}

	friendUsernameStr := c.Query("username")

	if authData.User.Username == friendUsernameStr {
		errorResponse.Send(c, http.StatusBadRequest, codes.CantSendRequestToYourselfErrCode, nil)
		return
	}

	friend, err := s.db.GetUserByUsername(c, friendUsernameStr)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			errorResponse.Send(c, http.StatusNotFound, codes.FriendNotFoundErrCode, nil)
			return
		}
		errorResponse.Send(c, http.StatusInternalServerError, codes.InternalErrCode, fmt.Errorf("cannot get user by username: %w", err))
		return
	}

	if !friend.OpenToRequests {
		errorResponse.Send(c, http.StatusNotFound, codes.FriendNotFoundErrCode, nil)
		return
	}

	count, err := s.db.CreateFriendsRequest(c, store.CreateFriendsRequestParams{
		UserIDFrom: authData.User.ID,
		UserIDTo:   friend.ID,
	})
	if err != nil {
		errorResponse.Send(c, http.StatusInternalServerError, codes.InternalErrCode, fmt.Errorf("error creating friendship: %w", err))
		return
	}
	if count == 0 {
		errorResponse.Send(c, http.StatusInternalServerError, codes.InternalErrCode, errors.New("error creating friendship: no rows updated"))
		return
	}

	c.Status(http.StatusNoContent)
}

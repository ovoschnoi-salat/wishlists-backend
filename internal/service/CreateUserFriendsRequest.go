package service

import (
	"backend/internal/errors"
	"backend/internal/errors/codes"
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
// @Failure 400 {object} errors.Response
// @Failure 401 {object} errors.Response
// @Failure 500 {object} errors.Response
// @Success 204
func (s *Service) CreateUserFriendsRequest(c *gin.Context) {
	authData, authorized := middlewares.GetInitDataFromContext(c)
	if !authorized {
		errors.SendResponse(c, http.StatusInternalServerError, codes.InternalErrCode, noInitDataErr)
		return
	}

	friendUsernameStr := c.Query("username")

	if authData.User.Username == friendUsernameStr {
		errors.SendResponse(c, http.StatusBadRequest, codes.CantSendRequestToYourselfErrCode, nil)
		return
	}

	friend, err := s.db.GetUserByUsername(c, friendUsernameStr)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			errors.SendResponse(c, http.StatusNotFound, codes.FriendNotFoundErrCode, nil)
			return
		}
		errors.SendResponse(c, http.StatusInternalServerError, codes.InternalErrCode, fmt.Errorf("cannot get user by username: %w", err))
		return
	}

	if !friend.OpenToRequests {
		errors.SendResponse(c, http.StatusNotFound, codes.FriendNotFoundErrCode, nil)
		return
	}

	count, err := s.db.CreateFriendsRequest(c, store.CreateFriendsRequestParams{
		UserIDFrom: authData.User.ID,
		UserIDTo:   friend.ID,
	})
	if err != nil {
		errors.SendResponse(c, http.StatusInternalServerError, codes.InternalErrCode, fmt.Errorf("error creating friendship: %w", err))
		return
	}
	if count == 0 {
		errors.SendResponse(c, http.StatusInternalServerError, codes.InternalErrCode, errors.New("error creating friendship: no rows updated"))
		return
	}

	c.Status(http.StatusNoContent)
}

package service

import (
	"backend/internal/errors"
	"backend/internal/errors/codes"
	"backend/internal/middlewares"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type IncomingFriendsRequestsCountResponse struct {
	Count int64 `json:"count"`
}

// GetUserIncomingFriendsRequestsCount godoc
// @Summary returns user's incoming friends requests count
// @Tags Friends requests
// @Router /api/user/friends/requests/incoming/count [get]
// @Security ApiKeyAuth
// @Produce json
// @Failure 401 {object} errors.Response
// @Failure 500 {object} errors.Response
// @Success 200 {object} IncomingFriendsRequestsCountResponse
func (s *Service) GetUserIncomingFriendsRequestsCount(c *gin.Context) {
	authData, authorized := middlewares.GetInitDataFromContext(c)
	if !authorized {
		errors.SendResponse(c, http.StatusInternalServerError, codes.InternalErrCode, noInitDataErr)
		return
	}

	count, err := s.db.GetIncomingFriendsRequestsCount(c, authData.User.ID)
	if err != nil {
		errors.SendResponse(c, http.StatusInternalServerError, codes.InternalErrCode, fmt.Errorf("failed to get incoming friend requests count: %w", err))
		return
	}
	c.JSON(http.StatusOK, IncomingFriendsRequestsCountResponse{count})
}

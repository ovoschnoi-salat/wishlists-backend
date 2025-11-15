package service

import (
	"backend/internal/errorResponse"
	"backend/internal/errorResponse/codes"
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
// @Failure 401 {object} Response
// @Failure 500 {object} Response
// @Success 200 {object} IncomingFriendsRequestsCountResponse
func (s *Service) GetUserIncomingFriendsRequestsCount(c *gin.Context) {
	authData, authorized := middlewares.GetInitDataFromContext(c)
	if !authorized {
		errorResponse.Send(c, http.StatusUnauthorized, codes.UnauthorizedErrCode, nil)
		return
	}

	count, err := s.db.GetIncomingFriendsRequestsCount(c, authData.User.ID)
	if err != nil {
		errorResponse.Send(c, http.StatusInternalServerError, codes.InternalErrCode, fmt.Errorf("failed to get incoming friend requests count: %w", err))
		return
	}
	c.JSON(http.StatusOK, IncomingFriendsRequestsCountResponse{count})
}

package service

import (
	"backend/internal/middlewares"
	"backend/internal/store"
	"backend/internal/subcodeErrors"
	"backend/internal/subcodeErrors/codes"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// UpdateUserSettings godoc
// @Summary updates user
// @Tags User
// @Router /api/user/settings [patch]
// @Security ApiKeyAuth
// @Accept	json
// @Param	wishlist	body	UserSettings	true	"request body"
// @Produce	json
// @Failure 400 {object} subcodeErrors.Response
// @Failure 401 {object} subcodeErrors.Response
// @Failure 500 {object} subcodeErrors.Response
// @Success 204
func (s *Service) UpdateUserSettings(c *gin.Context) {
	authData, authorized := middlewares.GetInitDataFromContext(c)
	if !authorized {
		subcodeErrors.SendResponse(c, http.StatusInternalServerError, codes.InternalErrCode, noInitDataErr)
		return
	}

	req := new(UserSettings)
	if err := c.ShouldBindJSON(req); err != nil {
		subcodeErrors.SendResponse(c, http.StatusBadRequest, codes.InvalidRequestParametersErrCode, err)
		return
	}

	count, err := s.db.UpdateUserSettings(c, store.UpdateUserSettingsParams{
		ID:             authData.User.ID,
		DisplayedName:  req.DisplayedName,
		OpenToRequests: req.OpenToRequests,
	})
	if err != nil {
		subcodeErrors.SendResponse(c, http.StatusInternalServerError, codes.InternalErrCode, fmt.Errorf("failed to update user: %w", err))
		return
	}
	if count == 0 {
		subcodeErrors.SendResponse(c, http.StatusInternalServerError, codes.InternalErrCode, errors.New("no user settings found"))
	}

	c.Status(http.StatusNoContent)
}

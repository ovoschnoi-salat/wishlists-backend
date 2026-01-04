package service

import (
	"backend/internal/middlewares"
	"backend/internal/store"
	"backend/internal/subcodeErrors"
	"backend/internal/subcodeErrors/codes"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserSettings struct {
	DisplayedName  string `json:"displayed_name"`
	OpenToRequests bool   `json:"open_to_requests"`
}

// GetUserSettings godoc
// @Summary returns user's settings
// @Tags User
// @Router /api/user/settings [get]
// @Security ApiKeyAuth
// @Produce json
// @Failure 500 {object} subcodeErrors.Response
// @Success 200 {object} UserSettings
func (s *Service) GetUserSettings(c *gin.Context) {
	user, found := middlewares.GetUserDataFromContext(c)
	if !found {
		subcodeErrors.SendResponse(c, http.StatusInternalServerError, codes.InternalErrCode, fmt.Errorf("error getting user from context"))
		return
	}

	c.JSON(http.StatusOK, mapStoreUserToUserSettings(user))
}

func mapStoreUserToUserSettings(user store.User) UserSettings {
	return UserSettings{
		DisplayedName:  user.DisplayedName,
		OpenToRequests: user.OpenToRequests,
	}
}

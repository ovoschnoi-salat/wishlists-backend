package service

import (
	"backend/internal/errors"
	"backend/internal/errors/codes"
	"backend/internal/middlewares"
	"backend/internal/store"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UpdateUser struct {
	PhotoUrl       string `json:"photo_url"`
	OpenToRequests bool   `json:"open_to_requests"`
	DisplayedName  string `json:"displayed_name"`
}

type User struct {
	DisplayedName  string `json:"displayed_name"`
	PhotoUrl       string `json:"photo_url"`
	OpenToRequests bool   `json:"open_to_requests"`
}

// UpdateUser godoc
// @Summary updates user
// @Tags User
// @Router /api/user [patch]
// @Security ApiKeyAuth
// @Accept	json
// @Param	wishlist	body	UpdateUser	true	"request body"
// @Produce	json
// @Failure 400 {object} errors.Response
// @Failure 401 {object} errors.Response
// @Failure 500 {object} errors.Response
// @Success 200 {object} User
func (s *Service) UpdateUser(c *gin.Context) {
	authData, authorized := middlewares.GetInitDataFromContext(c)
	if !authorized {
		errors.SendResponse(c, http.StatusInternalServerError, codes.InternalErrCode, noInitDataErr)
		return
	}

	req := new(UpdateUser)
	if err := c.ShouldBindJSON(req); err != nil {
		errors.SendResponse(c, http.StatusBadRequest, codes.InvalidRequestParametersErrCode, err)
		return
	}

	user, err := s.db.UpdateUser(c, store.UpdateUserParams{
		ID:             authData.User.ID,
		DisplayedName:  req.DisplayedName,
		PhotoUrl:       req.PhotoUrl,
		OpenToRequests: req.OpenToRequests,
	})
	if err != nil {
		errors.SendResponse(c, http.StatusInternalServerError, codes.InternalErrCode, fmt.Errorf("failed to update user: %w", err))
		return
	}

	c.JSON(http.StatusOK, mapStoreUserToUser(user))
}

func mapStoreUserToUser(u store.User) User {
	return User{
		DisplayedName:  u.DisplayedName,
		PhotoUrl:       u.PhotoUrl,
		OpenToRequests: u.OpenToRequests,
	}
}

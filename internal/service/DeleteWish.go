package service

import (
	"backend/internal/middlewares"
	"backend/internal/store"
	"backend/internal/subcodeErrors"
	"backend/internal/subcodeErrors/codes"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// DeleteUserWish godoc
// @Summary delete wish
// @Tags User
// @Router /api/user/wish [delete]
// @Security ApiKeyAuth
// @Accept	json
// @Param	wish_id	query	int	true	"Wish ID"
// @Produce	json
// @Failure 400 {object} subcodeErrors.Response
// @Failure 401 {object} subcodeErrors.Response
// @Failure 500 {object} subcodeErrors.Response
// @Success 204
func (s *Service) DeleteUserWish(c *gin.Context) {
	authData, authorized := middlewares.GetInitDataFromContext(c)
	if !authorized {
		subcodeErrors.SendResponse(c, http.StatusInternalServerError, codes.InternalErrCode, noInitDataErr)
		return
	}

	wishIDRaw := c.Query("wish_id")
	wishID, err := strconv.ParseInt(wishIDRaw, 10, 64)
	if err != nil {
		subcodeErrors.SendResponse(c, http.StatusBadRequest, codes.InvalidRequestParametersErrCode, fmt.Errorf("invalid wish_id: %w", err))
		return
	}

	count, err := s.db.DeleteWishlistItem(c, store.DeleteWishlistItemParams{
		ID:      wishID,
		OwnerID: authData.User.ID,
	})
	if err != nil {
		subcodeErrors.SendResponse(c, http.StatusInternalServerError, codes.InternalErrCode, fmt.Errorf("error deleting wish: %w", err))
		return
	}
	if count == 0 {
		subcodeErrors.SendResponse(c, http.StatusBadRequest, codes.InvalidRequestErrCode, errors.New("wish not found"))
		return
	}

	c.Status(http.StatusNoContent)
}

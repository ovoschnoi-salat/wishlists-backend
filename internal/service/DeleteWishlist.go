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

// DeleteWishlist godoc
// @Summary creates wishlist
// @Tags User
// @Router /api/user/wishlist [delete]
// @Security ApiKeyAuth
// @Accept	json
// @Param	wishlist_id	query	int	true	"Wishlist ID"
// @Produce	json
// @Failure 400 {object} subcodeErrors.Response
// @Failure 401 {object} subcodeErrors.Response
// @Failure 500 {object} subcodeErrors.Response
// @Success 204
func (s *Service) DeleteWishlist(c *gin.Context) {
	authData, authorized := middlewares.GetInitDataFromContext(c)
	if !authorized {
		subcodeErrors.SendResponse(c, http.StatusInternalServerError, codes.InternalErrCode, noInitDataErr)
		return
	}

	wishlistIDRaw := c.Query("wishlist_id")
	wishlistID, err := strconv.ParseInt(wishlistIDRaw, 10, 64)
	if err != nil {
		subcodeErrors.SendResponse(c, http.StatusBadRequest, codes.InvalidRequestParametersErrCode, fmt.Errorf("invalid wishlist_id: %w", err))
		return
	}

	count, err := s.db.CheckUserOwnsWishlist(c, store.CheckUserOwnsWishlistParams{
		ID:      wishlistID,
		OwnerID: authData.User.ID,
	})
	if err != nil {
		subcodeErrors.SendResponse(c, http.StatusInternalServerError, codes.InternalErrCode, fmt.Errorf("error checking access to wishlist: %w", err))
		return
	}
	if count == 0 {
		subcodeErrors.SendResponse(c, http.StatusBadRequest, codes.InvalidRequestErrCode, errors.New("no access to wishlist"))
		return
	}

	count, err = s.db.DeleteWishlist(c, wishlistID)
	if err != nil {
		subcodeErrors.SendResponse(c, http.StatusInternalServerError, codes.InternalErrCode, fmt.Errorf("error deleting wishlist: %w", err))
		return
	}
	if count == 0 {
		subcodeErrors.SendResponse(c, http.StatusBadRequest, codes.InvalidRequestErrCode, errors.New("wishlist not found"))
		return
	}

	c.Status(http.StatusNoContent)
}

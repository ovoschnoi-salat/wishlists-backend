package service

import (
	"backend/internal/errorResponse"
	"backend/internal/errorResponse/codes"
	"backend/internal/middlewares"
	"backend/internal/store"
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
// @Failure 400 {object} errorResponse.Response
// @Failure 401 {object} errorResponse.Response
// @Failure 500 {object} errorResponse.Response
// @Success 204
func (s *Service) DeleteWishlist(c *gin.Context) {
	authData, authorized := middlewares.GetInitDataFromContext(c)
	if !authorized {
		errorResponse.Send(c, http.StatusUnauthorized, codes.UnauthorizedErrCode, nil)
		return
	}

	wishlistIDRaw := c.Query("wishlist_id")
	wishlistID, err := strconv.ParseInt(wishlistIDRaw, 10, 64)
	if err != nil {
		errorResponse.Send(c, http.StatusBadRequest, codes.InvalidRequestParametersErrCode, fmt.Errorf("invalid wishlist_id: %w", err))
		return
	}

	count, err := s.db.DeleteWishlist(c, store.DeleteWishlistParams{
		ID:      wishlistID,
		OwnerID: authData.User.ID,
	})
	if err != nil {
		errorResponse.Send(c, http.StatusInternalServerError, codes.InternalErrCode, fmt.Errorf("error deleting wishlist: %w", err))
		return
	}
	if count == 0 {
		errorResponse.Send(c, http.StatusBadRequest, codes.InvalidRequestErrCode, errors.New("wishlist not found"))
		return
	}

	c.Status(http.StatusNoContent)
}

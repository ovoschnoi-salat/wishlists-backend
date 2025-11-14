package service

import (
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
// @Success 204
func (s *Service) DeleteWishlist(c *gin.Context) {
	authData := middlewares.GetInitDataFromContext(c)
	if authData == nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	wishlistIDRaw := c.Query("wishlist_id")
	wishlistID, err := strconv.ParseInt(wishlistIDRaw, 10, 64)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("invalid wishlist_id: %w", err))
	}

	count, err := s.db.DeleteWishlist(c, store.DeleteWishlistParams{
		ID:      wishlistID,
		OwnerID: authData.User.ID,
	})
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("error deleting wishlist: %w", err))
		return
	}
	if count == 0 {
		c.AbortWithError(http.StatusBadRequest, errors.New("wishlist not found"))
		return
	}

	c.Status(http.StatusNoContent)
}

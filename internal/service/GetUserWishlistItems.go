package service

import (
	"backend/internal/middlewares"
	"backend/internal/store"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetUserWishlistItems godoc
// @Summary returns wishlist items
// @Tags User
// @Router /api/user/wishlist/items [get]
// @Security ApiKeyAuth
// @Param wishlist_id query int true "Wishlist ID"
// @Produce json
// @Success 200 {array} WishlistItem
func (s *Service) GetUserWishlistItems(c *gin.Context) {
	authData := middlewares.GetInitDataFromContext(c)
	if authData == nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	wishlistIDStr := c.Query("wishlist_id")
	wishlistID, err := strconv.ParseInt(wishlistIDStr, 10, 64)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("error parsing wishlist_id: %w", err))
		return
	}

	items, err := s.db.GetWishlistItems(c, store.GetWishlistItemsParams{
		WishlistID: wishlistID,
		OwnerID:    authData.User.ID,
	})
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("error getting wishlist items: %w", err))
		return
	}

	c.JSON(http.StatusOK, mapStoreWishlistItemsToWishlistItems(items))
}

func mapStoreWishlistItemsToWishlistItems(items []store.WishlistItem) []WishlistItem {
	res := make([]WishlistItem, len(items))
	for i, item := range items {
		res[i] = mapStoreWishlistItemToWishlistItem(item)
	}
	return res
}

package service

import (
	"backend/internal/middlewares"
	"backend/internal/store"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetUserWishlistItems godoc
// @Summary returns wishlist items
// @Tags wishlist
// @Param wishlist_id query int true "Wishlist ID"
// @Accept json
// @Produce json
// @Success 200 {array} WishlistItem
// @Router /user/wishlist/items [get]
// @Security ApiKeyAuth
func (s *Service) GetUserWishlistItems(c *gin.Context) {
	authData := middlewares.GetInitDataFromContext(c)
	if authData == nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// Get wishlist ID from URL parameter
	wishlistIDStr := c.Query("wishlist_id")
	wishlistID, err := strconv.ParseInt(wishlistIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid wishlist ID"})
		return
	}

	// Get wishlist items
	items, err := s.db.GetWishlistItems(c, store.GetWishlistItemsParams{
		WishlistID: wishlistID,
		OwnerID:    authData.User.ID,
	})
	if err != nil {
		c.Error(err)
		c.Status(http.StatusInternalServerError)
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

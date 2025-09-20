package service

import (
	"backend/internal/store"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetWishlistItems godoc
// @Summary returns wishlist items
// @Tags wishlist
// @Param wishlist_id path int true "Wishlist ID"
// @Accept json
// @Produce json
// @Success 200 {array} WishlistItem
// @Router /user/wishlists/{wishlist_id}/items [get]
func (s *Service) GetWishlistItems(c *gin.Context) {
	// TODO: Extract user ID from authentication context
	// For now, using hardcoded user ID like in GetMyWishlists
	ownerID := int64(1)

	// Get wishlist ID from URL parameter
	wishlistIDStr := c.Param("wishlist_id")
	wishlistID, err := strconv.ParseInt(wishlistIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid wishlist ID"})
		return
	}

	// Get wishlist items
	items, err := s.db.GetWishlistItems(c, store.GetWishlistItemsParams{
		WishlistID: wishlistID,
		OwnerID:    ownerID,
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

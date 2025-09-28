package service

import (
	"backend/internal/middlewares"
	"backend/internal/store"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetUserWishlistItem godoc
// @Summary returns wishlist item
// @Tags User
// @Router /api/user/wishlist/item [get]
// @Security ApiKeyAuth
// @Param item_id query int true "Wishlist item ID"
// @Produce json
// @Success 200 {object} WishlistItem
func (s *Service) GetUserWishlistItem(c *gin.Context) {
	authData := middlewares.GetInitDataFromContext(c)
	if authData == nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// Get wishlist ID from URL parameter
	itemIDStr := c.Query("item_id")
	itemID, err := strconv.ParseInt(itemIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid wishlist ID"})
		return
	}

	// Get wishlist items
	item, err := s.db.GetWishlistItem(c, store.GetWishlistItemParams{
		ID:      itemID,
		OwnerID: authData.User.ID,
	})
	if err != nil {
		c.Error(err)
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, mapStoreWishlistItemToWishlistItem(item))
}

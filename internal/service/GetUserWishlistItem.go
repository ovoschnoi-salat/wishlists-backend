package service

import (
	"backend/internal/middlewares"
	"fmt"
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
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("invalid wishlist ID: %w", err))
		return
	}

	// Get wishlist items
	item, err := s.db.GetWishlistItem(c, itemID)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to get item: %w", err))
		return
	}
	if item.OwnerID != authData.User.ID {
		c.AbortWithError(http.StatusNotFound, fmt.Errorf("requested item does not belong to user"))
		return
	}

	c.JSON(http.StatusOK, mapStoreWishlistItemToWishlistItem(item))
}

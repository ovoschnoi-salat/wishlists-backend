package service

import (
	"backend/internal/middlewares"
	"backend/internal/store"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

// GetUserFriendWishlistItems godoc
// @Summary returns wishlist items
// @Tags Friends
// @Router /api/user/friend/wishlist/items [get]
// @Security ApiKeyAuth
// @Param wishlist_id query int true "Wishlist ID"
// @Produce json
// @Success 200 {array} WishlistItem
func (s *Service) GetUserFriendWishlistItems(c *gin.Context) {
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

	_, err = s.db.CheckIfUserHasAccessToWishlist(c, store.CheckIfUserHasAccessToWishlistParams{
		ID:     wishlistID,
		UserID: authData.User.ID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.Status(http.StatusNotFound)
			return
		}

		c.Error(err)
		c.Status(http.StatusInternalServerError)
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

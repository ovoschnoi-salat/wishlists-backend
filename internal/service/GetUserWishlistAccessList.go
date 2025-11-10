package service

import (
	"backend/internal/middlewares"
	"backend/internal/store"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetUserWishlistAccessList godoc
// @Summary returns wishlist access list
// @Tags User
// @Router /api/user/wishlist/access [get]
// @Security ApiKeyAuth
// @Param wishlist_id query int true "Wishlist ID"
// @Produce json
// @Success 200 {array} number
func (s *Service) GetUserWishlistAccessList(c *gin.Context) {
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

	accessList, err := s.db.GetWishlistAccessList(c, store.GetWishlistAccessListParams{
		ListID:  wishlistID,
		OwnerID: authData.User.ID,
	})
	if err != nil {
		c.Error(err)
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, mapStoreWishlistAccessListToWishlistAccessList(accessList))
}

func mapStoreWishlistAccessListToWishlistAccessList(items []store.WishlistAccessList) []int64 {
	res := make([]int64, len(items))
	for i, item := range items {
		res[i] = item.UserID
	}
	return res
}

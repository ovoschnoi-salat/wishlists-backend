package service

import (
	"backend/internal/middlewares"
	"backend/internal/store"
	"fmt"
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

	wishlistIDStr := c.Query("wishlist_id")
	wishlistID, err := strconv.ParseInt(wishlistIDStr, 10, 64)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("invalid wishlist ID: %w", err))
		return
	}

	accessList, err := s.db.GetWishlistAccessList(c, store.GetWishlistAccessListParams{
		ListID:  wishlistID,
		OwnerID: authData.User.ID,
	})
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to get access list: %w", err))
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

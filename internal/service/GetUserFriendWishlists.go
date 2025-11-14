package service

import (
	"backend/internal/middlewares"
	"backend/internal/store"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type FriendWishlist struct {
	ID    int64  `json:"id"`
	Title string `json:"title"`
}

// GetUserFriendWishlists godoc
// @Summary returns user's friend wishlists
// @Tags Friends
// @Router /api/user/friend/wishlists [get]
// @Security ApiKeyAuth
// @Param friend_id query int true "Friend ID"
// @Produce json
// @Success 200 {array} Wishlist
func (s *Service) GetUserFriendWishlists(c *gin.Context) {
	authData := middlewares.GetInitDataFromContext(c)
	if authData == nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	friendIDStr := c.Query("friend_id")
	friendID, err := strconv.ParseInt(friendIDStr, 10, 64)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("invalid friend ID: %w", err))
		return
	}

	wishlists, err := s.db.GetFriendWishlists(c, store.GetFriendWishlistsParams{
		OwnerID: friendID,
		UserID:  authData.User.ID,
	})
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to get friend wishlists: %w", err))
		return
	}
	c.JSON(http.StatusOK, mapStoreWishlistsToWishlists(wishlists))
}

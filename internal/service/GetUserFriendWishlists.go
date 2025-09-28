package service

import (
	"backend/internal/middlewares"
	"backend/internal/store"
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
// @Success 200 {array} FriendWishlist
func (s *Service) GetUserFriendWishlists(c *gin.Context) {
	authData := middlewares.GetInitDataFromContext(c)
	if authData == nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// Get wishlist ID from URL parameter
	friendIDStr := c.Query("friend_id")
	friendID, err := strconv.ParseInt(friendIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid friend ID"})
		return
	}

	wishlists, err := s.db.GetFriendWishlists(c, store.GetFriendWishlistsParams{
		OwnerID: friendID,
		UserID:  authData.User.ID,
	})
	if err != nil {
		c.Error(err)
		c.Status(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, mapStoreWishlistsToFriendWishlists(wishlists))
}

func mapStoreWishlistToFriendWishlists(wishlist store.Wishlist) FriendWishlist {
	return FriendWishlist{
		ID:    wishlist.ID,
		Title: wishlist.Title,
	}
}

func mapStoreWishlistsToFriendWishlists(wishlists []store.Wishlist) []FriendWishlist {
	res := make([]FriendWishlist, len(wishlists))
	for i, u := range wishlists {
		res[i] = mapStoreWishlistToFriendWishlists(u)
	}
	return res
}

package service

import (
	"backend/internal/middlewares"
	"backend/internal/store"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

type FriendWishlistItem struct {
	ID                       int64           `json:"id"`
	WishlistID               int64           `json:"wishlist_id"`
	Title                    string          `json:"title"`
	Description              string          `json:"description"`
	Price                    string          `json:"price"`
	Links                    json.RawMessage `json:"links"`
	Reservable               bool            `json:"reservable"`
	Reserved                 bool            `json:"reserved"`
	ReservationCanBeCanceled bool            `json:"reservation_can_be_canceled"`
}

// GetUserFriendWishlistItems godoc
// @Summary returns wishlist items
// @Tags Friends
// @Router /api/user/friend/wishlist/items [get]
// @Security ApiKeyAuth
// @Param wishlist_id query int true "Wishlist ID"
// @Produce json
// @Success 200 {array} FriendWishlistItem{links=[]WishlistItemLink}
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

	c.JSON(http.StatusOK, mapStoreWishlistItemsToFriendWishlistItems(items))
}

func mapStoreWishlistItemsToFriendWishlistItems(items []store.WishlistItem) []FriendWishlistItem {
	res := make([]FriendWishlistItem, len(items))
	for i, item := range items {
		res[i] = mapStoreWishlistItemToFriendWishlistItem(item)
	}
	return res
}

func mapStoreWishlistItemToFriendWishlistItem(item store.WishlistItem) FriendWishlistItem {
	return FriendWishlistItem{
		ID:          item.ID,
		WishlistID:  item.WishlistID,
		Title:       item.Title,
		Description: item.Description,
		Price:       item.Price,
		Links:       item.Links,
		Reservable:  item.Reservable,
	}
}

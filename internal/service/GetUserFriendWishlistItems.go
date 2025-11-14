package service

import (
	"backend/internal/middlewares"
	"backend/internal/store"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type FriendWishlistItem struct {
	ID                       int64              `json:"id"`
	WishlistID               int64              `json:"wishlist_id"`
	Title                    string             `json:"title"`
	Description              string             `json:"description"`
	Price                    string             `json:"price"`
	Links                    []WishlistItemLink `json:"links"`
	Reservable               bool               `json:"reservable"`
	Reserved                 bool               `json:"reserved"`
	ReservationCanBeCanceled bool               `json:"reservation_can_be_canceled"`
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

	wishlistIDStr := c.Query("wishlist_id")
	wishlistID, err := strconv.ParseInt(wishlistIDStr, 10, 64)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("invalid wishlist ID: %w", err))
		return
	}

	items, err := s.db.GetFriendWishlistItems(c, store.GetFriendWishlistItemsParams{
		WishlistID: wishlistID,
		FriendID:   authData.User.ID,
	})
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("error getting friend wishlist items: %w", err))
		return
	}

	c.JSON(http.StatusOK, mapStoreWishlistItemsToFriendWishlistItems(items, authData.User.ID))
}

func mapStoreWishlistItemsToFriendWishlistItems(items []store.WishlistItem, userID int64) []FriendWishlistItem {
	res := make([]FriendWishlistItem, len(items))
	for i, item := range items {
		res[i] = mapStoreWishlistItemToFriendWishlistItem(item, userID)
	}
	return res
}

func mapStoreWishlistItemToFriendWishlistItem(item store.WishlistItem, userID int64) FriendWishlistItem {
	// Parse links from JSON bytes
	var links []WishlistItemLink
	if len(item.Links) > 2 {
		json.Unmarshal(item.Links, &links)
	}

	return FriendWishlistItem{
		ID:                       item.ID,
		WishlistID:               item.WishlistID,
		Title:                    item.Title,
		Description:              item.Description,
		Price:                    item.Price,
		Links:                    links,
		Reservable:               item.Reservable,
		Reserved:                 item.ReservedBy.Valid,
		ReservationCanBeCanceled: item.ReservedBy.Valid && item.ReservedBy.Int64 == userID,
	}
}

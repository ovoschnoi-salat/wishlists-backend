package service

import (
	"backend/internal/middlewares"
	"backend/internal/store"
	"backend/internal/subcodeErrors"
	"backend/internal/subcodeErrors/codes"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetUserWishlistItems godoc
// @Summary returns wishlist items
// @Tags User
// @Router /api/user/wishlist/items [get]
// @Security ApiKeyAuth
// @Param wishlist_id query int true "Wishlist ID"
// @Produce json
// @Failure 400 {object} subcodeErrors.Response
// @Failure 401 {object} subcodeErrors.Response
// @Failure 500 {object} subcodeErrors.Response
// @Success 200 {array} WishlistItem
func (s *Service) GetUserWishlistItems(c *gin.Context) {
	authData, authorized := middlewares.GetInitDataFromContext(c)
	if !authorized {
		subcodeErrors.SendResponse(c, http.StatusInternalServerError, codes.InternalErrCode, noInitDataErr)
		return
	}

	wishlistIDStr := c.Query("wishlist_id")
	wishlistID, err := strconv.ParseInt(wishlistIDStr, 10, 64)
	if err != nil {
		subcodeErrors.SendResponse(c, http.StatusBadRequest, codes.InvalidRequestParametersErrCode, fmt.Errorf("invalid wishlist_id: %w", err))
		return
	}

	count, err := s.db.CheckUserOwnsWishlist(c, store.CheckUserOwnsWishlistParams{
		ID:      wishlistID,
		OwnerID: authData.User.ID,
	})
	if err != nil {
		subcodeErrors.SendResponse(c, http.StatusInternalServerError, codes.InternalErrCode, fmt.Errorf("error checking access to wish: %w", err))
		return
	}
	if count == 0 {
		subcodeErrors.SendResponse(c, http.StatusBadRequest, codes.InvalidRequestErrCode, errors.New("no access to wishlist"))
		return
	}

	items, err := s.db.GetWishlistItems(c, wishlistID)
	if err != nil {
		subcodeErrors.SendResponse(c, http.StatusInternalServerError, codes.InternalErrCode, fmt.Errorf("error getting wishlist items: %w", err))
		return
	}

	c.JSON(http.StatusOK, mapStoreWishlistItemsToWishlistItems(items))
}

func mapStoreWishlistItemsToWishlistItems(items []store.WishlistItem) []WishlistItem {
	res := make([]WishlistItem, len(items))
	for i, item := range items {
		res[i] = mapStoreWishToWishlistItem(item)
	}
	return res
}

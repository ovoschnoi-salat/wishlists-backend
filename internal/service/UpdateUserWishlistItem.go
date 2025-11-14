package service

import (
	"backend/internal/middlewares"
	"backend/internal/store"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

// UpdateUserWishlistItem godoc
// @Summary	updates a wishlist item
// @Tags	User
// @Router	/api/user/wishlist/item [patch]
// @Security	ApiKeyAuth
// @Accept	json
// @Param	item_id	query	int							true	"Item ID"
// @Param	item	body	CreateWishlistItemRequest	true	"Item"true "request body"
// @Success 200 {object} WishlistItem
func (s *Service) UpdateUserWishlistItem(c *gin.Context) {
	authData := middlewares.GetInitDataFromContext(c)
	if authData == nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// Get item ID from URL parameter
	itemIDRaw := c.Query("item_id")
	itemID, err := strconv.ParseInt(itemIDRaw, 10, 64)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("invalid item_id: %s", itemIDRaw))
		return
	}

	// Parse request body
	req := new(CreateWishlistItemRequest)
	err = c.BindJSON(req)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("invalid request body: %s", err))
		return
	}

	// Convert links to JSON bytes
	linksJSON, err := json.Marshal(req.Links)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("invalid link format: %s", err))
		return
	}

	// Update the wishlist item
	wishlistItem, err := s.db.UpdateWishlistItem(c, store.UpdateWishlistItemParams{
		ID:          itemID,
		OwnerID:     authData.User.ID,
		Title:       req.Title,
		Description: req.Description,
		Price:       req.Price,
		Links:       linksJSON,
		Reservable:  req.Reservable,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.AbortWithError(http.StatusBadRequest, fmt.Errorf("wish %s cannot be updated: %w", itemIDRaw, err))
			return
		}
		c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("error updating wishlist item: %w", err))
		return
	}

	if !wishlistItem.Reservable && wishlistItem.ReservedBy.Int64 != 0 {
		count, err := s.db.ResetWishlistItemReservation(c, wishlistItem.ID)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("error resetting wishlist item reservationfor item %d: %w", wishlistItem.ID, err))
			return
		}
		if count == 0 {
			c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("error resetting wishlist item reservation for item %d", wishlistItem.ID))
			return
		}
	}

	c.JSON(http.StatusOK, mapStoreWishlistItemToWishlistItem(wishlistItem))
}

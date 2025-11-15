package service

import (
	"backend/internal/middlewares"
	"backend/internal/store"
	"backend/internal/subcodeErrors"
	"backend/internal/subcodeErrors/codes"
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
// @Failure 400 {object} subcodeErrors.Response
// @Failure 401 {object} subcodeErrors.Response
// @Failure 500 {object} subcodeErrors.Response
// @Success 200 {object} WishlistItem
func (s *Service) UpdateUserWishlistItem(c *gin.Context) {
	authData, authorized := middlewares.GetInitDataFromContext(c)
	if !authorized {
		subcodeErrors.SendResponse(c, http.StatusInternalServerError, codes.InternalErrCode, noInitDataErr)
		return
	}

	// Get item ID from URL parameter
	itemIDRaw := c.Query("item_id")
	itemID, err := strconv.ParseInt(itemIDRaw, 10, 64)
	if err != nil {
		subcodeErrors.SendResponse(c, http.StatusBadRequest, codes.InvalidRequestParametersErrCode, fmt.Errorf("invalid item_id: %s", itemIDRaw))
		return
	}

	// Parse request body
	req := new(CreateWishlistItemRequest)
	err = c.ShouldBindJSON(req)
	if err != nil {
		subcodeErrors.SendResponse(c, http.StatusBadRequest, codes.InvalidRequestParametersErrCode, err)
		return
	}

	// Convert links to JSON bytes
	linksJSON, err := json.Marshal(req.Links)
	if err != nil {
		subcodeErrors.SendResponse(c, http.StatusBadRequest, codes.InvalidRequestParametersErrCode, fmt.Errorf("invalid link format: %s", err))
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
			subcodeErrors.SendResponse(c, http.StatusBadRequest, codes.InvalidRequestErrCode, fmt.Errorf("wish %s cannot be updated: %w", itemIDRaw, err))
			return
		}
		subcodeErrors.SendResponse(c, http.StatusInternalServerError, codes.InternalErrCode, fmt.Errorf("error updating wishlist item: %w", err))
		return
	}

	if !wishlistItem.Reservable && wishlistItem.ReservedBy.Int64 != 0 {
		count, err := s.db.ResetWishlistItemReservation(c, wishlistItem.ID)
		if err != nil {
			subcodeErrors.SendResponse(c, http.StatusInternalServerError, codes.InternalErrCode, fmt.Errorf("error resetting wishlist item reservationfor item %d: %w", wishlistItem.ID, err))
			return
		}
		if count == 0 {
			subcodeErrors.SendResponse(c, http.StatusInternalServerError, codes.InternalErrCode, fmt.Errorf("error resetting wishlist item reservation for item %d: no rows affected", wishlistItem.ID))
			return
		}
	}

	c.JSON(http.StatusOK, mapStoreWishlistItemToWishlistItem(wishlistItem))
}

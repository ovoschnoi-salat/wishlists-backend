package service

import (
	"backend/internal/errorResponse"
	"backend/internal/errorResponse/codes"
	"backend/internal/middlewares"
	"backend/internal/store"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

// GetUserWishlistItem godoc
// @Summary returns wishlist item
// @Tags User
// @Router /api/user/wishlist/item [get]
// @Security ApiKeyAuth
// @Param item_id query int true "Wishlist item ID"
// @Produce json
// @Failure 400 {object} errorResponse.Response
// @Failure 401 {object} errorResponse.Response
// @Failure 500 {object} errorResponse.Response
// @Success 200 {object} WishlistItem
func (s *Service) GetUserWishlistItem(c *gin.Context) {
	authData, authorized := middlewares.GetInitDataFromContext(c)
	if !authorized {
		errorResponse.Send(c, http.StatusUnauthorized, codes.UnauthorizedErrCode, nil)
		return
	}

	// Get wishlist ID from URL parameter
	itemIDStr := c.Query("item_id")
	itemID, err := strconv.ParseInt(itemIDStr, 10, 64)
	if err != nil {
		errorResponse.Send(c, http.StatusBadRequest, codes.UnauthorizedErrCode, fmt.Errorf("invalid item_id: %w", err))
		return
	}

	// Get wishlist items
	item, err := s.db.GetUserWishlistItem(c, store.GetUserWishlistItemParams{
		ID:      itemID,
		OwnerID: authData.User.ID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			errorResponse.Send(c, http.StatusNotFound, codes.WishNotFoundErrCode, fmt.Errorf("no items found"))
		}
		errorResponse.Send(c, http.StatusInternalServerError, codes.InternalErrCode, fmt.Errorf("failed to get item: %w", err))
		return
	}

	c.JSON(http.StatusOK, mapStoreWishlistItemToWishlistItem(item))
}

package service

import (
	"backend/internal/middlewares"
	"backend/internal/store"
	"backend/internal/subcodeErrors"
	"backend/internal/subcodeErrors/codes"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetUserFriendWishlistItem godoc
// @Summary returns friend wishlist item
// @Tags Friends
// @Router /api/user/friend/wishlist/item [get]
// @Security ApiKeyAuth
// @Param wish_id query int true "Wish ID"
// @Produce json
// @Failure 400 {object} subcodeErrors.Response
// @Failure 401 {object} subcodeErrors.Response
// @Failure 500 {object} subcodeErrors.Response
// @Success 200 {object} FriendWishlistItem
func (s *Service) GetUserFriendWishlistItem(c *gin.Context) {
	authData, authorized := middlewares.GetInitDataFromContext(c)
	if !authorized {
		subcodeErrors.SendResponse(c, http.StatusInternalServerError, codes.InternalErrCode, noInitDataErr)
		return
	}

	wishIDStr := c.Query("wish_id")
	wishID, err := strconv.ParseInt(wishIDStr, 10, 64)
	if err != nil {
		subcodeErrors.SendResponse(c, http.StatusBadRequest, codes.InvalidRequestParametersErrCode, fmt.Errorf("invalid wish_id: %w", err))
		return
	}

	item, err := s.db.GetWishlistItem(c, wishID)
	if err != nil {
		subcodeErrors.SendResponse(c, http.StatusInternalServerError, codes.InternalErrCode, fmt.Errorf("error getting friend wishlist item: %w", err))
		return
	}

	count, err := s.db.CheckIfUserHasAccessToWishlist(c, store.CheckIfUserHasAccessToWishlistParams{
		ID:       item.WishlistID,
		FriendID: authData.User.ID,
	})
	if err != nil {
		subcodeErrors.SendResponse(c, http.StatusInternalServerError, codes.InternalErrCode, fmt.Errorf("error checking if user has access to friend wishlist item: %w", err))
		return
	}
	if count == 0 {
		subcodeErrors.SendResponse(c, http.StatusNotFound, codes.WishNotFoundErrCode, fmt.Errorf("item not found"))
		return
	}

	c.JSON(http.StatusOK, mapStoreWishlistItemToFriendWishlistItem(item, authData.User.ID))
}

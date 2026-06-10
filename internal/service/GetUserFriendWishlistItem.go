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

	count, err := s.db.CheckUserHasAccessToWish(c, store.CheckUserHasAccessToWishParams{
		ID:       wishID,
		FriendID: authData.User.ID,
	})
	if err != nil {
		subcodeErrors.SendResponse(c, http.StatusInternalServerError, codes.InternalErrCode, fmt.Errorf("error checking access to wish: %w", err))
		return
	}
	if count == 0 {
		subcodeErrors.SendResponse(c, http.StatusBadRequest, codes.InvalidRequestErrCode, errors.New("no access to wish"))
		return
	}

	item, err := s.db.GetWish(c, wishID)
	if err != nil {
		subcodeErrors.SendResponse(c, http.StatusInternalServerError, codes.InternalErrCode, fmt.Errorf("error getting friend wishlist item: %w", err))
		return
	}

	c.JSON(http.StatusOK, mapStoreWishlistItemToFriendWishlistItem(item, authData.User.ID))
}

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
	"github.com/jackc/pgx/v5"
)

// CreateWishSplitRequest godoc
// @Summary creates a wishlist item split requests
// @Tags Friend's Wish
// @Router /api/user/friend/wishlist/wish/split-request [post]
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param wish_id query int true "wish id"
// @Failure 400 {object} subcodeErrors.Response
// @Failure 401 {object} subcodeErrors.Response
// @Failure 500 {object} subcodeErrors.Response
// @Success 204
func (s *Service) CreateWishSplitRequest(c *gin.Context) {
	authData, authorized := middlewares.GetInitDataFromContext(c)
	if !authorized {
		subcodeErrors.SendResponse(c, http.StatusInternalServerError, codes.InternalErrCode, noInitDataErr)
		return
	}

	// Get wishlist ID from URL parameter
	wishIDStr := c.Query("wish_id")
	wishID, err := strconv.ParseInt(wishIDStr, 10, 64)
	if err != nil {
		subcodeErrors.SendResponse(c, http.StatusBadRequest, codes.UnauthorizedErrCode, fmt.Errorf("invalid wish_id: %w", err))
		return
	}

	// Get wishlist item
	wish, err := s.db.GetWish(c, wishID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			subcodeErrors.SendResponse(c, http.StatusNotFound, codes.WishNotFoundErrCode, fmt.Errorf("no items found"))
			return
		}
		subcodeErrors.SendResponse(c, http.StatusInternalServerError, codes.InternalErrCode, fmt.Errorf("failed to get item: %w", err))
		return
	}

	if wish.OwnerID != authData.User.ID {
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
	}

	count, err := s.db.CreateWishSplitRequest(c, store.CreateWishSplitRequestParams{
		ListID:  wish.WishlistID,
		OwnerID: wish.OwnerID,
		WishID:  wish.ID,
		UserID:  authData.User.ID,
	})
	if err != nil {
		subcodeErrors.SendResponse(c, http.StatusInternalServerError, codes.InternalErrCode, fmt.Errorf("failed to create split request: %w", err))
		return
	}
	if count == 0 {
		subcodeErrors.SendResponse(c, http.StatusInternalServerError, codes.InternalErrCode, fmt.Errorf("failed to create split request: rows count 0"))
		return
	}

	c.Status(http.StatusNoContent)
}

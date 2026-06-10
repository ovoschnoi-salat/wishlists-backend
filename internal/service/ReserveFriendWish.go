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
	"github.com/jackc/pgx/v5/pgtype"
)

// ReserveFriendWish godoc
// @Summary		reserves friend's wish
// @Tags		Friend's Wish
// @Router		/api/user/friend/wishlist/wish/reservation/reserve [post]
// @Security	ApiKeyAuth
// @Param		wish_id query int true "Wish ID"
// @Failure 400 {object} subcodeErrors.Response
// @Failure 401 {object} subcodeErrors.Response
// @Failure 500 {object} subcodeErrors.Response
// @Success		204
func (s *Service) ReserveFriendWish(c *gin.Context) {
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

	count, err = s.db.ReserveWishlistItem(c, store.ReserveWishlistItemParams{
		ID:         wishID,
		ReservedBy: pgtype.Int8{Int64: authData.User.ID, Valid: true},
	})
	if err != nil {
		subcodeErrors.SendResponse(c, http.StatusInternalServerError, codes.InternalErrCode, fmt.Errorf("error reserving wish: %w", err))
		return
	}
	if count == 0 {
		subcodeErrors.SendResponse(c, http.StatusBadRequest, codes.InvalidRequestErrCode, errors.New("wish cannot be reserved"))
		return
	}

	c.Status(http.StatusNoContent)
}

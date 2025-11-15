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

// CancelFriendWishReservation godoc
// @Summary		reserves friend's wish
// @Tags		Friend's Wish
// @Router		/api/user/friend/wishlist/wish/reservation/cancel [post]
// @Security	ApiKeyAuth
// @Param		wish_id	query	int	true	"Wish ID"
// @Failure 400 {object} subcodeErrors.Response
// @Failure 401 {object} subcodeErrors.Response
// @Failure 500 {object} subcodeErrors.Response
// @Success		204
func (s *Service) CancelFriendWishReservation(c *gin.Context) {
	authData, authorized := middlewares.GetInitDataFromContext(c)
	if !authorized {
		subcodeErrors.SendResponse(c, http.StatusInternalServerError, codes.InternalErrCode, noInitDataErr)
		return
	}

	wishIDStr := c.Query("wish_id")
	wishID, err := strconv.ParseInt(wishIDStr, 10, 64)
	if err != nil {
		subcodeErrors.SendResponse(c, http.StatusBadRequest, codes.InvalidRequestParametersErrCode, fmt.Errorf("invalid wish ID: %w", err))
		return
	}

	count, err := s.db.CancelWishlistItemReservation(c, store.CancelWishlistItemReservationParams{
		ID:         wishID,
		ReservedBy: pgtype.Int8{Int64: authData.User.ID, Valid: true},
	})
	if err != nil {
		subcodeErrors.SendResponse(c, http.StatusInternalServerError, codes.InternalErrCode, fmt.Errorf("error cancelling reservation: %w", err))
		return
	}
	if count == 0 {
		subcodeErrors.SendResponse(c, http.StatusBadRequest, codes.InvalidRequestErrCode, errors.New("error cancelling reservation: no rows affected"))
		return
	}

	c.Status(http.StatusNoContent)
}

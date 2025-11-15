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
	"github.com/jackc/pgx/v5/pgtype"
)

// ReserveFriendWish godoc
// @Summary		reserves friend's wish
// @Tags		Friend's Wish
// @Router		/api/user/friend/wishlist/wish/reservation/reserve [post]
// @Security	ApiKeyAuth
// @Param		wish_id query int true "Wish ID"
// @Failure 400 {object} Response
// @Failure 401 {object} Response
// @Failure 500 {object} Response
// @Success		204
func (s *Service) ReserveFriendWish(c *gin.Context) {
	authData, authorized := middlewares.GetInitDataFromContext(c)
	if !authorized {
		errorResponse.Send(c, http.StatusUnauthorized, codes.UnauthorizedErrCode, nil)
		return
	}

	wishIDStr := c.Query("wish_id")
	wishID, err := strconv.ParseInt(wishIDStr, 10, 64)
	if err != nil {
		errorResponse.Send(c, http.StatusBadRequest, codes.InvalidRequestParametersErrCode, fmt.Errorf("invalid wish_id: %w", err))
		return
	}

	count, err := s.db.ReserveWishlistItem(c, store.ReserveWishlistItemParams{
		ID:         wishID,
		ReservedBy: pgtype.Int8{Int64: authData.User.ID, Valid: true},
	})
	if err != nil {
		errorResponse.Send(c, http.StatusInternalServerError, codes.InternalErrCode, fmt.Errorf("error reserving wish: %w", err))
		return
	}
	if count == 0 {
		errorResponse.Send(c, http.StatusBadRequest, codes.InvalidRequestErrCode, errors.New("wish cannot be reserved"))
		return
	}

	c.Status(http.StatusNoContent)
}

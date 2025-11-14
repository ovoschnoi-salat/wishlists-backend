package service

import (
	"backend/internal/middlewares"
	"backend/internal/store"
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
// @Success		204
func (s *Service) CancelFriendWishReservation(c *gin.Context) {
	authData := middlewares.GetInitDataFromContext(c)
	if authData == nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	wishIDStr := c.Query("wish_id")
	wishID, err := strconv.ParseInt(wishIDStr, 10, 64)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("invalid wish ID: %w", err))
		return
	}

	count, err := s.db.CancelWishlistItemReservation(c, store.CancelWishlistItemReservationParams{
		ID:         wishID,
		ReservedBy: pgtype.Int8{Int64: authData.User.ID, Valid: true},
	})
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("error cancelling reservation: %w", err))
		return
	}
	if count < 1 {
		c.AbortWithError(http.StatusBadRequest, errors.New("wish reservation cannot be cancelled"))
		return
	}

	c.Status(http.StatusNoContent)
}

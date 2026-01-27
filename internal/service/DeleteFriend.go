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

// DeleteFriend godoc
// @Summary creates wishlist
// @Tags Friends
// @Router /api/user/friend [delete]
// @Security ApiKeyAuth
// @Accept	json
// @Param	friend_id	query	int	true	"Friend ID"
// @Produce	json
// @Failure 400 {object} subcodeErrors.Response
// @Failure 401 {object} subcodeErrors.Response
// @Failure 500 {object} subcodeErrors.Response
// @Success 204
func (s *Service) DeleteFriend(c *gin.Context) {
	authData, authorized := middlewares.GetInitDataFromContext(c)
	if !authorized {
		subcodeErrors.SendResponse(c, http.StatusInternalServerError, codes.InternalErrCode, noInitDataErr)
		return
	}

	friendIDRaw := c.Query("friend_id")
	friendID, err := strconv.ParseInt(friendIDRaw, 10, 64)
	if err != nil {
		subcodeErrors.SendResponse(c, http.StatusBadRequest, codes.InvalidRequestParametersErrCode, fmt.Errorf("invalid friend_id: %w", err))
		return
	}

	count, err := s.db.DeleteFriendship(c, store.DeleteFriendshipParams{
		UserID:   authData.User.ID,
		FriendID: friendID,
	})
	if err != nil {
		subcodeErrors.SendResponse(c, http.StatusInternalServerError, codes.InternalErrCode, fmt.Errorf("error deleting friend: %w", err))
		return
	}
	if count == 0 {
		subcodeErrors.SendResponse(c, http.StatusBadRequest, codes.InvalidRequestErrCode, errors.New("friend not found"))
		return
	}

	_, err = s.db.ResetWishlistItemsReservationsForFriend(c, store.ResetWishlistItemsReservationsForFriendParams{
		OwnerID:    friendID,
		ReservedBy: pgtype.Int8{Int64: authData.User.ID, Valid: true},
	})
	if err != nil {
		subcodeErrors.SendResponse(c, http.StatusInternalServerError, codes.InternalErrCode, fmt.Errorf("error deleting friend: %w", err))
		return
	}

	c.Status(http.StatusNoContent)
}

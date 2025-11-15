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

type FriendWishlist struct {
	ID    int64  `json:"id"`
	Title string `json:"title"`
}

// GetUserFriendWishlists godoc
// @Summary returns user's friend wishlists
// @Tags Friends
// @Router /api/user/friend/wishlists [get]
// @Security ApiKeyAuth
// @Param friend_id query int true "Friend ID"
// @Produce json
// @Failure 400 {object} subcodeErrors.Response
// @Failure 401 {object} subcodeErrors.Response
// @Failure 500 {object} subcodeErrors.Response
// @Success 200 {array} Wishlist
func (s *Service) GetUserFriendWishlists(c *gin.Context) {
	authData, authorized := middlewares.GetInitDataFromContext(c)
	if !authorized {
		subcodeErrors.SendResponse(c, http.StatusInternalServerError, codes.InternalErrCode, noInitDataErr)
		return
	}

	friendIDStr := c.Query("friend_id")
	friendID, err := strconv.ParseInt(friendIDStr, 10, 64)
	if err != nil {
		subcodeErrors.SendResponse(c, http.StatusBadRequest, codes.InvalidRequestParametersErrCode, fmt.Errorf("invalid friend_id: %w", err))
		return
	}

	wishlists, err := s.db.GetFriendWishlists(c, store.GetFriendWishlistsParams{
		OwnerID: friendID,
		UserID:  authData.User.ID,
	})
	if err != nil {
		subcodeErrors.SendResponse(c, http.StatusInternalServerError, codes.InternalErrCode, fmt.Errorf("failed to get friend wishlists: %w", err))
		return
	}

	c.JSON(http.StatusOK, mapStoreWishlistsToWishlists(wishlists))
}

package service

import (
	"backend/internal/middlewares"
	"backend/internal/store"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

// UnreserveFriendWish godoc
// @Summary reserves friend's wish
// @Tags Friend's Wish
// @Router /api/user/friend/wishlist/wish/unreserve [post]
// @Security ApiKeyAuth
// @Param wish_id query int true "Wish ID"
// @Success 204
func (s *Service) UnreserveFriendWish(c *gin.Context) {
	authData := middlewares.GetInitDataFromContext(c)
	if authData == nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// Get item ID from URL parameter
	wishIDStr := c.Query("wish_id")
	wishID, err := strconv.ParseInt(wishIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid wish ID"})
		return
	}

	wishlist, err := s.db.GetWishlistByWishId(c, wishID)
	if err != nil {
		c.Error(err)
		c.Status(http.StatusInternalServerError)
		return
	}

	_, err = s.db.CheckIfFriends(c, store.CheckIfFriendsParams{
		UserID:   authData.User.ID,
		FriendID: wishlist.OwnerID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid wish ID"})
			return
		}
		c.Error(err)
		c.Status(http.StatusInternalServerError)
		return
	}

	if wishlist.IsPrivate {
		_, err = s.db.CheckUserHasAccessToPrivateWishlist(c, store.CheckUserHasAccessToPrivateWishlistParams{
			ListID: wishlist.ID,
			UserID: authData.User.ID,
		})
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid wish ID"})
				return
			}
			c.Error(err)
			c.Status(http.StatusInternalServerError)
			return
		}
	}

	count, err := s.db.UnreserveWishlistItem(c, store.UnreserveWishlistItemParams{
		ID:         wishID,
		ReservedBy: pgtype.Int8{Int64: authData.User.ID, Valid: true},
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	if count < 1 {
		c.Error(errors.New("no rows updated"))
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusNoContent)
}

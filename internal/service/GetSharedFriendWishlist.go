package service

import (
	"backend/internal/middlewares"
	"backend/internal/store"
	"backend/internal/subcodeErrors"
	"backend/internal/subcodeErrors/codes"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

// GetSharedWishlist godoc
// @Summary returns shared wishlist
// @Tags Friends
// @Router /api/shared/wishlist [get]
// @Security ApiKeyAuth
// @Param wishlist_uuid query string true "Wishlist UUID"
// @Produce json
// @Failure 400 {object} subcodeErrors.Response
// @Failure 401 {object} subcodeErrors.Response
// @Failure 500 {object} subcodeErrors.Response
// @Success 200 {object} FriendWishlist
func (s *Service) GetSharedWishlist(c *gin.Context) {
	authData, authorized := middlewares.GetInitDataFromContext(c)
	if !authorized {
		subcodeErrors.SendResponse(c, http.StatusInternalServerError, codes.InternalErrCode, noInitDataErr)
		return
	}

	wishlistUUIDStr := c.Query("wishlist_uuid")
	wishlistUUID, err := uuid.Parse(wishlistUUIDStr)
	if err != nil {
		subcodeErrors.SendResponse(c, http.StatusBadRequest, codes.InvalidRequestParametersErrCode, fmt.Errorf("invalid wishlist_uuid: %w", err))
		return
	}

	wishlist, err := s.db.GetSharedWishlist(c, UUIDToPgUUID(wishlistUUID))
	if err != nil {
		subcodeErrors.SendResponse(c, http.StatusInternalServerError, codes.InternalErrCode, fmt.Errorf("failed to get friend wishlists: %w", err))
		return
	}

	_, err = s.db.AddToFriends(c, store.AddToFriendsParams{
		UserID:   authData.User.ID,
		FriendID: wishlist.OwnerID,
	})
	if err != nil {
		subcodeErrors.SendResponse(c, http.StatusInternalServerError, codes.InternalErrCode, fmt.Errorf("failed to add friend: %w", err))
		return
	}

	if wishlist.IsPrivate {
		_, err := s.db.InsertWishlistAccessItem(c, store.InsertWishlistAccessItemParams{
			ListID:  wishlist.ID,
			OwnerID: wishlist.OwnerID,
			UserID:  authData.User.ID,
		})
		if err != nil {
			subcodeErrors.SendResponse(c, http.StatusInternalServerError, codes.InternalErrCode, fmt.Errorf("failed to grant access to wishlist: %w", err))
			return
		}
	}

	c.JSON(http.StatusOK, mapStoreWishlistToWishlist(wishlist))
}

func UUIDToPgUUID(uuid uuid.UUID) pgtype.UUID {
	return pgtype.UUID{
		Bytes: uuid,
		Valid: true,
	}
}

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
	"github.com/samber/lo"
)

type GetSplitRequestsResponse struct {
	SplitRequests         []Friend `json:"split_requests"`
	RequestFromUserExists bool     `json:"request_from_user_exists"`
}

// GetWishSplitRequests godoc
// @Summary returns wish split requests
// @Tags Wish
// @Router /api/user/wishlist/item/split-requests [get]
// @Security ApiKeyAuth
// @Param item_id query int true "Wishlist item ID"
// @Produce json
// @Failure 400 {object} subcodeErrors.Response
// @Failure 401 {object} subcodeErrors.Response
// @Failure 404 {object} subcodeErrors.Response
// @Failure 500 {object} subcodeErrors.Response
// @Success 200 {object} GetSplitRequestsResponse
func (s *Service) GetWishSplitRequests(c *gin.Context) {
	authData, authorized := middlewares.GetInitDataFromContext(c)
	if !authorized {
		subcodeErrors.SendResponse(c, http.StatusInternalServerError, codes.InternalErrCode, noInitDataErr)
		return
	}

	// Get wishlist ID from URL parameter
	wishIDStr := c.Query("item_id")
	wishID, err := strconv.ParseInt(wishIDStr, 10, 64)
	if err != nil {
		subcodeErrors.SendResponse(c, http.StatusBadRequest, codes.UnauthorizedErrCode, fmt.Errorf("invalid item_id: %w", err))
		return
	}

	// Get wishlist item
	wishlist, err := s.db.GetWishlistByWishId(c, wishID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			subcodeErrors.SendResponse(c, http.StatusNotFound, codes.WishNotFoundErrCode, fmt.Errorf("no wish found"))
			return
		}
		subcodeErrors.SendResponse(c, http.StatusInternalServerError, codes.InternalErrCode, fmt.Errorf("failed to get item: %w", err))
		return
	}

	if wishlist.OwnerID != authData.User.ID {
		access, err := s.db.CheckIfUserHasAccessToWishlist(c, store.CheckIfUserHasAccessToWishlistParams{
			ID:       wishlist.ID,
			FriendID: authData.User.ID,
		})
		if err != nil {
			subcodeErrors.SendResponse(c, http.StatusInternalServerError, codes.InternalErrCode, fmt.Errorf("failed to check item access: %w", err))
			return
		}
		if access == 0 {
			subcodeErrors.SendResponse(c, http.StatusNotFound, codes.WishNotFoundErrCode, fmt.Errorf("no items found"))
			return
		}
	} else if wishlist.SplitRequestPrivacy == store.SplitRequestPrivacyInvisibleToOwner {
		subcodeErrors.SendResponse(c, http.StatusNotFound, codes.WishNotFoundErrCode, fmt.Errorf("no access to split requests"))
		return
	}

	splitRequests, err := s.db.GetWishSplitRequests(c, wishID)
	if err != nil {
		subcodeErrors.SendResponse(c, http.StatusInternalServerError, codes.InternalErrCode, fmt.Errorf("failed to get wish split requests: %w", err))
		return
	}

	c.JSON(http.StatusOK, mapStoreUsersToSplitRequests(splitRequests, authData.User.ID))
}

func mapStoreUsersToSplitRequests(splitRequests []store.User, userID int64) GetSplitRequestsResponse {
	return GetSplitRequestsResponse{
		SplitRequests: mapStoreUsersToFriends(splitRequests),
		RequestFromUserExists: lo.ContainsBy(splitRequests, func(item store.User) bool {
			return item.ID == userID
		}),
	}
}

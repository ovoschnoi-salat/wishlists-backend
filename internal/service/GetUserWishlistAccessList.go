package service

import (
	"backend/internal/errorResponse"
	"backend/internal/errorResponse/codes"
	"backend/internal/middlewares"
	"backend/internal/store"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetUserWishlistAccessList godoc
// @Summary returns wishlist access list
// @Tags User
// @Router /api/user/wishlist/access [get]
// @Security ApiKeyAuth
// @Param wishlist_id query int true "Wishlist ID"
// @Produce json
// @Failure 400 {object} Response
// @Failure 401 {object} Response
// @Failure 500 {object} Response
// @Success 200 {array} number
func (s *Service) GetUserWishlistAccessList(c *gin.Context) {
	authData, authorized := middlewares.GetInitDataFromContext(c)
	if !authorized {
		errorResponse.Send(c, http.StatusUnauthorized, codes.UnauthorizedErrCode, nil)
		return
	}

	wishlistIDStr := c.Query("wishlist_id")
	wishlistID, err := strconv.ParseInt(wishlistIDStr, 10, 64)
	if err != nil {
		errorResponse.Send(c, http.StatusBadRequest, codes.InvalidRequestParametersErrCode, fmt.Errorf("invalid wishlist_id: %w", err))
		return
	}

	accessList, err := s.db.GetWishlistAccessList(c, store.GetWishlistAccessListParams{
		ListID:  wishlistID,
		OwnerID: authData.User.ID,
	})
	if err != nil {
		errorResponse.Send(c, http.StatusInternalServerError, codes.InternalErrCode, fmt.Errorf("failed to get access list: %w", err))
		return
	}

	c.JSON(http.StatusOK, mapStoreWishlistAccessListToWishlistAccessList(accessList))
}

func mapStoreWishlistAccessListToWishlistAccessList(items []store.WishlistAccessList) []int64 {
	res := make([]int64, len(items))
	for i, item := range items {
		res[i] = item.UserID
	}
	return res
}

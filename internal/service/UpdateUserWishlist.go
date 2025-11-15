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
	"github.com/jackc/pgx/v5"
)

// UpdateUserWishlist godoc
// @Summary creates wishlist
// @Tags User
// @Router /api/user/wishlist [patch]
// @Security ApiKeyAuth
// @Accept	json
// @Param	wishlist_id	query	int						true	"Wishlist ID"
// @Param	wishlist	body	CreateWishlistRequest	true	"request body"
// @Produce	json
// @Failure 400 {object} Response
// @Failure 401 {object} Response
// @Failure 500 {object} Response
// @Success 200 {object} Wishlist
func (s *Service) UpdateUserWishlist(c *gin.Context) {
	authData, authorized := middlewares.GetInitDataFromContext(c)
	if !authorized {
		errorResponse.Send(c, http.StatusUnauthorized, codes.UnauthorizedErrCode, nil)
		return
	}

	wishlistIDRaw := c.Query("wishlist_id")
	wishlistID, err := strconv.ParseInt(wishlistIDRaw, 10, 64)
	if err != nil {
		errorResponse.Send(c, http.StatusBadRequest, codes.InvalidRequestParametersErrCode, fmt.Errorf("invalid wishlist_id: %w", err))
		return
	}

	req := new(CreateWishlistRequest)
	if err := c.ShouldBindJSON(req); err != nil {
		errorResponse.Send(c, http.StatusBadRequest, codes.InvalidRequestParametersErrCode, err)
		return
	}

	wishlist, err := s.db.UpdateWishlist(c, store.UpdateWishlistParams{
		ID:          wishlistID,
		OwnerID:     authData.User.ID,
		Title:       req.Title,
		Description: req.Description,
		IsPrivate:   req.IsPrivate,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			errorResponse.Send(c, http.StatusBadRequest, codes.InvalidRequestErrCode, fmt.Errorf("can't update wishlist %s: %w", wishlistIDRaw, err))
			return
		}
		errorResponse.Send(c, http.StatusInternalServerError, codes.InternalErrCode, fmt.Errorf("failed to update wishlist %s: %w", wishlistIDRaw, err))
		return
	}

	if req.IsPrivate {
		accessList, err := s.db.GetWishlistAccessList(c, store.GetWishlistAccessListParams{
			ListID:  wishlistID,
			OwnerID: authData.User.ID,
		})
		if err != nil {
			errorResponse.Send(c, http.StatusInternalServerError, codes.InternalErrCode, fmt.Errorf("failed to get wishlist %s access list: %w", wishlistIDRaw, err))
			return
		}
		usersWithAccessOld := make(map[int64]struct{}, len(accessList))
		for _, accessItem := range accessList {
			usersWithAccessOld[accessItem.UserID] = struct{}{}
		}
		usersWithAccessNew := make(map[int64]struct{}, len(req.UsersWithAccess))
		for _, usersWithAccess := range req.UsersWithAccess {
			usersWithAccessNew[usersWithAccess] = struct{}{}
		}

		for userID := range usersWithAccessOld {
			if _, ok := usersWithAccessNew[userID]; !ok {
				count, err := s.db.DeleteWishlistAccessItem(c, store.DeleteWishlistAccessItemParams{
					ListID: wishlistID,
					UserID: userID,
				})
				if err != nil {
					errorResponse.Send(c, http.StatusInternalServerError, codes.InternalErrCode, fmt.Errorf("error deleting access for user %d from wishlist %s: %w", userID, wishlistIDRaw, err))
					return
				}
				if count == 0 {
					errorResponse.Send(c, http.StatusInternalServerError, codes.InternalErrCode, fmt.Errorf("error deleting access for user %d from wishlist %s: no rows affected", userID, wishlistIDRaw))
					return
				}
			}
		}
		for userID := range usersWithAccessNew {
			if _, ok := usersWithAccessOld[userID]; !ok {
				count, err := s.db.InsertWishlistAccessItem(c, store.InsertWishlistAccessItemParams{
					ListID:  wishlistID,
					OwnerID: authData.User.ID,
					UserID:  userID,
				})
				if err != nil {
					errorResponse.Send(c, http.StatusInternalServerError, codes.InternalErrCode, fmt.Errorf("error inserting access for user %d to wishlist %s: %w", userID, wishlistIDRaw, err))
					return
				}
				if count == 0 {
					errorResponse.Send(c, http.StatusInternalServerError, codes.InternalErrCode, fmt.Errorf("error inserting access for user %d to wishlist %s: no rows affected", userID, wishlistIDRaw))
					return
				}
			}
		}
	} else {
		err := s.db.DeleteWishlistAccessItems(c, wishlistID)
		if err != nil {
			errorResponse.Send(c, http.StatusInternalServerError, codes.InternalErrCode, fmt.Errorf("failed to delete wishlist %s access items: %w", wishlistIDRaw, err))
			return
		}
	}

	c.JSON(http.StatusOK, mapStoreWishlistToWishlist(wishlist))
}

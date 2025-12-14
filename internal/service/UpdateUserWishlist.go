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
// @Failure 400 {object} subcodeErrors.Response
// @Failure 401 {object} subcodeErrors.Response
// @Failure 500 {object} subcodeErrors.Response
// @Success 200 {object} Wishlist
func (s *Service) UpdateUserWishlist(c *gin.Context) {
	authData, authorized := middlewares.GetInitDataFromContext(c)
	if !authorized {
		subcodeErrors.SendResponse(c, http.StatusInternalServerError, codes.InternalErrCode, noInitDataErr)
		return
	}

	wishlistIDRaw := c.Query("wishlist_id")
	wishlistID, err := strconv.ParseInt(wishlistIDRaw, 10, 64)
	if err != nil {
		subcodeErrors.SendResponse(c, http.StatusBadRequest, codes.InvalidRequestParametersErrCode, fmt.Errorf("invalid wishlist_id: %w", err))
		return
	}

	req := new(CreateWishlistRequest)
	if err := c.ShouldBindJSON(req); err != nil {
		subcodeErrors.SendResponse(c, http.StatusBadRequest, codes.InvalidRequestParametersErrCode, err)
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
			subcodeErrors.SendResponse(c, http.StatusBadRequest, codes.InvalidRequestErrCode, fmt.Errorf("can't update wishlist %s: %w", wishlistIDRaw, err))
			return
		}
		subcodeErrors.SendResponse(c, http.StatusInternalServerError, codes.InternalErrCode, fmt.Errorf("failed to update wishlist %s: %w", wishlistIDRaw, err))
		return
	}

	if req.IsPrivate {
		err := s.db.RecreateWishlistAccessList(c, store.RecreateAccessListParams{
			WishlistId:    wishlist.ID,
			OwnerID:       authData.User.ID,
			NewFriendsIDs: req.UsersWithAccess,
		})
		if err != nil {
			_ = c.Error(fmt.Errorf("error updating wishlist access items: %w", err))
		}
	} else {
		err := s.db.DeleteWishlistAccessItems(c, wishlistID)
		if err != nil {
			subcodeErrors.SendResponse(c, http.StatusInternalServerError, codes.InternalErrCode, fmt.Errorf("failed to delete wishlist %s access items: %w", wishlistIDRaw, err))
			return
		}
	}

	c.JSON(http.StatusOK, mapStoreWishlistToWishlist(wishlist))
}

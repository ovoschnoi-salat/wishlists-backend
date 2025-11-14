package service

import (
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
// @Success 200 {object} Wishlist
func (s *Service) UpdateUserWishlist(c *gin.Context) {
	authData := middlewares.GetInitDataFromContext(c)
	if authData == nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	wishlistIDRaw := c.Query("wishlist_id")
	wishlistID, err := strconv.ParseInt(wishlistIDRaw, 10, 64)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("invalid wishlist_id: %w", err))
		return
	}

	req := new(CreateWishlistRequest)
	if err := c.BindJSON(req); err != nil {
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
			c.AbortWithError(http.StatusBadRequest, fmt.Errorf("can' update wishlist %s: %w", wishlistIDRaw, err))
			return
		}
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if req.IsPrivate {
		accessList, err := s.db.GetWishlistAccessList(c, store.GetWishlistAccessListParams{
			ListID:  wishlistID,
			OwnerID: authData.User.ID,
		})
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
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
					c.AbortWithError(http.StatusInternalServerError, err)
					return
				}
				if count == 0 {
					c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("error deleting access for user %d", userID))
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
					c.AbortWithError(http.StatusInternalServerError, err)
					return
				}
				if count == 0 {
					c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("error inserting access for user %d", userID))
					return
				}
			}
		}
	} else {
		err := s.db.DeleteWishlistAccessItems(c, wishlistID)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
	}

	c.JSON(http.StatusOK, mapStoreWishlistToWishlist(wishlist))
}

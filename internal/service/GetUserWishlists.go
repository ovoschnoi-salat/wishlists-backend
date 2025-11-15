package service

import (
	"backend/internal/errors"
	"backend/internal/errors/codes"
	"backend/internal/middlewares"
	"backend/internal/store"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Wishlist struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	IsPrivate   bool   `json:"is_private"`
}

// GetUserWishlists godoc
// @Summary returns user's wishlists
// @Tags User
// @Router /api/user/wishlists [get]
// @Security ApiKeyAuth
// @Produce json
// @Failure 401 {object} errors.Response
// @Failure 500 {object} errors.Response
// @Success 200 {array} Wishlist
func (s *Service) GetUserWishlists(c *gin.Context) {
	authData, authorized := middlewares.GetInitDataFromContext(c)
	if !authorized {
		errors.SendResponse(c, http.StatusInternalServerError, codes.InternalErrCode, noInitDataErr)
		return
	}

	wishlists, err := s.db.GetUserWishlists(c, authData.User.ID)
	if err != nil {
		errors.SendResponse(c, http.StatusInternalServerError, codes.InternalErrCode, fmt.Errorf("error getting wishlists: %w", err))
		return
	}
	c.JSON(http.StatusOK, mapStoreWishlistsToWishlists(wishlists))
}

func mapStoreWishlistToWishlist(wishlist store.Wishlist) Wishlist {
	return Wishlist{
		ID:          wishlist.ID,
		Title:       wishlist.Title,
		Description: wishlist.Description,
		IsPrivate:   wishlist.IsPrivate,
	}
}

func mapStoreWishlistsToWishlists(wishlist []store.Wishlist) []Wishlist {
	res := make([]Wishlist, len(wishlist))
	for i, w := range wishlist {
		res[i] = mapStoreWishlistToWishlist(w)
	}
	return res
}

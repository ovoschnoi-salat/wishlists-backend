package service

import (
	"backend/internal/store"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Wishlist struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	IsPrivate   bool   `json:"is_private"`
}

// GetMyWishlists godoc
// @Summary returns user's wishlists
// @Schemes
// @Description
// @Tags user
// @Accept json
// @Produce json
// @Success 200 {array} Wishlist
// @Router /user/wishlists [get]
func (s *Service) GetMyWishlists(c *gin.Context) {
	wishlists, err := s.db.GetWishlists(c, 1)
	if err != nil {
		c.Error(err)
		c.Status(http.StatusInternalServerError)
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

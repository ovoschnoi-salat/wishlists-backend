package service

import (
	"backend/internal/store"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CreateWishlistRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	IsPrivate   bool   `json:"is_private"`
}

// CreateWishlist godoc
// @Summary creates wishlist
// @Description
// @Tags user
// @Param			wishlist	body		CreateWishlistRequest	true "request body"
// @Accept json
// @Produce json
// @Success 200 {object} Wishlist
// @Router /user/wishlists [post]
func (s *Service) CreateWishlist(c *gin.Context) {
	req := new(CreateWishlistRequest)
	err := c.BindJSON(req)
	if err != nil {
		return
	}
	wishlist, err := s.db.CreateWishlist(c, store.CreateWishlistParams{
		OwnerID:     1,
		Title:       req.Title,
		Description: req.Description,
		IsPrivate:   req.IsPrivate,
	})
	if err != nil {
		c.Error(err)
		c.Status(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, mapStoreWishlistToWishlist(wishlist))
}

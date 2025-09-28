package service

import (
	"backend/internal/middlewares"
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
// @Tags User
// @Router /api/user/wishlist [post]
// @Security ApiKeyAuth
// @Accept	json
// @Param	wishlist	body	CreateWishlistRequest	true	"request body"
// @Produce	json
// @Success 200 {object} Wishlist
func (s *Service) CreateWishlist(c *gin.Context) {
	authData := middlewares.GetInitDataFromContext(c)
	if authData == nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	req := new(CreateWishlistRequest)
	err := c.BindJSON(req)
	if err != nil {
		return
	}
	wishlist, err := s.db.CreateWishlist(c, store.CreateWishlistParams{
		OwnerID:     authData.User.ID,
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

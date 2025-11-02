package service

import (
	"backend/internal/middlewares"
	"backend/internal/store"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// UpdateWishlist godoc
// @Summary creates wishlist
// @Tags User
// @Router /api/user/wishlist [patch]
// @Security ApiKeyAuth
// @Accept	json
// @Param	wishlist_id	query	int						true	"Wishlist ID"
// @Param	wishlist	body	CreateWishlistRequest	true	"request body"
// @Produce	json
// @Success 204
func (s *Service) UpdateWishlist(c *gin.Context) {
	authData := middlewares.GetInitDataFromContext(c)
	if authData == nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	wishlistIDRaw := c.Param("wishlist_id")
	wishlistID, err := strconv.ParseInt(wishlistIDRaw, 10, 64)
	if err != nil {
		c.Error(fmt.Errorf("invalid wishlist_id: %w", err))
		c.AbortWithStatus(http.StatusBadRequest)
	}

	req := new(CreateWishlistRequest)
	if err := c.BindJSON(req); err != nil {
		return
	}

	count, err := s.db.UpdateWishlist(c, store.UpdateWishlistParams{
		ID:          wishlistID,
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
	if count == 0 {
		c.Status(http.StatusNotFound)
		return
	}
	c.Status(http.StatusNoContent)
}

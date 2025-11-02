package service

import (
	"backend/internal/middlewares"
	"backend/internal/store"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// UpdateUserWishlistItem godoc
// @Summary	updates a wishlist item
// @Tags	User
// @Router	/api/user/wishlist/item [patch]
// @Security	ApiKeyAuth
// @Accept	json
// @Param	item_id	query	int							true	"Item ID"
// @Param	item	body	CreateWishlistItemRequest	true	"Item"true "request body"
// @Success	204
func (s *Service) UpdateUserWishlistItem(c *gin.Context) {
	authData := middlewares.GetInitDataFromContext(c)
	if authData == nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// Get item ID from URL parameter
	itemIDStr := c.Query("item_id")
	itemID, err := strconv.ParseInt(itemIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item ID"})
		return
	}

	// Parse request body
	req := new(CreateWishlistItemRequest)
	err = c.BindJSON(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Convert links to JSON bytes
	linksJSON, err := json.Marshal(req.Links)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid links format"})
		return
	}

	// Update the wishlist item
	count, err := s.db.UpdateWishlistItem(c, store.UpdateWishlistItemParams{
		ID:          itemID,
		OwnerID:     authData.User.ID,
		Title:       req.Title,
		Description: req.Description,
		Price:       req.Price,
		Links:       linksJSON,
		Reservable:  req.Reservable,
	})
	if err != nil {
		c.Error(err)
		c.Status(http.StatusInternalServerError)
		return
	}
	if count == 0 {
		c.Status(http.StatusNotFound)
	}

	c.Status(http.StatusNoContent)
}

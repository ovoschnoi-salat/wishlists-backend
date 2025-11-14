package service

import (
	"backend/internal/middlewares"
	"backend/internal/store"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CreateWishlistItemRequest struct {
	WishlistID  int64              `json:"wishlist_id"`
	Title       string             `json:"title"`
	Description string             `json:"description"`
	Price       string             `json:"price"`
	Links       []WishlistItemLink `json:"links"`
	Reservable  bool               `json:"reservable"`
}

type WishlistItemLink struct {
	Title string `json:"title"`
	URL   string `json:"url"`
}

type WishlistItem struct {
	ID          int64              `json:"id"`
	WishlistID  int64              `json:"wishlist_id"`
	Title       string             `json:"title"`
	Description string             `json:"description"`
	Price       string             `json:"price"`
	Links       []WishlistItemLink `json:"links"`
	Reservable  bool               `json:"reservable"`
}

// CreateUserWishlistItem godoc
// @Summary creates a wishlist item
// @Tags User
// @Router /api/user/wishlist/item [post]
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param item body CreateWishlistItemRequest true "request body"
// @Success 200 {object} WishlistItem
func (s *Service) CreateUserWishlistItem(c *gin.Context) {
	authData := middlewares.GetInitDataFromContext(c)
	if authData == nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	req := new(CreateWishlistItemRequest)
	if err := c.BindJSON(req); err != nil {
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("invalid request body: %w", err))
		return
	}

	linksJSON, err := json.Marshal(req.Links)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("invalid links format: %w", err))
		return
	}

	item, err := s.db.CreateWishlistItem(c, store.CreateWishlistItemParams{
		WishlistID:  req.WishlistID,
		OwnerID:     authData.User.ID,
		Title:       req.Title,
		Description: req.Description,
		Price:       req.Price,
		Links:       linksJSON,
		Reservable:  req.Reservable,
	})
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("error creating item: %w", err))
		return
	}

	c.JSON(http.StatusOK, mapStoreWishlistItemToWishlistItem(item))
}

func mapStoreWishlistItemToWishlistItem(item store.WishlistItem) WishlistItem {
	// Parse links from JSON bytes
	var links []WishlistItemLink
	if len(item.Links) > 2 {
		json.Unmarshal(item.Links, &links)
	}

	return WishlistItem{
		ID:          item.ID,
		WishlistID:  item.WishlistID,
		Title:       item.Title,
		Description: item.Description,
		Price:       item.Price,
		Links:       links,
		Reservable:  item.Reservable,
	}
}

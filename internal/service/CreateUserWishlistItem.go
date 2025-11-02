package service

import (
	"backend/internal/middlewares"
	"backend/internal/store"
	"encoding/json"
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
	Reserved    bool               `json:"reserved"`
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

	// Parse request body
	req := new(CreateWishlistItemRequest)
	if err := c.BindJSON(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Convert links to JSON bytes
	linksJSON, err := json.Marshal(req.Links)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid links format"})
		return
	}

	// Create the wishlist item
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
		c.Error(err)
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, mapStoreWishlistItemToWishlistItem(item))
}

func mapStoreWishlistItemToWishlistItem(item store.WishlistItem) WishlistItem {
	// Parse links from JSON bytes
	var links []WishlistItemLink
	if len(item.Links) > 0 {
		json.Unmarshal(item.Links, &links)
	}

	// Handle reserved_by field
	var reservedBy *int64
	if item.ReservedBy.Valid {
		reservedBy = &item.ReservedBy.Int64
	}

	return WishlistItem{
		ID:          item.ID,
		WishlistID:  item.WishlistID,
		Title:       item.Title,
		Description: item.Description,
		Price:       item.Price,
		Links:       links,
		Reservable:  item.Reservable,
		Reserved:    reservedBy != nil,
	}
}

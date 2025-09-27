package service

import (
	"backend/internal/middlewares"
	"backend/internal/store"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CreateWishlistItemRequest struct {
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
	ReservedBy  *int64             `json:"reserved_by"`
}

// CreateUserWishlistItem godoc
// @Summary creates a wishlist item
// @Tags wishlist
// @Param wishlist_id query int true "Wishlist ID"
// @Param item body CreateWishlistItemRequest true "request body"
// @Accept json
// @Produce json
// @Success 200 {object} WishlistItem
// @Router /user/wishlist/item [post]
// @Security ApiKeyAuth
func (s *Service) CreateUserWishlistItem(c *gin.Context) {
	authData := middlewares.GetInitDataFromContext(c)
	if authData == nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// Get wishlist ID from URL parameter
	wishlistIDStr := c.Query("wishlist_id")
	wishlistID, err := strconv.ParseInt(wishlistIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid wishlist ID"})
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

	// Create the wishlist item
	item, err := s.db.CreateWishlistItem(c, store.CreateWishlistItemParams{
		WishlistID:  wishlistID,
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
		ReservedBy:  reservedBy,
	}
}

package service

import (
	"backend/internal/middlewares"
	"backend/internal/store"
	"backend/internal/subcodeErrors"
	"backend/internal/subcodeErrors/codes"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CreateWishlistRequest struct {
	Title           string  `json:"title"`
	Description     string  `json:"description"`
	IsPrivate       bool    `json:"is_private"`
	UsersWithAccess []int64 `json:"users_with_access"`
}

// CreateWishlist godoc
// @Summary creates wishlist
// @Tags User
// @Router /api/user/wishlist [post]
// @Security ApiKeyAuth
// @Accept	json
// @Produce	json
// @Param	wishlist	body	CreateWishlistRequest	true	"request body"
// @Failure 400 {object} subcodeErrors.Response
// @Failure 401 {object} subcodeErrors.Response
// @Failure 500 {object} subcodeErrors.Response
// @Success 200 {object} Wishlist
func (s *Service) CreateWishlist(c *gin.Context) {
	authData, authorized := middlewares.GetInitDataFromContext(c)
	if !authorized {
		subcodeErrors.SendResponse(c, http.StatusInternalServerError, codes.InternalErrCode, noInitDataErr)
		return
	}

	req := new(CreateWishlistRequest)
	err := c.ShouldBindJSON(req)
	if err != nil {
		subcodeErrors.SendResponse(c, http.StatusBadRequest, codes.InvalidRequestParametersErrCode, err)
		return
	}

	wishlist, err := s.db.CreateWishlist(c, store.CreateWishlistParams{
		OwnerID:     authData.User.ID,
		Title:       req.Title,
		Description: req.Description,
		IsPrivate:   req.IsPrivate,
		ShareUuid:   uuid.New().String(),
	})
	if err != nil {
		subcodeErrors.SendResponse(c, http.StatusInternalServerError, codes.InternalErrCode, fmt.Errorf("error creating wishlist: %w", err))
		return
	}

	if req.IsPrivate {
		for _, UserID := range req.UsersWithAccess {
			count, err := s.db.InsertWishlistAccessItem(c, store.InsertWishlistAccessItemParams{
				ListID:  wishlist.ID,
				OwnerID: authData.User.ID,
				UserID:  UserID,
			})
			if err != nil {
				_ = c.Error(fmt.Errorf("error inserting wishlist access item: %w", err))
			} else if count == 0 {
				_ = c.Error(fmt.Errorf("error inserting wishlist access item: not inserted id %d", UserID))
			}
		}
	}

	c.JSON(http.StatusOK, mapStoreWishlistToWishlist(wishlist))
}

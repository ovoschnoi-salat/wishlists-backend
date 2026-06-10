package service

import (
	"backend/internal/middlewares"
	"backend/internal/store"
	"backend/internal/subcodeErrors"
	"backend/internal/subcodeErrors/codes"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// DeleteWishSplitRequest godoc
// @Summary		creates wish split request
// @Tags		Friend's Wish
// @Router		/api/user/friend/wishlist/wish/split-request [delete]
// @Security	ApiKeyAuth
// @Param		wish_id	query	int	true	"Wish ID"
// @Failure 400 {object} subcodeErrors.Response
// @Failure 401 {object} subcodeErrors.Response
// @Failure 500 {object} subcodeErrors.Response
// @Success	204
func (s *Service) DeleteWishSplitRequest(c *gin.Context) {
	authData, authorized := middlewares.GetInitDataFromContext(c)
	if !authorized {
		subcodeErrors.SendResponse(c, http.StatusInternalServerError, codes.InternalErrCode, noInitDataErr)
		return
	}

	wishIDStr := c.Query("wish_id")
	wishID, err := strconv.ParseInt(wishIDStr, 10, 64)
	if err != nil {
		subcodeErrors.SendResponse(c, http.StatusBadRequest, codes.InvalidRequestParametersErrCode, fmt.Errorf("invalid wish ID: %w", err))
		return
	}

	count, err := s.db.DeleteWishSplitRequest(c, store.DeleteWishSplitRequestParams{
		WishID: wishID,
		UserID: authData.User.ID,
	})
	if err != nil {
		subcodeErrors.SendResponse(c, http.StatusInternalServerError, codes.InternalErrCode, fmt.Errorf("error cancelling split request: %w", err))
		return
	}
	if count == 0 {
		subcodeErrors.SendResponse(c, http.StatusBadRequest, codes.InvalidRequestErrCode, errors.New("error cancelling split request: no rows affected"))
		return
	}

	c.Status(http.StatusNoContent)
}

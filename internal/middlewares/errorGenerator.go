package middlewares

import (
	"backend/internal/subcodeErrors"
	"backend/internal/subcodeErrors/codes"
	"errors"
	"math/rand"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ErrorGenerator(c *gin.Context) {
	if rand.Int()&1 != 0 {
		subcodeErrors.SendResponse(c, http.StatusInternalServerError, codes.TestErrCode, errors.New("test error"))
	}
}

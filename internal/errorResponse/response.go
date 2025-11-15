package errorResponse

import (
	"backend/internal/errorResponse/codes"
	uuidMiddleware "backend/internal/middlewares/uuid"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Response struct {
	Subcode     codes.ErrorCode
	RequestUUID uuid.UUID
}

func Send(c *gin.Context, httpCode int, errCode codes.ErrorCode, err error) {
	requestUUID := uuidMiddleware.GetUUIDFromContext(c)
	resp := Response{
		Subcode:     errCode,
		RequestUUID: requestUUID,
	}
	codes.SetErrorCodeToContext(c, errCode)
	if err != nil {
		_ = c.Error(err)
	}
	c.JSON(httpCode, resp)
}

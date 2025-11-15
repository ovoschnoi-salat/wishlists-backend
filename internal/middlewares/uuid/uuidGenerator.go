package uuid

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const uuidCtxKey = "request-uuid"

func Generator(c *gin.Context) {
	requestUuid := uuid.New()
	c.Set(uuidCtxKey, requestUuid)
}

func GetUUIDFromContext(c *gin.Context) uuid.UUID {
	if value, ok := c.Get(uuidCtxKey); ok {
		if requestUUID, ok := value.(uuid.UUID); ok {
			return requestUUID
		}
	}
	return uuid.Nil
}

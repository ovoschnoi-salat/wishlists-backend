package middlewares

import (
	uuid2 "backend/internal/middlewares/uuid"
	"backend/internal/subcodeErrors/codes"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func Logger(c *gin.Context) {
	start := time.Now()

	c.Next()

	latency := time.Now().Sub(start)
	if latency > time.Minute {
		latency = latency.Truncate(time.Second)
	}

	uuid := uuid2.GetUUIDFromContext(c)

	errCode, foundErrCode := codes.GetErrorCodeFromContext(c)

	var event *zerolog.Event
	if c.Writer.Status() >= 500 || len(c.Errors) != 0 {
		event = log.Error()
	} else {
		event = log.Info()
	}
	if len(c.Errors) != 0 {
		event = event.Strs("errors", c.Errors.Errors())
	}
	event = event.
		Str("client-ip", c.ClientIP()).
		Int("status", c.Writer.Status()).
		Str("latency", latency.String()).
		Str("method", c.Request.Method).
		Str("pattern", c.FullPath()).
		Str("path", c.Request.URL.Path).
		Str("uuid", uuid.String())
	if foundErrCode {
		event = event.Int32("error-code", int32(errCode))
	}
	authData, authorized := GetInitDataFromContext(c)
	if authorized {
		event = event.Any("username", authData.User.Username)
		event = event.Any("user_id", authData.User.ID)
	}
	event.Send()
}

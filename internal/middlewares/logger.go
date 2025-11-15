package middlewares

import (
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
	var event *zerolog.Event
	if c.Writer.Status() >= 500 {
		event = log.Error()
	} else if len(c.Errors) != 0 {
		event = log.Error()
		event = event.Strs("errors", c.Errors.Errors())
		c.JSON(0, c.Errors.JSON())
	} else {
		event = log.Info()
	}
	event = event.
		Str("client-ip", c.ClientIP()).
		Int("status", c.Writer.Status()).
		Str("latency", latency.String()).
		Str("method", c.Request.Method).
		Str("pattern", c.FullPath()).
		Str("path", c.Request.URL.Path)
	authData := GetInitDataFromContext(c)
	if authData != nil {
		event = event.Any("username", authData.User.Username)
		event = event.Any("user_id", authData.User.ID)
	}
	event.Send()
}

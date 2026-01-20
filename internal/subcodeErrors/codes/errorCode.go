package codes

import (
	"github.com/gin-gonic/gin"
)

type ErrorCode int32

const (
	UnknownErrCode                  ErrorCode = 0
	UnauthorizedErrCode             ErrorCode = 1001
	InternalErrCode                 ErrorCode = 1002
	NotFoundErrCode                 ErrorCode = 1003
	NoAccessErrCode                 ErrorCode = 1004
	InvalidRequestErrCode           ErrorCode = 1005
	InvalidRequestParametersErrCode ErrorCode = 1006

	FriendNotFoundErrCode            ErrorCode = 2001
	CantSendRequestToYourselfErrCode ErrorCode = 2002

	WishNotFoundErrCode ErrorCode = 3001

	TestErrCode ErrorCode = 6666
)

const errCodeCtxKey = "ErrorCode"

func SetErrorCodeToContext(c *gin.Context, errCode ErrorCode) {
	c.Set(errCodeCtxKey, errCode)
}

func GetErrorCodeFromContext(c *gin.Context) (ErrorCode, bool) {
	if value, ok := c.Get(errCodeCtxKey); ok {
		if code, ok := value.(ErrorCode); ok {
			return code, true
		}
	}
	return UnknownErrCode, false
}

func abs[T ~int | ~int8 | ~int16 | ~int32 | ~int64](v T) T {
	if v < 0 {
		return -v
	}
	return v
}

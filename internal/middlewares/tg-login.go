package middlewares

import (
	"backend/internal/store"
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	initdata "github.com/telegram-mini-apps/init-data-golang"
)

const initDataCtxKey = "init_data"

func NewTgAuthMiddleware(secretToken string, db *store.Queries) gin.HandlerFunc {
	// Define how long since init data generation date init data is valid.
	//expIn := 10 * time.Minute

	return func(c *gin.Context) {
		initData := c.GetHeader("authorization")
		if len(initData) == 0 || !strings.HasPrefix(initData, "tma") || len(strings.TrimPrefix(initData, "tma ")) == 0 {
			c.Error(errors.New("no init data found"))
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		initData = strings.TrimPrefix(initData, "tma ")
		//err := initdata.Validate(strings.TrimPrefix(initData[0], "tma "), secretToken, expIn)
		//if err != nil {
		//	return nil, err
		//}
		data, err := initdata.Parse(initData)
		if err != nil {
			c.Error(err)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		_, err = db.GetUser(c, data.User.ID)
		if errors.Is(err, pgx.ErrNoRows) {
			count, err := db.CreateUser(c, store.CreateUserParams{
				ID:       data.User.ID,
				Username: data.User.Username,
				Name:     data.User.FirstName + " " + data.User.LastName,
				PhotoUrl: data.User.PhotoURL,
			})
		}
		if err != nil {
			c.Error(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		c.Set(initDataCtxKey, &data)
	}
}

func GetInitDataFromContext(ctx context.Context) *initdata.InitData {
	if value, ok := ctx.Value(initDataCtxKey).(*initdata.InitData); ok {
		return value
	}
	return nil
}

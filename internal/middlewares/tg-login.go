package middlewares

import (
	"backend/internal/config"
	"backend/internal/store"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	initdata "github.com/telegram-mini-apps/init-data-golang"
)

const userDataHeaderKey = "x-user-data"
const userDataCtxKey = "user_data"
const initDataCtxKey = "init_data"

func NewTgAuthMiddleware(secretToken string, db *store.Queries, stage config.Stage) gin.HandlerFunc {
	// Define how long since init data generation date init data is valid.
	expIn := 10 * time.Minute

	return func(c *gin.Context) {
		initData := c.GetHeader("authorization")
		if len(initData) == 0 || !strings.HasPrefix(initData, "tma") || len(strings.TrimPrefix(initData, "tma ")) == 0 {
			c.AbortWithError(http.StatusUnauthorized, errors.New("no init data found"))
			return
		}
		initData = strings.TrimPrefix(initData, "tma ")

		if stage == config.PROD {
			err := initdata.Validate(initData, secretToken, expIn)
			if err != nil {
				c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("error validating init data: %w", err))
				return
			}
		}

		data, err := initdata.Parse(initData)
		if err != nil {
			c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("error parsing init data: %w", err))
			return
		}
		user, err := db.GetUser(c, data.User.ID)

		if err != nil {
			if !errors.Is(err, pgx.ErrNoRows) {
				c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("error getting user: %w", err))
				return
			}
			req := store.CreateUserParams{
				ID:       data.User.ID,
				Username: data.User.Username,
				PhotoUrl: data.User.PhotoURL,
			}
			if data.User.FirstName != "" {
				var name string
				if data.User.LastName != "" {
					name = data.User.FirstName + " " + data.User.LastName
				} else {
					name = data.User.FirstName
				}
				req.Name = pgtype.Text{
					String: name,
					Valid:  true,
				}
				req.DisplayedName = name
			} else {
				req.DisplayedName = "@" + req.Username
			}
			user, err = db.CreateUser(c, req)
			if err != nil {
				c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("error creating user: %w", err))
				return
			}
		}
		c.Set(userDataCtxKey, &user)
		c.Set(initDataCtxKey, &data)
		userData, err := json.Marshal(user)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("error marshaling user: %w", err))
			return
		}
		c.Header(userDataHeaderKey, string(userData))
	}
}

func GetUserDataFromContext(ctx context.Context) *store.User {
	if value, ok := ctx.Value(userDataCtxKey).(*store.User); ok {
		return value
	}
	return nil
}

func GetInitDataFromContext(ctx context.Context) *initdata.InitData {
	if value, ok := ctx.Value(initDataCtxKey).(*initdata.InitData); ok {
		return value
	}
	return nil
}

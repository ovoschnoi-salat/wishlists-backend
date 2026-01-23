package middlewares

import (
	"backend/internal/config"
	"backend/internal/store"
	"backend/internal/subcodeErrors"
	"backend/internal/subcodeErrors/codes"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	initdata "github.com/telegram-mini-apps/init-data-golang"
)

const userDataCtxKey = "user_data"
const initDataCtxKey = "init_data"

func NewTgAuthMiddleware(secretToken string, db *store.Queries, stage config.Stage) gin.HandlerFunc {
	// Define how long since init data generation date init data is valid.
	expIn := 365 * 24 * time.Hour

	return func(c *gin.Context) {
		initData := c.GetHeader("authorization")
		if len(initData) == 0 || !strings.HasPrefix(initData, "tma") || len(strings.TrimPrefix(initData, "tma ")) == 0 {
			subcodeErrors.SendResponse(c, http.StatusUnauthorized, codes.UnauthorizedErrCode, errors.New("no init data found"))
			return
		}
		initData = strings.TrimPrefix(initData, "tma ")

		if stage == config.PROD {
			err := initdata.Validate(initData, secretToken, expIn)
			if err != nil {
				subcodeErrors.SendResponse(c, http.StatusUnauthorized, codes.UnauthorizedErrCode, fmt.Errorf("error validating init data: %w", err))
				return
			}
		}

		data, err := initdata.Parse(initData)
		if err != nil {
			subcodeErrors.SendResponse(c, http.StatusUnauthorized, codes.UnauthorizedErrCode, fmt.Errorf("error parsing init data: %w", err))
			return
		}
		user, err := db.GetUser(c, data.User.ID)

		if err != nil {
			if !errors.Is(err, pgx.ErrNoRows) {
				subcodeErrors.SendResponse(c, http.StatusInternalServerError, codes.InternalErrCode, fmt.Errorf("error getting user: %w", err))
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
				req.DisplayedName = name
			}
			user, err = db.CreateUser(c, req)
			if err != nil {
				subcodeErrors.SendResponse(c, http.StatusInternalServerError, codes.InternalErrCode, fmt.Errorf("error creating user: %w", err))
				return
			}
		}
		c.Set(userDataCtxKey, user)
		c.Set(initDataCtxKey, data)
	}
}

func GetUserDataFromContext(c *gin.Context) (store.User, bool) {
	if value, ok := c.Get(userDataCtxKey); ok {
		if user, ok := value.(store.User); ok {
			return user, true
		}
	}
	return store.User{}, false
}

func GetInitDataFromContext(c *gin.Context) (initdata.InitData, bool) {
	if value, ok := c.Get(initDataCtxKey); ok {
		if initData, ok := value.(initdata.InitData); ok {
			return initData, true
		}
	}
	return initdata.InitData{}, false
}

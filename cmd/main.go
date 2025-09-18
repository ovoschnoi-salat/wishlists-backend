package main

import (
	"backend/internal/service"
	"backend/internal/store"
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"

	_ "backend/docs"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

//	@title			Swagger Example API
//	@version		1.0
//	@description	This is a sample server celler server.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@host		localhost:8080
//	@BasePath	/api

//	@securityDefinitions.basic	BasicAuth

//	@securityDefinitions.apikey	ApiKeyAuth
//	@in							header
//	@name						Authorization
//	@description				Description for what is this security definition being used

//	@securitydefinitions.oauth2.application	OAuth2Application
//	@tokenUrl								https://example.com/oauth/token
//	@scope.write							Grants write access
//	@scope.admin							Grants read and write access to administrative information

//	@securitydefinitions.oauth2.implicit	OAuth2Implicit
//	@authorizationUrl						https://example.com/oauth/authorize
//	@scope.write							Grants write access
//	@scope.admin							Grants read and write access to administrative information

//	@securitydefinitions.oauth2.password	OAuth2Password
//	@tokenUrl								https://example.com/oauth/token
//	@scope.read								Grants read access
//	@scope.write							Grants write access
//	@scope.admin							Grants read and write access to administrative information

//	@securitydefinitions.oauth2.accessCode	OAuth2AccessCode
//	@tokenUrl								https://example.com/oauth/token
//	@authorizationUrl						https://example.com/oauth/authorize
//	@scope.admin							Grants read and write access to administrative information

func main() {
	ctx := context.Background()
	dsn := os.Getenv("PG_DSN")

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		log.Fatal(err)
	}

	if err := upMigrations(pool); err != nil {
		log.Fatal(err)
	}

	storeObj := store.New(pool)

	serviceObj := service.NewService(storeObj)

	port := os.Getenv("PORT")

	router := gin.Default()

	router.Use(cors.Default()) // All origins allowed by default

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	g := router.Group("/api")
	serviceObj.RegisterHandlers(g)

	log.Printf("server listening at :%s", port)
	err = http.ListenAndServe(":"+port, router)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
}

func upMigrations(pool *pgxpool.Pool) error {
	db := stdlib.OpenDBFromPool(pool)

	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Println("error closing db: ", err)
		}
	}(db)

	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	if err := goose.Up(db, "migrations"); err != nil {
		return err
	}

	return nil
}

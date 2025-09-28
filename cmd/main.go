package main

import (
	"backend/internal/middlewares"
	"backend/internal/service"
	"backend/internal/store"
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	_ "backend/docs"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

//	@title			Wishlists backend API
//	@version		1.0
//	@description	Backend service for wishlists app

//	@host		localhost:8080

//	@securityDefinitions.apikey	ApiKeyAuth
//	@in			header
//	@name		Authorization

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

	router.Use(cors.New(cors.Config{
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
		AllowAllOrigins:  true,
	})) // All origins allowed by default

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	serviceObj.RegisterHandlers(router.Group("", middlewares.NewTgAuthMiddleware("", storeObj)))

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

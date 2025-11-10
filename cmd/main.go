package main

import (
	"backend/internal/config"
	"backend/internal/middlewares"
	"backend/internal/service"
	"backend/internal/store"
	"backend/pkg/graceful"
	"backend/pkg/http"
	"context"
	"database/sql"
	"os"
	"time"

	_ "backend/docs"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/rs/zerolog/log"
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
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("error loading config")
	}

	ctx := context.Background()
	dsn := os.Getenv("PG_DSN")

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		log.Fatal().Err(err).Msg("error connecting to database")
	}

	if err := upMigrations(pool); err != nil {
		log.Fatal().Err(err).Msg("error upgrading migrations")
	}

	storeObj := store.New(pool)

	serviceObj := service.NewService(storeObj)

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(middlewares.Logger, gin.Recovery())

	router.Use(cors.New(cors.Config{
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
		AllowAllOrigins:  true,
	})) // All origins allowed by default

	if cfg.Stage == config.DEV {
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	}

	serviceObj.RegisterHandlers(router.Group("", middlewares.NewTgAuthMiddleware("", storeObj, cfg.Stage)))

	httpServer := http.NewServer(cfg.HttpServer, router)

	runner := graceful.DefaultConfig().GetRunner()
	err = runner.Run(httpServer)
	if err != nil {
		log.Fatal().Err(err).Msg("error running app")
	}
}

func upMigrations(pool *pgxpool.Pool) error {
	db := stdlib.OpenDBFromPool(pool)

	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Error().Err(err).Msg("error closing db")
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

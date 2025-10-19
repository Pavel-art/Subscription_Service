package main

import (
	"SubscriptionService/configs"
	"SubscriptionService/internal/api"
	"SubscriptionService/internal/application/services"
	"SubscriptionService/internal/persistence"
	"SubscriptionService/pkg/db"
	"SubscriptionService/pkg/logger"
	"context"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	zerologgin "github.com/go-mods/zerolog-gin"
)

func main() {
	ctx := context.Background()

	// --- init configs ---
	configs.Init()
	dbConfig := configs.NewDataBaseConfig()
	logConfig := configs.NewLogConfig()

	// --- init logger ---
	customLogger := logger.NewLogger(logConfig)

	// --- init gin app ---
	app := gin.Default()
	app.Use(zerologgin.LoggerWithOptions(&zerologgin.Options{
		Name:   "server",
		Logger: customLogger,
	}))
	app.Use(gin.Recovery())

	// --- init database pool ---
	pool, err := db.NewPGXPool(ctx, dbConfig.Url, customLogger)
	if err != nil {
		log.Fatalf("failed to create db pool: %v", err)
	}
	defer pool.Close()

	// --- init repository ---
	subRepo := persistence.NewSubRepository(pool)

	// --- init service ---
	subService := services.NewSubService(subRepo, customLogger)

	// --- init handlers ---
	api.NewHandler(app, subService, customLogger)

	// --- run server ---
	port := ":8080"
	customLogger.Info().Msgf("Starting server on %s", port)
	if err := app.Run(port); err != nil {
		customLogger.Fatal().Err(err).Msg("failed to run server")
		os.Exit(1)
	}
}

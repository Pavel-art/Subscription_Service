package main

import (
	"SubscriptionService/configs"
	"SubscriptionService/internal/api"
	"SubscriptionService/internal/application/services"
	"SubscriptionService/internal/persistence"
	"SubscriptionService/pkg/db"
	"SubscriptionService/pkg/logger"
	"SubscriptionService/pkg/migrate"
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	zerologgin "github.com/go-mods/zerolog-gin"
)

func main() {
	ctx := context.Background()

	// --- init configs ---
	configs.Init()
	dbConfig := configs.NewDataBaseConfig()
	logConfig := configs.NewLogConfig()
	serverConfig := configs.NewServerConfig()

	// --- init logger ---
	customLogger := logger.NewLogger(logConfig)

	// --- init gin app ---
	gin.SetMode(gin.ReleaseMode)
	app := gin.New()
	app.Use(zerologgin.LoggerWithOptions(&zerologgin.Options{
		Name:   "server",
		Logger: customLogger,
	}))
	app.Use(gin.Recovery())

	// --- run migrations ---
	if err := migrate.Up(dbConfig.Url, customLogger); err != nil {
		customLogger.Warn().Err(err).Msg("migrations failed (continuing anyway, tables may already exist)")
	}

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
	api.RegisterSwagger(app)

	// --- run server ---
	addr := serverConfig.Port
	if addr == "" {
		addr = "8081"
	}
	if !strings.HasPrefix(addr, ":") {
		addr = ":" + addr
	}

	server := &http.Server{
		Addr:    addr,
		Handler: app,
	}

	customLogger.Info().Msgf("Starting server on %s", addr)

	// Запускаем сервер в отдельной горутине
	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			customLogger.Fatal().Err(err).Msg("failed to run server")
		}
	}()

	// Ожидаем сигналов завершения
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	customLogger.Info().Msg("Shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		customLogger.Error().Err(err).Msg("Server forced to shutdown")
	} else {
		customLogger.Info().Msg("Server exited gracefully")
	}
}

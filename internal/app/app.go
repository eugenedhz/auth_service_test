package app

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/eugenedhz/auth_service_test/internal/config"
	"github.com/eugenedhz/auth_service_test/internal/domain/models"
	"github.com/eugenedhz/auth_service_test/internal/lib/tokens"
	userrepo "github.com/eugenedhz/auth_service_test/internal/repository/mock"
	sessionrepo "github.com/eugenedhz/auth_service_test/internal/repository/postgres"
	"github.com/eugenedhz/auth_service_test/internal/service/auth"
	authhttp "github.com/eugenedhz/auth_service_test/internal/transport/http"
	"github.com/eugenedhz/auth_service_test/pkg/database"
	email "github.com/eugenedhz/auth_service_test/pkg/email/mock"
)

const (
	envDev  = "dev"
	envProd = "prod"
)

type App struct {
	httpServer  *http.Server
	serverPort  string
	authService *auth.AuthService
}

func NewApp() *App {
	cfg, err := config.LoadConfigs()
	if err != nil {
		panic(err)
	}
	host := cfg.DbHost
	port := cfg.DbPort
	user := cfg.DbUser
	password := cfg.DbPassword
	dbname := cfg.DbName
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	db, err := database.NewPostgresConnection(dsn)
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&models.Session{})

	userRepo := userrepo.NewUserRepository()
	sessionRepo := sessionrepo.NewSessionRepository(db)
	tokenManager := tokens.NewTokenManager(cfg.SigningKey, cfg.AccessTokenTTL)
	emailSender := email.NewEmailSender()
	logger := setupLogger(cfg.LogLevel)

	authService := auth.NewAuthService(
		sessionRepo,
		userRepo,
		tokenManager,
		emailSender,
		logger,
	)

	return &App{
		authService: authService,
		serverPort:  cfg.ServerPort,
	}
}

func (a *App) Run() error {
	router := gin.Default()
	router.Use(
		gin.Recovery(),
		gin.Logger(),
	)

	authhttp.RegisterHTTPEndpoints(router, a.authService)

	a.httpServer = &http.Server{
		Addr:    ":" + a.serverPort,
		Handler: router,
	}

	go func() {
		if err := a.httpServer.ListenAndServe(); err != nil {
			panic(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Interrupt)

	<-quit

	ctx, shutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdown()

	return a.httpServer.Shutdown(ctx)
}

func setupLogger(logLevel string) *slog.Logger {
	var log *slog.Logger

	switch logLevel {
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}

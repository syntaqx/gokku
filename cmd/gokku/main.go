package main

import (
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/sessions"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/syntaqx/gokku/internal/config"
	"github.com/syntaqx/gokku/internal/handlers"
	"github.com/syntaqx/gokku/internal/middleware"
	"github.com/syntaqx/gokku/internal/render"
	"github.com/syntaqx/gokku/internal/repository"
)

func main() {
	cfg := config.LoadConfig()

	var logger *zap.Logger
	if cfg.Debug {
		logger, _ = zap.NewDevelopment()
	} else {
		logger, _ = zap.NewProduction()
	}
	defer logger.Sync()

	dsn := cfg.Database.AssembleDSN()
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Fatal("failed to connect database", zap.Error(err))
	}

	err = db.AutoMigrate(&repository.User{})
	if err != nil {
		logger.Fatal("failed to migrate database", zap.Error(err))
	}

	userRepository := repository.NewUsersRepository(db)
	userRepository.Seed() // Seed the initial user

	store := sessions.NewCookieStore([]byte("your-secret-key")) // Use a secure key here

	renderer := render.New(render.Options{
		Layout:        "_layouts/base",
		Directory:     "templates",
		Extensions:    []string{".html"},
		IsDevelopment: cfg.Debug,
	})

	healthHandler := handlers.NewHealthHandler(renderer)
	authHandler := handlers.NewAuthHandler(renderer, userRepository, store)
	pageHandler := handlers.NewPageHandler(renderer)
	settingsHandler := handlers.NewSettingsHandler(renderer)
	appsHandler := handlers.NewApplicationsHandler(renderer)
	actionsHandler := handlers.NewActionsHandler(renderer)

	r := chi.NewRouter()

	// Base middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Public routes
	healthHandler.RegisterRoutes(r)
	authHandler.RegisterRoutes(r)

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(middleware.Auth(userRepository, store))
		pageHandler.RegisterRoutes(r)
		settingsHandler.RegisterRoutes(r)
		appsHandler.RegisterRoutes(r)
		actionsHandler.RegisterRoutes(r)
	})

	// Create a route along /files that will serve contents from
	// the ./data/ folder.
	workDir, _ := os.Getwd()
	filesDir := http.Dir(filepath.Join(workDir, "static"))
	FileServer(r, "/", filesDir)

	srv := &http.Server{
		Addr:         net.JoinHostPort(cfg.Server.Host, cfg.Server.Port),
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	logger.Info("Starting server", zap.String("address", srv.Addr))
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Fatal("Server failed", zap.Error(err))
	}
}

// FileServer conveniently sets up a http.FileServer handler to serve
// static files from a http.FileSystem.
func FileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit any URL parameters.")
	}

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
		fs.ServeHTTP(w, r)
	})
}

package main

import (
	"context"
	"net/http"
	"os"

	"github.com/IvanTime-Kai/url-shortener/internal/cache"
	"github.com/IvanTime-Kai/url-shortener/internal/handler"
	"github.com/IvanTime-Kai/url-shortener/internal/logger"
	"github.com/IvanTime-Kai/url-shortener/internal/repository"
	"github.com/IvanTime-Kai/url-shortener/internal/service"
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"

	"github.com/IvanTime-Kai/url-shortener/internal/middleware"
)

func main() {
	godotenv.Load()

	// Logger
	log := logger.New()

	// PostgreSQL
	db, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Error("cannot connect to database:", err)
	}
	defer db.Close()

	// Redis
	opt, err := redis.ParseURL(os.Getenv("REDIS_URL"))
	if err != nil {
		log.Error("cannot parse redis url", err)
	}
	rdb := redis.NewClient(opt)
	defer rdb.Close()

	// Wire up
	linkRepo := repository.NewLinkRepository(db)
	clickRepo := repository.NewClickRepository(db)
	counter := service.NewClickCounter()
	linkCache := cache.NewLinkCache(rdb)
	rateLimiter := cache.NewRateLimit(rdb)
	svc := service.NewLinkService(linkRepo, clickRepo, counter, linkCache, os.Getenv("BASE_URL"))
	h := handler.NewLinkHandler(svc)

	r := chi.NewRouter()
	r.Use(middleware.Logger(log))
	r.Use(chiMiddleware.Recoverer)
	r.Use(middleware.RateLimit(rateLimiter))

	r.Route("/api/links", func(r chi.Router) {
		r.Post("/", h.Shorten)
		r.Get("/", h.List)
		r.Delete("/{code}", h.Delete)
		r.Get("/{code}/stats", h.Stats)
	})

	r.Get("/{code}", h.Redirect)

	port := os.Getenv("APP_PORT")
	log.Info("server starting", "port", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Error("server error", "error", err)
		os.Exit(1)
	}
}

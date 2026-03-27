package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/IvanTime-Kai/url-shortener/internal/handler"
	"github.com/IvanTime-Kai/url-shortener/internal/repository"
	"github.com/IvanTime-Kai/url-shortener/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	db, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("cannot connect to database:", err)
	}

	defer db.Close()

	repo := repository.NewLinkRepository(db)
	svc := service.NewLinkService(repo, os.Getenv("BASE_URL"))
	h := handler.NewLinkHandler(svc)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/{code}", h.Redirect)
	r.Route("/api/links", func(r chi.Router) {
		r.Post("/", h.Shorten)
		r.Get("/", h.List)
		r.Delete("/{code}", h.Delete)
	})

	port := os.Getenv("APP_PORT")
	fmt.Printf("Server running on :%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}

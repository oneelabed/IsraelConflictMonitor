package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/oneelabed/IsraelConflictMonitor/internal/config"
	"github.com/oneelabed/IsraelConflictMonitor/internal/database"
	. "github.com/oneelabed/IsraelConflictMonitor/internal/handlers"
	. "github.com/oneelabed/IsraelConflictMonitor/internal/scraper"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load()

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT is not found in env")
	}

	DB_url := os.Getenv("DB_URL")
	if DB_url == "" {
		log.Fatal("DB url is not found in env")
	}

	conn, err := sql.Open("postgres", DB_url)
	if err != nil {
		log.Fatal("Can't connect to DB:", err)
	}

	apiCfg := config.ApiConfig{
		DB: database.New(conn),
	}

	go StartScraping(apiCfg.DB, 23, time.Minute)

	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	v1router := chi.NewRouter()
	v1router.Get("/healthz", HandlerReadiness)
	v1router.Get("/err", HandlerError)
	v1router.Post("/users", func(w http.ResponseWriter, r *http.Request) {
		HandlerCreateUser(&apiCfg, w, r)
	})
	v1router.Post("/login", func(w http.ResponseWriter, r *http.Request) {
		HandlerLogin(&apiCfg, w, r)
	})
	v1router.Get("/users", MiddlewareAuth(&apiCfg, HandlerGetUserByAPI))
	v1router.Post("/feeds", func(w http.ResponseWriter, r *http.Request) {
		HandlerCreateFeed(&apiCfg, w, r)
	})
	v1router.Get("/feeds", MiddlewareAuth(&apiCfg, HandlerGetFeedsForUser))
	v1router.Post("/feed_follows", MiddlewareAuth(&apiCfg, HandlerCreateFeedFollow))
	v1router.Get("/feed_follows", MiddlewareAuth(&apiCfg, HandlerGetFeedFollows))
	v1router.Delete("/feed_follows", MiddlewareAuth(&apiCfg, HandlerDeleteFeedFollow))
	v1router.Get("/posts", MiddlewareAuth(&apiCfg, HandlerGetPostsForUser))
	v1router.Get("/posts/search", MiddlewareAuth(&apiCfg, HandlerSearchPosts))
	v1router.Get("/posts/diverse", func(w http.ResponseWriter, r *http.Request) {
		HandlerGetDiversePosts(&apiCfg, w, r)
	})
	v1router.Get("/posts/check-new", MiddlewareAuth(&apiCfg, HandlerCheckNewPosts))

	router.Mount("/v1", v1router)

	srv := &http.Server{
		Handler: router,
		Addr:    ":" + port,
	}

	fmt.Printf("Server starting on port %v\n", port)
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

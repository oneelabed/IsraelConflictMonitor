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
	"github.com/oneelabed/RSSAggregator/internal/database"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB *database.Queries // connection to DB
}

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

	apiCfg := apiConfig{
		DB: database.New(conn),
	}

	go startScraping(apiCfg.DB, 10, time.Minute)

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
	v1router.Get("/healthz", handlerReadiness)
	v1router.Get("/err", handlerError)
	v1router.Post("/users", apiCfg.handlerCreateUser)
	v1router.Get("/users", apiCfg.middlewareAuth(apiCfg.handlerGetUserByAPI))
	v1router.Post("/feeds", apiCfg.middlewareAuth(apiCfg.handlerCreateFeed))
	v1router.Get("/feeds", apiCfg.handlerGetFeeds)
	v1router.Post("/feed_follows", apiCfg.middlewareAuth(apiCfg.handlerCreateFeedFollow))
	v1router.Get("/feed_follows", apiCfg.middlewareAuth(apiCfg.handlerGetFeedFollows))
	v1router.Delete("/feed_follows/{feedFollowId}", apiCfg.middlewareAuth(apiCfg.handlerDeleteFeedFollow))
	v1router.Get("/feed_follows", apiCfg.middlewareAuth(apiCfg.handlerGetFeedFollows))

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

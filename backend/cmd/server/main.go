package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/0xstxrless/punkt/backend/internal/db"
	"github.com/0xstxrless/punkt/backend/router"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	ctx := context.Background()

	pool, err := db.Connect(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(err)
	}

	logger := router.NewLogger()

	app := &router.App{
		Queries: db.New(pool),
		Logger:  logger,
	}

	srv := &http.Server{
		Addr:    ":" + os.Getenv("PORT"),
		Handler: router.NewRouter(app),
	}

	log.Println("Server listening on", srv.Addr)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
